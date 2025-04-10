package services

import (
	"errors"
	"go-gin-api/models"

	"gorm.io/gorm"
)

type UserService interface {
	FindAll() ([]models.User, error)
	FindOneByEmail(email string) (models.User, error)
	GetProfile(userID float64) (*models.User, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

func (us userService) FindAll() ([]models.User, error) {
	var users []models.User
	result := us.db.Find(&users)
	return users, result.Error
}

func (us userService) FindOneByEmail(email string) (models.User, error) {
	var user models.User
	result := us.db.Where("email = ?", email).First(&user)
	return user, result.Error
}

func (us userService) GetProfile(userID float64) (*models.User, error) {
	var user models.User
	if err := us.db.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}