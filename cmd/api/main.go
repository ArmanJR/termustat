package main

import (
	"github.com/armanjr/termustat/internal/config"
	"github.com/armanjr/termustat/internal/routes"
	"github.com/armanjr/termustat/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Database connection
	db, err := database.ConnectDB(&cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Database migration failed:", err)
	}

	// Gin setup
	router := gin.Default()
	routes.SetupRoutes(router)

	log.Fatal(router.Run(":" + cfg.Port))
}
