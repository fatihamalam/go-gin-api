package config

import (
	"fmt"
	"go-gin-api/models"
	"go-gin-api/utils"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func SetupValidator(db *gorm.DB) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("uniqueEmail", func(fl validator.FieldLevel) bool {
			email := fl.Field().String()
			var user models.User
			query := db.Where("email = ?", email)
			if utils.GinCtx != nil {
				currentID := utils.GinCtx.Param("id")
				fmt.Printf("currentID: %s", currentID)
				if currentID != "" {
					query.Where("id != ?", currentID)
				}
			}
			if err := query.First(&user).Error; err != nil {
				return err == gorm.ErrRecordNotFound
			}
			return false
		})
	}
}