package commercialprofilegrpc

import (
	"context"
	"errors"
	"strings"

	"common/identity"
	payment "payment/internal/domain/payment"
	commercialprofile "payment/internal/usecase/commercialprofile"
	paymentpb "pb/types/payment"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	paymentpb.UnimplementedPaymentProfileServiceServer
	service *commercialprofile.Service
}

func New(service *commercialprofile.Service) *Handler { return &Handler{service: service} }

func (h *Handler) GetMyPaymentOrder(ctx context.Context, req *paymentpb.GetMyPaymentOrderRequest) (*paymentpb.MyPaymentOrder, error) {
	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 || !strings.EqualFold(actor.TokenType, "ACCESS") {
		return nil, status.Error(codes.Unauthenticated, "authenticated profile is required")
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, commercialprofile.ErrUnavailable.Error())
	}
	if req == nil || req.GetOrderId() == 0 {
		return nil, status.Error(codes.InvalidArgument, commercialprofile.ErrInvalidView.Error())
	}
	projection, err := h.service.Get(ctx, commercialprofile.Viewer{
		ProfileID: actor.ProfileID, OrganizationID: actor.OrganizationID,
	}, req.GetOrderId())
	if err != nil {
		return nil, profileError(err)
	}
	return toProto(projection), nil
}

func profileError(err error) error {
	switch {
	case errors.Is(err, commercialprofile.ErrInvalidView):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, commercialprofile.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, payment.ErrOrderNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "payment order view canceled")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "payment order view timed out")
	default:
		return status.Error(codes.Unavailable, commercialprofile.ErrUnavailable.Error())
	}
}

func toProto(projection commercialprofile.Projection) *paymentpb.MyPaymentOrder {
	order := projection.Order
	response := &paymentpb.MyPaymentOrder{
		OrderId:             order.ID,
		Reference:           order.Reference,
		SubjectKind:         string(order.Subject.Kind),
		SubjectId:           order.Subject.ID,
		ProductCode:         order.Terms.ProductCode,
		PlanCode:            order.Terms.PlanCode,
		Currency:            string(order.Terms.Price.Currency),
		AmountMinor:         order.Terms.Price.AmountMinor,
		OrderStatus:         string(order.Status),
		LatestAttemptStatus: projection.LatestAttemptStatus,
		SettlementStatus:    projection.SettlementStatus,
		FulfillmentStatus:   projection.FulfillmentStatus,
		ProgressState:       projection.ProgressState,
	}
	if !order.ExpiresAt.IsZero() {
		response.ExpiresAt = timestamppb.New(order.ExpiresAt)
	}
	if order.FundsConfirmedAt != nil && !order.FundsConfirmedAt.IsZero() {
		response.FundsConfirmedAt = timestamppb.New(*order.FundsConfirmedAt)
	}
	if !order.UpdatedAt.IsZero() {
		response.UpdatedAt = timestamppb.New(order.UpdatedAt)
	}
	return response
}
