package models

import (
	"github.com/google/uuid"
	"time"
)

type Professor struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UniversityID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	Name           string     `gorm:"not null;size:255"`
	NormalizedName string     `gorm:"not null;size:255;index"`
	University     University `gorm:"foreignKey:UniversityID"`
	Courses        []Course   `gorm:"foreignKey:ProfessorID"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
}
