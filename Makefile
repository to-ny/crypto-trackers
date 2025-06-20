.PHONY: deploy verify cleanup help

help:
	@echo "Available targets:"
	@echo "  deploy   - Deploy crypto-trackers to Kubernetes"
	@echo "  verify   - Verify deployment is working"
	@echo "  cleanup  - Remove deployment and cleanup resources"

deploy:
	helm upgrade --install crypto-trackers ./helm/crypto-trackers --create-namespace --namespace crypto-trackers

verify:
	./scripts/verify.sh

cleanup:
	./scripts/cleanup.sh