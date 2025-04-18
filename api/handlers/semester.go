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

// Create a new semester
// @Summary      Create Semester
// @Description  Creates a new semester with year and term
// @Tags         semesters
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateSemesterRequest  true  "Create semester payload"
// @Success      201   {object}  dto.SemesterResponse
// @Failure      400   {object}  dto.ErrorResponse          "Invalid input or duplicate year+term"
// @Failure      409   {object}  dto.ErrorResponse          "Semester already exists"
// @Failure      500   {object}  dto.ErrorResponse          "Internal server error"
// @Router       /v1/admin/semesters [post]
// @Security     BearerAuth
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

// Get retrieves a semester by ID
// @Summary      Get Semester
// @Description  Retrieves a semester by its ID
// @Tags         semesters
// @Produce      json
// @Param        id   path      string              true  "Semester ID"
// @Success      200  {object}  dto.SemesterResponse
// @Failure      400  {object}  dto.ErrorResponse   "Invalid semester ID"
// @Failure      404  {object}  dto.ErrorResponse   "Semester not found"
// @Failure      500  {object}  dto.ErrorResponse   "Internal server error"
// @Router       /v1/admin/semesters/{id} [get]
// @Security     BearerAuth
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

// GetAll lists all semesters
// @Summary      List Semesters
// @Description  Retrieves all semesters, ordered by most recent
// @Tags         semesters
// @Produce      json
// @Success      200  {array}   dto.SemesterResponse
// @Failure      500  {object}  dto.ErrorResponse   "Internal server error"
// @Router       /v1/admin/semesters [get]
// @Security     BearerAuth
func (h *SemesterHandler) GetAll(c *gin.Context) {
	semesters, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Failed to fetch semesters", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, semesters)
}

// Update modifies an existing semester
// @Summary      Update Semester
// @Description  Updates the year or term of an existing semester
// @Tags         semesters
// @Accept       json
// @Produce      json
// @Param        id    path      string                     true  "Semester ID"
// @Param        body  body      dto.UpdateSemesterRequest  true  "Update semester payload"
// @Success      200   {object}  dto.SemesterResponse
// @Failure      400   {object}  dto.ErrorResponse          "Invalid input or ID"
// @Failure      404   {object}  dto.ErrorResponse          "Semester not found"
// @Failure      409   {object}  dto.ErrorResponse          "Duplicate semester"
// @Failure      500   {object}  dto.ErrorResponse          "Internal server error"
// @Router       /v1/admin/semesters/{id} [put]
// @Security     BearerAuth
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

// Delete removes a semester by ID
// @Summary      Delete Semester
// @Description  Deletes a semester by its ID
// @Tags         semesters
// @Produce      json
// @Param        id   path      string              true  "Semester ID"
// @Success      200  {object}  map[string]string   "message: Semester deleted successfully"
// @Failure      400  {object}  dto.ErrorResponse   "Invalid semester ID"
// @Failure      404  {object}  dto.ErrorResponse   "Semester not found"
// @Failure      500  {object}  dto.ErrorResponse   "Internal server error"
// @Router       /v1/admin/semesters/{id} [delete]
// @Security     BearerAuth
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
