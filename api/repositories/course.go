package repositories

import (
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *models.Course) (*models.Course, error)
	Find(id uuid.UUID) (*models.Course, error)
	FindAllBySemester(semesterID uuid.UUID) ([]*models.Course, error)
	FindAllByFaculty(facultyID uuid.UUID) ([]*models.Course, error)
	FindAllByProfessor(professorID uuid.UUID) ([]*models.Course, error)
	FindByUniversityAndCode(universityID uuid.UUID, code string) (*models.Course, error)
	Update(course *models.Course) (*models.Course, error)
	Delete(id uuid.UUID) error
	BatchCreate(courses []*models.Course) ([]*models.Course, error)
	Search(filters *dto.CourseSearchFilters) ([]models.Course, error)
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(course *models.Course) (*models.Course, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(course).Error; err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to create course")
	}

	if len(course.CourseTimes) > 0 {
		if err := tx.Create(&course.CourseTimes).Error; err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, "failed to create course times")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}

	var created models.Course
	if err := r.db.Preload("CourseTimes").First(&created, course.ID).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch created course")
	}

	return &created, nil
}

func (r *courseRepository) Find(id uuid.UUID) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("CourseTimes").First(&course, "id = ?", id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, errors.NewNotFoundError("course", id.String())
		default:
			return nil, errors.Wrap(err, "failed to find course")
		}
	}
	return &course, nil
}

func (r *courseRepository) FindAllBySemester(semesterID uuid.UUID) ([]*models.Course, error) {
	var courses []*models.Course
	err := r.db.Preload("CourseTimes").
		Where("semester_id = ?", semesterID).
		Find(&courses).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find courses by semester")
	}
	return courses, nil
}

func (r *courseRepository) FindAllByFaculty(facultyID uuid.UUID) ([]*models.Course, error) {
	var courses []*models.Course
	err := r.db.Preload("CourseTimes").
		Where("faculty_id = ?", facultyID).
		Find(&courses).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find courses by faculty")
	}
	return courses, nil
}

func (r *courseRepository) FindAllByProfessor(professorID uuid.UUID) ([]*models.Course, error) {
	var courses []*models.Course
	err := r.db.Preload("CourseTimes").
		Where("professor_id = ?", professorID).
		Find(&courses).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find courses by professor")
	}
	return courses, nil
}

func (r *courseRepository) FindByUniversityAndCode(universityID uuid.UUID, code string) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("CourseTimes").
		Where("university_id = ? AND code = ?", universityID, code).
		First(&course).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, errors.NewNotFoundError("course", code)
		default:
			return nil, errors.Wrap(err, "failed to find course by code")
		}
	}
	return &course, nil
}

func (r *courseRepository) Update(course *models.Course) (*models.Course, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete existing course times
	if err := tx.Where("course_id = ?", course.ID).Delete(&models.CourseTime{}).Error; err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to delete existing course times")
	}

	// Update course
	if err := tx.Save(course).Error; err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to update course")
	}

	// Create new course times
	if len(course.CourseTimes) > 0 {
		if err := tx.Create(&course.CourseTimes).Error; err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, "failed to create new course times")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}

	var updated models.Course
	if err := r.db.Preload("CourseTimes").First(&updated, course.ID).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch updated course")
	}

	return &updated, nil
}

func (r *courseRepository) Delete(id uuid.UUID) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("course_id = ?", id).Delete(&models.CourseTime{}).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete course times")
	}

	result := tx.Delete(&models.Course{}, "id = ?", id)
	if result.Error != nil {
		tx.Rollback()
		return errors.Wrap(result.Error, "failed to delete course")
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.NewNotFoundError("course", id.String())
	}

	return tx.Commit().Error
}

func (r *courseRepository) BatchCreate(courses []*models.Course) ([]*models.Course, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, course := range courses {
		if err := tx.Create(course).Error; err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, "failed to create course in batch")
		}

		if len(course.CourseTimes) > 0 {
			if err := tx.Create(&course.CourseTimes).Error; err != nil {
				tx.Rollback()
				return nil, errors.Wrap(err, "failed to create course times in batch")
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.Wrap(err, "failed to commit batch creation")
	}

	var created []*models.Course
	if err := r.db.Preload("CourseTimes").Find(&created).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch created courses")
	}

	return created, nil
}

func (r *courseRepository) Search(filters *dto.CourseSearchFilters) ([]models.Course, error) {
	var courses []models.Course

	query := r.db.Model(&models.Course{})

	query = query.Joins("Faculty").Joins("Professor")

	if filters.FacultyID != uuid.Nil {
		query = query.Where("faculty_id = ?", filters.FacultyID)
	}

	if filters.ProfessorID != uuid.Nil {
		query = query.Where("professor_id = ?", filters.ProfessorID)
	}

	if filters.Query != "" {
		searchQuery := "%" + filters.Query + "%"
		query = query.Where("(LOWER(name) LIKE LOWER(?) OR LOWER(code) LIKE LOWER(?))",
			searchQuery, searchQuery)
	}

	if err := query.Find(&courses).Error; err != nil {
		return nil, errors.Wrap(err, "failed to search courses")
	}

	return courses, nil
}
