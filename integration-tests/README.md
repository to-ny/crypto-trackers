# Integration Tests

Automated end-to-end tests that verify the complete signal detection pipeline.

## Test Scenarios

- **Golden Cross**: Tests SMA 20/50 bullish crossover detection
- **Death Cross**: Tests SMA 20/50 bearish crossover detection  
- **Volume Spike**: Tests volume spike detection above 7-day average

## Test Process

1. **Setup**: Scale data-ingestion to 0, restart detectors, reset consumer offsets
2. **Inject**: Generate 60 test events per scenario directly to Kafka
3. **Verify**: Monitor trading-signals topic for expected outputs within 60s
4. **Cleanup**: Restore data-ingestion to 1, cleanup test artifacts

## Usage

```bash
# Build test image
docker build -t crypto-trackers/integration-test-runner:latest .

# Deploy and run tests
helm template crypto-trackers ../helm/crypto-trackers --show-only templates/integration-tests/test-runner-job.yaml | kubectl apply -f -
kubectl wait --for=condition=complete --timeout=300s job/integration-test-runner -n crypto-trackers
kubectl logs job/integration-test-runner -n crypto-trackers
```

## Expected Results

- Golden cross: 1 bullish signal
- Death cross: 1 bearish signal
- Volume spike: 1 spike signal

Tests run in isolated environment without affecting live data ingestion.