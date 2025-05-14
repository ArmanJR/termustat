package main

import (
	"context"
	"fmt"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/armanjr/termustat/api/app"
	"github.com/armanjr/termustat/api/config"
	"github.com/armanjr/termustat/api/database"
	"github.com/armanjr/termustat/api/handlers"
	"github.com/armanjr/termustat/api/infrastructure/mailer"
	"github.com/armanjr/termustat/api/logger"
	"github.com/armanjr/termustat/api/repositories"
	"github.com/armanjr/termustat/api/routes"
	"github.com/armanjr/termustat/api/services"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title        Termustat API
// @version      1.0
// @host         localhost:8080
// @BasePath     /api
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

	// Third-party services
	mailerConfig := mailer.MailerConfig{
		Domain:  cfg.MailgunDomain,
		APIKey:  cfg.MailgunAPIKey,
		Sender:  "noreply@" + cfg.MailgunDomain,
		TplPath: "templates/email/",
	}
	mailerService := mailer.NewMailer(mailerConfig, log)

	// Initialize database
	db, err := database.NewDatabase(cfg.GetDatabaseConfig())
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run SQL migrations
	if err := database.RunMigrations(db, log); err != nil {
		log.Fatal("SQL migrations failed", zap.Error(err))
	}

	// Initialize repositories
	authRepo := repositories.NewAuthRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)
	professorRepo := repositories.NewProfessorRepository(db)
	universityRepo := repositories.NewUniversityRepository(db)
	semesterRepo := repositories.NewSemesterRepository(db)
	facultyRepo := repositories.NewFacultyRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	adminUserRepo := repositories.NewAdminUserRepository(db)
	userCourseRepo := repositories.NewUserCourseRepository(db)

	// Internal services
	authService := services.NewAuthService(
		authRepo,
		refreshTokenRepo,
		mailerService,
		log,
		cfg.JWTSecret,
		cfg.JWTTTL,
		cfg.RefreshTTL,
		cfg.FrontendURL,
	)
	universityService := services.NewUniversityService(universityRepo, log)
	professorService := services.NewProfessorService(professorRepo, universityService, log)
	semesterService := services.NewSemesterService(semesterRepo, log)
	facultyService := services.NewFacultyService(facultyRepo, universityService, log)
	courseService := services.NewCourseService(courseRepo, universityService, facultyService, professorService, semesterService, log)
	adminUserService := services.NewAdminUserService(adminUserRepo, universityService, facultyService, log)
	userCourseService := services.NewUserCourseService(userCourseRepo, courseService, adminUserService, semesterService, log)

	// Initialize router
	router := gin.New()

	// Setup middleware
	router.Use(ginzap.Ginzap(log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(log, true))

	// Allow frontend CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			cfg.FrontendURL,
			cfg.NginxURL,
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize application
	application := &app.App{
		DB:     db,
		Router: router,
		Config: cfg,
		Logger: log,
	}

	// Initialize handlers
	ginHandlers := &routes.Handlers{
		Auth:       handlers.NewAuthHandler(authService, universityService, facultyService, log),
		Professor:  handlers.NewProfessorHandler(professorService, log),
		University: handlers.NewUniversityHandler(universityService, log),
		Semester:   handlers.NewSemesterHandler(semesterService, log),
		Faculty:    handlers.NewFacultyHandler(facultyService, log),
		Course:     handlers.NewCourseHandler(courseService, log),
		AdminUser:  handlers.NewAdminUserHandler(adminUserService, log),
		UserCourse: handlers.NewUserCourseHandler(userCourseService, log),
		Health:     handlers.NewHealthHandler(log),
	}

	// Setup routes
	routes.SetupRoutes(
		application,
		ginHandlers,
		authService,
		adminUserService,
		log,
	)

	// Setup server
	serverAddr := ":" + cfg.Port
	log.Info("Starting server",
		zap.String("port", cfg.Port),
		zap.String("environment", cfg.Environment),
		zap.String("timezone", cfg.Timezone),
	)

	// Setup signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatal("Failed to start server",
				zap.String("address", serverAddr),
				zap.Error(err),
			)
		}
	}()

	// Wait for interrupt signal
	<-quit

	// Handle graceful shutdown
	gracefulShutdown(application)
}

func gracefulShutdown(app *app.App) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.Logger.Info("Initiating graceful shutdown...")

	// Close database connection
	if app.DB != nil {
		sqlDB, err := app.DB.DB()
		if err != nil {
			app.Logger.Error("Error getting underlying SQL DB instance", zap.Error(err))
		} else {
			if err := sqlDB.Close(); err != nil {
				app.Logger.Error("Error closing database connection", zap.Error(err))
			} else {
				app.Logger.Info("Database connection closed successfully")
			}
		}
	}

	// Shutdown the HTTP server
	if app.Router != nil {
		srv := &http.Server{
			Addr:    ":" + app.Config.Port,
			Handler: app.Router,
		}

		// Shutdown the server with context timeout
		if err := srv.Shutdown(ctx); err != nil {
			app.Logger.Error("Server forced to shutdown", zap.Error(err))
		} else {
			app.Logger.Info("Server shutdown completed successfully")
		}
	}

	if err := app.Logger.Sync(); err != nil {
		fmt.Printf("Error flushing logs: %v\n", err)
	}

	app.Logger.Info("Graceful shutdown completed")
}
