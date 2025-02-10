package dto

import (
	"github.com/google/uuid"
	"time"
)

type ProfessorResponse struct {
	ID             uuid.UUID
	Name           string
	NormalizedName string
	UniversityID   uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UpdateProfessorRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateProfessorRequest struct {
	UniversityID uuid.UUID `json:"university_id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
}
