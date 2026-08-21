package domain

import (
	"encoding/hex"
	"strings"
	"time"
)

type SubjectKind string

const (
	SubjectProfile      SubjectKind = "profile"
	SubjectOrganization SubjectKind = "organization"
)

type Subject struct {
	Kind SubjectKind
	ID   string
}

func (s Subject) IsValid() bool {
	return (s.Kind == SubjectProfile || s.Kind == SubjectOrganization) && strings.TrimSpace(s.ID) != ""
}

type CurrencyCode string

const CurrencyVND CurrencyCode = "VND"

type Money struct {
	Currency    CurrencyCode
	AmountMinor int64
}

func (m Money) IsValid() bool {
	return strings.TrimSpace(string(m.Currency)) != "" && m.AmountMinor > 0
}

type CommercialTerms struct {
	ProductCode          string
	PlanCode             string
	PlanVersionID        uint64
	PlanVersion          string
	TierRank             int32
	SubscriptionTermDays int32
	TermsChecksum        string
	Price                Money
}

func (t CommercialTerms) IsValid() bool {
	checksum := strings.TrimSpace(t.TermsChecksum)
	if len(checksum) != 64 {
		return false
	}
	if _, err := hex.DecodeString(checksum); err != nil {
		return false
	}
	return strings.TrimSpace(t.ProductCode) != "" &&
		strings.TrimSpace(t.PlanCode) != "" &&
		t.PlanVersionID != 0 &&
		strings.TrimSpace(t.PlanVersion) != "" &&
		t.TierRank > 0 &&
		t.SubscriptionTermDays > 0 &&
		t.Price.IsValid()
}

type OrderStatus string

const (
	OrderPendingFunds   OrderStatus = "pending_funds"
	OrderFundsConfirmed OrderStatus = "funds_confirmed"
	OrderRequiresReview OrderStatus = "requires_review"
)

type Order struct {
	ID                 uint64
	Subject            Subject
	Terms              CommercialTerms
	CommandKey         string
	CommandFingerprint string
	Reference          string
	Status             OrderStatus
	ExpiresAt          time.Time
	FundsConfirmedAt   *time.Time
	CreatedAt          time.Time
}

func (o Order) IsValid() bool {
	return o.Subject.IsValid() && o.Terms.IsValid() &&
		strings.TrimSpace(o.CommandKey) != "" &&
		strings.TrimSpace(o.CommandFingerprint) != "" &&
		strings.TrimSpace(o.Reference) != "" &&
		o.Status != "" && !o.CreatedAt.IsZero()
}

type ProviderEventType string

const (
	ProviderAuthorized     ProviderEventType = "authorized"
	ProviderCaptured       ProviderEventType = "captured"
	ProviderFundsConfirmed ProviderEventType = "funds_confirmed"
	ProviderReversed       ProviderEventType = "reversed"
	ProviderRefunded       ProviderEventType = "refunded"
)

type ProviderEvent struct {
	Provider      string
	TransactionID string
	Reference     string
	Amount        Money
	OccurredAt    time.Time
	Type          ProviderEventType
	EvidenceHash  string
}

func (e ProviderEvent) IsValid() bool {
	return strings.TrimSpace(e.Provider) != "" &&
		strings.TrimSpace(e.TransactionID) != "" &&
		strings.TrimSpace(e.Reference) != "" &&
		e.Amount.IsValid() && !e.OccurredAt.IsZero() && e.Type != ""
}

type SettlementStatus string

const (
	SettlementFundsConfirmed SettlementStatus = "funds_confirmed"
	SettlementUnmatched      SettlementStatus = "unmatched"
	SettlementRequiresReview SettlementStatus = "requires_review"
	SettlementLate           SettlementStatus = "late"
)

type Settlement struct {
	ID                    uint64
	OrderID               uint64
	Provider              string
	ProviderTransactionID string
	Reference             string
	Amount                Money
	OccurredAt            time.Time
	EvidenceHash          string
	Status                SettlementStatus
	ReviewReason          string
	CreatedAt             time.Time
}

type FulfillmentStatus string

const (
	FulfillmentPending        FulfillmentStatus = "pending"
	FulfillmentRunning        FulfillmentStatus = "running"
	FulfillmentCompleted      FulfillmentStatus = "completed"
	FulfillmentRetry          FulfillmentStatus = "retry"
	FulfillmentRequiresReview FulfillmentStatus = "requires_review"
)

type Fulfillment struct {
	ID           uint64
	OrderID      uint64
	Status       FulfillmentStatus
	AttemptCount int
	AvailableAt  time.Time
	LockedBy     string
	ClaimVersion uint64
	LeaseUntil   time.Time
	LastError    string
	CompletedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type FulfillmentClaim struct {
	Fulfillment Fulfillment
	Order       Order
}

type SettlementEffect struct {
	EffectKey            string
	OrderID              uint64
	Subject              Subject
	ProductCode          string
	PlanCode             string
	PlanVersionID        uint64
	PlanVersion          string
	TierRank             int32
	SubscriptionTermDays int32
	TermsChecksum        string
	OccurredAt           time.Time
}

type RedriveCommand struct {
	FulfillmentID uint64
	CommandKey    string
	ActorID       string
}

type PaymentMethod string

const (
	PaymentMethodCard         PaymentMethod = "card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodQR           PaymentMethod = "qr"
)

func (m PaymentMethod) IsValid() bool {
	switch m {
	case PaymentMethodCard, PaymentMethodBankTransfer, PaymentMethodQR:
		return true
	default:
		return false
	}
}

type AttemptStatus string

const (
	AttemptCreated       AttemptStatus = "created"
	AttemptPendingAction AttemptStatus = "pending_action"
	AttemptProcessing    AttemptStatus = "processing"
	AttemptAuthorized    AttemptStatus = "authorized"
	AttemptSucceeded     AttemptStatus = "succeeded"
	AttemptDeclined      AttemptStatus = "declined"
	AttemptExpired       AttemptStatus = "expired"
	AttemptCanceled      AttemptStatus = "canceled"
	AttemptFailed        AttemptStatus = "failed"
)

type NextActionKind string

const (
	NextActionNone      NextActionKind = "none"
	NextActionRedirect  NextActionKind = "redirect"
	NextActionDisplayQR NextActionKind = "display_qr"
	NextActionWait      NextActionKind = "wait"
)

type NextAction struct {
	Kind        NextActionKind
	RedirectURL string
	QRPayload   string
	ExpiresAt   *time.Time
}

func (a NextAction) IsValid() bool {
	switch a.Kind {
	case NextActionNone, NextActionWait:
		return true
	case NextActionRedirect:
		return strings.TrimSpace(a.RedirectURL) != ""
	case NextActionDisplayQR:
		return strings.TrimSpace(a.QRPayload) != ""
	default:
		return false
	}
}

type PaymentAttempt struct {
	ID                 uint64
	OrderID            uint64
	CommandKey         string
	CommandFingerprint string
	Method             PaymentMethod
	Provider           string
	ProviderReference  string
	Status             AttemptStatus
	NextAction         NextAction
	ExpiresAt          *time.Time
	FailureCode        string
	FailureDetail      string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (a PaymentAttempt) IsValid() bool {
	return a.OrderID != 0 &&
		strings.TrimSpace(a.CommandKey) != "" &&
		strings.TrimSpace(a.CommandFingerprint) != "" &&
		a.Method.IsValid() && strings.TrimSpace(a.Provider) != "" &&
		a.Status != "" && a.NextAction.IsValid() && !a.CreatedAt.IsZero()
}

type ProviderAttempt struct {
	ProviderReference string
	Status            AttemptStatus
	NextAction        NextAction
	ExpiresAt         *time.Time
}

func (a ProviderAttempt) IsValid() bool {
	return strings.TrimSpace(a.ProviderReference) != "" && a.Status != "" && a.NextAction.IsValid()
}
