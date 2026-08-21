package commercialpostgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"common/productcatalog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"user/internal/modules/commercial/checkout"
)

type checkoutCommandRow struct {
	ID                   uint64                      `gorm:"column:id;primaryKey"`
	SubjectKind          checkout.SubjectKind        `gorm:"column:subject_kind"`
	SubjectID            string                      `gorm:"column:subject_id"`
	CommandKey           string                      `gorm:"column:command_key"`
	ProductID            uint64                      `gorm:"column:product_id"`
	ProductCode          string                      `gorm:"column:product_code"`
	PlanID               uint64                      `gorm:"column:plan_id"`
	PlanCode             string                      `gorm:"column:plan_code"`
	PlanVersionID        uint64                      `gorm:"column:plan_version_id"`
	PlanVersion          string                      `gorm:"column:plan_version"`
	TierRank             int32                       `gorm:"column:tier_rank"`
	SubscriptionTermDays int32                       `gorm:"column:subscription_term_days"`
	TermsChecksum        string                      `gorm:"column:terms_checksum"`
	SubjectScope         productcatalog.SubjectScope `gorm:"column:subject_scope"`
	Currency             string                      `gorm:"column:currency"`
	AmountMinor          int64                       `gorm:"column:amount_minor"`
	CreatedAt            time.Time                   `gorm:"column:created_at"`
}

func (checkoutCommandRow) TableName() string { return "commercial_checkout_commands" }

func (s *CheckoutStore) FindByCommand(
	ctx context.Context,
	subject checkout.Subject,
	commandKey string,
) (checkout.CheckoutCommand, bool, error) {
	if s == nil || s.db == nil || !subject.IsValid() || commandKey == "" {
		return checkout.CheckoutCommand{}, false, checkout.ErrInvalidCommand
	}
	var row checkoutCommandRow
	err := s.db.GetDB(ctx).
		Where("subject_kind = ? AND subject_id = ? AND command_key = ?", subject.Kind, subject.ID, commandKey).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return checkout.CheckoutCommand{}, false, nil
	}
	if err != nil {
		return checkout.CheckoutCommand{}, false, fmt.Errorf("load checkout command: %w", err)
	}
	command := checkoutCommandFromRow(row)
	if !command.IsValid() {
		return checkout.CheckoutCommand{}, false, checkout.ErrCommercialTermsUnavailable
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
	row := checkoutCommandToRow(command)
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
		return checkoutCommandFromRow(row), true, nil
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

func checkoutCommandToRow(command checkout.CheckoutCommand) checkoutCommandRow {
	return checkoutCommandRow{
		SubjectKind:          command.Subject.Kind,
		SubjectID:            command.Subject.ID,
		CommandKey:           command.CommandKey,
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

func checkoutCommandFromRow(row checkoutCommandRow) checkout.CheckoutCommand {
	return checkout.CheckoutCommand{
		Subject:  checkout.Subject{Kind: row.SubjectKind, ID: row.SubjectID},
		PlanCode: row.PlanCode,
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
