package middlewares

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"soporte/internal/config"
)

type clientLimiter struct {
	tokens   float64
	lastSeen time.Time
}

type rateLimiterStore struct {
	mu      sync.Mutex
	clients map[string]*clientLimiter
	rate    float64
	burst   float64
	ttl     time.Duration
}

type allowResult struct {
	allowed   bool
	remaining int
	resetIn   int // segundos hasta que el balde esté lleno
}

func newRateLimiterStore(requests int, window time.Duration) *rateLimiterStore {
	return &rateLimiterStore{
		clients: make(map[string]*clientLimiter),
		rate:    float64(requests) / window.Seconds(),
		burst:   float64(requests),
		ttl:     window * 2,
	}
}

func (s *rateLimiterStore) allow(key string, now time.Time) allowResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ip, client := range s.clients {
		if now.Sub(client.lastSeen) > s.ttl {
			delete(s.clients, ip)
		}
	}

	client, ok := s.clients[key]
	if !ok {
		s.clients[key] = &clientLimiter{
			tokens:   s.burst - 1,
			lastSeen: now,
		}
		remaining := int(s.burst - 1)
		resetIn := int(math.Ceil(1.0 / s.rate))
		return allowResult{allowed: true, remaining: remaining, resetIn: resetIn}
	}

	elapsed := now.Sub(client.lastSeen).Seconds()
	client.tokens += elapsed * s.rate
	if client.tokens > s.burst {
		client.tokens = s.burst
	}
	client.lastSeen = now

	if client.tokens < 1 {
		// Segundos hasta tener 1 ficha disponible
		tokensNeeded := 1.0 - client.tokens
		resetIn := int(math.Ceil(tokensNeeded / s.rate))
		return allowResult{allowed: false, remaining: 0, resetIn: resetIn}
	}

	client.tokens--
	remaining := int(client.tokens)
	resetIn := int(math.Ceil((s.burst - client.tokens) / s.rate))
	return allowResult{allowed: true, remaining: remaining, resetIn: resetIn}
}

func RateLimit(cfg config.SecurityConfig) gin.HandlerFunc {
	store := newRateLimiterStore(cfg.RateLimitRequests, cfg.RateLimitWindow)
	limit := strconv.Itoa(cfg.RateLimitRequests)

	return func(c *gin.Context) {
		result := store.allow(c.ClientIP(), time.Now())

		c.Header("RateLimit-Limit", limit)
		c.Header("RateLimit-Remaining", strconv.Itoa(result.remaining))
		c.Header("RateLimit-Reset", strconv.Itoa(result.resetIn))

		if !result.allowed {
			c.Header("Retry-After", strconv.Itoa(result.resetIn))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    "rate_limited",
					"message": "too many requests",
				},
				"request_id": GetRequestID(c),
			})
			return
		}

		c.Next()
	}
}
