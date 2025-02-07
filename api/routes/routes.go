package routes

import (
	"github.com/armanjr/termustat/api/app"
	"github.com/armanjr/termustat/api/handlers"
	"github.com/armanjr/termustat/api/middlewares"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth       *handlers.AuthHandler
	University *handlers.UniversityHandler
	Professor  *handlers.ProfessorHandler
}

func SetupRoutes(app *app.App, h *Handlers) {
	// Public routes
	public := app.Router.Group("/v1")
	{
		public.POST("/auth/register", h.Auth.Register)
		public.POST("/auth/login", h.Auth.Login)
		public.POST("/auth/forgot-password", h.Auth.ForgotPassword)
		public.POST("/auth/reset-password", h.Auth.ResetPassword)
		public.POST("/auth/verify-email", h.Auth.VerifyEmail)
	}

	// Protected routes
	protected := app.Router.Group("/v1")
	protected.Use(middlewares.JWTAuthMiddleware())
	{
		protected.GET("/users/me", h.Auth.GetCurrentUser)
	}

	// Admin routes
	admin := app.Router.Group("/v1/admin")
	admin.Use(middlewares.JWTAuthMiddleware(), middlewares.IsAdminMiddleware())
	{
		// University routes
		admin.POST("/universities", h.University.Create)
		admin.GET("/universities", h.University.GetAll)
		admin.GET("/universities/:id", h.University.Get)
		admin.PUT("/universities/:id", h.University.Update)
		admin.DELETE("/universities/:id", h.University.Delete)
		admin.DELETE("/universities/:id/professors", h.Professor.GetByUniversity)

		// Professor routes
		admin.GET("/professors/:id", h.Professor.Get)
	}
}

func SetupRoutesLegacy(router *gin.Engine) {
	// an /api/ is already added by nginx
	// Public routes
	//public := router.Group("/v1")
	//{
	//	public.POST("/auth/register", handlers.Register)
	//	public.POST("/auth/login", handlers.Login)
	//	public.POST("/auth/forgot-password", handlers.ForgotPassword)
	//	public.POST("/auth/reset-password", handlers.ResetPassword)
	//}

	// Protected routes
	//protected := router.Group("/v1")
	//protected.Use(middlewares.JWTAuthMiddleware())
	//{
	//	protected.GET("/users/me", handlers.GetCurrentUser)
	//}

	// Admin routes
	admin := router.Group("/v1/admin")
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
		//admin.POST("/universities", handlers.CreateUniversity)
		//admin.GET("/universities", handlers.GetAllUniversities)
		//admin.GET("/universities/:id", handlers.GetUniversity)
		//admin.PUT("/universities/:id", handlers.UpdateUniversity)
		//admin.DELETE("/universities/:id", handlers.DeleteUniversity)
		//admin.GET("/universities/:id/professors", handlers.GetProfessorsByUniversity)

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
		//admin.PUT("/professors/:id", handlers.UpdateProfessor) // do we need this?

		// Course
		admin.POST("/courses", handlers.CreateCourse)
		admin.POST("/courses/batch", handlers.BatchCreateCourses)
		admin.GET("/courses/:id", handlers.GetCourse)
		admin.PUT("/courses/:id", handlers.UpdateCourse)
		admin.DELETE("/courses/:id", handlers.DeleteCourse)
	}
}
