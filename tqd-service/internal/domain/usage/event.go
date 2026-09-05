package usage

import (
	commonmetering "common/metering"
	"common/operation"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
	"tqd/internal/access"
)

// CommercialEvidence records the commercial decision that justified accepted usage.
// Empty evidence is allowed only for legacy rows during migration; new catalog-backed
// Report acceptance should populate it from user-service GetAccess.
type CommercialEvidence struct {
	SubscriptionID uint64
	PlanCode       string
	PlanVersion    string
	PolicyVersion  string
	LimitSnapshot  int64
}

func (e CommercialEvidence) IsEmpty() bool {
	return e.SubscriptionID == 0 &&
		strings.TrimSpace(e.PlanCode) == "" &&
		strings.TrimSpace(e.PlanVersion) == "" &&
		strings.TrimSpace(e.PolicyVersion) == "" &&
		e.LimitSnapshot == 0
}

func (e CommercialEvidence) IsValid() bool {
	if e.IsEmpty() {
		return true
	}
	return e.SubscriptionID > 0 &&
		strings.TrimSpace(e.PlanCode) != "" &&
		strings.TrimSpace(e.PlanVersion) != "" &&
		strings.TrimSpace(e.PolicyVersion) != "" &&
		e.LimitSnapshot >= 0
}

// NewEventInput contains the durable facts of one accepted consumption.
// ReservationID is runtime-admission evidence, not Usage identity.
type NewEventInput struct {
	Subject   access.Subject
	Operation operation.Code
	MeterCode commonmetering.Code

	OperationID    string
	IdempotencyKey string
	CommandKey     string
	ReservationID  string

	Amount int64

	PeriodStart time.Time
	PeriodEnd   time.Time
	CreatedAt   time.Time

	Commercial CommercialEvidence
}

/**
 * Event ghi lại một lần usage đã được business chấp nhận.
 *
 * Đây là dữ liệu bền vững để audit commercial decision và repair runtime projection.
 * ReservationID có thể rỗng khi admission không dùng Redis reservation.
 */
type Event struct {
	UsageKey string

	Subject   access.Subject
	Operation operation.Code
	MeterCode commonmetering.Code

	OperationID    string
	IdempotencyKey string
	CommandKey     string
	ReservationID  string

	Amount int64

	PeriodStart time.Time
	PeriodEnd   time.Time
	CreatedAt   time.Time

	Commercial CommercialEvidence
}

// NewEvent creates durable accepted-usage evidence for one logical command.
func NewEvent(input NewEventInput) Event {
	input.CommandKey = strings.TrimSpace(input.CommandKey)
	return Event{
		UsageKey:       BuildUsageKey(input.Subject, input.Operation, input.MeterCode, input.CommandKey),
		Subject:        input.Subject,
		Operation:      input.Operation,
		MeterCode:      input.MeterCode,
		OperationID:    strings.TrimSpace(input.OperationID),
		IdempotencyKey: strings.TrimSpace(input.IdempotencyKey),
		CommandKey:     input.CommandKey,
		ReservationID:  strings.TrimSpace(input.ReservationID),
		Amount:         input.Amount,
		PeriodStart:    input.PeriodStart,
		PeriodEnd:      input.PeriodEnd,
		CreatedAt:      input.CreatedAt,
		Commercial: CommercialEvidence{
			SubscriptionID: input.Commercial.SubscriptionID,
			PlanCode:       strings.TrimSpace(input.Commercial.PlanCode),
			PlanVersion:    strings.TrimSpace(input.Commercial.PlanVersion),
			PolicyVersion:  strings.TrimSpace(input.Commercial.PolicyVersion),
			LimitSnapshot:  input.Commercial.LimitSnapshot,
		},
	}
}

/**
 * BuildUsageKey chống ghi usage trùng theo subject + operation + meter + command.
 */
func BuildUsageKey(
	subject access.Subject,
	operationCode operation.Code,
	meterCode commonmetering.Code,
	commandKey string,
) string {
	raw := string(subject.Type) + "\n" + subject.ID + "\n" + string(operationCode) + "\n" + string(meterCode) + "\n" + commandKey
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (e Event) IsValid() bool {
	return e.UsageKey == BuildUsageKey(e.Subject, e.Operation, e.MeterCode, e.CommandKey) &&
		e.Subject.IsValid() &&
		e.Operation.IsValid() &&
		e.MeterCode.IsValid() &&
		e.OperationID != "" &&
		e.CommandKey != "" &&
		e.Amount > 0 &&
		!e.PeriodStart.IsZero() &&
		e.PeriodEnd.After(e.PeriodStart) &&
		!e.CreatedAt.IsZero() &&
		e.Commercial.IsValid()
}
