.PHONY: help build test lint deploy verify cleanup

help:
	@echo "build    - Build all images"
	@echo "test     - Run all tests"
	@echo "lint     - Format and lint code"
	@echo "deploy   - Deploy system"
	@echo "verify   - Check deployment"
	@echo "cleanup  - Remove deployment"

build:
	docker build -t crypto-trackers/data-ingestion:latest ./services/data-ingestion/
	cd ./services/ma-signal-detector && docker build -t crypto-trackers/ma-signal-detector:latest .
	cd ./services/volume-spike-detector && docker build -t crypto-trackers/volume-spike-detector:latest .
	cd ./services/alert-service && docker build -t crypto-trackers/alert-service:latest .

test:
	cd ./services/data-ingestion && python -m pytest
	cd ./services/ma-signal-detector && go test ./... -v
	cd ./services/volume-spike-detector && go test ./... -v
	cd ./services/alert-service && go test ./... -v

lint:
	cd ./services/data-ingestion && python -m black . && python -m ruff check . && python -m mypy .
	cd ./services/ma-signal-detector && go fmt ./... && go vet ./...
	cd ./services/volume-spike-detector && go fmt ./... && go vet ./...
	cd ./services/alert-service && go fmt ./... && go vet ./...

deploy:
	helm upgrade --install crypto-trackers ./helm/crypto-trackers --create-namespace --namespace crypto-trackers

verify:
	./scripts/verify.sh

cleanup:
	./scripts/cleanup.sh
