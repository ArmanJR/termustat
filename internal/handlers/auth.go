package handlers

import (
	"errors"
	"fmt"
	"github.com/armanjr/termustat/internal/config"
	"github.com/armanjr/termustat/pkg/utils"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"

	"github.com/armanjr/termustat/internal/models"
	"github.com/armanjr/termustat/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	mailer *services.Mailer
	config *config.Config
	logger *zap.Logger
}

func NewAuthHandler(mailer *services.Mailer, cfg *config.Config, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		mailer: mailer,
		config: cfg,
		logger: logger,
	}
}

// RegisterRequest Register Request
type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=8"`
	StudentID    string `json:"student_id" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	UniversityID string `json:"university_id" binding:"required,uuid4"`
	FacultyID    string `json:"faculty_id" binding:"required,uuid4"`
	Gender       string `json:"gender" binding:"required,oneof=male female"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid registration request",
			zap.Error(err),
			zap.Any("request", req))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var count int64
	config.DB.Model(&models.User{}).
		Where("email = ? OR student_id = ?", req.Email, req.StudentID).
		Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email or Student ID already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		StudentID:     req.StudentID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		UniversityID:  uuid.MustParse(req.UniversityID),
		FacultyID:     uuid.MustParse(req.FacultyID),
		Gender:        req.Gender,
		EmailVerified: false,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	go h.sendVerificationEmail(&user)

	h.logger.Info("New user registered",
		zap.String("email", req.Email),
		zap.String("student_id", req.StudentID))

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful. Please check your email to verify your account."})
}

func (h *AuthHandler) sendVerificationEmail(user *models.User) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	verification := models.EmailVerification{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	if err := config.DB.Create(&verification).Error; err != nil {
		h.logger.Error("Failed to create verification record", zap.Error(err))
		return
	}

	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", h.config.FrontendURL, token)

	tplData := struct {
		Name            string
		VerificationURL string
	}{
		Name:            user.FirstName,
		VerificationURL: verificationURL,
	}

	emailContent, err := h.mailer.RenderTemplate("verification_email.html", tplData)
	if err != nil {
		h.logger.Error("Failed to render verification email template", zap.Error(err))
		return
	}

	if err := h.mailer.SendEmail(user.Email, emailContent.Subject, emailContent.Body); err != nil {
		h.logger.Error("Failed to send verification email",
			zap.String("email", user.Email),
			zap.Error(err))
		return
	}

	h.logger.Info("Verification email sent successfully",
		zap.String("email", user.Email),
		zap.Time("expires_at", expiresAt))
}

// LoginRequest Login Request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.EmailVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Email not verified"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID.String(), h.config.JWTSecret, h.config.JWTTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ForgotPasswordRequest Forgot Password Request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Don't reveal if user exists
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link will be sent"})
		return
	}

	resetToken := uuid.New()
	resetExpiry := time.Now().Add(time.Hour)

	passwordReset := models.PasswordReset{
		Token:     resetToken,
		UserID:    user.ID,
		ExpiresAt: resetExpiry,
	}

	if err := config.DB.Create(&passwordReset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reset token"})
		return
	}

	go h.sendPasswordResetEmail(&user, resetToken.String())

	c.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent to your email"})
}

func (h *AuthHandler) sendPasswordResetEmail(user *models.User, token string) {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", h.config.FrontendURL, token)
	tplData := struct{ ResetURL string }{ResetURL: resetURL}

	emailContent, err := h.mailer.RenderTemplate("password_reset_email.html", tplData)
	if err != nil {
		return
	}

	err = h.mailer.SendEmail(user.Email, emailContent.Subject, emailContent.Body)
	if err != nil {
		h.logger.Error("Failed to send password reset email",
			zap.String("email", user.Email),
			zap.Error(err))
		return
	}
}

// ResetPasswordRequest Reset Password Request
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required,uuid4"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reset models.PasswordReset
	if err := config.DB.Where("token = ? AND expires_at > ?", req.Token, time.Now()).
		First(&reset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := config.DB.Model(&models.User{}).
		Where("id = ?", reset.UserID).
		Update("password_hash", string(hashedPassword)).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	config.DB.Delete(&reset)

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func (h *AuthHandler) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			return
		}

		claims, err := utils.ParseJWT(tokenString, h.config.JWTSecret)
		if err != nil {
			if errors.Is(err, utils.ErrExpiredToken) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set user ID in context for downstream handlers
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
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
