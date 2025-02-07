package services

import (
	"fmt"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/repositories"
	"github.com/armanjr/termustat/api/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProfessorService interface {
	GetOrCreateByName(universityID uuid.UUID, name string) (*dto.ProfessorResponse, error)
	GetAllByUniversity(universityID uuid.UUID) (*[]dto.ProfessorResponse, error)
	Get(id uuid.UUID) (*dto.ProfessorResponse, error)
}

type professorService struct {
	professorRepository repositories.ProfessorRepository
	universityService   UniversityService
	logger              *zap.Logger
}

func NewProfessorService(
	professorRepository repositories.ProfessorRepository,
	universityService UniversityService,
	logger *zap.Logger) ProfessorService {
	return &professorService{
		professorRepository: professorRepository,
		universityService:   universityService,
		logger:              logger,
	}
}

func mapProfessorToDTO(professor *models.Professor) dto.ProfessorResponse {
	return dto.ProfessorResponse{
		ID:             professor.ID,
		Name:           professor.Name,
		NormalizedName: professor.NormalizedName,
		UniversityID:   professor.UniversityID,
		CreatedAt:      professor.CreatedAt,
		UpdatedAt:      professor.UpdatedAt,
	}
}

func (s *professorService) GetAllByUniversity(universityID uuid.UUID) (*[]dto.ProfessorResponse, error) {
	professors, err := s.professorRepository.FindAllByUniversity(universityID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch professors: %w", err)
	}
	if len(*professors) == 0 {
		return &[]dto.ProfessorResponse{}, nil
	}
	response := make([]dto.ProfessorResponse, len(*professors))
	for i, prof := range *professors {
		response[i] = mapProfessorToDTO(&prof)
	}
	return &response, nil
}

func (s *professorService) Get(id uuid.UUID) (*dto.ProfessorResponse, error) {
	professor, err := s.professorRepository.Find(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return nil, err
		default:
			s.logger.Error("Failed to fetch professor",
				zap.String("id", id.String()),
				zap.String("service", "Professor"),
				zap.String("operation", "Get"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to get professor")
		}
	}
	response := dto.ProfessorResponse{
		ID:             professor.ID,
		Name:           professor.Name,
		NormalizedName: professor.NormalizedName,
		UniversityID:   professor.UniversityID,
		CreatedAt:      professor.CreatedAt,
		UpdatedAt:      professor.UpdatedAt,
	}
	return &response, nil
}

func (s *professorService) GetOrCreateByName(universityID uuid.UUID, name string) (*dto.ProfessorResponse, error) {
	university, err := s.universityService.Get(universityID)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return nil, err
		default:
			s.logger.Error("Failed to fetch university",
				zap.String("id", universityID.String()),
				zap.String("service", "Professor"),
				zap.String("failed_service", "University"),
				zap.String("operation", "GetOrCreateByName"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to get professor")
		}
	}

	normalizedName := utils.NormalizeProfessor(name)
	if normalizedName == "" {
		s.logger.Error("Invalid professor name after normalization",
			zap.String("name", name),
			zap.String("service", "Professor"),
			zap.String("operation", "GetOrCreateByName"),
			zap.Error(err))
		return nil, fmt.Errorf("invalid professor name after normalization")
	}

	professor, err := s.professorRepository.FindByUniversityAndNormalizedName(universityID, normalizedName)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			newProfessor := models.Professor{
				UniversityID:   university.ID,
				Name:           name,
				NormalizedName: normalizedName,
			}
			professor, err = s.professorRepository.Create(&newProfessor)
			if err != nil {
				s.logger.Error("Failed to create professor",
					zap.String("name", name),
					zap.String("university_id", university.ID.String()),
					zap.String("service", "Professor"),
					zap.String("operation", "GetOrCreateByName"),
					zap.Error(err))
				return nil, fmt.Errorf("failed to create professor: %w", err)
			}
		default:
			s.logger.Error("Failed to fetch professor by university and normalized name",
				zap.String("id", university.ID.String()),
				zap.String("normalized_name", normalizedName),
				zap.String("service", "Professor"),
				zap.String("operation", "GetOrCreateByName"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to find professor")
		}
	}

	response := dto.ProfessorResponse{
		ID:             professor.ID,
		Name:           professor.Name,
		NormalizedName: professor.NormalizedName,
		UniversityID:   professor.UniversityID,
		CreatedAt:      professor.CreatedAt,
		UpdatedAt:      professor.UpdatedAt,
	}

	return &response, nil
}
