package services_test

import (
	"testing"
	"time"

	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	infraMailer "github.com/armanjr/termustat/api/infrastructure/mailer"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// --- Mock Repository Implementation ---

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthRepository) FindUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) FindUserByEmailOrStudentID(email, studentID string) (*models.User, error) {
	args := m.Called(email, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) FindUserByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) CreateEmailVerification(verification *models.EmailVerification) error {
	args := m.Called(verification)
	return args.Error(0)
}

func (m *MockAuthRepository) CreatePasswordReset(reset *models.PasswordReset) error {
	args := m.Called(reset)
	return args.Error(0)
}

func (m *MockAuthRepository) FindPasswordResetByToken(token string) (*models.PasswordReset, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PasswordReset), args.Error(1)
}

func (m *MockAuthRepository) UpdateUserPassword(userID uuid.UUID, hashedPassword string) error {
	args := m.Called(userID, hashedPassword)
	return args.Error(0)
}

func (m *MockAuthRepository) DeletePasswordReset(reset *models.PasswordReset) error {
	args := m.Called(reset)
	return args.Error(0)
}

func (m *MockAuthRepository) FindEmailVerificationByToken(token string) (*models.EmailVerification, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EmailVerification), args.Error(1)
}

func (m *MockAuthRepository) VerifyUserEmail(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAuthRepository) DeleteEmailVerification(verification *models.EmailVerification) error {
	args := m.Called(verification)
	return args.Error(0)
}

// --- Mock Mailer Implementation ---

type MockMailerService struct {
	mock.Mock
}

func (m *MockMailerService) SendVerificationEmail(user *models.User, token string) error {
	args := m.Called(user, token)
	return args.Error(0)
}

func (m *MockMailerService) SendPasswordResetEmail(user *models.User, resetToken string) error {
	args := m.Called(user, resetToken)
	return args.Error(0)
}

func (m *MockMailerService) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func (m *MockMailerService) RenderTemplate(tplName string, data interface{}) (*infraMailer.EmailTemplate, error) {
	args := m.Called(tplName, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*infraMailer.EmailTemplate), args.Error(1)
}

// --- Test Cases ---

func TestRegisterService_Success(t *testing.T) {
	// Setup mocks and logger.
	mockRepo := new(MockAuthRepository)
	mockMailer := new(MockMailerService)
	logger, _ := zap.NewDevelopment()

	service := services.NewAuthService(
		mockRepo,
		mockMailer,
		logger,
		"test-secret",
		24*time.Hour,
		"http://localhost:3000",
	)

	req := &dto.RegisterServiceRequest{
		Email:        "new@example.com",
		Password:     "password123",
		StudentID:    "ST12345",
		FirstName:    "Test",
		LastName:     "User",
		UniversityID: uuid.New(),
		FacultyID:    uuid.New(),
		Gender:       "male",
	}

	// Expect repository to indicate user not found.
	mockRepo.On("FindUserByEmailOrStudentID", req.Email, req.StudentID).
		Return(nil, errors.NewNotFoundError("user", ""))
	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)
	// Expect CreateEmailVerification to be called.
	mockRepo.On("CreateEmailVerification", mock.AnythingOfType("*models.EmailVerification")).Return(nil)
	// Expect the mailer to send the verification email.
	mockMailer.On("SendVerificationEmail", mock.AnythingOfType("*models.User"), mock.AnythingOfType("string")).Return(nil)

	// Call Register.
	err := service.Register(req)

	// Assertions.
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockMailer.AssertExpectations(t)
}

func TestRegisterService_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	mockMailer := new(MockMailerService)
	logger, _ := zap.NewDevelopment()

	service := services.NewAuthService(
		mockRepo,
		mockMailer,
		logger,
		"test-secret",
		24*time.Hour,
		"http://localhost:3000",
	)

	req := &dto.RegisterServiceRequest{
		Email:        "existing@example.com",
		Password:     "password123",
		StudentID:    "ST12345",
		FirstName:    "Test",
		LastName:     "User",
		UniversityID: uuid.New(),
		FacultyID:    uuid.New(),
		Gender:       "male",
	}

	existingUser := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		StudentID: req.StudentID,
	}

	// Repository returns an existing user.
	mockRepo.On("FindUserByEmailOrStudentID", req.Email, req.StudentID).Return(existingUser, nil)

	err := service.Register(req)

	assert.Error(t, err)
	assert.Equal(t, "email or student ID already exists", err.Error())
	mockRepo.AssertExpectations(t)
	// Mailer should not be called.
	mockMailer.AssertNotCalled(t, "SendVerificationEmail", mock.Anything)
}

func TestRegisterService_MailerError(t *testing.T) {
	// Even if the mailer fails, registration should return no error.
	mockRepo := new(MockAuthRepository)
	mockMailer := new(MockMailerService)
	logger, _ := zap.NewDevelopment()

	service := services.NewAuthService(
		mockRepo,
		mockMailer,
		logger,
		"test-secret",
		24*time.Hour,
		"http://localhost:3000",
	)

	req := &dto.RegisterServiceRequest{
		Email:        "new@example.com",
		Password:     "password123",
		StudentID:    "ST12345",
		FirstName:    "Test",
		LastName:     "User",
		UniversityID: uuid.New(),
		FacultyID:    uuid.New(),
		Gender:       "male",
	}

	mockRepo.On("FindUserByEmailOrStudentID", req.Email, req.StudentID).
		Return(nil, errors.NewNotFoundError("user", ""))
	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)
	mockRepo.On("CreateEmailVerification", mock.AnythingOfType("*models.EmailVerification")).Return(nil)

	// Simulate mailer error.
	mockMailer.On("SendVerificationEmail", mock.AnythingOfType("*models.User"), mock.AnythingOfType("string")).Return(nil)

	err := service.Register(req)

	// Registration should succeed even if the mailer fails.
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockMailer.AssertExpectations(t)
}
