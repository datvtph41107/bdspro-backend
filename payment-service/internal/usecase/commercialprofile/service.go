// Package commercialprofile owns the authenticated customer's read model for
// Payment commercial progress. It does not mutate payment state or decide
// subscription/entitlement truth.
package commercialprofile

import (
	"context"
	"errors"
	"strconv"

	payment "payment/internal/domain/payment"
)

var (
	ErrUnavailable = errors.New("payment profile is unavailable")
	ErrInvalidView = errors.New("payment order view is invalid")
	ErrForbidden   = errors.New("payment order does not belong to the authenticated subject")
)

const (
	ProgressAwaitingPayment = "AWAITING_PAYMENT"
	ProgressActivating      = "ACTIVATING"
	ProgressActive          = "ACTIVE"
	ProgressActionRequired  = "ACTION_REQUIRED"
)

type Viewer struct {
	ProfileID      uint64
	OrganizationID *uint64
}

type Repository interface {
	GetAdminOrder(context.Context, uint64) (payment.AdminOrder, error)
}

type Projection struct {
	Order               payment.Order
	LatestAttemptStatus string
	SettlementStatus    string
	FulfillmentStatus   string
	ProgressState       string
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) Get(ctx context.Context, viewer Viewer, orderID uint64) (Projection, error) {
	if s == nil || s.repository == nil {
		return Projection{}, ErrUnavailable
	}
	if orderID == 0 || viewer.ProfileID == 0 {
		return Projection{}, ErrInvalidView
	}
	if err := ctx.Err(); err != nil {
		return Projection{}, err
	}
	item, err := s.repository.GetAdminOrder(ctx, orderID)
	if err != nil {
		return Projection{}, err
	}
	if !viewerOwns(viewer, item.Order.Subject) {
		return Projection{}, ErrForbidden
	}

	projection := Projection{Order: item.Order, ProgressState: ProgressAwaitingPayment}
	if len(item.Attempts) > 0 {
		projection.LatestAttemptStatus = string(item.Attempts[0].Status)
	}
	if item.Settlement != nil {
		projection.SettlementStatus = string(item.Settlement.Status)
	}
	if item.Fulfillment != nil {
		projection.FulfillmentStatus = string(item.Fulfillment.Status)
	}
	projection.ProgressState = progress(item)
	return projection, nil
}

func viewerOwns(viewer Viewer, subject payment.Subject) bool {
	switch subject.Kind {
	case payment.SubjectProfile:
		return subject.ID == strconv.FormatUint(viewer.ProfileID, 10)
	case payment.SubjectOrganization:
		return viewer.OrganizationID != nil && *viewer.OrganizationID != 0 &&
			subject.ID == strconv.FormatUint(*viewer.OrganizationID, 10)
	default:
		return false
	}
}

func progress(item payment.AdminOrder) string {
	if item.Order.Status == payment.OrderRequiresReview ||
		(item.Settlement != nil && item.Settlement.Status == payment.SettlementRequiresReview) ||
		(item.Fulfillment != nil && item.Fulfillment.Status == payment.FulfillmentRequiresReview) {
		return ProgressActionRequired
	}
	if item.Fulfillment != nil && item.Fulfillment.Status == payment.FulfillmentCompleted {
		return ProgressActive
	}
	if item.Order.Status == payment.OrderFundsConfirmed || item.Settlement != nil || item.Fulfillment != nil {
		return ProgressActivating
	}
	return ProgressAwaitingPayment
}
