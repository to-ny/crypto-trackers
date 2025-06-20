# Crypto Trading Signal Detector - Implementation Roadmap

## Implementation Guidelines

This roadmap is designed for AI coding agents. Each task is self-contained and includes clear acceptance criteria.

**Before starting any task:**
1. Review DESIGN.md for data schemas and service specifications
2. Follow PROJECT_STRUCTURE.md for organization
3. Implement the specific requirements in the task description
4. Test the functionality before marking complete

## Phase 1: Foundation Infrastructure

### Task 1.1: Helm Chart Foundation
**Objective:** Create the basic Helm chart structure and namespace setup.

**Requirements:**
- Create `helm/crypto-trackers/Chart.yaml` with chart metadata
- Create `helm/crypto-trackers/values.yaml` with default configuration
- Create `helm/crypto-trackers/templates/namespace.yaml` for `crypto-trackers` namespace
- Create `helm/crypto-trackers/templates/configmap.yaml` for shared configuration

**Acceptance Criteria:**
- Helm chart can be installed with `helm install crypto-trackers ./helm/crypto-trackers`
- Namespace `crypto-trackers` is created
- ConfigMap contains placeholder values for Kafka configuration

### Task 1.2: Kafka Infrastructure
**Objective:** Deploy Kafka cluster via Helm templates.

**Requirements:**
- Create StatefulSet for Kafka brokers in `helm/crypto-trackers/templates/kafka/`
- Create Service for Kafka broker discovery
- Create Job to initialize topics: `crypto-prices` and `trading-signals`
- Configure appropriate resource requests and storage

**Acceptance Criteria:**
- Kafka cluster deploys successfully
- Topics are created automatically
- Services can connect to Kafka via service DNS names

## Phase 2: Data Ingestion Service

### Task 2.1: Python Service Structure
**Objective:** Create the data ingestion service foundation.

**Requirements:**
- Create `services/data-ingestion/` directory structure
- Create `Dockerfile` with Python 3.11 base image
- Create `requirements.txt` with necessary dependencies: `requests`, `kafka-python`
- Create basic `src/main.py` with health check endpoints (`/health`, `/ready`)

**Acceptance Criteria:**
- Service starts and responds to health check endpoints
- Docker image builds successfully
- Service can be deployed to Kubernetes (without functionality)

### Task 2.2: CoinGecko API Integration
**Objective:** Implement CoinGecko API client and data fetching.

**Requirements:**
- Create `src/api/coingecko.py` with API client class
- Implement method to fetch BTC and ETH price data
- Handle API errors gracefully (timeouts, rate limits, invalid responses)
- Follow the Price Event schema from DESIGN.md exactly

**API Endpoint:** Use CoinGecko's `/simple/price` endpoint
**Error Handling:** Log errors and return None on failure

**Acceptance Criteria:**
- API client successfully fetches BTC and ETH data
- Returns data in exact Price Event schema format
- Handles network errors without crashing

### Task 2.3: Kafka Producer Integration
**Objective:** Publish price data to Kafka topic.

**Requirements:**
- Create `src/kafka/producer.py` with Kafka producer client
- Implement method to publish Price Events to `crypto-prices` topic
- Add retry logic for Kafka connection failures
- Integrate with main application loop (60-second polling)

**Acceptance Criteria:**
- Price events are published to `crypto-prices` topic in correct JSON format
- Service handles Kafka connection failures gracefully
- Messages are properly partitioned by symbol

### Task 2.4: Service Deployment
**Objective:** Create Helm deployment for data ingestion service.

**Requirements:**
- Create Helm deployment template in `helm/crypto-trackers/templates/services/data-ingestion/`
- Include Deployment, Service, and ConfigMap resources
- Configure environment variables for API polling interval and Kafka settings
- Add proper health checks and resource requests

**Acceptance Criteria:**
- Service deploys successfully via Helm
- Health checks pass in Kubernetes
- Service logs show successful API calls and Kafka publishing

## Phase 3: Moving Average Signal Detection

### Task 3.1: Go Service Foundation
**Objective:** Create the moving average detection service structure.

**Requirements:**
- Create `services/ma-signal-detector/` directory with Go module
- Create `cmd/main.go` with basic HTTP server and health endpoints
- Create `internal/config/` for configuration management
- Create `Dockerfile` with Go build process

**Acceptance Criteria:**
- Go service compiles and runs locally
- Health endpoints (`/health`, `/ready`) respond correctly
- Docker image builds successfully

### Task 3.2: Kafka Consumer Setup
**Objective:** Consume price events from Kafka.

**Requirements:**
- Create `internal/kafka/consumer.go` with Kafka consumer client
- Subscribe to `crypto-prices` topic
- Deserialize Price Event JSON messages
- Implement proper error handling and offset management

**Acceptance Criteria:**
- Consumer successfully reads messages from `crypto-prices` topic
- Messages are properly deserialized to Go structs
- Consumer handles Kafka connection issues gracefully

### Task 3.3: Moving Average Calculation
**Objective:** Implement SMA calculation and crossover detection.

**Requirements:**
- Create `internal/signals/` package for signal logic
- Implement in-memory price history storage (last 100 points per symbol)
- Calculate SMA 20 and SMA 50 for each price update
- Detect golden cross and death cross patterns
- Generate Trading Signal events following DESIGN.md schema

**Signal Logic:**
- Golden Cross: SMA 20 crosses above SMA 50 (bullish)
- Death Cross: SMA 20 crosses below SMA 50 (bearish)
- Require minimum 50 data points before generating signals

**Acceptance Criteria:**
- SMA calculations are mathematically correct
- Crossover detection triggers only on actual crosses (not continuous signals)
- Trading Signal events match exact schema from DESIGN.md

### Task 3.4: Signal Publishing and Deployment
**Objective:** Publish signals to Kafka and deploy service.

**Requirements:**
- Create Kafka producer to publish Trading Signals to `trading-signals` topic
- Create Helm deployment templates for the service
- Configure proper resource requests and scaling
- Integrate all components into working service

**Acceptance Criteria:**
- Trading signals are published in correct JSON format
- Service deploys successfully via Helm
- End-to-end flow works: price events → signal detection → signal publishing

## Phase 4: Volume Spike Detection

### Task 4.1: Volume Spike Service
**Objective:** Create volume spike detection service.

**Requirements:**
- Create `services/volume-spike-detector/` with same Go structure as MA service
- Implement 7-day rolling volume history per cryptocurrency
- Detect volume spikes above configurable threshold (default 1.3x average)
- Generate Trading Signal events with volume spike details

**Signal Logic:**
- Calculate 7-day average volume
- Trigger signal when current volume > threshold * average
- Include spike multiplier and volume metrics in signal details

**Acceptance Criteria:**
- Volume history is maintained correctly
- Spike detection algorithm works as specified
- Signals include proper volume metrics in details field

## Phase 5: Alert Service

### Task 5.1: Alert Service Implementation
**Objective:** Create final alert service to consume and output signals.

**Requirements:**
- Create `services/alert-service/` with Go service structure
- Consume from `trading-signals` topic
- Implement rate limiting (1 alert per symbol per 5 minutes)
- Output structured alerts to console and logs

**Rate Limiting Logic:**
- Track last alert time per symbol
- Skip alerts within cooldown period
- Log all decisions (sent/skipped) for debugging

**Acceptance Criteria:**
- Consumes all signal types correctly
- Rate limiting prevents spam
- Alert output is clearly formatted and informative

## Phase 6: System Integration

### Task 6.1: End-to-End Testing
**Objective:** Verify complete system functionality.

**Requirements:**
- Deploy all services via Helm
- Verify data flows through entire pipeline
- Test with live CoinGecko data
- Validate signal generation for both BTC and ETH

**Acceptance Criteria:**
- Complete pipeline from API → signals → alerts works
- All services remain healthy under normal operation
- Signals are generated for actual market conditions

### Task 6.2: Documentation and Cleanup
**Objective:** Finalize system documentation.

**Requirements:**
- Update README.md with deployment instructions
- Add service-specific READMEs where needed
- Document any configuration options or operational procedures
- Clean up any temporary files or unused code

**Acceptance Criteria:**
- System can be deployed by following README instructions
- All services have basic operational documentation
- Code is clean and production-ready

## Success Criteria

Upon completion of all tasks:
- [ ] Complete system detects moving average crossovers for BTC and ETH
- [ ] Volume spike detection works for both cryptocurrencies  
- [ ] Alerts are generated and properly rate-limited
- [ ] All services deploy via Helm charts
- [ ] System runs continuously without manual intervention
- [ ] Adding new cryptocurrencies requires only configuration changes