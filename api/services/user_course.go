package services

import (
	"context"
	"fmt"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserCourseService interface {
	AddCourse(userID, courseID, semesterID uuid.UUID) error
	RemoveCourse(userID, courseID uuid.UUID) error
	GetUserCourses(userID uuid.UUID, semesterID uuid.UUID) ([]dto.CourseResponse, error)
	ValidateTimeConflicts(userID, semesterID uuid.UUID, courseID uuid.UUID) error
	ValidateGenderRestriction(userID uuid.UUID, courseID uuid.UUID) error
	ValidateCapacity(courseID uuid.UUID) error
}

type userCourseService struct {
	userCourseRepo  repositories.UserCourseRepository
	courseService   CourseService
	userService     AdminUserService
	semesterService SemesterService
	logger          *zap.Logger
}

func NewUserCourseService(
	userCourseRepo repositories.UserCourseRepository,
	courseService CourseService,
	userService AdminUserService,
	semesterService SemesterService,
	logger *zap.Logger,
) UserCourseService {
	return &userCourseService{
		userCourseRepo:  userCourseRepo,
		courseService:   courseService,
		userService:     userService,
		semesterService: semesterService,
		logger:          logger,
	}
}

func (s *userCourseService) AddCourse(userID, courseID, semesterID uuid.UUID) error {
	// Check if course exists
	course, err := s.courseService.Get(courseID)
	if err != nil {
		return errors.Wrap(err, "failed to find course")
	}

	// Validate semester
	_, err = s.semesterService.Get(semesterID)
	if err != nil {
		return errors.Wrap(err, "failed to find semester")
	}

	// Check if already enrolled
	exists, err := s.userCourseRepo.ExistsByCourseAndSemester(userID, courseID, semesterID)
	if err != nil {
		return errors.Wrap(err, "failed to check enrollment")
	}
	if exists {
		return errors.NewConflictError("already enrolled in this course")
	}

	// Validate capacity
	if err := s.ValidateCapacity(courseID); err != nil {
		return err
	}

	// Validate gender restriction
	if err := s.ValidateGenderRestriction(userID, courseID); err != nil {
		return err
	}

	// Validate time conflicts
	if err := s.ValidateTimeConflicts(userID, semesterID, courseID); err != nil {
		return err
	}

	// Create enrollment
	userCourse := &models.UserCourse{
		UserID:     userID,
		CourseID:   courseID,
		SemesterID: semesterID,
	}

	if err := s.userCourseRepo.Create(userCourse); err != nil {
		return errors.Wrap(err, "failed to create enrollment")
	}

	s.logger.Info("Course added successfully",
		zap.String("user_id", userID.String()),
		zap.String("course_id", courseID.String()),
		zap.String("course_name", course.Name))

	return nil
}

func (s *userCourseService) RemoveCourse(userID, courseID uuid.UUID) error {
	if err := s.userCourseRepo.Delete(userID, courseID); err != nil {
		return errors.Wrap(err, "failed to remove course")
	}

	s.logger.Info("Course removed successfully",
		zap.String("user_id", userID.String()),
		zap.String("course_id", courseID.String()))

	return nil
}

func (s *userCourseService) GetUserCourses(userID uuid.UUID, semesterID uuid.UUID) ([]dto.CourseResponse, error) {
	userCourses, err := s.userCourseRepo.FindByUserAndSemester(userID, semesterID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user courses")
	}

	responses := make([]dto.CourseResponse, 0, len(userCourses))
	for _, uc := range userCourses {
		course, err := s.courseService.Get(uc.CourseID)
		if err != nil {
			s.logger.Error("Failed to fetch course details",
				zap.String("course_id", uc.CourseID.String()),
				zap.Error(err))
			continue
		}
		responses = append(responses, *course)
	}

	return responses, nil
}

func (s *userCourseService) ValidateTimeConflicts(userID, semesterID uuid.UUID, courseID uuid.UUID) error {
	// Get user's current courses
	userCourses, err := s.userCourseRepo.FindByUserAndSemester(userID, semesterID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch user courses")
	}

	// Get new course
	newCourse, err := s.courseService.Get(courseID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch course")
	}

	// Check for time conflicts
	for _, uc := range userCourses {
		course, err := s.courseService.Get(uc.CourseID)
		if err != nil {
			continue
		}

		if hasTimeConflict(course.CourseTimes, newCourse.CourseTimes) {
			return fmt.Errorf("time conflict with course: %s", course.Name)
		}
	}

	return nil
}

func (s *userCourseService) ValidateGenderRestriction(userID uuid.UUID, courseID uuid.UUID) error {
	ctx := context.Background() // todo: remove and pass request context
	user, err := s.userService.Get(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch user")
	}

	course, err := s.courseService.Get(courseID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch course")
	}

	if course.GenderRestriction != "mixed" && course.GenderRestriction != user.Gender {
		return fmt.Errorf("course is restricted to %s students", course.GenderRestriction)
	}

	return nil
}

func (s *userCourseService) ValidateCapacity(courseID uuid.UUID) error {
	course, err := s.courseService.Get(courseID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch course")
	}

	if course.Capacity <= 0 {
		return errors.New("course is full")
	}

	enrollments, err := s.userCourseRepo.FindByCourseAndSemester(courseID, course.SemesterID)
	if err != nil {
		return errors.Wrap(err, "failed to check course capacity")
	}

	if len(enrollments) >= course.Capacity {
		return errors.New("course is full")
	}

	return nil
}

func hasTimeConflict(times1, times2 []dto.CourseTimeResponse) bool {
	for _, t1 := range times1 {
		for _, t2 := range times2 {
			if t1.DayOfWeek == t2.DayOfWeek {
				if (t1.StartTime.Before(t2.EndTime) && t1.EndTime.After(t2.StartTime)) ||
					(t2.StartTime.Before(t1.EndTime) && t2.EndTime.After(t1.StartTime)) {
					return true
				}
			}
		}
	}
	return false
}
