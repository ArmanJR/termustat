package models

import (
	"github.com/google/uuid"
	"time"
)

type Semester struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Year      int       `gorm:"not null;index"`
	Term      string    `gorm:"not null;size:6;check:term IN ('spring', 'fall');index"`
	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
