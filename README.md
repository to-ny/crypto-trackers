# Crypto Trading Signal Detector

Kubernetes-based cryptocurrency trading signal detection system.

## Services

- **Data Ingestion** (Python): Fetches BTC/ETH prices from CoinGecko
- **MA Signal Detector** (Go): Detects SMA 20/50 crossovers  
- **Volume Spike Detector** (Go): Detects volume spikes
- **Alert Service** (Go): Rate-limited alerts

## Monitoring

- **Prometheus**: Metrics collection at http://localhost:9090

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
