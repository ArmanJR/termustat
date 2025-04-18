package handlers

import (
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type UserCourseHandler struct {
	service services.UserCourseService
	logger  *zap.Logger
}

func NewUserCourseHandler(service services.UserCourseService, logger *zap.Logger) *UserCourseHandler {
	return &UserCourseHandler{
		service: service,
		logger:  logger,
	}
}

// AddCourse handles course selection
// @Summary      Select Course
// @Description  Adds a course to the current user's schedule
// @Tags         user-courses
// @Accept       json
// @Produce      json
// @Param        body  body      map[string]string  true  "course_id and semester_id"
// @Success      200   {object}  map[string]string  "message: Course added successfully"
// @Failure      400   {object}  dto.ErrorResponse  "Invalid input"
// @Failure      404   {object}  dto.ErrorResponse  "Course or semester not found"
// @Failure      409   {object}  dto.ErrorResponse  "Conflict (e.g. already selected)"
// @Failure      500   {object}  dto.ErrorResponse  "Internal server error"
// @Router       /v1/user/courses/select [post]
// @Security     BearerAuth
func (h *UserCourseHandler) AddCourse(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))

	var req struct {
		CourseID   uuid.UUID `json:"course_id" binding:"required"`
		SemesterID uuid.UUID `json:"semester_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid course selection request",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.service.AddCourse(userID, req.CourseID, req.SemesterID); err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, errors.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to add course",
				zap.String("user_id", userID.String()),
				zap.String("course_id", req.CourseID.String()),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add course"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course added successfully"})
}

// RemoveCourse handles course removal
// @Summary      Remove Course
// @Description  Removes a course from the current user's schedule
// @Tags         user-courses
// @Produce      json
// @Param        courseId  path      string              true  "Course ID"
// @Success      200       {object}  map[string]string   "message: Course removed successfully"
// @Failure      400       {object}  dto.ErrorResponse   "Invalid course ID"
// @Failure      404       {object}  dto.ErrorResponse   "Course not found in user's schedule"
// @Failure      500       {object}  dto.ErrorResponse   "Internal server error"
// @Router       /v1/user/courses/select/{courseId} [delete]
// @Security     BearerAuth
func (h *UserCourseHandler) RemoveCourse(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		h.logger.Warn("Invalid course ID format",
			zap.String("course_id", c.Param("courseId")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if err := h.service.RemoveCourse(userID, courseID); err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found in user's schedule"})
		default:
			h.logger.Error("Failed to remove course",
				zap.String("user_id", userID.String()),
				zap.String("course_id", courseID.String()),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove course"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course removed successfully"})
}

// GetUserCourses returns all courses for a user in a specific semester
// @Summary      Get User Courses
// @Description  Retrieves all courses selected by the user for a given semester
// @Tags         user-courses
// @Produce      json
// @Param        semester_id  query     string  true  "Semester ID"
// @Success      200          {array}   dto.CourseResponse
// @Failure      400          {object}  dto.ErrorResponse  "Invalid semester ID"
// @Failure      500          {object}  dto.ErrorResponse  "Internal server error"
// @Router       /v1/user/courses/selected [get]
// @Security     BearerAuth
func (h *UserCourseHandler) GetUserCourses(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	semesterID, err := uuid.Parse(c.Query("semester_id"))
	if err != nil {
		h.logger.Warn("Invalid semester ID format",
			zap.String("semester_id", c.Query("semester_id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid semester ID"})
		return
	}

	courses, err := h.service.GetUserCourses(userID, semesterID)
	if err != nil {
		h.logger.Error("Failed to fetch user courses",
			zap.String("user_id", userID.String()),
			zap.String("semester_id", semesterID.String()),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// ValidateTimeConflicts checks for time conflicts with existing courses
// @Summary      Validate Time Conflicts
// @Description  Checks if adding the selected course causes any time conflicts
// @Tags         user-courses
// @Produce      json
// @Param        course_id    query     string  true  "Course ID"
// @Param        semester_id  query     string  true  "Semester ID"
// @Success      200          {object}  map[string]string  "message: No time conflicts found"
// @Failure      400          {object}  dto.ErrorResponse  "Invalid input"
// @Failure      409          {object}  dto.ErrorResponse  "Time conflict exists"
// @Router       /v1/user/courses/validate [get]
// @Security     BearerAuth
func (h *UserCourseHandler) ValidateTimeConflicts(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	courseID, err := uuid.Parse(c.Query("course_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	semesterID, err := uuid.Parse(c.Query("semester_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid semester ID"})
		return
	}

	if err := h.service.ValidateTimeConflicts(userID, semesterID, courseID); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "No time conflicts found"})
}
