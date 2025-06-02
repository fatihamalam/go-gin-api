package services

import (
	"errors"
	"go-gin-api/dto"
	"go-gin-api/models"

	"gorm.io/gorm"
)

type UnitService interface {
	FindAll(search string, order string, limit int, offset int) (*[]models.UnitResponse, *int64, error)
	CreateUnit(input dto.CreateUnitInput) (*models.UnitResponse, error)
	DeleteUnit(id string) error
}

type unitService struct {
	db *gorm.DB
}

func NewUnitService(db *gorm.DB) UnitService {
	return &unitService{
		db: db,
	}
}

func (us unitService) FindAll(search string, order string, limit int, offset int) (*[]models.UnitResponse, *int64, error) {
	var units []models.Unit
	var totalData int64

	query := us.db.Model(&models.Unit{})
	query.Order(order)

	if search != "" {
		if search != "" {
			likeSearch := "%" + search + "%"
			query.Where("name ILIKE ?", likeSearch)
		}
	}

	if err := query.Find(&units).Limit(limit).Offset(offset).Error; err != nil {
		return nil, nil, err
	}

	if err := query.Count(&totalData).Error; err != nil {
		return nil, nil, err
	}

	unitResponse := models.MapUnitsToResponse(units)

	return &unitResponse, &totalData, nil
}

func (us unitService) CreateUnit(input dto.CreateUnitInput) (*models.UnitResponse, error) {
	unit := models.Unit{
		Name: input.Name,
		Quantity: input.Quantity,
	}

	if err := us.db.Create(&unit).Error; err != nil {
		return nil, errors.New("failed to create unit")
	}

	unitResponse := unit.ToResponse()

	return &unitResponse, nil
}

func (us unitService) DeleteUnit(id string) error {
	var unit models.Unit
	if err := us.db.First(&unit, id).Error; err != nil {
		return errors.New("unit not found")
	}

	if err := us.db.Delete(&unit).Error; err != nil {
		return errors.New("failed to delete unit")
	}

	return nil
}