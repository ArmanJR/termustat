package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Timezone string `mapstructure:"TIMEZONE"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	SSLMode    string `mapstructure:"SSL_MODE"`

	JWTSecret string `mapstructure:"JWT_SECRET"`
	JWTTTL    int    `mapstructure:"JWT_TTL"`
	Port      string `mapstructure:"PORT"`

	MailgunAPIKey string `mapstructure:"MAILGUN_API_KEY"`
	MailgunDomain string `mapstructure:"MAILGUN_DOMAIN"`
	FrontendURL   string `mapstructure:"FRONTEND_URL"`
}

var Cfg Config

func LoadConfig() {
	var config Config

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load Configs: %v", zap.Error(err)))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Configs: %v", zap.Error(err)))
	}

	Cfg = config
}
