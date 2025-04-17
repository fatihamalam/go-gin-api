package routes

import (
	"go-gin-api/controllers"
	"go-gin-api/middleware"
	"go-gin-api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoleRoutes(api *gin.RouterGroup, db *gorm.DB) {
	roleService := services.NewRoleService(db)
	roleController := controllers.NewRoleController(roleService)

	api.Use(middleware.JWTAuthMiddleware())
	api.GET("/roles", middleware.PermissionMiddleware(db, "role:read"), roleController.GetAllRoles)
	api.GET("/roles/:id", middleware.PermissionMiddleware(db, "role:read"), roleController.GetRole)
	api.POST("/roles", middleware.PermissionMiddleware(db, "role:write"), roleController.CreateRole)
	api.PUT("/roles/:id", middleware.PermissionMiddleware(db, "role:write"), roleController.UpdateRole)
	api.DELETE("/roles/:id", middleware.PermissionMiddleware(db, "role:delete"), roleController.DeleteRole)
}