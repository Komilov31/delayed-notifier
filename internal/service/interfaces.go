package service

import (
	"context"

	"github.com/Komilov31/delayed-notifier/internal/model"
)

type Storage interface {
	CreateNotification(model.Notification) (*model.Notification, error)
	DeleteNotificationById(int) error
	GetNotificationById(int) (*model.Notification, error)
	GetAllNotifications() ([]model.Notification, error)
	GetReadyNotifications() ([]model.Notification, error)
	UpdateNotificationStatus(int, string) error
}

type Cache interface {
	Get(string) (string, error)
	Set(int, interface{}) error
}

type Queue interface {
	Publish(model.Notification) error
	Consume(ctx context.Context) (<-chan []byte, error)
}

type Sender interface {
	SendToTelegram(int, string) error
}
