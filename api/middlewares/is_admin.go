package middlewares

import (
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type AdminMiddleware struct {
	adminUserService services.AdminUserService
	logger           *zap.Logger
}

func NewAdminMiddleware(adminUserService services.AdminUserService, logger *zap.Logger) *AdminMiddleware {
	return &AdminMiddleware{
		adminUserService: adminUserService,
		logger:           logger,
	}
}

func (m *AdminMiddleware) IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		parsedID, err := uuid.Parse(userID.(string))
		if err != nil {
			m.logger.Warn("Invalid user ID format in middleware",
				zap.String("user_id", userID.(string)))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		user, err := m.adminUserService.Get(parsedID)
		if err != nil {
			switch {
			case errors.Is(err, errors.ErrNotFound):
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			default:
				m.logger.Error("Failed to fetch user in admin middleware",
					zap.String("user_id", parsedID.String()),
					zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			c.Abort()
			return
		}

		if !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
