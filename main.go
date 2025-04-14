package main

import (
	"go-gin-api/config"
	"go-gin-api/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}
	
	db := config.ConnectDB()
	dbcon, err := db.DB()
	if err != nil {
		panic("Failed to close database connection.")
	}
	defer dbcon.Close()

	config.SetupValidator(db)
	
	r := gin.Default()
	apiRoute := r.Group("/api")
	routes.SetupAuthRoutes(apiRoute, db)
	routes.SetupUserRoutes(apiRoute, db)
	r.Run(":8080")
}