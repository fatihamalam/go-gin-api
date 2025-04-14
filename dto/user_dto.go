package dto

type CreateUserInput struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email,uniqueEmail"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserInput struct {
	Name     string  `json:"name" binding:"required,min=3"`
	Email    string  `json:"email" binding:"required,email,uniqueEmail"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6"`
	RoleID   *uint   `josn:"role_id,omitempty"`
}