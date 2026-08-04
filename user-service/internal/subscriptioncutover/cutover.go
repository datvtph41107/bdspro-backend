package subscriptioncutover

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"user/internal/models"
	"user/internal/subscriptionbackfill"
	"user/internal/subscriptionshadow"
)

const (
	EventTypeCutoverActivated  = "migration.subscription_read_cutover_activated"
	EventTypeCutoverRolledBack = "migration.subscription_read_cutover_rolled_back"
	DefaultBatchEntries        = 500
	MaxBatchEntries            = 1000
)

var (
	ErrEvidenceInvalid  = errors.New("subscription cutover evidence is invalid")
	ErrConcurrentChange = errors.New("subscription cutover state changed concurrently")
	ErrActivationDenied = errors.New("subscription cutover activation is not enabled")
	ErrRollbackDenied   = errors.New("subscription cutover rollback is not enabled")
)

/** Entry freezes one concrete membership target. */
type Entry struct {
	LegacyMembershipID uint64                        `json:"legacyMembershipId"`
	SubscriptionID     uint64                        `json:"subscriptionId"`
	LegacyPlanID       uint64                        `json:"legacyPlanId"`
	PlanVersionID      uint64                        `json:"planVersionId"`
	SubscriptionKey    string                        `json:"subscriptionKey"`
	SourceChecksum     string                        `json:"sourceChecksum"`
	PlanTermsChecksum  string                        `json:"planTermsChecksum"`
	MappingStatus      models.CatalogMigrationStatus `json:"mappingStatus"`
	PlanMappingStatus  models.CatalogMigrationStatus `json:"planMappingStatus"`
}

/** Manifest binds one bounded batch to full-population evidence. */
type Manifest struct {
	SchemaVersion           string    `json:"schemaVersion"`
	GeneratedAt             time.Time `json:"generatedAt"`
	ReadinessReportChecksum string    `json:"readinessReportChecksum"`
	PopulationChecksum      string    `json:"populationChecksum"`
	TotalEntries            int       `json:"totalEntries"`
	BatchIndex              int       `json:"batchIndex"`
	BatchCount              int       `json:"batchCount"`
	Entries                 []Entry   `json:"entries"`
	EvidenceChecksum        string    `json:"evidenceChecksum"`
}

/** Validate verifies canonical batch evidence and every entry. */
func (m Manifest) Validate() error {
	if m.SchemaVersion != "1.1.0" || m.GeneratedAt.IsZero() {
		return fmt.Errorf("%w: manifest header is incomplete", ErrEvidenceInvalid)
	}
	if !validChecksum(m.ReadinessReportChecksum) ||
		!validChecksum(m.PopulationChecksum) || !validChecksum(m.EvidenceChecksum) {
		return fmt.Errorf("%w: manifest checksum is invalid", ErrEvidenceInvalid)
	}
	if m.TotalEntries <= 0 || m.BatchIndex <= 0 || m.BatchCount <= 0 ||
		m.BatchIndex > m.BatchCount || len(m.Entries) == 0 ||
		len(m.Entries) > MaxBatchEntries || len(m.Entries) > m.TotalEntries {
		return fmt.Errorf("%w: manifest batch shape is invalid", ErrEvidenceInvalid)
	}
	seen := make(map[uint64]struct{}, len(m.Entries))
	previous := uint64(0)
	for _, entry := range m.Entries {
		if err := entry.Validate(); err != nil {
			return err
		}
		if _, exists := seen[entry.LegacyMembershipID]; exists {
			return fmt.Errorf("%w: duplicate membership %d", ErrEvidenceInvalid, entry.LegacyMembershipID)
		}
		if previous > entry.LegacyMembershipID {
			return fmt.Errorf("%w: entries are not sorted", ErrEvidenceInvalid)
		}
		seen[entry.LegacyMembershipID] = struct{}{}
		previous = entry.LegacyMembershipID
	}
	expected, err := manifestChecksum(m)
	if err != nil {
		return err
	}
	if expected != m.EvidenceChecksum {
		return fmt.Errorf("%w: evidence checksum differs", ErrEvidenceInvalid)
	}
	return nil
}

/** Plan freezes the whole population before any batch is activated. */
type Plan struct {
	SchemaVersion           string     `json:"schemaVersion"`
	GeneratedAt             time.Time  `json:"generatedAt"`
	ReadinessReportChecksum string     `json:"readinessReportChecksum"`
	PopulationChecksum      string     `json:"populationChecksum"`
	TotalEntries            int        `json:"totalEntries"`
	BatchSize               int        `json:"batchSize"`
	Batches                 []Manifest `json:"batches"`
	EvidenceChecksum        string     `json:"evidenceChecksum"`
}

/** Validate proves complete, ordered, non-overlapping batches. */
func (p Plan) Validate() error {
	if p.SchemaVersion != "1.0.0" || p.GeneratedAt.IsZero() ||
		!validChecksum(p.ReadinessReportChecksum) ||
		!validChecksum(p.PopulationChecksum) || !validChecksum(p.EvidenceChecksum) ||
		p.TotalEntries <= 0 || p.BatchSize <= 0 || p.BatchSize > MaxBatchEntries ||
		len(p.Batches) == 0 {
		return fmt.Errorf("%w: cutover plan header is invalid", ErrEvidenceInvalid)
	}
	expectedBatchCount := (p.TotalEntries + p.BatchSize - 1) / p.BatchSize
	if len(p.Batches) != expectedBatchCount {
		return fmt.Errorf("%w: cutover plan batch count differs", ErrEvidenceInvalid)
	}
	entries := make([]Entry, 0, p.TotalEntries)
	previous := uint64(0)
	seen := make(map[uint64]struct{}, p.TotalEntries)
	for index, batch := range p.Batches {
		if err := batch.Validate(); err != nil {
			return err
		}
		if batch.GeneratedAt != p.GeneratedAt ||
			batch.ReadinessReportChecksum != p.ReadinessReportChecksum ||
			batch.PopulationChecksum != p.PopulationChecksum ||
			batch.TotalEntries != p.TotalEntries || batch.BatchIndex != index+1 ||
			batch.BatchCount != len(p.Batches) {
			return fmt.Errorf("%w: batch %d is not bound to its plan", ErrEvidenceInvalid, index+1)
		}
		for _, entry := range batch.Entries {
			if _, exists := seen[entry.LegacyMembershipID]; exists || previous > entry.LegacyMembershipID {
				return fmt.Errorf("%w: plan entries overlap or are unsorted", ErrEvidenceInvalid)
			}
			seen[entry.LegacyMembershipID] = struct{}{}
			previous = entry.LegacyMembershipID
			entries = append(entries, entry)
		}
	}
	if len(entries) != p.TotalEntries {
		return fmt.Errorf("%w: plan does not cover the full population", ErrEvidenceInvalid)
	}
	population, err := populationChecksum(entries)
	if err != nil {
		return err
	}
	if population != p.PopulationChecksum {
		return fmt.Errorf("%w: population checksum differs", ErrEvidenceInvalid)
	}
	expected, err := planChecksum(p)
	if err != nil {
		return err
	}
	if expected != p.EvidenceChecksum {
		return fmt.Errorf("%w: plan checksum differs", ErrEvidenceInvalid)
	}
	return nil
}

/** Batch returns one one-based, validated activation batch. */
func (p Plan) Batch(index int) (Manifest, error) {
	if err := p.Validate(); err != nil {
		return Manifest{}, err
	}
	if index <= 0 || index > len(p.Batches) {
		return Manifest{}, fmt.Errorf("%w: batch index %d is invalid", ErrEvidenceInvalid, index)
	}
	return p.Batches[index-1], nil
}

/** Validate checks one immutable target row. */
func (e Entry) Validate() error {
	if e.LegacyMembershipID == 0 || e.SubscriptionID == 0 ||
		e.LegacyPlanID == 0 || e.PlanVersionID == 0 {
		return fmt.Errorf("%w: entry identity is incomplete", ErrEvidenceInvalid)
	}
	if strings.TrimSpace(e.SubscriptionKey) == "" ||
		!validChecksum(e.SourceChecksum) || !validChecksum(e.PlanTermsChecksum) {
		return fmt.Errorf("%w: entry evidence is incomplete", ErrEvidenceInvalid)
	}
	if e.MappingStatus != models.CatalogMigrationShadowRead {
		return fmt.Errorf("%w: membership %d is not shadow-ready", ErrEvidenceInvalid, e.LegacyMembershipID)
	}
	if e.PlanMappingStatus != models.CatalogMigrationShadowRead &&
		e.PlanMappingStatus != models.CatalogMigrationCutover {
		return fmt.Errorf("%w: plan mapping is not ready", ErrEvidenceInvalid)
	}
	return nil
}

/** Run identifies one explicit operational change. */
type Run struct {
	RunID       string
	ActorID     string
	RequestID   string
	OperationID string
	ExecutedAt  time.Time
	Reason      string
}

/** Validate rejects anonymous or untraceable changes. */
func (r Run) Validate() error {
	for name, value := range map[string]string{
		"run id": r.RunID, "actor id": r.ActorID,
		"request id": r.RequestID, "operation id": r.OperationID,
	} {
		if !safeIdentifier(value) {
			return fmt.Errorf("%s is invalid", name)
		}
	}
	if r.ExecutedAt.IsZero() || strings.TrimSpace(r.Reason) == "" {
		return fmt.Errorf("run time and reason are required")
	}
	return nil
}

/** RollbackManifest proves exactly which cutover batch is reversed. */
type RollbackManifest struct {
	SchemaVersion           string    `json:"schemaVersion"`
	GeneratedAt             time.Time `json:"generatedAt"`
	CutoverEvidenceChecksum string    `json:"cutoverEvidenceChecksum"`
	PopulationChecksum      string    `json:"populationChecksum"`
	TotalEntries            int       `json:"totalEntries"`
	BatchIndex              int       `json:"batchIndex"`
	BatchCount              int       `json:"batchCount"`
	Reason                  string    `json:"reason"`
	Entries                 []Entry   `json:"entries"`
	EvidenceChecksum        string    `json:"evidenceChecksum"`
}

/** Validate verifies deterministic rollback evidence. */
func (m RollbackManifest) Validate() error {
	if m.SchemaVersion != "1.1.0" || m.GeneratedAt.IsZero() ||
		strings.TrimSpace(m.Reason) == "" ||
		!validChecksum(m.CutoverEvidenceChecksum) ||
		!validChecksum(m.PopulationChecksum) || !validChecksum(m.EvidenceChecksum) ||
		m.TotalEntries <= 0 || m.BatchIndex <= 0 || m.BatchCount <= 0 ||
		m.BatchIndex > m.BatchCount || len(m.Entries) == 0 ||
		len(m.Entries) > MaxBatchEntries || len(m.Entries) > m.TotalEntries {
		return fmt.Errorf("%w: rollback manifest is incomplete", ErrEvidenceInvalid)
	}
	previous := uint64(0)
	seen := make(map[uint64]struct{}, len(m.Entries))
	for _, entry := range m.Entries {
		copyEntry := entry
		copyEntry.MappingStatus = models.CatalogMigrationShadowRead
		if err := copyEntry.Validate(); err != nil {
			return err
		}
		if _, exists := seen[entry.LegacyMembershipID]; exists || previous > entry.LegacyMembershipID {
			return fmt.Errorf("%w: rollback entries overlap or are not sorted", ErrEvidenceInvalid)
		}
		seen[entry.LegacyMembershipID] = struct{}{}
		previous = entry.LegacyMembershipID
	}
	expected, err := rollbackChecksum(m)
	if err != nil {
		return err
	}
	if expected != m.EvidenceChecksum {
		return fmt.Errorf("%w: rollback checksum differs", ErrEvidenceInvalid)
	}
	return nil
}

/** Event is one append-only cutover or rollback proof. */
type Event struct {
	EventType      string
	SubscriptionID uint64
	PlanVersionID  uint64
	EffectiveAt    time.Time
	Run            Run
	Metadata       json.RawMessage
}

/** Repository is the operational cutover transaction boundary. */
type Repository interface {
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
	LoadCandidate(ctx context.Context, legacyMembershipID uint64, forUpdate bool) (subscriptionbackfill.Candidate, error)
	EnsureCutoverEvent(ctx context.Context, event Event) error
	MarkCutover(ctx context.Context, entry Entry, evidenceChecksum string, run Run) error
	EnsureRollbackEvent(ctx context.Context, event Event) error
	MarkRolledBack(ctx context.Context, entry Entry, evidenceChecksum string, run Run) error
}

/** Service prepares, activates, and rolls back explicit manifests. */
type Service struct {
	repository Repository
}

/** NewService creates the operator-only cutover service. */
func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

/** PreparePlan freezes all eligible rows before bounded activation. */
func (s *Service) PreparePlan(
	ctx context.Context,
	readiness subscriptionbackfill.ReadinessResult,
	generatedAt time.Time,
	batchSize int,
) (Plan, error) {
	if s == nil || s.repository == nil {
		return Plan{}, fmt.Errorf("subscription cutover repository is not configured")
	}
	if err := validateStrictReadiness(readiness); err != nil {
		return Plan{}, err
	}
	if generatedAt.IsZero() {
		return Plan{}, fmt.Errorf("%w: generated_at is required", ErrEvidenceInvalid)
	}
	if batchSize <= 0 {
		batchSize = DefaultBatchEntries
	}
	if batchSize > MaxBatchEntries {
		return Plan{}, fmt.Errorf("%w: batch size exceeds %d", ErrEvidenceInvalid, MaxBatchEntries)
	}
	if len(readiness.Report.Entries) == 0 ||
		len(readiness.Report.Entries) != readiness.Report.EligibleLegacy ||
		readiness.Report.Matched != readiness.Report.EligibleLegacy {
		return Plan{}, fmt.Errorf("%w: readiness report must include every matched entry", ErrEvidenceInvalid)
	}

	ids := make([]uint64, 0, len(readiness.Report.Entries))
	for _, observed := range readiness.Report.Entries {
		if observed.Status != subscriptionshadow.StatusMatched || observed.LegacyMembershipID == 0 {
			return Plan{}, fmt.Errorf("%w: readiness contains non-matched rows", ErrEvidenceInvalid)
		}
		ids = append(ids, observed.LegacyMembershipID)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	entries := make([]Entry, 0, len(ids))
	for _, id := range ids {
		candidate, err := s.repository.LoadCandidate(ctx, id, false)
		if err != nil {
			return Plan{}, err
		}
		entry, err := entryFromCandidate(candidate)
		if err != nil {
			return Plan{}, err
		}
		entries = append(entries, entry)
	}
	population, err := populationChecksum(entries)
	if err != nil {
		return Plan{}, err
	}
	batchCount := (len(entries) + batchSize - 1) / batchSize
	plan := Plan{
		SchemaVersion:           "1.0.0",
		GeneratedAt:             generatedAt.UTC(),
		ReadinessReportChecksum: readiness.Report.Checksum,
		PopulationChecksum:      population,
		TotalEntries:            len(entries),
		BatchSize:               batchSize,
		Batches:                 make([]Manifest, 0, batchCount),
	}
	for offset := 0; offset < len(entries); offset += batchSize {
		limit := offset + batchSize
		if limit > len(entries) {
			limit = len(entries)
		}
		batch := Manifest{
			SchemaVersion:           "1.1.0",
			GeneratedAt:             plan.GeneratedAt,
			ReadinessReportChecksum: plan.ReadinessReportChecksum,
			PopulationChecksum:      plan.PopulationChecksum,
			TotalEntries:            plan.TotalEntries,
			BatchIndex:              len(plan.Batches) + 1,
			BatchCount:              batchCount,
			Entries:                 append([]Entry(nil), entries[offset:limit]...),
		}
		checksum, err := manifestChecksum(batch)
		if err != nil {
			return Plan{}, err
		}
		batch.EvidenceChecksum = checksum
		plan.Batches = append(plan.Batches, batch)
	}
	checksum, err := planChecksum(plan)
	if err != nil {
		return Plan{}, err
	}
	plan.EvidenceChecksum = checksum
	return plan, plan.Validate()
}

/** Prepare keeps the original single-batch API for small migrations. */
func (s *Service) Prepare(
	ctx context.Context,
	readiness subscriptionbackfill.ReadinessResult,
	generatedAt time.Time,
) (Manifest, error) {
	plan, err := s.PreparePlan(ctx, readiness, generatedAt, MaxBatchEntries)
	if err != nil {
		return Manifest{}, err
	}
	if len(plan.Batches) != 1 {
		return Manifest{}, fmt.Errorf("%w: use PreparePlan for multi-batch cutover", ErrEvidenceInvalid)
	}
	return plan.Batches[0], nil
}

/** Activate atomically marks every manifest row for catalog reads. */
func (s *Service) Activate(ctx context.Context, manifest Manifest, run Run) error {
	if s == nil || s.repository == nil {
		return fmt.Errorf("subscription cutover repository is not configured")
	}
	if err := manifest.Validate(); err != nil {
		return err
	}
	if err := run.Validate(); err != nil {
		return err
	}
	return s.repository.WithTransaction(ctx, func(txCtx context.Context) error {
		for _, entry := range manifest.Entries {
			candidate, err := s.repository.LoadCandidate(txCtx, entry.LegacyMembershipID, true)
			if err != nil {
				return err
			}
			if candidate.Mapping.Status == models.CatalogMigrationCutover {
				current, currentErr := entryFromCandidateEvidence(candidate)
				if currentErr != nil || current != entry ||
					candidate.Mapping.CutoverEvidenceChecksum != manifest.EvidenceChecksum {
					return fmt.Errorf("%w: membership %d has different cutover evidence", ErrConcurrentChange, entry.LegacyMembershipID)
				}
				continue
			}
			current, err := entryFromCandidate(candidate)
			if err != nil {
				return err
			}
			if current != entry {
				return fmt.Errorf("%w: membership %d differs from manifest", ErrConcurrentChange, entry.LegacyMembershipID)
			}
			metadata, err := eventMetadata(entry, manifest.EvidenceChecksum, manifest.PopulationChecksum, manifest.BatchIndex, manifest.BatchCount)
			if err != nil {
				return err
			}
			if err := s.repository.EnsureCutoverEvent(txCtx, Event{
				EventType:      EventTypeCutoverActivated,
				SubscriptionID: entry.SubscriptionID,
				PlanVersionID:  entry.PlanVersionID,
				EffectiveAt:    run.ExecutedAt,
				Run:            run,
				Metadata:       metadata,
			}); err != nil {
				return err
			}
			if err := s.repository.MarkCutover(txCtx, entry, manifest.EvidenceChecksum, run); err != nil {
				return err
			}
		}
		return nil
	})
}

/** BuildRollbackManifest creates explicit rollback proof. */
func BuildRollbackManifest(
	cutover Manifest,
	reason string,
	generatedAt time.Time,
) (RollbackManifest, error) {
	if err := cutover.Validate(); err != nil {
		return RollbackManifest{}, err
	}
	manifest := RollbackManifest{
		SchemaVersion:           "1.1.0",
		GeneratedAt:             generatedAt.UTC(),
		CutoverEvidenceChecksum: cutover.EvidenceChecksum,
		PopulationChecksum:      cutover.PopulationChecksum,
		TotalEntries:            cutover.TotalEntries,
		BatchIndex:              cutover.BatchIndex,
		BatchCount:              cutover.BatchCount,
		Reason:                  strings.TrimSpace(reason),
		Entries:                 append([]Entry(nil), cutover.Entries...),
	}
	checksum, err := rollbackChecksum(manifest)
	if err != nil {
		return RollbackManifest{}, err
	}
	manifest.EvidenceChecksum = checksum
	return manifest, manifest.Validate()
}

/** Rollback marks cutover rows terminal after feature-flag rollback. */
func (s *Service) Rollback(ctx context.Context, manifest RollbackManifest, run Run) error {
	if s == nil || s.repository == nil {
		return fmt.Errorf("subscription cutover repository is not configured")
	}
	if err := manifest.Validate(); err != nil {
		return err
	}
	if err := run.Validate(); err != nil {
		return err
	}
	return s.repository.WithTransaction(ctx, func(txCtx context.Context) error {
		for _, entry := range manifest.Entries {
			candidate, err := s.repository.LoadCandidate(txCtx, entry.LegacyMembershipID, true)
			if err != nil {
				return err
			}
			if candidate.Mapping.Status != models.CatalogMigrationCutover ||
				candidate.Mapping.CutoverEvidenceChecksum != manifest.CutoverEvidenceChecksum ||
				candidate.Mapping.SubscriptionID == nil ||
				*candidate.Mapping.SubscriptionID != entry.SubscriptionID ||
				candidate.Mapping.SourceChecksum != entry.SourceChecksum ||
				candidate.Target.PlanTermsChecksum != entry.PlanTermsChecksum {
				return fmt.Errorf("%w: cutover row %d differs", ErrConcurrentChange, entry.LegacyMembershipID)
			}
			metadata, err := eventMetadata(entry, manifest.EvidenceChecksum, manifest.PopulationChecksum, manifest.BatchIndex, manifest.BatchCount)
			if err != nil {
				return err
			}
			if err := s.repository.EnsureRollbackEvent(txCtx, Event{
				EventType:      EventTypeCutoverRolledBack,
				SubscriptionID: entry.SubscriptionID,
				PlanVersionID:  entry.PlanVersionID,
				EffectiveAt:    run.ExecutedAt,
				Run:            run,
				Metadata:       metadata,
			}); err != nil {
				return err
			}
			if err := s.repository.MarkRolledBack(txCtx, entry, manifest.EvidenceChecksum, run); err != nil {
				return err
			}
		}
		return nil
	})
}

func validateStrictReadiness(readiness subscriptionbackfill.ReadinessResult) error {
	policy := readiness.Policy
	report := readiness.Report
	if !readiness.Ready || !validChecksum(report.Checksum) ||
		policy.MinimumCoverageRate != 1 || policy.MaximumMismatchRate != 0 ||
		policy.MinimumRollbackEvidenceRate != 1 || policy.MaximumReadErrors != 0 ||
		!policy.RequireNoExistingCutover ||
		report.EligibleLegacy < policy.MinimumEligible ||
		report.LegacyIncomplete > policy.MaximumLegacyIncomplete ||
		report.Observed != report.EligibleLegacy ||
		report.Matched != report.EligibleLegacy || report.Mismatched != 0 ||
		report.NotReady != 0 || report.ReadErrors != 0 || report.CutoverMappings != 0 ||
		report.RollbackEvidence != report.EligibleLegacy || report.MissingRollbackProof != 0 {
		return fmt.Errorf("%w: strict readiness evidence is required", ErrEvidenceInvalid)
	}
	return nil
}

func entryFromCandidate(candidate subscriptionbackfill.Candidate) (Entry, error) {
	if candidate.Mapping.Status != models.CatalogMigrationShadowRead {
		return Entry{}, fmt.Errorf("%w: candidate is not shadow-ready", ErrEvidenceInvalid)
	}
	return entryFromCandidateEvidence(candidate)
}

func entryFromCandidateEvidence(candidate subscriptionbackfill.Candidate) (Entry, error) {
	if candidate.Mapping.SubscriptionID == nil || *candidate.Mapping.SubscriptionID == 0 ||
		candidate.ExistingSubscription == nil {
		return Entry{}, fmt.Errorf("%w: candidate evidence is incomplete", ErrEvidenceInvalid)
	}
	if err := candidate.ExistingSubscription.Validate(); err != nil {
		return Entry{}, fmt.Errorf("%w: %v", ErrEvidenceInvalid, err)
	}
	entry := Entry{
		LegacyMembershipID: candidate.Membership.ID,
		SubscriptionID:     *candidate.Mapping.SubscriptionID,
		LegacyPlanID:       candidate.Membership.PlanID,
		PlanVersionID:      candidate.Target.PlanVersionID,
		SubscriptionKey:    candidate.ExistingSubscription.SubscriptionKey,
		SourceChecksum:     candidate.Mapping.SourceChecksum,
		PlanTermsChecksum:  candidate.Target.PlanTermsChecksum,
		MappingStatus:      models.CatalogMigrationShadowRead,
		PlanMappingStatus:  candidate.Target.MappingStatus,
	}
	return entry, entry.Validate()
}

func manifestChecksum(manifest Manifest) (string, error) {
	payload := struct {
		SchemaVersion           string  `json:"schemaVersion"`
		GeneratedAt             string  `json:"generatedAt"`
		ReadinessReportChecksum string  `json:"readinessReportChecksum"`
		PopulationChecksum      string  `json:"populationChecksum"`
		TotalEntries            int     `json:"totalEntries"`
		BatchIndex              int     `json:"batchIndex"`
		BatchCount              int     `json:"batchCount"`
		Entries                 []Entry `json:"entries"`
	}{
		SchemaVersion: manifest.SchemaVersion, GeneratedAt: manifest.GeneratedAt.UTC().Format(time.RFC3339Nano),
		ReadinessReportChecksum: manifest.ReadinessReportChecksum,
		PopulationChecksum:      manifest.PopulationChecksum, TotalEntries: manifest.TotalEntries,
		BatchIndex: manifest.BatchIndex, BatchCount: manifest.BatchCount, Entries: manifest.Entries,
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func populationChecksum(entries []Entry) (string, error) {
	encoded, err := json.Marshal(entries)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func planChecksum(plan Plan) (string, error) {
	batchChecksums := make([]string, len(plan.Batches))
	for index, batch := range plan.Batches {
		batchChecksums[index] = batch.EvidenceChecksum
	}
	payload := struct {
		SchemaVersion           string   `json:"schemaVersion"`
		GeneratedAt             string   `json:"generatedAt"`
		ReadinessReportChecksum string   `json:"readinessReportChecksum"`
		PopulationChecksum      string   `json:"populationChecksum"`
		TotalEntries            int      `json:"totalEntries"`
		BatchSize               int      `json:"batchSize"`
		BatchChecksums          []string `json:"batchChecksums"`
	}{plan.SchemaVersion, plan.GeneratedAt.UTC().Format(time.RFC3339Nano), plan.ReadinessReportChecksum,
		plan.PopulationChecksum, plan.TotalEntries, plan.BatchSize, batchChecksums}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func rollbackChecksum(manifest RollbackManifest) (string, error) {
	payload := struct {
		SchemaVersion           string  `json:"schemaVersion"`
		GeneratedAt             string  `json:"generatedAt"`
		CutoverEvidenceChecksum string  `json:"cutoverEvidenceChecksum"`
		PopulationChecksum      string  `json:"populationChecksum"`
		TotalEntries            int     `json:"totalEntries"`
		BatchIndex              int     `json:"batchIndex"`
		BatchCount              int     `json:"batchCount"`
		Reason                  string  `json:"reason"`
		Entries                 []Entry `json:"entries"`
	}{manifest.SchemaVersion, manifest.GeneratedAt.UTC().Format(time.RFC3339Nano), manifest.CutoverEvidenceChecksum,
		manifest.PopulationChecksum, manifest.TotalEntries, manifest.BatchIndex, manifest.BatchCount,
		manifest.Reason, manifest.Entries}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func eventMetadata(entry Entry, evidenceChecksum, populationChecksum string, batchIndex, batchCount int) (json.RawMessage, error) {
	return json.Marshal(map[string]any{
		"legacyMembershipId": entry.LegacyMembershipID,
		"subscriptionId":     entry.SubscriptionID,
		"legacyPlanId":       entry.LegacyPlanID,
		"planVersionId":      entry.PlanVersionID,
		"sourceChecksum":     entry.SourceChecksum,
		"planTermsChecksum":  entry.PlanTermsChecksum,
		"evidenceChecksum":   evidenceChecksum,
		"populationChecksum": populationChecksum,
		"batchIndex":         batchIndex,
		"batchCount":         batchCount,
	})
}

func validChecksum(value string) bool {
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil && strings.ToLower(value) == value
}

func safeIdentifier(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > 128 {
		return false
	}
	for _, character := range value {
		if character >= 'a' && character <= 'z' || character >= 'A' && character <= 'Z' || character >= '0' && character <= '9' || strings.ContainsRune("._:-", character) {
			continue
		}
		return false
	}
	return true
}
