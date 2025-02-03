package course

import (
	"github.com/armanjr/termustat/app/models"
	"gorm.io/gorm"
)

// CreateCourse creates a course and its associated times in a transaction
func CreateCourse(db *gorm.DB, course *models.Course, times []models.CourseTime) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(course).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range times {
		times[i].CourseID = course.ID
		if err := tx.Create(&times[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
