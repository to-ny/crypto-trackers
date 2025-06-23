# Crypto Trading Signal Detector - Software Design Document

## 1. Project Overview

The Crypto Trading Signal Detector is a real-time distributed system that monitors cryptocurrency market data and generates trading signals based on technical analysis indicators.

### Objectives
- Practice distributed systems architecture and event-driven design
- Build a scalable system for cryptocurrency signal detection
- Implement multi-language microservices with Kafka integration

### Success Criteria
- Detect and alert on moving average crossovers and volume spikes for BTC and ETH
- Add new cryptocurrencies through configuration only
- Full Kubernetes deployment with Helm charts

## 2. System Architecture

### High-Level Design
```
CoinGecko API → Data Ingestion → Kafka → Signal Detection → Alerts
```

### Core Components

**Data Ingestion Service (Python)**
- Polls CoinGecko API every 60 seconds for BTC/ETH data
- Publishes standardized price events to Kafka

**Signal Detection Services (Go)**
- Moving Average Service: Detects SMA 20/50 crossovers
- Volume Spike Service: Identifies volume above 7-day average threshold

**Alert Service (Go)**
- Consumes trading signals and generates notifications
- Rate-limited output (1 alert per symbol per 5 minutes)

**Kafka**
- Event streaming backbone with two topics:
  - `crypto-prices`: Raw market data
  - `trading-signals`: Generated trading signals

## 3. Data Models

### Price Event (crypto-prices topic)
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

### Trading Signal (trading-signals topic)
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

For volume spikes, details contains:
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

## 4. Service Specifications

### Data Ingestion Service
- **Language**: Python
- **Function**: Fetch data from CoinGecko API and publish to Kafka
- **Configuration**: Polling interval, supported symbols
- **Error Handling**: Continue on API failures, retry Kafka publishing

### Moving Average Service
- **Language**: Go
- **Function**: Calculate SMA 20/50 and detect crossovers
- **State**: In-memory price history (last 100 points per symbol)
- **Signals**: Golden cross (bullish), Death cross (bearish)
- **Requirement**: Minimum 50 data points before generating signals

### Volume Spike Service
- **Language**: Go
- **Function**: Detect volume spikes above 7-day average
- **State**: Rolling 7-day volume history per symbol
- **Threshold**: Configurable (default 1.3x average volume)
- **Signal Strength**: Based on spike magnitude

### Alert Service
- **Language**: Go
- **Function**: Consume signals and generate notifications
- **Output**: Console logs and structured JSON
- **Rate Limiting**: Per-symbol cooldown periods

## 5. Technology Stack

- **Languages**: Python (data processing), Go (signal processing)
- **Message Broker**: Apache Kafka
- **Deployment**: Kubernetes with Helm charts
- **External API**: CoinGecko (free tier)
- **Containerization**: Docker

## 6. Deployment

- **Kubernetes**: All services deployed as containers
- **Helm**: Template-based deployment and configuration
- **Namespace**: `crypto-trackers`
- **Health Checks**: `/health` and `/ready` endpoints for all services
- **Scaling**: Horizontal scaling for signal detection services

## 7. Monitoring

### Health Monitoring
- Kubernetes liveness and readiness probes
- Service health endpoints

### Logging
- Structured JSON logs to stdout
- Standard fields: timestamp, level, service, message
- Log levels: INFO, WARN, ERROR, DEBUG

### Key Metrics
- API call success rates
- Signal generation frequency
- Message processing rates
- Error rates per service
