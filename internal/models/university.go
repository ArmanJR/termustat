package models

import (
	"github.com/google/uuid"
	"time"
)

type University struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"not null"`
	IsActive  bool      `gorm:"default:true;index"`
	Faculties []Faculty
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
