# Crypto Signals Helm Chart

This Helm chart deploys the basic infrastructure for the crypto trading signal detection system.

## Current Status

This chart currently deploys:
- Namespace (`crypto-signals`)
- Global ConfigMap with system configuration

Future versions will include Kafka and microservices as they are implemented.

## Installation

```bash
# Install the chart
helm install crypto-signals ./helm/crypto-signals

# Verify installation
kubectl get all -n crypto-signals

# Check configuration
kubectl get configmap -n crypto-signals crypto-signals-config -o yaml
```

## Uninstallation

```bash
helm uninstall crypto-signals
```

## Configuration

The following table lists the configurable parameters and their default values:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `namespace` | Kubernetes namespace | `crypto-signals` |
| `config.kafkaBootstrapServers` | Kafka bootstrap servers | `kafka:9092` |
| `config.kafkaTopicPrefix` | Kafka topic prefix | `crypto-signals` |
| `config.logLevel` | Log level | `info` |
| `config.logFormat` | Log format | `json` |
| `config.apiTimeout` | API timeout seconds | `30` |
| `config.apiRetryAttempts` | API retry attempts | `3` |

## Customization

Create a custom values file:

```yaml
# custom-values.yaml
namespace: my-crypto-signals
config:
  logLevel: debug
  kafkaBootstrapServers: my-kafka:9092
```

Install with custom values:

```bash
helm install crypto-signals ./helm/crypto-signals -f custom-values.yaml
```