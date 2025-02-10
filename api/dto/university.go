package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateUniversityRequest struct {
	NameEn   string `json:"name_en" binding:"required"`
	NameFa   string `json:"name_fa" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

type UpdateUniversityRequest struct {
	NameEn   string `json:"name_en" binding:"required"`
	NameFa   string `json:"name_fa" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

type UniversityResponse struct {
	ID        uuid.UUID `json:"id"`
	NameEn    string    `json:"name_en"`
	NameFa    string    `json:"name_fa"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
