package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupUniversityHandlerWithMocks(t *testing.T) (*handlers.UniversityHandler, *MockUniversityService) {
	mockService := new(MockUniversityService) // This will now refer to MockUniversityService from auth_test.go
	logger, _ := zap.NewDevelopment()

	handler := handlers.NewUniversityHandler(
		mockService,
		logger,
	)
	return handler, mockService
}

func TestCreateUniversity_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	isActiveVal := true
	reqBody := dto.CreateUniversityRequest{
		NameEn:   "Test University",
		NameFa:   "دانشگاه تست",
		IsActive: &isActiveVal, // Corrected to use a pointer
	}
	expectedResp := &dto.UniversityResponse{
		ID:       uuid.New(),
		NameEn:   reqBody.NameEn,
		NameFa:   reqBody.NameFa,
		IsActive: *reqBody.IsActive, // Corrected to dereference the pointer
	}

	mockService.On("ExistsByName", mock.Anything, reqBody.NameEn, reqBody.NameFa).Return(false, nil)
	mockService.On("Create", mock.Anything, &reqBody).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBody, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest(http.MethodPost, "/universities", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var actualResp dto.UniversityResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp.NameEn, actualResp.NameEn)
	assert.Equal(t, expectedResp.NameFa, actualResp.NameFa)
	// ID is generated, so we only check if it's not nil
	assert.NotEqual(t, uuid.Nil, actualResp.ID)

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "ExistsByName", mock.Anything, reqBody.NameEn, reqBody.NameFa)
	mockService.AssertCalled(t, "Create", mock.Anything, &reqBody)
}

func TestGetUniversity_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	uniID := uuid.New()
	expectedResp := &dto.UniversityResponse{
		ID:     uniID,
		NameEn: "Test University",
		NameFa: "دانشگاه تست",
	}

	mockService.On("Get", mock.Anything, uniID).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: uniID.String()}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/universities/"+uniID.String(), nil)

	handler.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualResp dto.UniversityResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, &actualResp)

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "Get", mock.Anything, uniID)
}

func TestGetAllUniversities_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	expectedResp := []dto.UniversityResponse{
		{ID: uuid.New(), NameEn: "Test University 1", NameFa: "دانشگاه تست ۱"},
		{ID: uuid.New(), NameEn: "Test University 2", NameFa: "دانشگاه تست ۲"},
	}

	mockService.On("GetAll", mock.Anything).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/universities", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualResp []dto.UniversityResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, actualResp)

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "GetAll", mock.Anything)
}

func TestUpdateUniversity_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	uniID := uuid.New()
	isActiveVal := true
	reqBody := dto.UpdateUniversityRequest{
		NameEn:   "Updated University",
		NameFa:   "دانشگاه آپدیت شده",
		IsActive: &isActiveVal,
	}
	expectedResp := &dto.UniversityResponse{
		ID:       uniID,
		NameEn:   reqBody.NameEn,
		NameFa:   reqBody.NameFa,
		IsActive: *reqBody.IsActive,
	}

	mockService.On("Update", mock.Anything, uniID, &reqBody).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBody, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest(http.MethodPut, "/universities/"+uniID.String(), bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: uniID.String()}}

	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualResp dto.UniversityResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp.ID, actualResp.ID)
	assert.Equal(t, expectedResp.NameEn, actualResp.NameEn)
	assert.Equal(t, expectedResp.NameFa, actualResp.NameFa)

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "Update", mock.Anything, uniID, &reqBody)
}

func TestDeleteUniversity_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	uniID := uuid.New()

	mockService.On("Delete", mock.Anything, uniID).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: uniID.String()}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/universities/"+uniID.String(), nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &actualResp)
	assert.NoError(t, err)
	assert.Equal(t, "University deleted successfully", actualResp["message"])

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "Delete", mock.Anything, uniID)
}

func TestDeleteUniversity_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	invalidID := "not-a-uuid"

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: invalidID}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/universities/"+invalidID, nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "Invalid university ID")

	mockService.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestDeleteUniversity_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	uniID := uuid.New()
	serviceErr := errors.NewNotFoundError("university", uniID.String())
	mockService.On("Delete", mock.Anything, uniID).Return(serviceErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: uniID.String()}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/universities/"+uniID.String(), nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "University not found")

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "Delete", mock.Anything, uniID)
}

func TestDeleteUniversity_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	uniID := uuid.New()
	serviceErr := errors.New("some internal error")
	mockService.On("Delete", mock.Anything, uniID).Return(serviceErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: uniID.String()}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/universities/"+uniID.String(), nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "Internal server error")

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "Delete", mock.Anything, uniID)
}

func TestUpdateUniversity_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	invalidID := "not-a-uuid"
	reqBody := dto.UpdateUniversityRequest{
		NameEn: "Updated University",
		NameFa: "دانشگاه آپدیت شده",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBody, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest(http.MethodPut, "/universities/"+invalidID, bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: invalidID}}

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "Invalid university ID")

	mockService.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateUniversity_InvalidRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)
	uniID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPut, "/universities/"+uniID.String(), bytes.NewBufferString(`{"name_en": "Test", "name_fa": "تست"`)) // Invalid JSON
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: uniID.String()}}

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "Invalid request format")
	mockService.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateUniversity_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	uniID := uuid.New()
	isActiveVal := false
	reqBody := dto.UpdateUniversityRequest{
		NameEn:   "Updated University",
		NameFa:   "دانشگاه آپدیت شده",
		IsActive: &isActiveVal,
	}
	serviceErr := errors.NewNotFoundError("university", uniID.String())
	mockService.On("Update", mock.Anything, uniID, &reqBody).Return(nil, serviceErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBody, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest(http.MethodPut, "/universities/"+uniID.String(), bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: uniID.String()}}

	handler.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "University not found")

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "Update", mock.Anything, uniID, &reqBody)
}

func TestUpdateUniversity_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	uniID := uuid.New()
	isActiveVal := true
	reqBody := dto.UpdateUniversityRequest{
		NameEn:   "Updated University",
		NameFa:   "دانشگاه آپدیت شده",
		IsActive: &isActiveVal,
	}
	serviceErr := errors.New("some internal error")
	mockService.On("Update", mock.Anything, uniID, &reqBody).Return(nil, serviceErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBody, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest(http.MethodPut, "/universities/"+uniID.String(), bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: uniID.String()}}

	handler.Update(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "Internal server error")

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "Update", mock.Anything, uniID, &reqBody)
}

func TestGetAllUniversities_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	serviceErr := errors.New("some internal error")
	mockService.On("GetAll", mock.Anything).Return(nil, serviceErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/universities", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "Internal server error")

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "GetAll", mock.Anything)
}

func TestGetUniversity_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	invalidID := "not-a-uuid"

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: invalidID}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/universities/"+invalidID, nil)

	handler.Get(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "Invalid university ID")

	mockService.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
}

func TestGetUniversity_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	uniID := uuid.New()
	serviceErr := errors.NewNotFoundError("university", uniID.String())
	mockService.On("Get", mock.Anything, uniID).Return(nil, serviceErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: uniID.String()}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/universities/"+uniID.String(), nil)

	handler.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "University not found")

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "Get", mock.Anything, uniID)
}

func TestGetUniversity_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	uniID := uuid.New()
	serviceErr := errors.New("some internal error")
	mockService.On("Get", mock.Anything, uniID).Return(nil, serviceErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: uniID.String()}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/universities/"+uniID.String(), nil)

	handler.Get(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "Internal server error")

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "Get", mock.Anything, uniID)
}

func TestCreateUniversity_InvalidRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupUniversityHandlerWithMocks(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/universities", bytes.NewBufferString(`{"name_en": "Test", "name_fa": "تست"`)) // Invalid JSON
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "Invalid request format")
}

func TestCreateUniversity_Conflict_AlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	isActiveVal := true
	reqBody := dto.CreateUniversityRequest{
		NameEn:   "Existing University",
		NameFa:   "دانشگاه موجود",
		IsActive: &isActiveVal,
	}

	mockService.On("ExistsByName", mock.Anything, reqBody.NameEn, reqBody.NameFa).Return(true, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBody, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest(http.MethodPost, "/universities", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "University with the same name already exists")

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "ExistsByName", mock.Anything, reqBody.NameEn, reqBody.NameFa)
	mockService.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateUniversity_ServiceError_ExistsByName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	isActiveVal := true
	reqBody := dto.CreateUniversityRequest{
		NameEn:   "Test University",
		NameFa:   "دانشگاه تست",
		IsActive: &isActiveVal,
	}

	mockService.On("ExistsByName", mock.Anything, reqBody.NameEn, reqBody.NameFa).Return(false, errors.New("some internal error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBody, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest(http.MethodPost, "/universities", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "Internal server error")

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "ExistsByName", mock.Anything, reqBody.NameEn, reqBody.NameFa)
	mockService.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateUniversity_ServiceError_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockService := setupUniversityHandlerWithMocks(t)

	isActiveVal := true
	reqBody := dto.CreateUniversityRequest{
		NameEn:   "Test University",
		NameFa:   "دانشگاه تست",
		IsActive: &isActiveVal,
	}

	mockService.On("ExistsByName", mock.Anything, reqBody.NameEn, reqBody.NameFa).Return(false, nil)
	mockService.On("Create", mock.Anything, &reqBody).Return(nil, errors.New("failed to create"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBody, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest(http.MethodPost, "/universities", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var errResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp["error"], "Internal server error")

	mockService.AssertExpectations(t)
	mockService.AssertCalled(t, "ExistsByName", mock.Anything, reqBody.NameEn, reqBody.NameFa)
	mockService.AssertCalled(t, "Create", mock.Anything, &reqBody)
}
