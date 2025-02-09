package repositories

import (
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserCourseRepository interface {
	Create(userCourse *models.UserCourse) error
	Delete(userID, courseID uuid.UUID) error
	FindByUserAndSemester(userID, semesterID uuid.UUID) ([]models.UserCourse, error)
	FindByCourseAndSemester(courseID, semesterID uuid.UUID) ([]models.UserCourse, error)
	ExistsByCourseAndSemester(userID, courseID, semesterID uuid.UUID) (bool, error)
	GetCoursesForUser(userID uuid.UUID, pagination *dto.PaginationQuery) (*dto.PaginatedList[models.UserCourse], error)
}

type userCourseRepository struct {
	db *gorm.DB
}

func NewUserCourseRepository(db *gorm.DB) UserCourseRepository {
	return &userCourseRepository{db: db}
}

func (r *userCourseRepository) Create(userCourse *models.UserCourse) error {
	if err := r.db.Create(userCourse).Error; err != nil {
		return errors.Wrap(err, "failed to create user course")
	}
	return nil
}

func (r *userCourseRepository) Delete(userID, courseID uuid.UUID) error {
	result := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).
		Delete(&models.UserCourse{})

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete user course")
	}

	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("user course", "")
	}

	return nil
}

func (r *userCourseRepository) FindByUserAndSemester(userID, semesterID uuid.UUID) ([]models.UserCourse, error) {
	var userCourses []models.UserCourse

	err := r.db.Preload("Course").
		Preload("Course.CourseTimes").
		Where("user_id = ? AND semester_id = ?", userID, semesterID).
		Find(&userCourses).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to find user courses")
	}

	return userCourses, nil
}

func (r *userCourseRepository) FindByCourseAndSemester(courseID, semesterID uuid.UUID) ([]models.UserCourse, error) {
	var userCourses []models.UserCourse

	err := r.db.Preload("User").
		Where("course_id = ? AND semester_id = ?", courseID, semesterID).
		Find(&userCourses).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to find course enrollments")
	}

	return userCourses, nil
}

func (r *userCourseRepository) ExistsByCourseAndSemester(userID, courseID, semesterID uuid.UUID) (bool, error) {
	var count int64

	err := r.db.Model(&models.UserCourse{}).
		Where("user_id = ? AND course_id = ? AND semester_id = ?", userID, courseID, semesterID).
		Count(&count).Error

	if err != nil {
		return false, errors.Wrap(err, "failed to check user course existence")
	}

	return count > 0, nil
}

func (r *userCourseRepository) GetCoursesForUser(userID uuid.UUID, pagination *dto.PaginationQuery) (*dto.PaginatedList[models.UserCourse], error) {
	var userCourses []models.UserCourse
	var total int64

	query := r.db.Model(&models.UserCourse{}).Where("user_id = ?", userID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count user courses")
	}

	// Get paginated results with preloaded relations
	err := query.Preload("Course").
		Preload("Course.CourseTimes").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&userCourses).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user courses")
	}

	return &dto.PaginatedList[models.UserCourse]{
		Items: userCourses,
		Total: total,
		Page:  pagination.Page,
		Limit: pagination.Limit,
	}, nil
}
