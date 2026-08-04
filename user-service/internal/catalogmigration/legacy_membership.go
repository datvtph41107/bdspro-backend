package catalogmigration

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"common/productcatalog"
)

type LegacyMembershipSnapshot struct {
	ID           uint64
	ProfileID    uint64
	PlanID       uint64
	FinishedDate *time.Time
	OrderID      uint64
	PaidAt       *time.Time
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

type MembershipDecision struct {
	ProductCode       string
	PlanCode          string
	PlanVersion       string
	PlanTermsChecksum string
	PeriodHandling    string
}

/** MembershipSourceChecksum fingerprints one legacy membership row. */
func MembershipSourceChecksum(snapshot LegacyMembershipSnapshot) (string, error) {
	return checksum(snapshot)
}

type SubscriptionDraft struct {
	LegacyMembershipID uint64
	LegacyPlanID       uint64
	SourceChecksum     string
	SubscriptionKey    string
	SubjectKind        string
	SubjectID          string
	ProductCode        string
	PlanCode           string
	PlanVersion        string
	PlanTermsChecksum  string
	Status             string
	StartedAt          time.Time
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	AccessUntil        time.Time
	OrderReference     string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

/** BuildSubscriptionDraft converts one complete legacy membership. */
func BuildSubscriptionDraft(
	snapshot LegacyMembershipSnapshot,
	decision MembershipDecision,
	asOf time.Time,
) (SubscriptionDraft, error) {
	if snapshot.ID == 0 || snapshot.ProfileID == 0 || snapshot.PlanID == 0 {
		return SubscriptionDraft{}, fmt.Errorf("membership identity is incomplete")
	}
	if snapshot.PaidAt == nil || snapshot.FinishedDate == nil {
		return SubscriptionDraft{}, fmt.Errorf(
			"membership requires paid_at and finished_date",
		)
	}
	if snapshot.CreatedAt == nil || snapshot.UpdatedAt == nil ||
		snapshot.CreatedAt.IsZero() || snapshot.UpdatedAt.IsZero() ||
		snapshot.UpdatedAt.Before(*snapshot.CreatedAt) {
		return SubscriptionDraft{}, fmt.Errorf(
			"membership audit timestamps are incomplete",
		)
	}
	if !snapshot.FinishedDate.After(*snapshot.PaidAt) {
		return SubscriptionDraft{}, fmt.Errorf("membership period is invalid")
	}
	if asOf.IsZero() {
		return SubscriptionDraft{}, fmt.Errorf("migration as_of is required")
	}
	if decision.PeriodHandling != "paid_at_to_finished_date" {
		return SubscriptionDraft{}, fmt.Errorf(
			"membership period handling requires explicit approval",
		)
	}
	if strings.TrimSpace(decision.ProductCode) == "" ||
		strings.TrimSpace(decision.PlanCode) == "" ||
		strings.TrimSpace(decision.PlanVersion) == "" {
		return SubscriptionDraft{}, fmt.Errorf("target plan identity is required")
	}
	if len(decision.PlanTermsChecksum) != 64 {
		return SubscriptionDraft{}, fmt.Errorf("target plan checksum is required")
	}
	if _, err := hex.DecodeString(decision.PlanTermsChecksum); err != nil {
		return SubscriptionDraft{}, fmt.Errorf("target plan checksum is invalid")
	}
	if _, err := productcatalog.New(productcatalog.Catalog{
		Version: "1.0.0",
		Products: []productcatalog.Product{
			{
				Code:        decision.ProductCode,
				DisplayName: "Migration target",
				Plans: []productcatalog.PlanVersion{
					{
						ProductCode:  decision.ProductCode,
						PlanCode:     decision.PlanCode,
						Version:      decision.PlanVersion,
						DisplayName:  "Migration target",
						Status:       productcatalog.StatusDraft,
						SubjectScope: productcatalog.SubjectScopeProfile,
					},
				},
			},
		},
	}); err != nil {
		return SubscriptionDraft{}, fmt.Errorf("target plan identity is invalid: %w", err)
	}

	status := "active"
	if asOf.Before(*snapshot.PaidAt) {
		status = "pending"
	} else if !asOf.Before(*snapshot.FinishedDate) {
		status = "expired"
	}

	checksum, err := MembershipSourceChecksum(snapshot)
	if err != nil {
		return SubscriptionDraft{}, err
	}

	orderReference := ""
	if snapshot.OrderID > 0 {
		orderReference = strconv.FormatUint(snapshot.OrderID, 10)
	}

	return SubscriptionDraft{
		LegacyMembershipID: snapshot.ID,
		LegacyPlanID:       snapshot.PlanID,
		SourceChecksum:     checksum,
		SubscriptionKey:    "sub_legacy_membership_" + strconv.FormatUint(snapshot.ID, 10),
		SubjectKind:        "profile",
		SubjectID:          strconv.FormatUint(snapshot.ProfileID, 10),
		ProductCode:        decision.ProductCode,
		PlanCode:           decision.PlanCode,
		PlanVersion:        decision.PlanVersion,
		PlanTermsChecksum:  decision.PlanTermsChecksum,
		Status:             status,
		StartedAt:          *snapshot.PaidAt,
		CurrentPeriodStart: *snapshot.PaidAt,
		CurrentPeriodEnd:   *snapshot.FinishedDate,
		AccessUntil:        *snapshot.FinishedDate,
		OrderReference:     orderReference,
		CreatedAt:          *snapshot.CreatedAt,
		UpdatedAt:          *snapshot.UpdatedAt,
	}, nil
}
