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

func (Permission) TableName() string {
	return "master_data.permissions"
}

type PermissionResponse struct {
	ID          uint      		`json:"id"`
	Name        string    		`json:"name"`
	Description string    		`json:"description"`
	Roles		[]RoleResponse	`json:"roles,omitempty"`
	CreatedAt   time.Time 		`json:"created_at"`
	UpdatedAt   time.Time 		`json:"updated_at"`
}

func (p *Permission) ToResponse() PermissionResponse {
	return PermissionResponse{
		ID: p.ID,
		Name: p.Name,
		Description: p.Description,
		Roles: MapRolesToResponse(p.Roles),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func MapPermissionsToResponse(permissions []Permission) []PermissionResponse {
	responses := make([]PermissionResponse, len(permissions))
	for i, permission := range permissions {
		responses[i] = permission.ToResponse()
	}
	return responses
}