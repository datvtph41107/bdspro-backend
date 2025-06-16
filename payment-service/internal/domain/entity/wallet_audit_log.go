package entity

import "time"

type WalletAuditLog struct {
	Id       uint32
	WalletId uint32
	OldBalance float64
	NewBalance float64
	Reason string
	ChangedBy uint32
	ChangedAt time.Time
}