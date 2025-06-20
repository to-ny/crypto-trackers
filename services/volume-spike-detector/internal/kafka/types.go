package kafka

import "time"

type PriceEvent struct {
	Timestamp      time.Time `json:"timestamp"`
	Symbol         string    `json:"symbol"`
	PriceUSD       float64   `json:"price_usd"`
	Volume24h      float64   `json:"volume_24h"`
	MarketCap      float64   `json:"market_cap"`
	PriceChange24h float64   `json:"price_change_24h"`
	Source         string    `json:"source"`
}

type TradingSignal struct {
	Timestamp      time.Time              `json:"timestamp"`
	Symbol         string                 `json:"symbol"`
	SignalType     string                 `json:"signal_type"`
	SignalStrength string                 `json:"signal_strength"`
	Direction      string                 `json:"direction"`
	Details        map[string]interface{} `json:"details"`
	ServiceID      string                 `json:"service_id"`
}
