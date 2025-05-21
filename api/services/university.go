package services

import (
	"context"
	"fmt"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strings"
)

type UniversityService interface {
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateUniversityRequest) (*dto.UniversityResponse, error)
	Create(ctx context.Context, req *dto.CreateUniversityRequest) (*dto.UniversityResponse, error)
	Get(ctx context.Context, id uuid.UUID) (*dto.UniversityResponse, error)
	GetAll(ctx context.Context) ([]dto.UniversityResponse, error)
	ExistsByName(ctx context.Context, nameEn, nameFa string) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type universityService struct {
	repo   repositories.UniversityRepository
	logger *zap.Logger
}

func NewUniversityService(repo repositories.UniversityRepository, logger *zap.Logger) UniversityService {
	return &universityService{
		repo:   repo,
		logger: logger,
	}
}

func (s *universityService) Get(ctx context.Context, id uuid.UUID) (*dto.UniversityResponse, error) {
	university, err := s.repo.Find(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return nil, err
		default:
			s.logger.Error("Failed to fetch university",
				zap.String("id", id.String()),
				zap.String("operation", "Get"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to get university")
		}
	}
	response := dto.UniversityResponse{
		ID:        university.ID,
		NameEn:    university.NameEn,
		NameFa:    university.NameFa,
		IsActive:  university.IsActive,
		CreatedAt: university.CreatedAt,
		UpdatedAt: university.UpdatedAt,
	}
	return &response, nil
}

func (s *universityService) Create(ctx context.Context, req *dto.CreateUniversityRequest) (*dto.UniversityResponse, error) {
	if req.NameEn == "" {
		return nil, errors.NewValidationError("name_en")
	}
	if req.NameFa == "" {
		return nil, errors.NewValidationError("name_fa")
	}
	if req.IsActive == nil {
		return nil, errors.NewValidationError("is_active")
	}

	university := &models.University{
		NameEn:   strings.TrimSpace(req.NameEn),
		NameFa:   strings.TrimSpace(req.NameFa),
		IsActive: *req.IsActive,
	}

	created, err := s.repo.Create(ctx, university)
	if err != nil {
		s.logger.Error("Failed to create university",
			zap.String("name_en", req.NameEn),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create university: %w", err)
	}

	response := dto.UniversityResponse{
		ID:        created.ID,
		NameEn:    created.NameEn,
		NameFa:    created.NameFa,
		IsActive:  created.IsActive,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}
	return &response, nil
}

func (s *universityService) GetAll(ctx context.Context) ([]dto.UniversityResponse, error) {
	universities, err := s.repo.FindAll(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch universities", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch universities: %w", err)
	}

	response := make([]dto.UniversityResponse, len(universities))
	for i, univ := range universities {
		response[i] = dto.UniversityResponse{
			ID:        univ.ID,
			NameEn:    univ.NameEn,
			NameFa:    univ.NameFa,
			IsActive:  univ.IsActive,
			CreatedAt: univ.CreatedAt,
			UpdatedAt: univ.UpdatedAt,
		}
	}
	return response, nil
}

func (s *universityService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateUniversityRequest) (*dto.UniversityResponse, error) {
	if req.NameEn == "" {
		return nil, errors.NewValidationError("name_en")
	}
	if req.NameFa == "" {
		return nil, errors.NewValidationError("name_fa")
	}
	if req.IsActive == nil {
		return nil, errors.NewValidationError("is_active")
	}

	university, err := s.repo.Find(ctx, id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}
		s.logger.Error("Failed to fetch university for update",
			zap.String("id", id.String()),
			zap.String("service", "University"),
			zap.String("operation", "Update"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to update university: %w", err)
	}

	university.NameEn = strings.TrimSpace(req.NameEn)
	university.NameFa = strings.TrimSpace(req.NameFa)
	university.IsActive = *req.IsActive

	updated, err := s.repo.Update(ctx, university)
	if err != nil {
		s.logger.Error("Failed to update university",
			zap.String("id", id.String()),
			zap.String("service", "University"),
			zap.String("operation", "Update"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to update university: %w", err)
	}

	response := dto.UniversityResponse{
		ID:        updated.ID,
		NameEn:    updated.NameEn,
		NameFa:    updated.NameFa,
		IsActive:  updated.IsActive,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}
	return &response, nil
}

func (s *universityService) ExistsByName(ctx context.Context, nameEn, nameFa string) (bool, error) {
	exists, err := s.repo.ExistsByName(ctx, nameEn, nameFa)
	if err != nil {
		s.logger.Error("Failed to check university existence",
			zap.String("name_en", nameEn),
			zap.String("name_fa", nameFa),
			zap.Error(err))
		return false, fmt.Errorf("failed to check university existence: %w", err)
	}
	return exists, nil
}

func (s *universityService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.Find(ctx, id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return err
		}
		s.logger.Error("Failed to fetch university for deletion",
			zap.String("id", id.String()),
			zap.String("service", "University"),
			zap.String("operation", "Delete"),
			zap.Error(err))
		return fmt.Errorf("failed to delete university: %w", err)
	}

	if err = s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete university",
			zap.String("id", id.String()),
			zap.String("service", "University"),
			zap.String("operation", "Delete"),
			zap.Error(err))
		return fmt.Errorf("failed to delete university: %w", err)
	}

	return nil
}
