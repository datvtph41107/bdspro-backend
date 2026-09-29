package settlement

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	catalogdomain "user/internal/domain/plan"
	subscriptiondomain "user/internal/domain/subscription"
	"user/internal/models"
)

type Effect struct {
	EffectKey            string
	OrderID              uint64
	SubjectKind          models.SubscriptionSubjectKind
	SubjectID            string
	ProductCode          string
	PlanCode             string
	PlanVersionID        uint64
	PlanVersion          string
	TierRank             int32
	SubscriptionTermDays int32
	TermsChecksum        string
	OccurredAt           time.Time
}

// SubscriptionKey returns the durable aggregate identity derived from the
// cross-owner settlement identity. Payment's numeric order ID is only unique
// inside one Payment database lifetime, while EffectKey also contains the
// immutable order reference and survives independent database restore/reset.
func (e Effect) SubscriptionKey() string {
	if strings.TrimSpace(e.EffectKey) == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(e.EffectKey))
	return "sub_payment_" + hex.EncodeToString(sum[:])
}

func (e Effect) Validate() error {
	if strings.TrimSpace(e.EffectKey) == "" || e.OrderID == 0 ||
		(e.SubjectKind != models.SubscriptionSubjectProfile && e.SubjectKind != models.SubscriptionSubjectOrganization) ||
		strings.TrimSpace(e.SubjectID) == "" || strings.TrimSpace(e.ProductCode) == "" ||
		strings.TrimSpace(e.PlanCode) == "" || e.PlanVersionID == 0 || strings.TrimSpace(e.PlanVersion) == "" ||
		e.TierRank <= 0 || e.SubscriptionTermDays <= 0 || e.OccurredAt.IsZero() {
		return ErrInvalidEffect
	}
	checksum := strings.TrimSpace(e.TermsChecksum)
	if len(checksum) != 64 {
		return ErrInvalidEffect
	}
	if _, err := hex.DecodeString(checksum); err != nil {
		return ErrInvalidEffect
	}
	return nil
}

func (e Effect) Fingerprint() (string, error) {
	if err := e.Validate(); err != nil {
		return "", err
	}
	payload := struct {
		OrderID              uint64                         `json:"orderId"`
		SubjectKind          models.SubscriptionSubjectKind `json:"subjectKind"`
		SubjectID            string                         `json:"subjectId"`
		ProductCode          string                         `json:"productCode"`
		PlanCode             string                         `json:"planCode"`
		PlanVersionID        uint64                         `json:"planVersionId"`
		PlanVersion          string                         `json:"planVersion"`
		TierRank             int32                          `json:"tierRank"`
		SubscriptionTermDays int32                          `json:"subscriptionTermDays"`
		TermsChecksum        string                         `json:"termsChecksum"`
		OccurredAt           time.Time                      `json:"occurredAt"`
	}{
		OrderID: e.OrderID, SubjectKind: e.SubjectKind, SubjectID: e.SubjectID,
		ProductCode: e.ProductCode, PlanCode: e.PlanCode, PlanVersionID: e.PlanVersionID,
		PlanVersion: e.PlanVersion, TierRank: e.TierRank, SubscriptionTermDays: e.SubscriptionTermDays,
		TermsChecksum: e.TermsChecksum, OccurredAt: e.OccurredAt.UTC(),
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("encode subscription settlement fingerprint: %w", err)
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}

type Action string

const (
	ActionActivated         Action = "activated"
	ActionUpgraded          Action = "upgraded"
	ActionRejectedSameTier  Action = "rejected_same_tier"
	ActionRejectedDowngrade Action = "rejected_downgrade"
)

type Target struct {
	ProductID            uint64
	PlanID               uint64
	PlanVersionID        uint64
	ProductCode          string
	PlanCode             string
	PlanVersion          string
	TierRank             int32
	SubscriptionTermDays int32
	TermsChecksum        string
	Published            bool
	PlanStatus           catalogdomain.Status
}

func (t Target) Matches(effect Effect) bool {
	return t.ProductID != 0 && t.PlanID != 0 && t.Published &&
		t.PlanVersionID == effect.PlanVersionID && t.ProductCode == effect.ProductCode &&
		t.PlanCode == effect.PlanCode && t.PlanVersion == effect.PlanVersion &&
		t.TierRank == effect.TierRank && t.SubscriptionTermDays == effect.SubscriptionTermDays &&
		t.TermsChecksum == effect.TermsChecksum
}

type Current struct {
	Aggregate        subscriptiondomain.Aggregate
	TierRank         int32
	Found            bool
	HasPendingChange bool
}

type Decision struct {
	Action            Action
	Subscription      subscriptiondomain.Aggregate
	FromPlanVersionID *uint64
	ToPlanVersionID   uint64
}

type Result struct {
	Action         Action
	SubscriptionID uint64
	Replayed       bool
}

type Acceptance struct {
	Effect      Effect
	Fingerprint string
}
