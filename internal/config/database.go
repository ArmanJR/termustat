package config

import (
	"fmt"
	"github.com/armanjr/termustat/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.University{},
		&models.Faculty{},
		&models.User{},
		&models.Course{},
		&models.CourseTime{},
		&models.UserCourse{},
		&models.PasswordReset{},
		&models.Professor{},
		&models.Semester{},
	)
}
