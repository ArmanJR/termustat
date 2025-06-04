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
	"golang.org/x/crypto/bcrypt"
)

type AdminUserService interface {
	Create(ctx context.Context, req *dto.AdminCreateUserRequest) (*dto.AdminUserResponse, error)
	Get(ctx context.Context, id uuid.UUID) (*dto.AdminUserResponse, error)
	GetAll(ctx context.Context, pagination *dto.PaginationQuery) (*dto.PaginatedList[dto.AdminUserResponse], error)
	Update(ctx context.Context, id uuid.UUID, req *dto.AdminUpdateUserRequest) (*dto.AdminUserResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUniversity(ctx context.Context, universityID uuid.UUID, pagination *dto.PaginationQuery) (*dto.PaginatedList[dto.AdminUserResponse], error)
	GetByFaculty(ctx context.Context, facultyID uuid.UUID, pagination *dto.PaginationQuery) (*dto.PaginatedList[dto.AdminUserResponse], error)
	UpdatePassword(ctx context.Context, id uuid.UUID, req *dto.AdminUpdatePasswordRequest) error
	VerifyEmail(ctx context.Context, id uuid.UUID) error
}

type adminUserService struct {
	adminUserRepository repositories.AdminUserRepository
	universityService   UniversityService
	facultyService      FacultyService
	logger              *zap.Logger
}

func NewAdminUserService(
	adminUserRepository repositories.AdminUserRepository,
	universityService UniversityService,
	facultyService FacultyService,
	logger *zap.Logger,
) AdminUserService {
	return &adminUserService{
		adminUserRepository: adminUserRepository,
		universityService:   universityService,
		facultyService:      facultyService,
		logger:              logger,
	}
}

func (s *adminUserService) Create(ctx context.Context, req *dto.AdminCreateUserRequest) (*dto.AdminUserResponse, error) {
	if _, err := s.universityService.Get(ctx, req.UniversityID); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.NewValidationError("invalid university_id")
		}
		return nil, fmt.Errorf("failed to validate university: %w", err)
	}

	if _, err := s.facultyService.Get(req.FacultyID); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.NewValidationError("invalid faculty_id")
		}
		return nil, fmt.Errorf("failed to validate faculty: %w", err)
	}

	existing, err := s.adminUserRepository.FindByEmailOrStudentID(ctx, req.Email, req.StudentID)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, errors.NewConflictError("user with this email or student ID already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		StudentID:     req.StudentID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		UniversityID:  req.UniversityID,
		FacultyID:     req.FacultyID,
		Gender:        req.Gender,
		EmailVerified: false,
		IsAdmin:       false,
	}

	created, err := s.adminUserRepository.Create(ctx, user)
	if err != nil {
		s.logger.Error("Failed to create user",
			zap.String("email", req.Email),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.mapUserToDTO(created)
}

func (s *adminUserService) Get(ctx context.Context, id uuid.UUID) (*dto.AdminUserResponse, error) {
	user, err := s.adminUserRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}
		s.logger.Error("Failed to fetch user",
			zap.String("id", id.String()),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return s.mapUserToDTO(user)
}

func (s *adminUserService) GetAll(ctx context.Context, pagination *dto.PaginationQuery) (*dto.PaginatedList[dto.AdminUserResponse], error) {
	result, err := s.adminUserRepository.GetAll(ctx, pagination)
	if err != nil {
		s.logger.Error("Failed to fetch users", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	dtos := make([]dto.AdminUserResponse, 0, len(result.Items))
	for _, user := range result.Items {
		dto, err := s.mapUserToDTO(&user)
		if err != nil {
			return nil, fmt.Errorf("failed to map user data: %w", err)
		}
		dtos = append(dtos, *dto)
	}

	return &dto.PaginatedList[dto.AdminUserResponse]{
		Items: dtos,
		Total: result.Total,
		Page:  result.Page,
		Limit: result.Limit,
	}, nil
}

func (s *adminUserService) Update(ctx context.Context, id uuid.UUID, req *dto.AdminUpdateUserRequest) (*dto.AdminUserResponse, error) {
	user, err := s.adminUserRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.UniversityID != uuid.Nil && req.UniversityID != user.UniversityID {
		if _, err := s.universityService.Get(ctx, req.UniversityID); err != nil {
			if errors.Is(err, errors.ErrNotFound) {
				return nil, errors.NewValidationError("invalid university_id")
			}
			return nil, fmt.Errorf("failed to validate university: %w", err)
		}
		user.UniversityID = req.UniversityID
	}

	if req.FacultyID != uuid.Nil && req.FacultyID != user.FacultyID {
		if _, err := s.facultyService.Get(req.FacultyID); err != nil {
			if errors.Is(err, errors.ErrNotFound) {
				return nil, errors.NewValidationError("invalid faculty_id")
			}
			return nil, fmt.Errorf("failed to validate faculty: %w", err)
		}
		user.FacultyID = req.FacultyID
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}

	updated, err := s.adminUserRepository.Update(ctx, user)
	if err != nil {
		s.logger.Error("Failed to update user",
			zap.String("id", id.String()),
			zap.Error(err))
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.mapUserToDTO(updated)
}

func (s *adminUserService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.adminUserRepository.Delete(ctx, id); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return err
		}
		s.logger.Error("Failed to delete user",
			zap.String("id", id.String()),
			zap.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *adminUserService) UpdatePassword(ctx context.Context, id uuid.UUID, req *dto.AdminUpdatePasswordRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.adminUserRepository.UpdatePassword(ctx, id, string(hashedPassword)); err != nil {
		s.logger.Error("Failed to update password",
			zap.String("id", id.String()),
			zap.Error(err))
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *adminUserService) VerifyEmail(ctx context.Context, id uuid.UUID) error {
	if err := s.adminUserRepository.UpdateEmailVerification(ctx, id, true); err != nil {
		s.logger.Error("Failed to verify email",
			zap.String("id", id.String()),
			zap.Error(err))
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}

func (s *adminUserService) GetByUniversity(ctx context.Context, universityID uuid.UUID, pagination *dto.PaginationQuery) (*dto.PaginatedList[dto.AdminUserResponse], error) {
	result, err := s.adminUserRepository.FindByUniversity(ctx, universityID, pagination)
	if err != nil {
		s.logger.Error("Failed to fetch users by university",
			zap.String("university_id", universityID.String()),
			zap.Error(err))
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	dtos := make([]dto.AdminUserResponse, 0, len(result.Items))
	for _, user := range result.Items {
		dto, err := s.mapUserToDTO(&user)
		if err != nil {
			return nil, fmt.Errorf("failed to map user data: %w", err)
		}
		dtos = append(dtos, *dto)
	}

	return &dto.PaginatedList[dto.AdminUserResponse]{
		Items: dtos,
		Total: result.Total,
		Page:  result.Page,
		Limit: result.Limit,
	}, nil
}

func (s *adminUserService) GetByFaculty(ctx context.Context, facultyID uuid.UUID, pagination *dto.PaginationQuery) (*dto.PaginatedList[dto.AdminUserResponse], error) {
	result, err := s.adminUserRepository.FindByFaculty(ctx, facultyID, pagination)
	if err != nil {
		s.logger.Error("Failed to fetch users by faculty",
			zap.String("faculty_id", facultyID.String()),
			zap.Error(err))
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	dtos := make([]dto.AdminUserResponse, 0, len(result.Items))
	for _, user := range result.Items {
		dto, err := s.mapUserToDTO(&user)
		if err != nil {
			return nil, fmt.Errorf("failed to map user data: %w", err)
		}
		dtos = append(dtos, *dto)
	}

	return &dto.PaginatedList[dto.AdminUserResponse]{
		Items: dtos,
		Total: result.Total,
		Page:  result.Page,
		Limit: result.Limit,
	}, nil
}

func (s *adminUserService) mapUserToDTO(user *models.User) (*dto.AdminUserResponse, error) {
	return &dto.AdminUserResponse{
		ID:            user.ID,
		Email:         user.Email,
		StudentID:     user.StudentID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		UniversityID:  user.UniversityID,
		FacultyID:     user.FacultyID,
		Gender:        user.Gender,
		EmailVerified: user.EmailVerified,
		IsAdmin:       user.IsAdmin,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}
