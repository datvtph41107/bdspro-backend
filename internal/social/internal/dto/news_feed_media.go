package dto

import "time"

type NewsFeedMediaDTO struct {
	ID        uint64    `json:"id"`
	Url       string    `json:"url"`
	Type      string    `json:"type"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
