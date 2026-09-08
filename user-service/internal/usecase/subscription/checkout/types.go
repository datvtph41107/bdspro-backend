package checkout

import (
	"crypto/sha256"
	"encoding/hex"
	"net/mail"
	"strconv"
	"strings"
	"time"

	plan "user/internal/domain/plan"
)

type Contact struct {
	FullName string
	Email    string
	Phone    string
}

func (c Contact) IsValidInput() bool {
	name, email := strings.TrimSpace(c.FullName), strings.TrimSpace(c.Email)
	if len(name) < 2 || len(name) > 160 || len(email) > 254 {
		return false
	}
	address, err := mail.ParseAddress(email)
	return err == nil && strings.EqualFold(address.Address, email)
}

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
// TierRank belongs to the stable Plan; every other plan value is resolved
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
	SubjectScope         plan.SubjectScope
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
	case plan.SubjectScopeAny:
		return subject.IsValid()
	case plan.SubjectScopeProfile:
		return subject.Kind == SubjectProfile && subject.IsValid()
	case plan.SubjectScopeOrganization:
		return subject.Kind == SubjectOrganization && subject.IsValid()
	default:
		return false
	}
}

func validSubjectScope(scope plan.SubjectScope) bool {
	return scope == plan.SubjectScopeProfile ||
		scope == plan.SubjectScopeOrganization ||
		scope == plan.SubjectScopeAny
}

type PlanChange string

const (
	PlanChangeSame      PlanChange = "same"
	PlanChangeUpgrade   PlanChange = "upgrade"
	PlanChangeDowngrade PlanChange = "downgrade"
)

func ClassifyPlanChange(currentTier, targetTier int32) (PlanChange, error) {
	if currentTier <= 0 || targetTier <= 0 {
		return "", ErrPlanTermsUnavailable
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
	Subject  Subject
	PlanCode string
	// RequestHash is a server-derived digest of the immutable checkout input.
	// It is deliberately not accepted from the client: the handler derives the
	// subject from the verified actor and the service derives this value from
	// the normalized plan code. The durable command key arbitrates retries;
	// RequestHash detects reuse of that key for a different semantic command.
	RequestHash string
	Terms       PlanTerms
	CommandKey  string
	CreatedAt   time.Time
}

// RequestHashFor returns the canonical SHA-256 identity of a Checkout request.
// Length-delimited fields are used so no pair of valid values can produce the
// same preimage through separator ambiguity. Callers still validate Subject
// and planCode before accepting the resulting command.
func RequestHashFor(subject Subject, planCode string) string {
	kind := strings.TrimSpace(string(subject.Kind))
	id := strings.TrimSpace(subject.ID)
	code := strings.TrimSpace(planCode)
	payload := strings.Join([]string{
		strconv.Itoa(len(kind)), ":", kind,
		strconv.Itoa(len(id)), ":", id,
		strconv.Itoa(len(code)), ":", code,
	}, "")
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func (c CheckoutCommand) IsValid() bool {
	planCode := strings.TrimSpace(c.PlanCode)
	return c.Subject.IsValid() && planCode != "" &&
		strings.TrimSpace(c.RequestHash) == RequestHashFor(c.Subject, planCode) &&
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
