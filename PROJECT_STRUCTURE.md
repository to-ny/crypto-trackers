# Project Structure Specification

## Overview

This document defines the standardized project structure for the Crypto Trading Signal Detector system. This structure supports multiple services in different languages while maintaining clear separation and easy navigation for development and deployment.

## Root Directory Structure

```
crypto-trading-signals/
├── README.md                     # Project overview and quick start
├── DESIGN.md                     # Software design document
├── ROADMAP.md                    # Implementation roadmap
├── PROJECT_STRUCTURE.md          # This document
├── .gitignore                    # Git ignore patterns
│
├── services/                     # All microservices
│   ├── data-ingestion/          # Python data ingestion service
│   ├── ma-signal-detector/      # Go moving average service
│   ├── volume-spike-detector/   # Go volume spike service
│   └── alert-service/           # Go alert service
│
└── helm/                        # Helm charts for deployment
    └── crypto-signals/         # Main Helm chart
        ├── Chart.yaml          # Chart metadata
        ├── values.yaml         # Default configuration values
        └── templates/          # Kubernetes resource templates
            ├── _helpers.tpl    # Template helpers
            ├── namespace.yaml  # Namespace template
            └── configmap.yaml  # ConfigMap template
```

## Service Directory Structure

### Python Service (data-ingestion)

```
services/data-ingestion/
├── Dockerfile                   # Container image definition
├── requirements.txt             # Python dependencies
├── .dockerignore               # Docker build exclusions
├── README.md                   # Service-specific documentation
│
├── src/                        # Source code
│   ├── __init__.py
│   ├── main.py                 # Application entry point
│   ├── config.py               # Configuration management
│   ├── health.py               # Health check endpoints
│   │
│   ├── api/                    # External API integration
│   │   ├── __init__.py
│   │   ├── coingecko.py        # CoinGecko API client
│   │   └── models.py           # API response models
│   │
│   ├── kafka/                  # Kafka integration
│   │   ├── __init__.py
│   │   ├── producer.py         # Kafka producer client
│   │   └── serializers.py      # Message serialization
│   │
│   └── utils/                  # Utility functions
│       ├── __init__.py
│       ├── logging.py          # Logging configuration
│       └── errors.py           # Error handling
│
└── tests/                      # Unit tests
    ├── __init__.py
    ├── test_main.py
    ├── test_api/
    └── test_kafka/
```

### Go Service Template (ma-signal-detector, volume-spike-detector, alert-service)

```
services/ma-signal-detector/
├── Dockerfile                   # Container image definition
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── .dockerignore               # Docker build exclusions
├── README.md                   # Service-specific documentation
│
├── cmd/                        # Application entry points
│   └── main.go                 # Main application
│
├── internal/                   # Private application code
│   ├── config/                 # Configuration management
│   │   └── config.go
│   │
│   ├── handlers/               # HTTP handlers
│   │   ├── health.go           # Health check endpoints
│   │   └── metrics.go          # Metrics endpoints
│   │
│   ├── kafka/                  # Kafka integration
│   │   ├── consumer.go         # Kafka consumer
│   │   ├── producer.go         # Kafka producer
│   │   └── models.go           # Message models
│   │
│   ├── signals/                # Signal detection logic
│   │   ├── ma.go               # Moving average calculations
│   │   ├── history.go          # Price history management
│   │   └── detector.go         # Signal detection engine
│   │
│   └── utils/                  # Utility packages
│       ├── logging.go          # Logging utilities
│       └── errors.go           # Error handling
│
└── tests/                      # Test files
    ├── integration/            # Integration tests
    └── unit/                   # Unit tests
```

## Deployment Structure

### Helm Charts

All deployments are managed through Helm charts in the `helm/crypto-signals/` directory:

```
helm/crypto-signals/
├── Chart.yaml                  # Chart metadata and dependencies
├── values.yaml                 # Default configuration values
├── README.md                   # Chart documentation and usage
│
├── templates/
│   ├── _helpers.tpl           # Template helper functions
│   ├── namespace.yaml         # Crypto-signals namespace
│   ├── configmap.yaml         # Global configuration
│   │
│   ├── kafka/                 # Kafka cluster resources (planned)
│   │   ├── zookeeper.yaml     # Zookeeper StatefulSet
│   │   ├── kafka.yaml         # Kafka StatefulSet
│   │   ├── kafka-service.yaml # Kafka service discovery
│   │   └── topics.yaml        # Topic creation job
│   │
│   ├── services/              # Microservice deployments (planned)
│   │   ├── data-ingestion/
│   │   │   ├── deployment.yaml    # Service deployment
│   │   │   ├── service.yaml       # Kubernetes service
│   │   │   ├── configmap.yaml     # Service-specific config
│   │   │   └── hpa.yaml          # Horizontal Pod Autoscaler
│   │   │
│   │   ├── ma-signal-detector/
│   │   │   ├── deployment.yaml
│   │   │   ├── service.yaml
│   │   │   ├── configmap.yaml
│   │   │   └── hpa.yaml
│   │   │
│   │   ├── volume-spike-detector/
│   │   │   ├── deployment.yaml
│   │   │   ├── service.yaml
│   │   │   ├── configmap.yaml
│   │   │   └── hpa.yaml
│   │   │
│   │   └── alert-service/
│   │       ├── deployment.yaml
│   │       ├── service.yaml
│   │       ├── configmap.yaml
│   │       └── hpa.yaml
│   │
│   ├── monitoring/            # Monitoring resources (planned)
│   │   ├── health-checks.yaml     # Health monitoring
│   │   ├── service-monitor.yaml   # Prometheus metrics collection
│   │   └── dashboards.yaml        # Grafana dashboards
│   │
│   └── secrets/               # Secret management (planned)
│       ├── api-keys.yaml      # External API credentials
│       └── kafka-auth.yaml    # Kafka authentication
│
└── charts/                    # Subcharts (planned)
    ├── kafka/                 # Kafka subchart
    └── monitoring/            # Monitoring subchart
```

## File Naming Conventions

### General Rules
- Use kebab-case for directories: `ma-signal-detector`
- Use kebab-case for YAML files: `kafka-service.yaml`
- Use snake_case for Python files: `coingecko_client.py`
- Use camelCase for Go files: `signalDetector.go`

### Docker Images
- Format: `crypto-signals/{service-name}:{version}`
- Examples: 
  - `crypto-signals/data-ingestion:v1.0.0`
  - `crypto-signals/ma-signal-detector:v1.0.0`

### Kubernetes Resources
- Include service name in resource names
- Format: `{service-name}-{resource-type}`
- Examples:
  - `data-ingestion-deployment`
  - `ma-signal-detector-service`

## Development Workflow

### Service Development
1. Create service directory using appropriate template
2. Implement core business logic first
3. Add Kafka integration
4. Create Dockerfile and test locally
5. Add Kubernetes manifests
6. Update service README with specific instructions

### Infrastructure Changes
1. Update Helm chart templates in `helm/crypto-signals/templates/`
2. Update values in `helm/crypto-signals/values.yaml`
3. Test changes in local environment using `helm install/upgrade`

### Documentation Updates
1. Service-specific docs in service README
2. Architecture changes in DESIGN.md
3. Any additional documentation can be added to root or service directories as needed
