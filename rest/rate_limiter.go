package rest

import (
	"context"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	lastCall map[string]time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		lastCall: make(map[string]time.Time),
	}
}

func (r *RateLimiter) Wait(ctx context.Context, endpoint string, method string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := method + " " + endpoint
	last := r.lastCall[key]

	cooldown := 250 * time.Millisecond
	elapsed := time.Since(last)
	if elapsed < cooldown {
		wait := cooldown - elapsed
		timer := time.NewTimer(wait)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}

	r.lastCall[key] = time.Now()
	return nil
}
