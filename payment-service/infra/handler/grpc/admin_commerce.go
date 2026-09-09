package commercegrpc

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"common/identity"
	"common/request"
	payment "payment/internal/domain/payment"
	adminusecase "payment/internal/usecase/admin"
	paymentpb "pb/types/payment"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PaymentPermissionAuthorizer là boundary hẹp tới Auth-service. Payment không
// đọc bảng role/permission và không tự suy diễn quyền từ role name.
type PaymentPermissionAuthorizer interface {
	HasPermissions(context.Context, []string) error
}

type AdminCommerceHandler struct {
	paymentpb.UnimplementedAdminCommerceServiceServer
	service    *adminusecase.Service
	authorizer PaymentPermissionAuthorizer
}

func NewAdminCommerceHandler(service *adminusecase.Service, authorizer PaymentPermissionAuthorizer) *AdminCommerceHandler {
	return &AdminCommerceHandler{service: service, authorizer: authorizer}
}

func (h *AdminCommerceHandler) ListPaymentOrders(ctx context.Context, req *paymentpb.ListAdminPaymentOrdersRequest) (*paymentpb.ListAdminPaymentOrdersResponse, error) {
	actorID, err := h.authorize(ctx, adminusecase.PermissionOrderView)
	_ = actorID
	if err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, "payment admin projection is unavailable")
	}
	if req == nil {
		req = &paymentpb.ListAdminPaymentOrdersRequest{}
	}
	page, err := h.service.ListOrders(ctx, adminusecase.OrderQuery{
		Page: req.GetPage(), PageSize: req.GetPageSize(), SubjectKind: req.GetSubjectKind(),
		SubjectID: req.GetSubjectId(), Status: req.GetStatus(), Reference: req.GetReference(),
	})
	if err != nil {
		return nil, adminPaymentError(err)
	}
	response := &paymentpb.ListAdminPaymentOrdersResponse{Total: page.Total, Page: page.Page, PageSize: page.PageSize}
	for _, item := range page.Orders {
		response.Orders = append(response.Orders, adminOrderToProto(item))
	}
	return response, nil
}

func (h *AdminCommerceHandler) GetPaymentOrder(ctx context.Context, req *paymentpb.GetAdminPaymentOrderRequest) (*paymentpb.AdminPaymentOrder, error) {
	if _, err := h.authorize(ctx, adminusecase.PermissionOrderView); err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, "payment admin projection is unavailable")
	}
	if req == nil || req.GetOrderId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "order id is required")
	}
	item, err := h.service.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, adminPaymentError(err)
	}
	return adminOrderToProto(item), nil
}

func (h *AdminCommerceHandler) ConfirmPaymentFunds(ctx context.Context, req *paymentpb.ConfirmAdminPaymentFundsRequest) (*paymentpb.ConfirmAdminPaymentFundsResponse, error) {
	actorID, err := h.authorize(ctx, adminusecase.PermissionFundsConfirm)
	if err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, "payment operation is unavailable")
	}
	if req == nil || req.GetOrderId() == 0 || strings.TrimSpace(req.GetReason()) == "" {
		return nil, status.Error(codes.InvalidArgument, "order id and confirmation reason are required")
	}
	commandKey, ok := request.IdempotencyKeyFromContext(ctx)
	if !ok || strings.TrimSpace(commandKey) == "" {
		return nil, status.Error(codes.InvalidArgument, "Idempotency-Key is required")
	}
	result, err := h.service.ConfirmFunds(ctx, adminusecase.ConfirmFundsCommand{
		OrderID: req.GetOrderId(), CommandKey: commandKey,
		ActorID: strconv.FormatUint(actorID, 10), Reason: req.GetReason(),
	})
	if err != nil {
		return nil, adminPaymentError(err)
	}
	return &paymentpb.ConfirmAdminPaymentFundsResponse{
		Order: adminOrderToProto(result.Order), Changed: result.Changed, Replay: result.Replay,
	}, nil
}

func (h *AdminCommerceHandler) ListFulfillments(ctx context.Context, req *paymentpb.ListAdminFulfillmentsRequest) (*paymentpb.ListAdminFulfillmentsResponse, error) {
	if _, err := h.authorize(ctx, adminusecase.PermissionFulfillmentView); err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, "payment admin projection is unavailable")
	}
	if req == nil {
		req = &paymentpb.ListAdminFulfillmentsRequest{}
	}
	page, err := h.service.ListFulfillments(ctx, adminusecase.FulfillmentQuery{Page: req.GetPage(), PageSize: req.GetPageSize(), Status: req.GetStatus(), OrderID: req.GetOrderId()})
	if err != nil {
		return nil, adminPaymentError(err)
	}
	response := &paymentpb.ListAdminFulfillmentsResponse{Total: page.Total, Page: page.Page, PageSize: page.PageSize}
	for _, item := range page.Fulfillments {
		response.Fulfillments = append(response.Fulfillments, &paymentpb.AdminFulfillmentItem{Fulfillment: fulfillmentToProto(item.Fulfillment), Order: adminOrderToProto(item.Order)})
	}
	return response, nil
}

func (h *AdminCommerceHandler) RedriveFulfillment(ctx context.Context, req *paymentpb.RedriveAdminFulfillmentRequest) (*paymentpb.RedriveAdminFulfillmentResponse, error) {
	actorID, err := h.authorize(ctx, adminusecase.PermissionFulfillmentRedrive)
	if err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, "payment admin projection is unavailable")
	}
	if req == nil || req.GetFulfillmentId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "fulfillment id is required")
	}
	if strings.TrimSpace(req.GetReason()) == "" {
		return nil, status.Error(codes.InvalidArgument, "redrive reason is required")
	}
	commandKey, ok := request.IdempotencyKeyFromContext(ctx)
	if !ok || strings.TrimSpace(commandKey) == "" {
		return nil, status.Error(codes.InvalidArgument, "Idempotency-Key is required")
	}
	changed, replay, current, err := h.service.Redrive(ctx, adminusecase.RedriveCommand{
		FulfillmentID: req.GetFulfillmentId(), CommandKey: commandKey, ActorID: strconv.FormatUint(actorID, 10), Reason: req.GetReason(),
	})
	if err != nil {
		return nil, adminPaymentError(err)
	}
	return &paymentpb.RedriveAdminFulfillmentResponse{FulfillmentId: req.GetFulfillmentId(), Changed: changed, Replay: replay, Status: string(current)}, nil
}

func (h *AdminCommerceHandler) authorize(ctx context.Context, permission string) (uint64, error) {
	if h == nil || h.authorizer == nil {
		return 0, status.Error(codes.Unavailable, "payment permission authority unavailable")
	}
	if _, ok := identity.ServiceCallerFromContext(ctx); !ok {
		return 0, status.Error(codes.Unauthenticated, "trusted gateway caller is required")
	}
	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 || !strings.EqualFold(actor.TokenType, "ACCESS") {
		return 0, status.Error(codes.Unauthenticated, "trusted admin actor identity is required")
	}
	if err := h.authorizer.HasPermissions(ctx, []string{permission}); err != nil {
		return 0, status.Error(codes.PermissionDenied, permission+" permission is required")
	}
	return actor.ProfileID, nil
}

func adminPaymentError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, payment.ErrInvalidCommand), strings.Contains(err.Error(), "page size"):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, payment.ErrOrderNotFound), errors.Is(err, payment.ErrFulfillmentNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, payment.ErrRedriveNotAllowed):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, payment.ErrFundsConfirmationNotAllowed):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, payment.ErrCommandEvidenceConflict), errors.Is(err, payment.ErrProviderEvidenceConflict):
		return status.Error(codes.Aborted, err.Error())
	default:
		return status.Error(codes.Internal, "payment admin projection failed")
	}
}

func adminOrderToProto(item payment.AdminOrder) *paymentpb.AdminPaymentOrder {
	o := item.Order
	response := &paymentpb.AdminPaymentOrder{
		OrderId: o.ID, Reference: o.Reference, SubjectKind: string(o.Subject.Kind), SubjectId: o.Subject.ID,
		ProductCode: o.Terms.ProductCode, PlanCode: o.Terms.PlanCode, PlanVersionId: o.Terms.PlanVersionID,
		PlanVersion: o.Terms.PlanVersion, TierRank: o.Terms.TierRank, SubscriptionTermDays: o.Terms.SubscriptionTermDays,
		TermsChecksum: o.Terms.TermsChecksum, Currency: string(o.Terms.Price.Currency), AmountMinor: o.Terms.Price.AmountMinor,
		Status: string(o.Status), ExpiresAt: timestampOrNil(o.ExpiresAt), FundsConfirmedAt: timestampPtrOrNil(o.FundsConfirmedAt),
		CreatedAt: timestampOrNil(o.CreatedAt), UpdatedAt: timestampOrNil(o.UpdatedAt),
	}
	for _, attempt := range item.Attempts {
		response.Attempts = append(response.Attempts, attemptToProto(attempt))
	}
	if item.Settlement != nil {
		response.Settlement = settlementToProto(*item.Settlement)
	}
	if item.Fulfillment != nil {
		response.Fulfillment = fulfillmentToProto(*item.Fulfillment)
	}
	if item.ManualConfirmation != nil {
		m := item.ManualConfirmation
		response.ManualConfirmation = &paymentpb.AdminManualPaymentConfirmation{
			ActorId: m.ActorID, Reason: m.Reason, CommandKey: m.CommandKey,
			Outcome: m.Outcome, CreatedAt: timestampOrNil(m.CreatedAt),
		}
	}
	return response
}

func attemptToProto(a payment.PaymentAttempt) *paymentpb.AdminPaymentAttempt {
	return &paymentpb.AdminPaymentAttempt{AttemptId: a.ID, OrderId: a.OrderID, CommandKey: a.CommandKey, Method: string(a.Method), Provider: a.Provider, ProviderReference: a.ProviderReference, Status: string(a.Status), NextActionKind: string(a.NextAction.Kind), FailureCode: a.FailureCode, FailureDetail: a.FailureDetail, ExpiresAt: timestampPtrOrNil(a.ExpiresAt), CreatedAt: timestampOrNil(a.CreatedAt), UpdatedAt: timestampOrNil(a.UpdatedAt)}
}

func settlementToProto(s payment.Settlement) *paymentpb.AdminSettlement {
	return &paymentpb.AdminSettlement{SettlementId: s.ID, OrderId: s.OrderID, Provider: s.Provider, ProviderTransactionId: s.ProviderTransactionID, Reference: s.Reference, Currency: string(s.Amount.Currency), AmountMinor: s.Amount.AmountMinor, OccurredAt: timestampOrNil(s.OccurredAt), Status: string(s.Status), ReviewReason: s.ReviewReason, CreatedAt: timestampOrNil(s.CreatedAt)}
}

func fulfillmentToProto(f payment.Fulfillment) *paymentpb.AdminFulfillment {
	return &paymentpb.AdminFulfillment{FulfillmentId: f.ID, OrderId: f.OrderID, Status: string(f.Status), AttemptCount: int32(f.AttemptCount), AvailableAt: timestampOrNil(f.AvailableAt), LockedBy: f.LockedBy, ClaimVersion: f.ClaimVersion, LeaseUntil: zeroTimeNil(f.LeaseUntil), LastError: f.LastError, CompletedAt: timestampPtrOrNil(f.CompletedAt), CreatedAt: timestampOrNil(f.CreatedAt), UpdatedAt: timestampOrNil(f.UpdatedAt)}
}

func timestampOrNil(value time.Time) *timestamppb.Timestamp {
	if value.IsZero() {
		return nil
	}
	return timestamppb.New(value.UTC())
}

func timestampPtrOrNil(value *time.Time) *timestamppb.Timestamp {
	if value == nil || value.IsZero() {
		return nil
	}
	return timestamppb.New(value.UTC())
}

func zeroTimeNil(value time.Time) *timestamppb.Timestamp { return timestampOrNil(value) }
