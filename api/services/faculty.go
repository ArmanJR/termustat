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

type FacultyService interface {
	Create(dto dto.CreateFacultyDTO) (*dto.FacultyResponse, error)
	Get(id uuid.UUID) (*dto.FacultyResponse, error)
	GetAllByUniversity(universityID uuid.UUID) ([]*dto.FacultyResponse, error)
	GetByUniversityAndShortCode(universityID uuid.UUID, shortCode string) (*dto.FacultyResponse, error)
	Update(id uuid.UUID, dto dto.UpdateFacultyDTO) (*dto.FacultyResponse, error)
	Delete(id uuid.UUID) error
}

type facultyService struct {
	facultyRepo       repositories.FacultyRepository
	universityService UniversityService
	logger            *zap.Logger
}

func NewFacultyService(
	facultyRepo repositories.FacultyRepository,
	universityService UniversityService,
	logger *zap.Logger,
) FacultyService {
	return &facultyService{
		facultyRepo:       facultyRepo,
		universityService: universityService,
		logger:            logger,
	}
}

func (s *facultyService) Create(dto dto.CreateFacultyDTO) (*dto.FacultyResponse, error) {
	university, err := s.universityService.Get(context.Background(), dto.UniversityID)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return nil, err
		default:
			s.logger.Error("Failed to fetch university",
				zap.String("university_id", dto.UniversityID.String()),
				zap.String("service", "Faculty"),
				zap.String("operation", "Create"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to create faculty")
		}
	}

	existing, err := s.facultyRepo.FindByUniversityAndShortCode(dto.UniversityID, dto.ShortCode)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		s.logger.Error("Failed to check faculty short code",
			zap.String("short_code", dto.ShortCode),
			zap.String("service", "Faculty"),
			zap.String("operation", "Create"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create faculty")
	}
	if existing != nil {
		return nil, errors.NewConflictError("faculty with this short code already exists")
	}

	faculty := &models.Faculty{
		UniversityID: university.ID,
		NameEn:       strings.TrimSpace(dto.NameEn),
		NameFa:       strings.TrimSpace(dto.NameFa),
		ShortCode:    strings.ToUpper(strings.TrimSpace(dto.ShortCode)),
		IsActive:     dto.IsActive,
	}

	created, err := s.facultyRepo.Create(faculty)
	if err != nil {
		s.logger.Error("Failed to create faculty",
			zap.String("service", "Faculty"),
			zap.String("operation", "Create"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create faculty")
	}

	return mapFacultyToResponse(created), nil
}

func (s *facultyService) Get(id uuid.UUID) (*dto.FacultyResponse, error) {
	faculty, err := s.facultyRepo.Find(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return nil, err
		default:
			s.logger.Error("Failed to fetch faculty",
				zap.String("id", id.String()),
				zap.String("service", "Faculty"),
				zap.String("operation", "Get"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to get faculty")
		}
	}
	return mapFacultyToResponse(faculty), nil
}

func (s *facultyService) GetAllByUniversity(universityID uuid.UUID) ([]*dto.FacultyResponse, error) {
	_, err := s.universityService.Get(context.Background(), universityID)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return nil, err
		default:
			s.logger.Error("Failed to fetch university",
				zap.String("university_id", universityID.String()),
				zap.String("service", "Faculty"),
				zap.String("operation", "GetByUniversity"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to get faculties")
		}
	}

	faculties, err := s.facultyRepo.FindAllByUniversityID(universityID)
	if err != nil {
		s.logger.Error("Failed to fetch faculties",
			zap.String("university_id", universityID.String()),
			zap.String("service", "Faculty"),
			zap.String("operation", "GetByUniversity"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get faculties")
	}

	return mapFacultiesToResponse(faculties), nil
}

func (s *facultyService) GetByUniversityAndShortCode(universityID uuid.UUID, shortCode string) (*dto.FacultyResponse, error) {
	faculty, err := s.facultyRepo.FindByUniversityAndShortCode(universityID, strings.ToUpper(strings.TrimSpace(shortCode)))
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return nil, err
		default:
			s.logger.Error("Failed to fetch faculty by short code",
				zap.String("university_id", universityID.String()),
				zap.String("short_code", shortCode),
				zap.String("service", "Faculty"),
				zap.String("operation", "GetByUniversityAndShortCode"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to get faculty")
		}
	}

	return mapFacultyToResponse(faculty), nil
}

func (s *facultyService) Update(id uuid.UUID, dto dto.UpdateFacultyDTO) (*dto.FacultyResponse, error) {
	faculty, err := s.facultyRepo.Find(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return nil, err
		default:
			s.logger.Error("Failed to fetch faculty",
				zap.String("id", id.String()),
				zap.String("service", "Faculty"),
				zap.String("operation", "Update"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to update faculty")
		}
	}

	existing, err := s.facultyRepo.FindByUniversityAndShortCode(dto.UniversityID, dto.ShortCode)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		s.logger.Error("Failed to check faculty short code",
			zap.String("short_code", dto.ShortCode),
			zap.String("service", "Faculty"),
			zap.String("operation", "Update"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to update faculty")
	}
	if existing != nil && existing.ID != id {
		return nil, errors.NewConflictError("faculty with this short code already exists")
	}

	faculty.UniversityID = dto.UniversityID
	faculty.NameEn = strings.TrimSpace(dto.NameEn)
	faculty.NameFa = strings.TrimSpace(dto.NameFa)
	faculty.ShortCode = strings.ToUpper(strings.TrimSpace(dto.ShortCode))
	faculty.IsActive = dto.IsActive

	updated, err := s.facultyRepo.Update(faculty)
	if err != nil {
		s.logger.Error("Failed to update faculty",
			zap.String("id", id.String()),
			zap.String("service", "Faculty"),
			zap.String("operation", "Update"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to update faculty")
	}

	return mapFacultyToResponse(updated), nil
}

func (s *facultyService) Delete(id uuid.UUID) error {
	_, err := s.facultyRepo.Find(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return err
		default:
			s.logger.Error("Failed to fetch faculty",
				zap.String("id", id.String()),
				zap.String("service", "Faculty"),
				zap.String("operation", "Delete"),
				zap.Error(err))
			return fmt.Errorf("failed to delete faculty")
		}
	}

	if err := s.facultyRepo.Delete(id); err != nil {
		s.logger.Error("Failed to delete faculty",
			zap.String("id", id.String()),
			zap.String("service", "Faculty"),
			zap.String("operation", "Delete"),
			zap.Error(err))
		return fmt.Errorf("failed to delete faculty")
	}

	return nil
}

func mapFacultyToResponse(faculty *models.Faculty) *dto.FacultyResponse {
	return &dto.FacultyResponse{
		ID:           faculty.ID,
		UniversityID: faculty.UniversityID,
		NameEn:       faculty.NameEn,
		NameFa:       faculty.NameFa,
		ShortCode:    faculty.ShortCode,
		IsActive:     faculty.IsActive,
		CreatedAt:    faculty.CreatedAt,
		UpdatedAt:    faculty.UpdatedAt,
	}
}

func mapFacultiesToResponse(faculties []*models.Faculty) []*dto.FacultyResponse {
	response := make([]*dto.FacultyResponse, len(faculties))
	for i, faculty := range faculties {
		response[i] = mapFacultyToResponse(faculty)
	}
	return response
}
