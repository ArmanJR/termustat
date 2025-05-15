package repositories

import (
	"context"
	"github.com/armanjr/termustat/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByEmailOrStudentID(ctx context.Context, email, studentID string) (*models.User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	CreateEmailVerification(ctx context.Context, verification *models.EmailVerification) error
	CreatePasswordReset(ctx context.Context, reset *models.PasswordReset) error
	FindPasswordResetByToken(ctx context.Context, token string) (*models.PasswordReset, error)
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error
	DeletePasswordReset(ctx context.Context, reset *models.PasswordReset) error
	FindEmailVerificationByToken(ctx context.Context, token string) (*models.EmailVerification, error)
	VerifyUserEmail(ctx context.Context, userID uuid.UUID) error
	DeleteEmailVerification(ctx context.Context, verification *models.EmailVerification) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *authRepository) FindUserByEmailOrStudentID(ctx context.Context, email, studentID string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ? OR student_id = ?", email, studentID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	return &user, err
}

func (r *authRepository) CreateEmailVerification(ctx context.Context, verification *models.EmailVerification) error {
	return r.db.WithContext(ctx).Create(verification).Error
}

func (r *authRepository) CreatePasswordReset(ctx context.Context, reset *models.PasswordReset) error {
	return r.db.WithContext(ctx).Create(reset).Error
}

func (r *authRepository) FindPasswordResetByToken(ctx context.Context, token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.WithContext(ctx).Where("token = ? AND expires_at > ?", token, time.Now()).First(&reset).Error
	return &reset, err
}

func (r *authRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("password_hash", hashedPassword).Error
}

func (r *authRepository) DeletePasswordReset(ctx context.Context, reset *models.PasswordReset) error {
	return r.db.WithContext(ctx).Delete(reset).Error
}

func (r *authRepository) FindEmailVerificationByToken(ctx context.Context, token string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&verification).Error
	return &verification, err
}

func (r *authRepository) VerifyUserEmail(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("email_verified", true).Error
}

func (r *authRepository) DeleteEmailVerification(ctx context.Context, verification *models.EmailVerification) error {
	return r.db.WithContext(ctx).Delete(verification).Error
}
