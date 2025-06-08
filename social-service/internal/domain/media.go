package domain

import _models "common/models"

type Media struct {
	_models.BaseEntity
	PostID      *uint64 `json:"post_id"`
	Link        string  `json:"link"`
	Type        string  `json:"type"`
	OrderNumber int     `json:"order_number"`
}
