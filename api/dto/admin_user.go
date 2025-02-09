package dto

import "github.com/google/uuid"

type AdminUpdateUserRequest struct {
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	UniversityID uuid.UUID `json:"university_id"`
	FacultyID    uuid.UUID `json:"faculty_id"`
	Gender       string    `json:"gender" binding:"omitempty,oneof=male female"`
	Password     string    `json:"password" binding:"omitempty,min=8"`
}
