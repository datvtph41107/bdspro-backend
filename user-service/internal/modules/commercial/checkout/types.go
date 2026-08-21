package checkout

import (
	"encoding/hex"
	"strings"
	"time"

	"common/productcatalog"
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

type Money struct {
	Currency    string
	AmountMinor int64
}

func (m Money) IsValid() bool {
	return strings.TrimSpace(m.Currency) != "" && m.AmountMinor > 0
}

// PlanTerms is the server-resolved immutable checkout snapshot sent to Payment.
// TierRank belongs to the stable Plan; every other commercial value is resolved
// from the selected immutable PlanVersion/price at checkout time.
type PlanTerms struct {
	ProductID            uint64
	ProductCode          string
	PlanID               uint64
	PlanCode             string
	PlanVersionID        uint64
	PlanVersion          string
	TierRank             int32
	SubscriptionTermDays int32
	TermsChecksum        string
	SubjectScope         productcatalog.SubjectScope
	Price                Money
}

func (t PlanTerms) IsValid() bool {
	checksum := strings.TrimSpace(t.TermsChecksum)
	if len(checksum) != 64 {
		return false
	}
	if _, err := hex.DecodeString(checksum); err != nil {
		return false
	}
	return t.ProductID != 0 && strings.TrimSpace(t.ProductCode) != "" &&
		t.PlanID != 0 && strings.TrimSpace(t.PlanCode) != "" &&
		t.PlanVersionID != 0 && strings.TrimSpace(t.PlanVersion) != "" &&
		t.TierRank > 0 && t.SubscriptionTermDays > 0 &&
		validSubjectScope(t.SubjectScope) && t.Price.IsValid()
}

func (t PlanTerms) Allows(subject Subject) bool {
	switch t.SubjectScope {
	case productcatalog.SubjectScopeAny:
		return subject.IsValid()
	case productcatalog.SubjectScopeProfile:
		return subject.Kind == SubjectProfile && subject.IsValid()
	case productcatalog.SubjectScopeOrganization:
		return subject.Kind == SubjectOrganization && subject.IsValid()
	default:
		return false
	}
}

func validSubjectScope(scope productcatalog.SubjectScope) bool {
	return scope == productcatalog.SubjectScopeProfile ||
		scope == productcatalog.SubjectScopeOrganization ||
		scope == productcatalog.SubjectScopeAny
}

type PlanChange string

const (
	PlanChangeSame      PlanChange = "same"
	PlanChangeUpgrade   PlanChange = "upgrade"
	PlanChangeDowngrade PlanChange = "downgrade"
)

func ClassifyPlanChange(currentTier, targetTier int32) (PlanChange, error) {
	if currentTier <= 0 || targetTier <= 0 {
		return "", ErrCommercialTermsUnavailable
	}
	switch {
	case targetTier == currentTier:
		return PlanChangeSame, nil
	case targetTier > currentTier:
		return PlanChangeUpgrade, nil
	default:
		return PlanChangeDowngrade, nil
	}
}

type CurrentSubscription struct {
	ID                   uint64
	PlanVersionID        uint64
	PlanCode             string
	TierRank             int32
	Status               string
	PendingPlanVersionID uint64
}

type CheckoutCommand struct {
	Subject    Subject
	PlanCode   string
	Terms      PlanTerms
	CommandKey string
	CreatedAt  time.Time
}

func (c CheckoutCommand) IsValid() bool {
	planCode := strings.TrimSpace(c.PlanCode)
	return c.Subject.IsValid() && planCode != "" &&
		c.Terms.IsValid() && c.Terms.Allows(c.Subject) &&
		planCode == strings.TrimSpace(c.Terms.PlanCode) &&
		strings.TrimSpace(c.CommandKey) != "" && !c.CreatedAt.IsZero()
}

type PaymentOrder struct {
	ID        uint64
	Reference string
	Status    string
}

type Result struct {
	Order   PaymentOrder
	Created bool
	Terms   PlanTerms
}
