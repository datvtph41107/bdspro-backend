package domain

type TagFriend struct {
	NewsFeedID uint64 `json:"news_feed_id"`
	UserID     uint64 `json:"user_id"`
}

func (TagFriend) TableName() string {
	return "tag_friend"
}
