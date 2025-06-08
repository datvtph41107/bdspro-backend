package domain

import _models "common/models"

type NewsFeedMedia struct {
	_models.BaseEntity
	NewsFeedID  *uint64 `json:"news_feed_id"`
	Url         string  `json:"url"`
	Type        string  `json:"type"`
	OrderNumber int     `json:"order_number"`
}

func (m *NewsFeedMedia) TableName() string {
	return "news_feed_medias"
}
