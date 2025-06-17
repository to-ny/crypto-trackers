# Crypto Trading Signal Detection System

A distributed microservices-based system for detecting cryptocurrency trading signals and sending alerts. The system processes real-time market data through multiple signal detection algorithms and delivers actionable insights.

## Architecture Overview

This system consists of multiple microservices running on Kubernetes:

- **Data Ingestion Service** (Python): Collects real-time cryptocurrency market data from various exchanges
- **Moving Average Signal Detector** (Go): Detects trading signals based on moving average crossovers
- **Volume Spike Detector** (Go): Identifies unusual volume spikes that may indicate significant price movements
- **Alert Service** (Go): Manages and delivers trading alerts to users via multiple channels

## Technology Stack

- **Container Orchestration**: Kubernetes
- **Message Streaming**: Apache Kafka
- **Languages**: Python (data ingestion), Go (signal processing and alerts)
- **Monitoring**: Prometheus & Grafana (planned)

## Quick Start

1. **Prerequisites**
   - Kubernetes cluster (minikube, kind, or cloud provider)
   - kubectl configured
   - Helm (for Kafka installation)

2. **Deploy the namespace and base configuration**
   ```bash
   kubectl apply -f k8s/namespace/
   ```

3. **Deploy Kafka**
   ```bash
   kubectl apply -f k8s/kafka/
   ```

4. **Deploy services**
   ```bash
   kubectl apply -f k8s/services/
   ```

## Project Structure

```
crypto-trading-signals/
├── services/                   # Microservices source code
│   ├── data-ingestion/        # Python service for data collection
│   ├── ma-signal-detector/    # Go service for MA signal detection
│   ├── volume-spike-detector/ # Go service for volume spike detection
│   └── alert-service/         # Go service for alert management
├── k8s/                       # Kubernetes manifests
│   ├── namespace/            # Namespace and global config
│   ├── kafka/               # Kafka deployment
│   ├── services/            # Service deployments
│   └── monitoring/          # Monitoring stack
└── docs/                    # Additional documentation
```

## Development

Each service directory contains its own README with specific development instructions.

## Documentation

- [System Design](DESIGN.md) - Detailed system architecture and design decisions
- [Project Structure](PROJECT_STRUCTURE.md) - Detailed project organization
- [Roadmap](ROADMAP.md) - Development roadmap and milestones

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

[License information to be added]