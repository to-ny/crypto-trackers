package signals

import (
	"context"
	"testing"
	"time"
	"volume-spike-detector/internal/kafka"
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

func TestVolumeDetector_ProcessPriceEvent(t *testing.T) {
	producer := &mockProducer{}
	threshold := 1.5
	detector := NewVolumeDetector(producer, threshold)

	t.Run("first volume event creates history", func(t *testing.T) {
		event := &kafka.PriceEvent{
			Timestamp: time.Now(),
			Symbol:    "BTC",
			Volume24h: 1000000000.0,
		}

		err := detector.ProcessPriceEvent(event)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(detector.volumeHistory) != 1 {
			t.Errorf("expected 1 symbol in history, got %d", len(detector.volumeHistory))
		}

		history := detector.volumeHistory["BTC"]
		if len(history.Volumes) != 1 {
			t.Errorf("expected 1 volume in history, got %d", len(history.Volumes))
		}

		if history.Volumes[0] != 1000000000.0 {
			t.Errorf("expected volume 1000000000.0, got %f", history.Volumes[0])
		}
	})

	t.Run("insufficient data no spike detection", func(t *testing.T) {
		event := &kafka.PriceEvent{
			Timestamp: time.Now(),
			Symbol:    "ETH",
			Volume24h: 500000000.0,
		}

		detector.ProcessPriceEvent(event)

		if len(producer.signals) != 0 {
			t.Errorf("expected no signals with insufficient data, got %d", len(producer.signals))
		}
	})

	t.Run("volume spike detection", func(t *testing.T) {
		producer.signals = nil

		baseVolume := 1000000000.0
		for i := 0; i < 5; i++ {
			event := &kafka.PriceEvent{
				Timestamp: time.Now().Add(time.Duration(i) * time.Hour),
				Symbol:    "SPIKE",
				Volume24h: baseVolume,
			}
			detector.ProcessPriceEvent(event)
		}

		spikeEvent := &kafka.PriceEvent{
			Timestamp: time.Now().Add(6 * time.Hour),
			Symbol:    "SPIKE",
			Volume24h: baseVolume * 2.0,
		}
		detector.ProcessPriceEvent(spikeEvent)

		if len(producer.signals) == 0 {
			t.Error("expected volume spike signal to be generated")
		}

		signal := producer.signals[0]
		if signal.SignalType != "volume_spike" {
			t.Errorf("expected signal type 'volume_spike', got %s", signal.SignalType)
		}

		if signal.Symbol != "SPIKE" {
			t.Errorf("expected symbol 'SPIKE', got %s", signal.Symbol)
		}

		details := signal.Details
		if details["current_volume"] != baseVolume*2.0 {
			t.Errorf("expected current_volume %f, got %v", baseVolume*2.0, details["current_volume"])
		}

		if details["threshold_exceeded"] != threshold {
			t.Errorf("expected threshold_exceeded %f, got %v", threshold, details["threshold_exceeded"])
		}
	})
}

func TestCalculateAverage(t *testing.T) {
	tests := []struct {
		name     string
		volumes  []float64
		expected float64
	}{
		{
			name:     "empty volumes",
			volumes:  []float64{},
			expected: 0,
		},
		{
			name:     "single volume",
			volumes:  []float64{100},
			expected: 100,
		},
		{
			name:     "multiple volumes",
			volumes:  []float64{100, 200, 300},
			expected: 200,
		},
		{
			name:     "decimal volumes",
			volumes:  []float64{100.5, 200.5, 300.0},
			expected: 200.33333333333334,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateAverage(tt.volumes)
			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestVolumeDetector_HistoryManagement(t *testing.T) {
	producer := &mockProducer{}
	detector := NewVolumeDetector(producer, 1.3)

	now := time.Now()
	cutoff := now.Add(-time.Duration(VolumeDays+1) * 24 * time.Hour)

	oldEvent := &kafka.PriceEvent{
		Timestamp: cutoff,
		Symbol:    "OLD",
		Volume24h: 1000.0,
	}
	detector.ProcessPriceEvent(oldEvent)

	newEvent := &kafka.PriceEvent{
		Timestamp: now,
		Symbol:    "OLD",
		Volume24h: 2000.0,
	}
	detector.ProcessPriceEvent(newEvent)

	history := detector.volumeHistory["OLD"]
	if len(history.Volumes) != 1 {
		t.Errorf("expected old volume to be removed, got %d volumes", len(history.Volumes))
	}

	if history.Volumes[0] != 2000.0 {
		t.Errorf("expected remaining volume 2000.0, got %f", history.Volumes[0])
	}
}
