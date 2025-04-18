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
	adminUserService services.AdminUserService
	logger           *zap.Logger
}

func NewAdminUserHandler(adminUserService services.AdminUserService, logger *zap.Logger) *AdminUserHandler {
	return &AdminUserHandler{
		adminUserService: adminUserService,
		logger:           logger,
	}
}

// GetAll returns a paginated list of all users
// @Summary      List users
// @Description  Returns a paginated list of all admin users
// @Tags         users
// @Produce      json
// @Param        page   query     int                false  	"Page number"        default(1)
// @Param        limit  query     int                false  	"Items per page"     default(10)
// @Success      200    {object}  dto.AdminUserListResponse  	"paginated list of AdminUserResponse"
// @Failure      500    {object}  dto.ErrorResponse  			"Failed to fetch users"
// @Router       /v1/admin/users [get]
// @Security     BearerAuth
func (h *AdminUserHandler) GetAll(c *gin.Context) {
	pagination := &dto.PaginationQuery{
		Page:  parseInt(c.DefaultQuery("page", "1")),
		Limit: parseInt(c.DefaultQuery("limit", "10")),
	}

	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.Limit < 1 || pagination.Limit > 100 {
		pagination.Limit = 10
	}
	pagination.Offset = (pagination.Page - 1) * pagination.Limit

	result, err := h.adminUserService.GetAll(pagination)
	if err != nil {
		h.logger.Error("Failed to fetch users",
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Get returns a single user by ID
// @Summary      Get user
// @Description  Retrieves a single admin user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      string              true  "User ID"
// @Success      200  {object}  dto.AdminUserResponse
// @Failure      400  {object}  dto.ErrorResponse   "Invalid user ID"
// @Failure      404  {object}  dto.ErrorResponse   "User not found"
// @Failure      500  {object}  dto.ErrorResponse   "Failed to get user"
// @Router       /v1/admin/users/{id} [get]
// @Security     BearerAuth
func (h *AdminUserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid user ID format",
			zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.adminUserService.Get(id)
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

// Update handles user updates by admin
// @Summary      Update user
// @Description  Updates an existing admin user's fields (and optionally password)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string                     true  "User ID"
// @Param        body  body      dto.AdminUpdateUserRequest true  "Updated user data"
// @Success      200   {object}  dto.AdminUserResponse
// @Failure      400   {object}  dto.ErrorResponse          "Invalid request or user ID"
// @Failure      404   {object}  dto.ErrorResponse          "User not found"
// @Failure      409   {object}  dto.ErrorResponse          "Conflict (e.g. duplicate student ID)"
// @Failure      500   {object}  dto.ErrorResponse          "Failed to update user"
// @Router       /v1/admin/users/{id} [put]
// @Security     BearerAuth
func (h *AdminUserHandler) Update(c *gin.Context) {
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

	updateReq := &dto.AdminUpdateUserRequest{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		UniversityID: req.UniversityID,
		FacultyID:    req.FacultyID,
		Gender:       req.Gender,
	}

	user, err := h.adminUserService.Update(id, updateReq)
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

	if req.Password != "" {
		if err := h.adminUserService.UpdatePassword(id, &dto.AdminUpdatePasswordRequest{
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

// Delete handles user deletion by admin
// @Summary      Delete user
// @Description  Deletes an admin user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      string              true  "User ID"
// @Success      200  {object}  map[string]string   "message: User deleted successfully"
// @Failure      400  {object}  dto.ErrorResponse   "Invalid user ID"
// @Failure      404  {object}  dto.ErrorResponse   "User not found"
// @Failure      500  {object}  dto.ErrorResponse   "Failed to delete user"
// @Router       /v1/admin/users/{id} [delete]
// @Security     BearerAuth
func (h *AdminUserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid user ID format",
			zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.adminUserService.Delete(id); err != nil {
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

func parseInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}
