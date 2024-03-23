package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/1chickin/go-cache-session-by-redis/model"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Auth(redisClient *redis.Client, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No session token found"})
			c.Abort()
			return
		}

		// check session in Redis
		userID, err := redisClient.Get(context.Background(), sessionToken).Result()
		if err == redis.Nil {
			// if not exist, check in database
			var session model.Session
			if result := db.Where("session_token = ?", sessionToken).First(&session); result.Error != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized session token in database " + sessionToken})
				return
			}
			// if database has session, set to Redis
			redisClient.Set(context.Background(), sessionToken, session.UserID, 30*time.Minute)
			userID = strconv.Itoa(int(session.UserID))
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized session token in Redis"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
