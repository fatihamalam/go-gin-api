package routes

import (
	"go-gin-api/controllers"
	"go-gin-api/middleware"
	"go-gin-api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(api *gin.RouterGroup, db *gorm.DB) {
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	api.GET("/profile", middleware.JWTAuthMiddleware(), userController.Profile)
}