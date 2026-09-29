package domain

import "time"

// ContactTagEntity là bảng nối giữa contact và tag.
type ContactTagEntity struct {
	ContactID uint64    `gorm:"column:contact_id;primaryKey" json:"contactId"`
	TagID     uint64    `gorm:"column:tag_id;primaryKey" json:"tagId"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (ContactTagEntity) TableName() string {
	return "tb_contact_tag_relation"
}