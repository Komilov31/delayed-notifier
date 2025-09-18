package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Komilov31/delayed-notifier/internal/dto"
	"github.com/Komilov31/delayed-notifier/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNotifierService is a mock implementation of NotifierService
type MockNotifierService struct {
	mock.Mock
}

func (m *MockNotifierService) CreateNotification(notification model.Notification) (*model.Notification, error) {
	args := m.Called(notification)
	return args.Get(0).(*model.Notification), args.Error(1)
}

func (m *MockNotifierService) GetNotificationStatus(id int) (*dto.NotificationStatus, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.NotificationStatus), args.Error(1)
}

func (m *MockNotifierService) GetAllNotifications() ([]model.Notification, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Notification), args.Error(1)
}

func (m *MockNotifierService) UpdateNotificationStatus(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockNotifierService) PublishReadyNotifications(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockNotifierService) ConsumeMessages(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestHandler_CreateNotification_Success(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	notificationDTO := dto.NotificationDTO{
		Text:       "Test notification",
		TelegramId: 123,
		SendAt:     time.Now().Add(time.Hour),
	}
	body, _ := json.Marshal(notificationDTO)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	expectedNotification := &model.Notification{
		Id:         1,
		Text:       notificationDTO.Text,
		TelegramId: notificationDTO.TelegramId,
		SendAt:     int(notificationDTO.SendAt.UnixMilli()),
		Status:     "active",
	}

	mockService.On("CreateNotification", mock.AnythingOfType("model.Notification")).Return(expectedNotification, nil)

	handler.CreateNotification(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_CreateNotification_InvalidPayload(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateNotification(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "CreateNotification")
}

func TestHandler_CreateNotification_TimeInPast(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	notificationDTO := dto.NotificationDTO{
		Text:       "Test notification",
		TelegramId: 123,
		SendAt:     time.Now().Add(-time.Hour),
	}
	body, _ := json.Marshal(notificationDTO)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateNotification(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "CreateNotification")
}

func TestHandler_CreateNotification_ServiceError(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	notificationDTO := dto.NotificationDTO{
		Text:       "Test notification",
		TelegramId: 123,
		SendAt:     time.Now().Add(time.Hour),
	}
	body, _ := json.Marshal(notificationDTO)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockService.On("CreateNotification", mock.AnythingOfType("model.Notification")).Return((*model.Notification)(nil), assert.AnError)

	handler.CreateNotification(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetNotificationStatus_Success(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	expectedStatus := &dto.NotificationStatus{
		Id:     1,
		Status: "active",
	}

	mockService.On("GetNotificationStatus", 1).Return(expectedStatus, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.GetNotificationStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetNotificationStatus_InvalidID(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	handler.GetNotificationStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "GetNotificationStatus")
}

func TestHandler_GetNotificationStatus_NotFound(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	mockService.On("GetNotificationStatus", 1).Return((*dto.NotificationStatus)(nil), assert.AnError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.GetNotificationStatus(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetAllNotifications_Success(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	expectedNotifications := []model.Notification{
		{Id: 1, Text: "Test 1", TelegramId: 123, SendAt: 1234567890, Status: "active"},
		{Id: 2, Text: "Test 2", TelegramId: 456, SendAt: 1234567891, Status: "active"},
	}

	mockService.On("GetAllNotifications").Return(expectedNotifications, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetAllNotifications(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetAllNotifications_ServiceError(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	mockService.On("GetAllNotifications").Return(nil, assert.AnError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetAllNotifications(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetMainPage_Success(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.LoadHTMLGlob("../../static/*.html")

	handler.GetMainPage(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_UpdateNotificationStatus_Success(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	mockService.On("UpdateNotificationStatus", 1, "canceled").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.UpdateNotificationStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_UpdateNotificationStatus_InvalidID(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	handler.UpdateNotificationStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "UpdateNotificationStatus")
}

func TestHandler_UpdateNotificationStatus_NotFound(t *testing.T) {
	mockService := new(MockNotifierService)
	handler := New(mockService)

	mockService.On("UpdateNotificationStatus", 1, "canceled").Return(assert.AnError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.UpdateNotificationStatus(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
