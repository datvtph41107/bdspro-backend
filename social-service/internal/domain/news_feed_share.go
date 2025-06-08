package domain

import (
	_models "common/models"
	"social/internal/enums"
)

type NewsFeedShare struct {
	_models.BaseEntity
	NewsFeedID uint64          `json:"news_feed_id"`
	TargetID   *uint64         `json:"target_id"`
	UserID     *uint64         `json:"user_id"`
	ShareType  enums.ShareType `json:"share_type"`
	// todo: cân nhắc tổ chức hoặc nhóm có thể share
	// thêm userType: 10: user, 20: organization, 30: group
}

func (n *NewsFeedShare) TableName() string {
	return "news_feed_shares"
}
