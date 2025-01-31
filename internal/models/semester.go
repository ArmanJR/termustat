package models

import (
	"github.com/google/uuid"
	"time"
)

type Semester struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Year      int       `gorm:"not null;index"`
	Term      string    `gorm:"not null;size:6;check:term IN ('spring', 'fall');index"`
	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Unique constraint for year + term combination
	UniqueConstraint struct {
		Year int    `gorm:"uniqueIndex:idx_year_term"`
		Term string `gorm:"uniqueIndex:idx_year_term"`
	}
}
