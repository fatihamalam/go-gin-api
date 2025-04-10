package routes

import (
	"go-gin-api/controllers"
	"go-gin-api/middleware"
	"go-gin-api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(api *gin.RouterGroup, db *gorm.DB) {
	authService := services.NewAuthService(db)
	authController := controllers.NewAuthController(authService)

	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.RefreshToken)
		auth.POST("/logout", middleware.JWTAuthMiddleware(), authController.Logout)
	}
}