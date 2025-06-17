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

2. **Deployment**
   ```bash
   # Install the system
   helm install crypto-signals ./helm/crypto-signals
   
   # Verify deployment
   kubectl get all -n crypto-signals
   ```

## Project Structure

```
crypto-trading-signals/
├── services/                   # Microservices source code (to be implemented)
│   ├── data-ingestion/        # Python service for data collection
│   ├── ma-signal-detector/    # Go service for MA signal detection
│   ├── volume-spike-detector/ # Go service for volume spike detection
│   └── alert-service/         # Go service for alert management
├── helm/                      # Helm charts for deployment
│   └── crypto-signals/       # Main chart with Kafka cluster
│       ├── templates/        # Kubernetes resource templates
│       │   ├── kafka-*.yaml # Kafka StatefulSet and Services
│       │   ├── zookeeper-*.yaml # Zookeeper StatefulSet and Services
│       │   ├── namespace.yaml # Namespace creation
│       │   └── configmap.yaml # Global configuration
│       ├── values.yaml      # Configuration values
│       ├── test-deployment.sh # Deployment testing script
│       └── README.md        # Chart documentation
├── DESIGN.md                  # System architecture and design
├── ROADMAP.md                # Implementation roadmap
└── PROJECT_STRUCTURE.md      # Project organization details
```

## Development

Each service directory contains its own README with specific development instructions.

## Documentation

- [System Design](DESIGN.md) - Detailed system architecture and design decisions
- [Project Structure](PROJECT_STRUCTURE.md) - Detailed project organization
- [Roadmap](ROADMAP.md) - Development roadmap and milestones
