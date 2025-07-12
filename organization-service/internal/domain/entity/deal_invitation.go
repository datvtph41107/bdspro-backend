package entity

import (
	"organization/internal/enums"
	"time"
)

type DealInvitationStatus uint32

const (
	DealInvitationStatusInvited     DealInvitationStatus = 1 // Đã mời
	DealInvitationStatusPending     DealInvitationStatus = 2 // Chờ xác nhận
	DealInvitationStatusRejected    DealInvitationStatus = 3 // Từ chối
	DealInvitationStatusAccepted    DealInvitationStatus = 4 // Đã tham gia
	DealInvitationStatusWithdrawn   DealInvitationStatus = 5 // Đã rút khỏi thương vụ
)

var DealInvitationStatusMap = map[DealInvitationStatus]string{
	DealInvitationStatusInvited:   "Đã mời",
	DealInvitationStatusPending:   "Chờ xác nhận",
	DealInvitationStatusRejected:  "Từ chối",
	DealInvitationStatusAccepted:  "Đã tham gia",
	DealInvitationStatusWithdrawn: "Đã rút khỏi thương vụ",
}

type DealInvitation struct {
	ID           uint64                `gorm:"primaryKey;autoIncrement"`
	DealID       uint64                `gorm:"not null;index"`
	InviterID    uint64                `gorm:"not null"` // Người gửi lời mời
	InviteeID    uint64                `gorm:"not null;index"` // Người được mời
	MemberType   enums.DealMemberType  `gorm:"not null;default:30"` // Loại thành viên (Member, Customer, Partner)
	Status       DealInvitationStatus  `gorm:"not null;default:1"`
	Message      string                `gorm:"type:text"` // Lời nhắn khi gửi lời mời
	AmountCommit *float64              `gorm:"default:0"` // Số tiền cam kết (nếu có)
	InvitedAt    time.Time             `gorm:"not null"`
	RespondedAt  *time.Time            // Thời gian phản hồi
	WithdrawnAt  *time.Time            // Thời gian rút khỏi
	CreatedAt    time.Time             `gorm:"not null"`
	UpdatedAt    time.Time             `gorm:"not null"`
	CreatedBy    uint64                `gorm:"not null"`
	UpdatedBy    uint64                `gorm:"not null"`
	
	// Relations
	Deal         *Deal                 `gorm:"foreignKey:DealID;references:ID"`
	Inviter      interface{}           `gorm:"-"` // Sẽ được populate từ User service
	Invitee      interface{}           `gorm:"-"` // Sẽ được populate từ User service
}

func (DealInvitation) TableName() string {
	return "deal_invitations"
}

// Helper methods
func (di *DealInvitation) IsPending() bool {
	return di.Status == DealInvitationStatusPending
}

func (di *DealInvitation) IsAccepted() bool {
	return di.Status == DealInvitationStatusAccepted
}

func (di *DealInvitation) IsRejected() bool {
	return di.Status == DealInvitationStatusRejected
}

func (di *DealInvitation) IsWithdrawn() bool {
	return di.Status == DealInvitationStatusWithdrawn
}

func (di *DealInvitation) CanResend() bool {
	return di.Status == DealInvitationStatusRejected || di.Status == DealInvitationStatusWithdrawn
} 