package controllers

import (
	"go-gin-api/dto"
	"go-gin-api/services"
	"go-gin-api/utils"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ItemCategoryController interface {
	GetAllItemCategories(c *gin.Context)
	CreateItemCategory(c *gin.Context)
	DeleteItemCategory(c *gin.Context)
}

type itemCategoryController struct {
	itemCategoryService services.ItemCategoryService
}

func NewItemCategoryController(itemCategoryService services.ItemCategoryService) ItemCategoryController {
	return &itemCategoryController{
		itemCategoryService: itemCategoryService,
	}
}

func (icc itemCategoryController) GetAllItemCategories(c *gin.Context) {
	allowedFields := map[string]bool{
		"name": true,
	}
	search := c.Query("query")
	includes := utils.BuildRelation(c, []string{"Items"})
	order := utils.Sorting(c, "name", "asc", allowedFields)
	limit, page, offset := utils.Paginate(c)
	categories, totalData, err := icc.itemCategoryService.FindAll(search, order, limit, offset, includes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"categories": categories,
			"currentPage": page,
			"totalPage": math.Ceil(float64(*totalData) / float64(limit)),
			"totalData": totalData,
			"limit": limit,
		},
		"message": "Item categories fetched successfully",
	})
}

func (icc itemCategoryController) CreateItemCategory(c *gin.Context) {
	var input dto.CreateItemCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.ValidationErrorToText(err)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
			return
		}
		
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := icc.itemCategoryService.CreateItemCategory(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"category": category,
		},
		"message": "Item category created successfully",
	})
}

func (icc itemCategoryController) DeleteItemCategory(c *gin.Context) {
	id := c.Param("id")

	if err := icc.itemCategoryService.DeleteItemCategory(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Item category deleted successfully",
	})
}