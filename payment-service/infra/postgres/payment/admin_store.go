package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	payment "payment/internal/domain/payment"
	adminusecase "payment/internal/usecase/admin"

	"gorm.io/gorm"
)

var _ adminusecase.Repository = (*Store)(nil)

func (s *Store) ListAdminOrders(ctx context.Context, query adminusecase.OrderQuery) (adminusecase.OrderPage, error) {
	if s == nil || s.db == nil {
		return adminusecase.OrderPage{}, errors.New("commerce database is not configured")
	}
	var total int64
	base := s.db.WithContext(ctx).Model(&orderModel{})
	base = applyOrderFilters(base, query)
	if err := base.Count(&total).Error; err != nil {
		return adminusecase.OrderPage{}, fmt.Errorf("count commerce orders: %w", err)
	}
	var rows []orderModel
	page := int(query.Page)
	pageSize := int(query.PageSize)
	if err := base.Order("created_at DESC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		return adminusecase.OrderPage{}, fmt.Errorf("list commerce orders: %w", err)
	}
	orders := make([]payment.AdminOrder, 0, len(rows))
	for _, row := range rows {
		projection, err := s.loadAdminOrder(ctx, row)
		if err != nil {
			return adminusecase.OrderPage{}, err
		}
		orders = append(orders, projection)
	}
	return adminusecase.OrderPage{Orders: orders, Total: uint64(maxInt64(total)), Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *Store) GetAdminOrder(ctx context.Context, orderID uint64) (payment.AdminOrder, error) {
	if orderID == 0 || s == nil || s.db == nil {
		return payment.AdminOrder{}, payment.ErrInvalidCommand
	}
	var row orderModel
	if err := s.db.WithContext(ctx).First(&row, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return payment.AdminOrder{}, payment.ErrOrderNotFound
		}
		return payment.AdminOrder{}, fmt.Errorf("get commerce order: %w", err)
	}
	return s.loadAdminOrder(ctx, row)
}

func (s *Store) ListAdminFulfillments(ctx context.Context, query adminusecase.FulfillmentQuery) (adminusecase.FulfillmentPage, error) {
	if s == nil || s.db == nil {
		return adminusecase.FulfillmentPage{}, errors.New("commerce database is not configured")
	}
	base := s.db.WithContext(ctx).Model(&fulfillmentModel{})
	if query.Status != "" {
		base = base.Where("status = ?", query.Status)
	}
	if query.OrderID != 0 {
		base = base.Where("order_id = ?", query.OrderID)
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return adminusecase.FulfillmentPage{}, fmt.Errorf("count commerce fulfillments: %w", err)
	}
	var rows []fulfillmentModel
	page, pageSize := int(query.Page), int(query.PageSize)
	if err := base.Order("updated_at DESC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		return adminusecase.FulfillmentPage{}, fmt.Errorf("list commerce fulfillments: %w", err)
	}
	items := make([]payment.AdminFulfillment, 0, len(rows))
	for _, row := range rows {
		var orderRow orderModel
		if err := s.db.WithContext(ctx).First(&orderRow, row.OrderID).Error; err != nil {
			return adminusecase.FulfillmentPage{}, fmt.Errorf("load order for fulfillment %d: %w", row.ID, err)
		}
		order, err := s.loadAdminOrder(ctx, orderRow)
		if err != nil {
			return adminusecase.FulfillmentPage{}, err
		}
		items = append(items, payment.AdminFulfillment{Fulfillment: fulfillmentFromModel(row), Order: order})
	}
	return adminusecase.FulfillmentPage{Fulfillments: items, Total: uint64(maxInt64(total)), Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *Store) RedriveAdminFulfillment(ctx context.Context, command payment.RedriveCommand) (bool, bool, payment.FulfillmentStatus, error) {
	changed, replay, err := s.RedriveFulfillment(ctx, command, time.Now().UTC())
	if err != nil {
		return false, replay, "", err
	}
	var row fulfillmentModel
	if loadErr := s.db.WithContext(ctx).First(&row, command.FulfillmentID).Error; loadErr != nil {
		return changed, replay, "", fmt.Errorf("reload fulfillment after redrive: %w", loadErr)
	}
	return changed, replay, payment.FulfillmentStatus(row.Status), nil
}

func (s *Store) loadAdminOrder(ctx context.Context, row orderModel) (payment.AdminOrder, error) {
	projection := payment.AdminOrder{Order: orderFromModel(row), Attempts: []payment.PaymentAttempt{}}
	var attempts []paymentAttemptModel
	if err := s.db.WithContext(ctx).Where("order_id = ?", row.ID).Order("created_at DESC, id DESC").Find(&attempts).Error; err != nil {
		return payment.AdminOrder{}, fmt.Errorf("load payment attempts for order %d: %w", row.ID, err)
	}
	for _, attempt := range attempts {
		projection.Attempts = append(projection.Attempts, paymentAttemptFromModel(attempt))
	}
	var settlement settlementModel
	err := s.db.WithContext(ctx).Where("order_id = ?", row.ID).Order("occurred_at DESC, id DESC").First(&settlement).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Unmatched evidence chưa có order_id nhưng vẫn phải hiển thị theo reference.
		err = s.db.WithContext(ctx).Where("reference = ?", row.Reference).Order("occurred_at DESC, id DESC").First(&settlement).Error
	}
	if err == nil {
		value := settlementFromModel(settlement)
		projection.Settlement = &value
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return payment.AdminOrder{}, fmt.Errorf("load settlement for order %d: %w", row.ID, err)
	}
	var manual commandEffectModel
	err = s.db.WithContext(ctx).
		Where("effect_type = ? AND scope_id = ?", payment.CommandEffectManualFundsConfirmation, row.ID).
		Order("created_at DESC, id DESC").First(&manual).Error
	if err == nil {
		projection.ManualConfirmation = &payment.CommandEffect{
			EffectType: manual.EffectType, ScopeID: manual.ScopeID, CommandKey: manual.CommandKey,
			ActorID: manual.ActorID, Reason: manual.Reason, Outcome: manual.Outcome, CreatedAt: manual.CreatedAt,
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return payment.AdminOrder{}, fmt.Errorf("load manual confirmation for order %d: %w", row.ID, err)
	}
	var fulfillment fulfillmentModel
	err = s.db.WithContext(ctx).Where("order_id = ?", row.ID).First(&fulfillment).Error
	if err == nil {
		value := fulfillmentFromModel(fulfillment)
		projection.Fulfillment = &value
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return payment.AdminOrder{}, fmt.Errorf("load fulfillment for order %d: %w", row.ID, err)
	}
	return projection, nil
}

func applyOrderFilters(db *gorm.DB, query adminusecase.OrderQuery) *gorm.DB {
	if value := strings.TrimSpace(query.SubjectKind); value != "" {
		db = db.Where("subject_kind = ?", value)
	}
	if value := strings.TrimSpace(query.SubjectID); value != "" {
		db = db.Where("subject_id = ?", value)
	}
	if value := strings.TrimSpace(query.Status); value != "" {
		db = db.Where("status = ?", value)
	}
	if value := strings.TrimSpace(query.Reference); value != "" {
		db = db.Where("reference ILIKE ?", "%"+value+"%")
	}
	return db
}

func maxInt64(value int64) int64 {
	if value < 0 {
		return 0
	}
	return value
}
