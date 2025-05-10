package models

import (
	"github.com/google/uuid"
	"time"
)

type Faculty struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UniversityID uuid.UUID `gorm:"type:uuid;not null;index"`
	NameEn       string    `gorm:"not null"`
	NameFa       string    `gorm:"not null"`
	ShortCode    string    `gorm:"not null;size:10"`
	IsActive     bool      `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
