# Crypto Trading Signal Detector - Software Design Document

## 1. Introduction and Overview

### 1.1 Project Description

The Crypto Trading Signal Detector is a real-time distributed system that monitors cryptocurrency market data and generates trading signals based on technical analysis indicators. The system continuously ingests price and volume data from external APIs, processes this data through configurable signal detection algorithms, and delivers actionable alerts to users.

### 1.2 Project Objectives and Key Requirements

**Primary Objectives:**
Practice essential software engineering concepts including distributed systems architecture, event-driven design patterns, and multi-language development using real-world cryptocurrency market data.

**Key Requirements:**
- **Real-time Processing**: System must process market data with minimal latency (< 5 seconds from data ingestion to signal generation)
- **Scalability**: Architecture must support easy addition of new cryptocurrencies and signal types
- **Reliability**: System should handle API failures gracefully and maintain data consistency
- **Observability**: Comprehensive logging, metrics, and health monitoring
- **Cloud-Native**: Full Kubernetes deployment with proper resource management

**Success Criteria:**
- Successfully detect and alert on moving average crossovers and volume spikes for BTC and ETH
- System runs continuously without manual intervention
- Adding a new cryptocurrency requires only configuration changes
- Complete Kubernetes deployment that can be reproduced in any cluster

### 1.3 Document Purpose and Audience

This document serves as the comprehensive design specification for the Crypto Trading Signal Detector system. It is intended for:
- **Personal Reference** to maintain clear architectural foundations and design decisions
- **AI Code Assistants** to ensure alignment and consistency during development sessions

### 1.4 Document Contents

This document covers:
- **System Architecture**: High-level component design and data flow
- **Technology Stack**: Framework and tool selections with rationale
- **Data Models**: Message schemas and contracts
- **Service Specifications**: Detailed component responsibilities
- **Deployment Strategy**: Kubernetes configuration and scaling approach
- **Monitoring and Observability**: Logging, metrics, and alerting strategy

### 1.5 Background Information

**Market Context:**
Cryptocurrency markets operate 24/7 with high volatility, making them ideal for real-time signal detection systems. Traditional technical analysis indicators like moving averages and volume analysis remain relevant in crypto trading.

**Technical Context:**
This system leverages Apache Kafka for event streaming, providing the foundation for a scalable, fault-tolerant architecture. The microservices approach allows for independent scaling and deployment of different system components.

## 2. System Architecture

### 2.1 High-Level Architecture

The system follows an event-driven microservices architecture with the following data flow:

```
External APIs → Data Ingestion → Kafka Topics → Signal Detection → Alert Generation
```

### 2.2 Core Components

**Data Ingestion Service**
- Polls CoinGecko API for BTC and ETH price/volume data
- Transforms raw API responses into standardized events
- Publishes price events to Kafka topics
- Language: Python (leveraging rich data processing libraries)

**Kafka Event Streaming Platform**
- Central message broker for all system communication
- Topics: `crypto-prices`, `trading-signals`
- Provides durability, scalability, and decoupling between services

**Signal Detection Services**
- Moving Average Service: Detects crossover patterns (SMA 20/50)
- Volume Spike Service: Identifies unusual trading volume
- Each service consumes price events and publishes signal events
- Language: Go (for high-performance concurrent processing)

**Alert Service**
- Consumes trading signals from Kafka
- Generates notifications (console output, structured logs)
- Language: Go (consistent with signal detection services)

### 2.3 Data Flow

1. **Ingestion Phase**: Data Ingestion Service fetches market data every 60 seconds
2. **Event Publishing**: Raw price data published to `crypto-prices` topic
3. **Signal Processing**: Detection services consume price events, apply algorithms
4. **Signal Publishing**: Detected signals published to `trading-signals` topic
5. **Alert Generation**: Alert Service consumes signals and generates notifications

### 2.4 Scalability Design

- **Horizontal Scaling**: Each service can be independently scaled based on load
- **Topic Partitioning**: Kafka topics partitioned by cryptocurrency symbol
- **Stateless Services**: All services designed to be stateless for easy scaling

## 3. Technology Stack

### 3.1 Core Technologies

**Apache Kafka**
- Event streaming platform for service communication
- Provides durability, fault tolerance, and high throughput
- Chosen for its proven scalability in financial data processing

**Kubernetes**
- Container orchestration platform for deployment and scaling
- Provides service discovery, load balancing, and resource management
- Industry standard for cloud-native applications

### 3.2 Programming Languages

**Python (Data Ingestion Service)**
- Rich ecosystem for API integration and data processing
- Libraries: `requests` for HTTP calls, `kafka-python` for Kafka integration
- Rapid development for data transformation tasks

**Go (Signal Detection & Alert Services)**
- High-performance concurrent processing capabilities
- Native Kafka client with excellent performance characteristics
- Compiled binaries reduce container overhead

### 3.3 Infrastructure Components

**Docker**
- Containerization for consistent deployment across environments
- Base images: `python:3.11-slim`, `golang:1.21-alpine`

**Helm Charts**
- Kubernetes deployment templating and management
- Simplified configuration management across environments

### 3.4 External Dependencies

**CoinGecko API**
- Free tier provides sufficient rate limits (10-50 calls/minute)
- Reliable cryptocurrency market data source
- RESTful interface with consistent data format

## 4. Data Models

### 4.1 Kafka Message Schemas

All Kafka messages follow a consistent JSON structure with metadata headers for traceability and routing.

### 4.2 Price Event Schema

Published to `crypto-prices` topic by the Data Ingestion Service.

```json
{
  "timestamp": "2024-06-16T14:30:00Z",
  "symbol": "BTC",
  "price_usd": 67450.23,
  "volume_24h": 28450000000,
  "market_cap": 1330000000000,
  "price_change_24h": 2.35,
  "source": "coingecko"
}
```

**Field Descriptions:**
- `timestamp`: ISO 8601 UTC timestamp of the price data
- `symbol`: Cryptocurrency symbol (BTC, ETH)
- `price_usd`: Current price in USD
- `volume_24h`: 24-hour trading volume in USD
- `market_cap`: Total market capitalization in USD
- `price_change_24h`: 24-hour price change percentage
- `source`: Data provider identifier

### 4.3 Trading Signal Schema

Published to `trading-signals` topic by Signal Detection Services.

```json
{
  "timestamp": "2024-06-16T14:30:05Z",
  "symbol": "BTC",
  "signal_type": "moving_average_crossover",
  "signal_strength": "strong",
  "direction": "bullish",
  "details": {
    "sma_20": 67200.45,
    "sma_50": 66800.12,
    "crossover_type": "golden_cross"
  },
  "service_id": "ma-detector-v1"
}
```

**Field Descriptions:**
- `timestamp`: Signal generation timestamp
- `symbol`: Cryptocurrency symbol
- `signal_type`: Type of signal (`moving_average_crossover`, `volume_spike`)
- `signal_strength`: Signal confidence (`weak`, `moderate`, `strong`)
- `direction`: Market direction (`bullish`, `bearish`, `neutral`)
- `details`: Signal-specific metadata (varies by signal type)
- `service_id`: Identifying the generating service for debugging

### 4.4 Volume Spike Signal Details

For `volume_spike` signal type, the details object contains:

```json
{
  "details": {
    "current_volume": 35000000000,
    "avg_volume_7d": 25000000000,
    "spike_multiplier": 1.4,
    "threshold_exceeded": 1.3
  }
}
```

## 5. Service Specifications

### 5.1 Data Ingestion Service

**Responsibilities:**
- Fetch cryptocurrency market data from CoinGecko API
- Transform raw API responses into standardized price events
- Publish price events to Kafka with proper error handling

**Configuration:**
- Polling interval: 60 seconds
- Supported symbols: BTC, ETH (configurable via environment variables)
- API timeout: 30 seconds with exponential backoff retry

**Error Handling:**
- API failures: Log error and continue with next polling cycle
- Kafka publishing failures: Retry up to 3 times with backoff
- Invalid API responses: Log and skip malformed data

**Health Checks:**
- `/health` endpoint for Kubernetes liveness probe
- `/ready` endpoint for readiness probe (checks Kafka connectivity)

### 5.2 Moving Average Signal Detection Service

**Responsibilities:**
- Consume price events from `crypto-prices` topic
- Calculate Simple Moving Averages (SMA 20 and SMA 50)
- Detect crossover patterns and generate trading signals
- Maintain price history for moving average calculations

**Signal Logic:**
- Golden Cross: SMA 20 crosses above SMA 50 (bullish signal)
- Death Cross: SMA 20 crosses below SMA 50 (bearish signal)
- Minimum 50 data points required before generating signals

**State Management:**
- In-memory price history (last 100 data points per symbol)
- No persistent storage required (stateless on restart)

### 5.3 Volume Spike Detection Service

**Responsibilities:**
- Consume price events from `crypto-prices` topic
- Calculate 7-day average volume for each cryptocurrency
- Detect volume spikes above configurable threshold
- Generate volume spike trading signals

**Signal Logic:**
- Spike threshold: Current volume > 1.3x average 7-day volume
- Signal strength based on spike magnitude (1.3x = weak, 2.0x+ = strong)

**State Management:**
- Rolling 7-day volume history per cryptocurrency
- In-memory storage with circular buffer implementation

### 5.4 Alert Service

**Responsibilities:**
- Consume all trading signals from `trading-signals` topic
- Format and output alerts to configured destinations
- Aggregate signals to prevent spam (rate limiting)

**Output Formats:**
- Structured JSON logs for observability
- Console output for development/debugging
- Extensible for future notification channels (webhooks, email)

**Rate Limiting:**
- Maximum 1 alert per symbol per 5-minute window
- Configurable cooldown periods per signal type

## 6. Deployment Strategy

### 6.1 Kubernetes Architecture

The system deploys as a collection of Kubernetes resources organized in a dedicated namespace (`crypto-signals`).

**Resource Types:**
- **StatefulSet**: Kafka cluster with persistent storage
- **Deployments**: All application services (stateless)
- **Services**: Internal service discovery and load balancing
- **ConfigMaps**: Environment-specific configuration
- **Secrets**: API keys and sensitive configuration
- **PersistentVolumes**: Kafka data storage

### 6.2 Kafka Deployment

**StatefulSet Configuration:**
- 3 Kafka brokers for high availability
- Persistent storage: 10GB per broker
- Resource requests: 1 CPU, 2GB RAM per broker
- Anti-affinity rules to distribute brokers across nodes

**Topic Configuration:**
- `crypto-prices`: 2 partitions, replication factor 3
- `trading-signals`: 2 partitions, replication factor 3
- Auto-create disabled for better control

### 6.3 Application Services Deployment

**Data Ingestion Service:**
- Single replica (polling nature doesn't require multiple instances)
- Resource requests: 0.2 CPU, 256MB RAM
- Environment variables for API configuration and polling interval

**Signal Detection Services:**
- 2 replicas each for availability
- Resource requests: 0.3 CPU, 512MB RAM
- Horizontal Pod Autoscaler based on CPU utilization (70% threshold)

**Alert Service:**
- Single replica with readiness for scaling
- Resource requests: 0.1 CPU, 128MB RAM

### 6.4 Service Discovery and Networking

**Internal Services:**
- ClusterIP services for Kafka brokers
- Service discovery via Kubernetes DNS
- No external exposure required for application services

**Health Checks:**
- Liveness probes on `/health` endpoints
- Readiness probes on `/ready` endpoints
- Initial delay: 30 seconds, period: 10 seconds

### 6.5 Configuration Management

**ConfigMaps:**
- `crypto-config`: Supported symbols, polling intervals, thresholds
- `kafka-config`: Broker addresses, topic names, consumer group IDs

**Secrets:**
- `api-secrets`: CoinGecko API key (if premium tier used)
- Base64 encoded and mounted as environment variables

### 6.6 Scaling Strategy

**Vertical Scaling:**
- Resource limits set to 2x requests for burst capacity
- Memory limits prevent OOM kills during data processing spikes

**Horizontal Scaling:**
- HPA configured for signal detection services
- Kafka partitioning enables parallel processing
- Stateless service design supports easy replica scaling

## 7. Monitoring and Observability

### 7.1 Logging Strategy

**Structured Logging:**
- JSON format for all application logs
- Standard fields: timestamp, level, service, message, correlation_id
- Kubernetes stdout/stderr collection for centralized logging

**Log Levels:**
- INFO: Normal operations, signal generations, API calls
- WARN: Recoverable errors, API timeouts, retry attempts
- ERROR: Service failures, unrecoverable errors, data corruption
- DEBUG: Detailed processing information (disabled in production)

### 7.2 Health Monitoring

**Kubernetes Health Checks:**
- Liveness probes prevent stuck containers
- Readiness probes manage traffic routing
- Startup probes handle slow initialization

**Service Health Endpoints:**
- `/health`: Basic service status
- `/ready`: Dependencies check (Kafka connectivity, API availability)
- HTTP 200 for healthy, 503 for unhealthy states

### 7.3 Application Metrics

**Key Performance Indicators:**
- API call success rate and response times
- Kafka message production/consumption rates
- Signal generation frequency per cryptocurrency
- Alert delivery success rates

**Resource Metrics:**
- CPU and memory utilization per service
- Kafka topic partition usage and lag
- Message throughput and processing latency

### 7.4 Error Tracking

**Error Categories:**
- External API failures (rate limits, timeouts, invalid responses)
- Kafka connectivity issues (broker unavailable, topic errors)
- Data processing errors (malformed messages, calculation failures)
- Resource constraints (memory exhaustion, CPU throttling)

**Error Response:**
- Graceful degradation for non-critical failures
- Circuit breaker pattern for external API calls
- Dead letter queues for unprocessable messages

### 7.5 Operational Dashboards

**System Overview:**
- Service status and replica counts
- Message flow rates through Kafka topics
- Error rates and response times across services

**Business Metrics:**
- Signal generation trends by type and cryptocurrency
- Alert delivery patterns and effectiveness
- System uptime and availability metrics

### 7.6 Alerting Strategy

**Critical Alerts:**
- Service unavailability (all replicas down)
- Kafka cluster issues (broker failures, topic unavailable)
- External API rate limit exceeded

**Warning Alerts:**
- High error rates (>5% over 10 minutes)
- Resource utilization approaching limits (>80%)
- Message processing lag increasing
