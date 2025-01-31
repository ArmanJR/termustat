package models

import (
	"github.com/google/uuid"
	"time"
)

type EmailVerification struct {
	Token     string    `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (EmailVerification) TableName() string {
	return "email_verifications"
}
