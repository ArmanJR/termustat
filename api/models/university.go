package models

import (
	"github.com/google/uuid"
	"time"
)

type University struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	NameEn    string    `gorm:"not null"`
	NameFa    string    `gorm:"not null"`
	IsActive  bool      `gorm:"not null"`
	Faculties []Faculty
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
