package commercegrpc

import (
	"context"
	"errors"

	"common/identity"
	payment "payment/internal/domain/payment"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"payment/internal/usecase/attempt"
	"payment/internal/usecase/order"
	paymentpb "pb/types/payment"
	"time"
)

type Handler struct {
	paymentpb.UnimplementedInternalCommerceServiceServer
	orders   *order.Service
	attempts *attempt.Service
}

func NewCommerceHandler(o *order.Service, a *attempt.Service) *Handler {
	return &Handler{orders: o, attempts: a}
}
func requireUser(ctx context.Context) error {
	caller, ok := identity.ServiceCallerFromContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "verified user-service caller required")
	}
	if caller.ServiceID != "user-service" {
		return status.Error(codes.PermissionDenied, "verified user-service caller required")
	}
	return nil
}
func (h *Handler) CreateCommercialOrder(ctx context.Context, r *paymentpb.CreateCommercialOrderRequest) (*paymentpb.CreateCommercialOrderResponse, error) {
	if err := requireUser(ctx); err != nil {
		return nil, err
	}
	if h == nil || h.orders == nil || r == nil {
		return nil, status.Error(codes.FailedPrecondition, "commerce unavailable")
	}
	subj := payment.Subject{Kind: payment.SubjectKind(r.GetSubjectKind()), ID: r.GetSubjectId()}
	terms := payment.CommercialTerms{ProductCode: r.GetProductCode(), PlanCode: r.GetPlanCode(), PlanVersionID: r.GetPlanVersionId(), PlanVersion: r.GetPlanVersion(), TierRank: r.GetTierRank(), SubscriptionTermDays: r.GetSubscriptionTermDays(), TermsChecksum: r.GetTermsChecksum(), Price: payment.Money{Currency: payment.CurrencyCode(r.GetCurrency()), AmountMinor: r.GetAmountMinor()}}
	o, created, err := h.orders.Create(ctx, order.CreateCommand{Subject: subj, Terms: terms, CommandKey: r.GetCommandKey()})
	if err != nil {
		return nil, commerceError(err)
	}
	return &paymentpb.CreateCommercialOrderResponse{OrderId: o.ID, Reference: o.Reference, Status: string(o.Status), Created: created}, nil
}
func (h *Handler) CreatePaymentAttempt(ctx context.Context, r *paymentpb.CreatePaymentAttemptRequest) (*paymentpb.CreatePaymentAttemptResponse, error) {
	if err := requireUser(ctx); err != nil {
		return nil, err
	}
	if h == nil || h.attempts == nil || r == nil {
		return nil, status.Error(codes.FailedPrecondition, "commerce unavailable")
	}
	a, created, err := h.attempts.Create(ctx, attempt.CreateCommand{Subject: payment.Subject{Kind: payment.SubjectKind(r.GetSubjectKind()), ID: r.GetSubjectId()}, OrderID: r.GetOrderId(), Method: payment.PaymentMethod(r.GetMethod()), CommandKey: r.GetCommandKey()})
	if err != nil {
		return nil, commerceError(err)
	}
	n := &paymentpb.PaymentNextAction{Kind: string(a.NextAction.Kind), RedirectUrl: a.NextAction.RedirectURL, QrPayload: a.NextAction.QRPayload}
	exp := ""
	if a.ExpiresAt != nil {
		exp = a.ExpiresAt.UTC().Format(time.RFC3339Nano)
	}
	if a.NextAction.ExpiresAt != nil {
		n.ExpiresAt = a.NextAction.ExpiresAt.UTC().Format(time.RFC3339Nano)
	}
	return &paymentpb.CreatePaymentAttemptResponse{AttemptId: a.ID, OrderId: a.OrderID, Method: string(a.Method), Provider: a.Provider, ProviderReference: a.ProviderReference, Status: string(a.Status), NextAction: n, ExpiresAt: exp, Created: created}, nil
}

// commerceError keeps the wire failure model stable. Callers can retry
// dependency failures and resolve conflicts, while malformed commands never
// look like transient infrastructure errors.
func commerceError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, payment.ErrInvalidCommand), errors.Is(err, payment.ErrInvalidCommercialTerms),
		errors.Is(err, payment.ErrUnsupportedPaymentMethod):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, payment.ErrOrderCommandConflict), errors.Is(err, payment.ErrAttemptCommandConflict),
		errors.Is(err, payment.ErrProviderEvidenceConflict):
		return status.Error(codes.Aborted, err.Error())
	case errors.Is(err, payment.ErrOrderNotFound), errors.Is(err, payment.ErrAttemptNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, payment.ErrAttemptNotAllowed), errors.Is(err, payment.ErrRedriveNotAllowed):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, payment.ErrProviderUnavailable), errors.Is(err, payment.ErrSubscriptionTransient):
		return status.Error(codes.Unavailable, err.Error())
	default:
		return status.Error(codes.Internal, "payment operation failed")
	}
}
