package services

import (
	"fmt"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type ProfessorService interface {
	GetProfessorsByUniversity(universityID uuid.UUID) (*[]dto.ProfessorResponse, error)
	GetOrCreateProfessor(universityID uuid.UUID, request *dto.CreateProfessorRequest) (*dto.ProfessorResponse, error)
	GetProfessor(id uuid.UUID) (*dto.ProfessorResponse, error)
	UpdateProfessor(id uuid.UUID, req *dto.UpdateProfessorRequest) (*dto.ProfessorResponse, error)
}

type professorService struct {
	repo   repositories.ProfessorRepository
	logger *zap.Logger
}

func NewProfessorService(repo repositories.ProfessorRepository, logger *zap.Logger) ProfessorService {
	return &professorService{
		repo:   repo,
		logger: logger,
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

func (s *professorService) GetProfessorsByUniversity(universityID uuid.UUID) (*[]dto.ProfessorResponse, error) {
	professors, err := s.repo.FindAllByUniversity(universityID)
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

func (s *professorService) GetProfessor(id uuid.UUID) (*dto.ProfessorResponse, error) {
	professor, err := s.repo.FindByID(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return nil, err
		default:
			s.logger.Error("Failed to fetch professor",
				zap.String("id", id.String()),
				zap.String("operation", "GetProfessor"),
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
