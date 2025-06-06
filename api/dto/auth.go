package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=8"`
	StudentID    string `json:"student_id" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	UniversityID string `json:"university_id" binding:"omitempty,uuid4"`
	FacultyID    string `json:"faculty_id" binding:"omitempty,uuid4"`
	Gender       string `json:"gender" binding:"required,oneof=male female"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required,uuid4"`
	Password string `json:"password" binding:"required,min=8"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// Service DTOs (for internal use)
type RegisterServiceRequest struct {
	Email        string
	Password     string
	StudentID    string
	FirstName    string
	LastName     string
	UniversityID uuid.UUID
	FacultyID    uuid.UUID
	Gender       string
}
