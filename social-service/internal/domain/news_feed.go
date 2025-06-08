package domain

import (
	_models "common/models"
	"social/internal/enums"
)

type NewsFeed struct {
	_models.BaseEntity
	Title          string           `json:"title"`
	Content        string           `json:"content"`
	Image          string           `json:"image"`
	Link           string           `json:"link"`
	Visibility     enums.Visibility `gorm:"default:10" json:"visibility"`
	ParentID       *uint64          `json:"parent_id"`
	PostID         *uint64          `json:"post_id"`
	NumLike        int              `json:"num_like"`
	NumComment     int              `json:"num_comment"`
	NumShare       int              `json:"num_share"`
	NumView        int              `json:"num_view"`
	IsReel         bool             `json:"is_reel"`
	NewsFeedMedias []*NewsFeedMedia `gorm:"-" json:"news_feed_medias"`
	FriendTagIds   []uint64         `gorm:"-" json:"friend_tag_ids"`
}
