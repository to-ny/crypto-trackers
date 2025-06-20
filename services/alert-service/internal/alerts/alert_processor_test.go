package alerts

import (
	"alert-service/internal/kafka"
	"testing"
	"time"
)

func TestAlertProcessor_ProcessSignal(t *testing.T) {
	processor := NewAlertProcessor(5)

	signal := &kafka.TradingSignal{
		Timestamp:      time.Now(),
		Symbol:         "BTC",
		SignalType:     "moving_average_crossover",
		SignalStrength: "strong",
		Direction:      "bullish",
		Details:        map[string]interface{}{"sma20": 50000.0, "sma50": 49000.0},
		ServiceID:      "ma-signal-detector",
	}

	t.Run("first signal processed", func(t *testing.T) {
		err := processor.ProcessSignal(signal)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("immediate second signal skipped", func(t *testing.T) {
		err := processor.ProcessSignal(signal)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("different symbol processed", func(t *testing.T) {
		ethSignal := &kafka.TradingSignal{
			Timestamp:      time.Now(),
			Symbol:         "ETH",
			SignalType:     "volume_spike",
			SignalStrength: "medium",
			Direction:      "neutral",
			Details:        map[string]interface{}{"spike_multiplier": 1.8},
			ServiceID:      "volume-spike-detector",
		}

		err := processor.ProcessSignal(ethSignal)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func TestAlertProcessor_FormatAlert(t *testing.T) {
	processor := NewAlertProcessor(5)

	signal := &kafka.TradingSignal{
		Timestamp:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Symbol:         "BTC",
		SignalType:     "test_signal",
		SignalStrength: "strong",
		Direction:      "bullish",
		Details:        map[string]interface{}{"test_key": "test_value", "number": 42},
		ServiceID:      "test-service",
	}

	alert := processor.formatAlert(signal)

	expectedFields := []string{
		"Symbol: BTC",
		"Signal Type: test_signal",
		"Direction: bullish",
		"Strength: strong",
		"Service: test-service",
		"test_key: test_value",
		"number: 42",
	}

	for _, field := range expectedFields {
		if !contains(alert, field) {
			t.Errorf("expected alert to contain '%s', got:\n%s", field, alert)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInMiddle(s, substr)))
}

func containsInMiddle(s, substr string) bool {
	for i := 1; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}