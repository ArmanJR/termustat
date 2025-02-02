package handlers

import (
	"github.com/armanjr/termustat/app/config"
	"github.com/armanjr/termustat/app/logger"
	"github.com/armanjr/termustat/app/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

// UserCourseRequest represents user course assignment request
type UserCourseRequest struct {
	CourseID   string `json:"course_id" binding:"required,uuid4"`
	SemesterID string `json:"semester_id" binding:"required,uuid4"`
}

// GetUserCourses returns all courses for a user
func GetUserCourses(c *gin.Context) {
	userID := c.Param("id")
	var courses []models.UserCourse

	if err := config.DB.Where("user_id = ?", userID).
		Preload("Course").
		Preload("Semester").
		Find(&courses).Error; err != nil {
		logger.Log.Error("Failed to fetch user courses", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// AddUserCourse adds a course to user's schedule
func AddUserCourse(c *gin.Context) {
	userID := c.Param("id")
	var req UserCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid user course request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if course exists
	var course models.Course
	if err := config.DB.First(&course, "id = ?", req.CourseID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course not found"})
		return
	}

	// Check if semester exists
	var semester models.Semester
	if err := config.DB.First(&semester, "id = ?", req.SemesterID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semester not found"})
		return
	}

	// Check for existing enrollment
	var existing models.UserCourse
	if err := config.DB.Where(
		"user_id = ? AND course_id = ? AND semester_id = ?",
		userID,
		req.CourseID,
		req.SemesterID,
	).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Course already added for this semester"})
		return
	}

	userCourse := models.UserCourse{
		UserID:     uuid.MustParse(userID),
		CourseID:   uuid.MustParse(req.CourseID),
		SemesterID: uuid.MustParse(req.SemesterID),
	}

	if err := config.DB.Create(&userCourse).Error; err != nil {
		logger.Log.Error("Failed to add course", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add course"})
		return
	}

	c.JSON(http.StatusCreated, userCourse)
}

// RemoveUserCourse removes a course from user's schedule
func RemoveUserCourse(c *gin.Context) {
	userID := c.Param("id")
	courseID := c.Param("course_id")

	var userCourse models.UserCourse
	if err := config.DB.Where(
		"user_id = ? AND course_id = ?",
		userID,
		courseID,
	).First(&userCourse).Error; err != nil {
		logger.Log.Warn("Course not found in user schedule", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found in schedule"})
		return
	}

	if err := config.DB.Delete(&userCourse).Error; err != nil {
		logger.Log.Error("Failed to remove course", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course removed successfully"})
}
