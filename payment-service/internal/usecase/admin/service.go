// Package admin chứa application service cho màn hình vận hành Payment.
// Package này chỉ điều phối projection và command; chính sách thanh toán vẫn
// nằm ở các usecase order/attempt/settlement/fulfillment tương ứng.
package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"

	payment "payment/internal/domain/payment"
)

const (
	PermissionOrderView          = "PAYMENT_ORDER_VIEW"
	PermissionFulfillmentView    = "PAYMENT_FULFILLMENT_VIEW"
	PermissionFulfillmentRedrive = "PAYMENT_FULFILLMENT_REDRIVE"
)

type OrderQuery struct {
	Page        uint32
	PageSize    uint32
	SubjectKind string
	SubjectID   string
	Status      string
	Reference   string
}

type OrderPage struct {
	Orders   []payment.AdminOrder
	Total    uint64
	Page     uint32
	PageSize uint32
}

type FulfillmentQuery struct {
	Page     uint32
	PageSize uint32
	Status   string
	OrderID  uint64
}

type FulfillmentPage struct {
	Fulfillments []payment.AdminFulfillment
	Total        uint64
	Page         uint32
	PageSize     uint32
}

type RedriveCommand struct {
	FulfillmentID uint64
	CommandKey    string
	ActorID       string
	Reason        string
}

// Repository là boundary duy nhất tới Payment commercial projection.
type Repository interface {
	ListAdminOrders(context.Context, OrderQuery) (OrderPage, error)
	GetAdminOrder(context.Context, uint64) (payment.AdminOrder, error)
	ListAdminFulfillments(context.Context, FulfillmentQuery) (FulfillmentPage, error)
	RedriveAdminFulfillment(context.Context, payment.RedriveCommand) (changed bool, replay bool, status payment.FulfillmentStatus, err error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) ListOrders(ctx context.Context, query OrderQuery) (OrderPage, error) {
	if s == nil || s.repository == nil {
		return OrderPage{}, errors.New("payment admin projection is not configured")
	}
	var err error
	query.Page, query.PageSize, err = normalizePage(query.Page, query.PageSize)
	if err != nil {
		return OrderPage{}, err
	}
	query.SubjectKind = strings.TrimSpace(query.SubjectKind)
	query.SubjectID = strings.TrimSpace(query.SubjectID)
	query.Status = strings.TrimSpace(query.Status)
	query.Reference = strings.TrimSpace(query.Reference)
	return s.repository.ListAdminOrders(ctx, query)
}

func (s *Service) GetOrder(ctx context.Context, orderID uint64) (payment.AdminOrder, error) {
	if orderID == 0 {
		return payment.AdminOrder{}, errors.New("order id is required")
	}
	if s == nil || s.repository == nil {
		return payment.AdminOrder{}, errors.New("payment admin projection is not configured")
	}
	return s.repository.GetAdminOrder(ctx, orderID)
}

func (s *Service) ListFulfillments(ctx context.Context, query FulfillmentQuery) (FulfillmentPage, error) {
	if s == nil || s.repository == nil {
		return FulfillmentPage{}, errors.New("payment admin projection is not configured")
	}
	var err error
	query.Page, query.PageSize, err = normalizePage(query.Page, query.PageSize)
	if err != nil {
		return FulfillmentPage{}, err
	}
	query.Status = strings.TrimSpace(query.Status)
	return s.repository.ListAdminFulfillments(ctx, query)
}

func (s *Service) Redrive(ctx context.Context, command RedriveCommand) (bool, bool, payment.FulfillmentStatus, error) {
	if s == nil || s.repository == nil {
		return false, false, "", errors.New("payment admin projection is not configured")
	}
	if command.FulfillmentID == 0 || strings.TrimSpace(command.CommandKey) == "" || strings.TrimSpace(command.ActorID) == "" || strings.TrimSpace(command.Reason) == "" {
		return false, false, "", payment.ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return false, false, "", err
	}
	return s.repository.RedriveAdminFulfillment(ctx, payment.RedriveCommand{
		FulfillmentID: command.FulfillmentID,
		CommandKey:    strings.TrimSpace(command.CommandKey),
		ActorID:       strings.TrimSpace(command.ActorID),
		Reason:        strings.TrimSpace(command.Reason),
	})
}

func normalizePage(page, pageSize uint32) (uint32, uint32, error) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		return 0, 0, fmt.Errorf("page size must be between 1 and 100")
	}
	return page, pageSize, nil
}
