package alerts

import (
	"alert-service/internal/kafka"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestAlertProcessor_EndToEndIntegration(t *testing.T) {
	alertsReceived := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "test_alerts_received", Help: "test"},
		[]string{"symbol", "signal_type"},
	)
	alertsSent := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "test_alerts_sent", Help: "test"},
		[]string{"symbol"},
	)
	alertsRateLimited := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "test_alerts_rate_limited", Help: "test"},
		[]string{"symbol"},
	)
	
	processor := NewAlertProcessor(5, *alertsReceived, *alertsSent, *alertsRateLimited)

	t.Run("rate limiting integration", func(t *testing.T) {
		signal := &kafka.TradingSignal{
			Timestamp:      time.Now(),
			Symbol:         "BTC",
			SignalType:     "moving_average_crossover",
			SignalStrength: "strong",
			Direction:      "bullish",
			Details:        map[string]interface{}{"crossover_type": "golden_cross"},
			ServiceID:      "ma-signal-detector",
		}

		err := processor.ProcessSignal(signal)
		if err != nil {
			t.Errorf("expected no error on first signal, got %v", err)
		}

		err = processor.ProcessSignal(signal)
		if err != nil {
			t.Errorf("expected no error on rate limited signal, got %v", err)
		}
	})

	t.Run("multiple symbols processing", func(t *testing.T) {
		symbols := []string{"BTC", "ETH", "ADA", "SOL"}
		
		for _, symbol := range symbols {
			signal := &kafka.TradingSignal{
				Timestamp:      time.Now(),
				Symbol:         symbol,
				SignalType:     "volume_spike",
				SignalStrength: "medium",
				Direction:      "neutral",
				Details:        map[string]interface{}{"spike_multiplier": 1.8},
				ServiceID:      "volume-spike-detector",
			}

			err := processor.ProcessSignal(signal)
			if err != nil {
				t.Errorf("expected no error for symbol %s, got %v", symbol, err)
			}
		}
	})

	t.Run("signal type variety handling", func(t *testing.T) {
		signalTypes := []struct {
			signalType string
			direction  string
			strength   string
		}{
			{"moving_average_crossover", "bullish", "strong"},
			{"moving_average_crossover", "bearish", "medium"},
			{"volume_spike", "neutral", "weak"},
		}

		for i, st := range signalTypes {
			signal := &kafka.TradingSignal{
				Timestamp:      time.Now(),
				Symbol:         "TEST",
				SignalType:     st.signalType,
				SignalStrength: st.strength,
				Direction:      st.direction,
				Details:        map[string]interface{}{"test_id": i},
				ServiceID:      "test-service",
			}

			err := processor.ProcessSignal(signal)
			if err != nil {
				t.Errorf("expected no error for signal type %s, got %v", st.signalType, err)
			}
		}
	})

	t.Run("cooldown period verification", func(t *testing.T) {
		shortProcessor := NewAlertProcessor(0, *alertsReceived, *alertsSent, *alertsRateLimited)

		signal := &kafka.TradingSignal{
			Timestamp:      time.Now(),
			Symbol:         "COOLDOWN",
			SignalType:     "test",
			SignalStrength: "medium",
			Direction:      "neutral",
			ServiceID:      "test",
		}

		err := shortProcessor.ProcessSignal(signal)
		if err != nil {
			t.Errorf("expected no error on first signal, got %v", err)
		}

		time.Sleep(1 * time.Millisecond)

		err = shortProcessor.ProcessSignal(signal)
		if err != nil {
			t.Errorf("expected no error after cooldown, got %v", err)
		}
	})
}