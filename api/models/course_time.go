package models

import (
	"github.com/google/uuid"
	"time"
)

type CourseTime struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CourseID  uuid.UUID `gorm:"type:uuid;not null;index"`
	DayOfWeek int       `gorm:"check:day_of_week BETWEEN 0 AND 6"`
	StartTime time.Time `gorm:"type:time;not null"`
	EndTime   time.Time `gorm:"type:time;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (CourseTime) TableName() string {
	return "course_times"
}
