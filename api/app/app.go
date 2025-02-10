package app

import (
	"github.com/armanjr/termustat/api/config"
	"github.com/armanjr/termustat/api/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *gin.Engine
	Mailer services.MailerService
	Config *config.Config
	Logger *zap.Logger
}
