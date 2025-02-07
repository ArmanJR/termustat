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

type SemesterHandler struct {
	service services.SemesterService
	logger  *zap.Logger
}

func NewSemesterHandler(service services.SemesterService, logger *zap.Logger) *SemesterHandler {
	return &SemesterHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SemesterHandler) Create(c *gin.Context) {
	var req dto.CreateSemesterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid semester request",
			zap.Error(err),
			zap.String("handler", "CreateSemester"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	semester, err := h.service.Create(&req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, errors.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to create semester", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, semester)
}

func (h *SemesterHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid semester ID format",
			zap.String("id", c.Param("id")),
			zap.String("handler", "GetSemester"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid semester ID"})
		return
	}

	semester, err := h.service.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Semester not found"})
		default:
			h.logger.Error("Failed to get semester",
				zap.String("id", id.String()),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, semester)
}

func (h *SemesterHandler) GetAll(c *gin.Context) {
	semesters, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Failed to fetch semesters", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, semesters)
}

func (h *SemesterHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid semester ID format",
			zap.String("id", c.Param("id")),
			zap.String("handler", "UpdateSemester"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid semester ID"})
		return
	}

	var req dto.UpdateSemesterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid semester update request",
			zap.Error(err),
			zap.String("handler", "UpdateSemester"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	semester, err := h.service.Update(id, &req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Semester not found"})
		case errors.Is(err, errors.ErrInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, errors.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to update semester",
				zap.String("id", id.String()),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, semester)
}

func (h *SemesterHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid semester ID format",
			zap.String("id", c.Param("id")),
			zap.String("handler", "DeleteSemester"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid semester ID"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Semester not found"})
		default:
			h.logger.Error("Failed to delete semester",
				zap.String("id", id.String()),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Semester deleted successfully"})
}
