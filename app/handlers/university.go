package handlers

import (
	"net/http"
	"strings"

	"github.com/armanjr/termustat/app/config"
	"github.com/armanjr/termustat/app/logger"
	"github.com/armanjr/termustat/app/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UniversityRequest represents university create/update request
type UniversityRequest struct {
	NameEn   string `json:"name_en" binding:"required"`
	NameFa   string `json:"name_fa" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

// CreateUniversity creates a new university
func CreateUniversity(c *gin.Context) {
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

	if err := config.DB.Create(&university).Error; err != nil {
		logger.Log.Error("Failed to create university", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create university"})
		return
	}

	c.JSON(http.StatusCreated, university)
}

// GetUniversity returns a single university
func GetUniversity(c *gin.Context) {
	id := c.Param("id")
	var university models.University

	if err := config.DB.Preload("Faculties").First(&university, "id = ?", id).Error; err != nil {
		logger.Log.Warn("University not found", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "University not found"})
		return
	}

	c.JSON(http.StatusOK, university)
}

// GetAllUniversities returns all universities
func GetAllUniversities(c *gin.Context) {
	var universities []models.University
	if err := config.DB.Find(&universities).Error; err != nil {
		logger.Log.Error("Failed to fetch universities", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch universities"})
		return
	}
	c.JSON(http.StatusOK, universities)
}

// UpdateUniversity updates university details
func UpdateUniversity(c *gin.Context) {
	id := c.Param("id")
	var university models.University

	if err := config.DB.First(&university, "id = ?", id).Error; err != nil {
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

	if err := config.DB.Save(&university).Error; err != nil {
		logger.Log.Error("Failed to update university", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update university"})
		return
	}

	c.JSON(http.StatusOK, university)
}

// DeleteUniversity soft deletes a university
func DeleteUniversity(c *gin.Context) {
	id := c.Param("id")
	var university models.University

	if err := config.DB.First(&university, "id = ?", id).Error; err != nil {
		logger.Log.Warn("University not found for deletion", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "University not found"})
		return
	}

	if err := config.DB.Delete(&university).Error; err != nil {
		logger.Log.Error("Failed to delete university", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete university"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "University deleted successfully"})
}
