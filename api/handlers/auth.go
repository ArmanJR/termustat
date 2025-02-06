package handlers

import (
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	authService services.AuthService
	logger      *zap.Logger
}

func NewAuthHandler(authService services.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid registration request",
			zap.Error(err),
			zap.Any("request", req))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serviceReq := &dto.RegisterServiceRequest{
		Email:        req.Email,
		Password:     req.Password,
		StudentID:    req.StudentID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		UniversityID: uuid.MustParse(req.UniversityID),
		FacultyID:    uuid.MustParse(req.FacultyID),
		Gender:       req.Gender,
	}

	if err := h.authService.Register(serviceReq); err != nil {
		status := http.StatusInternalServerError
		message := "Failed to register user"

		if err.Error() == "email or student ID already exists" {
			status = http.StatusConflict
			message = err.Error()
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	h.logger.Info("New user registered",
		zap.String("email", req.Email),
		zap.String("student_id", req.StudentID))

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Please check your email to verify your account.",
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.VerifyEmail(req.Token); err != nil {
		status := http.StatusInternalServerError
		message := "Failed to verify email"

		if err.Error() == "invalid or expired token" {
			status = http.StatusBadRequest
			message = err.Error()
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to login"

		switch err.Error() {
		case "invalid credentials":
			status = http.StatusUnauthorized
		case "email not verified":
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ForgotPassword(req.Email); err != nil {
		h.logger.Error("Failed to process forgot password request",
			zap.String("email", req.Email),
			zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a reset link will be sent",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ResetPassword(req.Token, req.Password); err != nil {
		status := http.StatusInternalServerError
		message := "Failed to reset password"

		if err.Error() == "invalid or expired token" {
			status = http.StatusBadRequest
			message = err.Error()
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	parsedID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.authService.GetCurrentUser(parsedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"student_id": user.StudentID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"university": user.UniversityID,
		"faculty":    user.FacultyID,
		"verified":   user.EmailVerified,
	})
}

//func (h *AuthHandler) UpdateUser(c *gin.Context) {
//	userID, exists := c.Get("userID")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
//		return
//	}
//
//	var req dto.UpdateUserRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	parsedID, err := uuid.Parse(userID.(string))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
//		return
//	}
//
//	if err := h.authService.UpdateUser(parsedID, &req); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
//}
//
//func (h *AuthHandler) ChangePassword(c *gin.Context) {
//	userID, exists := c.Get("userID")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
//		return
//	}
//
//	var req dto.ChangePasswordRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	parsedID, err := uuid.Parse(userID.(string))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
//		return
//	}
//
//	if err := h.authService.ChangePassword(parsedID, req.OldPassword, req.NewPassword); err != nil {
//		status := http.StatusInternalServerError
//		message := "Failed to change password"
//
//		if err.Error() == "invalid old password" {
//			status = http.StatusBadRequest
//			message = err.Error()
//		}
//
//		c.JSON(status, gin.H{"error": message})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
//}
