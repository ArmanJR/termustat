package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email         string    `gorm:"uniqueIndex;not null;size:255"`
	PasswordHash  string    `gorm:"not null"`
	StudentID     string    `gorm:"not null;uniqueIndex;size:20"`
	FirstName     string    `gorm:"size:100"`
	LastName      string    `gorm:"size:100"`
	UniversityID  uuid.UUID `gorm:"type:uuid;not null;index"`
	FacultyID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Gender        string    `gorm:"size:6;check:gender IN ('male', 'female')"`
	EmailVerified bool      `gorm:"default:false"`
	IsAdmin       bool      `gorm:"default:false"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
