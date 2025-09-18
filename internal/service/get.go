package service

import (
	"fmt"
	"strconv"

	"github.com/Komilov31/delayed-notifier/internal/dto"
	"github.com/Komilov31/delayed-notifier/internal/model"
	"github.com/go-redis/redis/v8"
)

func (s *Service) GetAllNotifications() ([]model.Notification, error) {
	return s.storage.GetAllNotifications()
}

func (s *Service) GetNotificationStatus(id int) (*dto.NotificationStatus, error) {
	statusString, err := s.cache.Get(strconv.Itoa(id))
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("could not get notif status from redis: " + err.Error())
	}

	if err == redis.Nil {
		notification, err := s.storage.GetNotificationById(id)
		if err != nil {
			return nil, err
		}

		statusString = notification.Status
	}

	var status dto.NotificationStatus
	status.Id = id
	status.Status = statusString

	return &status, nil
}
