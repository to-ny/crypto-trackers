# Data Ingestion Service

Fetches BTC/ETH price data from CoinGecko API and publishes events to Kafka.

## Development

```bash
# Install dependencies
pip install -r requirements.txt

# Run service
python src/main.py

# Run tests
python -m pytest tests/ -v
```

## Environment Variables

- `KAFKA_BOOTSTRAP_SERVERS`: Kafka cluster address (default: `kafka:9092`)
- `POLLING_INTERVAL`: Seconds between API calls (default: `60`)
- `PORT`: HTTP server port (default: `8080`)

## Build

```bash
docker build -t data-ingestion .
```