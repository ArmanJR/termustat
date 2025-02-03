package handlers

import (
	"github.com/armanjr/termustat/app/config"
	"github.com/armanjr/termustat/app/logger"
	"github.com/armanjr/termustat/app/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strconv"

	"net/http"
	"time"
)

// UserResponse represents the user data sent in responses
type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	StudentID     string    `json:"student_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	UniversityID  uuid.UUID `json:"university_id"`
	FacultyID     uuid.UUID `json:"faculty_id"`
	Gender        string    `json:"gender"`
	EmailVerified bool      `json:"email_verified"`
	IsAdmin       bool      `json:"is_admin"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Pagination struct to hold pagination parameters
type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

// PaginatedResponse struct to wrap the response with pagination info
type PaginatedResponse struct {
	Data       []UserResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

// AdminUserRequest represents admin user update request
type AdminUserRequest struct {
	Email        string `json:"email"`
	StudentID    string `json:"student_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	UniversityID string `json:"university_id" binding:"uuid4"`
	FacultyID    string `json:"faculty_id" binding:"uuid4"`
	Gender       string `json:"gender" binding:"oneof=male female"`
	Password     string `json:"password"`
}

// GetUser returns a single user (Admin only)
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		logger.Log.Warn("User not found", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		StudentID:     user.StudentID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		UniversityID:  user.UniversityID,
		FacultyID:     user.FacultyID,
		Gender:        user.Gender,
		EmailVerified: user.EmailVerified,
		IsAdmin:       user.IsAdmin,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	})
}

// GetAllUsers returns all users with pagination (Admin only)
func GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var users []models.User
	var total int64

	if err := config.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		logger.Log.Error("Failed to count users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
		return
	}

	if err := config.DB.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		logger.Log.Error("Failed to fetch users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			StudentID:     user.StudentID,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			UniversityID:  user.UniversityID,
			FacultyID:     user.FacultyID,
			Gender:        user.Gender,
			EmailVerified: user.EmailVerified,
			IsAdmin:       user.IsAdmin,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		}
	}

	paginatedResponse := PaginatedResponse{
		Data: response,
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}

	c.JSON(http.StatusOK, paginatedResponse)
}

// UpdateUser updates user details (Admin only)
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		logger.Log.Warn("User not found for update", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req AdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid user update request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Update other fields
	user.Email = req.Email
	user.StudentID = req.StudentID
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.UniversityID = uuid.MustParse(req.UniversityID)
	user.FacultyID = uuid.MustParse(req.FacultyID)
	user.Gender = req.Gender
	user.UpdatedAt = time.Now()

	if err := config.DB.Save(&user).Error; err != nil {
		logger.Log.Error("Failed to update user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		StudentID:    user.StudentID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UniversityID: user.UniversityID,
		FacultyID:    user.FacultyID,
		Gender:       user.Gender,
		IsAdmin:      user.IsAdmin,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	})
}

// DeleteUser deletes a user (Admin only)
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		logger.Log.Warn("User not found for deletion", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		logger.Log.Error("Failed to delete user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
