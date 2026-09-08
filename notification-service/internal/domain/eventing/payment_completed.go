package eventing

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrEventIdentityConflict = errors.New("integration event identity conflict")

// PaymentCompleted is the Notification-owned normalized integration fact.
// It deliberately contains no RabbitMQ or protobuf types.
type PaymentCompleted struct {
	EventID          string
	OrderID          uint64
	SubjectKind      string
	SubjectID        string
	ProductCode      string
	PlanCode         string
	PlanVersionID    uint64
	PlanVersion      string
	Currency         string
	AmountMinor      int64
	FundsConfirmedAt time.Time
	CompletedAt      time.Time
}

// PayloadHash binds a durable EventID to the normalized business fact. A
// redelivery with the same EventID and a different hash is not a replay: it is
// an integration-contract conflict that must not be silently accepted.
func (f PaymentCompleted) PayloadHash() string {
	canonical := fmt.Sprintf(
		"%s\x1f%d\x1f%s\x1f%s\x1f%s\x1f%s\x1f%d\x1f%s\x1f%s\x1f%d\x1f%s\x1f%s",
		strings.TrimSpace(f.EventID),
		f.OrderID,
		strings.TrimSpace(f.SubjectKind),
		strings.TrimSpace(f.SubjectID),
		strings.TrimSpace(f.ProductCode),
		strings.TrimSpace(f.PlanCode),
		f.PlanVersionID,
		strings.TrimSpace(f.PlanVersion),
		strings.ToUpper(strings.TrimSpace(f.Currency)),
		f.AmountMinor,
		f.FundsConfirmedAt.UTC().Format(time.RFC3339Nano),
		f.CompletedAt.UTC().Format(time.RFC3339Nano),
	)
	sum := sha256.Sum256([]byte(canonical))
	return hex.EncodeToString(sum[:])
}
