package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter manages rate limiters for different IP addresses.
type IPRateLimiter struct {
	ips      map[string]*limiterEntry
	mu       sync.RWMutex
	r        rate.Limit
	b        int
	cleanTTL time.Duration
}

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewIPRateLimiter creates a new IP-based rate limiter with automatic cleanup.
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	rl := &IPRateLimiter{
		ips:      make(map[string]*limiterEntry),
		r:        r,
		b:        b,
		cleanTTL: 10 * time.Minute,
	}

	// Background goroutine to clean up stale entries
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

// cleanup removes IP entries not seen within the TTL window.
func (i *IPRateLimiter) cleanup() {
	i.mu.Lock()
	defer i.mu.Unlock()
	now := time.Now()
	for ip, entry := range i.ips {
		if now.Sub(entry.lastSeen) > i.cleanTTL {
			delete(i.ips, ip)
		}
	}
}

// GetLimiter returns the rate limiter for the given IP, creating it if it doesn't exist.
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	entry, exists := i.ips[ip]
	i.mu.RUnlock()

	if exists {
		i.mu.Lock()
		entry.lastSeen = time.Now()
		i.mu.Unlock()
		return entry.limiter
	}

	i.mu.Lock()
	// Double check after acquiring write lock
	entry, exists = i.ips[ip]
	if !exists {
		entry = &limiterEntry{
			limiter:  rate.NewLimiter(i.r, i.b),
			lastSeen: time.Now(),
		}
		i.ips[ip] = entry
	}
	i.mu.Unlock()

	return entry.limiter
}

// RateLimit middleware limits requests based on client IP.
func RateLimit(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)
		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Çok fazla istek gönderildi. Lütfen bir süre sonra tekrar deneyin.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
