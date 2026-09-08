package postgres

import (
	"context"
	"errors"
	"fmt"

	domain "payment/internal/domain/payment"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Store) FindOrderByCommand(ctx context.Context, subject domain.Subject, productCode, commandKey string) (domain.Order, bool, error) {
	if s == nil || s.db == nil {
		return domain.Order{}, false, errors.New("commerce database is not configured")
	}
	var row orderModel
	err := s.db.WithContext(ctx).
		Where("subject_kind = ? AND subject_id = ? AND product_code = ? AND command_key = ?", string(subject.Kind), subject.ID, productCode, commandKey).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Order{}, false, nil
	}
	if err != nil {
		return domain.Order{}, false, fmt.Errorf("find commerce order by command: %w", err)
	}
	return orderFromModel(row), true, nil
}

func (s *Store) FindOrderByReference(ctx context.Context, reference string) (domain.Order, bool, error) {
	if s == nil || s.db == nil {
		return domain.Order{}, false, errors.New("commerce database is not configured")
	}
	var row orderModel
	err := s.db.WithContext(ctx).Where("UPPER(reference) = UPPER(?)", reference).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Order{}, false, nil
	}
	if err != nil {
		return domain.Order{}, false, fmt.Errorf("find commerce order by reference: %w", err)
	}
	return orderFromModel(row), true, nil
}

func (s *Store) CreateOrder(ctx context.Context, order domain.Order) (domain.Order, bool, error) {
	if s == nil || s.db == nil {
		return domain.Order{}, false, errors.New("commerce database is not configured")
	}
	row := orderToModel(order)
	result := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "subject_kind"}, {Name: "subject_id"}, {Name: "product_code"}, {Name: "command_key"},
		},
		DoNothing: true,
	}).Create(&row)
	if result.Error != nil {
		return domain.Order{}, false, fmt.Errorf("create commerce order: %w", result.Error)
	}
	if result.RowsAffected == 1 {
		return orderFromModel(row), true, nil
	}

	existing, found, err := s.FindOrderByCommand(ctx, order.Subject, order.Terms.ProductCode, order.CommandKey)
	if err != nil {
		return domain.Order{}, false, err
	}
	if !found {
		return domain.Order{}, false, errors.New("commerce order conflict winner could not be loaded")
	}
	if existing.CommandFingerprint != order.CommandFingerprint {
		return domain.Order{}, false, domain.ErrOrderCommandConflict
	}
	return existing, false, nil
}
