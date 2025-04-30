package middlewares

import (
	"github.com/armanjr/termustat/api/services"
	"github.com/armanjr/termustat/api/utils"
	"github.com/gin-gonic/gin"
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
		// Check 1: Ensure user is authenticated (should be guaranteed by JWTMiddleware running first)
		userID, userExists := c.Get("userID")
		if !userExists {
			m.logger.Error("IsAdmin middleware called without userID in context. Check middleware order.")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Check 2: Get scopes from context set by JWTMiddleware
		scopesVal, scopesExist := c.Get("userScopes")
		if !scopesExist {
			m.logger.Warn("No scopes found in context for authenticated user", zap.Any("userID", userID))
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		// Check 3: Assert type and check for the required scope
		scopes, ok := scopesVal.([]string)
		if !ok {
			m.logger.Error("User scopes in context are not of type []string", zap.Any("userID", userID), zap.Any("scopesType", scopesVal))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Check 4: Verify the "admin-dashboard" scope exists
		requiredScope := "admin-dashboard"
		if !utils.ContainsScope(scopes, requiredScope) {
			m.logger.Warn("Admin access denied: missing required scope",
				zap.Any("userID", userID),
				zap.Strings("userScopes", scopes),
				zap.String("requiredScope", requiredScope))
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		//m.logger.Debug("Admin access granted", zap.Any("userID", userID), zap.Strings("scopes", scopes))
		c.Next()
	}
}
