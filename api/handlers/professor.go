package handlers

import (
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/services"
	"net/http"
	"time"

	"github.com/armanjr/termustat/api/config"
	"github.com/armanjr/termustat/api/logger"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProfessorHandler struct {
	professorService services.ProfessorService
	logger           *zap.Logger
}

func NewProfessorHandler(professorService services.ProfessorService, logger *zap.Logger) *ProfessorHandler {
	return &ProfessorHandler{
		professorService: professorService,
		logger:           logger,
	}
}

// GetByUniversity returns all professors for a university
func (h *ProfessorHandler) GetByUniversity(c *gin.Context) {
	universityID := c.Param("id")

	parsedUniversityID, err := uuid.Parse(universityID)
	if err != nil {
		h.logger.Warn("Invalid university ID format", zap.String("university_id", universityID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	professors, err := h.professorService.GetProfessorsByUniversity(parsedUniversityID)
	if err != nil {
		h.logger.Error("Get professors error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, *professors)
}

// Get returns a single professor
func (h *ProfessorHandler) Get(c *gin.Context) {
	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Warn("Invalid professor ID format", zap.String("professor_id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid professor ID"})
		return
	}

	professor, err := h.professorService.GetProfessor(parsedID)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Professor not found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, professor)
}

// UpdateProfessorRequest represents professor update request
type UpdateProfessorRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateProfessor updates a professor's name
func UpdateProfessor(c *gin.Context) {
	id := c.Param("id")
	var professor models.Professor

	if err := config.DB.First(&professor, "id = ?", id).Error; err != nil {
		logger.Log.Warn("Professor not found for update", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Professor not found"})
		return
	}

	var req UpdateProfessorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid professor update request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	normalizedName := utils.NormalizeProfessor(req.Name)
	if normalizedName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid professor name after normalization"})
		return
	}

	// Check for existing professor with same normalized name
	var existing models.Professor
	err := config.DB.Where(
		"university_id = ? AND normalized_name = ? AND id != ?",
		professor.UniversityID,
		normalizedName,
		professor.ID,
	).First(&existing).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Professor with this name already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Error("Database error checking professor existence", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update professor"})
		return
	}

	professor.Name = req.Name
	professor.NormalizedName = normalizedName
	professor.UpdatedAt = time.Now()

	if err := config.DB.Save(&professor).Error; err != nil {
		logger.Log.Error("Failed to update professor", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update professor"})
		return
	}

	c.JSON(http.StatusOK, professor)
}
