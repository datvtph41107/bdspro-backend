package entity

import "time"

type GroupNotification struct {
	Id        uint32
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint32
	UpdatedBy uint32

	GroupId uint32
	Type    string
	Content string
}
