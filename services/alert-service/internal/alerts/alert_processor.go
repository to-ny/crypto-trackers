package alerts

import (
	"alert-service/internal/kafka"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type AlertProcessor struct {
	rateLimiter       *RateLimiter
	alertsReceived    prometheus.CounterVec
	alertsSent        prometheus.CounterVec
	alertsRateLimited prometheus.CounterVec
}

func NewAlertProcessor(cooldownMinutes int, alertsReceived, alertsSent, alertsRateLimited prometheus.CounterVec) *AlertProcessor {
	return &AlertProcessor{
		rateLimiter:       NewRateLimiter(cooldownMinutes),
		alertsReceived:    alertsReceived,
		alertsSent:        alertsSent,
		alertsRateLimited: alertsRateLimited,
	}
}

func (a *AlertProcessor) ProcessSignal(signal *kafka.TradingSignal) error {
	a.alertsReceived.WithLabelValues(signal.Symbol, signal.SignalType).Inc()

	if !a.rateLimiter.CanSendAlert(signal.Symbol) {
		a.alertsRateLimited.WithLabelValues(signal.Symbol).Inc()
		log.Printf("SKIPPED: Alert for %s within cooldown period (last sent < 5 min ago)", signal.Symbol)
		return nil
	}

	a.sendAlert(signal)
	a.rateLimiter.RecordAlert(signal.Symbol)
	a.alertsSent.WithLabelValues(signal.Symbol).Inc()
	log.Printf("SENT: Alert recorded for %s", signal.Symbol)

	return nil
}

func (a *AlertProcessor) sendAlert(signal *kafka.TradingSignal) {
	alert := a.formatAlert(signal)
	fmt.Println("🚨 TRADING SIGNAL ALERT 🚨")
	fmt.Println(alert)
	fmt.Println("=" + strings.Repeat("=", 50))

	log.Printf("ALERT: %s %s signal for %s (strength: %s)",
		signal.SignalType, signal.Direction, signal.Symbol, signal.SignalStrength)
}

func (a *AlertProcessor) formatAlert(signal *kafka.TradingSignal) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Symbol: %s\n", signal.Symbol))
	builder.WriteString(fmt.Sprintf("Signal Type: %s\n", signal.SignalType))
	builder.WriteString(fmt.Sprintf("Direction: %s\n", signal.Direction))
	builder.WriteString(fmt.Sprintf("Strength: %s\n", signal.SignalStrength))
	builder.WriteString(fmt.Sprintf("Time: %s\n", signal.Timestamp.Format(time.RFC3339)))
	builder.WriteString(fmt.Sprintf("Service: %s\n", signal.ServiceID))

	if len(signal.Details) > 0 {
		builder.WriteString("Details:\n")
		for key, value := range signal.Details {
			builder.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
	}

	return builder.String()
}
