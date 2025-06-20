package signals

import (
	"log"
	"ma-signal-detector/internal/kafka"
	"sync"
	"time"
)

const (
	MaxHistorySize = 100
	MinSignalSize  = 50
	SMA20Period    = 20
	SMA50Period    = 50
)

type PriceHistory struct {
	Prices []float64
	mutex  sync.RWMutex
}

type MADetector struct {
	priceHistory map[string]*PriceHistory
	lastSignals  map[string]string
	mutex        sync.RWMutex
	producer     *kafka.Producer
}

func NewMADetector(producer *kafka.Producer) *MADetector {
	return &MADetector{
		priceHistory: make(map[string]*PriceHistory),
		lastSignals:  make(map[string]string),
		producer:     producer,
	}
}

func (ma *MADetector) ProcessPriceEvent(event *kafka.PriceEvent) error {
	ma.mutex.Lock()
	defer ma.mutex.Unlock()

	history, exists := ma.priceHistory[event.Symbol]
	if !exists {
		history = &PriceHistory{
			Prices: make([]float64, 0, MaxHistorySize),
		}
		ma.priceHistory[event.Symbol] = history
		log.Printf("Started tracking price history for %s", event.Symbol)
	}

	history.mutex.Lock()
	history.Prices = append(history.Prices, event.PriceUSD)
	if len(history.Prices) > MaxHistorySize {
		history.Prices = history.Prices[1:]
	}
	priceCount := len(history.Prices)
	history.mutex.Unlock()

	log.Printf("Processed price event for %s: $%.2f (history: %d points)", event.Symbol, event.PriceUSD, priceCount)

	if priceCount >= MinSignalSize {
		return ma.checkForCrossover(event.Symbol, event.Timestamp)
	}

	return nil
}

func (ma *MADetector) checkForCrossover(symbol string, timestamp time.Time) error {
	history := ma.priceHistory[symbol]
	history.mutex.RLock()
	prices := make([]float64, len(history.Prices))
	copy(prices, history.Prices)
	history.mutex.RUnlock()

	if len(prices) < SMA50Period {
		return nil
	}

	currentSMA20 := calculateSMA(prices, SMA20Period)
	currentSMA50 := calculateSMA(prices, SMA50Period)

	if len(prices) < SMA50Period+1 {
		return nil
	}

	prevPrices := prices[:len(prices)-1]
	prevSMA20 := calculateSMA(prevPrices, SMA20Period)
	prevSMA50 := calculateSMA(prevPrices, SMA50Period)

	var signalType string
	var direction string

	if prevSMA20 <= prevSMA50 && currentSMA20 > currentSMA50 {
		signalType = "golden_cross"
		direction = "bullish"
	} else if prevSMA20 >= prevSMA50 && currentSMA20 < currentSMA50 {
		signalType = "death_cross"
		direction = "bearish"
	}

	if signalType != "" && ma.lastSignals[symbol] != signalType {
		ma.lastSignals[symbol] = signalType
		return ma.publishSignal(symbol, timestamp, signalType, direction, currentSMA20, currentSMA50)
	}

	return nil
}

func (ma *MADetector) publishSignal(symbol string, timestamp time.Time, crossoverType, direction string, sma20, sma50 float64) error {
	signal := &kafka.TradingSignal{
		Timestamp:      timestamp,
		Symbol:         symbol,
		SignalType:     "moving_average_crossover",
		SignalStrength: "strong",
		Direction:      direction,
		Details: map[string]interface{}{
			"sma_20":         sma20,
			"sma_50":         sma50,
			"crossover_type": crossoverType,
		},
		ServiceID: "ma-detector-v1",
	}

	if err := ma.producer.PublishSignal("trading-signals", signal); err != nil {
		log.Printf("Failed to publish signal for %s: %v", symbol, err)
		return err
	}

	log.Printf("Published %s signal for %s (SMA20: %.2f, SMA50: %.2f)", crossoverType, symbol, sma20, sma50)
	return nil
}

func calculateSMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	sum := 0.0
	start := len(prices) - period
	for i := start; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(period)
}