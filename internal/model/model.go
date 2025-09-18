package model

import "time"

type Notification struct {
	Id         int       `json:"id"`
	Text       string    `json:"text"`
	Status     string    `json:"status"`
	TelegramId int       `json:"telegram_id"`
	SendAt     int       `json:"send_at"`
	CreatedAt  time.Time `json:"created_at"`
}
