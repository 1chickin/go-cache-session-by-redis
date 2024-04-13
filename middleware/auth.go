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
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized session token!"})
				return
			}
			// check if the session has expired
			if time.Now().After(session.ExpiresAt) {
				// delete the session from the database
				db.Delete(&session)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session expired!"})
				return
			}
			// if database has valid session, set to Redis with the remaining time until expiration
			expireTime := time.Until(session.ExpiresAt)
			redisClient.Set(context.Background(), sessionToken, session.UserID, expireTime)
			userID = strconv.Itoa(int(session.UserID))
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized session token!"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
