package shadow

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"common/productcatalog"
	"user/internal/models"
	"user/internal/modules/catalog/migration"
	"user/internal/modules/subscription/domain"
)

const (
	ModeEnvironment    = "QHPRO_SUBSCRIPTION_SHADOW_READ_MODE"
	TimeoutEnvironment = "QHPRO_SUBSCRIPTION_SHADOW_READ_TIMEOUT"
)

var ErrNotFound = errors.New("subscription shadow record not found")

type Mode string

const (
	ModeOff Mode = "off"
	ModeLog Mode = "log"
)

type Status string

const (
	StatusMatched             Status = "matched"
	StatusMissingMapping      Status = "missing_mapping"
	StatusMappingNotEnabled   Status = "mapping_not_enabled"
	StatusMappingIncomplete   Status = "mapping_incomplete"
	StatusLegacyIncomplete    Status = "legacy_incomplete"
	StatusLegacyChanged       Status = "legacy_changed"
	StatusPlanMappingMissing  Status = "plan_mapping_missing"
	StatusPlanMappingMismatch Status = "plan_mapping_mismatch"
	StatusSubscriptionMissing Status = "subscription_missing"
	StatusSubscriptionInvalid Status = "subscription_invalid"
	StatusSubjectMismatch     Status = "subject_mismatch"
	StatusStatusMismatch      Status = "status_mismatch"
	StatusPeriodMismatch      Status = "period_mismatch"
	StatusOrderMismatch       Status = "order_mismatch"
	StatusAuditMismatch       Status = "audit_mismatch"
	StatusReadError           Status = "read_error"
)

/** Mapping links one legacy membership to one subscription. */
type Mapping struct {
	LegacyMembershipID uint64
	SubscriptionID     *uint64
	SourceChecksum     string
	Status             models.CatalogMigrationStatus
}

/** PlanMapping links the legacy plan to the approved version. */
type PlanMapping struct {
	LegacyPlanID  uint64
	PlanVersionID *uint64
	Status        models.CatalogMigrationStatus
}

/** Repository loads subscription shadow-read data. */
type Repository interface {
	GetLegacyMembershipMapping(ctx context.Context, legacyMembershipID uint64) (Mapping, error)
	GetLegacyPlanMappingForSubscription(ctx context.Context, legacyPlanID uint64) (PlanMapping, error)
	GetSubscription(ctx context.Context, subscriptionID uint64) (domain.Aggregate, error)
}

/** Observer compares one legacy membership without changing output. */
type Observer interface {
	ObserveMembership(ctx context.Context, membership *models.Membership)
}

/** Entry is one membership-to-subscription comparison. */
type Entry struct {
	LegacyMembershipID uint64                        `json:"legacyMembershipId"`
	SubscriptionID     uint64                        `json:"subscriptionId,omitempty"`
	MappingStatus      models.CatalogMigrationStatus `json:"mappingStatus,omitempty"`
	PlanMappingStatus  models.CatalogMigrationStatus `json:"planMappingStatus,omitempty"`
	Status             Status                        `json:"status"`
	Detail             string                        `json:"detail,omitempty"`
}

/** Logger is the minimal log dependency. */
type Logger interface {
	Printf(format string, values ...any)
}

/** Config controls bounded subscription shadow reads. */
type Config struct {
	Mode    Mode
	Timeout time.Duration
	Now     func() time.Time
	Logger  Logger
}

/** Reader compares subscription v2 without changing production output. */
type Reader struct {
	repository Repository
	mode       Mode
	timeout    time.Duration
	now        func() time.Time
	logger     Logger
}

/** NewReader creates an environment-controlled subscription reader. */
func NewReader(repository Repository) *Reader {
	return BuildReaderWithConfig(repository, Config{
		Mode:    parseMode(os.Getenv(ModeEnvironment)),
		Timeout: parseTimeout(os.Getenv(TimeoutEnvironment)),
		Now:     time.Now,
		Logger:  log.Default(),
	})
}

/** BuildReaderWithConfig creates a deterministic subscription reader. */
func BuildReaderWithConfig(repository Repository, config Config) *Reader {
	if config.Mode != ModeLog {
		config.Mode = ModeOff
	}
	if config.Timeout <= 0 {
		config.Timeout = 200 * time.Millisecond
	}
	if config.Now == nil {
		config.Now = time.Now
	}
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	return &Reader{
		repository: repository,
		mode:       config.Mode,
		timeout:    config.Timeout,
		now:        config.Now,
		logger:     config.Logger,
	}
}

/** ObserveMembership records one comparison and never changes output. */
func (r *Reader) ObserveMembership(ctx context.Context, membership *models.Membership) {
	if r == nil || r.mode == ModeOff || r.repository == nil || membership == nil {
		return
	}

	shadowCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	entry := r.CompareMembership(shadowCtx, *membership, r.now().UTC())
	r.logger.Printf(
		"subscription-shadow membership_id=%d subscription_id=%d status=%s detail=%q",
		entry.LegacyMembershipID,
		entry.SubscriptionID,
		entry.Status,
		entry.Detail,
	)
}

/** CompareMembership compares one legacy row with an approved subscription. */
func (r *Reader) CompareMembership(
	ctx context.Context,
	membership models.Membership,
	at time.Time,
) Entry {
	entry := Entry{LegacyMembershipID: membership.ID}
	if r == nil || r.repository == nil {
		entry.Status = StatusReadError
		entry.Detail = "subscription shadow repository is not configured"
		return entry
	}
	if membership.ID == 0 || membership.ProfileID == 0 || membership.PlanID == 0 ||
		membership.PaidAt == nil || membership.FinishedDate == nil {
		entry.Status = StatusLegacyIncomplete
		return entry
	}

	mapping, err := r.repository.GetLegacyMembershipMapping(ctx, membership.ID)
	if errors.Is(err, ErrNotFound) {
		entry.Status = StatusMissingMapping
		return entry
	}
	if err != nil {
		entry.Status = StatusReadError
		entry.Detail = err.Error()
		return entry
	}
	if mapping.SubscriptionID == nil || *mapping.SubscriptionID == 0 {
		entry.Status = StatusMappingIncomplete
		return entry
	}
	entry.SubscriptionID = *mapping.SubscriptionID
	entry.MappingStatus = mapping.Status
	if !shadowEnabled(mapping.Status) {
		entry.Status = StatusMappingNotEnabled
		entry.Detail = string(mapping.Status)
		return entry
	}

	snapshot := toLegacySnapshot(membership)
	checksum, err := migration.MembershipSourceChecksum(snapshot)
	if err != nil {
		entry.Status = StatusReadError
		entry.Detail = err.Error()
		return entry
	}
	if checksum != mapping.SourceChecksum {
		entry.Status = StatusLegacyChanged
		return entry
	}

	planMapping, err := r.repository.GetLegacyPlanMappingForSubscription(ctx, membership.PlanID)
	if errors.Is(err, ErrNotFound) {
		entry.Status = StatusPlanMappingMissing
		return entry
	}
	if err != nil {
		entry.Status = StatusReadError
		entry.Detail = err.Error()
		return entry
	}
	entry.PlanMappingStatus = planMapping.Status
	if planMapping.PlanVersionID == nil || *planMapping.PlanVersionID == 0 ||
		!shadowEnabled(planMapping.Status) {
		entry.Status = StatusPlanMappingMissing
		entry.Detail = string(planMapping.Status)
		return entry
	}

	aggregate, err := r.repository.GetSubscription(ctx, *mapping.SubscriptionID)
	if errors.Is(err, ErrNotFound) {
		entry.Status = StatusSubscriptionMissing
		return entry
	}
	if err != nil {
		entry.Status = StatusReadError
		entry.Detail = err.Error()
		return entry
	}
	if err := aggregate.Validate(); err != nil {
		entry.Status = StatusSubscriptionInvalid
		entry.Detail = err.Error()
		return entry
	}
	if aggregate.PlanVersionID != *planMapping.PlanVersionID ||
		aggregate.PlanStatus != productcatalog.StatusActive {
		entry.Status = StatusPlanMappingMismatch
		return entry
	}
	if aggregate.SubjectKind != models.SubscriptionSubjectProfile ||
		aggregate.SubjectID != strconv.FormatUint(membership.ProfileID, 10) {
		entry.Status = StatusSubjectMismatch
		return entry
	}

	expectedStatus := legacyStatus(membership, at)
	if aggregate.Status != expectedStatus {
		entry.Status = StatusStatusMismatch
		entry.Detail = fmt.Sprintf("legacy=%s catalog=%s", expectedStatus, aggregate.Status)
		return entry
	}
	if !sameInstant(aggregate.StartedAt, *membership.PaidAt) ||
		!sameInstant(aggregate.CurrentPeriodStart, *membership.PaidAt) ||
		!sameInstant(aggregate.CurrentPeriodEnd, *membership.FinishedDate) ||
		aggregate.AccessUntil == nil ||
		!sameInstant(*aggregate.AccessUntil, *membership.FinishedDate) {
		entry.Status = StatusPeriodMismatch
		return entry
	}

	expectedOrder := ""
	if membership.OrderID > 0 {
		expectedOrder = strconv.FormatUint(membership.OrderID, 10)
	}
	if aggregate.OrderReference != expectedOrder {
		entry.Status = StatusOrderMismatch
		return entry
	}
	if membership.CreatedAt == nil || membership.UpdatedAt == nil ||
		!sameInstant(aggregate.CreatedAt, *membership.CreatedAt) ||
		!sameInstant(aggregate.UpdatedAt, *membership.UpdatedAt) {
		entry.Status = StatusAuditMismatch
		return entry
	}
	if aggregate.SubscriptionKey != "sub_legacy_membership_"+strconv.FormatUint(membership.ID, 10) ||
		aggregate.AutoRenew {
		entry.Status = StatusSubscriptionInvalid
		entry.Detail = "legacy subscription identity or renewal policy differs"
		return entry
	}

	entry.Status = StatusMatched
	return entry
}

/** ValidateMode checks the supported rollout modes. */
func ValidateMode(value string) error {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" || normalized == string(ModeOff) || normalized == string(ModeLog) {
		return nil
	}
	return fmt.Errorf("unsupported subscription shadow mode %q", value)
}

func parseMode(value string) Mode {
	if strings.EqualFold(strings.TrimSpace(value), string(ModeLog)) {
		return ModeLog
	}
	return ModeOff
}

func parseTimeout(value string) time.Duration {
	if strings.TrimSpace(value) == "" {
		return 200 * time.Millisecond
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return 200 * time.Millisecond
	}
	return parsed
}

func shadowEnabled(status models.CatalogMigrationStatus) bool {
	return status == models.CatalogMigrationShadowRead ||
		status == models.CatalogMigrationCutover
}

func legacyStatus(membership models.Membership, at time.Time) models.SubscriptionStatus {
	if at.Before(*membership.PaidAt) {
		return models.SubscriptionPending
	}
	if !at.Before(*membership.FinishedDate) {
		return models.SubscriptionExpired
	}
	return models.SubscriptionActive
}

func sameInstant(left time.Time, right time.Time) bool {
	return left.UTC().Equal(right.UTC())
}

func toLegacySnapshot(membership models.Membership) migration.LegacyMembershipSnapshot {
	return migration.LegacyMembershipSnapshot{
		ID:           membership.ID,
		ProfileID:    membership.ProfileID,
		PlanID:       membership.PlanID,
		FinishedDate: copyTime(membership.FinishedDate),
		OrderID:      membership.OrderID,
		PaidAt:       copyTime(membership.PaidAt),
		CreatedAt:    copyTime(membership.CreatedAt),
		UpdatedAt:    copyTime(membership.UpdatedAt),
	}
}

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}
