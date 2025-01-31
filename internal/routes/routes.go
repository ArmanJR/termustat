package routes

import (
	"github.com/armanjr/termustat/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/forgot-password", authHandler.ForgotPassword)
		public.POST("/auth/reset-password", authHandler.ResetPassword)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(authHandler.JWTAuthMiddleware())
	{
		protected.GET("/users/me", authHandler.GetCurrentUser)
	}
}
