package controllers

import (
	"go-gin-api/dto"
	"go-gin-api/services"
	"go-gin-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface{
	Register(*gin.Context)
	Login(*gin.Context)
	Logout(*gin.Context)
	RefreshToken(*gin.Context)
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (ac authController) Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.ValidationErrorToText(err)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.authService.Register(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"user": user,
		},
		"message": "User registered successfully",
	})
}

func (ac authController) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := ac.authService.Login(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refresh_token", *refreshToken, 7*24*3600, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"access_token": accessToken,
		},
		"message": "User logged in successfully",
	})
}

func (ac authController) Logout(c *gin.Context) {
	email, _ := c.Get("email")

	if err := ac.authService.Logout(email.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "User logged out successfully",
	})
}

func (ac authController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
		return
	}

	accessToken, err := ac.authService.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"access_token": accessToken,
		},
		"message": "Access token refreshed successfully",
	})
}