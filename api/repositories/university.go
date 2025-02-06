package repositories

import (
	"github.com/armanjr/termustat/api/models"
	"gorm.io/gorm"
)

type UniversityRepository interface {
	Create(university *models.University) error
	GetByID(id string) (*models.University, error)
	GetAll() ([]models.University, error)
	Update(university *models.University) error
	Delete(id string) error
}

type universityRepository struct {
	db *gorm.DB
}

func NewUniversityRepository(db *gorm.DB) UniversityRepository {
	return &universityRepository{db: db}
}

func (r *universityRepository) Create(university *models.University) error {
	return r.db.Create(university).Error
}

func (r *universityRepository) GetByID(id string) (*models.University, error) {
	var university models.University
	err := r.db.Preload("Faculties").First(&university, "id = ?", id).Error
	return &university, err
}

func (r *universityRepository) GetAll() ([]models.University, error) {
	var universities []models.University
	err := r.db.Find(&universities).Error
	return universities, err
}

func (r *universityRepository) Update(university *models.University) error {
	return r.db.Save(university).Error
}

func (r *universityRepository) Delete(id string) error {
	return r.db.Delete(&models.University{}, "id = ?", id).Error
}
