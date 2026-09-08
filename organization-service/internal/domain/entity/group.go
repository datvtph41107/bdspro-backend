package entity

import "time"

type GroupStatus string

const (
	GroupStatusActive    GroupStatus = "active"    // Đang hoạt động
	GroupStatusPaused    GroupStatus = "paused"    // Tạm ngưng
	GroupStatusDisbanded GroupStatus = "disbanded" // Đã giải tán
)

type Group struct {
	Id        uint32
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint32
	UpdatedBy uint32

	Name        string
	Description string
	AvatarUrl   string
	Status      GroupStatus
}
