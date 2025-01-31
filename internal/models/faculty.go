package models

import (
	"github.com/google/uuid"
	"time"
)

type Faculty struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UniversityID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name         string    `gorm:"not null"`
	ShortCode    string    `gorm:"not null;size:10"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
