# Volume Spike Detector

Detects volume spikes from Kafka price events and publishes trading signals.

## Development

```bash
# Build service
go build ./cmd/main.go

# Run locally
KAFKA_BOOTSTRAP_SERVERS=localhost:9092 ./main
```

## Environment Variables

- `KAFKA_BOOTSTRAP_SERVERS`: Kafka cluster address (default: `kafka-service:9092`)
- `KAFKA_GROUP_ID`: Consumer group ID (default: `volume-spike-detector`)
- `PORT`: HTTP server port (default: `8080`)
- `SPIKE_THRESHOLD`: Volume spike threshold multiplier (default: `1.3`)

## Build

```bash
# Build Docker image
docker build -t crypto-trackers/volume-spike-detector:latest .

# Run container locally
docker run -p 8080:8080 \
  -e KAFKA_BOOTSTRAP_SERVERS=localhost:9092 \
  crypto-trackers/volume-spike-detector:latest
```

## Deployment

Service is deployed as part of the main crypto-trackers Helm chart located at `/helm/crypto-trackers/`.