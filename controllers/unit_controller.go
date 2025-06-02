package controllers

import (
	"go-gin-api/services"
	"go-gin-api/utils"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UnitController interface {
	GetAllUnits(c *gin.Context)
}

type unitController struct {
	unitService services.UnitService
}

func NewUnitController(unitService services.UnitService) UnitController {
	return &unitController{
		unitService: unitService,
	}
}

func (uc unitController) GetAllUnits(c *gin.Context) {
	allowedFields := map[string]bool{
		"name": true,
	}
	search := c.Query("query")
	order := utils.Sorting(c, "name", "asc", allowedFields)
	limit, page, offset := utils.Paginate(c)
	units, totalData, err := uc.unitService.FindAll(search, order, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"units": units,
			"currentPage": page,
			"totalPage": math.Ceil(float64(*totalData) / float64(limit)),
			"totalData": totalData,
			"limit": limit,
		},
		"message": "Units fetched successfully",
	})
}