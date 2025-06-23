# Moving Average Signal Detector

Detects SMA 20/50 crossovers from Kafka price events and publishes trading signals.

## Development

```bash
# Build service
go build ./cmd/main.go

# Run locally
KAFKA_BOOTSTRAP_SERVERS=localhost:9092 ./main

# Run tests
go test ./... -v

# Code quality
go fmt ./...
go vet ./...
```

## Environment Variables

- `KAFKA_BOOTSTRAP_SERVERS`: Kafka cluster address (default: `kafka-service:9092`)
- `KAFKA_GROUP_ID`: Consumer group ID (default: `ma-signal-detector`)
- `PORT`: HTTP server port (default: `8080`)

## Build

```bash
# Build Docker image
docker build -t crypto-trackers/ma-signal-detector:latest .

# Run container locally
docker run -p 8080:8080 \
  -e KAFKA_BOOTSTRAP_SERVERS=localhost:9092 \
  crypto-trackers/ma-signal-detector:latest
```

## Deployment

Service is deployed as part of the main crypto-trackers Helm chart located at `/helm/crypto-trackers/`.
