package repositories

import (
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfessorRepository interface {
	FindByUniversityAndNormalizedName(universityID uuid.UUID, normalizedName string) (*models.Professor, error)
	FindAllByUniversity(universityID uuid.UUID) (*[]models.Professor, error)
	Create(professor *models.Professor) (*models.Professor, error)
	Find(id uuid.UUID) (*models.Professor, error)
}

type professorRepository struct {
	db *gorm.DB
}

func (r *professorRepository) Find(id uuid.UUID) (*models.Professor, error) {
	var professor models.Professor
	if err := r.db.First(&professor, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("professor", id.String())
		}
		return nil, errors.Wrap(err, "database error: failed to find professor")
	}
	return &professor, nil
}

func (r *professorRepository) FindByUniversityAndNormalizedName(universityID uuid.UUID, normalizedName string) (*models.Professor, error) {
	var professor models.Professor
	err := r.db.Where(
		"normalized_name = ? AND university_id = ?",
		normalizedName,
		universityID,
	).First(&professor).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, errors.NewNotFoundError("professor", normalizedName)
		default:
			return nil, errors.Wrap(err, "database error: failed to find professor")
		}
	}
	return &professor, nil
}

func (r *professorRepository) Create(professor *models.Professor) (*models.Professor, error) {
	if err := r.db.Create(professor).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create professor")
	}

	var created models.Professor
	if err := r.db.First(&created, professor.ID).Error; err != nil {
		return nil, errors.Wrap(err, "database error: failed to fetch created professor")
	}

	return &created, nil
}

func NewProfessorRepository(db *gorm.DB) ProfessorRepository {
	return &professorRepository{db: db}
}

func (r *professorRepository) FindAllByUniversity(universityID uuid.UUID) (*[]models.Professor, error) {
	var professors []models.Professor
	if err := r.db.Where("university_id = ?", universityID).Find(&professors).Error; err != nil {
		return nil, errors.Wrap(err, "database error: failed to find all professors")
	}
	return &professors, nil
}
