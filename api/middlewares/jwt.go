package middlewares

import (
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/services"
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			return
		}

		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			switch {
			case errors.Is(err, errors.ErrExpiredToken):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			default:
				m.logger.Warn("Invalid token",
					zap.Error(err),
					zap.String("token", tokenString))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			return
		}

		// Store user ID in context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
