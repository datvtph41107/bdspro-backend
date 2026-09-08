package base_dto

import "time"

type NotificationDTO struct {
	Title      string     `json:"title"`
	Message    string     `json:"message"`
	Type       int        `json:"type"`
	UserID     uint64     `json:"user_id"`
	VisibleAt  *time.Time `json:"visible_at"`
	AttachData []string   `json:"attach_data"`
}

type NotificationBatchDTO struct {
	Datas []NotificationDTO `json:"datas"`
}
