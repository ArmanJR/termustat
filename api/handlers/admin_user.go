package handlers

import (
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type AdminUserHandler struct {
	userService services.UserService
	logger      *zap.Logger
}

func NewAdminUserHandler(userService services.UserService, logger *zap.Logger) *AdminUserHandler {
	return &AdminUserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetAllUsers returns a paginated list of all users
func (h *AdminUserHandler) GetAllUsers(c *gin.Context) {
	pagination := &dto.PaginationQuery{
		Page:  parseInt(c.DefaultQuery("page", "1")),
		Limit: parseInt(c.DefaultQuery("limit", "10")),
	}

	// Validate pagination
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.Limit < 1 || pagination.Limit > 100 {
		pagination.Limit = 10
	}
	pagination.Offset = (pagination.Page - 1) * pagination.Limit

	result, err := h.userService.GetAll(pagination)
	if err != nil {
		h.logger.Error("Failed to fetch users",
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUser returns a single user by ID
func (h *AdminUserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid user ID format",
			zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			h.logger.Error("Failed to get user",
				zap.String("id", id.String()),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles user updates by admin
func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid user ID format",
			zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid user update request",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Convert admin request to service request
	updateReq := &dto.UpdateUserRequest{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		UniversityID: req.UniversityID,
		FacultyID:    req.FacultyID,
		Gender:       req.Gender,
	}

	user, err := h.userService.Update(id, updateReq)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case errors.Is(err, errors.ErrInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, errors.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			h.logger.Error("Failed to update user",
				zap.String("id", id.String()),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		}
		return
	}

	// If password update is requested
	if req.Password != "" {
		if err := h.userService.UpdatePassword(id, &dto.UpdatePasswordRequest{
			NewPassword: req.Password,
		}); err != nil {
			h.logger.Error("Failed to update user password",
				zap.String("id", id.String()),
				zap.Error(err))
			// Return success anyway since user data was updated
			c.JSON(http.StatusOK, gin.H{
				"user":    user,
				"message": "User updated but password update failed",
			})
			return
		}
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles user deletion by admin
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid user ID format",
			zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userService.Delete(id); err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			h.logger.Error("Failed to delete user",
				zap.String("id", id.String()),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Helper function to parse int with default value
func parseInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}
