package services

import (
	"fmt"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type CourseService interface {
	Create(dto dto.CreateCourseDTO) (*dto.CourseResponse, error)
	Get(id uuid.UUID) (*dto.CourseResponse, error)
	GetAllBySemester(semesterID uuid.UUID) ([]*dto.CourseResponse, error)
	GetAllByFaculty(facultyID uuid.UUID) ([]*dto.CourseResponse, error)
	Update(id uuid.UUID, dto dto.UpdateCourseDTO) (*dto.CourseResponse, error)
	Delete(id uuid.UUID) error
	BatchCreate(dtos []dto.CreateCourseDTO) ([]*dto.CourseResponse, error)
	Search(filters *dto.CourseSearchFilters) ([]dto.CourseResponse, error)
}

type courseService struct {
	courseRepo        repositories.CourseRepository
	universityService UniversityService
	facultyService    FacultyService
	professorService  ProfessorService
	semesterService   SemesterService
	logger            *zap.Logger
}

func NewCourseService(
	courseRepo repositories.CourseRepository,
	universityService UniversityService,
	facultyService FacultyService,
	professorService ProfessorService,
	semesterService SemesterService,
	logger *zap.Logger,
) CourseService {
	return &courseService{
		courseRepo:        courseRepo,
		universityService: universityService,
		facultyService:    facultyService,
		professorService:  professorService,
		semesterService:   semesterService,
		logger:            logger,
	}
}

func (s *courseService) Create(dto dto.CreateCourseDTO) (*dto.CourseResponse, error) {
	if _, err := s.universityService.Get(dto.UniversityID); err != nil {
		return nil, err
	}

	if _, err := s.facultyService.Get(dto.FacultyID); err != nil {
		return nil, err
	}

	if _, err := s.semesterService.Get(dto.SemesterID); err != nil {
		return nil, err
	}

	professor, err := s.professorService.GetOrCreateByName(dto.UniversityID, dto.ProfessorName)
	if err != nil {
		s.logger.Error("Failed to get/create professor",
			zap.String("professor_name", dto.ProfessorName),
			zap.String("service", "Course"),
			zap.String("operation", "Create"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to process professor")
	}

	examStart, examEnd, err := s.parseExamDateTime(dto.DateExam, dto.TimeExam)
	if err != nil {
		return nil, errors.NewValidationError("invalid exam time: " + err.Error())
	}

	courseTimes, err := s.parseCourseTimes(dto.Times)
	if err != nil {
		return nil, errors.NewValidationError("invalid course time: " + err.Error())
	}

	// Check for existing course with same code
	existing, err := s.courseRepo.FindByUniversityAndCode(dto.UniversityID, dto.Code)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		s.logger.Error("Failed to check existing course",
			zap.String("code", dto.Code),
			zap.String("service", "Course"),
			zap.String("operation", "Create"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create course")
	}
	if existing != nil {
		return nil, errors.NewConflictError("course with this code already exists")
	}

	course := &models.Course{
		UniversityID:      dto.UniversityID,
		FacultyID:         dto.FacultyID,
		ProfessorID:       professor.ID,
		SemesterID:        dto.SemesterID,
		Code:              strings.TrimSpace(dto.Code),
		Name:              strings.TrimSpace(dto.Name),
		Weight:            dto.Weight,
		Capacity:          dto.Capacity,
		GenderRestriction: dto.GenderRestriction,
		ExamStart:         examStart,
		ExamEnd:           examEnd,
		CourseTimes:       courseTimes,
	}

	created, err := s.courseRepo.Create(course)
	if err != nil {
		s.logger.Error("Failed to create course",
			zap.String("service", "Course"),
			zap.String("operation", "Create"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create course")
	}

	return mapCourseToResponse(created), nil
}

func (s *courseService) Get(id uuid.UUID) (*dto.CourseResponse, error) {
	course, err := s.courseRepo.Find(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return nil, err
		default:
			s.logger.Error("Failed to fetch course",
				zap.String("id", id.String()),
				zap.String("service", "Course"),
				zap.String("operation", "Get"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to get course")
		}
	}
	return mapCourseToResponse(course), nil
}

func (s *courseService) GetAllBySemester(semesterID uuid.UUID) ([]*dto.CourseResponse, error) {
	if _, err := s.semesterService.Get(semesterID); err != nil {
		return nil, err
	}

	courses, err := s.courseRepo.FindAllBySemester(semesterID)
	if err != nil {
		s.logger.Error("Failed to fetch courses",
			zap.String("semester_id", semesterID.String()),
			zap.String("service", "Course"),
			zap.String("operation", "GetAllBySemester"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get courses")
	}

	return mapCoursesToResponse(courses), nil
}

func (s *courseService) GetAllByFaculty(facultyID uuid.UUID) ([]*dto.CourseResponse, error) {
	if _, err := s.facultyService.Get(facultyID); err != nil {
		return nil, err
	}

	courses, err := s.courseRepo.FindAllBySemester(facultyID)
	if err != nil {
		s.logger.Error("Failed to fetch courses",
			zap.String("faculty_id", facultyID.String()),
			zap.String("service", "Course"),
			zap.String("operation", "GetByFaculty"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get courses")
	}

	return mapCoursesToResponse(courses), nil
}

func (s *courseService) Update(id uuid.UUID, dto dto.UpdateCourseDTO) (*dto.CourseResponse, error) {
	existing, err := s.courseRepo.Find(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return nil, err
		default:
			s.logger.Error("Failed to fetch course for update",
				zap.String("id", id.String()),
				zap.String("service", "Course"),
				zap.String("operation", "Update"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to update course")
		}
	}

	if _, err := s.universityService.Get(dto.UniversityID); err != nil {
		return nil, err
	}

	if _, err := s.facultyService.Get(dto.FacultyID); err != nil {
		return nil, err
	}

	if _, err := s.semesterService.Get(dto.SemesterID); err != nil {
		return nil, err
	}

	professor, err := s.professorService.GetOrCreateByName(dto.UniversityID, dto.ProfessorName)
	if err != nil {
		s.logger.Error("Failed to get/create professor",
			zap.String("professor_name", dto.ProfessorName),
			zap.String("service", "Course"),
			zap.String("operation", "Create"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to process professor")
	}

	examStart, examEnd, err := s.parseExamDateTime(dto.DateExam, dto.TimeExam)
	if err != nil {
		return nil, errors.NewValidationError("invalid exam time: " + err.Error())
	}

	courseTimes, err := s.parseCourseTimes(dto.Times)
	if err != nil {
		return nil, errors.NewValidationError("invalid course time: " + err.Error())
	}

	existing.UniversityID = dto.UniversityID
	existing.FacultyID = dto.FacultyID
	existing.ProfessorID = professor.ID
	existing.Code = strings.TrimSpace(dto.Code)
	existing.Name = strings.TrimSpace(dto.Name)
	existing.Weight = dto.Weight
	existing.Capacity = dto.Capacity
	existing.GenderRestriction = dto.GenderRestriction
	existing.ExamStart = examStart
	existing.ExamEnd = examEnd
	existing.CourseTimes = courseTimes

	updated, err := s.courseRepo.Update(existing)
	if err != nil {
		s.logger.Error("Failed to update course",
			zap.String("id", id.String()),
			zap.String("service", "Course"),
			zap.String("operation", "Update"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to update course")
	}

	return mapCourseToResponse(updated), nil
}

func (s *courseService) Delete(id uuid.UUID) error {
	err := s.courseRepo.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFound):
			return err
		default:
			s.logger.Error("Failed to delete course",
				zap.String("id", id.String()),
				zap.String("service", "Course"),
				zap.String("operation", "Delete"),
				zap.Error(err))
			return fmt.Errorf("failed to delete course")
		}
	}
	return nil
}

func (s *courseService) BatchCreate(dtos []dto.CreateCourseDTO) ([]*dto.CourseResponse, error) {
	if len(dtos) == 0 {
		return nil, errors.NewValidationError("no courses provided")
	}

	// Pre-validate university and semester to avoid redundant checks
	universityID := dtos[0].UniversityID
	semesterID := dtos[0].SemesterID

	// Validate university exists
	if _, err := s.universityService.Get(universityID); err != nil {
		return nil, err
	}

	// Validate semester exists
	if _, err := s.semesterService.Get(semesterID); err != nil {
		return nil, err
	}

	// Validate all courses belong to the same university and semester
	for i, dto := range dtos {
		if dto.UniversityID != universityID {
			return nil, errors.NewValidationError(fmt.Sprintf("course at index %d has different university ID", i))
		}
		if dto.SemesterID != semesterID {
			return nil, errors.NewValidationError(fmt.Sprintf("course at index %d has different semester ID", i))
		}
	}

	// Track course codes to prevent duplicates within the batch
	courseCodes := make(map[string]bool)

	courses := make([]*models.Course, len(dtos))
	for i, dto := range dtos {
		// Check for duplicate course codes within the batch
		if _, exists := courseCodes[dto.Code]; exists {
			return nil, errors.NewValidationError(fmt.Sprintf("duplicate course code %s at index %d", dto.Code, i))
		}
		courseCodes[dto.Code] = true

		// Check if course code already exists in database
		existing, err := s.courseRepo.FindByUniversityAndCode(universityID, dto.Code)
		if err != nil && !errors.Is(err, errors.ErrNotFound) {
			s.logger.Error("Failed to check existing course",
				zap.String("code", dto.Code),
				zap.String("service", "Course"),
				zap.String("operation", "BatchCreate"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to validate course at index %d", i)
		}
		if existing != nil {
			return nil, errors.NewConflictError(fmt.Sprintf("course with code %s already exists", dto.Code))
		}

		// Prepare and validate each course
		course, err := s.prepareCourse(dto)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to prepare course at index %d", i))
		}
		courses[i] = course
	}

	// Begin transaction for batch creation
	created, err := s.courseRepo.BatchCreate(courses)
	if err != nil {
		s.logger.Error("Failed to batch create courses",
			zap.Int("count", len(courses)),
			zap.String("service", "Course"),
			zap.String("operation", "BatchCreate"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create courses")
	}

	s.logger.Info("Successfully created courses in batch",
		zap.Int("count", len(created)),
		zap.String("service", "Course"),
		zap.String("operation", "BatchCreate"))

	return mapCoursesToResponse(created), nil
}

func (s *courseService) Search(filters *dto.CourseSearchFilters) ([]dto.CourseResponse, error) {
	// Validate filters if necessary
	if err := s.validateSearchFilters(filters); err != nil {
		return nil, errors.Wrap(err, "invalid search filters")
	}

	// Perform search
	courses, err := s.courseRepo.Search(filters)
	if err != nil {
		s.logger.Error("Failed to search courses",
			zap.Any("filters", filters),
			zap.Error(err))
		return nil, fmt.Errorf("failed to search courses: %w", err)
	}

	// Map to response DTOs
	response := make([]dto.CourseResponse, 0, len(courses))
	for _, course := range courses {
		dto, err := s.mapCourseToDTO(&course)
		if err != nil {
			s.logger.Error("Failed to map course to DTO",
				zap.String("course_id", course.ID.String()),
				zap.Error(err))
			return nil, fmt.Errorf("failed to process course data: %w", err)
		}
		response = append(response, *dto)
	}

	return response, nil
}

func (s *courseService) mapCourseToDTO(course *models.Course) (*dto.CourseResponse, error) {
	// Get faculty details from faculty service
	faculty, err := s.facultyService.Get(course.FacultyID)
	if err != nil {
		s.logger.Error("Failed to fetch faculty details for course",
			zap.String("course_id", course.ID.String()),
			zap.String("faculty_id", course.FacultyID.String()),
			zap.Error(err))
		return nil, fmt.Errorf("failed to fetch faculty details: %w", err)
	}

	// Get professor details from professor service
	professor, err := s.professorService.Get(course.ProfessorID)
	if err != nil {
		s.logger.Error("Failed to fetch professor details for course",
			zap.String("course_id", course.ID.String()),
			zap.String("professor_id", course.ProfessorID.String()),
			zap.Error(err))
		return nil, fmt.Errorf("failed to fetch professor details: %w", err)
	}

	return &dto.CourseResponse{
		ID:                course.ID,
		UniversityID:      course.UniversityID,
		FacultyID:         course.FacultyID,
		FacultyNameEn:     faculty.NameEn,
		FacultyNameFa:     faculty.NameFa,
		ProfessorID:       course.ProfessorID,
		ProfessorName:     professor.Name,
		SemesterID:        course.SemesterID,
		Code:              course.Code,
		Name:              course.Name,
		Weight:            course.Weight,
		Capacity:          course.Capacity,
		GenderRestriction: course.GenderRestriction,
		ExamStart:         course.ExamStart,
		ExamEnd:           course.ExamEnd,
		CourseTimes:       nil,
		CreatedAt:         course.CreatedAt,
		UpdatedAt:         course.UpdatedAt,
	}, nil
}

// Helper methods

func (s *courseService) validateSearchFilters(filters *dto.CourseSearchFilters) error {
	if filters.FacultyID != uuid.Nil {
		_, err := s.facultyService.Get(filters.FacultyID)
		if err != nil {
			if errors.Is(err, errors.ErrNotFound) {
				return errors.NewValidationError("invalid faculty_id")
			}
			return fmt.Errorf("failed to validate faculty: %w", err)
		}
	}

	if filters.ProfessorID != uuid.Nil {
		_, err := s.professorService.Get(filters.ProfessorID)
		if err != nil {
			if errors.Is(err, errors.ErrNotFound) {
				return errors.NewValidationError("invalid professor_id")
			}
			return fmt.Errorf("failed to validate professor: %w", err)
		}
	}

	return nil
}

func (s *courseService) parseExamDateTime(dateStr, timeStr string) (time.Time, time.Time, error) {
	examDate, err := time.Parse("2006/01/02", dateStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.NewValidationError("invalid exam date format")
	}

	timeParts := strings.Split(timeStr, "-")
	if len(timeParts) != 2 {
		return time.Time{}, time.Time{}, errors.NewValidationError("invalid exam time format")
	}

	startTime, err := time.Parse("15:04", timeParts[0])
	if err != nil {
		return time.Time{}, time.Time{}, errors.NewValidationError("invalid exam start time")
	}

	endTime, err := time.Parse("15:04", timeParts[1])
	if err != nil {
		return time.Time{}, time.Time{}, errors.NewValidationError("invalid exam end time")
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

	if examEnd.Before(examStart) {
		return time.Time{}, time.Time{}, errors.NewValidationError("exam end time must be after start time")
	}

	return examStart, examEnd, nil
}

func (s *courseService) parseTimeSlot(timeStr string) (*models.CourseTime, error) {
	if timeStr == "" {
		return nil, nil
	}

	parts := strings.Split(timeStr, "/")
	if len(parts) != 2 {
		return nil, errors.NewValidationError("invalid time format")
	}

	dayPart := strings.TrimPrefix(parts[0], "d")
	day, err := strconv.Atoi(dayPart)
	if err != nil || day < 0 || day > 6 {
		return nil, errors.NewValidationError("invalid day value")
	}

	timeRange := strings.Split(parts[1], "-")
	if len(timeRange) != 2 {
		return nil, errors.NewValidationError("invalid time range")
	}

	startTime, err := time.Parse("15:04", timeRange[0])
	if err != nil {
		return nil, errors.NewValidationError("invalid start time")
	}

	endTime, err := time.Parse("15:04", timeRange[1])
	if err != nil {
		return nil, errors.NewValidationError("invalid end time")
	}

	if endTime.Before(startTime) {
		return nil, errors.NewValidationError("end time must be after start time")
	}

	return &models.CourseTime{
		DayOfWeek: day,
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

func (s *courseService) parseCourseTimes(times []string) ([]models.CourseTime, error) {
	var courseTimes []models.CourseTime
	for _, ts := range times {
		if ts == "" {
			continue
		}
		ct, err := s.parseTimeSlot(ts)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("invalid time slot: %s", ts))
		}
		if ct != nil {
			courseTimes = append(courseTimes, *ct)
		}
	}
	return courseTimes, nil
}

func (s *courseService) prepareCourse(dto dto.CreateCourseDTO) (*models.Course, error) {
	if _, err := s.universityService.Get(dto.UniversityID); err != nil {
		return nil, err
	}

	if _, err := s.facultyService.Get(dto.FacultyID); err != nil {
		return nil, err
	}

	if _, err := s.semesterService.Get(dto.SemesterID); err != nil {
		return nil, err
	}

	professor, err := s.professorService.GetOrCreateByName(dto.UniversityID, dto.ProfessorName)
	if err != nil {
		s.logger.Error("Failed to get/create professor",
			zap.String("professor_name", dto.ProfessorName),
			zap.String("service", "Course"),
			zap.String("operation", "prepareCourse"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to process professor")
	}

	examStart, examEnd, err := s.parseExamDateTime(dto.DateExam, dto.TimeExam)
	if err != nil {
		return nil, err
	}

	courseTimes, err := s.parseCourseTimes(dto.Times)
	if err != nil {
		return nil, err
	}

	// Check for existing course with same code
	existing, err := s.courseRepo.FindByUniversityAndCode(dto.UniversityID, dto.Code)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		s.logger.Error("Failed to check existing course",
			zap.String("code", dto.Code),
			zap.String("service", "Course"),
			zap.String("operation", "prepareCourse"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to prepare course")
	}
	if existing != nil {
		return nil, errors.NewConflictError("course with this code already exists")
	}

	return &models.Course{
		UniversityID:      dto.UniversityID,
		FacultyID:         dto.FacultyID,
		ProfessorID:       professor.ID,
		SemesterID:        dto.SemesterID,
		Code:              strings.TrimSpace(dto.Code),
		Name:              strings.TrimSpace(dto.Name),
		Weight:            dto.Weight,
		Capacity:          dto.Capacity,
		GenderRestriction: dto.GenderRestriction,
		ExamStart:         examStart,
		ExamEnd:           examEnd,
		CourseTimes:       courseTimes,
	}, nil
}

func (s *courseService) CreateFromEngine(reqDto dto.CourseEngineDTO) (*dto.CourseResponse, error) {
	// Convert engine times to standard format
	var times []string
	for _, t := range []string{reqDto.Time1, reqDto.Time2, reqDto.Time3, reqDto.Time4, reqDto.Time5} {
		if t != "" {
			times = append(times, t)
		}
	}

	createDTO := dto.CreateCourseDTO{
		UniversityID:      reqDto.UniversityID,
		SemesterID:        reqDto.SemesterID,
		Code:              reqDto.CourseID,
		Name:              reqDto.Name,
		ProfessorName:     reqDto.Professor,
		Weight:            reqDto.Weight,
		Capacity:          reqDto.Capacity,
		GenderRestriction: reqDto.Gender,
		Times:             times,
		TimeExam:          reqDto.TimeExam,
		DateExam:          reqDto.DateExam,
	}

	return s.Create(createDTO)
}

func mapCourseTimeToResponse(courseTime models.CourseTime) dto.CourseTimeResponse {
	return dto.CourseTimeResponse{
		ID:        courseTime.ID,
		CourseID:  courseTime.CourseID,
		DayOfWeek: courseTime.DayOfWeek,
		StartTime: courseTime.StartTime,
		EndTime:   courseTime.EndTime,
	}
}

func mapCourseTimesToResponse(courseTimes []models.CourseTime) []dto.CourseTimeResponse {
	if courseTimes == nil {
		return nil
	}
	response := make([]dto.CourseTimeResponse, len(courseTimes))
	for i, ct := range courseTimes {
		response[i] = mapCourseTimeToResponse(ct)
	}
	return response
}

func mapCourseToResponse(course *models.Course) *dto.CourseResponse {
	if course == nil {
		return nil
	}
	return &dto.CourseResponse{
		ID:                course.ID,
		UniversityID:      course.UniversityID,
		FacultyID:         course.FacultyID,
		ProfessorID:       course.ProfessorID,
		SemesterID:        course.SemesterID,
		Code:              course.Code,
		Name:              course.Name,
		Weight:            course.Weight,
		Capacity:          course.Capacity,
		GenderRestriction: course.GenderRestriction,
		ExamStart:         course.ExamStart,
		ExamEnd:           course.ExamEnd,
		CourseTimes:       mapCourseTimesToResponse(course.CourseTimes),
		CreatedAt:         course.CreatedAt,
		UpdatedAt:         course.UpdatedAt,
	}
}

func mapCoursesToResponse(courses []*models.Course) []*dto.CourseResponse {
	if courses == nil {
		return nil
	}
	responses := make([]*dto.CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = mapCourseToResponse(course)
	}
	return responses
}
