package domain

import (
	_models "common/models"
)

type AppInvitedEntity struct {
	_models.BaseEntity
	ContactID *uint64        `gorm:"column:contact_id" json:"contactId"`
	Content   string         `gorm:"column:content" json:"content"`
	Link      string         `gorm:"column:link" json:"link"`
	Channel   string         `gorm:"column:channel" json:"channel"`
	Contact   *ContactEntity `gorm:"foreignKey:ContactID;references:ID" json:"contact"`
}

func (AppInvitedEntity) TableName() string {
	return "app_invited"
}