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

type FacultyHandler struct {
	facultyService services.FacultyService
	logger         *zap.Logger
}

func NewFacultyHandler(facultyService services.FacultyService, logger *zap.Logger) *FacultyHandler {
	return &FacultyHandler{
		facultyService: facultyService,
		logger:         logger,
	}
}

func (h *FacultyHandler) Create(c *gin.Context) {
	var req dto.CreateFacultyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid create faculty request",
			zap.String("handler", "Faculty"),
			zap.String("operation", "Create"),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	faculty, err := h.facultyService.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, errors.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to create faculty",
				zap.String("handler", "Faculty"),
				zap.String("operation", "Create"),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create faculty"})
		}
		return
	}

	c.JSON(http.StatusCreated, faculty)
}

func (h *FacultyHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty ID"})
		return
	}

	faculty, err := h.facultyService.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to get faculty",
				zap.String("id", id.String()),
				zap.String("handler", "Faculty"),
				zap.String("operation", "GetByID"),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get faculty"})
		}
		return
	}

	c.JSON(http.StatusOK, faculty)
}

func (h *FacultyHandler) GetAllByUniversity(c *gin.Context) {
	universityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	faculties, err := h.facultyService.GetAllByUniversity(universityID)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to get faculties",
				zap.String("id", universityID.String()),
				zap.String("handler", "Faculty"),
				zap.String("operation", "GetByUniversity"),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get faculties"})
		}
		return
	}

	c.JSON(http.StatusOK, faculties)
}

func (h *FacultyHandler) GetByUniversityAndShortCode(c *gin.Context) {
	universityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	shortCode := c.Param("short_code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
		return
	}

	faculty, err := h.facultyService.GetByUniversityAndShortCode(universityID, shortCode)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to get faculty by short code",
				zap.String("id", universityID.String()),
				zap.String("short_code", shortCode),
				zap.String("handler", "Faculty"),
				zap.String("operation", "GetByUniversityAndShortCode"),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get faculty"})
		}
		return
	}

	c.JSON(http.StatusOK, faculty)
}

func (h *FacultyHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty ID"})
		return
	}

	var req dto.UpdateFacultyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid update faculty request",
			zap.String("handler", "Faculty"),
			zap.String("operation", "Update"),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	faculty, err := h.facultyService.Update(id, req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, errors.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to update faculty",
				zap.String("id", id.String()),
				zap.String("handler", "Faculty"),
				zap.String("operation", "Update"),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update faculty"})
		}
		return
	}

	c.JSON(http.StatusOK, faculty)
}

func (h *FacultyHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty ID"})
		return
	}

	if err := h.facultyService.Delete(id); err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to delete faculty",
				zap.String("id", id.String()),
				zap.String("handler", "Faculty"),
				zap.String("operation", "Delete"),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete faculty"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Faculty deleted successfully"})
}
