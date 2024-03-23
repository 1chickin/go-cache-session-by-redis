package main

import (
	"github.com/1chickin/go-cache-session-by-redis/config"
	"github.com/1chickin/go-cache-session-by-redis/controller"
	"github.com/1chickin/go-cache-session-by-redis/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
}

func main() {
	db := config.ConnectDatabase()
	config.MigrateDB()
	redisClient := config.ConnectRedis()

	router := gin.Default()
	router.POST("/signup", controller.Signup)
	router.POST("/login", func(c *gin.Context) {
		controller.Login(c, redisClient)
	})
	router.GET("/ping", middleware.Auth(redisClient, db), middleware.RateLimit(), func(c *gin.Context) {
		controller.Ping(c, redisClient)
	})

	// router.GET("/top", controller.Top)
	// router.GET("/count", controller.Count)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
