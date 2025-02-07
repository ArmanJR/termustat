package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateSemesterRequest struct {
	Year int    `json:"year" binding:"required,min=1900,max=2200"`
	Term string `json:"term" binding:"required,oneof=spring fall"`
}

type UpdateSemesterRequest struct {
	Year int    `json:"year" binding:"required,min=1900,max=2200"`
	Term string `json:"term" binding:"required,oneof=spring fall"`
}

type SemesterResponse struct {
	ID        uuid.UUID `json:"id"`
	Year      int       `json:"year"`
	Term      string    `json:"term"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
