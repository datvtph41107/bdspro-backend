// Package checkoutpostgres implements User Subscription Checkout persistence.
// It resolves authoritative Plan/current Subscription state and owns the
// durable frozen Checkout snapshot. It does not own Payment transport or
// Subscription mutation policy.
package checkoutpostgres

import (
	"context"
	"errors"
	"fmt"

	_db "common/db"
	"user/internal/usecase/subscription/checkout"

	"gorm.io/gorm"
)

type CheckoutStore struct {
	db *_db.TransactionRepo
}

func NewCheckoutStore(db *_db.TransactionRepo) *CheckoutStore {
	return &CheckoutStore{db: db}
}

var (
	_ checkout.PlanResolver       = (*CheckoutStore)(nil)
	_ checkout.SubscriptionReader = (*CheckoutStore)(nil)
	_ checkout.CommandStore       = (*CheckoutStore)(nil)
	_ checkout.ContactStore       = (*CheckoutStore)(nil)
)

func (s *CheckoutStore) UpdateCheckoutContact(ctx context.Context, profileID uint64, contact checkout.Contact) (checkout.Contact, error) {
	if s == nil || s.db == nil || profileID == 0 || !contact.IsValidInput() {
		return checkout.Contact{}, checkout.ErrInvalidCommand
	}
	db := s.db.GetDB(ctx)
	result := db.Table("user_profile").Where("profile_id = ? AND deleted_at IS NULL", profileID).Updates(map[string]any{
		"full_name":  contact.FullName,
		"email":      contact.Email,
		"updated_at": gorm.Expr("NOW()"),
	})
	if result.Error != nil {
		return checkout.Contact{}, fmt.Errorf("update checkout contact: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return checkout.Contact{}, errors.New("checkout profile not found")
	}
	var saved struct {
		FullName string `gorm:"column:full_name"`
		Email    string `gorm:"column:email"`
		Phone    string `gorm:"column:phone"`
	}
	if err := db.Table("user_profile").Select("full_name, email, phone").Where("profile_id = ? AND deleted_at IS NULL", profileID).Take(&saved).Error; err != nil {
		return checkout.Contact{}, fmt.Errorf("reload checkout contact: %w", err)
	}
	return checkout.Contact{FullName: saved.FullName, Email: saved.Email, Phone: saved.Phone}, nil
}
