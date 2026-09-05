package wallet

import "time"

type WalletAuditLog struct {
	Id         uint32
	WalletId   uint32
	OldBalance int64
	NewBalance int64
	Reason     string
	ChangedBy  uint32
	ChangedAt  time.Time
}
