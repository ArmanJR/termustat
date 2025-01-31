package config

import (
	"github.com/spf13/viper"
)

type Config struct {
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

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
