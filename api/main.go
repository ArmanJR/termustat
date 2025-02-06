package main

import (
	"github.com/armanjr/termustat/api/app"
	"github.com/armanjr/termustat/api/config"
	"github.com/armanjr/termustat/api/database"
	"github.com/armanjr/termustat/api/handlers"
	"github.com/armanjr/termustat/api/logger"
	"github.com/armanjr/termustat/api/repositories"
	"github.com/armanjr/termustat/api/routes"
	"github.com/armanjr/termustat/api/services"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Set application timezone
	if err := os.Setenv("TZ", cfg.Timezone); err != nil {
		log.Fatal("Failed to set timezone:", err)
	}

	// Initialize logger
	logger.InitLogger()
	defer logger.Log.Sync()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.GetDatabaseConfig())
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := database.AutoMigrate(db); err != nil {
		logger.Log.Fatal("Database migration failed", zap.Error(err))
	}

	// Initialize repositories
	authRepo := repositories.NewAuthRepository(db)
	//universityRepo := repositories.NewUniversityRepository(db)
	// other repositories

	// Initialize services
	mailerConfig := services.MailerConfig{
		Domain:  cfg.MailgunDomain,
		APIKey:  cfg.MailgunAPIKey,
		Sender:  "noreply@" + cfg.MailgunDomain,
		TplPath: "templates/email/",
	}
	mailerService := services.NewMailerService(mailerConfig, logger.Log)

	authService := services.NewAuthService(
		authRepo,
		mailerService,
		logger.Log,
		cfg.JWTSecret,
		cfg.JWTTTL,
		cfg.FrontendURL,
	)

	//universityService := services.NewUniversityService()

	// Initialize router
	router := gin.New()

	// Setup middleware
	router.Use(ginzap.Ginzap(logger.Log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger.Log, true))

	// Initialize application
	application := &app.App{
		DB:     db,
		Router: router,
		Config: cfg,
		Logger: logger.Log,
	}

	// Initialize handlers
	ginHandlers := &routes.Handlers{
		Auth: handlers.NewAuthHandler(authService, logger.Log),
		//University: routes.NewUniversityHandler(universityRepo),
		// Add other handlers as needed
	}

	// Setup routes
	routes.SetupRoutes(application, ginHandlers)

	// Start server
	serverAddr := ":" + cfg.Port
	logger.Log.Info("Starting server",
		zap.String("port", cfg.Port),
		zap.String("environment", cfg.Environment),
		zap.String("timezone", cfg.Timezone),
	)

	if err := router.Run(serverAddr); err != nil {
		logger.Log.Fatal("Failed to start server",
			zap.String("address", serverAddr),
			zap.Error(err),
		)
	}
}

// gracefulShutdown handles graceful shutdown of the server
func gracefulShutdown(app *app.App) {}
