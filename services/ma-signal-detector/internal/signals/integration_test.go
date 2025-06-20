package signals

import (
	"ma-signal-detector/internal/kafka"
	"testing"
	"time"
)

func TestMADetector_EndToEndIntegration(t *testing.T) {
	producer := &mockProducer{}
	detector := NewMADetector(producer)

	t.Run("golden cross signal generation", func(t *testing.T) {
		producer.signals = nil
		goldDetector := NewMADetector(producer)

		basePrice := 50000.0
		
		for i := 0; i < SMA50Period; i++ {
			price := basePrice
			event := &kafka.PriceEvent{
				Timestamp: time.Now().Add(time.Duration(i) * time.Minute),
				Symbol:    "GOLD",
				PriceUSD:  price,
			}
			goldDetector.ProcessPriceEvent(event)
		}

		for i := 0; i < SMA20Period+5; i++ {
			price := basePrice + float64(i+1)*100
			event := &kafka.PriceEvent{
				Timestamp: time.Now().Add(time.Duration(SMA50Period+i) * time.Minute),
				Symbol:    "GOLD",
				PriceUSD:  price,
			}
			goldDetector.ProcessPriceEvent(event)
		}

		if len(producer.signals) == 0 {
			t.Fatal("expected at least one signal")
		}

		signal := producer.signals[0]
		if signal.Direction != "bullish" {
			t.Errorf("expected bullish direction for golden cross, got %s", signal.Direction)
		}

		details := signal.Details
		crossoverType, ok := details["crossover_type"].(string)
		if !ok || crossoverType != "golden_cross" {
			t.Errorf("expected golden_cross, got %v", crossoverType)
		}
	})

	t.Run("duplicate signal prevention", func(t *testing.T) {
		initialSignalCount := len(producer.signals)

		event := &kafka.PriceEvent{
			Timestamp: time.Now(),
			Symbol:    "BTC",
			PriceUSD:  50000.0,
		}
		detector.ProcessPriceEvent(event)

		if len(producer.signals) != initialSignalCount {
			t.Error("duplicate signal should not be generated")
		}
	})

	t.Run("concurrent access safety", func(t *testing.T) {
		done := make(chan bool)

		for i := 0; i < 10; i++ {
			go func(i int) {
				event := &kafka.PriceEvent{
					Timestamp: time.Now(),
					Symbol:    "CONCURRENT",
					PriceUSD:  float64(50000 + i),
				}
				detector.ProcessPriceEvent(event)
				done <- true
			}(i)
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		history := detector.priceHistory["CONCURRENT"]
		if len(history.Prices) != 10 {
			t.Errorf("expected 10 prices after concurrent access, got %d", len(history.Prices))
		}
	})
}