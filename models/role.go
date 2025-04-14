package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name		string			`json:"name"`
	Description	string			`json:"description"`
	Permissions	[]Permission 	`json:"permissions" gorm:"many2many:role_permissions;"`
}

type RoleResponse struct {
	ID          uint      				`json:"id"`
	Name        string    				`json:"name"`
	Description string    				`json:"description"`
	Permissions	[]PermissionResponse	`json:"permissions,omitempty"`
	CreatedAt   time.Time 				`json:"created_at"`
	UpdatedAt   time.Time 				`json:"updated_at"`
}