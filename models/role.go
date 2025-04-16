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

func (Role) TableName() string {
	return "master_data.roles"
}

type RoleResponse struct {
	ID          uint      				`json:"id"`
	Name        string    				`json:"name"`
	Description string    				`json:"description"`
	Permissions	[]PermissionResponse	`json:"permissions,omitempty"`
	CreatedAt   time.Time 				`json:"created_at"`
	UpdatedAt   time.Time 				`json:"updated_at"`
}

func (r *Role) ToResponse() RoleResponse {
	return RoleResponse{
		ID: r.ID,
		Name: r.Name,
		Description: r.Description,
		Permissions: MapPermissionsToResponse(r.Permissions),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func MapRolesToResponse(roles []Role) []RoleResponse {
	responses := make([]RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = role.ToResponse()
	}
	return responses
}