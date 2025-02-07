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

type UniversityHandler struct {
	service services.UniversityService
	logger  *zap.Logger
}

func NewUniversityHandler(service services.UniversityService, logger *zap.Logger) *UniversityHandler {
	return &UniversityHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UniversityHandler) Create(c *gin.Context) {
	var req dto.CreateUniversityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid university request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	university, err := h.service.Create(&req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrValidation):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, errors.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to create university", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, university)
}

func (h *UniversityHandler) Get(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Warn("Invalid university ID format", zap.String("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	university, err := h.service.GetByID(parsedID)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "University not found"})
		default:
			h.logger.Error("Failed to get university", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, university)
}

func (h *UniversityHandler) GetAll(c *gin.Context) {
	universities, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Failed to fetch universities", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, universities)
}

func (h *UniversityHandler) Update(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Warn("Invalid university ID format", zap.String("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	var req dto.UpdateUniversityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid university update request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	university, err := h.service.Update(parsedID, &req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "University not found"})
		case errors.Is(err, errors.ErrValidation):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to update university", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, university)
}

func (h *UniversityHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Warn("Invalid university ID format", zap.String("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	if err := h.service.Delete(parsedID); err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "University not found"})
		default:
			h.logger.Error("Failed to delete university", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "University deleted successfully"})
}
