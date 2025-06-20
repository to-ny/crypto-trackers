# Crypto Trading Signal Detector

Kubernetes-based system for cryptocurrency trading signal detection.

## Prerequisites

- Kubernetes cluster
- Helm 3.x
- kubectl configured

## Components

- Namespace: crypto-trackers
- Kafka cluster (`kafka-service:9092`)
- ZooKeeper (`zookeeper-service:2181`)
- Topics: `crypto-prices`, `trading-signals`
- Data Ingestion Service: Fetches BTC/ETH prices from CoinGecko API
- MA Signal Detector: Detects SMA 20/50 crossovers and generates trading signals
- Volume Spike Detector: Detects volume spikes above 7-day average threshold

## Usage

```bash
# Deploy the system
make deploy

# Verify deployment
make verify

# Cleanup resources
make cleanup

# Show available commands
make help
```
