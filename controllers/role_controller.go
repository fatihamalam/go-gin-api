package controllers

import (
	"go-gin-api/dto"
	"go-gin-api/services"
	"go-gin-api/utils"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleController interface {
	GetAllRoles(c *gin.Context)
	GetRole(c *gin.Context)
	CreateRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	DeleteRole(c *gin.Context)
	GetPermissionsByRoleID(c *gin.Context)
}

type roleController struct{
	roleService services.RoleService
}

func NewRoleController(roleService services.RoleService) RoleController {
	return &roleController{
		roleService: roleService,
	}
}

func (rc roleController) GetAllRoles(c *gin.Context) {
	allowedFields := map[string]bool{
		"name": true, "created_at": true,
	}
	search := c.Query("query")
	includes := utils.BuildRelation(c, []string{"Permissions"})
	order := utils.Sorting(c, "name", "asc", allowedFields)
	limit, page, offset := utils.Paginate(c)
	roles, totalData, err := rc.roleService.FindAll(search, order, limit, offset, includes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"roles": roles,
			"currentPage": page,
			"totalPage": math.Ceil(float64(*totalData) / float64(limit)),
			"totalData": totalData,
			"limit": limit,
		},
		"message": "Roles fetched successfully",
	})
}

func (rc roleController) GetRole(c *gin.Context) {
	id := c.Param("id")
	role, err := rc.roleService.FindOneByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"role": role,
		},
		"message": "Role fetched successfully",
	})
}

func (rc roleController) CreateRole(c *gin.Context) {
	var input dto.CreateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.ValidationErrorToText(err)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
			return
		}
		
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := rc.roleService.CreateRole(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"role": role,
		},
		"message": "Role created successfully",
	})
}

func (rc roleController) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	
	var input dto.UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.ValidationErrorToText(err)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
			return
		}
		
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := rc.roleService.UpdateRole(roleID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"role": role,
		},
		"message": "Role updated successfully",
	})
}

func (rc roleController) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")

	if err := rc.roleService.DeleteRole(roleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Role deleted successfully",
	})
}

func (rc roleController) GetPermissionsByRoleID(c *gin.Context) {
	roleID := c.Param("id")

	permissions, err := rc.roleService.GetPermissionsByRoleID(roleID);
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"permissions": permissions,
		},
		"message": "Permissions fetched successfully",
	})
}