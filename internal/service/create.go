package service

import (
	"github.com/Komilov31/delayed-notifier/internal/model"
)

func (s *Service) CreateNotification(notification model.Notification) (*model.Notification, error) {
	notification.Status = "active"
	notif, err := s.storage.CreateNotification(notification)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(notif.Id, notif.Status); err != nil {
		return nil, err
	}

	return notif, nil
}
