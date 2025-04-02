package app

import (
	"github.com/armanjr/termustat/api/config"
	"github.com/armanjr/termustat/api/infrastructure/mailer"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *gin.Engine
	Mailer mailer.Mailer
	Config *config.Config
	Logger *zap.Logger
}
