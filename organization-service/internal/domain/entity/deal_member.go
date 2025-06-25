package entity

import "organization/internal/enums"

type DealMember struct {
	DealID     uint64               `gorm:"primaryKey;not null"`
	MemberID   uint64               `gorm:"primaryKey;not null"`
	MemberType enums.DealMemberType `gorm:"primaryKey;not null;default:10"`
}

func (DealMember) TableName() string {
	return "deal_members"
}
