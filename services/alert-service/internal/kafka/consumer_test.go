package kafka

import (
	"errors"
	"testing"
	"time"
)

type mockSignalHandler struct {
	signals   []*TradingSignal
	shouldErr bool
}

func (m *mockSignalHandler) handleSignal(signal *TradingSignal) error {
	if m.shouldErr {
		return errors.New("mock error")
	}
	m.signals = append(m.signals, signal)
	return nil
}

func TestTradingSignal_Structure(t *testing.T) {
	signal := &TradingSignal{
		Timestamp:      time.Now(),
		Symbol:         "BTC",
		SignalType:     "test",
		SignalStrength: "strong",
		Direction:      "bullish",
		Details:        map[string]interface{}{"test": "value"},
		ServiceID:      "test-service",
	}

	if signal.Symbol != "BTC" {
		t.Errorf("expected symbol BTC, got %s", signal.Symbol)
	}

	if signal.SignalType != "test" {
		t.Errorf("expected signal type test, got %s", signal.SignalType)
	}

	if signal.Details["test"] != "value" {
		t.Errorf("expected details test value, got %v", signal.Details["test"])
	}
}

func TestConsumerGroupHandler_ProcessMessage(t *testing.T) {
	handler := &mockSignalHandler{}

	t.Run("successful signal processing", func(t *testing.T) {
		signal := &TradingSignal{
			Timestamp:      time.Now(),
			Symbol:         "BTC",
			SignalType:     "test",
			SignalStrength: "strong",
			Direction:      "bullish",
			Details:        map[string]interface{}{},
			ServiceID:      "test",
		}

		if len(handler.signals) != 0 {
			t.Errorf("expected 0 signals initially, got %d", len(handler.signals))
		}

		err := handler.handleSignal(signal)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(handler.signals) != 1 {
			t.Errorf("expected 1 signal, got %d", len(handler.signals))
		}

		if handler.signals[0].Symbol != "BTC" {
			t.Errorf("expected symbol BTC, got %s", handler.signals[0].Symbol)
		}
	})

	t.Run("handler error handling", func(t *testing.T) {
		errorHandler := &mockSignalHandler{shouldErr: true}

		signal := &TradingSignal{
			Symbol:    "ETH",
			ServiceID: "test",
		}

		err := errorHandler.handleSignal(signal)
		if err == nil {
			t.Error("expected error from handler")
		}

		if len(errorHandler.signals) != 0 {
			t.Errorf("expected no signals on error, got %d", len(errorHandler.signals))
		}
	})
}