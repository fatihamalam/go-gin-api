package dto

type CreateUnitInput struct {
	Name     string `json:"name" binding:"required,min=3"`
	Quantity uint   `json:"quantity" binding:"required,gt=0"`
}