.PHONY: help build test lint deploy verify monitor cleanup integration-test

help:
	@echo "build            - Build all images"
	@echo "test             - Run all tests"
	@echo "lint             - Format and lint code"
	@echo "deploy           - Deploy system"
	@echo "verify           - Check deployment"
	@echo "monitor          - Access monitoring dashboard"
	@echo "integration-test - Build & Run integration tests"
	@echo "cleanup          - Remove deployment"

build:
	docker build -t crypto-trackers/data-ingestion:latest ./services/data-ingestion/
	docker build -t crypto-trackers/ma-signal-detector:latest ./services/ma-signal-detector/
	docker build -t crypto-trackers/volume-spike-detector:latest ./services/volume-spike-detector/
	docker build -t crypto-trackers/alert-service:latest ./services/alert-service/

test:
	@echo "Testing data-ingestion..."
	@cd ./services/data-ingestion && (test -d venv || (python -m venv venv && ./venv/bin/pip install -r requirements.txt -r requirements-dev.txt)) && ./venv/bin/python -m pytest; echo $$? > /tmp/test_data_ingestion_exit
	@echo "Testing ma-signal-detector..."
	@cd ./services/ma-signal-detector && go test ./... -v; echo $$? > /tmp/test_ma_signal_exit
	@echo "Testing volume-spike-detector..."
	@cd ./services/volume-spike-detector && go test ./... -v; echo $$? > /tmp/test_volume_spike_exit
	@echo "Testing alert-service..."
	@cd ./services/alert-service && go test ./... -v; echo $$? > /tmp/test_alert_service_exit
	@data_exit=$$(cat /tmp/test_data_ingestion_exit); ma_exit=$$(cat /tmp/test_ma_signal_exit); volume_exit=$$(cat /tmp/test_volume_spike_exit); alert_exit=$$(cat /tmp/test_alert_service_exit); \
	total_exit=$$(($$data_exit + $$ma_exit + $$volume_exit + $$alert_exit)); \
	rm -f /tmp/test_*_exit; \
	if [ $$total_exit -eq 0 ]; then echo "All tests completed successfully"; else echo "Tests failed in one or more services"; exit 1; fi

lint:
	@echo "Linting data-ingestion..."
	@cd ./services/data-ingestion && (test -d venv || (python -m venv venv && ./venv/bin/pip install -r requirements.txt -r requirements-dev.txt)) && ./venv/bin/python -m black . && ./venv/bin/python -m ruff check . --fix && ./venv/bin/python -m mypy .; echo $$? > /tmp/lint_data_ingestion_exit
	@echo "Linting ma-signal-detector..."
	@cd ./services/ma-signal-detector && go fmt ./... && go vet ./...; echo $$? > /tmp/lint_ma_signal_exit
	@echo "Linting volume-spike-detector..."
	@cd ./services/volume-spike-detector && go fmt ./... && go vet ./...; echo $$? > /tmp/lint_volume_spike_exit
	@echo "Linting alert-service..."
	@cd ./services/alert-service && go fmt ./... && go vet ./...; echo $$? > /tmp/lint_alert_service_exit
	@data_exit=$$(cat /tmp/lint_data_ingestion_exit); ma_exit=$$(cat /tmp/lint_ma_signal_exit); volume_exit=$$(cat /tmp/lint_volume_spike_exit); alert_exit=$$(cat /tmp/lint_alert_service_exit); \
	total_exit=$$(($$data_exit + $$ma_exit + $$volume_exit + $$alert_exit)); \
	rm -f /tmp/lint_*_exit; \
	if [ $$total_exit -eq 0 ]; then echo "All linting completed successfully"; else echo "Linting failed in one or more services"; exit 1; fi

deploy:
	helm upgrade --install crypto-trackers ./helm/crypto-trackers --create-namespace --namespace crypto-trackers

verify:
	./scripts/verify.sh

monitor:
	@echo "Starting port-forward for Prometheus (9090) and Grafana (3000)..."
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3000 (admin/admin)"
	@echo "Press Ctrl+C to stop"
	kubectl port-forward -n crypto-trackers svc/prometheus-service 9090:9090 & \
	kubectl port-forward -n crypto-trackers svc/grafana 3000:3000 & \
	wait

cleanup:
	./scripts/cleanup.sh

build-integration-test:
	docker build -t crypto-trackers/integration-test-runner:latest ./integration-tests/

integration-test: build-integration-test
	@echo "Running integration tests..."
	kubectl delete job integration-test-runner -n crypto-trackers --ignore-not-found=true
	kubectl apply -f ./integration-tests/test-runner-job.yaml
	kubectl wait --for=condition=complete --timeout=300s job/integration-test-runner -n crypto-trackers
	kubectl logs job/integration-test-runner -n crypto-trackers
