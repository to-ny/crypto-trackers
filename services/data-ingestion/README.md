# Data Ingestion Service

Fetches BTC/ETH price data from CoinGecko API and publishes events to Kafka.

## Development

```bash
# Setup virtual environment
python -m venv venv
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt
pip install -r requirements-dev.txt

# Run service
python src/main.py

# Run tests
python -m pytest tests/ -v

# Code quality
ruff check src/ tests/
black src/ tests/
mypy src/
```

## Environment Variables

- `KAFKA_BOOTSTRAP_SERVERS`: Kafka cluster address (default: `kafka:9092`)
- `POLLING_INTERVAL`: Seconds between API calls (default: `60`)
- `PORT`: HTTP server port (default: `8080`)

## Build

```bash
# Build Docker image
docker build -t crypto-trackers/data-ingestion:latest .

# Run container locally
docker run -p 8080:8080 \
  -e KAFKA_BOOTSTRAP_SERVERS=localhost:9092 \
  -e POLLING_INTERVAL=60 \
  crypto-trackers/data-ingestion:latest
```

# Deployment

Service is deployed as part of the main crypto-trackers Helm chart located at `/helm/crypto-trackers/`.
