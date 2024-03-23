package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

/*
function Count
- Method: GET
- Description: Provides an approximate count of users who have called the /ping API, leveraging HyperLogLog.
- Response 200 OK:
```json

	{
	  "estimatedCount": 150
	}
*/
func Count(c *gin.Context, redisClient *redis.Client) {
	// get count from HyperLogLog
	count, err := redisClient.PFCount(context.Background(), "pingUsersHyperLogLog").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get estimated count!",
		})
		return
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		"estimatedCount": count,
	})
}
