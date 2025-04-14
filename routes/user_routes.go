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

	api.Use(middleware.JWTAuthMiddleware())
	api.GET("/profile", userController.Profile)
	api.GET("/users", middleware.PermissionMiddleware(db, "user:read"), userController.GetAllUsers)
	api.GET("/users/:id", middleware.PermissionMiddleware(db, "user:read"), userController.GetUser)
	api.POST("/users", middleware.PermissionMiddleware(db, "user:write"), userController.CreateUser)
	api.PUT("/users/:id", middleware.PermissionMiddleware(db, "user:write"), userController.UpdateUser)
	api.DELETE("/users/:id", middleware.PermissionMiddleware(db, "user:delete"), userController.DeleteUser)
}