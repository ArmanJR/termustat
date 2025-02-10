package services

import (
	"fmt"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SemesterService interface {
	Create(req *dto.CreateSemesterRequest) (*dto.SemesterResponse, error)
	Get(id uuid.UUID) (*dto.SemesterResponse, error)
	GetAll() ([]dto.SemesterResponse, error)
	Update(id uuid.UUID, req *dto.UpdateSemesterRequest) (*dto.SemesterResponse, error)
	Delete(id uuid.UUID) error
}

type semesterService struct {
	repo   repositories.SemesterRepository
	logger *zap.Logger
}

func NewSemesterService(repo repositories.SemesterRepository, logger *zap.Logger) SemesterService {
	return &semesterService{
		repo:   repo,
		logger: logger,
	}
}

func (s *semesterService) Create(req *dto.CreateSemesterRequest) (*dto.SemesterResponse, error) {
	if !isValidTerm(req.Term) {
		return nil, errors.NewValidationError("term must be either 'spring' or 'fall'")
	}

	if !isValidYear(req.Year) {
		return nil, errors.NewValidationError("year must be between 1900 and 2200")
	}

	existing, err := s.repo.FindByYearAndTerm(req.Year, req.Term)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		s.logger.Error("Database error while checking existing semester",
			zap.Int("year", req.Year),
			zap.String("term", req.Term),
			zap.Error(err))
		return nil, fmt.Errorf("failed to check existing semester: %w", err)
	}
	if existing != nil {
		return nil, errors.NewConflictError("semester already exists for this year and term")
	}

	semester := &models.Semester{
		Year: req.Year,
		Term: req.Term,
	}

	created, err := s.repo.Create(semester)
	if err != nil {
		s.logger.Error("Failed to create semester",
			zap.Int("year", req.Year),
			zap.String("term", req.Term),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create semester: %w", err)
	}

	return mapSemesterToDTO(created), nil
}

func (s *semesterService) Get(id uuid.UUID) (*dto.SemesterResponse, error) {
	semester, err := s.repo.Find(id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}
		s.logger.Error("Failed to fetch semester",
			zap.String("id", id.String()),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get semester: %w", err)
	}

	return mapSemesterToDTO(semester), nil
}

func (s *semesterService) GetAll() ([]dto.SemesterResponse, error) {
	semesters, err := s.repo.FindAll()
	if err != nil {
		s.logger.Error("Failed to fetch semesters", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch semesters: %w", err)
	}

	response := make([]dto.SemesterResponse, len(semesters))
	for i, semester := range semesters {
		response[i] = *mapSemesterToDTO(&semester)
	}

	return response, nil
}

func (s *semesterService) Update(id uuid.UUID, req *dto.UpdateSemesterRequest) (*dto.SemesterResponse, error) {
	if !isValidTerm(req.Term) {
		return nil, errors.NewValidationError("term must be either 'spring' or 'fall'")
	}

	if !isValidYear(req.Year) {
		return nil, errors.NewValidationError("year must be between 1900 and 2200")
	}

	existing, err := s.repo.Find(id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}
		s.logger.Error("Failed to fetch semester for update",
			zap.String("id", id.String()),
			zap.Error(err))
		return nil, fmt.Errorf("failed to update semester: %w", err)
	}

	if existing.Year != req.Year || existing.Term != req.Term {
		duplicate, err := s.repo.FindByYearAndTerm(req.Year, req.Term)
		if err != nil && !errors.Is(err, errors.ErrNotFound) {
			s.logger.Error("Database error while checking existing semester",
				zap.Int("year", req.Year),
				zap.String("term", req.Term),
				zap.Error(err))
			return nil, fmt.Errorf("failed to check existing semester: %w", err)
		}
		if duplicate != nil {
			return nil, errors.NewConflictError("semester already exists for this year and term")
		}
	}

	existing.Year = req.Year
	existing.Term = req.Term

	updated, err := s.repo.Update(existing)
	if err != nil {
		s.logger.Error("Failed to update semester",
			zap.String("id", id.String()),
			zap.Error(err))
		return nil, fmt.Errorf("failed to update semester: %w", err)
	}

	return mapSemesterToDTO(updated), nil
}

func (s *semesterService) Delete(id uuid.UUID) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return err
		}
		s.logger.Error("Failed to delete semester",
			zap.String("id", id.String()),
			zap.Error(err))
		return fmt.Errorf("failed to delete semester: %w", err)
	}

	return nil
}

func isValidTerm(term string) bool {
	return term == "spring" || term == "fall"
}

func isValidYear(year int) bool {
	return year >= 1900 && year <= 2200
}

func mapSemesterToDTO(semester *models.Semester) *dto.SemesterResponse {
	return &dto.SemesterResponse{
		ID:        semester.ID,
		Year:      semester.Year,
		Term:      semester.Term,
		CreatedAt: semester.CreatedAt,
		UpdatedAt: semester.UpdatedAt,
	}
}
