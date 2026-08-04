package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	_db "common/db"
	userpostgres "user/infra/postgres"
	"user/internal/subscriptionbackfill"
	"user/internal/subscriptioncutover"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type options struct {
	action      string
	dsn         string
	readiness   string
	manifest    string
	output      string
	runID       string
	actorID     string
	requestID   string
	operationID string
	reason      string
	executedAt  string
	batchSize   int
	batchIndex  int
}

func main() {
	if err := run(context.Background(), parseFlags()); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() options {
	var value options
	flag.StringVar(&value.action, "action", "prepare", "prepare|activate|rollback-plan|rollback")
	flag.StringVar(&value.dsn, "database-dsn", os.Getenv("QHPRO_DATABASE_DSN"), "PostgreSQL DSN")
	flag.StringVar(&value.readiness, "readiness", "", "Task 4.5 readiness JSON")
	flag.StringVar(&value.manifest, "manifest", "", "cutover or rollback manifest JSON")
	flag.StringVar(&value.output, "output", "", "output JSON path")
	flag.StringVar(&value.runID, "run-id", "", "stable operation run ID")
	flag.StringVar(&value.actorID, "actor-id", "", "operator identity")
	flag.StringVar(&value.requestID, "request-id", "", "request correlation ID")
	flag.StringVar(&value.operationID, "operation-id", "", "operation correlation ID")
	flag.StringVar(&value.reason, "reason", "", "approved change or rollback reason")
	flag.StringVar(&value.executedAt, "executed-at", "", "RFC3339 execution time; defaults to now")
	flag.IntVar(&value.batchSize, "batch-size", subscriptioncutover.DefaultBatchEntries, "prepare batch size; max 1000")
	flag.IntVar(&value.batchIndex, "batch-index", 1, "one-based batch index for activate or rollback-plan")
	flag.Parse()
	return value
}

func run(ctx context.Context, value options) error {
	executedAt, err := parseTime(value.executedAt)
	if err != nil {
		return err
	}
	switch value.action {
	case "prepare":
		if strings.TrimSpace(value.readiness) == "" {
			return fmt.Errorf("readiness file is required")
		}
	case "activate":
		if os.Getenv("QHPRO_ALLOW_SUBSCRIPTION_READ_CUTOVER") != "1" {
			return subscriptioncutover.ErrActivationDenied
		}
	case "rollback-plan":
		if strings.TrimSpace(value.reason) == "" {
			return fmt.Errorf("rollback reason is required")
		}
	case "rollback":
		if os.Getenv("QHPRO_ALLOW_SUBSCRIPTION_READ_ROLLBACK") != "1" {
			return subscriptioncutover.ErrRollbackDenied
		}
		if os.Getenv("QHPRO_CONFIRM_SUBSCRIPTION_READ_MODE_LEGACY") != "1" {
			return fmt.Errorf("confirm deployed read mode is legacy before database rollback")
		}
	default:
		return fmt.Errorf("unsupported action %q", value.action)
	}
	if strings.TrimSpace(value.dsn) == "" {
		return fmt.Errorf("database DSN is required")
	}
	db, err := gorm.Open(gormpostgres.Open(value.dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	repository := userpostgres.NewCatalogRepo(_db.NewTransactionRepo(db))
	service := subscriptioncutover.NewService(repository)

	switch value.action {
	case "prepare":
		var readiness subscriptionbackfill.ReadinessResult
		if err := readJSON(value.readiness, &readiness); err != nil {
			return err
		}
		plan, err := service.PreparePlan(ctx, readiness, executedAt, value.batchSize)
		if err != nil {
			return err
		}
		return writeJSON(value.output, plan)
	case "activate":
		var plan subscriptioncutover.Plan
		if err := readJSON(value.manifest, &plan); err != nil {
			return err
		}
		manifest, err := plan.Batch(value.batchIndex)
		if err != nil {
			return err
		}
		if err := service.Activate(ctx, manifest, operationalRun(value, executedAt)); err != nil {
			return err
		}
		return writeJSON(value.output, map[string]any{
			"action": "activate", "planEvidenceChecksum": plan.EvidenceChecksum,
			"batchEvidenceChecksum": manifest.EvidenceChecksum,
			"batchIndex":            manifest.BatchIndex, "batchCount": manifest.BatchCount,
			"entries": len(manifest.Entries), "executedAt": executedAt,
		})
	case "rollback-plan":
		var plan subscriptioncutover.Plan
		if err := readJSON(value.manifest, &plan); err != nil {
			return err
		}
		manifest, err := plan.Batch(value.batchIndex)
		if err != nil {
			return err
		}
		rollback, err := subscriptioncutover.BuildRollbackManifest(manifest, value.reason, executedAt)
		if err != nil {
			return err
		}
		return writeJSON(value.output, rollback)
	case "rollback":
		var manifest subscriptioncutover.RollbackManifest
		if err := readJSON(value.manifest, &manifest); err != nil {
			return err
		}
		if err := service.Rollback(ctx, manifest, operationalRun(value, executedAt)); err != nil {
			return err
		}
		return writeJSON(value.output, map[string]any{
			"action": "rollback", "evidenceChecksum": manifest.EvidenceChecksum,
			"entries": len(manifest.Entries), "executedAt": executedAt,
		})
	default:
		return fmt.Errorf("unsupported action %q", value.action)
	}
}

func operationalRun(value options, executedAt time.Time) subscriptioncutover.Run {
	return subscriptioncutover.Run{
		RunID: value.runID, ActorID: value.actorID,
		RequestID: value.requestID, OperationID: value.operationID,
		ExecutedAt: executedAt, Reason: value.reason,
	}
}

func parseTime(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Now().UTC(), nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid executed-at: %w", err)
	}
	return parsed.UTC(), nil
}

func readJSON(path string, target any) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("manifest file is required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
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
