package signals

import (
	"log"
	"sync"
	"time"
	"volume-spike-detector/internal/kafka"
)

const (
	VolumeDays     = 7
	HoursInDay     = 24
	MaxHistorySize = VolumeDays * HoursInDay
)

type VolumeHistory struct {
	Volumes   []float64
	Timestamps []time.Time
	mutex     sync.RWMutex
}

type VolumeDetector struct {
	volumeHistory map[string]*VolumeHistory
	threshold     float64
	mutex         sync.RWMutex
	producer      *kafka.Producer
}

func NewVolumeDetector(producer *kafka.Producer, threshold float64) *VolumeDetector {
	return &VolumeDetector{
		volumeHistory: make(map[string]*VolumeHistory),
		threshold:     threshold,
		producer:      producer,
	}
}

func (vd *VolumeDetector) ProcessPriceEvent(event *kafka.PriceEvent) error {
	vd.mutex.Lock()
	defer vd.mutex.Unlock()

	history, exists := vd.volumeHistory[event.Symbol]
	if !exists {
		history = &VolumeHistory{
			Volumes:    make([]float64, 0, MaxHistorySize),
			Timestamps: make([]time.Time, 0, MaxHistorySize),
		}
		vd.volumeHistory[event.Symbol] = history
		log.Printf("Started tracking volume history for %s", event.Symbol)
	}

	history.mutex.Lock()
	defer history.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-time.Duration(VolumeDays) * 24 * time.Hour)

	for i := 0; i < len(history.Timestamps); {
		if history.Timestamps[i].Before(cutoff) {
			history.Volumes = history.Volumes[1:]
			history.Timestamps = history.Timestamps[1:]
		} else {
			break
		}
	}

	history.Volumes = append(history.Volumes, event.Volume24h)
	history.Timestamps = append(history.Timestamps, event.Timestamp)

	if len(history.Volumes) > MaxHistorySize {
		history.Volumes = history.Volumes[1:]
		history.Timestamps = history.Timestamps[1:]
	}

	volumeCount := len(history.Volumes)
	log.Printf("Processed volume event for %s: %.0f (history: %d points)", event.Symbol, event.Volume24h, volumeCount)

	if volumeCount >= 2 {
		return vd.checkForVolumeSpike(event.Symbol, event.Timestamp, event.Volume24h, history.Volumes)
	}

	return nil
}

func (vd *VolumeDetector) checkForVolumeSpike(symbol string, timestamp time.Time, currentVolume float64, volumeHistory []float64) error {
	if len(volumeHistory) < 2 {
		return nil
	}

	avg7Day := calculateAverage(volumeHistory[:len(volumeHistory)-1])
	
	if avg7Day == 0 {
		return nil
	}

	spikeMultiplier := currentVolume / avg7Day

	if spikeMultiplier > vd.threshold {
		return vd.publishVolumeSpike(symbol, timestamp, currentVolume, avg7Day, spikeMultiplier)
	}

	return nil
}

func (vd *VolumeDetector) publishVolumeSpike(symbol string, timestamp time.Time, currentVolume, avg7Day, spikeMultiplier float64) error {
	signalStrength := "medium"
	if spikeMultiplier > 2.0 {
		signalStrength = "strong"
	} else if spikeMultiplier < 1.5 {
		signalStrength = "weak"
	}

	signal := &kafka.TradingSignal{
		Timestamp:      timestamp,
		Symbol:         symbol,
		SignalType:     "volume_spike",
		SignalStrength: signalStrength,
		Direction:      "bullish",
		Details: map[string]interface{}{
			"current_volume":     currentVolume,
			"avg_volume_7d":      avg7Day,
			"spike_multiplier":   spikeMultiplier,
			"threshold_exceeded": vd.threshold,
		},
		ServiceID: "volume-detector-v1",
	}

	if err := vd.producer.PublishSignal("trading-signals", signal); err != nil {
		log.Printf("Failed to publish volume spike signal for %s: %v", symbol, err)
		return err
	}

	log.Printf("Published volume spike signal for %s (%.1fx spike: current=%.0f, avg=%.0f)", 
		symbol, spikeMultiplier, currentVolume, avg7Day)
	return nil
}

func calculateAverage(volumes []float64) float64 {
	if len(volumes) == 0 {
		return 0
	}

	sum := 0.0
	for _, volume := range volumes {
		sum += volume
	}

	return sum / float64(len(volumes))
}