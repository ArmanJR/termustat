package dto

import (
	"github.com/google/uuid"
	"time"
)

type ProfessorMinimalResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	NormalizedName string    `json:"normalized_name"`
}

type ProfessorDetailResponse struct {
	ID             uuid.UUID          `json:"id"`
	Name           string             `json:"name"`
	NormalizedName string             `json:"normalized_name"`
	University     UniversityResponse `json:"university"`
	Courses        []CourseResponse   `json:"courses"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type CreateProfessorRequest struct {
	UniversityID uuid.UUID `json:"university_id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
}
