package services

import (
	"errors"
	"go-gin-api/dto"
	"go-gin-api/models"
	"go-gin-api/utils"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(input dto.RegisterInput) (*models.RegisterResponse, error)
	Login(input dto.LoginInput) (*string, *string, error)
	Logout(email string) error
	RefreshToken(refreshToken string) (*string, error)
}

type authService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{
		db: db,
	}
}

func (as authService) Register(input dto.RegisterInput) (*models.RegisterResponse, error) {
	user := models.User{
		Name: input.Name,
		Email: input.Email,
		Password: input.Password,
	}

	if err := user.HashPassword(); err != nil {
		return nil, errors.New("failed to hash password")
	}

	if result := as.db.Create(&user); result.Error != nil {
		return nil, errors.New("failed to create user")
	}

	return &models.RegisterResponse{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		RoleID: user.RoleID,
		IsActive: user.IsActive,
	}, nil
}

func (as authService) Login(input dto.LoginInput) (*string, *string, error) {
	var user models.User
	result := as.db.Where(&models.User{ Email: input.Email, IsActive: true }).First(&user)
	if result.Error != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if err := user.CheckPassword(input.Password); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, nil, err
	}

	user.RefreshToken = *refreshToken
	result = as.db.Save(&user)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return accessToken, refreshToken, nil
}

func (as authService) Logout(email string) error {
	var user models.User
	result := as.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return result.Error
	}
	user.RefreshToken = ""
	if err := as.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (as authService) RefreshToken(refreshToken string) (*string, error) {
	token, err := utils.ParseToken(refreshToken, true)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	
	var user models.User
	if err := as.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if user.RefreshToken != refreshToken {
		return nil, errors.New("token mismatch")
	}

	newAccessToken, _, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return newAccessToken, nil
}