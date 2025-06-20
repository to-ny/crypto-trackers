package alerts

import (
	"sync"
	"time"
)

type RateLimiter struct {
	lastAlertTime map[string]time.Time
	cooldown      time.Duration
	mutex         sync.RWMutex
}

func NewRateLimiter(cooldownMinutes int) *RateLimiter {
	return &RateLimiter{
		lastAlertTime: make(map[string]time.Time),
		cooldown:      time.Duration(cooldownMinutes) * time.Minute,
		mutex:         sync.RWMutex{},
	}
}

func (r *RateLimiter) CanSendAlert(symbol string) bool {
	r.mutex.RLock()
	lastTime, exists := r.lastAlertTime[symbol]
	r.mutex.RUnlock()

	if !exists {
		return true
	}

	return time.Since(lastTime) >= r.cooldown
}

func (r *RateLimiter) RecordAlert(symbol string) {
	r.mutex.Lock()
	r.lastAlertTime[symbol] = time.Now()
	r.mutex.Unlock()
}
