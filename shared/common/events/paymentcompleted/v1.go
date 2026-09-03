package paymentcompleted

import (
	"strings"
	"time"
)

const (
	EventTypeV1   = "payment.completed.v1"
	SchemaVersion = 1
)

// V1 is a past-fact integration contract. Consumers may react independently;
// it is never the command path that activates a Subscription.
type V1 struct {
	EventID          string    `json:"eventId"`
	OrderID          uint64    `json:"orderId"`
	SubjectKind      string    `json:"subjectKind"`
	SubjectID        string    `json:"subjectId"`
	ProductCode      string    `json:"productCode"`
	PlanCode         string    `json:"planCode"`
	PlanVersionID    uint64    `json:"planVersionId"`
	PlanVersion      string    `json:"planVersion"`
	Currency         string    `json:"currency"`
	AmountMinor      int64     `json:"amountMinor"`
	FundsConfirmedAt time.Time `json:"fundsConfirmedAt"`
	CompletedAt      time.Time `json:"completedAt"`
}

// IsValid protects the consumer Inbox from accepting a syntactically valid
// JSON envelope that cannot describe a Payment-owned durable fact. It is
// intentionally independent from transport authentication and provider
// verification; those checks belong to their respective adapters.
func (e V1) IsValid() bool {
	if strings.TrimSpace(e.EventID) == "" || e.OrderID == 0 || strings.TrimSpace(e.SubjectID) == "" ||
		strings.TrimSpace(e.ProductCode) == "" || strings.TrimSpace(e.PlanCode) == "" ||
		e.PlanVersionID == 0 || strings.TrimSpace(e.PlanVersion) == "" ||
		strings.TrimSpace(e.Currency) == "" || e.AmountMinor <= 0 ||
		e.FundsConfirmedAt.IsZero() || e.CompletedAt.IsZero() {
		return false
	}
	switch strings.TrimSpace(e.SubjectKind) {
	case "profile", "organization":
		return true
	default:
		return false
	}
}
