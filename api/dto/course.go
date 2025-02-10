package dto

import (
	"github.com/google/uuid"
	"time"
)

// Request DTOs
type CreateCourseDTO struct {
	UniversityID      uuid.UUID `json:"university_id" binding:"required,uuid4"`
	FacultyID         uuid.UUID `json:"faculty_id" binding:"required,uuid4"`
	SemesterID        uuid.UUID `json:"semester_id" binding:"required,uuid4"`
	ProfessorName     string    `json:"professor_name" binding:"required"`
	Code              string    `json:"code" binding:"required"`
	Name              string    `json:"name" binding:"required"`
	Weight            int       `json:"weight" binding:"required,min=1"`
	Capacity          int       `json:"capacity" binding:"min=0"`
	GenderRestriction string    `json:"gender" binding:"required,oneof=male female mixed"`
	Times             []string  `json:"times" binding:"required"`
	TimeExam          string    `json:"time_exam" binding:"required"`
	DateExam          string    `json:"date_exam" binding:"required"`
}

type UpdateCourseDTO struct {
	UniversityID      uuid.UUID `json:"university_id" binding:"required,uuid4"`
	FacultyID         uuid.UUID `json:"faculty_id" binding:"required,uuid4"`
	SemesterID        uuid.UUID `json:"semester_id" binding:"required,uuid4"`
	ProfessorName     string    `json:"professor_name" binding:"required"`
	Code              string    `json:"code" binding:"required"`
	Name              string    `json:"name" binding:"required"`
	Weight            int       `json:"weight" binding:"required,min=1"`
	Capacity          int       `json:"capacity" binding:"min=0"`
	GenderRestriction string    `json:"gender" binding:"required,oneof=male female mixed"`
	Times             []string  `json:"times" binding:"required"`
	TimeExam          string    `json:"time_exam" binding:"required"`
	DateExam          string    `json:"date_exam" binding:"required"`
}

type CourseEngineDTO struct {
	UniversityID uuid.UUID `json:"university_id" binding:"required,uuid4"`
	SemesterID   uuid.UUID `json:"semester_id" binding:"required,uuid4"`
	CourseID     string    `json:"course_id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Weight       int       `json:"weight" binding:"required,min=1"`
	Capacity     int       `json:"capacity" binding:"min=0"`
	Gender       string    `json:"gender" binding:"required,oneof=male female mixed"`
	Professor    string    `json:"professor" binding:"required"`
	Faculty      string    `json:"faculty" binding:"required"`
	Time1        string    `json:"time1"`
	Time2        string    `json:"time2"`
	Time3        string    `json:"time3"`
	Time4        string    `json:"time4"`
	Time5        string    `json:"time5"`
	TimeExam     string    `json:"time_exam" binding:"required"`
	DateExam     string    `json:"date_exam" binding:"required"`
}

type BatchCreateCoursesDTO struct {
	Courses []CreateCourseDTO `json:"courses" binding:"required,min=1,dive"`
}

type BatchEngineCoursesDTO struct {
	UniversityID uuid.UUID         `json:"university_id" binding:"required,uuid4"`
	SemesterID   uuid.UUID         `json:"semester_id" binding:"required,uuid4"`
	Courses      []CourseEngineDTO `json:"courses" binding:"required,min=1,dive"`
}

// Response DTOs
type CourseTimeResponse struct {
	ID        uuid.UUID `json:"id"`
	CourseID  uuid.UUID `json:"course_id"`
	DayOfWeek int       `json:"day_of_week"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type CourseResponse struct {
	ID                uuid.UUID            `json:"id"`
	UniversityID      uuid.UUID            `json:"university_id"`
	FacultyID         uuid.UUID            `json:"faculty_id"`
	FacultyNameEn     string               `json:"faculty_name_en"`
	FacultyNameFa     string               `json:"faculty_name_fa"`
	ProfessorID       uuid.UUID            `json:"professor_id"`
	ProfessorName     string               `json:"professor_name"`
	SemesterID        uuid.UUID            `json:"semester_id"`
	Code              string               `json:"code"`
	Name              string               `json:"name"`
	Weight            int                  `json:"weight"`
	Capacity          int                  `json:"capacity"`
	GenderRestriction string               `json:"gender_restriction"`
	ExamStart         time.Time            `json:"exam_start"`
	ExamEnd           time.Time            `json:"exam_end"`
	CourseTimes       []CourseTimeResponse `json:"course_times"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}

// List Response
type CourseListResponse struct {
	Courses []*CourseResponse `json:"courses"`
	Total   int64             `json:"total"`
}

type CourseSearchFilters struct {
	FacultyID   uuid.UUID `form:"faculty_id"`
	ProfessorID uuid.UUID `form:"professor_id"`
	Query       string    `form:"q"`
}
