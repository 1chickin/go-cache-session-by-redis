package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit() gin.HandlerFunc {
	// save the last time a user called an endpoint
	var lastCalled = make(map[string]time.Time)

	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		if lastTimeCalled, exists := lastCalled[userID.(string)]; exists {
			if time.Since(lastTimeCalled) < 5*time.Second {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
				return
			}
		}

		lastCalled[userID.(string)] = time.Now()
		c.Next()
	}
}
