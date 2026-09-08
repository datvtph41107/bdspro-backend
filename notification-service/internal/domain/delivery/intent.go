package delivery

import "time"

type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusRetry   Status = "retry"
	StatusSent    Status = "sent"
	StatusFailed  Status = "failed"
	StatusUnknown Status = "unknown"
)

type Intent struct {
	ID           uint64
	EventID      string
	SubjectKind  string
	SubjectID    string
	Channel      string
	TemplateCode string
	Payload      string
	Status       Status
	AttemptCount int
	AvailableAt  time.Time
	LockedBy     string
	ClaimVersion uint64
	LeaseUntil   *time.Time
	LastError    string
	SentAt       *time.Time
}
