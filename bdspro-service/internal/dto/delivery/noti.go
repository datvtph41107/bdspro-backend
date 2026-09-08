package delivery_dto

import (
	"time"
)

type NotiDTO struct {
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Type      int        `json:"type"`
	UserID    uint64     `json:"user_id"`
	VisibleAt *time.Time `json:"visible_at"`
}

type NotiBatchDTO struct {
	Datas []NotiDTO `json:"datas"`
}
