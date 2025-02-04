package handlers

import (
	"fmt"
	"github.com/armanjr/termustat/app/logger"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/armanjr/termustat/app/config"
	"github.com/armanjr/termustat/app/models"
	courseRepo "github.com/armanjr/termustat/app/repositories/course"
	professorRepo "github.com/armanjr/termustat/app/repositories/professor"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CourseRequest represents course creation/update request
type CourseRequest struct {
	UniversityID string `json:"university_id" binding:"required,uuid4"`
	FacultyID    string `json:"faculty_id" binding:"required,uuid4"`
	SemesterID   string `json:"semester_id" binding:"required,uuid4"`

	ProfessorName     string   `json:"professor_name" binding:"required"`
	Code              string   `json:"code" binding:"required"`
	Name              string   `json:"name" binding:"required"`
	Weight            int      `json:"weight" binding:"required,min=1"`
	Capacity          int      `json:"capacity" binding:"min=0"`
	GenderRestriction string   `json:"gender" binding:"required,oneof=male female mixed"`
	Times             []string `json:"times" binding:"required"`
	TimeExam          string   `json:"time_exam" binding:"required"`
	DateExam          string   `json:"date_exam" binding:"required"`
}

type CourseEngineRequest struct {
	CourseID  string `json:"course_id"`
	Name      string `json:"name"`
	Weight    string `json:"weight"`
	Capacity  string `json:"capacity"`
	Gender    string `json:"gender"`
	Professor string `json:"professor"` // professor name
	Faculty   string `json:"faculty"`   // faculty code
	Time1     string `json:"time1"`
	Time2     string `json:"time2"`
	Time3     string `json:"time3"`
	Time4     string `json:"time4"`
	Time5     string `json:"time5"`
	TimeExam  string `json:"time_exam"`
	DateExam  string `json:"date_exam"`
}

// CoursesEngineRequest represents courses creation/update request from engine's processing
type CoursesEngineRequest struct {
	UniversityID string                `json:"university_id" binding:"required"`
	SemesterID   string                `json:"semester_id" binding:"required"`
	Courses      []CourseEngineRequest `json:"courses" binding:"required"`
}

// CourseProcessingError represents a custom error for course processing
type CourseProcessingError struct {
	StatusCode int
	Err        error
}

func (e *CourseProcessingError) Error() string {
	return e.Err.Error()
}

// CreateCourse is used to create course via admin dashboard
func CreateCourse(c *gin.Context) {
	var req CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid course request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := processCourse(req)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	if err := courseRepo.CreateCourse(config.DB, &course, course.CourseTimes); err != nil {
		logger.Log.Error("Failed to create course", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// CreateCourseFromEngine is used to process engine's response and call CreateCourse
func CreateCourseFromEngine(c *gin.Context) {
	var req CoursesEngineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid course request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//to be implemented

	c.JSON(http.StatusNotImplemented, nil)
}

// parseTimeSlot parses a time string into CourseTime
func parseTimeSlot(timeStr string) (*models.CourseTime, error) {
	if timeStr == "" {
		return nil, nil
	}

	parts := strings.Split(timeStr, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid time format")
	}

	dayPart := strings.TrimPrefix(parts[0], "d")
	day, err := strconv.Atoi(dayPart)
	if err != nil || day < 0 || day > 6 {
		return nil, fmt.Errorf("invalid day value")
	}

	timeRange := strings.Split(parts[1], "-")
	if len(timeRange) != 2 {
		return nil, fmt.Errorf("invalid time range")
	}

	startTime, err := time.Parse("15:04", timeRange[0])
	if err != nil {
		return nil, fmt.Errorf("invalid start time")
	}

	endTime, err := time.Parse("15:04", timeRange[1])
	if err != nil {
		return nil, fmt.Errorf("invalid end time")
	}

	return &models.CourseTime{
		DayOfWeek: day,
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// parseExamDateTime parses exam date and time
func parseExamDateTime(dateStr, timeStr string) (time.Time, time.Time, error) {
	examDate, err := time.Parse("2006/01/02", dateStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	timeParts := strings.Split(timeStr, "-")
	if len(timeParts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid exam time format")
	}

	startTime, err := time.Parse("15:04", timeParts[0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endTime, err := time.Parse("15:04", timeParts[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	examStart := time.Date(
		examDate.Year(),
		examDate.Month(),
		examDate.Day(),
		startTime.Hour(),
		startTime.Minute(),
		0, 0, time.UTC,
	)

	examEnd := time.Date(
		examDate.Year(),
		examDate.Month(),
		examDate.Day(),
		endTime.Hour(),
		endTime.Minute(),
		0, 0, time.UTC,
	)

	return examStart, examEnd, nil
}

func processCourse(courseRequest CourseRequest) (models.Course, *CourseProcessingError) {
	var course models.Course

	// Parse professor
	universityID := uuid.MustParse(courseRequest.UniversityID)
	professorID, err := professorRepo.GetOrCreate(universityID, courseRequest.ProfessorName)
	if err != nil {
		logger.Log.Error("Failed to get/create professor", zap.Error(err))
		return course, &CourseProcessingError{StatusCode: http.StatusInternalServerError, Err: fmt.Errorf("failed to process professor")}
	}

	// Parse exam times
	examStart, examEnd, err := parseExamDateTime(courseRequest.DateExam, courseRequest.TimeExam)
	if err != nil {
		logger.Log.Warn("Invalid exam time format", zap.Error(err))
		return course, &CourseProcessingError{StatusCode: http.StatusBadRequest, Err: fmt.Errorf("invalid exam time format")}
	}

	// Parse course times
	var courseTimes []models.CourseTime
	for _, ts := range courseRequest.Times {
		if ts == "" {
			continue
		}
		ct, err := parseTimeSlot(ts)
		if err != nil {
			logger.Log.Warn("Invalid time slot format", zap.String("slot", ts), zap.Error(err))
			return course, &CourseProcessingError{StatusCode: http.StatusBadRequest, Err: fmt.Errorf("invalid time slot: %s", ts)}
		}
		courseTimes = append(courseTimes, *ct)
	}

	course = models.Course{
		UniversityID:      universityID,
		FacultyID:         uuid.MustParse(courseRequest.FacultyID),
		ProfessorID:       professorID,
		SemesterID:        uuid.MustParse(courseRequest.SemesterID),
		Code:              courseRequest.Code,
		Name:              courseRequest.Name,
		Weight:            courseRequest.Weight,
		Capacity:          courseRequest.Capacity,
		GenderRestriction: courseRequest.GenderRestriction,
		ExamStart:         examStart,
		ExamEnd:           examEnd,
		CourseTimes:       courseTimes,
	}

	return course, nil
}

// BatchCourseRequest represents batch course creation
type BatchCourseRequest struct {
	Courses []CourseRequest `json:"courses" binding:"required,min=1"`
}

func BatchCreateCourses(c *gin.Context) {
	var req BatchCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid batch request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()
	defer tx.Rollback()

	var createdCourses []models.Course
	for _, courseReq := range req.Courses {
		course, err := processCourse(courseReq)
		if err != nil {
			c.JSON(err.StatusCode, gin.H{"error": err.Error()})
			return
		}

		if err := courseRepo.CreateCourse(config.DB, &course, course.CourseTimes); err != nil {
			logger.Log.Error("Failed to create course", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create course: %s", err)})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"created": len(createdCourses)})
}

// GetCourse returns a single course
func GetCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course

	if err := config.DB.Preload("CourseTimes").
		First(&course, "id = ?", id).Error; err != nil {
		logger.Log.Warn("Course not found", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// UpdateCourse updates course details
func UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	var req CourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid course update request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var course models.Course
	if err := config.DB.First(&course, "id = ?", id).Error; err != nil {
		logger.Log.Warn("Course not found", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	updatedCourse, err := processCourse(req)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	updatedCourse.ID = course.ID

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete old course times
	if err := tx.Where("course_id = ?", id).Delete(&models.CourseTime{}).Error; err != nil {
		logger.Log.Error("Failed to delete old course times", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course schedule"})
		tx.Rollback()
		return
	}

	// Update course details and insert new course times
	if err := courseRepo.CreateCourse(tx, &updatedCourse, updatedCourse.CourseTimes); err != nil {
		logger.Log.Error("Failed to update course", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		tx.Rollback()
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, updatedCourse)
}

// DeleteCourse deletes a course
func DeleteCourse(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Where("course_id = ?", id).Delete(&models.CourseTime{}).Error; err != nil {
		logger.Log.Error("Failed to delete course times", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course schedule"})
		return
	}

	if err := config.DB.Delete(&models.Course{}, "id = ?", id).Error; err != nil {
		logger.Log.Error("Failed to delete course", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}
