package backfill

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"common/productcatalog"
	"user/internal/models"
	"user/internal/modules/catalog/migration"
	"user/internal/modules/subscription/domain"
)

const (
	EventTypeLegacyMembershipBackfilled = "migration.legacy_membership_backfilled"
	PeriodPaidAtToFinishedDate          = "paid_at_to_finished_date"
	KeyStrategyLegacyMembershipID       = "legacy_membership_id"
)

var (
	ErrNotFound           = errors.New("subscription backfill record not found")
	ErrConcurrentChange   = errors.New("subscription backfill record changed concurrently")
	ErrConflictingTarget  = errors.New("subscription backfill target conflicts with existing data")
	ErrCandidateBlocked   = errors.New("subscription backfill candidate is blocked")
	ErrExecutionForbidden = errors.New("subscription backfill execution is not enabled")
)

/** ReviewDecision records explicit legacy interpretation. */
type ReviewDecision struct {
	PeriodHandling          string `json:"periodHandling"`
	SubscriptionKeyStrategy string `json:"subscriptionKeyStrategy"`
	AutoRenew               bool   `json:"autoRenew"`
	Reason                  string `json:"reason"`
}

/** Validate rejects implicit or unsupported migration semantics. */
func (d ReviewDecision) Validate() error {
	if d.PeriodHandling != PeriodPaidAtToFinishedDate {
		return fmt.Errorf("unsupported period handling %q", d.PeriodHandling)
	}
	if d.SubscriptionKeyStrategy != KeyStrategyLegacyMembershipID {
		return fmt.Errorf("unsupported subscription key strategy %q", d.SubscriptionKeyStrategy)
	}
	if d.AutoRenew {
		return fmt.Errorf("legacy membership does not prove auto renewal")
	}
	if strings.TrimSpace(d.Reason) == "" {
		return fmt.Errorf("review reason is required")
	}
	return nil
}

/** ReviewRequest approves one current legacy membership snapshot. */
type ReviewRequest struct {
	LegacyMembershipID uint64
	ReviewedBy         uint64
	ReviewedAt         time.Time
	Decision           ReviewDecision
}

/** Mapping is the durable membership migration state. */
type Mapping struct {
	LegacyMembershipID       uint64
	SubscriptionID           *uint64
	SourceChecksum           string
	Status                   models.CatalogMigrationStatus
	Decision                 json.RawMessage
	ReviewedBy               *uint64
	ReviewedAt               *time.Time
	MigratedAt               *time.Time
	BlockedReason            string
	LastRunID                string
	LastError                string
	CutoverEvidenceChecksum  string
	CutoverAt                *time.Time
	CutoverBy                string
	RollbackEvidenceChecksum string
	RolledBackAt             *time.Time
	RolledBackBy             string
}

/** PlanTarget is the approved durable plan-version target. */
type PlanTarget struct {
	LegacyPlanID      uint64
	ProductID         uint64
	ProductCode       string
	PlanVersionID     uint64
	PlanCode          string
	PlanVersion       string
	PlanTermsChecksum string
	PlanStatus        productcatalog.Status
	PlanPublishedAt   *time.Time
	SubjectScope      productcatalog.SubjectScope
	EffectiveFrom     *time.Time
	EffectiveUntil    *time.Time
	MappingStatus     models.CatalogMigrationStatus
}

/** Candidate is one locked migration aggregate. */
type Candidate struct {
	Membership           models.Membership
	Mapping              Mapping
	Target               PlanTarget
	ExistingSubscription *domain.Aggregate
}

/** Run identifies one explicit migration execution. */
type Run struct {
	RunID       string
	ActorID     string
	RequestID   string
	OperationID string
	AsOf        time.Time
	ExecutedAt  time.Time
}

/** Validate rejects anonymous or nondeterministic runs. */
func (r Run) Validate() error {
	for name, value := range map[string]string{
		"run id":   r.RunID,
		"actor id": r.ActorID,
	} {
		if !isSafeIdentifier(value, false) {
			return fmt.Errorf("%s is invalid", name)
		}
	}
	for name, value := range map[string]string{
		"request id":   r.RequestID,
		"operation id": r.OperationID,
	} {
		if !isSafeIdentifier(value, true) {
			return fmt.Errorf("%s is invalid", name)
		}
	}
	if r.AsOf.IsZero() || r.ExecutedAt.IsZero() {
		return fmt.Errorf("run timestamps are required")
	}
	return nil
}

func isSafeIdentifier(value string, optional bool) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return optional
	}
	if len(value) > 128 {
		return false
	}
	for _, character := range value {
		if character >= 'a' && character <= 'z' ||
			character >= 'A' && character <= 'Z' ||
			character >= '0' && character <= '9' ||
			strings.ContainsRune("._:-", character) {
			continue
		}
		return false
	}
	return true
}

type ResultStatus string

const (
	ResultPlanned    ResultStatus = "planned"
	ResultMigrated   ResultStatus = "migrated"
	ResultIdempotent ResultStatus = "idempotent"
	ResultBlocked    ResultStatus = "blocked"
	ResultFailed     ResultStatus = "failed"
)

/** Result describes one membership migration attempt. */
type Result struct {
	LegacyMembershipID uint64       `json:"legacyMembershipId"`
	SubscriptionID     uint64       `json:"subscriptionId,omitempty"`
	Status             ResultStatus `json:"status"`
	Detail             string       `json:"detail,omitempty"`
}

/** Event records the append-only migration proof. */
type Event struct {
	SubscriptionID uint64
	PlanVersionID  uint64
	EffectiveAt    time.Time
	Run            Run
	Metadata       json.RawMessage
}

/** Repository is the transactional migration boundary. */
type Repository interface {
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
	LoadCandidate(ctx context.Context, legacyMembershipID uint64, forUpdate bool) (Candidate, error)
	SaveReview(ctx context.Context, request ReviewRequest, checksum string, decision json.RawMessage) error
	ListReviewedMembershipIDs(ctx context.Context, afterID uint64, limit int) ([]uint64, error)
	FindSubscriptionByKey(ctx context.Context, key string) (domain.Aggregate, error)
	FindLiveSubscription(
		ctx context.Context,
		subjectKind models.SubscriptionSubjectKind,
		subjectID string,
		productID uint64,
	) (domain.Aggregate, error)
	CreateSubscription(
		ctx context.Context,
		draft migration.SubscriptionDraft,
		target PlanTarget,
		decision ReviewDecision,
		run Run,
	) (uint64, error)
	EnsureMigrationEvent(ctx context.Context, event Event) error
	MarkShadowRead(
		ctx context.Context,
		legacyMembershipID uint64,
		subscriptionID uint64,
		run Run,
	) error
	MarkBlocked(ctx context.Context, legacyMembershipID uint64, run Run, reason string) error
}

/** LegacyMembershipSource provides bounded coverage scans. */
type LegacyMembershipSource interface {
	ListLegacyMemberships(ctx context.Context, afterID uint64, limit int) ([]models.Membership, error)
}

/** EventEvidence checks durable rollback proof. */
type EventEvidence interface {
	HasLegacyBackfillEvent(
		ctx context.Context,
		subscriptionID uint64,
		legacyMembershipID uint64,
		sourceChecksum string,
	) (bool, error)
}
