package domain

import (
	_models "common/models"
	shared_enum "pb/enums"

	"github.com/lib/pq"
)

type HistoryEntity struct {
	_models.BaseEntity
	TargetId   uint64                     `json:"targetId"`
	TargetType shared_enum.ETargetHistory `json:"targetType"`
	ActionType shared_enum.EHistory       `json:"actionType"`
	Title      string                     `json:"title,omitempty"`
	Note       pq.StringArray             `gorm:"type:text[]" json:"note,omitempty"`
	PreStage   string                     `json:"preStage,omitempty"`
	AfterStage string                     `json:"afterStage,omitempty"`
	OwnerID    *uint64                    `json:"ownerId"`
	OwnerType  shared_enum.EOwnerType     `json:"ownerType"`
	IsInternal bool                       `json:"isInternal"`
}

func (HistoryEntity) TableName() string {
	return "noti_histories"
}
