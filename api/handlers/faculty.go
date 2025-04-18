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

// Create a new faculty
// @Summary      Create Faculty
// @Description  Creates a new faculty under a university
// @Tags         faculties
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateFacultyDTO  true  "Create faculty payload"
// @Success      201   {object}  dto.FacultyResponse
// @Failure      400   {object}  dto.ErrorResponse     "Invalid request body"
// @Failure      404   {object}  dto.ErrorResponse     "University not found"
// @Failure      409   {object}  dto.ErrorResponse     "Conflict (e.g. short code exists)"
// @Failure      500   {object}  dto.ErrorResponse     "Failed to create faculty"
// @Router       /v1/admin/faculties [post]
// @Security     BearerAuth
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

// GetByID retrieves a faculty by its ID
// @Summary      Get Faculty
// @Description  Retrieves a faculty by its unique ID
// @Tags         faculties
// @Produce      json
// @Param        id   path      string              true  "Faculty ID"
// @Success      200  {object}  dto.FacultyResponse
// @Failure      400  {object}  dto.ErrorResponse   "Invalid faculty ID"
// @Failure      404  {object}  dto.ErrorResponse   "Faculty not found"
// @Failure      500  {object}  dto.ErrorResponse   "Failed to get faculty"
// @Router       /v1/admin/faculties/{id} [get]
// @Security     BearerAuth
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

// GetAllByUniversity lists all faculties for a university
// @Summary      List Faculties
// @Description  Retrieves all faculties under the specified university
// @Tags         faculties
// @Produce      json
// @Param        id   path      string                true  "University ID"
// @Success      200  {array}   dto.FacultyResponse
// @Failure      400  {object}  dto.ErrorResponse     "Invalid university ID"
// @Failure      404  {object}  dto.ErrorResponse     "University not found"
// @Failure      500  {object}  dto.ErrorResponse     "Failed to get faculties"
// @Router       /v1/admin/universities/{id}/faculties [get]
// @Security     BearerAuth
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

// GetByUniversityAndShortCode retrieves a faculty by its short code
// @Summary      Get Faculty by Short Code
// @Description  Retrieves a faculty within a university by its short code
// @Tags         faculties
// @Produce      json
// @Param        id          path      string              true  "University ID"
// @Param        short_code  path      string              true  "Faculty short code"
// @Success      200         {object}  dto.FacultyResponse
// @Failure      400         {object}  dto.ErrorResponse   "Invalid university ID or short code missing"
// @Failure      404         {object}  dto.ErrorResponse   "Faculty not found"
// @Failure      500         {object}  dto.ErrorResponse   "Failed to get faculty"
// @Router       /v1/admin/universities/{id}/faculty/{short_code} [get]
// @Security     BearerAuth
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

// Update modifies an existing faculty
// @Summary      Update Faculty
// @Description  Updates the details of an existing faculty
// @Tags         faculties
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "Faculty ID"
// @Param        body  body      dto.UpdateFacultyDTO true  "Update faculty payload"
// @Success      200   {object}  dto.FacultyResponse
// @Failure      400   {object}  dto.ErrorResponse    "Invalid request body or ID"
// @Failure      404   {object}  dto.ErrorResponse    "Faculty not found"
// @Failure      409   {object}  dto.ErrorResponse    "Conflict updating faculty"
// @Failure      500   {object}  dto.ErrorResponse    "Failed to update faculty"
// @Router       /v1/admin/faculties/{id} [put]
// @Security     BearerAuth
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

// Delete removes a faculty by ID
// @Summary      Delete Faculty
// @Description  Deletes a faculty from the system
// @Tags         faculties
// @Produce      json
// @Param        id   path      string              true  "Faculty ID"
// @Success      200  {object}  map[string]string   "message: Faculty deleted successfully"
// @Failure      400  {object}  dto.ErrorResponse   "Invalid faculty ID"
// @Failure      404  {object}  dto.ErrorResponse   "Faculty not found"
// @Failure      500  {object}  dto.ErrorResponse   "Failed to delete faculty"
// @Router       /v1/admin/faculties/{id} [delete]
// @Security     BearerAuth
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
