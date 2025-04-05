package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace request context with the new timeout context
		c.Request = c.Request.WithContext(ctx)

		// Create a done channel to wait for handler completion
		done := make(chan struct{})

		go func() {
			c.Next() // Call the next handler(s)
			close(done)
		}()

		select {
		case <-done:
			// All good
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "request timed out",
			})
		}
	}
}
