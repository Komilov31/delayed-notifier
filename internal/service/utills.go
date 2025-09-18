package service

import (
	"encoding/json"
	"fmt"

	"github.com/Komilov31/delayed-notifier/internal/model"
	"github.com/wb-go/wbf/zlog"
)

func (s *Service) handleMessage(msg []byte, notification model.Notification) error {
	if err := json.Unmarshal(msg, &notification); err != nil {
		return fmt.Errorf("could not unmarshal notification from queue: " + err.Error())
	}

	if err := s.sender.SendToTelegram(notification.TelegramId, notification.Text); err != nil {
		return fmt.Errorf("could not send notification to Telegram: " + err.Error())
	}

	if err := s.storage.UpdateNotificationStatus(notification.Id, "completed"); err != nil {
		return fmt.Errorf("could not update notification  status in db: " + err.Error())
	}

	if err := s.cache.Set(notification.Id, "completed"); err != nil {
		return fmt.Errorf("could not update notification  status in redis: " + err.Error())
	}

	zlog.Logger.Info().Msg("succesfully handled message from queue")
	return nil
}
