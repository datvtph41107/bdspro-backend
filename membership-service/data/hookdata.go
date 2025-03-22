package data

import "time"

type NotiComplete struct {
	ID        uint64
	ProfileId uint64
	OrderId   uint64
	Type      int8
	TargetId  uint64
	PaidAt    *time.Time
	Signature string
}
