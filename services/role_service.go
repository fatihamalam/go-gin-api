package services

import (
	"errors"
	"go-gin-api/dto"
	"go-gin-api/models"

	"gorm.io/gorm"
)

type RoleService interface {
	FindAll(search string, order string, limit int, offset int) (*[]models.RoleResponse, *int64, error)
	FindOneByID(roleID string) (*models.RoleResponse, error)
	CreateRole(input dto.CreateRoleInput) (*models.RoleResponse, error)
	UpdateRole(roleID string, input dto.UpdateRoleInput) (*models.RoleResponse, error)
	DeleteRole(roleID string) error
}

type roleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) RoleService {
	return &roleService{
		db: db,
	}
}

func (rs roleService) FindAll(search string, order string, limit int, offset int) (*[]models.RoleResponse, *int64, error) {
	var roles []models.Role
	var totalData int64

	query := rs.db.Model(&models.Role{})
	query.Order(order)

	if search != "" {
		likeSearch := "%" + search + "%"
		query.Where("name ILIKE ?", likeSearch)
	}

	result := query.Find(&roles).Limit(limit).Offset(offset)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	result = query.Count(&totalData)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	roleResponse := models.MapRolesToResponse(roles)

	return &roleResponse, &totalData, nil
}

func (rs roleService) FindOneByID(roleID string) (*models.RoleResponse, error) {
	var role models.Role
	if err := rs.db.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, errors.New("role not found")
	}

	roleResponse := role.ToResponse()

	return &roleResponse, nil
}

func (rs roleService) CreateRole(input dto.CreateRoleInput) (*models.RoleResponse, error) {
	role := models.Role{
		Name: input.Name,
		Description: input.Description,
	}

	if err := rs.db.Create(&role).Error; err != nil {
		return nil, errors.New("failed to create role")
	}

	roleResponse := role.ToResponse()

	return &roleResponse, nil
}

func (rs roleService) UpdateRole(roleID string, input dto.UpdateRoleInput) (*models.RoleResponse, error) {
	var role models.Role
	if err := rs.db.First(&role, roleID).Error; err != nil {
		return nil, errors.New("role not found")
	}

	role.Name = input.Name
	role.Description = input.Description

	if err := rs.db.Save(&role).Error; err != nil {
		return nil, errors.New("failed to update role")
	}

	roleResponse := role.ToResponse()

	return &roleResponse, nil
}

func (rs roleService) DeleteRole(roleID string) error {
	var role models.Role
	if err := rs.db.First(&role, roleID).Error; err != nil {
		return errors.New("role not found")
	}

	if err := rs.db.Delete(&role).Error; err != nil {
		return errors.New("failed to delete role")
	}

	return nil
}