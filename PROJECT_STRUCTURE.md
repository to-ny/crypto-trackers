# Project Structure

## Overview

High-level organization for the system.

## Root Structure

```
crypto-trackers/
├── README.md
├── DESIGN.md
├── ROADMAP.md
├── services/           # All microservices
└── helm/              # Kubernetes deployment
```

## Service Organization

Each service is self-contained in `services/{service-name}/`:

```
services/
├── data-ingestion/          # Python - CoinGecko API integration
├── ma-signal-detector/      # Go - Moving average signals
├── volume-spike-detector/   # Go - Volume spike detection
└── alert-service/          # Go - Signal notifications
```

## Deployment Structure

```
helm/crypto-trackers/
├── Chart.yaml
├── values.yaml
└── templates/
    ├── namespace.yaml
    ├── configmap.yaml
    ├── kafka/           # Kafka cluster resources
    └── services/        # Microservice deployments
```

## Key Principles

- **Service Independence**: Each service has its own Docker image and deployment
- **Configuration Management**: Environment-specific values in Helm values files