package routes

import (
	"github.com/armanjr/termustat/app/handlers"
	"github.com/armanjr/termustat/app/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", handlers.Register)
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/forgot-password", handlers.ForgotPassword)
		public.POST("/auth/reset-password", handlers.ResetPassword)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middlewares.JWTAuthMiddleware())
	{
		protected.GET("/users/me", handlers.GetCurrentUser)
	}

	// Admin routes
	admin := router.Group("/api/v1/admin")
	admin.Use(middlewares.JWTAuthMiddleware(), middlewares.IsAdminMiddleware())
	{
		// University
		admin.POST("/universities", handlers.CreateUniversity)
		admin.GET("/universities", handlers.GetAllUniversities)
		admin.GET("/universities/:id", handlers.GetUniversity)
		admin.PUT("/universities/:id", handlers.UpdateUniversity)
		admin.DELETE("/universities/:id", handlers.DeleteUniversity)

		// Faculty
		admin.POST("/faculties", handlers.CreateFaculty)
		admin.GET("/faculties/:id", handlers.GetFaculty)
		admin.PUT("/faculties/:id", handlers.UpdateFaculty)
		admin.DELETE("/faculties/:id", handlers.DeleteFaculty)

		// Semester
		admin.POST("/semesters", handlers.CreateSemester)
		admin.GET("/semesters", handlers.GetAllSemesters)
		admin.GET("/semesters/:id", handlers.GetSemester)
		admin.PUT("/semesters/:id", handlers.UpdateSemester)
		admin.DELETE("/semesters/:id", handlers.DeleteSemester)
	}
}
