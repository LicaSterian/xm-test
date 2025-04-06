package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Limiter middleware config
type ClientLimiter struct {
	clients map[string]*rate.Limiter
	mu      sync.Mutex
	r       rate.Limit
	b       int
}

func NewClientLimiter(r rate.Limit, b int) *ClientLimiter {
	return &ClientLimiter{
		clients: make(map[string]*rate.Limiter),
		r:       r,
		b:       b,
	}
}

func (cl *ClientLimiter) getLimiter(ip string) *rate.Limiter {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	limiter, exists := cl.clients[ip]
	if !exists {
		limiter = rate.NewLimiter(cl.r, cl.b)
		cl.clients[ip] = limiter
	}
	return limiter
}

func RateLimitMiddleware(cl *ClientLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := cl.getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
