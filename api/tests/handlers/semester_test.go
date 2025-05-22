package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// --- Mock SemesterService ---
type MockSemesterService struct {
	mock.Mock
}

func (m *MockSemesterService) Create(req *dto.CreateSemesterRequest) (*dto.SemesterResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SemesterResponse), args.Error(1)
}

func (m *MockSemesterService) Get(id uuid.UUID) (*dto.SemesterResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SemesterResponse), args.Error(1)
}

func (m *MockSemesterService) GetAll() ([]dto.SemesterResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SemesterResponse), args.Error(1)
}

func (m *MockSemesterService) Update(id uuid.UUID, req *dto.UpdateSemesterRequest) (*dto.SemesterResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SemesterResponse), args.Error(1)
}

func (m *MockSemesterService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// --- Test Setup Helper for Handler ---
func setupSemesterHandlerWithMocks(t *testing.T) (*handlers.SemesterHandler, *MockSemesterService) {
	mockService := new(MockSemesterService)
	logger, _ := zap.NewDevelopment()

	handler := handlers.NewSemesterHandler(
		mockService,
		logger,
	)
	return handler, mockService
}

// --- Handler Test Cases ---

func TestCreateSemester_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)

	reqBody := dto.CreateSemesterRequest{
		Year: 1404,
		Term: "fall",
	}
	expectedResponse := &dto.SemesterResponse{
		ID:        uuid.New(),
		Year:      reqBody.Year,
		Term:      reqBody.Term,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("Create", &reqBody).Return(expectedResponse, nil)

	jsonBody, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/semesters", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var actualResp dto.SemesterResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Year, actualResp.Year)
	assert.Equal(t, expectedResponse.Term, actualResp.Term)
	assert.NotEqual(t, uuid.Nil, actualResp.ID)
	assert.WithinDuration(t, expectedResponse.CreatedAt, actualResp.CreatedAt, 2*time.Second)

	mockService.AssertExpectations(t)
}

func TestCreateSemester_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupSemesterHandlerWithMocks(t)

	reqBody := dto.CreateSemesterRequest{
		Year: 1404,
		Term: "summer", // Invalid term
	}

	jsonBody, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/semesters", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var respBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request format", respBody["error"])
}

func TestCreateSemester_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)

	reqBody := dto.CreateSemesterRequest{
		Year: 1404,
		Term: "fall",
	}

	// Provide the base message to NewConflictError
	mockService.On("Create", &reqBody).Return(nil, errors.NewConflictError("semester already exists for this year and term"))

	jsonBody, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/semesters", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	// Assert the full message returned by the handler
	assert.Equal(t, "semester already exists for this year and term is conflicting", respBody["error"])
	mockService.AssertExpectations(t)
}

func TestGetSemester_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)

	semesterID := uuid.New()
	expectedResponse := &dto.SemesterResponse{
		ID:        semesterID,
		Year:      1403,
		Term:      "spring",
		CreatedAt: time.Now().Add(-10 * 24 * time.Hour),
		UpdatedAt: time.Now().Add(-5 * 24 * time.Hour),
	}

	mockService.On("Get", semesterID).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: semesterID.String()}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/semesters/"+semesterID.String(), nil)

	handler.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualResp dto.SemesterResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.ID, actualResp.ID)
	assert.Equal(t, expectedResponse.Year, actualResp.Year)
	assert.Equal(t, expectedResponse.Term, actualResp.Term)
	assert.WithinDuration(t, expectedResponse.CreatedAt, actualResp.CreatedAt, time.Second)
	assert.WithinDuration(t, expectedResponse.UpdatedAt, actualResp.UpdatedAt, time.Second)

	mockService.AssertExpectations(t)
}

func TestGetSemester_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)

	semesterID := uuid.New()
	mockService.On("Get", semesterID).Return(nil, errors.NewNotFoundError("Semester", semesterID.String()))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: semesterID.String()}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/semesters/"+semesterID.String(), nil)

	handler.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Semester not found", respBody["error"])
	mockService.AssertExpectations(t)
}

func TestGetSemester_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupSemesterHandlerWithMocks(t)

	invalidID := "not-a-uuid"
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: invalidID}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/semesters/"+invalidID, nil)

	handler.Get(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid semester ID", respBody["error"])
}

func TestGetAllSemesters_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)

	expectedResponse := []dto.SemesterResponse{
		{ID: uuid.New(), Year: 1403, Term: "fall", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Year: 1404, Term: "spring", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	mockService.On("GetAll").Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/semesters", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualResp []dto.SemesterResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedResponse), len(actualResp))
	if len(expectedResponse) > 0 && len(actualResp) > 0 {
		assert.Equal(t, expectedResponse[0].Year, actualResp[0].Year)
	}
	mockService.AssertExpectations(t)
}

func TestGetAllSemesters_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)

	expectedResponse := []dto.SemesterResponse{}
	mockService.On("GetAll").Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/semesters", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualResp []dto.SemesterResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Empty(t, actualResp)
	mockService.AssertExpectations(t)
}

func TestUpdateSemester_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)

	semesterID := uuid.New()
	reqBody := dto.UpdateSemesterRequest{
		Year: 1404,
		Term: "spring",
	}
	expectedResponse := &dto.SemesterResponse{
		ID:        semesterID,
		Year:      reqBody.Year,
		Term:      reqBody.Term,
		UpdatedAt: time.Now(),
	}

	mockService.On("Update", semesterID, &reqBody).Return(expectedResponse, nil)

	jsonBody, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: semesterID.String()}}
	c.Request, _ = http.NewRequest(http.MethodPut, "/semesters/"+semesterID.String(), bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualResp dto.SemesterResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.ID, actualResp.ID)
	assert.Equal(t, expectedResponse.Year, actualResp.Year)
	assert.Equal(t, expectedResponse.Term, actualResp.Term)
	assert.WithinDuration(t, expectedResponse.UpdatedAt, actualResp.UpdatedAt, time.Second)
	mockService.AssertExpectations(t)
}

func TestUpdateSemester_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupSemesterHandlerWithMocks(t)
	semesterID := uuid.New()

	reqBody := dto.UpdateSemesterRequest{
		Year: 1404,
		Term: "winter", // Invalid term
	}

	jsonBody, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: semesterID.String()}}
	c.Request, _ = http.NewRequest(http.MethodPut, "/semesters/"+semesterID.String(), bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var respBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request format", respBody["error"])
}

func TestUpdateSemester_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)
	semesterID := uuid.New()
	reqBody := dto.UpdateSemesterRequest{Year: 2025, Term: "fall"}

	mockService.On("Update", semesterID, &reqBody).Return(nil, errors.NewNotFoundError("Semester", semesterID.String()))

	jsonBody, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: semesterID.String()}}
	c.Request, _ = http.NewRequest(http.MethodPut, "/semesters/"+semesterID.String(), bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Semester not found", respBody["error"])
	mockService.AssertExpectations(t)
}

func TestDeleteSemester_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)
	semesterID := uuid.New()

	mockService.On("Delete", semesterID).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: semesterID.String()}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/semesters/"+semesterID.String(), nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Semester deleted successfully", respBody["message"])
	mockService.AssertExpectations(t)
}

func TestDeleteSemester_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)
	semesterID := uuid.New()

	mockService.On("Delete", semesterID).Return(errors.NewNotFoundError("Semester", semesterID.String()))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: semesterID.String()}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/semesters/"+semesterID.String(), nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Semester not found", respBody["error"])
	mockService.AssertExpectations(t)
}

func TestDeleteSemester_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupSemesterHandlerWithMocks(t)

	invalidID := "not-a-uuid"
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: invalidID}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/semesters/"+invalidID, nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid semester ID", respBody["error"])
}

func TestCreateSemester_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)

	reqBody := dto.CreateSemesterRequest{Year: 1404, Term: "fall"}
	mockService.On("Create", &reqBody).Return(nil, fmt.Errorf("database error"))

	jsonBody, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/semesters", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Internal server error", respBody["error"])
	mockService.AssertExpectations(t)
}

func TestGetSemester_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)
	semesterID := uuid.New()

	mockService.On("Get", semesterID).Return(nil, fmt.Errorf("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: semesterID.String()}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/semesters/"+semesterID.String(), nil)

	handler.Get(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetAllSemesters_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)

	mockService.On("GetAll").Return(nil, fmt.Errorf("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/semesters", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateSemester_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)
	semesterID := uuid.New()
	reqBody := dto.UpdateSemesterRequest{Year: 2025, Term: "spring"}

	mockService.On("Update", semesterID, &reqBody).Return(nil, fmt.Errorf("database error"))

	jsonBody, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: semesterID.String()}}
	c.Request, _ = http.NewRequest(http.MethodPut, "/semesters/"+semesterID.String(), bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Update(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteSemester_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)
	semesterID := uuid.New()

	mockService.On("Delete", semesterID).Return(fmt.Errorf("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: semesterID.String()}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/semesters/"+semesterID.String(), nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateSemester_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupSemesterHandlerWithMocks(t)

	semesterID := uuid.New()
	reqBody := dto.UpdateSemesterRequest{
		Year: 1404,
		Term: "fall",
	}

	// Provide the base message to NewConflictError
	mockService.On("Update", semesterID, &reqBody).Return(nil, errors.NewConflictError("semester already exists for this year and term"))

	jsonBody, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: semesterID.String()}}
	c.Request, _ = http.NewRequest(http.MethodPut, "/semesters/"+semesterID.String(), bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Update(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	// Assert the full message returned by the handler
	assert.Equal(t, "semester already exists for this year and term is conflicting", respBody["error"])
	mockService.AssertExpectations(t)
}
