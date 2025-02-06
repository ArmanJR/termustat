package handlers

import (
	"github.com/armanjr/termustat/api/repositories"
	"net/http"
	"strings"

	"github.com/armanjr/termustat/api/logger"
	"github.com/armanjr/termustat/api/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UniversityHandler struct {
	repo repositories.UniversityRepository
}

func NewUniversityHandler(repo repositories.UniversityRepository) *UniversityHandler {
	return &UniversityHandler{repo: repo}
}

type UniversityRequest struct {
	NameEn   string `json:"name_en" binding:"required"`
	NameFa   string `json:"name_fa" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

func (h *UniversityHandler) Create(c *gin.Context) {
	var req UniversityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid university request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	university := models.University{
		NameEn:   strings.TrimSpace(req.NameEn),
		NameFa:   strings.TrimSpace(req.NameFa),
		IsActive: *req.IsActive,
	}

	if err := h.repo.Create(&university); err != nil {
		logger.Log.Error("Failed to create university", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create university"})
		return
	}

	c.JSON(http.StatusCreated, university)
}

func (h *UniversityHandler) Get(c *gin.Context) {
	id := c.Param("id")
	university, err := h.repo.GetByID(id)
	if err != nil {
		logger.Log.Warn("University not found", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "University not found"})
		return
	}

	c.JSON(http.StatusOK, university)
}

func (h *UniversityHandler) GetAll(c *gin.Context) {
	universities, err := h.repo.GetAll()
	if err != nil {
		logger.Log.Error("Failed to fetch universities", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch universities"})
		return
	}
	c.JSON(http.StatusOK, universities)
}

func (h *UniversityHandler) Update(c *gin.Context) {
	id := c.Param("id")
	university, err := h.repo.GetByID(id)
	if err != nil {
		logger.Log.Warn("University not found for update", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "University not found"})
		return
	}

	var req UniversityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid university update request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	university.NameEn = strings.TrimSpace(req.NameEn)
	university.NameFa = strings.TrimSpace(req.NameFa)
	university.IsActive = *req.IsActive

	if err := h.repo.Update(university); err != nil {
		logger.Log.Error("Failed to update university", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update university"})
		return
	}

	c.JSON(http.StatusOK, university)
}

func (h *UniversityHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(id); err != nil {
		logger.Log.Error("Failed to delete university", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete university"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "University deleted successfully"})
}
