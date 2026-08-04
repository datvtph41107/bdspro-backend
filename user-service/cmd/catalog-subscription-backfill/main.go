package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_db "common/db"
	userpostgres "user/infra/postgres"
	"user/internal/subscriptionbackfill"
	"user/internal/subscriptionshadow"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type options struct {
	action               string
	dsn                  string
	membershipIDs        string
	reviewedBy           uint64
	reviewReason         string
	runID                string
	actorID              string
	requestID            string
	operationID          string
	asOf                 string
	batchSize            int
	output               string
	includeEntries       bool
	minimumEligible      int
	maximumIncomplete    int
	minimumCoverage      float64
	maximumMismatch      float64
	minimumRollback      float64
	maximumReadErrors    int
	allowExistingCutover bool
}

func main() {
	if err := run(context.Background(), parseFlags()); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() options {
	var value options
	flag.StringVar(&value.action, "action", "report", "review|plan|execute|report|readiness")
	flag.StringVar(&value.dsn, "database-dsn", os.Getenv("QHPRO_DATABASE_DSN"), "PostgreSQL DSN")
	flag.StringVar(&value.membershipIDs, "membership-ids", "", "comma-separated legacy membership IDs")
	flag.Uint64Var(&value.reviewedBy, "reviewed-by", 0, "reviewer profile ID")
	flag.StringVar(&value.reviewReason, "review-reason", "", "explicit review reason")
	flag.StringVar(&value.runID, "run-id", "", "stable migration run ID")
	flag.StringVar(&value.actorID, "actor-id", "", "migration actor identity")
	flag.StringVar(&value.requestID, "request-id", "", "optional request correlation ID")
	flag.StringVar(&value.operationID, "operation-id", "", "optional operation correlation ID")
	flag.StringVar(&value.asOf, "as-of", "", "RFC3339 evaluation time; defaults to now")
	flag.IntVar(&value.batchSize, "batch-size", 100, "bounded scan batch size")
	flag.StringVar(&value.output, "output", "", "optional JSON output path")
	flag.BoolVar(&value.includeEntries, "include-entries", false, "include per-membership report entries")
	flag.IntVar(&value.minimumEligible, "minimum-eligible", 1, "minimum eligible membership sample")
	flag.IntVar(&value.maximumIncomplete, "maximum-legacy-incomplete", 0, "allowed incomplete legacy memberships")
	flag.Float64Var(&value.minimumCoverage, "minimum-coverage-rate", 1, "required shadow coverage from 0 to 1")
	flag.Float64Var(&value.maximumMismatch, "maximum-mismatch-rate", 0, "allowed mismatch rate from 0 to 1")
	flag.Float64Var(&value.minimumRollback, "minimum-rollback-evidence-rate", 1, "required rollback evidence from 0 to 1")
	flag.IntVar(&value.maximumReadErrors, "maximum-read-errors", 0, "allowed read errors")
	flag.BoolVar(&value.allowExistingCutover, "allow-existing-cutover", false, "allow readiness with existing cutover mappings")
	flag.Parse()
	return value
}

func run(ctx context.Context, value options) error {
	switch value.action {
	case "review":
		if os.Getenv("QHPRO_ALLOW_SUBSCRIPTION_BACKFILL_REVIEW") != "1" {
			return fmt.Errorf("set QHPRO_ALLOW_SUBSCRIPTION_BACKFILL_REVIEW=1")
		}
	case "execute":
		if os.Getenv("QHPRO_ALLOW_SUBSCRIPTION_BACKFILL_EXECUTE") != "1" {
			return subscriptionbackfill.ErrExecutionForbidden
		}
	case "plan", "report", "readiness":
	default:
		return fmt.Errorf("unsupported action %q", value.action)
	}
	if strings.TrimSpace(value.dsn) == "" {
		return fmt.Errorf("database DSN is required")
	}
	if value.batchSize <= 0 || value.batchSize > 1000 {
		return fmt.Errorf("batch size must be between 1 and 1000")
	}
	asOf, err := parseAsOf(value.asOf)
	if err != nil {
		return err
	}
	db, err := gorm.Open(gormpostgres.Open(value.dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	repository := userpostgres.NewCatalogRepo(_db.NewTransactionRepo(db))
	service := subscriptionbackfill.NewService(repository)
	shadowReader := subscriptionshadow.NewReaderWithConfig(
		repository,
		subscriptionshadow.Config{
			Mode:    subscriptionshadow.ModeLog,
			Timeout: 5 * time.Second,
			Now:     func() time.Time { return asOf },
			Logger:  log.New(io.Discard, "", 0),
		},
	)
	auditor := subscriptionbackfill.NewAuditor(repository, shadowReader, repository)

	switch value.action {
	case "review":
		return review(ctx, service, value, asOf)
	case "plan":
		return plan(ctx, service, value, asOf)
	case "execute":
		return execute(ctx, service, repository, value, asOf)
	case "report":
		report, err := auditor.Build(ctx, asOf, value.batchSize, value.includeEntries)
		if err != nil {
			return err
		}
		return writeJSON(value.output, report)
	case "readiness":
		report, err := auditor.Build(ctx, asOf, value.batchSize, value.includeEntries)
		if err != nil {
			return err
		}
		result, err := subscriptionbackfill.EvaluateReadiness(
			report,
			subscriptionbackfill.ReadinessPolicy{
				MinimumEligible:             value.minimumEligible,
				MaximumLegacyIncomplete:     value.maximumIncomplete,
				MinimumCoverageRate:         value.minimumCoverage,
				MaximumMismatchRate:         value.maximumMismatch,
				MinimumRollbackEvidenceRate: value.minimumRollback,
				MaximumReadErrors:           value.maximumReadErrors,
				RequireNoExistingCutover:    !value.allowExistingCutover,
			},
		)
		if err != nil {
			return err
		}
		if err := writeJSON(value.output, result); err != nil {
			return err
		}
		if !result.Ready {
			return fmt.Errorf("cutover readiness gate failed")
		}
		return nil
	default:
		return fmt.Errorf("unsupported action %q", value.action)
	}
}

func review(
	ctx context.Context,
	service *subscriptionbackfill.Service,
	value options,
	asOf time.Time,
) error {
	ids, err := parseIDs(value.membershipIDs)
	if err != nil {
		return err
	}
	if value.reviewedBy == 0 || strings.TrimSpace(value.reviewReason) == "" {
		return fmt.Errorf("reviewed-by and review-reason are required")
	}
	for _, id := range ids {
		if err := service.Review(ctx, subscriptionbackfill.ReviewRequest{
			LegacyMembershipID: id,
			ReviewedBy:         value.reviewedBy,
			ReviewedAt:         asOf,
			Decision: subscriptionbackfill.ReviewDecision{
				PeriodHandling:          subscriptionbackfill.PeriodPaidAtToFinishedDate,
				SubscriptionKeyStrategy: subscriptionbackfill.KeyStrategyLegacyMembershipID,
				AutoRenew:               false,
				Reason:                  value.reviewReason,
			},
		}); err != nil {
			return fmt.Errorf("review membership %d: %w", id, err)
		}
	}
	return writeJSON(value.output, map[string]any{
		"action": "review", "reviewed": ids, "asOf": asOf,
	})
}

func plan(
	ctx context.Context,
	service *subscriptionbackfill.Service,
	value options,
	asOf time.Time,
) error {
	ids, err := parseIDs(value.membershipIDs)
	if err != nil {
		return err
	}
	results := make([]subscriptionbackfill.Result, 0, len(ids))
	for _, id := range ids {
		result, err := service.Plan(ctx, id, asOf)
		if err != nil {
			return fmt.Errorf("plan membership %d: %w", id, err)
		}
		results = append(results, result)
	}
	return writeJSON(value.output, results)
}

func execute(
	ctx context.Context,
	service *subscriptionbackfill.Service,
	repository subscriptionbackfill.Repository,
	value options,
	asOf time.Time,
) error {
	if strings.TrimSpace(value.runID) == "" || strings.TrimSpace(value.actorID) == "" {
		return fmt.Errorf("run-id and actor-id are required")
	}
	ids, err := parseOptionalIDs(value.membershipIDs)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		ids, err = listAllReviewed(ctx, repository, value.batchSize)
		if err != nil {
			return err
		}
	}
	run := subscriptionbackfill.Run{
		RunID: value.runID, ActorID: value.actorID,
		RequestID: value.requestID, OperationID: value.operationID,
		AsOf: asOf, ExecutedAt: time.Now().UTC(),
	}
	results := make([]subscriptionbackfill.Result, 0, len(ids))
	var executionErrors []error
	for _, id := range ids {
		result, err := service.Execute(ctx, id, run)
		results = append(results, result)
		if err != nil {
			executionErrors = append(executionErrors, fmt.Errorf("membership %d: %w", id, err))
		}
	}
	if err := writeJSON(value.output, results); err != nil {
		return err
	}
	return errors.Join(executionErrors...)
}

func listAllReviewed(
	ctx context.Context,
	repository subscriptionbackfill.Repository,
	batchSize int,
) ([]uint64, error) {
	var all []uint64
	afterID := uint64(0)
	for {
		ids, err := repository.ListReviewedMembershipIDs(ctx, afterID, batchSize)
		if err != nil {
			return nil, err
		}
		if len(ids) == 0 {
			break
		}
		all = append(all, ids...)
		afterID = ids[len(ids)-1]
		if len(ids) < batchSize {
			break
		}
	}
	return all, nil
}

func parseAsOf(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Now().UTC(), nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid as-of: %w", err)
	}
	return parsed.UTC(), nil
}

func parseIDs(value string) ([]uint64, error) {
	ids, err := parseOptionalIDs(value)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("membership-ids is required")
	}
	return ids, nil
}

func parseOptionalIDs(value string) ([]uint64, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	seen := make(map[uint64]struct{})
	ids := make([]uint64, 0)
	for _, raw := range strings.Split(value, ",") {
		id, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
		if err != nil || id == 0 {
			return nil, fmt.Errorf("invalid membership id %q", raw)
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}

func writeJSON(path string, value any) error {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	encoded = append(encoded, '\n')
	if strings.TrimSpace(path) == "" {
		_, err = os.Stdout.Write(encoded)
		return err
	}
	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, encoded, 0o600); err != nil {
		return err
	}
	return os.Rename(temporary, path)
}
