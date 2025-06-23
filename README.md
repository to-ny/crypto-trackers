# Crypto Trading Signal Detector

Kubernetes-based cryptocurrency trading signal detection system.

## Context

This project has been developed during a week of June 2025. The main goal was to familiarize myself working with an AI coding agent and see how far I could push it.

I chose a somewhat complex system (multiple services in several languages, async communication using Kafka, deployment using Kubernetes & Helm) to have a more "real-world" project to work on.

The result has been very positive for me. Aside from a few initial adjustements, it was a bliss to make Claude build this project and to review the changes and guide it step by step.

## Services

- [**Data Ingestion**](services/data-ingestion/README.md) (Python): Fetches BTC/ETH prices from CoinGecko
- [**MA Signal Detector**](services/ma-signal-detector/README.md) (Go): Detects SMA 20/50 crossovers
- [**Volume Spike Detector**](services/volume-spike-detector/README.md) (Go): Detects volume spikes
- [**Alert Service**](services/alert-service/README.md) (Go): Rate-limited alerts

## Development

```bash
make build    # Build images
make test     # Run tests
make lint     # Format code
```

## Deployment

```bash
make deploy   # Deploy system
make verify   # Check status
make monitor  # Access monitoring dashboard
make cleanup  # Remove deployment
```

## Integration Testing

```bash
make integration-test
```

See [integration-tests/README.md](integration-tests/README.md) for details.
