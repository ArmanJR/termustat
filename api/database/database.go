package database

import (
	"fmt"
	"github.com/armanjr/termustat/api/config"
	"github.com/armanjr/termustat/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(config config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s timezone=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode, config.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
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
}
