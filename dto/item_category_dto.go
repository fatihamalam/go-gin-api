package dto

type CreateItemCategoryInput struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"omitempty"`
}