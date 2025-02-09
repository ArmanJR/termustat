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

// Create handles the creation of a new course
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
func (h *CourseHandler) GetByFaculty(c *gin.Context) {
	facultyID, err := uuid.Parse(c.Param("facultyID"))
	if err != nil {
		h.logger.Warn("Invalid faculty ID format",
			zap.String("faculty_id", c.Param("facultyID")))
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
