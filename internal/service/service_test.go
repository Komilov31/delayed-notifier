package service

import (
	"context"
	"testing"

	"github.com/Komilov31/delayed-notifier/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock implementation of Storage
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreateNotification(notification model.Notification) (*model.Notification, error) {
	args := m.Called(notification)
	return args.Get(0).(*model.Notification), args.Error(1)
}

func (m *MockStorage) DeleteNotificationById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStorage) GetNotificationById(id int) (*model.Notification, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Notification), args.Error(1)
}

func (m *MockStorage) GetAllNotifications() ([]model.Notification, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Notification), args.Error(1)
}

func (m *MockStorage) GetReadyNotifications() ([]model.Notification, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Notification), args.Error(1)
}

func (m *MockStorage) UpdateNotificationStatus(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

// MockCache is a mock implementation of Cache
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Set(key int, value interface{}) error {
	args := m.Called(key, value)
	return args.Error(0)
}

// MockQueue is a mock implementation of Queue
type MockQueue struct {
	mock.Mock
}

func (m *MockQueue) Publish(notification model.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockQueue) Consume(ctx context.Context) (<-chan []byte, error) {
	args := m.Called(ctx)
	return args.Get(0).(<-chan []byte), args.Error(1)
}

// MockSender is a mock implementation of Sender
type MockSender struct {
	mock.Mock
}

func (m *MockSender) SendToTelegram(id int, text string) error {
	args := m.Called(id, text)
	return args.Error(0)
}

func TestService_CreateNotification_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	notification := model.Notification{
		Text:       "Test notification",
		TelegramId: 123,
		SendAt:     1234567890,
	}

	expectedNotification := &model.Notification{
		Id:         1,
		Text:       notification.Text,
		TelegramId: notification.TelegramId,
		SendAt:     notification.SendAt,
		Status:     "active",
	}

	mockStorage.On("CreateNotification", mock.MatchedBy(func(n model.Notification) bool {
		return n.Text == notification.Text && n.Status == "active"
	})).Return(expectedNotification, nil)

	mockCache.On("Set", expectedNotification.Id, expectedNotification.Status).Return(nil)

	result, err := service.CreateNotification(notification)

	assert.NoError(t, err)
	assert.Equal(t, expectedNotification, result)
	mockStorage.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestService_CreateNotification_StorageError(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	notification := model.Notification{
		Text:       "Test notification",
		TelegramId: 123,
		SendAt:     1234567890,
	}

	mockStorage.On("CreateNotification", mock.AnythingOfType("model.Notification")).Return((*model.Notification)(nil), assert.AnError)

	result, err := service.CreateNotification(notification)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStorage.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "Set")
}

func TestService_CreateNotification_CacheError(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	notification := model.Notification{
		Text:       "Test notification",
		TelegramId: 123,
		SendAt:     1234567890,
	}

	expectedNotification := &model.Notification{
		Id:         1,
		Text:       notification.Text,
		TelegramId: notification.TelegramId,
		SendAt:     notification.SendAt,
		Status:     "active",
	}

	mockStorage.On("CreateNotification", mock.AnythingOfType("model.Notification")).Return(expectedNotification, nil)

	mockCache.On("Set", expectedNotification.Id, expectedNotification.Status).Return(assert.AnError)

	result, err := service.CreateNotification(notification)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStorage.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestService_GetAllNotifications_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	expectedNotifications := []model.Notification{
		{Id: 1, Text: "Test 1", TelegramId: 123, SendAt: 1234567890, Status: "active"},
		{Id: 2, Text: "Test 2", TelegramId: 456, SendAt: 1234567891, Status: "active"},
	}

	mockStorage.On("GetAllNotifications").Return(expectedNotifications, nil)

	result, err := service.GetAllNotifications()

	assert.NoError(t, err)
	assert.Equal(t, expectedNotifications, result)
	mockStorage.AssertExpectations(t)
}

func TestService_GetAllNotifications_StorageError(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	mockStorage.On("GetAllNotifications").Return(nil, assert.AnError)

	result, err := service.GetAllNotifications()

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStorage.AssertExpectations(t)
}

func TestService_GetNotificationStatus_CacheHit(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	mockCache.On("Get", "1").Return("active", nil)

	result, err := service.GetNotificationStatus(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.Id)
	assert.Equal(t, "active", result.Status)
	mockCache.AssertExpectations(t)
	mockStorage.AssertNotCalled(t, "GetNotificationById")
}

func TestService_GetNotificationStatus_CacheMiss(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	notification := &model.Notification{
		Id:         1,
		Text:       "Test",
		TelegramId: 123,
		SendAt:     1234567890,
		Status:     "active",
	}

	mockCache.On("Get", "1").Return("", redis.Nil)
	mockStorage.On("GetNotificationById", 1).Return(notification, nil)

	result, err := service.GetNotificationStatus(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.Id)
	assert.Equal(t, "active", result.Status)
	mockCache.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestService_GetNotificationStatus_CacheError(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	mockCache.On("Get", "1").Return("", assert.AnError)

	result, err := service.GetNotificationStatus(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockCache.AssertExpectations(t)
	mockStorage.AssertNotCalled(t, "GetNotificationById")
}

func TestService_GetNotificationStatus_StorageError(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	mockCache.On("Get", "1").Return("", redis.Nil)
	mockStorage.On("GetNotificationById", 1).Return((*model.Notification)(nil), assert.AnError)

	result, err := service.GetNotificationStatus(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockCache.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestService_UpdateNotificationStatus_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	mockCache.On("Set", 1, "canceled").Return(nil)
	mockStorage.On("UpdateNotificationStatus", 1, "canceled").Return(nil)

	err := service.UpdateNotificationStatus(1, "canceled")

	assert.NoError(t, err)
	mockCache.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestService_UpdateNotificationStatus_CacheError(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	mockCache.On("Set", 1, "canceled").Return(assert.AnError)

	err := service.UpdateNotificationStatus(1, "canceled")

	assert.Error(t, err)
	mockCache.AssertExpectations(t)
	mockStorage.AssertNotCalled(t, "UpdateNotificationStatus")
}

func TestService_UpdateNotificationStatus_StorageError(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	mockQueue := new(MockQueue)
	mockSender := new(MockSender)
	service := New(mockStorage, mockCache, mockQueue, mockSender)

	mockCache.On("Set", 1, "canceled").Return(nil)
	mockStorage.On("UpdateNotificationStatus", 1, "canceled").Return(assert.AnError)

	err := service.UpdateNotificationStatus(1, "canceled")

	assert.Error(t, err)
	mockCache.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}
