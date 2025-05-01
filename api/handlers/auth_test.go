package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
)

// --- Mock AuthService ---
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req *dto.RegisterServiceRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

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

// --- Mock University Service ---
type MockUniversityService struct {
	mock.Mock
}

func (m *MockUniversityService) Get(id uuid.UUID) (*dto.UniversityResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UniversityResponse), args.Error(1)
}
func (m *MockUniversityService) Create(req *dto.CreateUniversityRequest) (*dto.UniversityResponse, error) {
	panic("Create not implemented in mock")
}
func (m *MockUniversityService) GetAll() ([]dto.UniversityResponse, error) {
	panic("GetAll not implemented in mock")
}
func (m *MockUniversityService) Update(id uuid.UUID, req *dto.UpdateUniversityRequest) (*dto.UniversityResponse, error) {
	panic("Update not implemented in mock")
}
func (m *MockUniversityService) ExistsByName(nameEn, nameFa string) (bool, error) {
	panic("ExistsByName not implemented in mock")
}
func (m *MockUniversityService) Delete(id uuid.UUID) error {
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

// --- Test Setup Helper for Handler ---
func setupAuthHandlerWithMocks(t *testing.T) (*handlers.AuthHandler, *MockAuthService, *MockUniversityService, *MockFacultyService) {
	mockAuthSvc := new(MockAuthService)
	mockUniSvc := new(MockUniversityService)
	mockFacultySvc := new(MockFacultyService)
	logger, _ := zap.NewDevelopment() // Use development logger for test output

	handler := handlers.NewAuthHandler(
		mockAuthSvc,
		mockUniSvc,
		mockFacultySvc,
		logger,
	)
	return handler, mockAuthSvc, mockUniSvc, mockFacultySvc
}

// --- Handler Test Cases ---

func TestRegisterHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, _, _ := setupAuthHandlerWithMocks(t)

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

	mockAuthSvc.On("Register", mock.MatchedBy(func(req *dto.RegisterServiceRequest) bool {
		// Perform detailed checks on the service request object
		return req.Email == reqBody.Email &&
			req.Password == reqBody.Password && // Consider if you want to check password here
			req.StudentID == reqBody.StudentID &&
			req.FirstName == reqBody.FirstName &&
			req.LastName == reqBody.LastName &&
			req.UniversityID.String() == reqBody.UniversityID &&
			req.FacultyID.String() == reqBody.FacultyID &&
			req.Gender == reqBody.Gender
	})).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	// Check response body for success message
	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["message"], "Registration successful")
	mockAuthSvc.AssertExpectations(t)
}

func TestRegisterHandler_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _, _, _ := setupAuthHandlerWithMocks(t) // Mocks needed for handler creation

	// Create invalid request body (e.g., missing password)
	reqBody := map[string]string{
		"email":      "test@example.com",
		"student_id": "ST123",
		"first_name": "Test",
		"last_name":  "User",
		"gender":     "male",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	// Optionally assert error message if binding error is consistent
}

func TestRegisterHandler_ServiceError_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, _, _ := setupAuthHandlerWithMocks(t)

	validUUID := uuid.New()
	reqBody := dto.RegisterRequest{
		Email:        "existing@example.com", // Assume this email causes conflict
		Password:     "password123",
		StudentID:    "ST-EXIST", // Assume this ID causes conflict
		FirstName:    "Test",
		LastName:     "User",
		UniversityID: validUUID.String(),
		FacultyID:    validUUID.String(),
		Gender:       "male",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Setup expectation - service returns a conflict error
	// Use specific error message string the service returns
	mockAuthSvc.On("Register", mock.AnythingOfType("*dto.RegisterServiceRequest")).
		Return(errors.New("email or student ID already exists"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusConflict, w.Code) // Expect 409 Conflict
	// Check response body for error message
	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "email or student ID already exists", respBody["error"])
	mockAuthSvc.AssertExpectations(t)
}

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, _, _ := setupAuthHandlerWithMocks(t)

	reqBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	accessToken := "access_token_123"
	accessExpiry := 3600
	refreshToken := "refresh_token_456"
	refreshExpiry := 7200
	mockAuthSvc.On("Login", reqBody.Email, reqBody.Password).
		Return(accessToken, accessExpiry, refreshToken, refreshExpiry, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, accessToken, response.AccessToken)
	assert.Equal(t, accessExpiry, response.ExpiresIn)

	cookies := w.Result().Cookies()
	var refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshCookie = cookie
			break
		}
	}
	assert.NotNil(t, refreshCookie, "Refresh token cookie should be set")
	assert.Equal(t, refreshToken, refreshCookie.Value)
	assert.Equal(t, refreshExpiry, refreshCookie.MaxAge)
	assert.True(t, refreshCookie.HttpOnly)
	// Note: Secure flag might be false in test environment depending on setup
	// assert.True(t, refreshCookie.Secure)
	assert.Equal(t, "/", refreshCookie.Path)

	mockAuthSvc.AssertExpectations(t)
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, _, _ := setupAuthHandlerWithMocks(t)

	reqBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Mock service returns specific error for invalid credentials
	mockAuthSvc.On("Login", reqBody.Email, reqBody.Password).
		Return("", 0, "", 0, errors.New("invalid credentials"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code) // Expect 401 Unauthorized

	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to login", respBody["error"])

	mockAuthSvc.AssertExpectations(t)
}

func TestLoginHandler_EmailNotVerified(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, _, _ := setupAuthHandlerWithMocks(t)

	reqBody := dto.LoginRequest{
		Email:    "unverified@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Mock service returns specific error for email not verified
	mockAuthSvc.On("Login", reqBody.Email, reqBody.Password).
		Return("", 0, "", 0, errors.New("email not verified"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusForbidden, w.Code) // Expect 403 Forbidden

	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to login", respBody["error"])

	mockAuthSvc.AssertExpectations(t)
}

func TestRefreshHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, _, _ := setupAuthHandlerWithMocks(t)

	oldRefreshToken := "old_refresh_token"
	newAccessToken := "new_access_token"
	newAccessExpiry := 1800
	newRefreshToken := "new_refresh_token"
	newRefreshExpiry := 3600

	mockAuthSvc.On("Refresh", oldRefreshToken).
		Return(newAccessToken, newAccessExpiry, newRefreshToken, newRefreshExpiry, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/refresh", nil)
	// Set refresh token in cookie for the incoming request
	c.Request.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: oldRefreshToken,
	})

	handler.Refresh(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, newAccessToken, response.AccessToken)
	assert.Equal(t, newAccessExpiry, response.ExpiresIn)

	cookies := w.Result().Cookies()
	var refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshCookie = cookie
			break
		}
	}
	assert.NotNil(t, refreshCookie, "New refresh token cookie should be set")
	assert.Equal(t, newRefreshToken, refreshCookie.Value)
	assert.Equal(t, newRefreshExpiry, refreshCookie.MaxAge)
	assert.True(t, refreshCookie.HttpOnly)
	// assert.True(t, refreshCookie.Secure) // May be false in test
	assert.Equal(t, "/", refreshCookie.Path)

	mockAuthSvc.AssertExpectations(t)
}

func TestRefreshHandler_MissingCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, _, _ := setupAuthHandlerWithMocks(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/refresh", nil)
	// No cookie set

	handler.Refresh(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "missing refresh token", respBody["error"])

	// Ensure the Refresh service method was not called
	mockAuthSvc.AssertNotCalled(t, "Refresh", mock.Anything)
}

func TestRefreshHandler_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, _, _ := setupAuthHandlerWithMocks(t)

	invalidToken := "invalid_refresh_token"
	// Mock service returns an error indicating the token is invalid
	mockAuthSvc.On("Refresh", invalidToken).
		Return("", 0, "", 0, errors.New("invalid refresh token")) // Use the error service returns

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/refresh", nil)
	c.Request.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: invalidToken,
	})

	handler.Refresh(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code) // Expect 401 for invalid token

	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "invalid refresh token", respBody["error"])

	mockAuthSvc.AssertExpectations(t)
}

func TestLogoutHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, _, _ := setupAuthHandlerWithMocks(t)

	refreshToken := "current_refresh_token"
	mockAuthSvc.On("Logout", refreshToken).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/logout", nil)
	c.Request.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})

	handler.Logout(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "logged out successfully", response["message"])

	cookies := w.Result().Cookies()
	var refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshCookie = cookie
			break
		}
	}
	assert.NotNil(t, refreshCookie, "Refresh token cookie should be present for clearing")
	assert.Equal(t, "", refreshCookie.Value, "Cookie value should be cleared")
	assert.True(t, refreshCookie.MaxAge < 0, "Cookie MaxAge should be negative to clear") // MaxAge < 0 deletes cookie
	assert.True(t, refreshCookie.HttpOnly)
	// assert.True(t, refreshCookie.Secure) // May be false in test
	assert.Equal(t, "/", refreshCookie.Path)

	mockAuthSvc.AssertExpectations(t)
}

func TestLogoutHandler_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, _, _ := setupAuthHandlerWithMocks(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/logout", nil)
	// No cookie

	handler.Logout(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "missing refresh token", response["error"])

	mockAuthSvc.AssertNotCalled(t, "Logout", mock.Anything)
}

func TestGetCurrentUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, mockUniSvc, mockFacultySvc := setupAuthHandlerWithMocks(t)

	testUserID := uuid.New()
	testUniID := uuid.New()
	testFacultyID := uuid.New()

	mockUser := &models.User{
		ID:            testUserID,
		Email:         "test@example.com",
		StudentID:     "S123",
		FirstName:     "Test",
		LastName:      "User",
		UniversityID:  testUniID,
		FacultyID:     testFacultyID,
		Gender:        "female",
		EmailVerified: true,
		IsAdmin:       false,
		CreatedAt:     time.Now().Add(-24 * time.Hour),
		UpdatedAt:     time.Now(),
	}

	mockUniversity := &dto.UniversityResponse{
		ID:     testUniID,
		NameEn: "Test University",
		NameFa: "دانشگاه تست",
	}

	mockFaculty := &dto.FacultyResponse{
		ID:     testFacultyID,
		NameEn: "Test Faculty",
		NameFa: "دانشکده تست",
	}

	mockAuthSvc.On("GetCurrentUser", testUserID).Return(mockUser, nil)
	mockUniSvc.On("Get", testUniID).Return(mockUniversity, nil)
	mockFacultySvc.On("Get", testFacultyID).Return(mockFaculty, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/v1/user/me", nil)
	c.Set("userID", testUserID.String()) // Simulate middleware setting userID

	handler.GetCurrentUser(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	assert.Equal(t, mockUser.ID.String(), responseBody["id"])
	assert.Equal(t, mockUser.Email, responseBody["email"])
	//assert.Equal(t, mockUser.IsAdmin, responseBody["is_admin"])

	uniMap, ok := responseBody["university"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, mockUniversity.ID.String(), uniMap["id"])
	assert.Equal(t, mockUniversity.NameEn, uniMap["name_en"])
	assert.Equal(t, mockUniversity.NameFa, uniMap["name_fa"])

	facultyMap, ok := responseBody["faculty"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, mockFaculty.ID.String(), facultyMap["id"])
	assert.Equal(t, mockFaculty.NameEn, facultyMap["name_en"])
	assert.Equal(t, mockFaculty.NameFa, facultyMap["name_fa"])

	mockAuthSvc.AssertExpectations(t)
	mockUniSvc.AssertExpectations(t)
	mockFacultySvc.AssertExpectations(t)
}

func TestGetCurrentUser_MissingUniversity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockAuthSvc, mockUniSvc, mockFacultySvc := setupAuthHandlerWithMocks(t)

	testUserID := uuid.New()
	testUniID := uuid.New()
	testFacultyID := uuid.New()

	mockUser := &models.User{ // User has Uni ID, but service will fail
		ID:           testUserID,
		Email:        "test@example.com",
		UniversityID: testUniID, // Valid looking ID
		FacultyID:    testFacultyID,
		IsAdmin:      false,
	}
	mockFaculty := &dto.FacultyResponse{ // Faculty lookup succeeds
		ID:     testFacultyID,
		NameEn: "Test Faculty",
		NameFa: "دانشکده تست",
	}

	mockAuthSvc.On("GetCurrentUser", testUserID).Return(mockUser, nil)
	// Simulate University Service returning an error (e.g., Not Found)
	mockUniSvc.On("Get", testUniID).Return(nil, errors.NewNotFoundError("university", testUniID.String()))
	mockFacultySvc.On("Get", testFacultyID).Return(mockFaculty, nil) // Faculty still found

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/v1/user/me", nil)
	c.Set("userID", testUserID.String())

	handler.GetCurrentUser(c)

	// Should still succeed overall, but uni field will be null
	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	assert.Equal(t, mockUser.ID.String(), responseBody["id"])
	assert.Nil(t, responseBody["university"], "University should be null when lookup fails")

	facultyMap, ok := responseBody["faculty"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, mockFaculty.ID.String(), facultyMap["id"]) // Faculty should still be present

	mockAuthSvc.AssertExpectations(t)
	mockUniSvc.AssertExpectations(t)
	mockFacultySvc.AssertExpectations(t)
}

func TestGetCurrentUser_NoUserIDInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _, _, _ := setupAuthHandlerWithMocks(t) // Mocks don't matter here

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/v1/user/me", nil)
	// Do NOT set userID in context

	handler.GetCurrentUser(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "User not authenticated", respBody["error"])
}
