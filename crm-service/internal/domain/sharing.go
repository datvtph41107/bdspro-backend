package domain

import (
	_models "common/models"
	"crm/internal/enums"

	"github.com/lib/pq"
)

type SharingEntity struct {
	_models.BaseEntity
	ContactID    uint64         `gorm:"primaryKey" json:"contactId"`
	ReceiverID   uint64         `gorm:"primaryKey" json:"receiverId"`
	ReceiverType enums.EOwnerOf `gorm:"primaryKey" json:"receiverType"`
	Permissions  pq.Int32Array  `gorm:"type:integer[]" json:"permissions"`
	Note         string         `json:"note"`
}

func (s *SharingEntity) TableName() string {
	return "sharing_access"
}