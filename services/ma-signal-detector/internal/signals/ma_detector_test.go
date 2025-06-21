package signals

import (
	"context"
	"ma-signal-detector/internal/kafka"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type mockProducer struct {
	signals   []*kafka.TradingSignal
	shouldErr bool
}

func (m *mockProducer) PublishSignal(ctx context.Context, topic string, signal *kafka.TradingSignal) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if m.shouldErr {
		return context.DeadlineExceeded
	}
	m.signals = append(m.signals, signal)
	return nil
}

func (m *mockProducer) Close() error {
	return nil
}

func TestMADetector_ProcessPriceEvent(t *testing.T) {
	producer := &mockProducer{}

	priceEventsProcessed := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "test_price_events_processed", Help: "test"},
		[]string{"symbol"},
	)
	signalsGenerated := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "test_signals_generated", Help: "test"},
		[]string{"symbol", "signal_type"},
	)
	processingTime := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "test_processing_time", Help: "test"},
		[]string{"symbol"},
	)

	detector := NewMADetector(producer, *priceEventsProcessed, *signalsGenerated, *processingTime)

	t.Run("first price event creates history", func(t *testing.T) {
		event := &kafka.PriceEvent{
			Timestamp: time.Now(),
			Symbol:    "BTC",
			PriceUSD:  50000.0,
		}

		err := detector.ProcessPriceEvent(event)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(detector.priceHistory) != 1 {
			t.Errorf("expected 1 symbol in history, got %d", len(detector.priceHistory))
		}

		history := detector.priceHistory["BTC"]
		if len(history.Prices) != 1 {
			t.Errorf("expected 1 price in history, got %d", len(history.Prices))
		}

		if history.Prices[0] != 50000.0 {
			t.Errorf("expected price 50000.0, got %f", history.Prices[0])
		}
	})

	t.Run("insufficient data points no signal", func(t *testing.T) {
		for i := 0; i < MinSignalSize-1; i++ {
			event := &kafka.PriceEvent{
				Timestamp: time.Now(),
				Symbol:    "ETH",
				PriceUSD:  float64(3000 + i),
			}
			detector.ProcessPriceEvent(event)
		}

		if len(producer.signals) != 0 {
			t.Errorf("expected no signals with insufficient data, got %d", len(producer.signals))
		}
	})
}

func TestCalculateSMA(t *testing.T) {
	tests := []struct {
		name     string
		prices   []float64
		period   int
		expected float64
	}{
		{
			name:     "empty prices",
			prices:   []float64{},
			period:   20,
			expected: 0,
		},
		{
			name:     "insufficient data",
			prices:   []float64{100, 110, 120},
			period:   5,
			expected: 0,
		},
		{
			name:     "exact period",
			prices:   []float64{100, 110, 120, 130, 140},
			period:   5,
			expected: 120,
		},
		{
			name:     "more than period",
			prices:   []float64{100, 110, 120, 130, 140, 150},
			period:   3,
			expected: 140,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateSMA(tt.prices, tt.period)
			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestMADetector_CrossoverDetection(t *testing.T) {
	producer := &mockProducer{}

	priceEventsProcessed := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "test_price_events_processed", Help: "test"},
		[]string{"symbol"},
	)
	signalsGenerated := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "test_signals_generated", Help: "test"},
		[]string{"symbol", "signal_type"},
	)
	processingTime := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "test_processing_time", Help: "test"},
		[]string{"symbol"},
	)

	detector := NewMADetector(producer, *priceEventsProcessed, *signalsGenerated, *processingTime)

	prices := make([]float64, SMA50Period+5)
	for i := 0; i < SMA50Period; i++ {
		prices[i] = 100.0
	}
	for i := SMA50Period; i < len(prices); i++ {
		prices[i] = 110.0
	}

	for i, price := range prices {
		event := &kafka.PriceEvent{
			Timestamp: time.Now().Add(time.Duration(i) * time.Minute),
			Symbol:    "TEST",
			PriceUSD:  price,
		}
		detector.ProcessPriceEvent(event)
	}

	if len(producer.signals) == 0 {
		t.Error("expected at least one signal to be generated")
	}

	signal := producer.signals[0]
	if signal.SignalType != "moving_average_crossover" {
		t.Errorf("expected signal type 'moving_average_crossover', got %s", signal.SignalType)
	}

	if signal.Symbol != "TEST" {
		t.Errorf("expected symbol 'TEST', got %s", signal.Symbol)
	}
}
