package main

import (
	"github.com/armanjr/termustat/app/config"
	"github.com/armanjr/termustat/app/handlers"
	"github.com/armanjr/termustat/app/routes"
	"github.com/armanjr/termustat/app/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

func initLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return config.Build()
}

func main() {
	logger, err := initLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading config", zap.Error(err))
	}

	err = config.ConnectDB(&cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := config.AutoMigrate(); err != nil {
		logger.Fatal("Database migration failed", zap.Error(err))
	}

	mailer := services.NewMailer(&cfg)

	authHandler := handlers.NewAuthHandler(mailer, &cfg, logger)

	router := gin.Default()
	routes.SetupRoutes(router, authHandler)

	logger.Info("Starting server", zap.String("port", cfg.Port))
	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
