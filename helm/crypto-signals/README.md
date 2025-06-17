# Crypto Signals Helm Chart

This Helm chart deploys the basic infrastructure for the crypto trading signal detection system.

## Current Status

This chart currently deploys:
- Namespace (`crypto-signals`)
- Global ConfigMap with system configuration
- Kafka cluster with 3 brokers
- Zookeeper ensemble with 3 nodes
- Persistent storage for Kafka and Zookeeper
- Anti-affinity rules for proper broker distribution

Future versions will include microservices as they are implemented.

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
| `config.logLevel` | Log level | `info` |
| `config.logFormat` | Log format | `json` |
| `config.apiTimeout` | API timeout seconds | `30` |
| `config.apiRetryAttempts` | API retry attempts | `3` |
| `kafka.replicas` | Number of Kafka brokers | `3` |
| `kafka.image.repository` | Kafka container image | `confluentinc/cp-kafka` |
| `kafka.image.tag` | Kafka image tag | `7.4.0` |
| `kafka.resources.requests.memory` | Kafka memory request | `1Gi` |
| `kafka.resources.requests.cpu` | Kafka CPU request | `500m` |
| `kafka.storage.size` | Kafka persistent storage size | `10Gi` |
| `kafka.jvmHeap` | Kafka JVM heap size | `1G` |
| `zookeeper.replicas` | Number of Zookeeper nodes | `3` |
| `zookeeper.image.repository` | Zookeeper container image | `confluentinc/cp-zookeeper` |
| `zookeeper.image.tag` | Zookeeper image tag | `7.4.0` |
| `zookeeper.storage.size` | Zookeeper persistent storage size | `5Gi` |

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

## Kafka Configuration

The Kafka cluster is configured with:
- 3 brokers for high availability
- Replication factor of 3 for data durability
- Anti-affinity rules to distribute brokers across nodes
- Persistent storage for data retention
- Production-ready resource limits

The Zookeeper ensemble provides:
- 3-node cluster for consensus
- Persistent storage for metadata
- Anti-affinity for node distribution