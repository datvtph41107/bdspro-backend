package checkoutpostgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	plan "user/internal/domain/plan"
	"user/internal/usecase/subscription/checkout"
)

type checkoutSnapshotRow struct {
	ID                   uint64               `gorm:"column:id;primaryKey"`
	SubjectKind          checkout.SubjectKind `gorm:"column:subject_kind"`
	SubjectID            string               `gorm:"column:subject_id"`
	CommandKey           string               `gorm:"column:command_key"`
	RequestHash          string               `gorm:"column:request_hash"`
	ProductID            uint64               `gorm:"column:product_id"`
	ProductCode          string               `gorm:"column:product_code"`
	PlanID               uint64               `gorm:"column:plan_id"`
	PlanCode             string               `gorm:"column:plan_code"`
	PlanVersionID        uint64               `gorm:"column:plan_version_id"`
	PlanVersion          string               `gorm:"column:plan_version"`
	TierRank             int32                `gorm:"column:tier_rank"`
	SubscriptionTermDays int32                `gorm:"column:subscription_term_days"`
	TermsChecksum        string               `gorm:"column:terms_checksum"`
	SubjectScope         plan.SubjectScope    `gorm:"column:subject_scope"`
	Currency             string               `gorm:"column:currency"`
	AmountMinor          int64                `gorm:"column:amount_minor"`
	CreatedAt            time.Time            `gorm:"column:created_at"`
}

func (checkoutSnapshotRow) TableName() string { return "subscription_checkout_snapshots" }

func (s *CheckoutStore) FindByCommand(
	ctx context.Context,
	subject checkout.Subject,
	commandKey string,
) (checkout.CheckoutCommand, bool, error) {
	if s == nil || s.db == nil || !subject.IsValid() || commandKey == "" {
		return checkout.CheckoutCommand{}, false, checkout.ErrInvalidCommand
	}
	var row checkoutSnapshotRow
	err := s.db.GetDB(ctx).
		Where("subject_kind = ? AND subject_id = ? AND command_key = ?", subject.Kind, subject.ID, commandKey).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return checkout.CheckoutCommand{}, false, nil
	}
	if err != nil {
		return checkout.CheckoutCommand{}, false, fmt.Errorf("load checkout command: %w", err)
	}
	command := checkoutSnapshotFromRow(row)
	if !command.IsValid() {
		return checkout.CheckoutCommand{}, false, checkout.ErrPlanTermsUnavailable
	}
	return command, true, nil
}

func (s *CheckoutStore) CreateCommand(
	ctx context.Context,
	command checkout.CheckoutCommand,
) (checkout.CheckoutCommand, bool, error) {
	if s == nil || s.db == nil || !command.IsValid() {
		return checkout.CheckoutCommand{}, false, checkout.ErrInvalidCommand
	}
	row := checkoutSnapshotToRow(command)
	result := s.db.GetDB(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "subject_kind"},
			{Name: "subject_id"},
			{Name: "command_key"},
		},
		DoNothing: true,
	}).Create(&row)
	if result.Error != nil {
		return checkout.CheckoutCommand{}, false, fmt.Errorf("create checkout command: %w", result.Error)
	}
	if result.RowsAffected == 1 {
		return checkoutSnapshotFromRow(row), true, nil
	}

	// Concurrent same-command loser returns the winner's frozen snapshot. The
	// application decides whether the replay planCode matches the original input.
	existing, found, err := s.FindByCommand(ctx, command.Subject, command.CommandKey)
	if err != nil {
		return checkout.CheckoutCommand{}, false, err
	}
	if !found {
		return checkout.CheckoutCommand{}, false, fmt.Errorf("checkout command conflict winner is missing")
	}
	return existing, false, nil
}

func checkoutSnapshotToRow(command checkout.CheckoutCommand) checkoutSnapshotRow {
	return checkoutSnapshotRow{
		SubjectKind:          command.Subject.Kind,
		SubjectID:            command.Subject.ID,
		CommandKey:           command.CommandKey,
		RequestHash:          command.RequestHash,
		ProductID:            command.Terms.ProductID,
		ProductCode:          command.Terms.ProductCode,
		PlanID:               command.Terms.PlanID,
		PlanCode:             command.PlanCode,
		PlanVersionID:        command.Terms.PlanVersionID,
		PlanVersion:          command.Terms.PlanVersion,
		TierRank:             command.Terms.TierRank,
		SubscriptionTermDays: command.Terms.SubscriptionTermDays,
		TermsChecksum:        command.Terms.TermsChecksum,
		SubjectScope:         command.Terms.SubjectScope,
		Currency:             command.Terms.Price.Currency,
		AmountMinor:          command.Terms.Price.AmountMinor,
		CreatedAt:            command.CreatedAt.UTC(),
	}
}

func checkoutSnapshotFromRow(row checkoutSnapshotRow) checkout.CheckoutCommand {
	return checkout.CheckoutCommand{
		Subject:     checkout.Subject{Kind: row.SubjectKind, ID: row.SubjectID},
		PlanCode:    row.PlanCode,
		RequestHash: row.RequestHash,
		Terms: checkout.PlanTerms{
			ProductID:            row.ProductID,
			ProductCode:          row.ProductCode,
			PlanID:               row.PlanID,
			PlanCode:             row.PlanCode,
			PlanVersionID:        row.PlanVersionID,
			PlanVersion:          row.PlanVersion,
			TierRank:             row.TierRank,
			SubscriptionTermDays: row.SubscriptionTermDays,
			TermsChecksum:        row.TermsChecksum,
			SubjectScope:         row.SubjectScope,
			Price: checkout.Money{
				Currency:    row.Currency,
				AmountMinor: row.AmountMinor,
			},
		},
		CommandKey: row.CommandKey,
		CreatedAt:  row.CreatedAt.UTC(),
	}
}
