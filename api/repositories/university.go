package repositories

import (
	"fmt"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UniversityRepository interface {
	Create(university *models.University) (*models.University, error)
	Update(university *models.University) (*models.University, error)
	Find(id uuid.UUID) (*models.University, error)
	FindAll() ([]models.University, error)
	Delete(id uuid.UUID) error
}

type universityRepository struct {
	db *gorm.DB
}

func NewUniversityRepository(db *gorm.DB) UniversityRepository {
	return &universityRepository{db: db}
}

func (r *universityRepository) Create(university *models.University) (*models.University, error) {
	if err := r.db.Create(university).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create university")
	}

	var created models.University
	if err := r.db.First(&created, university.ID).Error; err != nil {
		return nil, errors.Wrap(err, "database error: failed to fetch created university")
	}

	return &created, nil
}

func (r *universityRepository) Find(id uuid.UUID) (*models.University, error) {
	var university models.University

	if err := r.db.First(&university, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("university", id.String())
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &university, nil
}

func (r *universityRepository) FindAll() ([]models.University, error) {
	var universities []models.University

	if err := r.db.Find(&universities).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch universities")
	}

	return universities, nil
}

func (r *universityRepository) Update(university *models.University) (*models.University, error) {
	if err := r.db.First(&models.University{}, university.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("university", university.ID.String())
		}
		return nil, errors.Wrap(err, "database error")
	}

	if err := r.db.Save(university).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update university")
	}

	var updated models.University
	if err := r.db.First(&updated, university.ID).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch updated university")
	}

	return &updated, nil
}

func (r *universityRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.University{}, "id = ?", id)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete university")
	}

	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("university", id.String())
	}

	return nil
}
