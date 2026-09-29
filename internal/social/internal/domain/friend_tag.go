package domain

import "time"

type FriendTag struct {
	NewsFeedID uint64    `gorm:"primaryKey" json:"news_feed_id"`
	UserID     uint64    `gorm:"primaryKey" json:"user_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (FriendTag) TableName() string {
	return "friend_tag"
}
