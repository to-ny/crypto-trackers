# Crypto Trading Signal Detector - Implementation Roadmap

## Implementation Order

This roadmap organizes features by implementation priority and logical dependencies. Each section can be tackled independently with coding agents, building incrementally toward a complete system.

## 1. Kubernetes Infrastructure

### 1.1 Basic Cluster Setup
- [ ] Namespace creation (`crypto-signals`)
- [ ] ConfigMaps for application configuration
- [ ] Secrets for API keys (even if unused initially)

### 1.2 Kafka Deployment
- [ ] Kafka StatefulSet with 3 brokers
- [ ] Persistent volumes for Kafka data
- [ ] Zookeeper coordination service
- [ ] Anti-affinity rules for broker distribution

### 1.3 Topic Configuration
- [ ] Create `crypto-prices` topic (2 partitions, replication factor 3)
- [ ] Create `trading-signals` topic (2 partitions, replication factor 3)
- [ ] Verify topic creation and partition assignment

## 2. Data Ingestion Service

### 2.1 Core API Integration
- [ ] CoinGecko API client implementation
- [ ] Support for BTC and ETH data fetching
- [ ] JSON response parsing and validation
- [ ] Error handling for API failures

### 2.2 Data Transformation
- [ ] Convert API response to standardized price event schema
- [ ] Add timestamp and source metadata
- [ ] Validate required fields before publishing

### 2.3 Kafka Producer
- [ ] Kafka producer client configuration
- [ ] Publish price events to `crypto-prices` topic
- [ ] Handle Kafka connection failures and retries
- [ ] Partitioning by cryptocurrency symbol

### 2.4 Service Operations
- [ ] Configurable polling interval (default 60 seconds)
- [ ] Health check endpoint (`/health`)
- [ ] Readiness check endpoint (`/ready`)
- [ ] Graceful shutdown handling

### 2.5 Containerization and Deployment
- [ ] Dockerfile with Python runtime
- [ ] Kubernetes Deployment manifest
- [ ] Service configuration for health checks
- [ ] Environment variable configuration

## 3. Moving Average Signal Detection Service

### 3.1 Kafka Consumer Setup
- [ ] Consumer client for `crypto-prices` topic
- [ ] Consumer group configuration
- [ ] Message deserialization and validation
- [ ] Offset management and error handling

### 3.2 Price History Management
- [ ] In-memory storage for price history (last 100 data points)
- [ ] Circular buffer implementation for memory efficiency
- [ ] Per-symbol price tracking
- [ ] Data structure thread safety

### 3.3 Moving Average Calculation
- [ ] Simple Moving Average (SMA) calculation logic
- [ ] SMA 20 and SMA 50 implementation
- [ ] Minimum data points validation (50 required)
- [ ] Handle insufficient data gracefully

### 3.4 Crossover Detection
- [ ] Golden cross detection (SMA 20 > SMA 50)
- [ ] Death cross detection (SMA 20 < SMA 50)
- [ ] Signal strength calculation based on crossover magnitude
- [ ] Prevent duplicate signals for same crossover

### 3.5 Signal Publishing
- [ ] Trading signal event creation
- [ ] Kafka producer for `trading-signals` topic
- [ ] Signal metadata and details formatting
- [ ] Publishing error handling

### 3.6 Service Infrastructure
- [ ] Health and readiness endpoints
- [ ] Structured logging with correlation IDs
- [ ] Graceful shutdown and resource cleanup
- [ ] Dockerfile and Kubernetes deployment

## 4. Volume Spike Detection Service

### 4.1 Volume History Tracking
- [ ] 7-day rolling volume history per cryptocurrency
- [ ] Circular buffer for volume data storage
- [ ] Average volume calculation over rolling window
- [ ] Memory-efficient data structures

### 4.2 Spike Detection Logic
- [ ] Configurable spike threshold (default 1.3x average)
- [ ] Current volume vs. average comparison
- [ ] Spike magnitude calculation
- [ ] Signal strength classification (weak, moderate, strong)

### 4.3 Signal Generation
- [ ] Volume spike signal event creation
- [ ] Signal details with volume metrics
- [ ] Direction determination (bullish volume spike)
- [ ] Publishing to `trading-signals` topic

### 4.4 Service Operations
- [ ] Kafka consumer for price events
- [ ] Health monitoring endpoints
- [ ] Error handling and recovery
- [ ] Deployment configuration

## 5. Alert Service

### 5.1 Signal Consumption
- [ ] Consumer for `trading-signals` topic
- [ ] Multi-signal type handling
- [ ] Message deserialization and routing
- [ ] Consumer group management

### 5.2 Rate Limiting
- [ ] Per-symbol alert rate limiting (1 per 5 minutes)
- [ ] Configurable cooldown periods
- [ ] In-memory rate limit state tracking
- [ ] Signal aggregation logic

### 5.3 Alert Generation
- [ ] Structured alert formatting
- [ ] Console output for development
- [ ] JSON log output for production
- [ ] Alert metadata and context

### 5.4 Service Infrastructure
- [ ] Health and readiness checks
- [ ] Graceful shutdown handling
- [ ] Kubernetes deployment
- [ ] Resource configuration

## 6. System Integration and Testing

### 6.1 End-to-End Validation
- [ ] Complete data flow verification
- [ ] Signal generation from live market data
- [ ] Alert delivery confirmation
- [ ] Performance and latency testing

### 6.2 Error Scenarios
- [ ] API failure handling
- [ ] Kafka broker failures
- [ ] Service restart recovery
- [ ] Invalid data handling

### 6.3 Monitoring and Observability
- [ ] Health check validation across all services
- [ ] Log aggregation and correlation
- [ ] Resource utilization monitoring
- [ ] Error rate tracking

## Implementation Guidelines

### Service Development Pattern
1. **Core Logic First:** Implement business logic without Kafka
2. **Add Kafka Integration:** Connect to message streams
3. **Health Checks:** Add monitoring endpoints
4. **Containerize:** Create Docker images
5. **Deploy:** Add Kubernetes manifests

### Coding Agent Sessions
- **Single Feature Focus:** Tackle one checkbox item per session
- **Reference Documentation:** Always link back to DESIGN.md schemas
- **Incremental Testing:** Verify functionality before moving to next item
- **Clean Commits:** Each completed feature should be a separate commit

### Dependencies
- Kafka infrastructure must be ready before any service development
- Data Ingestion Service must be working before signal detection
- Signal detection services must work before Alert Service
- Within each service, core logic before Kafka integration
