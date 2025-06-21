package signals

import (
	"testing"
	"time"
	"volume-spike-detector/internal/kafka"

	"github.com/prometheus/client_golang/prometheus"
)

func TestVolumeDetector_EndToEndIntegration(t *testing.T) {
	producer := &mockProducer{}
	threshold := 1.5
	
	eventsProcessed := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "test_events_processed", Help: "test"},
		[]string{"symbol"},
	)
	spikesDetected := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "test_spikes_detected", Help: "test"},
		[]string{"symbol"},
	)
	processingTime := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "test_processing_time", Help: "test"},
		[]string{"symbol"},
	)
	
	detector := NewVolumeDetector(producer, threshold, *eventsProcessed, *spikesDetected, *processingTime)

	baseVolume := 1000000000.0

	t.Run("medium spike detection", func(t *testing.T) {
		producer.signals = nil
		mediumDetector := NewVolumeDetector(producer, threshold, *eventsProcessed, *spikesDetected, *processingTime)

		for i := 0; i < 5; i++ {
			event := &kafka.PriceEvent{
				Timestamp: time.Now().Add(time.Duration(i) * time.Hour),
				Symbol:    "MEDIUM",
				Volume24h: baseVolume,
			}
			mediumDetector.ProcessPriceEvent(event)
		}

		spikeEvent := &kafka.PriceEvent{
			Timestamp: time.Now().Add(6 * time.Hour),
			Symbol:    "MEDIUM",
			Volume24h: baseVolume * 1.8,
		}
		mediumDetector.ProcessPriceEvent(spikeEvent)

		if len(producer.signals) == 0 {
			t.Fatal("expected volume spike signal")
		}

		signal := producer.signals[0]
		if signal.SignalStrength != "medium" {
			t.Errorf("expected strength medium, got %s", signal.SignalStrength)
		}
	})

	t.Run("strong spike detection", func(t *testing.T) {
		producer.signals = nil
		strongDetector := NewVolumeDetector(producer, threshold, *eventsProcessed, *spikesDetected, *processingTime)

		for i := 0; i < 5; i++ {
			event := &kafka.PriceEvent{
				Timestamp: time.Now().Add(time.Duration(i) * time.Hour),
				Symbol:    "STRONG",
				Volume24h: baseVolume,
			}
			strongDetector.ProcessPriceEvent(event)
		}

		spikeEvent := &kafka.PriceEvent{
			Timestamp: time.Now().Add(6 * time.Hour),
			Symbol:    "STRONG",
			Volume24h: baseVolume * 2.5,
		}
		strongDetector.ProcessPriceEvent(spikeEvent)

		if len(producer.signals) == 0 {
			t.Fatal("expected volume spike signal")
		}

		signal := producer.signals[0]
		if signal.SignalStrength != "strong" {
			t.Errorf("expected strength strong, got %s", signal.SignalStrength)
		}
	})

	t.Run("historical data cleanup", func(t *testing.T) {
		now := time.Now()

		for i := 0; i < VolumeDays*2; i++ {
			event := &kafka.PriceEvent{
				Timestamp: now.Add(time.Duration(i-VolumeDays*2) * 24 * time.Hour),
				Symbol:    "CLEANUP",
				Volume24h: 1000.0,
			}
			detector.ProcessPriceEvent(event)
		}

		history := detector.volumeHistory["CLEANUP"]
		if len(history.Volumes) > VolumeDays {
			t.Errorf("expected history to be cleaned up to max %d days, got %d entries", VolumeDays, len(history.Volumes))
		}

		oldestTime := history.Timestamps[0]
		cutoff := now.Add(-time.Duration(VolumeDays) * 24 * time.Hour)
		if oldestTime.Before(cutoff) {
			t.Errorf("expected oldest timestamp to be after cutoff %v, got %v", cutoff, oldestTime)
		}
	})

	t.Run("concurrent volume processing", func(t *testing.T) {
		done := make(chan bool)

		for i := 0; i < 10; i++ {
			go func(i int) {
				event := &kafka.PriceEvent{
					Timestamp: time.Now().Add(time.Duration(i) * time.Hour),
					Symbol:    "CONCURRENT",
					Volume24h: float64(1000000000 + i*100000000),
				}
				detector.ProcessPriceEvent(event)
				done <- true
			}(i)
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		history := detector.volumeHistory["CONCURRENT"]
		if len(history.Volumes) != 10 {
			t.Errorf("expected 10 volumes after concurrent access, got %d", len(history.Volumes))
		}
	})

	t.Run("no spike below threshold", func(t *testing.T) {
		producer.signals = nil

		for i := 0; i < 5; i++ {
			event := &kafka.PriceEvent{
				Timestamp: time.Now().Add(time.Duration(i) * time.Hour),
				Symbol:    "NOSPIKE",
				Volume24h: baseVolume,
			}
			detector.ProcessPriceEvent(event)
		}

		noSpikeEvent := &kafka.PriceEvent{
			Timestamp: time.Now().Add(6 * time.Hour),
			Symbol:    "NOSPIKE",
			Volume24h: baseVolume * 1.2,
		}
		detector.ProcessPriceEvent(noSpikeEvent)

		if len(producer.signals) != 0 {
			t.Errorf("expected no signals below threshold, got %d", len(producer.signals))
		}
	})
}
