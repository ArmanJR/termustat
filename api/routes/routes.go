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
	Semester   *handlers.SemesterHandler
}

func SetupRoutes(app *app.App, h *Handlers) {
	// Public routes
	public := app.Router.Group("/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
			auth.POST("/forgot-password", h.Auth.ForgotPassword)
			auth.POST("/reset-password", h.Auth.ResetPassword)
			auth.POST("/verify-email", h.Auth.VerifyEmail)
		}
	}

	// Protected routes
	protected := app.Router.Group("/v1")
	protected.Use(middlewares.JWTAuthMiddleware())
	{
		users := protected.Group("/users")
		{
			users.GET("/me", h.Auth.GetCurrentUser)
		}
	}

	// Admin routes
	admin := app.Router.Group("/v1/admin")
	admin.Use(middlewares.JWTAuthMiddleware(), middlewares.IsAdminMiddleware())
	{
		// University routes
		universities := admin.Group("/universities")
		{
			universities.POST("", h.University.Create)
			universities.GET("", h.University.GetAll)
			universities.GET("/:id", h.University.Get)
			universities.PUT("/:id", h.University.Update)
			universities.DELETE("/:id", h.University.Delete)
			universities.GET("/:id/professors", h.Professor.GetByUniversity)
		}

		// Professor routes
		professors := admin.Group("/professors")
		{
			professors.POST("", h.Professor.Create)
			professors.GET("/:id", h.Professor.Get)
		}

		// Semester routes
		semesters := admin.Group("/semesters")
		{
			semesters.POST("", h.Semester.Create)
			semesters.GET("", h.Semester.GetAll)
			semesters.GET("/:id", h.Semester.Get)
			semesters.PUT("/:id", h.Semester.Update)
			semesters.DELETE("/:id", h.Semester.Delete)
		}
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
		//admin.GET("/universities/:id/professors", handlers.GetAllByUniversity)

		// Faculty
		admin.POST("/faculties", handlers.CreateFaculty)
		admin.GET("/faculties/:id", handlers.GetFaculty)
		admin.PUT("/faculties/:id", handlers.UpdateFaculty)
		admin.DELETE("/faculties/:id", handlers.DeleteFaculty)

		// Semester
		//admin.POST("/semesters", handlers.CreateSemester)
		//admin.GET("/semesters", handlers.GetAllSemesters)
		//admin.GET("/semesters/:id", handlers.GetSemester)
		//admin.PUT("/semesters/:id", handlers.UpdateSemester)
		//admin.DELETE("/semesters/:id", handlers.DeleteSemester)

		// Professor
		//admin.GET("/professors/:id", handlers.GetProfessor)
		//admin.PUT("/professors/:id", handlers.UpdateProfessor) // do we need this?

		// Course
		admin.POST("/courses", handlers.CreateCourse)
		admin.POST("/courses/batch", handlers.BatchCreateCourses)
		admin.GET("/courses/:id", handlers.GetCourse)
		admin.PUT("/courses/:id", handlers.UpdateCourse)
		admin.DELETE("/courses/:id", handlers.DeleteCourse)
	}
}
