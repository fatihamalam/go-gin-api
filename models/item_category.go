package models

import (
	"time"

	"gorm.io/gorm"
)

type ItemCategory struct {
	gorm.Model
	Name       			string			`json:"name" gorm:"uniqueIndex; not null"`
	Description			string			`json:"description"`
	IsActive			bool			`json:"is_active" gorm:"default:true"`
	Items				[]Item			`json:"items" gorm:"foreignKey:ItemCategoryID"`
}

func (ItemCategory) TableName() string {
	return "master_data.item_categories"
}

type ItemCategoryResponse struct {
	ID			uint			`json:"id"`
	Name		string			`json:"name"`
	Description	string			`json:"description"`
	IsActive	bool			`json:"is_active"`
	Items		[]ItemResponse	`json:"items,omitempty"`
	CreatedAt 	time.Time		`json:"created_at"`
	UpdatedAt 	time.Time		`json:"updated_at"`
}

func (ic *ItemCategory) ToResponse() ItemCategoryResponse {
	return ItemCategoryResponse{
		ID: ic.ID,
		Name: ic.Name,
		Description: ic.Description,
		IsActive: ic.IsActive,
		Items: MapItemsToResponse(ic.Items),
		CreatedAt: ic.CreatedAt,
		UpdatedAt: ic.UpdatedAt,
	}
}

func MapItemCategoryToResponse(itemCategories []ItemCategory) []ItemCategoryResponse {
	responses := make([]ItemCategoryResponse, len(itemCategories))
	for i, itemCategory := range itemCategories {
		responses[i] = itemCategory.ToResponse()
	}
	return responses
}