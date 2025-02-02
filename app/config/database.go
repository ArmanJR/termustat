package config

import (
	"fmt"
	"github.com/armanjr/termustat/app/logger"
	"github.com/armanjr/termustat/app/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		Cfg.DBHost, Cfg.DBUser, Cfg.DBPassword, Cfg.DBName, Cfg.DBPort, Cfg.SSLMode, Cfg.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	DB = db
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.University{},
		&models.Faculty{},
		&models.User{},
		&models.EmailVerification{},
		&models.Course{},
		&models.CourseTime{},
		&models.UserCourse{},
		&models.PasswordReset{},
		&models.Professor{},
		&models.Semester{},
	)
	if err != nil {
		logger.Log.Fatal("Database migration failed", zap.Error(err))
	}
}
