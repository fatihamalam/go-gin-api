package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name			string 	`json:"name" gorm:"not null"`
	Email			string 	`json:"email" gorm:"uniqueIndex;not null"`
	Password		string 	`json:"-" gorm:"not null"`
	RefreshToken	string 	`json:"refresh_token,omitempty"`
	RoleID			uint   	`json:"role_id" gorm:"default:2"`
	Role			Role	`json:"role" gorm:"foreignKey:RoleID"`
	IsActive		bool	`json:"is_active" gorm:"default:true"`
}

func (User) TableName() string {
	return "master_data.users"
}

type UserResponse struct {
	ID        uint          `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	RoleID    uint          `json:"role_id"`
	Role      *RoleResponse `json:"role,omitempty"`
	IsActive  bool          `json:"is_active"`
	CreatedAt time.Time		`json:"created_at"`
	UpdatedAt time.Time		`json:"updated_at"`
}

type RegisterResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	RoleID   uint   `json:"role_id"`
	IsActive bool   `json:"is_active"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) ToResponse() UserResponse {
	var roleResponse RoleResponse
	if u.Role.ID != 0 {
		roleResponse = u.Role.ToResponse()
	}
	
	return UserResponse{
		ID: u.ID,
		Name: u.Name,
		Email: u.Email,
		RoleID: u.RoleID,
		Role: &roleResponse,
		IsActive: u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func MapUsersToResponse(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}
	return responses
}