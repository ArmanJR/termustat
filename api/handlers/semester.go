package handlers

import (
	"github.com/armanjr/termustat/api/config"
	"github.com/armanjr/termustat/api/logger"
	"github.com/armanjr/termustat/api/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// SemesterRequest represents semester create/update request
type SemesterRequest struct {
	Year int    `json:"year" binding:"required,min=1900,max=2200"`
	Term string `json:"term" binding:"required,oneof=spring fall"`
}

// CreateSemester creates a new semester
func CreateSemester(c *gin.Context) {
	var req SemesterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid semester request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	semester := models.Semester{
		Year: req.Year,
		Term: req.Term,
	}

	if err := config.DB.Create(&semester).Error; err != nil {
		logger.Log.Error("Failed to create semester", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create semester"})
		return
	}

	c.JSON(http.StatusCreated, semester)
}

// GetSemester returns a single semester
func GetSemester(c *gin.Context) {
	id := c.Param("id")
	var semester models.Semester

	if err := config.DB.First(&semester, "id = ?", id).Error; err != nil {
		logger.Log.Warn("Semester not found", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Semester not found"})
		return
	}

	c.JSON(http.StatusOK, semester)
}

// GetAllSemesters returns all semesters
func GetAllSemesters(c *gin.Context) {
	var semesters []models.Semester
	if err := config.DB.Order("year DESC, term DESC").Find(&semesters).Error; err != nil {
		logger.Log.Error("Failed to fetch semesters", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch semesters"})
		return
	}
	c.JSON(http.StatusOK, semesters)
}

// UpdateSemester updates semester details
func UpdateSemester(c *gin.Context) {
	id := c.Param("id")
	var semester models.Semester

	if err := config.DB.First(&semester, "id = ?", id).Error; err != nil {
		logger.Log.Warn("Semester not found for update", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Semester not found"})
		return
	}

	var req SemesterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid semester update request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	semester.Year = req.Year
	semester.Term = req.Term

	if err := config.DB.Save(&semester).Error; err != nil {
		logger.Log.Error("Failed to update semester", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update semester"})
		return
	}

	c.JSON(http.StatusOK, semester)
}

// DeleteSemester deletes a semester
func DeleteSemester(c *gin.Context) {
	id := c.Param("id")
	var semester models.Semester

	if err := config.DB.First(&semester, "id = ?", id).Error; err != nil {
		logger.Log.Warn("Semester not found for deletion", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Semester not found"})
		return
	}

	if err := config.DB.Delete(&semester).Error; err != nil {
		logger.Log.Error("Failed to delete semester", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete semester"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Semester deleted successfully"})
}
