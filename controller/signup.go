package controller

import (
	"net/http"

	"github.com/1chickin/go-cache-session-by-redis/config"
	"github.com/1chickin/go-cache-session-by-redis/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// get username & password from request
	var requestBody struct {
		Username string
		Password string
	}

	if c.Bind(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to load request body!",
		})
		return
	}

	// check username exist
	var existingUser model.User
	resultCheckExist := config.DB.Where("username = ?", requestBody.Username).First(&existingUser)
	if resultCheckExist.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username already exists!",
		})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password!",
		})
		return
	}

	// create user
	user := &model.User{Username: requestBody.Username, Password: string(hashedPassword)}
	result := config.DB.Create(&user) // pass pointer of data to Create
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user!",
		})
		return
	}

	// response
	c.JSON(http.StatusOK, gin.H{})
}
