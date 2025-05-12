package handlers

import (
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type CourseHandler struct {
	service services.CourseService
	logger  *zap.Logger
}

func NewCourseHandler(service services.CourseService, logger *zap.Logger) *CourseHandler {
	return &CourseHandler{
		service: service,
		logger:  logger,
	}
}

// Create a new course
// @Summary      Create a course
// @Description  Creates a new course in the system
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        course  body      dto.CreateCourseDTO  true   "Course payload"
// @Success      201     {object}  dto.CourseResponse
// @Failure      400     {object}  dto.ErrorResponse     "Invalid request or not found"
// @Failure      409     {object}  dto.ErrorResponse     "Conflict (e.g. duplicate code)"
// @Failure      500     {object}  dto.ErrorResponse     "Internal server error"
// @Router       /courses [post]
func (h *CourseHandler) Create(c *gin.Context) {
	var req dto.CreateCourseDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid course request",
			zap.Error(err),
			zap.String("handler", "CreateCourse"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	course, err := h.service.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, errors.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to create course", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, course)
}

// Get retrieves a course by ID
// @Summary      Get a course
// @Description  Retrieves a course by its ID
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        id   path      string            true  "Course ID"
// @Success      200  {object}  dto.CourseResponse
// @Failure      400  {object}  dto.ErrorResponse  "Invalid ID format"
// @Failure      404  {object}  dto.ErrorResponse  "Course not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /courses/{id} [get]
func (h *CourseHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid course ID format",
			zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := h.service.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		default:
			h.logger.Error("Failed to get course",
				zap.String("id", id.String()),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, course)
}

// GetByFaculty retrieves all courses for a faculty
// @Summary      List courses by faculty
// @Description  Retrieves all courses under the specified faculty
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        id         path      string             true  "Faculty ID"
// @Success      200        {array}   dto.CourseResponse
// @Failure      400        {object}  dto.ErrorResponse  "Invalid faculty ID"
// @Failure      500        {object}  dto.ErrorResponse  "Internal server error"
// @Router       /faculties/{id}/courses [get]
func (h *CourseHandler) GetByFaculty(c *gin.Context) {
	facultyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid faculty ID format",
			zap.String("faculty_id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty ID"})
		return
	}

	courses, err := h.service.GetAllByFaculty(facultyID)
	if err != nil {
		h.logger.Error("Failed to fetch courses by faculty",
			zap.String("faculty_id", facultyID.String()),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// Update handles course updates
// @Summary      Update a course
// @Description  Updates the course identified by its ID
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        id      path      string              true  "Course ID"
// @Param        course  body      dto.UpdateCourseDTO true  "Updated course payload"
// @Success      200     {object}  dto.CourseResponse
// @Failure      400     {object}  dto.ErrorResponse     "Invalid request or not found"
// @Failure      404     {object}  dto.ErrorResponse     "Course not found"
// @Failure      409     {object}  dto.ErrorResponse     "Conflict"
// @Failure      500     {object}  dto.ErrorResponse     "Internal server error"
// @Router       /courses/{id} [put]
func (h *CourseHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid course ID format",
			zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var req dto.UpdateCourseDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid course update request",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	course, err := h.service.Update(id, req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		case errors.Is(err, errors.ErrInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, errors.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to update course",
				zap.String("id", id.String()),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, course)
}

// Delete handles course deletion
// @Summary      Delete a course
// @Description  Deletes the course identified by its ID
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        id   path      string            true  "Course ID"
// @Success      200  {object}  map[string]string  "message: Course deleted successfully"
// @Failure      400  {object}  dto.ErrorResponse  "Invalid ID format"
// @Failure      404  {object}  dto.ErrorResponse  "Course not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /courses/{id} [delete]
func (h *CourseHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid course ID format",
			zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		default:
			h.logger.Error("Failed to delete course",
				zap.String("id", id.String()),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

// Search handles course search with filters
// @Summary      Search courses
// @Description  Searches for courses by faculty, professor, or keyword
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        faculty_id    query     string  false  "Filter by Faculty ID"
// @Param        professor_id  query     string  false  "Filter by Professor ID"
// @Param        q             query     string  false  "Full‐text search query"
// @Success      200           {array}   dto.CourseResponse
// @Failure      400           {object}  dto.ErrorResponse  "Invalid query parameters"
// @Failure      500           {object}  dto.ErrorResponse  "Internal server error"
// @Router       /courses [get]
func (h *CourseHandler) Search(c *gin.Context) {
	var filters dto.CourseSearchFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		h.logger.Warn("Invalid search filters",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search parameters"})
		return
	}

	courses, err := h.service.Search(&filters)
	if err != nil {
		h.logger.Error("Failed to search courses", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, courses)
}
