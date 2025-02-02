package main

import (
	"github.com/armanjr/termustat/app/config"
	"github.com/armanjr/termustat/app/logger"
	"github.com/armanjr/termustat/app/routes"
	"github.com/armanjr/termustat/app/services"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

func main() {
	// Configs
	config.LoadConfig()

	// Application Timezone
	if err := os.Setenv("TZ", config.Cfg.Timezone); err != nil {
		log.Fatal("Failed to set timezone:", err)
	}

	// Logger
	logger.InitLogger()
	defer logger.Log.Sync()

	// Database
	config.ConnectDB()
	config.AutoMigrate()

	// Services
	services.RegisterMailer()

	// HTTP
	router := gin.New()
	router.Use(ginzap.Ginzap(logger.Log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger.Log, true))
	routes.SetupRoutes(router)
	logger.Log.Info("Starting server", zap.String("port", config.Cfg.Port))
	if err := router.Run(":" + config.Cfg.Port); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}
