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
