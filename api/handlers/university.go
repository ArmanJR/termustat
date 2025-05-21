package handlers

import (
	"context"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
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

// Create a new university
// @Summary      Create University
// @Description  Creates a new university
// @Tags         universities
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateUniversityRequest  true  "Create university payload"
// @Success      201   {object}  dto.UniversityResponse
// @Failure      400   {object}  dto.ErrorResponse           "Invalid input"
// @Failure      409   {object}  dto.ErrorResponse           "Conflict (name already exists)"
// @Failure      500   {object}  dto.ErrorResponse           "Internal server error"
// @Router       /v1/admin/universities [post]
// @Security     BearerAuth
func (h *UniversityHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var req dto.CreateUniversityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid university request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	exists, err := h.service.ExistsByName(ctx, req.NameEn, req.NameFa)
	if err != nil {
		h.logger.Error("Failed to check university existence", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "University with the same name already exists"})
		return
	}

	university, err := h.service.Create(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrInvalid):
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

// Get retrieves a university by ID
// @Summary      Get University
// @Description  Retrieves a university by its ID
// @Tags         universities
// @Produce      json
// @Param        id   path      string              true  "University ID"
// @Success      200  {object}  dto.UniversityResponse
// @Failure      400  {object}  dto.ErrorResponse   "Invalid university ID"
// @Failure      404  {object}  dto.ErrorResponse   "University not found"
// @Failure      500  {object}  dto.ErrorResponse   "Internal server error"
// @Router       /v1/admin/universities/{id} [get]
// @Security     BearerAuth
func (h *UniversityHandler) Get(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Warn("Invalid university ID format", zap.String("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	university, err := h.service.Get(ctx, parsedID)
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

// GetAll lists all universities
// @Summary      List Universities
// @Description  Retrieves all universities
// @Tags         universities
// @Produce      json
// @Success      200  {array}   dto.UniversityResponse
// @Failure      500  {object}  dto.ErrorResponse   "Internal server error"
// @Router       /v1/admin/universities [get]
// @Security     BearerAuth
func (h *UniversityHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	universities, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error("Failed to fetch universities", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, universities)
}

// Update modifies an existing university
// @Summary      Update University
// @Description  Updates the specified university’s details
// @Tags         universities
// @Accept       json
// @Produce      json
// @Param        id    path      string                     true  "University ID"
// @Param        body  body      dto.UpdateUniversityRequest true  "Update university payload"
// @Success      200   {object}  dto.UniversityResponse
// @Failure      400   {object}  dto.ErrorResponse          "Invalid input or ID"
// @Failure      404   {object}  dto.ErrorResponse          "University not found"
// @Failure      500   {object}  dto.ErrorResponse          "Internal server error"
// @Router       /v1/admin/universities/{id} [put]
// @Security     BearerAuth
func (h *UniversityHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

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

	university, err := h.service.Update(ctx, parsedID, &req)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "University not found"})
		case errors.Is(err, errors.ErrInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to update university", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, university)
}

// Delete removes a university by ID
// @Summary      Delete University
// @Description  Deletes the specified university
// @Tags         universities
// @Produce      json
// @Param        id   path      string              true  "University ID"
// @Success      200  {object}  map[string]string   "message: University deleted successfully"
// @Failure      400  {object}  dto.ErrorResponse   "Invalid university ID"
// @Failure      404  {object}  dto.ErrorResponse   "University not found"
// @Failure      500  {object}  dto.ErrorResponse   "Internal server error"
// @Router       /v1/admin/universities/{id} [delete]
// @Security     BearerAuth
func (h *UniversityHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Warn("Invalid university ID format", zap.String("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	if err := h.service.Delete(ctx, parsedID); err != nil {
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
