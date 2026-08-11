package usage

import (
	commonaccess "common/access"
	"common/businessoperation"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

/**
 * Event ghi lại một lần usage đã được business chấp nhận.
 *
 * Đây là dữ liệu bền vững để sau này có thể kiểm tra lại Redis runtime.
 */
type Event struct {
	UsageKey string

	Subject   commonaccess.Subject
	Operation businessoperation.Code

	OperationID    string
	IdempotencyKey string
	CommandKey     string
	ReservationID  string

	Amount int64

	PeriodStart time.Time
	PeriodEnd   time.Time
	CreatedAt   time.Time
}

/**
 * NewEvent tạo usage event với key ổn định cho cùng một command.
 */
func NewEvent(
	subject commonaccess.Subject,
	operation businessoperation.Code,
	operationID string,
	idempotencyKey string,
	commandKey string,
	reservationID string,
	amount int64,
	periodStart time.Time,
	periodEnd time.Time,
	createdAt time.Time,
) Event {
	return Event{
		UsageKey:       BuildUsageKey(subject, operation, commandKey),
		Subject:        subject,
		Operation:      operation,
		OperationID:    strings.TrimSpace(operationID),
		IdempotencyKey: strings.TrimSpace(idempotencyKey),
		CommandKey:     strings.TrimSpace(commandKey),
		ReservationID:  strings.TrimSpace(reservationID),
		Amount:         amount,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		CreatedAt:      createdAt,
	}
}

/**
 * BuildUsageKey tạo key chống ghi usage trùng cho cùng subject + operation + command.
 */
func BuildUsageKey(
	subject commonaccess.Subject,
	operation businessoperation.Code,
	commandKey string,
) string {
	raw := string(subject.Type) + "\n" + subject.ID + "\n" + string(operation) + "\n" + commandKey
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (e Event) IsValid() bool {
	return e.UsageKey == BuildUsageKey(e.Subject, e.Operation, e.CommandKey) &&
		e.Subject.IsValid() &&
		e.Operation.IsValid() &&
		e.OperationID != "" &&
		e.CommandKey != "" &&
		e.ReservationID != "" &&
		e.Amount > 0 &&
		!e.PeriodStart.IsZero() &&
		e.PeriodEnd.After(e.PeriodStart) &&
		!e.CreatedAt.IsZero()
}
