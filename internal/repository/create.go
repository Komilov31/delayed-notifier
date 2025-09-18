package repository

import (
	"fmt"

	"github.com/Komilov31/delayed-notifier/internal/model"
)

func (r *Repository) CreateNotification(notification model.Notification) (*model.Notification, error) {
	query := `INSERT INTO notifications(text, status, telegram_id, send_at)
	VALUES($1, $2, $3, $4) RETURNING id, created_at`

	err := r.db.Master.QueryRow(
		query,
		notification.Text,
		notification.Status,
		notification.TelegramId,
		notification.SendAt,
	).Scan(&notification.Id, &notification.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not scan notification info from db: %w", err)
	}

	return &notification, nil
}
