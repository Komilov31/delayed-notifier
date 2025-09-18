package repository

import (
	"database/sql"
	"fmt"

	"github.com/Komilov31/delayed-notifier/internal/model"
)

func (r *Repository) GetNotificationById(id int) (*model.Notification, error) {
	query := "SELECT * FROM notifications WHERE id = $1"

	var notification model.Notification
	err := r.db.Master.QueryRow(query, id).Scan(
		&notification.Id,
		&notification.Text,
		&notification.Status,
		&notification.TelegramId,
		&notification.SendAt,
		&notification.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoSuchNotification
		}
		return nil, fmt.Errorf("could not get notification from db: %w", err)
	}

	return &notification, nil
}

func (r *Repository) GetAllNotifications() ([]model.Notification, error) {
	query := "SELECT * FROM notifications"

	rows, err := r.db.Master.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not get all notifications from db: %w", err)
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		err := rows.Scan(
			&notification.Id,
			&notification.Text,
			&notification.Status,
			&notification.TelegramId,
			&notification.SendAt,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan row to model: %w", err)
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *Repository) GetReadyNotifications() ([]model.Notification, error) {
	query := `SELECT *
	FROM notifications
	WHERE (send_at - (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT) < 30000 
	AND status='active';
	`

	rows, err := r.db.Master.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not get all notifications from db: %w", err)
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		err := rows.Scan(
			&notification.Id,
			&notification.Text,
			&notification.Status,
			&notification.TelegramId,
			&notification.SendAt,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan row to model: %w", err)
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}
