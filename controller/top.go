package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

/*
function Top.
- Method: GET
- Description: Returns the top 10 users based on the frequency of API calls from Redis.
- Responses 200 OK:
```json

	{
		"topUsersCallingAPIAllTime": [
			"CallingPingAPI userID:1 called 1 times",
			"CallingPingAPI userID:3 called 4 times"
		]
	}
*/
func TopUsersCallingAPI(c *gin.Context, redisClient *redis.Client) {
	// get top users with the number of calling api from Redis
	topUsers, err := getTopUsers(redisClient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get top users!",
		})
		return
	}

	// loop through top users to get the number of calling api
	for i, key := range topUsers {
		count, err := redisClient.Get(context.Background(), key).Result()
		if err != nil {
			fmt.Printf("Failed to get count of %s: %v\n", key, err)
			os.Exit(1)
		}
		// add to response the user id and the coressponding number of calling api
		topUsers[i] = key + " called " + count + " times"
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		"topUsersCallingAPIAllTime": topUsers,
	})
}

func getTopUsers(redisClient *redis.Client) ([]string, error) {
	keys, err := redisClient.Keys(context.Background(), "CallingPingAPI userID:*").Result()
	if err != nil {
		fmt.Printf("Failed to get top users: %v\n", err)
		os.Exit(1)
	}

	return keys, err
}
