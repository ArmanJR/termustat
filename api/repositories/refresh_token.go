package repositories

import (
	"github.com/armanjr/termustat/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type RefreshTokenRepository interface {
	Create(t *models.RefreshToken) error
	Find(token string) (*models.RefreshToken, error)
	Revoke(id uuid.UUID) error
	RevokeAllForUser(userID uuid.UUID) error
	CleanupExpired() error
}

type refreshTokenRepository struct{ db *gorm.DB }

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db}
}

func (r *refreshTokenRepository) Create(t *models.RefreshToken) error {
	return r.db.Create(t).Error
}

func (r *refreshTokenRepository) Find(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.Preload("User").
		Where("token = ? AND revoked = FALSE AND expires_at > ?", token, time.Now()).
		First(&rt).Error
	return &rt, err
}

func (r *refreshTokenRepository) Revoke(id uuid.UUID) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}

func (r *refreshTokenRepository) RevokeAllForUser(userID uuid.UUID) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

func (r *refreshTokenRepository) CleanupExpired() error {
	return r.db.Where("expires_at <= ?", time.Now()).
		Delete(&models.RefreshToken{}).Error
}
