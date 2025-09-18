package service

import (
	"context"
	"time"

	"github.com/Komilov31/delayed-notifier/internal/model"
	"github.com/wb-go/wbf/zlog"
)

const (
	workersNum = 3
)

func (s *Service) PublishReadyNotifications(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			notifications, err := s.storage.GetReadyNotifications()
			if err != nil {
				return err
			}

			for _, notif := range notifications {
				if (int64(notif.SendAt) - time.Now().UnixMilli()) <= time.Minute.Milliseconds()/2 {
					if err := s.queue.Publish(notif); err != nil {
						return err
					}
					zlog.Logger.Info().Msg("successfully published message")
				}
			}
		}

		time.Sleep(time.Minute)
	}
}

func (s *Service) ConsumeMessages(ctx context.Context) error {
	messages, err := s.queue.Consume(ctx)
	if err != nil {
		return err
	}

	for i := range workersNum {
		go func(i int) {
			zlog.Logger.Info().Msgf("consumer with index %d started", i)
			for msg := range messages {
				var notification model.Notification
				if err := s.handleMessage(msg, notification); err != nil {
					zlog.Logger.Error().Msg(err.Error())
					continue
				}
			}
		}(i)
	}

	return nil
}
