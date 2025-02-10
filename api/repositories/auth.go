package repositories

import (
	"github.com/armanjr/termustat/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserByEmailOrStudentID(email, studentID string) (*models.User, error)
	FindUserByID(id uuid.UUID) (*models.User, error)
	CreateEmailVerification(verification *models.EmailVerification) error
	CreatePasswordReset(reset *models.PasswordReset) error
	FindPasswordResetByToken(token string) (*models.PasswordReset, error)
	UpdateUserPassword(userID uuid.UUID, hashedPassword string) error
	DeletePasswordReset(reset *models.PasswordReset) error
	FindEmailVerificationByToken(token string) (*models.EmailVerification, error)
	VerifyUserEmail(userID uuid.UUID) error
	DeleteEmailVerification(verification *models.EmailVerification) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *authRepository) FindUserByEmailOrStudentID(email, studentID string) (*models.User, error) {
	var count int64
	err := r.db.Model(&models.User{}).
		Where("email = ? OR student_id = ?", email, studentID).
		Count(&count).Error
	return nil, err
}

func (r *authRepository) FindUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *authRepository) CreateEmailVerification(verification *models.EmailVerification) error {
	return r.db.Create(verification).Error
}

func (r *authRepository) CreatePasswordReset(reset *models.PasswordReset) error {
	return r.db.Create(reset).Error
}

func (r *authRepository) FindPasswordResetByToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&reset).Error
	return &reset, err
}

func (r *authRepository) UpdateUserPassword(userID uuid.UUID, hashedPassword string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("password_hash", hashedPassword).Error
}

func (r *authRepository) DeletePasswordReset(reset *models.PasswordReset) error {
	return r.db.Delete(reset).Error
}

func (r *authRepository) FindEmailVerificationByToken(token string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.Where("token = ?", token).First(&verification).Error
	return &verification, err
}

func (r *authRepository) VerifyUserEmail(userID uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("email_verified", true).Error
}

func (r *authRepository) DeleteEmailVerification(verification *models.EmailVerification) error {
	return r.db.Delete(verification).Error
}
