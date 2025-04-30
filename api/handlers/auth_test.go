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
func (m *MockAuthService) Login(email, password string) (string, int, string, int, error) {
	args := m.Called(email, password)
	return args.String(0), args.Int(1), args.String(2), args.Int(3), args.Error(4)
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

func (m *MockAuthService) Refresh(old string) (string, int, string, int, error) {
	args := m.Called(old)
	return args.String(0), args.Int(1), args.String(2), args.Int(3), args.Error(4)
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

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewAuthHandler(mockService, logger)

	reqBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Mock service to return tokens and expiry times
	accessToken := "access_token_123"
	accessExpiry := 3600
	refreshToken := "refresh_token_456"
	refreshExpiry := 7200
	mockService.On("Login", reqBody.Email, reqBody.Password).
		Return(accessToken, accessExpiry, refreshToken, refreshExpiry, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	// Check status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response body contains access token and expiry
	var response dto.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, accessToken, response.AccessToken)
	assert.Equal(t, accessExpiry, response.ExpiresIn)

	// Verify refresh token cookie
	cookies := w.Result().Cookies()
	var refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshCookie = cookie
			break
		}
	}
	assert.NotNil(t, refreshCookie)
	assert.Equal(t, refreshToken, refreshCookie.Value)
	assert.Equal(t, refreshExpiry, refreshCookie.MaxAge)
	assert.True(t, refreshCookie.HttpOnly)
	assert.True(t, refreshCookie.Secure)

	mockService.AssertExpectations(t)
}

func TestRefreshHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewAuthHandler(mockService, logger)

	// Mock service to return new tokens and expiry times
	newAccessToken := "new_access_token"
	newAccessExpiry := 3600
	newRefreshToken := "new_refresh_token"
	newRefreshExpiry := 7200
	mockService.On("Refresh", "old_refresh_token").
		Return(newAccessToken, newAccessExpiry, newRefreshToken, newRefreshExpiry, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/refresh", nil)

	// Set refresh token in cookie
	c.Request.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "old_refresh_token",
	})

	handler.Refresh(c)

	// Check status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response body contains new access token and expiry
	var response dto.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, newAccessToken, response.AccessToken)
	assert.Equal(t, newAccessExpiry, response.ExpiresIn)

	// Verify new refresh token cookie
	cookies := w.Result().Cookies()
	var refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshCookie = cookie
			break
		}
	}
	assert.NotNil(t, refreshCookie)
	assert.Equal(t, newRefreshToken, refreshCookie.Value)
	assert.Equal(t, newRefreshExpiry, refreshCookie.MaxAge)
	assert.True(t, refreshCookie.HttpOnly)
	assert.True(t, refreshCookie.Secure)

	mockService.AssertExpectations(t)
}

func TestLogoutHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewAuthHandler(mockService, logger)

	// Mock service expectations
	mockService.On("Logout", "current_refresh_token").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/logout", nil)

	// Set refresh token in cookie
	c.Request.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "current_refresh_token",
	})

	handler.Logout(c)

	// Check status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "logged out successfully", response["message"])

	// Verify cookie was cleared
	cookies := w.Result().Cookies()
	var refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshCookie = cookie
			break
		}
	}
	assert.NotNil(t, refreshCookie)
	assert.Equal(t, "", refreshCookie.Value)
	assert.True(t, refreshCookie.MaxAge < 0)
	assert.True(t, refreshCookie.HttpOnly)
	assert.True(t, refreshCookie.Secure)

	mockService.AssertExpectations(t)
}

func TestLogoutHandler_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewAuthHandler(mockService, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/logout", nil)

	handler.Logout(c)

	// Check status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verify response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "missing refresh token", response["error"])

	mockService.AssertNotCalled(t, "Logout")
}
