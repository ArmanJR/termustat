package database

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

func RunMigrationsWithPathResolution(db *gorm.DB, logger *zap.Logger) error {
	paths := []string{
		"api/database/migrations",
		"database/migrations",
		"/api/database/migrations",
	}

	var lastErr error
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			logger.Info("Running SQL migrations via golang-migrate", zap.String("path", p))
			if err := RunMigrations(db, p); err != nil {
				logger.Warn("Migrations failed", zap.String("path", p), zap.Error(err))
				lastErr = err
			} else {
				return nil
			}
		}
	}

	return fmt.Errorf("all migration paths failed: %w", lastErr)
}
