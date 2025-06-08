package domain

import _models "common/models"

type Report struct {
	_models.BaseEntity
	Content    string `json:"content"`
	NewsFeedID uint64 `json:"newsFeedId"`
	ReasonID   uint64 `json:"reasonId"`
	UserID     uint64 `json:"userId"`
}

func (Report) TableName() string {
	return "report"
}
