package services

import (
	"errors"
	"go-gin-api/dto"
	"go-gin-api/models"

	"gorm.io/gorm"
)

type UserService interface {
	FindAll(search string, order string, limit int, offset int) (*[]models.UserResponse, *int64, error)
	FindOneByID(userID string) (*models.UserResponse, error)
	FindOneByEmail(email string) (*models.UserResponse, error)
	GetProfile(userID float64) (*models.UserResponse, error)
	CreateUser(input dto.CreateUserInput) (*models.UserResponse, error)
	UpdateUser(userID string, input dto.UpdateUserInput) (*models.UserResponse, error)
	DeleteUser(userID string) error
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

func (us userService) FindAll(search string, order string, limit int, offset int) (*[]models.UserResponse, *int64, error) {
	var users []models.User
	var totalData int64

	query := us.db.Model(&models.User{})
	query.Order(order)
	
	if search != "" {
		likeSearch := "%" + search + "%"
		query.Where("name ILIKE ? OR email ILIKE ?", likeSearch, likeSearch)
	}

	result := query.Find(&users).Limit(limit).Offset(offset)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	result = query.Count(&totalData)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	userResponse := models.MapUsersToResponse(users)
	
	return &userResponse, &totalData, nil
}

func (us userService) FindOneByID(userID string) (*models.UserResponse, error) {
	var user models.User
	if err := us.db.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	userResponse := user.ToResponse()

	return &userResponse, nil
}

func (us userService) FindOneByEmail(email string) (*models.UserResponse, error) {
	var user models.User
	if err := us.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	userResponse := user.ToResponse()

	return &userResponse, nil
}

func (us userService) GetProfile(userID float64) (*models.UserResponse, error) {
	var user models.User
	if err := us.db.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	userResponse := user.ToResponse()

	return &userResponse, nil
}

func (us userService) CreateUser(input dto.CreateUserInput) (*models.UserResponse, error) {
	user := models.User{
		Name: input.Name,
		Email: input.Email,
		Password: input.Password,
	}

	if err := user.HashPassword(); err != nil {
		return nil, errors.New("failed to hash password")
	}

	if err := us.db.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	userResponse := user.ToResponse()
	
	return &userResponse, nil
}

func (us userService) UpdateUser(userID string, input dto.UpdateUserInput) (*models.UserResponse, error) {
	var user models.User
	if err := us.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = input.Name
	user.Email = input.Email
	if input.Password != nil {
		user.Password = *input.Password
		user.HashPassword()
	}
	if input.RoleID != nil {
		user.RoleID = *input.RoleID
	}

	if err := us.db.Save(&user).Error; err != nil {
		return nil, errors.New("failed to update user")
	}

	userResponse := user.ToResponse()

	return &userResponse, nil
}

func (us userService) DeleteUser(userID string) error {
	var user models.User
	if err := us.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	if err := us.db.Delete(&user).Error; err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}