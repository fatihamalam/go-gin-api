package config

import (
	"fmt"
	"go-gin-api/models"
	"go-gin-api/utils"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	host := utils.GetEnv("DB_HOST", "localhost")
	user := utils.GetEnv("DB_USER", "postgres")
	// password := utils.GetEnv("DB_PASSWORD", "")
	dbname := utils.GetEnv("DB_NAME", "go_gin_api")
	port := utils.GetEnv("DB_PORT", "5432")
	
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable", host, user, dbname, port)
	log.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully.")

	err = db.AutoMigrate(
		&models.Permission{},
		&models.Role{},
		&models.User{},
		&models.ItemCategory{},
		&models.Unit{},
		&models.Item{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedInitialData(db)

	return db
}

func seedInitialData(db *gorm.DB) {
	permissions := []models.Permission{
		{Name: "user:read", Description: "Can read user information"},
		{Name: "user:write", Description: "Can create and update users"},
		{Name: "user:delete", Description: "Can delete users"},
		{Name: "role:read", Description: "Can read role information"},
		{Name: "role:write", Description: "Can create and update roles"},
		{Name: "role:delete", Description: "Can delete roles"},
	}

	for _, p := range permissions {
		var permission models.Permission
		if result := db.Where("name = ?", p.Name).First(&permission); result.RowsAffected == 0 {
			db.Create(&p)
		}
	}

	var adminRole models.Role
	if result := db.Where("name = ?", "admin").First(&adminRole); result.RowsAffected == 0 {
		adminRole = models.Role{
			Name: "admin",
			Description: "Administrator",
		}
		db.Create(&adminRole)

		var allPermissions []models.Permission
		db.Find(&allPermissions)
		db.Model(&adminRole).Association("Permissions").Append(allPermissions)
	}

	var userRole models.Role
	if result := db.Where("name = ?", "user").First(&userRole); result.RowsAffected == 0 {
		userRole = models.Role{
			Name: "user",
			Description: "User",
		}
		db.Create(&userRole)

		var userReadPermission []models.Permission
		db.Where("name = ?", "user:read").Find(&userReadPermission)
		db.Model(&userRole).Association("Permissions").Append(userReadPermission)
	}

	var adminUser models.User
	if result := db.Where("email = ?", "admin@example.com").First(&adminUser); result.RowsAffected == 0 {
		adminUser = models.User{
			Name: "Admin",
			Email: "admin@example.com",
			Password: "admin123",
			RoleID: adminRole.ID,
		}
		adminUser.HashPassword()
		db.Create(&adminUser)
	}
}
