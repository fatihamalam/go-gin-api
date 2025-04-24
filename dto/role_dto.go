package dto

type CreateRoleInput struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"omitempty"`
}

type UpdateRoleInput struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"omitempty"`
}

type UpdatePermissionsInput struct {
	PermissionIds []string `json:"permissionIds" binding:"required"`
}