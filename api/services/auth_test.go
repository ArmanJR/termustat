package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	infraMailer "github.com/armanjr/termustat/api/infrastructure/mailer"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/services"
	"github.com/armanjr/termustat/api/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// --- Mock Repository Implementation ---

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) FindUserByEmailOrStudentID(ctx context.Context, email, studentID string) (*models.User, error) {
	args := m.Called(ctx, email, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) FindUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) CreateEmailVerification(ctx context.Context, verification *models.EmailVerification) error {
	args := m.Called(ctx, verification)
	return args.Error(0)
}

func (m *MockAuthRepository) CreatePasswordReset(ctx context.Context, reset *models.PasswordReset) error {
	args := m.Called(ctx, reset)
	return args.Error(0)
}

func (m *MockAuthRepository) FindPasswordResetByToken(ctx context.Context, token string) (*models.PasswordReset, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PasswordReset), args.Error(1)
}

func (m *MockAuthRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	args := m.Called(ctx, userID, hashedPassword)
	return args.Error(0)
}

func (m *MockAuthRepository) DeletePasswordReset(ctx context.Context, reset *models.PasswordReset) error {
	args := m.Called(ctx, reset)
	return args.Error(0)
}

func (m *MockAuthRepository) FindEmailVerificationByToken(ctx context.Context, token string) (*models.EmailVerification, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EmailVerification), args.Error(1)
}

func (m *MockAuthRepository) VerifyUserEmail(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthRepository) DeleteEmailVerification(ctx context.Context, verification *models.EmailVerification) error {
	args := m.Called(ctx, verification)
	return args.Error(0)
}

// --- Mock Refresh Token Repository ---

type MockRefreshRepo struct {
	mock.Mock
}

func (m *MockRefreshRepo) Create(rt *models.RefreshToken) error {
	return m.Called(rt).Error(0)
}

func (m *MockRefreshRepo) Find(token string) (*models.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshRepo) Revoke(id uuid.UUID) error {
	return m.Called(id).Error(0)
}

func (m *MockRefreshRepo) RevokeAllForUser(userID uuid.UUID) error {
	return m.Called(userID).Error(0)
}

func (m *MockRefreshRepo) CleanupExpired() error {
	return m.Called().Error(0)
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

// --- Mock University Service ---
type MockUniversityService struct {
	mock.Mock
}

func (m *MockUniversityService) Get(ctx context.Context, id uuid.UUID) (*dto.UniversityResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UniversityResponse), args.Error(1)
}

// Implement other methods if needed by tests, otherwise leave empty or panic
func (m *MockUniversityService) Create(ctx context.Context, req *dto.CreateUniversityRequest) (*dto.UniversityResponse, error) {
	panic("Create not implemented in mock")
}
func (m *MockUniversityService) GetAll(ctx context.Context) ([]dto.UniversityResponse, error) {
	panic("GetAll not implemented in mock")
}
func (m *MockUniversityService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateUniversityRequest) (*dto.UniversityResponse, error) {
	panic("Update not implemented in mock")
}
func (m *MockUniversityService) ExistsByName(ctx context.Context, nameEn, nameFa string) (bool, error) {
	panic("ExistsByName not implemented in mock")
}
func (m *MockUniversityService) Delete(ctx context.Context, id uuid.UUID) error {
	panic("Delete not implemented in mock")
}

// --- Mock Faculty Service ---
type MockFacultyService struct {
	mock.Mock
}

func (m *MockFacultyService) Get(id uuid.UUID) (*dto.FacultyResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FacultyResponse), args.Error(1)
}

// Implement other methods if needed by tests
func (m *MockFacultyService) Create(dto dto.CreateFacultyDTO) (*dto.FacultyResponse, error) {
	panic("Create not implemented in mock")
}
func (m *MockFacultyService) GetAllByUniversity(universityID uuid.UUID) ([]*dto.FacultyResponse, error) {
	panic("GetAllByUniversity not implemented in mock")
}
func (m *MockFacultyService) GetByUniversityAndShortCode(universityID uuid.UUID, shortCode string) (*dto.FacultyResponse, error) {
	panic("GetByUniversityAndShortCode not implemented in mock")
}
func (m *MockFacultyService) Update(id uuid.UUID, dto dto.UpdateFacultyDTO) (*dto.FacultyResponse, error) {
	panic("Update not implemented in mock")
}
func (m *MockFacultyService) Delete(id uuid.UUID) error {
	panic("Delete not implemented in mock")
}

// --- Mock AuthService ---
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req *dto.RegisterServiceRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
func (m *MockAuthService) Login(ctx context.Context, email, password string) (string, int, string, int, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Int(1), args.String(2), args.Int(3), args.Error(4)
}
func (m *MockAuthService) ForgotPassword(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}
func (m *MockAuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}
func (m *MockAuthService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockAuthService) VerifyEmail(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*utils.JWTClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.JWTClaims), args.Error(1)
}
func (m *MockAuthService) Refresh(ctx context.Context, oldToken string) (string, int, string, int, error) {
	args := m.Called(ctx, oldToken)
	return args.String(0), args.Int(1), args.String(2), args.Int(3), args.Error(4)
}
func (m *MockAuthService) Logout(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

// --- Constants ---
const testJWTSecret = "test-secret-key-for-jwt"

// --- Test Setup Helper ---
func setupAuthService(t *testing.T) (services.AuthService, *MockAuthRepository, *MockRefreshRepo, *MockMailerService) {
	mockRepo := new(MockAuthRepository)
	mockRTRepo := new(MockRefreshRepo)
	mockMailer := new(MockMailerService)
	logger, _ := zap.NewDevelopment()

	service := services.NewAuthService(
		mockRepo,
		mockRTRepo,
		mockMailer,
		logger,
		testJWTSecret,
		15*time.Minute, // Short TTL for testing
		1*time.Hour,    // Short TTL for testing
		"http://localhost:3000",
	)
	return service, mockRepo, mockRTRepo, mockMailer
}

// Helper to hash password for mock user setup
func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// --- Test Cases ---

func TestRegisterService_Success(t *testing.T) {
	// Setup mocks and logger.
	mockRepo := new(MockAuthRepository)
	mockRTRepo := new(MockRefreshRepo)
	mockMailer := new(MockMailerService)
	logger, _ := zap.NewDevelopment()

	service := services.NewAuthService(
		mockRepo,
		mockRTRepo,
		mockMailer,
		logger,
		"test-secret",
		24*time.Hour,
		720*time.Hour,
		"http://localhost:3000",
	)

	ctx := context.Background()
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
	mockRepo.On("FindUserByEmailOrStudentID", ctx, req.Email, req.StudentID).
		Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(nil)
	// Expect CreateEmailVerification to be called.
	mockRepo.On("CreateEmailVerification", ctx, mock.AnythingOfType("*models.EmailVerification")).Return(nil)
	// Expect the mailer to send the verification email.
	mockMailer.On("SendVerificationEmail", mock.AnythingOfType("*models.User"), mock.AnythingOfType("string")).Return(nil)

	// Call Register.
	err := service.Register(ctx, req)

	// Assertions.
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockMailer.AssertExpectations(t)
}

func TestRegisterService_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	mockRTRepo := new(MockRefreshRepo)
	mockMailer := new(MockMailerService)
	logger, _ := zap.NewDevelopment()

	service := services.NewAuthService(
		mockRepo,
		mockRTRepo,
		mockMailer,
		logger,
		"test-secret",
		24*time.Hour,
		720*time.Hour,
		"http://localhost:3000",
	)

	ctx := context.Background()
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
	mockRepo.On("FindUserByEmailOrStudentID", ctx, req.Email, req.StudentID).Return(existingUser, nil)

	err := service.Register(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "email or student ID already exists", err.Error())
	mockRepo.AssertExpectations(t)
	// Mailer should not be called.
	mockMailer.AssertNotCalled(t, "SendVerificationEmail", mock.Anything)
}

func TestRegisterService_MailerError(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	mockRTRepo := new(MockRefreshRepo)
	mockMailer := new(MockMailerService)
	logger, _ := zap.NewDevelopment()

	service := services.NewAuthService(
		mockRepo,
		mockRTRepo,
		mockMailer,
		logger,
		"test-secret",
		24*time.Hour,
		720*time.Hour,
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

	mockRepo.
		On("FindUserByEmailOrStudentID", mock.Anything, req.Email, req.StudentID).
		Return(nil, gorm.ErrRecordNotFound)

	mockRepo.
		On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
		Return(nil)

	mockRepo.
		On("CreateEmailVerification", mock.Anything, mock.AnythingOfType("*models.EmailVerification")).
		Return(nil)

	mockMailer.
		On("SendVerificationEmail",
			mock.AnythingOfType("*models.User"),
			mock.AnythingOfType("string")).
		Return(assert.AnError)

	err := service.Register(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockMailer.AssertExpectations(t)
}

func TestLoginService_AdminScope(t *testing.T) {
	service, mockRepo, mockRTRepo, _ := setupAuthService(t)
	email := "admin@example.com"
	password := "password123"
	userID := uuid.New()

	adminUser := &models.User{
		ID:            userID,
		Email:         email,
		PasswordHash:  hashPassword(password),
		EmailVerified: true,
		IsAdmin:       true, // User is admin
	}

	mockRepo.On("FindUserByEmail", mock.Anything, email).Return(adminUser, nil)
	mockRTRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	accessToken, _, refreshToken, _, err := service.Login(context.Background(), email, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// Verify scope in access token
	claims, parseErr := utils.ParseJWT(accessToken, testJWTSecret)
	assert.NoError(t, parseErr)
	assert.NotNil(t, claims)
	assert.Contains(t, claims.Scopes, "admin-dashboard", "Admin user should have admin-dashboard scope")
	assert.Len(t, claims.Scopes, 1)
	assert.Equal(t, userID.String(), claims.UserID)

	mockRepo.AssertExpectations(t)
	mockRTRepo.AssertExpectations(t)
}

func TestLoginService_NoAdminScope(t *testing.T) {
	service, mockRepo, mockRTRepo, _ := setupAuthService(t)
	email := "user@example.com"
	password := "password123"
	userID := uuid.New()

	nonAdminUser := &models.User{
		ID:            userID,
		Email:         email,
		PasswordHash:  hashPassword(password),
		EmailVerified: true,
		IsAdmin:       false, // User is NOT admin
	}

	mockRepo.On("FindUserByEmail", mock.Anything, email).Return(nonAdminUser, nil)
	mockRTRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	accessToken, _, refreshToken, _, err := service.Login(context.Background(), email, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// Verify scope in access token
	claims, parseErr := utils.ParseJWT(accessToken, testJWTSecret)
	assert.NoError(t, parseErr)
	assert.NotNil(t, claims)
	assert.Empty(t, claims.Scopes, "Non-admin user should have no scopes") // Check for empty scopes
	assert.Equal(t, userID.String(), claims.UserID)

	mockRepo.AssertExpectations(t)
	mockRTRepo.AssertExpectations(t)
}

func TestRefreshService_AdminScopePreserved(t *testing.T) {
	service, _, mockRTRepo, _ := setupAuthService(t)
	oldRefreshToken := "old-refresh-token-admin"
	adminUserID := uuid.New()

	adminUser := models.User{ // Define the admin user explicitly
		ID:      adminUserID,
		IsAdmin: true,
	}

	mockOldRT := &models.RefreshToken{
		ID:        uuid.New(),
		Token:     oldRefreshToken,
		UserID:    adminUserID,
		User:      adminUser, // Associate the user with the token
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockRTRepo.On("Find", oldRefreshToken).Return(mockOldRT, nil)
	mockRTRepo.On("Revoke", mockOldRT.ID).Return(nil)
	mockRTRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	newAccessToken, _, newRefreshToken, _, err := service.Refresh(context.Background(), oldRefreshToken)

	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)
	assert.NotEqual(t, oldRefreshToken, newRefreshToken)

	// Verify scope in the new access token
	claims, parseErr := utils.ParseJWT(newAccessToken, testJWTSecret)
	assert.NoError(t, parseErr)
	assert.NotNil(t, claims)
	assert.Contains(t, claims.Scopes, "admin-dashboard", "Admin scope should be preserved on refresh")
	assert.Len(t, claims.Scopes, 1)
	assert.Equal(t, adminUserID.String(), claims.UserID)

	mockRTRepo.AssertExpectations(t)
}

func TestRefreshService_NoAdminScopePreserved(t *testing.T) {
	service, _, mockRTRepo, _ := setupAuthService(t)
	oldRefreshToken := "old-refresh-token-user"
	nonAdminUserID := uuid.New()

	nonAdminUser := models.User{ // Define the non-admin user
		ID:      nonAdminUserID,
		IsAdmin: false,
	}

	mockOldRT := &models.RefreshToken{
		ID:        uuid.New(),
		Token:     oldRefreshToken,
		UserID:    nonAdminUserID,
		User:      nonAdminUser, // Associate the user
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockRTRepo.On("Find", oldRefreshToken).Return(mockOldRT, nil)
	mockRTRepo.On("Revoke", mockOldRT.ID).Return(nil)
	mockRTRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	newAccessToken, _, newRefreshToken, _, err := service.Refresh(context.Background(), oldRefreshToken)

	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)
	assert.NotEqual(t, oldRefreshToken, newRefreshToken)

	// Verify scope in the new access token
	claims, parseErr := utils.ParseJWT(newAccessToken, testJWTSecret)
	assert.NoError(t, parseErr)
	assert.NotNil(t, claims)
	assert.Empty(t, claims.Scopes, "Non-admin user should still have no scopes after refresh")
	assert.Equal(t, nonAdminUserID.String(), claims.UserID)

	mockRTRepo.AssertExpectations(t)
}

func TestRefreshService_InvalidToken(t *testing.T) {
	service, _, mockRTRepo, _ := setupAuthService(t)
	invalidToken := "invalid-token"

	mockRTRepo.On("Find", invalidToken).Return(nil, errors.New("simulated find error")) // Simulate token not found or other error

	_, _, _, _, err := service.Refresh(context.Background(), invalidToken)

	assert.Error(t, err)
	assert.Equal(t, "invalid refresh token", err.Error())
	mockRTRepo.AssertExpectations(t)
	// Ensure Revoke and Create are not called
	mockRTRepo.AssertNotCalled(t, "Revoke", mock.Anything)
	mockRTRepo.AssertNotCalled(t, "Create", mock.Anything)
}
