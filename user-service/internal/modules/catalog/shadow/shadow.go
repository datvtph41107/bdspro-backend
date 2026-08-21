package shadow

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"common/productcatalog"
	"user/internal/models"
	"user/internal/modules/catalog/domain"
	"user/internal/modules/catalog/migration"
)

const (
	ModeEnvironment    = "QHPRO_CATALOG_SHADOW_READ_MODE"
	TimeoutEnvironment = "QHPRO_CATALOG_SHADOW_READ_TIMEOUT"
)

var ErrNotFound = errors.New("catalog shadow record not found")

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
	StatusLegacyChanged       Status = "legacy_changed"
	StatusCatalogMissing      Status = "catalog_missing"
	StatusCatalogNotActive    Status = "catalog_not_active"
	StatusCatalogNotEffective Status = "catalog_not_effective"
	StatusCatalogTermsDrift   Status = "catalog_terms_drift"
	StatusReadError           Status = "read_error"
)

/** Mapping links one legacy plan to one catalog version. */
type Mapping struct {
	LegacyPlanID   uint64
	PlanVersionID  *uint64
	SourceChecksum string
	Status         models.CatalogMigrationStatus
}

/** Repository loads shadow-read catalog data. */
type Repository interface {
	GetLegacyPlanMapping(ctx context.Context, legacyPlanID uint64) (Mapping, error)
	GetPlanVersion(ctx context.Context, planVersionID uint64) (
		domain.PlanVersionAggregate,
		error,
	)
}

/** Observer runs shadow reads without changing legacy responses. */
type Observer interface {
	ObserveVisiblePlans(ctx context.Context, plans []models.Plan)
}

/** Entry is one legacy-to-catalog comparison. */
type Entry struct {
	LegacyPlanID  uint64
	PlanVersionID uint64
	Status        Status
	Detail        string
}

/** Report summarizes one shadow-read sample. */
type Report struct {
	ComparedAt time.Time
	Entries    []Entry
	Counts     map[Status]int
}

/** Logger is the minimal structured log dependency. */
type Logger interface {
	Printf(format string, values ...any)
}

/** Config controls bounded shadow reads. */
type Config struct {
	Mode    Mode
	Timeout time.Duration
	Now     func() time.Time
	Logger  Logger
}

/** Reader compares catalog v2 without changing production output. */
// @bind: user/internal/modules/catalog/shadow.Observer
type Reader struct {
	repository Repository
	mode       Mode
	timeout    time.Duration
	now        func() time.Time
	logger     Logger
}

/** NewReader builds an environment-controlled reader. */
func NewReader(repository Repository) *Reader {
	return BuildReaderWithConfig(repository, Config{
		Mode:    parseMode(os.Getenv(ModeEnvironment)),
		Timeout: parseTimeout(os.Getenv(TimeoutEnvironment)),
		Now:     time.Now,
		Logger:  log.Default(),
	})
}

/** BuildReaderWithConfig creates a deterministic reader. */
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

/** ObserveVisiblePlans records comparison results and never changes output. */
func (r *Reader) ObserveVisiblePlans(ctx context.Context, plans []models.Plan) {
	if r == nil || r.mode == ModeOff || r.repository == nil {
		return
	}

	shadowCtx, cancel := context.WithTimeout(
		context.WithoutCancel(ctx),
		r.timeout,
	)
	defer cancel()

	report := r.CompareVisiblePlans(shadowCtx, plans, r.now().UTC())
	r.logger.Printf(
		"catalog-shadow visible_plans compared=%d matched=%d mismatched=%d counts=%v",
		len(report.Entries),
		report.Counts[StatusMatched],
		len(report.Entries)-report.Counts[StatusMatched],
		report.Counts,
	)
}

/** CompareVisiblePlans compares legacy rows with approved catalog mappings. */
func (r *Reader) CompareVisiblePlans(
	ctx context.Context,
	plans []models.Plan,
	at time.Time,
) Report {
	report := Report{
		ComparedAt: at,
		Entries:    make([]Entry, 0, len(plans)),
		Counts:     make(map[Status]int),
	}
	if r == nil || r.repository == nil {
		return report
	}

	ordered := append([]models.Plan(nil), plans...)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].ID < ordered[j].ID
	})

	for _, plan := range ordered {
		entry := r.comparePlan(ctx, plan, at)
		report.Entries = append(report.Entries, entry)
		report.Counts[entry.Status]++
	}
	return report
}

func (r *Reader) comparePlan(ctx context.Context, plan models.Plan, at time.Time) Entry {
	entry := Entry{LegacyPlanID: plan.ID}
	mapping, err := r.repository.GetLegacyPlanMapping(ctx, plan.ID)
	if errors.Is(err, ErrNotFound) {
		entry.Status = StatusMissingMapping
		return entry
	}
	if err != nil {
		entry.Status = StatusReadError
		entry.Detail = err.Error()
		return entry
	}
	if mapping.PlanVersionID == nil || *mapping.PlanVersionID == 0 {
		entry.Status = StatusMappingIncomplete
		return entry
	}
	entry.PlanVersionID = *mapping.PlanVersionID
	if mapping.Status != models.CatalogMigrationShadowRead &&
		mapping.Status != models.CatalogMigrationCutover {
		entry.Status = StatusMappingNotEnabled
		entry.Detail = string(mapping.Status)
		return entry
	}

	checksum, err := migration.SourceChecksum(toLegacySnapshot(plan))
	if err != nil {
		entry.Status = StatusReadError
		entry.Detail = err.Error()
		return entry
	}
	if checksum != mapping.SourceChecksum {
		entry.Status = StatusLegacyChanged
		return entry
	}

	aggregate, err := r.repository.GetPlanVersion(ctx, *mapping.PlanVersionID)
	if errors.Is(err, ErrNotFound) {
		entry.Status = StatusCatalogMissing
		return entry
	}
	if err != nil {
		entry.Status = StatusReadError
		entry.Detail = err.Error()
		return entry
	}
	if aggregate.Status != productcatalog.StatusActive ||
		aggregate.EffectiveFrom == nil || aggregate.PublishedAt == nil {
		entry.Status = StatusCatalogNotActive
		return entry
	}
	if at.Before(*aggregate.EffectiveFrom) ||
		(aggregate.EffectiveUntil != nil && !at.Before(*aggregate.EffectiveUntil)) {
		entry.Status = StatusCatalogNotEffective
		return entry
	}

	checksum, err = productcatalog.TermsChecksum(aggregate.CurrentContract())
	if err != nil {
		entry.Status = StatusReadError
		entry.Detail = err.Error()
		return entry
	}
	if checksum != aggregate.TermsChecksum {
		entry.Status = StatusCatalogTermsDrift
		return entry
	}

	entry.Status = StatusMatched
	return entry
}

func toLegacySnapshot(plan models.Plan) migration.LegacyPlanSnapshot {
	return migration.LegacyPlanSnapshot{
		ID:                 plan.ID,
		PackageName:        plan.PackageName,
		DurationDays:       plan.Duration,
		Price:              plan.Price,
		CreateProduct:      int64(plan.CreateProduct),
		LimitPost:          plan.LimitPost,
		PublicProductLimit: int64(plan.PublicProductLimit),
		CreateProductAI:    int64(plan.CreateProductAI),
		FriendLimit:        int64(plan.FriendLimit),
		CustomerLimit:      int64(plan.CustomerLimit),
		JoinOrg:            int64(plan.JoinOrg),
		CreateOrganization: int64(plan.CreateOrganization),
		OrgMemberLimit:     int64(plan.OrgMemberLimit),
		OrgProductLimit:    int64(plan.OrgProductLimit),
		Color:              plan.Color,
		Features:           append([]string(nil), plan.Features...),
		Visibility:         plan.Visibility,
	}
}

func parseMode(value string) Mode {
	if strings.EqualFold(strings.TrimSpace(value), string(ModeLog)) {
		return ModeLog
	}
	return ModeOff
}

func parseTimeout(value string) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 200 * time.Millisecond
	}
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 || duration > 2*time.Second {
		return 200 * time.Millisecond
	}
	return duration
}

/** ValidateMode checks operator configuration. */
func ValidateMode(value string) error {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || value == string(ModeOff) || value == string(ModeLog) {
		return nil
	}
	return fmt.Errorf("invalid catalog shadow mode %q", value)
}
