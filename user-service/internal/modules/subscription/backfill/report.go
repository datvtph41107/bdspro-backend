package backfill

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"user/internal/models"
	"user/internal/modules/catalog/migration"
	"user/internal/modules/subscription/shadow"
)

/** ShadowComparator evaluates one legacy membership against subscription v2. */
type ShadowComparator interface {
	CompareMembership(
		ctx context.Context,
		membership models.Membership,
		at time.Time,
	) shadow.Entry
}

/** CoverageReport is point-in-time migration evidence. */
type CoverageReport struct {
	AsOf                 time.Time      `json:"asOf"`
	TotalLegacy          int            `json:"totalLegacy"`
	EligibleLegacy       int            `json:"eligibleLegacy"`
	LegacyIncomplete     int            `json:"legacyIncomplete"`
	Observed             int            `json:"observed"`
	Matched              int            `json:"matched"`
	Mismatched           int            `json:"mismatched"`
	NotReady             int            `json:"notReady"`
	ReadErrors           int            `json:"readErrors"`
	CutoverMappings      int            `json:"cutoverMappings"`
	RollbackEvidence     int            `json:"rollbackEvidence"`
	MissingRollbackProof int            `json:"missingRollbackProof"`
	Statuses             map[string]int `json:"statuses"`
	Entries              []shadow.Entry `json:"entries,omitempty"`
	Checksum             string         `json:"checksum"`
}

/** CoverageRate returns observed eligible membership coverage. */
func (r CoverageReport) CoverageRate() float64 {
	if r.EligibleLegacy == 0 {
		return 0
	}
	return float64(r.Observed) / float64(r.EligibleLegacy)
}

/** MismatchRate returns mismatches among observed subscriptions. */
func (r CoverageReport) MismatchRate() float64 {
	if r.Observed == 0 {
		return 0
	}
	return float64(r.Mismatched) / float64(r.Observed)
}

/** RollbackEvidenceRate returns durable proof among observed rows. */
func (r CoverageReport) RollbackEvidenceRate() float64 {
	if r.Observed == 0 {
		return 0
	}
	return float64(r.RollbackEvidence) / float64(r.Observed)
}

/** Auditor builds bounded point-in-time coverage evidence. */
type Auditor struct {
	memberships LegacyMembershipSource
	comparator  ShadowComparator
	events      EventEvidence
}

/** NewAuditor creates a read-only migration auditor. */
func NewAuditor(
	memberships LegacyMembershipSource,
	comparator ShadowComparator,
	events EventEvidence,
) *Auditor {
	return &Auditor{
		memberships: memberships,
		comparator:  comparator,
		events:      events,
	}
}

/** Build scans legacy memberships without changing runtime state. */
func (a *Auditor) Build(
	ctx context.Context,
	asOf time.Time,
	batchSize int,
	includeEntries bool,
) (CoverageReport, error) {
	if a == nil || a.memberships == nil || a.comparator == nil || a.events == nil {
		return CoverageReport{}, fmt.Errorf("subscription backfill auditor is not configured")
	}
	if asOf.IsZero() {
		return CoverageReport{}, fmt.Errorf("coverage as_of is required")
	}
	if batchSize <= 0 || batchSize > 1000 {
		return CoverageReport{}, fmt.Errorf("batch size must be between 1 and 1000")
	}

	report := CoverageReport{
		AsOf:     asOf.UTC(),
		Statuses: make(map[string]int),
	}
	afterID := uint64(0)
	for {
		memberships, err := a.memberships.ListLegacyMemberships(ctx, afterID, batchSize)
		if err != nil {
			return CoverageReport{}, err
		}
		if len(memberships) == 0 {
			break
		}
		for _, membership := range memberships {
			report.TotalLegacy++
			afterID = membership.ID
			entry := a.comparator.CompareMembership(ctx, membership, report.AsOf)
			report.Statuses[string(entry.Status)]++
			if includeEntries {
				report.Entries = append(report.Entries, entry)
			}
			if entry.Status == shadow.StatusLegacyIncomplete {
				report.LegacyIncomplete++
				continue
			}
			report.EligibleLegacy++
			classifyEntry(&report, entry)
			if entry.MappingStatus == models.CatalogMigrationCutover {
				report.CutoverMappings++
			}
			if entry.SubscriptionID == 0 || !observedStatus(entry.Status) {
				continue
			}
			sourceChecksum, err := migration.MembershipSourceChecksum(
				toLegacySnapshot(membership),
			)
			if err != nil {
				report.ReadErrors++
				report.MissingRollbackProof++
				continue
			}
			hasEvent, err := a.events.HasLegacyBackfillEvent(
				ctx,
				entry.SubscriptionID,
				membership.ID,
				sourceChecksum,
			)
			if err != nil {
				report.ReadErrors++
				report.MissingRollbackProof++
				continue
			}
			if hasEvent && entry.Status != shadow.StatusLegacyChanged {
				report.RollbackEvidence++
			} else {
				report.MissingRollbackProof++
			}
		}
		if len(memberships) < batchSize {
			break
		}
	}
	checksum, err := coverageChecksum(report)
	if err != nil {
		return CoverageReport{}, err
	}
	report.Checksum = checksum
	return report, nil
}

func classifyEntry(report *CoverageReport, entry shadow.Entry) {
	switch entry.Status {
	case shadow.StatusMatched:
		report.Observed++
		report.Matched++
	case shadow.StatusLegacyChanged,
		shadow.StatusPlanMappingMismatch,
		shadow.StatusSubscriptionInvalid,
		shadow.StatusSubjectMismatch,
		shadow.StatusStatusMismatch,
		shadow.StatusPeriodMismatch,
		shadow.StatusOrderMismatch,
		shadow.StatusAuditMismatch:
		report.Observed++
		report.Mismatched++
	case shadow.StatusReadError:
		report.ReadErrors++
	case shadow.StatusMappingNotEnabled:
		report.NotReady++
	default:
		report.NotReady++
	}
}

func observedStatus(status shadow.Status) bool {
	switch status {
	case shadow.StatusMatched,
		shadow.StatusLegacyChanged,
		shadow.StatusPlanMappingMismatch,
		shadow.StatusSubscriptionInvalid,
		shadow.StatusSubjectMismatch,
		shadow.StatusStatusMismatch,
		shadow.StatusPeriodMismatch,
		shadow.StatusOrderMismatch,
		shadow.StatusAuditMismatch:
		return true
	default:
		return false
	}
}

/** ReadinessPolicy defines an explicit point-in-time gate. */
type ReadinessPolicy struct {
	MinimumEligible             int     `json:"minimumEligible"`
	MaximumLegacyIncomplete     int     `json:"maximumLegacyIncomplete"`
	MinimumCoverageRate         float64 `json:"minimumCoverageRate"`
	MaximumMismatchRate         float64 `json:"maximumMismatchRate"`
	MinimumRollbackEvidenceRate float64 `json:"minimumRollbackEvidenceRate"`
	MaximumReadErrors           int     `json:"maximumReadErrors"`
	RequireNoExistingCutover    bool    `json:"requireNoExistingCutover"`
}

/** Validate rejects permissive or nonsensical thresholds. */
func (p ReadinessPolicy) Validate() error {
	if p.MinimumEligible <= 0 {
		return fmt.Errorf("minimum eligible sample must be positive")
	}
	for name, value := range map[string]float64{
		"minimum coverage rate":          p.MinimumCoverageRate,
		"maximum mismatch rate":          p.MaximumMismatchRate,
		"minimum rollback evidence rate": p.MinimumRollbackEvidenceRate,
	} {
		if value < 0 || value > 1 {
			return fmt.Errorf("%s must be between 0 and 1", name)
		}
	}
	if p.MaximumLegacyIncomplete < 0 {
		return fmt.Errorf("maximum legacy incomplete cannot be negative")
	}
	if p.MaximumReadErrors < 0 {
		return fmt.Errorf("maximum read errors cannot be negative")
	}
	return nil
}

/** StrictReadinessPolicy requires complete evidence. */
func StrictReadinessPolicy() ReadinessPolicy {
	return ReadinessPolicy{
		MinimumEligible:             1,
		MaximumLegacyIncomplete:     0,
		MinimumCoverageRate:         1,
		MaximumMismatchRate:         0,
		MinimumRollbackEvidenceRate: 1,
		MaximumReadErrors:           0,
		RequireNoExistingCutover:    true,
	}
}

/** ReadinessResult never changes mapping or read ownership. */
type ReadinessResult struct {
	Ready   bool            `json:"ready"`
	Policy  ReadinessPolicy `json:"policy"`
	Report  CoverageReport  `json:"report"`
	Reasons []string        `json:"reasons,omitempty"`
}

/** EvaluateReadiness evaluates evidence without performing cutover. */
func EvaluateReadiness(
	report CoverageReport,
	policy ReadinessPolicy,
) (ReadinessResult, error) {
	if err := policy.Validate(); err != nil {
		return ReadinessResult{}, err
	}
	result := ReadinessResult{Ready: true, Policy: policy, Report: report}
	if report.AsOf.IsZero() {
		result.Reasons = append(result.Reasons, "report as_of is missing")
	}
	if strings.TrimSpace(report.Checksum) == "" {
		result.Reasons = append(result.Reasons, "report checksum is missing")
	}
	if report.EligibleLegacy < policy.MinimumEligible {
		result.Reasons = append(result.Reasons, "eligible sample is below the required minimum")
	}
	if report.LegacyIncomplete > policy.MaximumLegacyIncomplete {
		result.Reasons = append(result.Reasons, "legacy incomplete rows exceed threshold")
	}
	if report.CoverageRate() < policy.MinimumCoverageRate {
		result.Reasons = append(result.Reasons, "shadow coverage is below threshold")
	}
	if report.MismatchRate() > policy.MaximumMismatchRate {
		result.Reasons = append(result.Reasons, "mismatch rate exceeds threshold")
	}
	if report.RollbackEvidenceRate() < policy.MinimumRollbackEvidenceRate {
		result.Reasons = append(result.Reasons, "rollback evidence is below threshold")
	}
	if report.ReadErrors > policy.MaximumReadErrors {
		result.Reasons = append(result.Reasons, "read errors exceed threshold")
	}
	if policy.RequireNoExistingCutover && report.CutoverMappings > 0 {
		result.Reasons = append(result.Reasons, "cutover mappings already exist")
	}
	result.Ready = len(result.Reasons) == 0
	return result, nil
}

func coverageChecksum(report CoverageReport) (string, error) {
	statuses := make([]string, 0, len(report.Statuses))
	for key := range report.Statuses {
		statuses = append(statuses, key)
	}
	sort.Strings(statuses)
	stableStatuses := make([]struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}, 0, len(statuses))
	for _, status := range statuses {
		stableStatuses = append(stableStatuses, struct {
			Status string `json:"status"`
			Count  int    `json:"count"`
		}{Status: status, Count: report.Statuses[status]})
	}
	payload := struct {
		AsOf                 time.Time `json:"asOf"`
		TotalLegacy          int       `json:"totalLegacy"`
		EligibleLegacy       int       `json:"eligibleLegacy"`
		LegacyIncomplete     int       `json:"legacyIncomplete"`
		Observed             int       `json:"observed"`
		Matched              int       `json:"matched"`
		Mismatched           int       `json:"mismatched"`
		NotReady             int       `json:"notReady"`
		ReadErrors           int       `json:"readErrors"`
		CutoverMappings      int       `json:"cutoverMappings"`
		RollbackEvidence     int       `json:"rollbackEvidence"`
		MissingRollbackProof int       `json:"missingRollbackProof"`
		Statuses             any       `json:"statuses"`
	}{
		AsOf: report.AsOf, TotalLegacy: report.TotalLegacy,
		EligibleLegacy: report.EligibleLegacy, LegacyIncomplete: report.LegacyIncomplete,
		Observed: report.Observed, Matched: report.Matched, Mismatched: report.Mismatched,
		NotReady: report.NotReady, ReadErrors: report.ReadErrors,
		CutoverMappings: report.CutoverMappings, RollbackEvidence: report.RollbackEvidence,
		MissingRollbackProof: report.MissingRollbackProof, Statuses: stableStatuses,
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}
