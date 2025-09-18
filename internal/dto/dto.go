package dto

import "time"

type NotificationStatus struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

type NotificationDTO struct {
	Id         int       `json:"id"`
	Text       string    `json:"text"`
	Status     string    `json:"status"`
	TelegramId int       `json:"telegram_id"`
	SendAt     time.Time `json:"send_at"`
	CreatedAt  time.Time `json:"created_at"`
}
