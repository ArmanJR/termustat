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

// GetAllByUniversity returns all professors for a university
func (h *ProfessorHandler) GetAllByUniversity(c *gin.Context) {
	universityID := c.Param("id")

	parsedUniversityID, err := uuid.Parse(universityID)
	if err != nil {
		h.logger.Warn("Invalid university ID format", zap.String("university_id", universityID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	professors, err := h.professorService.GetAllByUniversity(parsedUniversityID)
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

	professor, err := h.professorService.Get(parsedID)
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

func (h *ProfessorHandler) Create(c *gin.Context) {
	var req dto.CreateProfessorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request format", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	professor, err := h.professorService.GetOrCreateByName(req.UniversityID, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "University not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, professor)
}
