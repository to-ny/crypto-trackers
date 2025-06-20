package kafka

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockSignalProducer struct {
	signals   []*TradingSignal
	shouldErr bool
}

func (m *mockSignalProducer) PublishSignal(ctx context.Context, topic string, signal *TradingSignal) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if m.shouldErr {
		return errors.New("mock error")
	}
	m.signals = append(m.signals, signal)
	return nil
}

func (m *mockSignalProducer) Close() error {
	return nil
}

func TestProducer_PublishSignal(t *testing.T) {
	signal := &TradingSignal{
		Timestamp:      time.Now(),
		Symbol:         "BTC",
		SignalType:     "test",
		SignalStrength: "strong",
		Direction:      "bullish",
		Details:        map[string]interface{}{"test": "value"},
		ServiceID:      "test-service",
	}

	t.Run("successful publish", func(t *testing.T) {
		producer := &mockSignalProducer{}
		ctx := context.Background()

		err := producer.PublishSignal(ctx, "test-topic", signal)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(producer.signals) != 1 {
			t.Errorf("expected 1 signal, got %d", len(producer.signals))
		}

		if producer.signals[0].Symbol != "BTC" {
			t.Errorf("expected symbol BTC, got %s", producer.signals[0].Symbol)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		producer := &mockSignalProducer{}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		cancel()

		err := producer.PublishSignal(ctx, "test-topic", signal)
		if err == nil {
			t.Error("expected context timeout error")
		}
	})

	t.Run("producer error", func(t *testing.T) {
		producer := &mockSignalProducer{shouldErr: true}
		ctx := context.Background()

		err := producer.PublishSignal(ctx, "test-topic", signal)
		if err == nil {
			t.Error("expected producer error")
		}
	})
}