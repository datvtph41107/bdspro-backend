package domain

// ─────────────────────────────────────────────────────────────────────────────
// hub/internal/domain/applink.go  —  Table declaration
// ─────────────────────────────────────────────────────────────────────────────

import (
	"time"

	"hub/internal/enums"
)

// Applink là bảng ánh xạ deep link code → refId + action.
//
// Flow:
//  1. QR share gọi CreateApplink → nhận base64(code) → embed vào QR URL
//  2. User scan → web mở URL → gọi GetApplinkByCode(base64) → nhận refId + action
//  3. Web redirect vào app://...
type Applink struct {
	// Primary key — uint64 nhất quán với các bảng khác trong project
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;not null"`

	// code: bigint random [10^6, 10^10], unique index, query field chính.
	// Không expose số thô ra ngoài — client chỉ thấy base64(code).
	Code uint64 `gorm:"column:code;type:bigint;uniqueIndex:idx_applink_code;not null"`

	// ref_id: ID của entity tương ứng (user_id, property_id, room_id, ...)
	RefID uint64 `gorm:"column:ref_id;not null;index:idx_applink_ref_action"`

	// action: XXYY enum — service(XX) + event(YY)
	Action enums.ApplinkAction `gorm:"column:action;type:int;not null;index:idx_applink_ref_action"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;not null"`
}

// TableName khai báo tên bảng — GORM convention override
func (Applink) TableName() string {
	return "applink"
}
