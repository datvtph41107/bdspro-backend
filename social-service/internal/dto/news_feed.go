package dto

import (
	_dto "common/dto"
	"encoding/json"
	"time"
)

type PostSearch struct {
	_dto.Pagable
	Text string `json:"text" form:"text"`
}

type NewsFeedGlobalSearch struct {
	_dto.Pagable
	Text   string `json:"text" form:"text"`
	IsReel *bool  `json:"isReel" form:"isReel"`
	IsPost *bool  `json:"isPost" form:"isPost"`
}

type NewsFeedPublic struct {
	ID             uint64          `db:"id" json:"id"`
	Title          string          `db:"title" json:"title"`
	Content        string          `db:"content" json:"content"`
	Image          string          `db:"image" json:"image"`
	Link           string          `db:"link" json:"link"`
	Visibility     int             `db:"visibility" json:"visibility"`
	ParentID       *uint64         `db:"parent_id" json:"parent_id"`
	PostID         *uint64         `db:"post_id" json:"post_id"`
	NumLike        int             `db:"num_like" json:"num_like"`
	NumComment     int             `db:"num_comment" json:"num_comment"`
	NumShare       int             `db:"num_share" json:"num_share"`
	NumView        int             `db:"num_view" json:"num_view"`
	IsReel         bool            `db:"is_reel" json:"is_reel"`
	CreatedAt      time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at" json:"updated_at"`
	CreatedBy      *uint64         `db:"created_by" json:"created_by"`
	NewsFeedMedias json.RawMessage `db:"news_feed_medias" json:"news_feed_medias"`
	FriendTagIds   json.RawMessage `db:"friend_tag_ids" json:"friend_tag_ids"`
	Author         *UserProfile    `gorm:"-" db:"-" json:"author"`
	Post           *BdsproPost     `gorm:"-" db:"-" json:"post"`
	LikedAt        *time.Time      `gorm:"liked_at" db:"liked_at" json:"liked_at"`
	DisLike        *bool           `gorm:"dis_like" db:"dis_like" json:"is_like"`
}
