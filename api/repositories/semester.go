package repositories

import (
	"fmt"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SemesterRepository interface {
	Find(id uuid.UUID) (*models.Semester, error)
	FindAll() ([]models.Semester, error)
	Create(semester *models.Semester) (*models.Semester, error)
	Update(semester *models.Semester) (*models.Semester, error)
	Delete(id uuid.UUID) error
	FindByYearAndTerm(year int, term string) (*models.Semester, error)
}

type semesterRepository struct {
	db *gorm.DB
}

func NewSemesterRepository(db *gorm.DB) SemesterRepository {
	return &semesterRepository{db: db}
}

func (r *semesterRepository) Find(id uuid.UUID) (*models.Semester, error) {
	var semester models.Semester
	if err := r.db.First(&semester, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("semester", id.String())
		}
		return nil, errors.Wrap(err, "database error")
	}
	return &semester, nil
}

func (r *semesterRepository) FindAll() ([]models.Semester, error) {
	var semesters []models.Semester
	if err := r.db.Order("year DESC, term DESC").Find(&semesters).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch semesters")
	}
	return semesters, nil
}

func (r *semesterRepository) Create(semester *models.Semester) (*models.Semester, error) {
	if err := r.db.Create(semester).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create semester")
	}

	var created models.Semester
	if err := r.db.First(&created, semester.ID).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch created semester")
	}

	return &created, nil
}

func (r *semesterRepository) Update(semester *models.Semester) (*models.Semester, error) {
	if err := r.db.First(&models.Semester{}, semester.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("semester", semester.ID.String())
		}
		return nil, errors.Wrap(err, "database error")
	}

	if err := r.db.Save(semester).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update semester")
	}

	var updated models.Semester
	if err := r.db.First(&updated, semester.ID).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch updated semester")
	}

	return &updated, nil
}

func (r *semesterRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Semester{}, "id = ?", id)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete semester")
	}

	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("semester", id.String())
	}

	return nil
}

func (r *semesterRepository) FindByYearAndTerm(year int, term string) (*models.Semester, error) {
	var semester models.Semester
	if err := r.db.Where("year = ? AND term = ?", year, term).First(&semester).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("semester",
				fmt.Sprintf("year: %d, term: %s", year, term))
		}
		return nil, errors.Wrap(err, "database error")
	}
	return &semester, nil
}
