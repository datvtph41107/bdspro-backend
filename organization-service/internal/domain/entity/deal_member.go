package entity

import "organization/internal/enums"

type DealMember struct {
	// ID              uint64               `gorm:"primaryKey;autoIncrement"`
	DealID          uint64               `gorm:"primaryKey;not null"`
	MemberID        uint64               `gorm:"primaryKey;not null"`
	MemberType      enums.DealMemberType `gorm:"primaryKey;not null;default:10"`
	RoleID          uint64               ``
	AmountCommit    float64              `gorm:"default:0"`
	CommissionValue float64              `gorm:"default:0"`     // Giá trị hoa hồng (tiền hoặc %)
	CommissionType  enums.CommissionType `gorm:"default:10"`    // Loại hoa hồng (% hoặc VND)
	Note            string               `gorm:"type:text"`     // Ghi chú
	IsUnilateral    bool                 `gorm:"default:false"` // Đánh dấu gỡ khỏi thương vụ một cách đơn phương
}

func (DealMember) TableName() string {
	return "deal_members"
}
