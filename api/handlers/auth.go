package handlers

import (
	"context"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AuthHandler struct {
	authService       services.AuthService
	universityService services.UniversityService
	facultyService    services.FacultyService
	logger            *zap.Logger
}

func NewAuthHandler(
	authService services.AuthService,
	universityService services.UniversityService,
	facultyService services.FacultyService,
	logger *zap.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		universityService: universityService,
		facultyService:    facultyService,
		logger:            logger,
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

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

	if err := h.authService.Register(ctx, serviceReq); err != nil {
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.VerifyEmail(ctx, req.Token); err != nil {
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
// @Description  Authenticates user and returns access token in response body and refresh token as HTTP-only cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest   true  "Login payload"
// @Success      200   {object}  dto.LoginResponse  "Contains access_token and expires_in"
// @Header       200   {string}  Set-Cookie         "refresh_token=<token>; Path=/; HttpOnly; Secure"
// @Failure      400   {object}  dto.ErrorResponse  "Invalid payload"
// @Failure      401   {object}  dto.ErrorResponse  "Invalid credentials"
// @Failure      403   {object}  dto.ErrorResponse  "Email not verified"
// @Failure      500   {object}  dto.ErrorResponse  "Failed to login"
// @Router       /v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, accessExpiry, refresh, refreshExpiry, err := h.authService.Login(ctx, req.Email, req.Password)
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
		refreshExpiry,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: access,
		ExpiresIn:   accessExpiry,
	})
}

// Refresh provides a new access-token / refresh-token pair
// @Summary      Refresh token
// @Description  Generate new access token using refresh token from HTTP-only cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200   {object}  dto.LoginResponse    "Contains access_token and expires_in"
// @Header       200   {string}  Set-Cookie           "refresh_token=<token>; Path=/; HttpOnly; Secure"
// @Failure      400   {object}  dto.ErrorResponse    "Missing refresh token cookie"
// @Failure      401   {object}  dto.ErrorResponse    "Invalid or expired refresh token"
// @Router       /v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh token"})
		return
	}

	access, accessExpiry, newRefresh, newRefreshExpiry, err := h.authService.Refresh(ctx, refresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"refresh_token",
		newRefresh,
		newRefreshExpiry,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: access,
		ExpiresIn:   accessExpiry,
	})
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh token"})
		return
	}

	if err := h.authService.Logout(ctx, refresh); err != nil {
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ForgotPassword(ctx, req.Email); err != nil {
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ResetPassword(ctx, req.Token, req.Password); err != nil {
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

// GetCurrentUser returns the authenticated user's detailed info
// @Summary      Get Current User Details
// @Description  Retrieves detailed information about the authenticated user, including university and faculty names (if available)
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]any     "Detailed user information, university/faculty may be null if lookup failed"
// @Failure      401  {object}  dto.ErrorResponse  "User not authenticated"
// @Failure      404  {object}  dto.ErrorResponse  "User not found" // Only if user lookup fails
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /v1/user/me [get]
// @Security     BearerAuth
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userIDVal, exists := c.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context for /me endpoint")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		h.logger.Error("userID in context is not a string", zap.Any("userID", userIDVal))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error processing user ID"})
		return
	}

	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Error("Failed to parse userID from context", zap.String("userIDStr", userIDStr), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := h.authService.GetCurrentUser(ctx, parsedID)
	if err != nil {
		// If the user associated with a valid token isn't found, it's usually a 404.
		logMsg := "User associated with token not found"
		respStatus := http.StatusNotFound
		respError := "User not found"
		if !errors.Is(err, errors.ErrNotFound) {
			logMsg = "Failed to fetch user"
			respStatus = http.StatusInternalServerError
			respError = "Failed to retrieve user data"
		}
		h.logger.Error(logMsg, zap.String("userID", parsedID.String()), zap.Error(err))
		c.JSON(respStatus, gin.H{"error": respError})
		return
	}

	var universityResponse *dto.UniversityResponse
	// Check if UniversityID is valid before attempting fetch
	if user.UniversityID != uuid.Nil {
		universityResponse, err = h.universityService.Get(ctx, user.UniversityID)
		if err != nil {
			h.logger.Warn("Failed to fetch university details for user (continuing)",
				zap.String("userID", user.ID.String()),
				zap.String("universityID", user.UniversityID.String()),
				zap.Error(err))
		}
	} else {
		h.logger.Info("User has nil UniversityID", zap.String("userID", user.ID.String()))
	}

	var facultyResponse *dto.FacultyResponse
	// Check if FacultyID is valid before attempting fetch
	if user.FacultyID != uuid.Nil {
		facultyResponse, err = h.facultyService.Get(user.FacultyID)
		if err != nil {
			h.logger.Warn("Failed to fetch faculty details for user (continuing)",
				zap.String("userID", user.ID.String()),
				zap.String("facultyID", user.FacultyID.String()),
				zap.Error(err))
		}
	} else {
		h.logger.Info("User has nil FacultyID", zap.String("userID", user.ID.String()))
	}

	response := gin.H{
		"id":             user.ID,
		"email":          user.Email,
		"student_id":     user.StudentID,
		"first_name":     user.FirstName,
		"last_name":      user.LastName,
		"email_verified": user.EmailVerified,
		//"is_admin":       user.IsAdmin,
		"university": nil,
		"faculty":    nil,
	}

	if universityResponse != nil {
		response["university"] = gin.H{
			"id":      universityResponse.ID,
			"name_en": universityResponse.NameEn,
			"name_fa": universityResponse.NameFa,
		}
	}

	if facultyResponse != nil {
		response["faculty"] = gin.H{
			"id":      facultyResponse.ID,
			"name_en": facultyResponse.NameEn,
			"name_fa": facultyResponse.NameFa,
		}
	}

	c.JSON(http.StatusOK, response)
}
