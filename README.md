# Crypto Trading Signal Detection System

A distributed microservices-based system for detecting cryptocurrency trading signals and sending alerts. The system processes real-time market data through multiple signal detection algorithms and delivers actionable insights.

## Architecture Overview

This system consists of multiple microservices deployed via Helm charts on Kubernetes:

- **Data Ingestion Service** (Python): Collects real-time cryptocurrency market data from various exchanges
- **Moving Average Signal Detector** (Go): Detects trading signals based on moving average crossovers
- **Volume Spike Detector** (Go): Identifies unusual volume spikes that may indicate significant price movements
- **Alert Service** (Go): Manages and delivers trading alerts to users via multiple channels

## Technology Stack

- **Container Orchestration**: Kubernetes
- **Deployment**: Helm v3+ charts
- **Message Streaming**: Apache Kafka
- **Languages**: Python (data ingestion), Go (signal processing and alerts)
- **Monitoring**: Prometheus & Grafana (planned)

## Quick Start

1. **Prerequisites**
   - Kubernetes cluster (minikube, kind, or cloud provider)
   - kubectl configured
   - Helm v3+ installed

2. **Deploy using Helm**
   ```bash
   # Install the base infrastructure (namespace and configuration)
   helm install crypto-signals ./helm/crypto-signals
   
   # Verify deployment
   kubectl get all -n crypto-signals
   ```

3. **Deploy additional components (when implemented)**
   ```bash
   # All components including Kafka and services will be deployed via Helm
   # Additional Helm values can be customized as needed
   helm upgrade crypto-signals ./helm/crypto-signals --set kafka.enabled=true
   ```

## Project Structure

```
crypto-trading-signals/
├── services/                   # Microservices source code (planned)
│   ├── data-ingestion/        # Python service for data collection
│   ├── ma-signal-detector/    # Go service for MA signal detection
│   ├── volume-spike-detector/ # Go service for volume spike detection
│   └── alert-service/         # Go service for alert management
├── helm/                      # Helm charts for deployment
│   └── crypto-signals/       # Main chart for the system
├── k8s/                       # Legacy Kubernetes manifests (reference)
│   ├── namespace/            # Namespace and global config
│   ├── kafka/               # Kafka deployment (planned)
│   ├── services/            # Service deployments (planned)
│   └── monitoring/          # Monitoring stack (planned)
└── docs/                    # Additional documentation
```

## Development

Each service directory contains its own README with specific development instructions.

## Documentation

- [System Design](DESIGN.md) - Detailed system architecture and design decisions
- [Project Structure](PROJECT_STRUCTURE.md) - Detailed project organization
- [Roadmap](ROADMAP.md) - Development roadmap and milestones
