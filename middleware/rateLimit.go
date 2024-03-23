package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var PingProcessTime = 5 * time.Second

func RateLimit() gin.HandlerFunc {
	// save the last time a user called an endpoint
	var lastCalled = make(map[string]time.Time)

	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		if lastTimeCalled, exists := lastCalled[userID.(string)]; exists {
			// API /ping only allows 1 caller at a time
			if time.Since(lastTimeCalled) < PingProcessTime {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded in 5s period 1 time calling ping API!"})
				return
			}
		}

		lastCalled[userID.(string)] = time.Now()
		c.Next()
	}
}
