package controllers

import (
	"go-gin-api/dto"
	"go-gin-api/services"
	"go-gin-api/utils"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Profile(*gin.Context)
	GetAllUsers(c *gin.Context)
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (uc userController) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	profile, err := uc.userService.GetProfile(userID.(float64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"profile": profile,
		},
		"message": "Profile fetched successfully",
	})
}

func (uc userController) GetAllUsers(c *gin.Context) {
	allowedFields := map[string]bool{
		"name": true, "email": true, "created_at": true,
	}
	search := c.Query("query")
	order := utils.Sorting(c, "name", "asc", allowedFields)
	limit, page, offset := utils.Paginate(c)
	users, totalData, err := uc.userService.FindAll(search, order, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"users": users,
			"currentPage": page,
			"totalPage": math.Ceil(float64(*totalData) / float64(limit)),
			"totalData": totalData,
			"limit": limit,
		},
		"message": "Users fetched successfully",
	})
}

func (uc userController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.userService.FindOneByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": user,
		},
		"message": "User fetched successfully",
	})
}

func (uc userController) CreateUser(c *gin.Context) {
	var input dto.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.ValidationErrorToText(err)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
			return
		}
		
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"user": user,
		},
		"message": "User created successfully",
	})
}

func (uc userController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	utils.GinCtx = c
	defer func () { utils.GinCtx = nil }()
	
	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.ValidationErrorToText(err)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
			return
		}
		
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.UpdateUser(userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": user,
		},
		"message": "User updated successfully",
	})
}

func (uc userController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if err := uc.userService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "User deleted successfully",
	})
}