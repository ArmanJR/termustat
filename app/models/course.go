package models

import (
	"github.com/google/uuid"
	"time"
)

type Course struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UniversityID      uuid.UUID `gorm:"type:uuid;not null;index"`
	FacultyID         uuid.UUID `gorm:"type:uuid;not null;index"`
	ProfessorID       uuid.UUID `gorm:"type:uuid;not null;index"`
	SemesterID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Code              string    `gorm:"not null;size:50"`
	Name              string    `gorm:"not null;size:255"`
	Weight            int       `gorm:"not null"`
	Capacity          int
	GenderRestriction string `gorm:"size:6;check:gender_restriction IN ('male', 'female', 'mixed')"`
	ExamStart         time.Time
	ExamEnd           time.Time
	CourseTimes       []CourseTime
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}
