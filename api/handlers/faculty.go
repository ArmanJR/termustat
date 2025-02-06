package handlers

import (
	"net/http"
	"strings"

	"github.com/armanjr/termustat/api/config"
	"github.com/armanjr/termustat/api/logger"
	"github.com/armanjr/termustat/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// FacultyRequest represents faculty create/update request
type FacultyRequest struct {
	UniversityID string `json:"university_id" binding:"required,uuid4"`
	NameEn       string `json:"name_en" binding:"required"`
	NameFa       string `json:"name_fa" binding:"required"`
	ShortCode    string `json:"short_code" binding:"required,max=10"`
	IsActive     *bool  `json:"is_active" binding:"required"`
}

// CreateFaculty creates a new faculty
func CreateFaculty(c *gin.Context) {
	var req FacultyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid faculty request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	universityID, err := uuid.Parse(req.UniversityID)
	if err != nil {
		logger.Log.Warn("Invalid university ID format", zap.String("university_id", req.UniversityID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	faculty := models.Faculty{
		UniversityID: universityID,
		NameEn:       strings.TrimSpace(req.NameEn),
		NameFa:       strings.TrimSpace(req.NameFa),
		ShortCode:    strings.ToUpper(strings.TrimSpace(req.ShortCode)),
		IsActive:     *req.IsActive,
	}

	if err := config.DB.Create(&faculty).Error; err != nil {
		logger.Log.Error("Failed to create faculty", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create faculty"})
		return
	}

	c.JSON(http.StatusCreated, faculty)
}

// GetFaculty returns a single faculty
func GetFaculty(c *gin.Context) {
	id := c.Param("id")
	var faculty models.Faculty

	if err := config.DB.First(&faculty, "id = ?", id).Error; err != nil {
		logger.Log.Warn("Faculty not found", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Faculty not found"})
		return
	}

	c.JSON(http.StatusOK, faculty)
}

// GetAllFaculties returns all faculties for a university
func GetAllFaculties(c *gin.Context) {
	universityID := c.Param("university_id")
	var faculties []models.Faculty

	_, err := uuid.Parse(universityID)
	if err != nil {
		logger.Log.Warn("Invalid university ID format", zap.String("university_id", universityID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	if err := config.DB.Where("university_id = ?", universityID).Find(&faculties).Error; err != nil {
		logger.Log.Error("Failed to fetch faculties", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch faculties"})
		return
	}

	c.JSON(http.StatusOK, faculties)
}

// UpdateFaculty updates faculty details
func UpdateFaculty(c *gin.Context) {
	id := c.Param("id")
	var faculty models.Faculty

	if err := config.DB.First(&faculty, "id = ?", id).Error; err != nil {
		logger.Log.Warn("Faculty not found for update", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Faculty not found"})
		return
	}

	var req FacultyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid faculty update request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	universityID, err := uuid.Parse(req.UniversityID)
	if err != nil {
		logger.Log.Warn("Invalid university ID format", zap.String("university_id", req.UniversityID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	faculty.UniversityID = universityID
	faculty.NameEn = strings.TrimSpace(req.NameEn)
	faculty.NameFa = strings.TrimSpace(req.NameFa)
	faculty.ShortCode = strings.ToUpper(strings.TrimSpace(req.ShortCode))
	faculty.IsActive = *req.IsActive

	if err := config.DB.Save(&faculty).Error; err != nil {
		logger.Log.Error("Failed to update faculty", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update faculty"})
		return
	}

	c.JSON(http.StatusOK, faculty)
}

// DeleteFaculty deletes a faculty
func DeleteFaculty(c *gin.Context) {
	id := c.Param("id")
	var faculty models.Faculty

	if err := config.DB.First(&faculty, "id = ?", id).Error; err != nil {
		logger.Log.Warn("Faculty not found for deletion", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Faculty not found"})
		return
	}

	if err := config.DB.Delete(&faculty).Error; err != nil {
		logger.Log.Error("Failed to delete faculty", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete faculty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Faculty deleted successfully"})
}
