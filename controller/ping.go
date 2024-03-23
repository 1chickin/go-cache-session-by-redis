package controller

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	RateLimitCount  = "2"
	RateLimitExpiry = 60 * time.Second
	PingProcessTime = 5 * time.Second
)

var (
	pingMutex sync.Mutex
)

func Ping(c *gin.Context, redisClient *redis.Client) {
	pingMutex.Lock()
	defer pingMutex.Unlock()

	sessionToken, err := c.Cookie("session_token")
	if sessionToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing session token"})
		return
	}

	userID, err := redisClient.Get(context.Background(), sessionToken).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session token"})
		return
	}

	// check rate limit
	rateKey := "rate:" + userID
	val, err := redisClient.Get(context.Background(), rateKey).Result()
	if err == nil && val >= RateLimitCount {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded in 1 minute period 2 time calling ping API!"})
		return
	}

	// set rate limit
	if err == redis.Nil {
		redisClient.Set(context.Background(), rateKey, "1", RateLimitExpiry)
	} else {
		redisClient.Incr(context.Background(), rateKey)
	}

	// count the number of calling ping API
	pingCountKey := "CallingPingAPI userID:" + userID
	redisClient.Incr(context.Background(), pingCountKey)

	// add userID to HyperLogLog for counting the number of users calling ping API
	redisClient.PFAdd(context.Background(), "pingUsersHyperLogLog", userID)

	// sleep the specified time
	time.Sleep(PingProcessTime)

	c.JSON(http.StatusOK, gin.H{"message": "ping pong"})
}
