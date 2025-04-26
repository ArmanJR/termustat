package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	// Environment
	Environment string `mapstructure:"ENVIRONMENT"`
	Port        string `mapstructure:"PORT"`
	Timezone    string `mapstructure:"TIMEZONE"`

	// Database
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	SSLMode    string `mapstructure:"SSL_MODE"`

	// JWT
	JWTSecret  string        `mapstructure:"JWT_SECRET"`
	JWTTTL     time.Duration `mapstructure:"JWT_TTL"`
	RefreshTTL time.Duration `mapstructure:"REFRESH_TTL"`

	// Mailgun
	MailgunAPIKey string `mapstructure:"MAILGUN_API_KEY"`
	MailgunDomain string `mapstructure:"MAILGUN_DOMAIN"`

	// Frontend
	FrontendURL string `mapstructure:"FRONTEND_URL"`

	// Nginx
	NginxURL string `mapstructure:"NGINX_URL"`
}

// DatabaseConfig Database configuration struct
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timezone string
}

// GetDatabaseConfig returns database configuration
func (c *Config) GetDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     c.DBHost,
		Port:     c.DBPort,
		User:     c.DBUser,
		Password: c.DBPassword,
		DBName:   c.DBName,
		SSLMode:  c.SSLMode,
		Timezone: c.Timezone,
	}
}

// LoadConfig loads configuration from environment file
func LoadConfig(configPath string) (*Config, error) {
	var config Config

	viper.AddConfigPath(configPath)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to load configs: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configs: %w", err)
	}

	// Set defaults if not provided
	if config.Environment == "" {
		config.Environment = "development"
	}

	if config.JWTTTL == 0 {
		config.JWTTTL = 48 * time.Hour // Default to 48 hours
	}

	if config.RefreshTTL == 0 {
		config.RefreshTTL = 720 * time.Hour // Default to 30 days
	}

	// Validate required fields
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig validates required configuration fields
func validateConfig(config *Config) error {
	if config.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if config.DBHost == "" || config.DBName == "" {
		return fmt.Errorf("database configuration is incomplete")
	}

	if config.MailgunAPIKey == "" || config.MailgunDomain == "" {
		return fmt.Errorf("mailgun configuration is incomplete")
	}

	return nil
}
