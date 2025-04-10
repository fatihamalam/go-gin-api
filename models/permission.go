package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name 		string	`json:"name" gorm:"uniqueIndex;not null"`
	Description	string	`json:"description"`
	Roles		[]Role	`json:"-" gorm:"many2many:role_permissions;"`
}