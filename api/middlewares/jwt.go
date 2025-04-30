package middlewares

import (
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/services"
	"github.com/armanjr/termustat/api/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type JWTMiddleware struct {
	authService services.AuthService
	logger      *zap.Logger
}

func NewJWTMiddleware(authService services.AuthService, logger *zap.Logger) *JWTMiddleware {
	return &JWTMiddleware{
		authService: authService,
		logger:      logger,
	}
}

func (m *JWTMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Debug("Authorization header missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			m.logger.Debug("Bearer prefix missing in Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			return
		}

		// Use ValidateToken which now internally uses ParseJWT
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			status := http.StatusUnauthorized
			errMsg := "Invalid token"
			if errors.Is(err, utils.ErrExpiredToken) {
				errMsg = "Token expired"
				m.logger.Info("Token expired", zap.String("token_prefix", tokenString[:min(10, len(tokenString))]))
			} else {
				m.logger.Warn("Token validation failed", zap.Error(err), zap.String("token_prefix", tokenString[:min(10, len(tokenString))]))
			}
			c.AbortWithStatusJSON(status, gin.H{"error": errMsg})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userScopes", claims.Scopes)
		//m.logger.Debug("Token validated successfully", zap.String("userID", claims.UserID), zap.Strings("scopes", claims.Scopes))
		c.Next()
	}
}
