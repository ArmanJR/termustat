package models

import (
	"github.com/google/uuid"
	"time"
)

type UserCourse struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	CourseID   uuid.UUID `gorm:"type:uuid;not null;index"`
	SemesterID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
