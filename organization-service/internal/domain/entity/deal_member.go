package entity

import (
	_models "common/models"
	"organization/internal/enums"
	"time"
)

type DealMemberStatus uint32

const (
	DealMemberStatusInvited   DealMemberStatus = 10 // Đã mời
	DealMemberStatusAccepted  DealMemberStatus = 20 // Đã tham gia
	DealMemberStatusRejected  DealMemberStatus = 30 // Từ chối
	DealMemberStatusWithdrawn DealMemberStatus = 40 // Đã rút khỏi thương vụ
)

var DealMemberStatusMap = map[DealMemberStatus]string{
	DealMemberStatusInvited:   "Đã mời",
	DealMemberStatusRejected:  "Từ chối",
	DealMemberStatusAccepted:  "Đã tham gia",
	DealMemberStatusWithdrawn: "Đã rút khỏi thương vụ",
}

type DealMember struct {
	_models.BaseEntity
	// ID       uint64 `gorm:"primaryKey;autoIncrement"`
	DealID   uint64 `gorm:"primaryKey;not null"`
	MemberID uint64 `gorm:"primaryKey;index;not null"`
	// MemberType      enums.DealMemberType `gorm:"primaryKey;not null;default:10"`
	RoleID          uint64               `gorm:"column:role_id;not null" json:"role_id"`
	AmountCommit    float64              `gorm:"default:0"`
	CommissionValue float64              `gorm:"default:0"`     // Giá trị hoa hồng (tiền hoặc %)
	CommissionType  enums.CommissionType `gorm:"default:10"`    // Loại hoa hồng (% hoặc VND)
	Note            string               `gorm:"type:text"`     // Ghi chú
	IsUnilateral    bool                 `gorm:"default:false"` // Đánh dấu gỡ khỏi thương vụ một cách đơn phương
	DoneInvestment  bool                 `gorm:"default:false"` // Đánh dấu đã đủ vốn góp

	RoleKey enums.RoleKey `gorm:"column:role_key;not null;default:430" json:"role_key"`

	// Thêm các trường từ DealInvitation
	Status      DealMemberStatus `gorm:"not null;default:10"` // Trạng thái
	Message     string           `gorm:"type:text"`           // Lời nhắn khi gửi lời mời
	InvitedAt   time.Time        `gorm:""`                    // Thời gian mời
	RespondedAt *time.Time       // Thời gian phản hồi
	WithdrawnAt *time.Time       // Thời gian rút khỏi
	InviterID   uint64           `gorm:""`                           // Người gửi lời mời
	ColorId     *uint32          `gorm:"index"`                      // Màu sắc
	IsOwner     bool             `gorm:"type:boolean;default:false"` // Đánh dấu là chủ thương vụ

	// Relations
	Deal  *Deal  `gorm:"foreignKey:DealID;references:ID"`
	Color *Color `gorm:"foreignKey:ColorId;references:ID"`
	// Role    *OrganizationRole `gorm:"foreignKey:RoleId;references:ID"`
	Inviter interface{} `gorm:"-"` // Sẽ được populate từ User service
	Member  interface{} `gorm:"-"` // Sẽ được populate từ User service
}

func (DealMember) TableName() string {
	return "deal_members"
}

// Helper methods
func (dm *DealMember) IsInvited() bool {
	return dm.Status == DealMemberStatusInvited
}

func (dm *DealMember) IsAccepted() bool {
	return dm.Status == DealMemberStatusAccepted
}

func (dm *DealMember) IsRejected() bool {
	return dm.Status == DealMemberStatusRejected
}

func (dm *DealMember) IsWithdrawn() bool {
	return dm.Status == DealMemberStatusWithdrawn
}

func (dm *DealMember) CanResend() bool {
	return dm.Status == DealMemberStatusRejected || dm.Status == DealMemberStatusWithdrawn
}
