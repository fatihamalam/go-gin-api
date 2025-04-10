package controllers

import (
	"go-gin-api/dto"
	"go-gin-api/services"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.authService.Register(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data": user,
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
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (ac authController) Logout(c *gin.Context) {
	email, _ := c.Get("email")

	if err := ac.authService.Logout(email.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}

func (ac authController) RefreshToken(c *gin.Context) {
	oldRefreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
		return
	}

	accessToken, refreshToken, err := ac.authService.RefreshToken(oldRefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refresh_token", *refreshToken, 7*24*3600, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}