package domain

type NewsFeedOfUser struct {
	ID     uint64 `gorm:"primaryKey;column:id" json:"id"`
	PostID uint64 `gorm:"column:post_id;not null" json:"postId"`
	UserID uint64 `gorm:"column:user_id;not null" json:"userId"`
}

func (NewsFeedOfUser) TableName() string {
	return "news_feed_of_user"
}
