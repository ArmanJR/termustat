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
	stdLog "log"
	"os"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		stdLog.Fatal("Failed to load configuration:", err)
	}

	// Set application timezone
	if err := os.Setenv("TZ", cfg.Timezone); err != nil {
		stdLog.Fatal("Failed to set timezone:", err)
	}

	// Initialize logger
	log := logger.NewLogger()
	defer log.Sync()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.GetDatabaseConfig())
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Database migration failed", zap.Error(err))
	}

	// Initialize repositories
	authRepo := repositories.NewAuthRepository(db)
	professorRepo := repositories.NewProfessorRepository(db)
	//universityRepo := repositories.NewUniversityRepository(db)
	// other repositories

	// Third-party services
	mailerConfig := services.MailerConfig{
		Domain:  cfg.MailgunDomain,
		APIKey:  cfg.MailgunAPIKey,
		Sender:  "noreply@" + cfg.MailgunDomain,
		TplPath: "templates/email/",
	}
	mailerService := services.NewMailerService(mailerConfig, log)

	// Internal services
	authService := services.NewAuthService(
		authRepo,
		mailerService,
		log,
		cfg.JWTSecret,
		cfg.JWTTTL,
		cfg.FrontendURL,
	)
	professorService := services.NewProfessorService(professorRepo, log)

	//universityService := services.NewUniversityService()

	// Initialize router
	router := gin.New()

	// Setup middleware
	router.Use(ginzap.Ginzap(log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(log, true))

	// Initialize application
	application := &app.App{
		DB:     db,
		Router: router,
		Config: cfg,
		Logger: log,
	}

	// Initialize handlers
	ginHandlers := &routes.Handlers{
		Auth:      handlers.NewAuthHandler(authService, log),
		Professor: handlers.NewProfessorHandler(professorService, log),
		//University: routes.NewUniversityHandler(universityRepo),
		// Add other handlers as needed
	}

	// Setup routes
	routes.SetupRoutes(application, ginHandlers)

	// Start server
	serverAddr := ":" + cfg.Port
	log.Info("Starting server",
		zap.String("port", cfg.Port),
		zap.String("environment", cfg.Environment),
		zap.String("timezone", cfg.Timezone),
	)

	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server",
			zap.String("address", serverAddr),
			zap.Error(err),
		)
	}
}

// gracefulShutdown handles graceful shutdown of the server
func gracefulShutdown(app *app.App) {}
