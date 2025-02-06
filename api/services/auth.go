package services

import (
	"fmt"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/repositories"
	"github.com/armanjr/termustat/api/utils"
	"github.com/google/uuid"
	"github.com/mailgun/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	Register(req *dto.RegisterServiceRequest) error
	Login(email, password string) (string, error)
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
	GetCurrentUser(userID uuid.UUID) (*models.User, error)
	VerifyEmail(token string) error
	//UpdateUser(userID uuid.UUID, req *dto.UpdateUserRequest) error
	//ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error
}

type authService struct {
	repo        repositories.AuthRepository
	mailer      MailerService
	logger      *zap.Logger
	jwtSecret   string
	jwtTTL      time.Duration
	frontendURL string
}

func NewAuthService(
	repo repositories.AuthRepository,
	mailer MailerService,
	logger *zap.Logger,
	jwtSecret string,
	jwtTTL time.Duration,
	frontendURL string,
) AuthService {
	return &authService{
		repo:        repo,
		mailer:      mailer,
		logger:      logger,
		jwtSecret:   jwtSecret,
		jwtTTL:      jwtTTL,
		frontendURL: frontendURL,
	}
}

func (s *authService) Register(req *dto.RegisterServiceRequest) error {
	_, err := s.repo.FindUserByEmailOrStudentID(req.Email, req.StudentID)
	if err == nil {
		return errors.New("email or student ID already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrapf(err, "failed to hash password")
	}

	user := &models.User{
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

	if err := s.sendVerificationEmail(user); err != nil {
		s.logger.Error("Failed to send verification email",
			zap.String("email", user.Email),
			zap.Error(err))
	}

	return nil
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.EmailVerified {
		return "", errors.New("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID.String(), s.jwtSecret, int(s.jwtTTL.Seconds()))
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate token")
	}

	return token, nil
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

//func (s *authService) UpdateUser(userID uuid.UUID, req *dto.UpdateUserRequest) error {
//	user, err := s.repo.FindUserByID(userID)
//	if err != nil {
//		return errors.Wrapf(err, "failed to find user")
//	}
//
//	user.FirstName = req.FirstName
//	user.LastName = req.LastName
//	user.Gender = req.Gender
//
//	if err := s.repo.UpdateUser(user); err != nil {
//		return errors.Wrapf(err, "failed to update user")
//	}
//
//	return nil
//}
//
//func (s *authService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
//	user, err := s.repo.FindUserByID(userID)
//	if err != nil {
//		return errors.Wrapf(err, "failed to find user")
//	}
//
//	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
//		return errors.New("invalid old password")
//	}
//
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
//	if err != nil {
//		return errors.Wrapf(err, "failed to hash password")
//	}
//
//	if err := s.repo.UpdateUserPassword(user.ID, string(hashedPassword)); err != nil {
//		return errors.Wrapf(err, "failed to update password")
//	}
//
//	return nil
//}

func (s *authService) sendVerificationEmail(user *models.User) error {
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	verification := &models.EmailVerification{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	if err := s.repo.CreateEmailVerification(verification); err != nil {
		return errors.Wrapf(err, "failed to create verification record")
	}

	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, token)
	tplData := struct {
		Name            string
		VerificationURL string
	}{
		Name:            user.FirstName,
		VerificationURL: verificationURL,
	}

	emailContent, err := s.mailer.RenderTemplate("verification_email.html", tplData)
	if err != nil {
		return errors.Wrapf(err, "failed to render verification email template")
	}

	if err := s.mailer.SendEmail(user.Email, emailContent.Subject, emailContent.Body); err != nil {
		return errors.Wrapf(err, "failed to send verification email")
	}

	return nil
}

func (s *authService) sendPasswordResetEmail(user *models.User, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, token)
	tplData := struct {
		ResetURL string
	}{ResetURL: resetURL}

	emailContent, err := s.mailer.RenderTemplate("password_reset_email.html", tplData)
	if err != nil {
		return errors.Wrapf(err, "failed to render password reset email template")
	}

	if err := s.mailer.SendEmail(user.Email, emailContent.Subject, emailContent.Body); err != nil {
		return errors.Wrapf(err, "failed to send password reset email")
	}

	return nil
}
