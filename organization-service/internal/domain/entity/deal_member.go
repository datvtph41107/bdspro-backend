package entity

import "organization/internal/enums"

type DealMember struct {
	ID           uint64               `gorm:"primaryKey;autoIncrement"`
	DealID       uint64               `gorm:"primaryKey;not null"`
	MemberID     uint64               `gorm:"primaryKey;not null"`
	MemberType   enums.DealMemberType `gorm:"primaryKey;not null;default:10"`
	AmountCommit float64
}

func (DealMember) TableName() string {
	return "deal_members"
}
