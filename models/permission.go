package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name 		string	`json:"name" gorm:"uniqueIndex;not null"`
	Description	string	`json:"description"`
	Roles		[]Role	`json:"-" gorm:"many2many:role_permissions;"`
}

type PermissionResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}