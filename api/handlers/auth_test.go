package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/handlers"
	"github.com/armanjr/termustat/api/models"
	"github.com/armanjr/termustat/api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req *dto.RegisterServiceRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

// Implement other AuthService methods to satisfy the interface
func (m *MockAuthService) Login(email, password string) (string, string, error) {
	args := m.Called(email, password)
	return args.String(0), args.String(1), args.Error(1)
}

func (m *MockAuthService) ForgotPassword(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(token, password string) error {
	args := m.Called(token, password)
	return args.Error(0)
}

func (m *MockAuthService) VerifyEmail(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthService) GetCurrentUser(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (*utils.JWTClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.JWTClaims), args.Error(1)
}

func (m *MockAuthService) Refresh(old string) (string, string, error) {
	args := m.Called(old)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthService) Logout(refresh string) error {
	args := m.Called(refresh)
	return args.Error(0)
}

func TestRegisterHandler_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	logger, _ := zap.NewDevelopment()

	handler := handlers.NewAuthHandler(mockService, logger)

	// Create request body
	validUUID := uuid.New()
	reqBody := dto.RegisterRequest{
		Email:        "test@example.com",
		Password:     "password123",
		StudentID:    "ST12345",
		FirstName:    "Test",
		LastName:     "User",
		UniversityID: validUUID.String(),
		FacultyID:    validUUID.String(),
		Gender:       "male",
	}

	jsonBody, _ := json.Marshal(reqBody)

	// Setup expectation
	mockService.On("Register", mock.MatchedBy(func(req *dto.RegisterServiceRequest) bool {
		return req.Email == reqBody.Email &&
			req.StudentID == reqBody.StudentID &&
			req.UniversityID.String() == reqBody.UniversityID
	})).Return(nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	handler.Register(c)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestRegisterHandler_InvalidRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	logger, _ := zap.NewDevelopment()

	handler := handlers.NewAuthHandler(mockService, logger)

	// Create invalid request body (missing required fields)
	reqBody := map[string]string{
		"email": "test@example.com",
		// Missing other required fields
	}

	jsonBody, _ := json.Marshal(reqBody)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	handler.Register(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterHandler_ServiceError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	logger, _ := zap.NewDevelopment()

	handler := handlers.NewAuthHandler(mockService, logger)

	// Create request body
	validUUID := uuid.New()
	reqBody := dto.RegisterRequest{
		Email:        "existing@example.com",
		Password:     "password123",
		StudentID:    "ST12345",
		FirstName:    "Test",
		LastName:     "User",
		UniversityID: validUUID.String(),
		FacultyID:    validUUID.String(),
		Gender:       "male",
	}

	jsonBody, _ := json.Marshal(reqBody)

	// Setup expectation - service returns error
	mockService.On("Register", mock.Anything).Return(errors.New("email or student ID already exists"))

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	handler.Register(c)

	// Assertions
	assert.Equal(t, http.StatusConflict, w.Code)
	mockService.AssertExpectations(t)
}
