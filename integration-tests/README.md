# Integration Tests

End-to-end tests for the signal detection pipeline.

## Test Scenarios

- **Death Cross**: SMA 20/50 bearish crossover detection  
- **Volume Spike**: Volume spike detection above 7-day average

## Development

```bash
# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Format code
python -m black .

# Lint code
python -m ruff check .
python -m mypy .

# Build test image
docker build -t crypto-trackers/integration-test-runner:latest .

# Deploy and run tests
kubectl apply -f test-runner-job.yaml
kubectl wait --for=condition=complete --timeout=300s job/integration-test-runner -n crypto-trackers
kubectl logs job/integration-test-runner -n crypto-trackers
```