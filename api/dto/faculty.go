package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateFacultyDTO struct {
	UniversityID uuid.UUID `json:"university_id" binding:"required"`
	NameEn       string    `json:"name_en" binding:"required"`
	NameFa       string    `json:"name_fa" binding:"required"`
	ShortCode    string    `json:"short_code" binding:"required,max=10"`
	IsActive     bool      `json:"is_active"`
}

type UpdateFacultyDTO struct {
	UniversityID uuid.UUID `json:"university_id" binding:"required"`
	NameEn       string    `json:"name_en" binding:"required"`
	NameFa       string    `json:"name_fa" binding:"required"`
	ShortCode    string    `json:"short_code" binding:"required,max=10"`
	IsActive     bool      `json:"is_active"`
}

type FacultyResponse struct {
	ID           uuid.UUID `json:"id"`
	UniversityID uuid.UUID `json:"university_id"`
	NameEn       string    `json:"name_en"`
	NameFa       string    `json:"name_fa"`
	ShortCode    string    `json:"short_code"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
