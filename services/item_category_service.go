package services

import (
	"errors"
	"go-gin-api/dto"
	"go-gin-api/models"

	"gorm.io/gorm"
)

type ItemCategoryService interface {
	FindAll(search string, order string, limit int, offset int, includes []string) (*[]models.ItemCategoryResponse, *int64, error)
	CreateItemCategory(input dto.CreateItemCategoryInput) (*models.ItemCategoryResponse, error)
	DeleteItemCategory(id string) error
}

type itemCategoryService struct {
	db *gorm.DB
}

func NewItemCategoryService(db *gorm.DB) ItemCategoryService {
	return &itemCategoryService{
		db: db,
	}
}

func (ics itemCategoryService) FindAll(search string, order string, limit int, offset int, includes []string) (*[]models.ItemCategoryResponse, *int64, error) {
	var itemCategories []models.ItemCategory
	var totalData int64

	query := ics.db.Model(&models.ItemCategory{})
	query.Order(order)

	if search != "" {
		likeSearch := "%" + search + "%"
		query.Where("name ILIKE ?", likeSearch)
	}

	for _, include := range includes {
		query.Preload(include)
	}

	result := query.Find(&itemCategories).Limit(limit).Offset(offset)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	result = query.Count(&totalData)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	itemCategoryResponse := models.MapItemCategoryToResponse(itemCategories)

	return &itemCategoryResponse, &totalData, nil
}

func (ics itemCategoryService) CreateItemCategory(input dto.CreateItemCategoryInput) (*models.ItemCategoryResponse, error) {
	category := models.ItemCategory{
		Name: input.Name,
		Description: input.Description,
	}

	if err := ics.db.Create(&category).Error; err != nil {
		return nil, errors.New("failed to create item category")
	}

	itemCategoryResponse := category.ToResponse()

	return &itemCategoryResponse, nil
}

func (ics itemCategoryService) DeleteItemCategory(id string) error {
	var category models.ItemCategory
	if err := ics.db.First(&category, id).Error; err != nil {
		return errors.New("item category not found")
	}

	if err := ics.db.Delete(&category).Error; err != nil {
		return errors.New("failed to delete item category")
	}

	return nil
}