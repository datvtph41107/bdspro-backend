package paymentcompleted

import "time"

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
