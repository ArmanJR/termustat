package services

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/infrastructure/mailer"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/repositories"
	"github.com/armanjr/termustat/api/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type AuthService interface {
	Register(req *dto.RegisterServiceRequest) error
	Login(email, password string) (string, int, string, int, error)
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
	GetCurrentUser(userID uuid.UUID) (*models.User, error)
	VerifyEmail(token string) error
	ValidateToken(token string) (*utils.JWTClaims, error)
	Refresh(oldToken string) (string, int, string, int, error)
	Logout(refreshToken string) error
}

type authService struct {
	repo        repositories.AuthRepository
	mailer      mailer.Mailer
	logger      *zap.Logger
	jwtSecret   string
	jwtTTL      time.Duration
	frontendURL string
	refreshRepo repositories.RefreshTokenRepository
	refreshTTL  time.Duration
}

func NewAuthService(
	repo repositories.AuthRepository,
	refreshRepo repositories.RefreshTokenRepository,
	mailer mailer.Mailer,
	logger *zap.Logger,
	jwtSecret string,
	jwtTTL time.Duration,
	refreshTTL time.Duration,
	frontendURL string,
) AuthService {
	return &authService{
		repo:        repo,
		mailer:      mailer,
		logger:      logger,
		jwtSecret:   jwtSecret,
		jwtTTL:      jwtTTL,
		frontendURL: frontendURL,
		refreshRepo: refreshRepo,
		refreshTTL:  refreshTTL,
	}
}

func (s *authService) Register(req *dto.RegisterServiceRequest) error {
	user, err := s.repo.FindUserByEmailOrStudentID(req.Email, req.StudentID)
	if err == nil && user != nil {
		return errors.New("email or student ID already exists")
	}
	// Handle err not being nil. If err is gorm.ErrRecordNotFound, then no duplicate exists.
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrapf(err, "failed to hash password")
	}

	user = &models.User{
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		StudentID:     req.StudentID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		UniversityID:  req.UniversityID,
		FacultyID:     req.FacultyID,
		Gender:        req.Gender,
		EmailVerified: false,
		IsAdmin:       false,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return errors.Wrapf(err, "failed to create user")
	}

	// Create and store the email verification record.
	token := uuid.NewString()
	verification := &models.EmailVerification{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.repo.CreateEmailVerification(verification); err != nil {
		return errors.Wrapf(err, "failed to create email verification")
	}

	// Pass the token to the mailer to include it in the verification email.
	if err := s.mailer.SendVerificationEmail(user, token); err != nil {
		s.logger.Error("Failed to send verification email",
			zap.String("email", user.Email),
			zap.Error(err))
	}

	return nil
}

func (s *authService) Login(email, password string) (string, int, string, int, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Login attempt failed: user not found", zap.String("email", email))
			return "", 0, "", 0, errors.New("invalid credentials")
		}
		s.logger.Error("Database error during login", zap.String("email", email), zap.Error(err))
		return "", 0, "", 0, errors.New("failed to login")
	}

	if !user.EmailVerified {
		s.logger.Warn("Login attempt failed: email not verified", zap.String("email", email), zap.String("user_id", user.ID.String()))
		return "", 0, "", 0, errors.New("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.logger.Warn("Login attempt failed: invalid password", zap.String("email", email), zap.String("user_id", user.ID.String()))
		return "", 0, "", 0, errors.New("invalid credentials")
	}

	var scopes []string
	if user.IsAdmin {
		scopes = append(scopes, "admin-dashboard")
	}

	accessExpirySeconds := int(s.jwtTTL.Seconds())
	access, err := utils.GenerateJWT(user.ID.String(), scopes, s.jwtSecret, s.jwtTTL)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.String("user_id", user.ID.String()), zap.Error(err))
		return "", 0, "", 0, errors.New("failed to generate access token")
	}

	refreshStr, err := generateRefreshString()
	if err != nil {
		s.logger.Error("Failed to generate refresh token string", zap.String("user_id", user.ID.String()), zap.Error(err))
		return "", 0, "", 0, errors.New("failed to generate refresh token")
	}

	refreshExpirySeconds := int(s.refreshTTL.Seconds())
	rt := &models.RefreshToken{
		Token:     refreshStr,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.refreshRepo.Create(rt); err != nil {
		s.logger.Error("Failed to store refresh token", zap.String("user_id", user.ID.String()), zap.Error(err))
		return "", 0, "", 0, errors.New("failed to store refresh token")
	}

	return access, accessExpirySeconds, refreshStr, refreshExpirySeconds, nil
}

func (s *authService) ForgotPassword(email string) error {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil // Don't reveal if email exists
	}

	resetToken := uuid.New()
	resetExpiry := time.Now().Add(time.Hour)

	passwordReset := &models.PasswordReset{
		Token:     resetToken,
		UserID:    user.ID,
		ExpiresAt: resetExpiry,
	}

	if err := s.repo.CreatePasswordReset(passwordReset); err != nil {
		return errors.Wrapf(err, "failed to create password reset")
	}

	if err := s.sendPasswordResetEmail(user, resetToken.String()); err != nil {
		s.logger.Error("Failed to send password reset email",
			zap.String("email", user.Email),
			zap.Error(err))
		return errors.Wrapf(err, "failed to send password reset email")
	}

	return nil
}

func (s *authService) ResetPassword(token, newPassword string) error {
	reset, err := s.repo.FindPasswordResetByToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrapf(err, "failed to hash password")
	}

	if err := s.repo.UpdateUserPassword(reset.UserID, string(hashedPassword)); err != nil {
		return errors.Wrapf(err, "failed to update password")
	}

	if err := s.repo.DeletePasswordReset(reset); err != nil {
		s.logger.Error("Failed to delete password reset",
			zap.String("token", token),
			zap.Error(err))
	}

	return nil
}

func (s *authService) GetCurrentUser(userID uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find user")
	}
	return user, nil
}

func (s *authService) VerifyEmail(token string) error {
	verification, err := s.repo.FindEmailVerificationByToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if verification.ExpiresAt.Before(time.Now()) {
		return errors.New("invalid or expired token")
	}

	if err := s.repo.VerifyUserEmail(verification.UserID); err != nil {
		return errors.Wrapf(err, "failed to verify user email")
	}

	if err := s.repo.DeleteEmailVerification(verification); err != nil {
		s.logger.Error("Failed to delete email verification",
			zap.String("token", token),
			zap.Error(err))
	}

	return nil
}

func (s *authService) ValidateToken(token string) (*utils.JWTClaims, error) {
	claims, err := utils.ParseJWT(token, s.jwtSecret)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrExpiredToken):
			s.logger.Warn("Token has expired")
			return nil, errors.Wrap(err, "token has expired")
		default:
			s.logger.Warn("Invalid token",
				zap.Error(err))
			return nil, errors.Wrap(err, "invalid token")
		}
	}

	// Verify user still exists and is active
	_, err = s.repo.FindUserByID(uuid.MustParse(claims.UserID))
	if err != nil {
		s.logger.Warn("User from token not found",
			zap.String("user_id", claims.UserID),
			zap.Error(err))
		return nil, errors.New("invalid token: user not found")
	}

	return claims, nil
}

func (s *authService) Refresh(old string) (string, int, string, int, error) {
	rt, err := s.refreshRepo.Find(old)
	if err != nil {
		s.logger.Warn("Invalid or expired refresh token provided", zap.String("token_prefix", old[:min(10, len(old))]))
		return "", 0, "", 0, errors.New("invalid refresh token")
	}

	if err := s.refreshRepo.Revoke(rt.ID); err != nil {
		s.logger.Error("Failed to revoke old refresh token", zap.String("token_id", rt.ID.String()), zap.Error(err))
		return "", 0, "", 0, errors.Wrap(err, "failed to revoke token")
	}

	user := rt.User // Or user, err := s.repo.FindUserByID(rt.UserID); if err != nil { ... }

	var scopes []string
	if user.IsAdmin {
		scopes = append(scopes, "admin-dashboard")
	}

	accessExpirySeconds := int(s.jwtTTL.Seconds())
	newAccess, err := utils.GenerateJWT(rt.UserID.String(), scopes, s.jwtSecret, s.jwtTTL)
	if err != nil {
		s.logger.Error("Failed to generate new access token during refresh", zap.String("user_id", rt.UserID.String()), zap.Error(err))
		return "", 0, "", 0, errors.Wrap(err, "failed to generate access token")
	}

	newRefresh, err := generateRefreshString()
	if err != nil {
		s.logger.Error("Failed to generate new refresh token string during refresh", zap.String("user_id", rt.UserID.String()), zap.Error(err))
		return "", 0, "", 0, errors.Wrap(err, "failed to generate refresh token")
	}

	refreshExpirySeconds := int(s.refreshTTL.Seconds())
	newRT := &models.RefreshToken{
		Token:     newRefresh,
		UserID:    rt.UserID,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}

	if err := s.refreshRepo.Create(newRT); err != nil {
		s.logger.Error("Failed to store new refresh token during refresh", zap.String("user_id", rt.UserID.String()), zap.Error(err))
		return "", 0, "", 0, errors.Wrap(err, "failed to store refresh token")
	}

	s.logger.Info("Token refreshed successfully", zap.String("user_id", rt.UserID.String()))
	return newAccess, accessExpirySeconds, newRefresh, refreshExpirySeconds, nil
}

func (s *authService) Logout(refreshToken string) error {
	rt, err := s.refreshRepo.Find(refreshToken)
	if err != nil {
		return nil
	} // idempotent
	return s.refreshRepo.Revoke(rt.ID)
}

// Helpers
const refreshByteLen = 64

func generateRefreshString() (string, error) {
	b := make([]byte, refreshByteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
