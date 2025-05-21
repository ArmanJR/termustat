package repositories

import (
	"context"
	"fmt"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UniversityRepository interface {
	Create(ctx context.Context, university *models.University) (*models.University, error)
	Update(ctx context.Context, university *models.University) (*models.University, error)
	Find(ctx context.Context, id uuid.UUID) (*models.University, error)
	FindAll(ctx context.Context) ([]models.University, error)
	ExistsByName(ctx context.Context, nameEn, nameFa string) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type universityRepository struct {
	db *gorm.DB
}

func NewUniversityRepository(db *gorm.DB) UniversityRepository {
	return &universityRepository{db: db}
}

func (r *universityRepository) Create(ctx context.Context, university *models.University) (*models.University, error) {
	if err := r.db.WithContext(ctx).Create(university).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create university")
	}

	var created models.University
	if err := r.db.WithContext(ctx).First(&created, university.ID).Error; err != nil {
		return nil, errors.Wrap(err, "database error: failed to fetch created university")
	}

	return &created, nil
}

func (r *universityRepository) Find(ctx context.Context, id uuid.UUID) (*models.University, error) {
	var university models.University

	if err := r.db.WithContext(ctx).First(&university, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("university", id.String())
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &university, nil
}

func (r *universityRepository) FindAll(ctx context.Context) ([]models.University, error) {
	var universities []models.University

	if err := r.db.WithContext(ctx).Find(&universities).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch universities")
	}

	return universities, nil
}

func (r *universityRepository) Update(ctx context.Context, university *models.University) (*models.University, error) {
	if err := r.db.WithContext(ctx).First(&models.University{}, university.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("university", university.ID.String())
		}
		return nil, errors.Wrap(err, "database error")
	}

	if err := r.db.WithContext(ctx).Save(university).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update university")
	}

	var updated models.University
	if err := r.db.WithContext(ctx).First(&updated, university.ID).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch updated university")
	}

	return &updated, nil
}

func (r *universityRepository) ExistsByName(ctx context.Context, nameEn, nameFa string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.University{}).
		Where("name_en = ? OR name_fa = ?", nameEn, nameFa).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check university existence: %w", err)
	}
	return count > 0, nil
}

func (r *universityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.University{}, "id = ?", id)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete university")
	}

	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("university", id.String())
	}

	return nil
}
