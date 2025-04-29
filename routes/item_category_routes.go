package routes

import (
	"go-gin-api/controllers"
	"go-gin-api/middleware"
	"go-gin-api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupItemCategoryRoutes(api *gin.RouterGroup, db *gorm.DB) {
	itemCategoryService := services.NewItemCategoryService(db)
	itemCategoryController := controllers.NewItemCategoryController(itemCategoryService)

	api.Use(middleware.JWTAuthMiddleware())
	api.GET("/item-categories", middleware.PermissionMiddleware(db, "item_category:read"), itemCategoryController.GetAllItemCategories)
	api.POST("/item-categories", middleware.PermissionMiddleware(db, "item_category:write"), itemCategoryController.CreateItemCategory)
	api.DELETE("/item-categories/:id", middleware.PermissionMiddleware(db, "item_category:delete"), itemCategoryController.DeleteItemCategory)
}