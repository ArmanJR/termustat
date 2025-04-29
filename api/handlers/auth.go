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

// Register a new user
// @Summary      Register
// @Description  Creates a new user account and sends a verification email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "Registration payload"
// @Success      201   {object}  map[string]string    "message: Registration successful"
// @Failure      400   {object}  dto.ErrorResponse    "Invalid payload or IDs"
// @Failure      409   {object}  dto.ErrorResponse    "email or student ID already exists"
// @Failure      500   {object}  dto.ErrorResponse    "Failed to register user"
// @Router       /v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid registration request",
			zap.Error(err),
			zap.Any("request", req))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse UniversityID if provided
	var universityID uuid.UUID
	if req.UniversityID != "" {
		uid, err := uuid.Parse(req.UniversityID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university_id"})
			return
		}
		universityID = uid
	}

	// Parse FacultyID if provided
	var facultyID uuid.UUID
	if req.FacultyID != "" {
		fid, err := uuid.Parse(req.FacultyID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty_id"})
			return
		}
		facultyID = fid
	}

	serviceReq := &dto.RegisterServiceRequest{
		Email:        req.Email,
		Password:     req.Password,
		StudentID:    req.StudentID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		UniversityID: universityID,
		FacultyID:    facultyID,
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

// VerifyEmail confirms a user's email
// @Summary      Verify Email
// @Description  Verifies a user's email using the provided token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.VerifyEmailRequest  true  "Verification payload"
// @Success      200   {object}  map[string]string       "message: Email verified successfully"
// @Failure      400   {object}  dto.ErrorResponse       "Invalid or expired token"
// @Failure      500   {object}  dto.ErrorResponse       "Failed to verify email"
// @Router       /v1/auth/verify-email [post]
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

// Login authenticates a user
// @Summary      Login
// @Description  Authenticates user and returns an access token in Authorization header and refresh token as HTTP-only cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest   true  "Login payload"
// @Success      200   {object}  map[string]string  "message: logged in successfully"
// @Header       200   {string}  Authorization      "Bearer <access_token>"
// @Header       200   {string}  Set-Cookie         "refresh_token=<token>; Path=/; HttpOnly; Secure"
// @Failure      400   {object}  dto.ErrorResponse  "Invalid payload"
// @Failure      401   {object}  dto.ErrorResponse  "Invalid credentials"
// @Failure      403   {object}  dto.ErrorResponse  "Email not verified"
// @Failure      500   {object}  dto.ErrorResponse  "Failed to login"
// @Router       /v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := h.authService.Login(req.Email, req.Password)
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

	c.SetCookie(
		"refresh_token",
		refresh,
		7*24*60*60, // 7 days in seconds
		"/",
		"",
		true,
		true,
	)

	c.Header("Authorization", "Bearer "+access)
	c.JSON(http.StatusOK, gin.H{
		"message": "logged in successfully",
	})
}

// Refresh provides a new access-token / refresh-token pair
// @Summary      Refresh token
// @Description  Generate new access token using refresh token from HTTP-only cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200   {object}  map[string]string    "message: token refreshed"
// @Header       200   {string}  Authorization        "Bearer <access_token>"
// @Header       200   {string}  Set-Cookie           "refresh_token=<token>; Path=/; HttpOnly; Secure"
// @Failure      400   {object}  dto.ErrorResponse    "Missing refresh token cookie"
// @Failure      401   {object}  dto.ErrorResponse    "Invalid or expired refresh token"
// @Router       /v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh token"})
		return
	}

	access, newRefresh, err := h.authService.Refresh(refresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"refresh_token",
		newRefresh,
		7*24*60*60,
		"/",
		"",
		true,
		true,
	)

	c.Header("Authorization", "Bearer "+access)
	c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
}

// Logout revokes the current session's refresh token
// @Summary      Logout
// @Description  Revokes the current refresh token and clears the cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200   {object}  map[string]string   "message: logged out successfully"
// @Header       200   {string}  Set-Cookie          "refresh_token=; Path=/; HttpOnly; Secure; MaxAge=0"
// @Failure      400   {object}  dto.ErrorResponse   "Missing refresh token"
// @Router       /v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh token"})
		return
	}

	if err := h.authService.Logout(refresh); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	// Clear the refresh token cookie
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// ForgotPassword starts a password reset
// @Summary      Forgot Password
// @Description  Sends a password reset link if the email exists
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.ForgotPasswordRequest  true  "Forgot password payload"
// @Success      200   {object}  map[string]string          "message: If the email exists, a reset link will be sent"
// @Failure      400   {object}  dto.ErrorResponse          "Invalid payload"
// @Failure      500   {object}  dto.ErrorResponse          "Failed to process request"
// @Router       /v1/auth/forgot-password [post]
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

// ResetPassword completes a password reset
// @Summary      Reset Password
// @Description  Resets the user's password using the provided token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.ResetPasswordRequest  true  "Reset password payload"
// @Success      200   {object}  map[string]string        "message: Password reset successful"
// @Failure      400   {object}  dto.ErrorResponse        "Invalid payload or token"
// @Failure      500   {object}  dto.ErrorResponse        "Failed to reset password"
// @Router       /v1/auth/reset-password [post]
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

// GetCurrentUser returns the authenticated user's info
// @Summary      Get Current User
// @Description  Retrieves information about the authenticated user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dto.AdminUserResponse
// @Failure      401  {object}  dto.ErrorResponse  "User not authenticated"
// @Failure      404  {object}  dto.ErrorResponse  "User not found"
// @Router       /v1/user/me [get]
// @Security     BearerAuth
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
