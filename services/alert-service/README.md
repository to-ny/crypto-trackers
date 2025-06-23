# Alert Service

Consumes trading signals from Kafka and outputs structured alerts with rate limiting.

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
- `KAFKA_GROUP_ID`: Consumer group ID (default: `alert-service`)
- `PORT`: HTTP server port (default: `8080`)
- `COOLDOWN_MINUTES`: Rate limiting cooldown period per symbol (default: `5`)

## Build

```bash
# Build Docker image
docker build -t crypto-trackers/alert-service:latest .

# Run container locally
docker run -p 8080:8080 \
  -e KAFKA_BOOTSTRAP_SERVERS=localhost:9092 \
  crypto-trackers/alert-service:latest
```

## Deployment

Service is deployed as part of the main crypto-trackers Helm chart located at `/helm/crypto-trackers/`.
