package web

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type client struct {
	requests  int
	windowEnd time.Time
}

type RateLimiter struct {
	mu        sync.RWMutex
	clients   map[string]*client
	maxReqs   int
	window    time.Duration
	cleanupInterval time.Duration
}

func NewRateLimiter(maxReqs int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:         make(map[string]*client),
		maxReqs:         maxReqs,
		window:          window,
		cleanupInterval: window * 2,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := rl.clientIP(r)

		rl.mu.Lock()
		c, ok := rl.clients[ip]
		now := time.Now()
		if !ok || now.After(c.windowEnd) {
			rl.clients[ip] = &client{requests: 1, windowEnd: now.Add(rl.window)}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		c.requests++
		if c.requests > rl.maxReqs {
			rl.mu.Unlock()
			w.Header().Set("Retry-After", time.Until(c.windowEnd).String())
			RespondWithError(w, http.StatusTooManyRequests, "rate_limited", "Too many requests, please try again later")
			return
		}
		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, c := range rl.clients {
			if now.After(c.windowEnd) {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}
