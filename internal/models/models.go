package models

import (
	"github.com/google/uuid"
	"time"
)

type University struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"not null"`
	Faculties []Faculty
	CreatedAt time.Time
}

type Faculty struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UniversityID uuid.UUID `gorm:"type:uuid;not null"`
	Name         string    `gorm:"not null"`
	ShortCode    string    `gorm:"not null"`
	CreatedAt    time.Time
}

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email         string    `gorm:"uniqueIndex;not null"`
	PasswordHash  string    `gorm:"not null"`
	StudentID     string    `gorm:"not null"`
	FirstName     string
	LastName      string
	UniversityID  uuid.UUID `gorm:"type:uuid;not null"`
	FacultyID     uuid.UUID `gorm:"type:uuid;not null"`
	Gender        string    `gorm:"check:gender IN ('male', 'female')"`
	EmailVerified bool      `gorm:"default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Course struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UniversityID      uuid.UUID `gorm:"type:uuid;not null"`
	FacultyID         uuid.UUID `gorm:"type:uuid;not null"`
	Code              string    `gorm:"not null"`
	Name              string    `gorm:"not null"`
	Weight            int       `gorm:"not null"`
	Capacity          int
	GenderRestriction string `gorm:"check:gender_restriction IN ('male', 'female', 'mixed')"`
	Professor         string
	ExamStart         time.Time
	ExamEnd           time.Time
	Schedules         []CourseSchedule
	CreatedAt         time.Time
}

type CourseSchedule struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CourseID  uuid.UUID `gorm:"type:uuid;not null"`
	DayOfWeek int       `gorm:"check:day_of_week BETWEEN 0 AND 6"`
	StartTime time.Time `gorm:"type:time;not null"`
	EndTime   time.Time `gorm:"type:time;not null"`
}

type UserCourse struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	CourseID  uuid.UUID `gorm:"type:uuid;not null"`
	Semester  string    `gorm:"not null"`
	CreatedAt time.Time
}

type PasswordReset struct {
	Token     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
