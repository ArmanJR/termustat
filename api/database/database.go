package database

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/armanjr/termustat/api/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func NewDatabase(config config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s timezone=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode, config.Timezone,
	)

	db, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("getting raw DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	return db, nil
}

func RunMigrations(gdb *gorm.DB, logger *zap.Logger) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		logger.Error("Failed to get raw DB", zap.Error(err))
		return fmt.Errorf("getting raw DB: %w", err)
	}
	return runMigrations(sqlDB)
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Wrap the embed.FS in the iofs driver
	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs", d,
		"postgres", driver,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
