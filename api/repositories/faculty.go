package repositories

import (
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FacultyRepository interface {
	Create(faculty *models.Faculty) (*models.Faculty, error)
	Find(id uuid.UUID) (*models.Faculty, error)
	FindAllByUniversityID(universityID uuid.UUID) ([]*models.Faculty, error)
	FindByUniversityAndShortCode(universityID uuid.UUID, shortCode string) (*models.Faculty, error)
	Update(faculty *models.Faculty) (*models.Faculty, error)
	Delete(id uuid.UUID) error
}

type facultyRepository struct {
	db *gorm.DB
}

func NewFacultyRepository(db *gorm.DB) FacultyRepository {
	return &facultyRepository{db: db}
}

func (r *facultyRepository) Create(faculty *models.Faculty) (*models.Faculty, error) {
	if err := r.db.Create(faculty).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create faculty")
	}

	var created models.Faculty
	if err := r.db.First(&created, faculty.ID).Error; err != nil {
		return nil, errors.Wrap(err, "database error: failed to fetch created faculty")
	}

	return &created, nil
}

func (r *facultyRepository) Find(id uuid.UUID) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.db.First(&faculty, "id = ?", id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, errors.NewNotFoundError("faculty", id.String())
		default:
			return nil, errors.Wrap(err, "database error: failed to find faculty")
		}
	}
	return &faculty, nil
}

func (r *facultyRepository) FindAllByUniversityID(universityID uuid.UUID) ([]*models.Faculty, error) {
	var faculties []*models.Faculty
	if err := r.db.Where("university_id = ?", universityID).Find(&faculties).Error; err != nil {
		return nil, errors.Wrap(err, "database error: failed to find faculties")
	}
	return faculties, nil
}

func (r *facultyRepository) FindByUniversityAndShortCode(universityID uuid.UUID, shortCode string) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.db.Where(
		"university_id = ? AND short_code = ?",
		universityID,
		shortCode,
	).First(&faculty).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, errors.NewNotFoundError("faculty", shortCode)
		default:
			return nil, errors.Wrap(err, "database error: failed to find faculty by short code")
		}
	}
	return &faculty, nil
}

func (r *facultyRepository) Update(faculty *models.Faculty) (*models.Faculty, error) {
	if err := r.db.Save(faculty).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update faculty")
	}

	var updated models.Faculty
	if err := r.db.First(&updated, faculty.ID).Error; err != nil {
		return nil, errors.Wrap(err, "database error: failed to fetch updated faculty")
	}

	return &updated, nil
}

func (r *facultyRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Faculty{}, "id = ?", id)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete faculty")
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("faculty", id.String())
	}
	return nil
}
