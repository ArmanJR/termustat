package repositories

import (
	"fmt"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfessorRepository interface {
	GetOrCreate(universityID uuid.UUID, rawName string) (uuid.UUID, error)
	FindAllByUniversity(universityID uuid.UUID) (*[]models.Professor, error)
	FindByID(id uuid.UUID) (*models.Professor, error)
	FindByNameAndUniversity(universityID uuid.UUID, name string) (*models.Professor, error)
	Create(professor *models.Professor) (*models.Professor, error)
}

type professorRepository struct {
	db *gorm.DB
}

func (r *professorRepository) FindByID(id uuid.UUID) (*models.Professor, error) {
	var professor models.Professor
	if err := r.db.First(&professor, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("professor", id.String())
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &professor, nil
}

func (r *professorRepository) FindByNameAndUniversity(universityID uuid.UUID, name string) (*models.Professor, error) {
	//TODO implement me
	panic("implement me")
}

func (r *professorRepository) Create(professor *models.Professor) (*models.Professor, error) {
	//TODO implement me
	panic("implement me")
}

func NewProfessorRepository(db *gorm.DB) ProfessorRepository {
	return &professorRepository{db: db}
}

func (r *professorRepository) FindAllByUniversity(universityID uuid.UUID) (*[]models.Professor, error) {
	var professors []models.Professor
	if err := r.db.Where("university_id = ?", universityID).Find(&professors).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &professors, nil
}

func (r *professorRepository) GetOrCreate(universityID uuid.UUID, rawName string) (uuid.UUID, error) {
	normalizedName := utils.NormalizeProfessor(rawName)
	if normalizedName == "" {
		return uuid.Nil, errors.New("invalid professor name after normalization")
	}

	var prof models.Professor
	err := r.db.Where(
		"university_id = ? AND normalized_name = ?",
		universityID,
		normalizedName,
	).First(&prof).Error

	if err == nil {
		return prof.ID, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newProf := models.Professor{
			UniversityID:   universityID,
			Name:           rawName,
			NormalizedName: normalizedName,
		}

		if err := r.db.Create(&newProf).Error; err != nil {
			return uuid.Nil, err
		}
		return newProf.ID, nil
	}

	return uuid.Nil, err
}
