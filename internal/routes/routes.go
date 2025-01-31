package routes

import (
	"github.com/gin-gonic/gin"
	"internal/handlers"
	"internal/middleware"
)

func SetupRoutes(router *gin.Engine) {
	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", handlers.Register)
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/forgot-password", handlers.ForgotPassword)
		public.POST("/auth/reset-password", handlers.ResetPassword)

		public.GET("/universities", handlers.GetUniversities)
		public.GET("/universities/:id/faculties", handlers.GetFaculties)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/users/me", handlers.GetCurrentUser)
		protected.POST("/users/me/courses", handlers.AddCourse)
		protected.DELETE("/users/me/courses/:id", handlers.RemoveCourse)
		protected.GET("/users/me/timetable", handlers.GetTimetable)
		protected.POST("/users/me/timetable", handlers.SaveTimetable)
		protected.GET("/users/me/courses", handlers.GetUserCourses)
	}
}
