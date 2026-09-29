package subscription

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	catalogdomain "user/internal/domain/plan"
	"user/internal/models"
)

/** Aggregate is one durable subscription with its plan identity. */
type Aggregate struct {
	ID                 uint64
	SubscriptionKey    string
	SubjectKind        models.SubscriptionSubjectKind
	SubjectID          string
	ProductID          uint64
	ProductCode        string
	PlanVersionID      uint64
	PlanCode           string
	PlanVersion        string
	PlanTermsChecksum  string
	PlanStatus         catalogdomain.Status
	Status             models.SubscriptionStatus
	StartedAt          time.Time
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	AccessUntil        *time.Time
	AutoRenew          bool
	OrderReference     string
	CanceledAt         *time.Time
	EndedAt            *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

/** Copy returns an aggregate that does not share pointer fields. */
func (a Aggregate) Copy() Aggregate {
	copied := a
	copied.AccessUntil = copyTime(a.AccessUntil)
	copied.CanceledAt = copyTime(a.CanceledAt)
	copied.EndedAt = copyTime(a.EndedAt)
	return copied
}

/** Validate verifies the persisted subscription shape. */
func (a Aggregate) Validate() error {
	if a.ID == 0 || a.ProductID == 0 || a.PlanVersionID == 0 {
		return fmt.Errorf("subscription identity is incomplete")
	}
	if strings.TrimSpace(a.SubscriptionKey) == "" {
		return fmt.Errorf("subscription key is required")
	}
	if a.SubjectKind != models.SubscriptionSubjectProfile &&
		a.SubjectKind != models.SubscriptionSubjectOrganization {
		return fmt.Errorf("unsupported subject kind %q", a.SubjectKind)
	}
	if strings.TrimSpace(a.SubjectID) == "" {
		return fmt.Errorf("subscription subject is required")
	}
	switch a.Status {
	case models.SubscriptionPending,
		models.SubscriptionActive,
		models.SubscriptionPastDue,
		models.SubscriptionCanceled,
		models.SubscriptionExpired:
	default:
		return fmt.Errorf("unsupported subscription status %q", a.Status)
	}
	if strings.TrimSpace(a.ProductCode) == "" ||
		strings.TrimSpace(a.PlanCode) == "" ||
		strings.TrimSpace(a.PlanVersion) == "" {
		return fmt.Errorf("subscription plan identity is incomplete")
	}
	if len(a.PlanTermsChecksum) != 64 {
		return fmt.Errorf("plan terms checksum is invalid")
	}
	if _, err := hex.DecodeString(a.PlanTermsChecksum); err != nil {
		return fmt.Errorf("plan terms checksum is invalid")
	}
	if a.CurrentPeriodEnd.Before(a.CurrentPeriodStart) ||
		a.CurrentPeriodEnd.Equal(a.CurrentPeriodStart) {
		return fmt.Errorf("subscription period is invalid")
	}
	if a.StartedAt.IsZero() || a.CurrentPeriodStart.IsZero() ||
		a.CurrentPeriodEnd.IsZero() {
		return fmt.Errorf("subscription dates are incomplete")
	}
	if a.AccessUntil != nil && a.AccessUntil.Before(a.StartedAt) {
		return fmt.Errorf("subscription access window is invalid")
	}
	if a.CreatedAt.IsZero() || a.UpdatedAt.IsZero() || a.UpdatedAt.Before(a.CreatedAt) {
		return fmt.Errorf("subscription audit timestamps are invalid")
	}
	return nil
}

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}
