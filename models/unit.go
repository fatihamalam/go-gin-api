package models

import (
	"time"

	"gorm.io/gorm"
)

type Unit struct {
	gorm.Model
	Name     string `json:"name" gorm:"uniqueIndex;not null"`
	Quantity uint   `json:"quantity" gorm:"not null;default:0"`
	Items    []Item `json:"items" gorm:"foreignKey:UnitID"`
}

func (Unit) TableName() string {
	return "master_data.units"
}

type UnitResponse struct {
	ID          uint      		`json:"id"`
	Name    	string 			`json:"name"`
	Quantity	uint   			`json:"quantity"`
	Items		[]ItemResponse	`json:"items,omitempty"`
	CreatedAt   time.Time 		`json:"created_at"`
	UpdatedAt   time.Time 		`json:"updated_at"`
}

func (u *Unit) ToResponse() UnitResponse {
	return UnitResponse{
		ID: u.ID,
		Name: u.Name,
		Quantity: u.Quantity,
		Items: MapItemsToResponse(u.Items),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func MapUnitsToResponse(units []Unit) []UnitResponse {
	responses := make([]UnitResponse, len(units))
	for i, unit := range units {
		responses[i] = unit.ToResponse()
	}
	return responses
}