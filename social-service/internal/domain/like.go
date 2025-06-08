package domain

import "time"

type Like struct {
	NewsFeedID uint64    `json:"newsFeedId"`
	UserID     uint64    `json:"userId"`
	Type       uint32    `json:"likeType"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (Like) TableName() string {
	return "like"
}
