package domain

import (
	_models "common/models"
	"social/internal/enums"
	"time"
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
	Rank           uint64           `gorm:"column:rank_feed" json:"rank"`
	NewsFeedMedias []NewsFeedMedia  `gorm:"foreignKey:NewsFeedID" json:"news_feed_medias"`
	FriendTags     []FriendTag      `gorm:"foreignKey:NewsFeedID;references:ID" json:"friend_tag"`
	FriendTagIds   []uint64         `gorm:"-" json:"friend_tag_ids"`
	LikedAt        *time.Time       `gorm:"-" json:"likedAt"`
	DisLike        *bool            `gorm:"-" json:"disLike"`
	OwnerOf        enums.OwnerOf    `gorm:"column:owner_of" json:"ownerOf"`
	OwnerID        *uint64          `gorm:"column:owner_id" json:"ownerId"`

	RemovedAt *time.Time `gorm:"column:removed_at" json:"removedAt"`
	// Author         *dto.UserProfile `gorm:"-" json:"author"`
}

type ReelEntity struct {
	NewsFeed
	ReelUrl string `gorm:"column:reel_url" json:"reel_url"`
}
