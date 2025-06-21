package alerts

import (
	"testing"
	"time"
)

func TestRateLimiter_CanSendAlert(t *testing.T) {
	limiter := NewRateLimiter(5)

	t.Run("first alert allowed", func(t *testing.T) {
		if !limiter.CanSendAlert("BTC") {
			t.Error("expected first alert to be allowed")
		}
	})

	t.Run("immediate second alert blocked", func(t *testing.T) {
		limiter.RecordAlert("BTC")
		if limiter.CanSendAlert("BTC") {
			t.Error("expected immediate second alert to be blocked")
		}
	})

	t.Run("different symbol allowed", func(t *testing.T) {
		limiter.RecordAlert("BTC")
		if !limiter.CanSendAlert("ETH") {
			t.Error("expected different symbol to be allowed")
		}
	})

	t.Run("alert allowed after cooldown", func(t *testing.T) {
		shortLimiter := NewRateLimiter(0)
		shortLimiter.RecordAlert("TEST")
		time.Sleep(1 * time.Millisecond)
		if !shortLimiter.CanSendAlert("TEST") {
			t.Error("expected alert to be allowed after cooldown")
		}
	})
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	limiter := NewRateLimiter(5)
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(i int) {
			symbol := "CONCURRENT"
			limiter.CanSendAlert(symbol)
			limiter.RecordAlert(symbol)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if !limiter.CanSendAlert("DIFFERENT") {
		t.Error("concurrent access should not affect different symbols")
	}
}
