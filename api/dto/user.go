package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateUserRequest struct {
	Email        string    `json:"email" binding:"required,email"`
	Password     string    `json:"password" binding:"required,min=8"`
	StudentID    string    `json:"student_id" binding:"required"`
	FirstName    string    `json:"first_name" binding:"required"`
	LastName     string    `json:"last_name" binding:"required"`
	UniversityID uuid.UUID `json:"university_id" binding:"required"`
	FacultyID    uuid.UUID `json:"faculty_id" binding:"required"`
	Gender       string    `json:"gender" binding:"required,oneof=male female"`
}

type UpdateUserRequest struct {
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	UniversityID uuid.UUID `json:"university_id"`
	FacultyID    uuid.UUID `json:"faculty_id"`
	Gender       string    `json:"gender" binding:"omitempty,oneof=male female"`
}

type UpdatePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	StudentID     string    `json:"student_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	UniversityID  uuid.UUID `json:"university_id"`
	FacultyID     uuid.UUID `json:"faculty_id"`
	Gender        string    `json:"gender"`
	EmailVerified bool      `json:"email_verified"`
	IsAdmin       bool      `json:"is_admin"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
