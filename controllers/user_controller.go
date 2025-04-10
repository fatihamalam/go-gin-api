package controllers

import (
	"go-gin-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Profile(*gin.Context)
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (uc userController) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	profile, err := uc.userService.GetProfile(userID.(float64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": profile})
}