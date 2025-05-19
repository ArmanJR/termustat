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
	GetOrCreateByName(universityID uuid.UUID, name string) (*dto.ProfessorMinimalResponse, error)
	GetAllByUniversity(universityID uuid.UUID) ([]dto.ProfessorMinimalResponse, error)
	Get(id uuid.UUID) (*dto.ProfessorDetailResponse, error)
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

func mapProfessorToListDTO(professor *models.Professor) dto.ProfessorMinimalResponse {
	return dto.ProfessorMinimalResponse{
		ID:             professor.ID,
		Name:           professor.Name,
		NormalizedName: professor.NormalizedName,
	}
}

func mapCoursesToProfessorResponse(courses []models.Course) []dto.CourseResponse {
	if len(courses) == 0 {
		return []dto.CourseResponse{}
	}

	response := make([]dto.CourseResponse, len(courses))
	for i, course := range courses {
		response[i] = dto.CourseResponse{
			ID:                course.ID,
			Code:              course.Code,
			Name:              course.Name,
			Weight:            course.Weight,
			Capacity:          course.Capacity,
			GenderRestriction: course.GenderRestriction,
			ExamStart:         course.ExamStart,
			ExamEnd:           course.ExamEnd,
		}
	}
	return response
}

func (s *professorService) GetAllByUniversity(universityID uuid.UUID) ([]dto.ProfessorMinimalResponse, error) {
	professors, err := s.professorRepository.FindAllByUniversity(universityID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch professors: %w", err)
	}

	if len(*professors) == 0 {
		return []dto.ProfessorMinimalResponse{}, nil
	}

	response := make([]dto.ProfessorMinimalResponse, len(*professors))
	for i, prof := range *professors {
		response[i] = mapProfessorToListDTO(&prof)
	}
	return response, nil
}

func (s *professorService) Get(id uuid.UUID) (*dto.ProfessorDetailResponse, error) {
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

	// Get university details
	university, err := s.universityService.Get(professor.UniversityID)
	if err != nil {
		s.logger.Error("Failed to fetch university for professor",
			zap.String("professor_id", id.String()),
			zap.String("university_id", professor.UniversityID.String()),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get professor details")
	}

	response := &dto.ProfessorDetailResponse{
		ID:             professor.ID,
		Name:           professor.Name,
		NormalizedName: professor.NormalizedName,
		University:     *university,
		Courses:        mapCoursesToProfessorResponse(professor.Courses),
		CreatedAt:      professor.CreatedAt,
		UpdatedAt:      professor.UpdatedAt,
	}
	return response, nil
}

func (s *professorService) GetOrCreateByName(universityID uuid.UUID, name string) (*dto.ProfessorMinimalResponse, error) {
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
			zap.String("operation", "GetOrCreateByName"))
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

	response := mapProfessorToListDTO(professor)
	return &response, nil
}
