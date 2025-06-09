package domain

import _models "common/models"

type NewsFeedMedia struct {
	_models.BaseEntity
	NewsFeedID  *uint64 `json:"newsFeedId"`
	Url         string  `json:"url"`
	Type        string  `json:"type"`
	OrderNumber int     `json:"orderNumber"`
}

func (m *NewsFeedMedia) TableName() string {
	return "news_feed_medias"
}
