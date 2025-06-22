# Crypto Trading Signal Detector

Kubernetes-based cryptocurrency trading signal detection system.

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
