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
		// User
		admin.GET("/users", handlers.GetAllUsers)
		admin.GET("/users/:id", handlers.GetUser)
		admin.PUT("/users/:id", handlers.UpdateUser)
		admin.DELETE("/users/:id", handlers.DeleteUser)

		admin.GET("/users/:id/courses", handlers.GetUserCourses)
		admin.POST("/users/:id/courses", handlers.AddUserCourse)
		admin.DELETE("/users/:id/courses/:course_id", handlers.RemoveUserCourse)

		// University
		admin.POST("/universities", handlers.CreateUniversity)
		admin.GET("/universities", handlers.GetAllUniversities)
		admin.GET("/universities/:id", handlers.GetUniversity)
		admin.PUT("/universities/:id", handlers.UpdateUniversity)
		admin.DELETE("/universities/:id", handlers.DeleteUniversity)
		admin.GET("/universities/:id/professors", handlers.GetProfessorsByUniversity)

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

		// Professor
		admin.GET("/professors/:id", handlers.GetProfessor)
		admin.PUT("/professors/:id", handlers.UpdateProfessor)
	}
}
