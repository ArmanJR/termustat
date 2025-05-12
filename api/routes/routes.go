package routes

import (
	"github.com/armanjr/termustat/api/app"
	_ "github.com/armanjr/termustat/api/docs"
	"github.com/armanjr/termustat/api/handlers"
	"github.com/armanjr/termustat/api/middlewares"
	"github.com/armanjr/termustat/api/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Handlers struct {
	Auth       *handlers.AuthHandler
	University *handlers.UniversityHandler
	Professor  *handlers.ProfessorHandler
	Semester   *handlers.SemesterHandler
	Faculty    *handlers.FacultyHandler
	Course     *handlers.CourseHandler
	AdminUser  *handlers.AdminUserHandler
	UserCourse *handlers.UserCourseHandler
	Health     *handlers.HealthHandler
}

type Middlewares struct {
	JWT   *middlewares.JWTMiddleware
	Admin *middlewares.AdminMiddleware
}

func SetupRoutes(app *app.App, h *Handlers, authService services.AuthService, adminUserService services.AdminUserService, logger *zap.Logger) {
	// Initialize middlewares
	mw := &Middlewares{
		JWT:   middlewares.NewJWTMiddleware(authService, logger),
		Admin: middlewares.NewAdminMiddleware(adminUserService, logger),
	}

	// Public routes
	public := app.Router.Group("/v1")
	{
		// Health check route
		public.GET("/health", h.Health.HealthCheck)

		// Swagger UI route
		public.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL("/api/v1/swagger/doc.json"),
		))

		// User routes
		auth := public.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
			auth.POST("/forgot-password", h.Auth.ForgotPassword)
			auth.POST("/reset-password", h.Auth.ResetPassword)
			auth.POST("/verify-email", h.Auth.VerifyEmail)
			auth.POST("/refresh", h.Auth.Refresh)
			auth.POST("/logout", h.Auth.Logout)
		}
	}

	// Protected routes
	protected := app.Router.Group("/v1")
	protected.Use(mw.JWT.AuthRequired())
	{
		// User routes
		user := protected.Group("/user")
		{
			user.GET("/me", h.Auth.GetCurrentUser)
		}

		// User Course routes
		userCourses := protected.Group("/courses")
		{
			userCourses.POST("/select", h.UserCourse.AddCourse)
			userCourses.DELETE("/select/:courseId", h.UserCourse.RemoveCourse)
			userCourses.GET("/selected", h.UserCourse.GetUserCourses)
			userCourses.GET("/validate", h.UserCourse.ValidateTimeConflicts)
		}
	}

	// Admin routes
	admin := app.Router.Group("/v1/admin")
	admin.Use(mw.JWT.AuthRequired(), mw.Admin.IsAdmin())
	{
		// University routes
		universities := admin.Group("/universities")
		{
			universities.POST("", h.University.Create)
			universities.GET("", h.University.GetAll)
			universities.GET("/:id", h.University.Get)
			universities.PUT("/:id", h.University.Update)
			universities.DELETE("/:id", h.University.Delete)
			universities.GET("/:id/professors", h.Professor.GetAllByUniversity)
			universities.GET("/:id/faculties", h.Faculty.GetAllByUniversity)
			universities.GET("/:id/faculties/:short_code", h.Faculty.GetByUniversityAndShortCode)
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

		// Faculty routes
		faculties := admin.Group("/faculties")
		{
			faculties.POST("", h.Faculty.Create)
			faculties.GET("/:id", h.Faculty.GetByID)
			faculties.PUT("/:id", h.Faculty.Update)
			faculties.DELETE("/:id", h.Faculty.Delete)
			faculties.GET("/:id/courses", h.Course.GetByFaculty)
		}

		// Course routes
		courses := admin.Group("/courses")
		{
			courses.POST("", h.Course.Create)
			courses.GET("", h.Course.Search)
			courses.GET("/:id", h.Course.Get)
			courses.PUT("/:id", h.Course.Update)
			courses.DELETE("/:id", h.Course.Delete)
		}

		// Admin User routes
		users := admin.Group("/users")
		{
			users.GET("", h.AdminUser.GetAll)
			users.GET("/:id", h.AdminUser.Get)
			users.PUT("/:id", h.AdminUser.Update)
			users.DELETE("/:id", h.AdminUser.Delete)
		}
	}
}
