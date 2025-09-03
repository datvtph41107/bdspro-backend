package domain

type NewsFeedOfGroup struct {
	ID      uint64 `gorm:"primaryKey;column:id" json:"id"`
	PostID  uint64 `gorm:"column:post_id;not null" json:"postId"`
	GroupID uint64 `gorm:"column:group_id;not null" json:"groupId"`
}

func (NewsFeedOfGroup) TableName() string {
	return "news_feed_of_group"
}
