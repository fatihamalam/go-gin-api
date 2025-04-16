package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name       			string 			`json:"name" gorm:"uniqueIndex; not null"`
	Description			string 			`json:"description"`
	SKU					string			`json:"sku" gorm:"uniqueIndex"`
	Barcode				string			`json:"barcode" gorm:"uniqueIndex"`
	SellPrice			float64			`json:"sell_price" gorm:"not null;default:0"`
	BuyPrice			float64			`json:"buy_price" gorm:"not null;default:0"`
	BuyPriceAvg			float64			`json:"buy_price_avg" gorm:"not null;default:0"`
	Quantity			uint			`json:"quantity" gorm:"not null;default:0"`
	Cost				float64			`json:"cost" gorm:"not null;default:0"`
	MinStock			uint			`json:"min_stock" gorm:"not null; default:0"`
	IsActive			bool			`json:"is_active" gorm:"not null;default:true"`
	ItemCategoryID 		uint 			`json:"item_category_id"`
	ItemCategory		ItemCategory	`json:"item_category" gorm:"foreignKey:ItemCategoryID"`
	UnitID				uint			`json:"unit_id"`
	Unit				Unit			`json:"unit" gorm:"foreignKey:UnitID"`
}

func (Item) TableName() string {
	return "master_data.items"
}

type ItemResponse struct {
	ID					uint					`json:"id"`
	Name       			string 					`json:"name"`
	Description			string 					`json:"description"`
	SKU					string					`json:"sku"`
	Barcode				string					`json:"barcode"`
	SellPrice			float64					`json:"sell_price"`
	BuyPrice			float64					`json:"buy_price"`
	BuyPriceAvg			float64					`json:"buy_price_avg"`
	Quantity			uint					`json:"quantity"`
	Cost				float64					`json:"cost"`
	MinStock			uint					`json:"min_stock"`
	IsActive			bool					`json:"is_active"`
	ItemCategoryID 		uint 					`json:"item_category_id"`
	ItemCategory		*ItemCategoryResponse	`json:"item_category,omitempty"`
	UnitID				uint					`json:"unit_id"`
	Unit				*UnitResponse			`json:"unit,omitempty"`
	CreatedAt 			time.Time				`json:"created_at"`
	UpdatedAt 			time.Time				`json:"updated_at"`
}

func (i *Item) ToResponse() ItemResponse {
	var itemCategoryResponse ItemCategoryResponse
	if i.ItemCategory.ID != 0 {
		itemCategoryResponse = i.ItemCategory.ToResponse()
	}

	var unitResponse UnitResponse
	if i.Unit.ID != 0 {
		unitResponse = i.Unit.ToResponse()
	}
	
	return ItemResponse{
		ID: i.ID,
		Name: i.Name,
		Description: i.Description,
		SKU: i.SKU,
		Barcode: i.Barcode,
		SellPrice: i.SellPrice,
		BuyPrice: i.BuyPrice,
		BuyPriceAvg: i.BuyPriceAvg,
		Quantity: i.Quantity,
		Cost: i.Cost,
		MinStock: i.MinStock,
		IsActive: i.IsActive,
		ItemCategoryID: i.ItemCategoryID,
		ItemCategory: &itemCategoryResponse,
		UnitID: i.UnitID,
		Unit: &unitResponse,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

func MapItemsToResponse(items []Item) []ItemResponse {
	responses := make([]ItemResponse, len(items))
	for i, item := range items {
		responses[i] = item.ToResponse()
	}
	return responses
}