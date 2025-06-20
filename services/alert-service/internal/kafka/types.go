package kafka

import "time"

type TradingSignal struct {
	Timestamp      time.Time              `json:"timestamp"`
	Symbol         string                 `json:"symbol"`
	SignalType     string                 `json:"signal_type"`
	SignalStrength string                 `json:"signal_strength"`
	Direction      string                 `json:"direction"`
	Details        map[string]interface{} `json:"details"`
	ServiceID      string                 `json:"service_id"`
}
