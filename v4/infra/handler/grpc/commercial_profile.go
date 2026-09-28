package grpc

import (
	"context"
	"errors"
	"strings"

	"common/identity"
	userpb "pb/types/user"
	"userv4/internal/usecase/commercialprofile"
	cancelapp "userv4/internal/usecase/subscription/cancel"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommercialProfileHandler struct {
	userpb.UnimplementedCommercialProfileServiceServer
	service      *commercialprofile.Service
	cancellation *cancelapp.Service
}

func NewCommercialProfileHandler(service *commercialprofile.Service, cancellation *cancelapp.Service) *CommercialProfileHandler {
	return &CommercialProfileHandler{service: service, cancellation: cancellation}
}

func (h *CommercialProfileHandler) CancelMySubscription(ctx context.Context, req *userpb.CancelMySubscriptionRequest) (*userpb.CancelMySubscriptionResponse, error) {
	if h == nil || h.cancellation == nil {
		return nil, status.Error(codes.Unavailable, "subscription cancellation is unavailable")
	}
	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 || !strings.EqualFold(actor.TokenType, "ACCESS") {
		return nil, status.Error(codes.Unauthenticated, "authenticated profile is required")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "cancellation request is required")
	}
	result, err := h.cancellation.Cancel(ctx, actor.ProfileID, req.GetProductCode(), req.GetReason())
	if err != nil {
		switch {
		case errors.Is(err, cancelapp.ErrInvalid):
			return nil, status.Error(codes.InvalidArgument, "product code and cancellation reason are required")
		case errors.Is(err, cancelapp.ErrNotFound):
			return nil, status.Error(codes.FailedPrecondition, "active paid subscription was not found")
		case errors.Is(err, cancelapp.ErrAmbiguous):
			return nil, status.Error(codes.FailedPrecondition, "multiple active paid subscriptions require operator review")
		case errors.Is(err, context.Canceled):
			return nil, status.Error(codes.Canceled, "subscription cancellation canceled")
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, "subscription cancellation timed out")
		default:
			return nil, status.Error(codes.Internal, "subscription cancellation failed")
		}
	}
	return &userpb.CancelMySubscriptionResponse{
		Subscription: subscriptionProjectionToProto(result.Subscription),
		Changed:      result.Changed,
	}, nil
}

func (h *CommercialProfileHandler) ListAvailablePlans(ctx context.Context, req *userpb.ListAvailablePlansRequest) (*userpb.ListAvailablePlansResponse, error) {
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, "commercial catalog is unavailable")
	}
	if req == nil {
		req = &userpb.ListAvailablePlansRequest{}
	}
	plans, err := h.service.ListAvailablePlans(ctx, req.GetProductCode(), req.GetSubjectKind())
	if err != nil {
		switch {
		case errors.Is(err, commercialprofile.ErrAmbiguousActiveTerms):
			return nil, status.Error(codes.FailedPrecondition, "commercial catalog has overlapping active terms")
		case err.Error() == "invalid subject kind":
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Unavailable, "commercial catalog is unavailable")
		}
	}
	response := &userpb.ListAvailablePlansResponse{Plans: make([]*userpb.CatalogPlanVersion, 0, len(plans))}
	for _, item := range plans {
		response.Plans = append(response.Plans, planVersionToProto(item))
	}
	return response, nil
}

func (h *CommercialProfileHandler) GetMyCommercialProfile(ctx context.Context, _ *userpb.GetMyCommercialProfileRequest) (*userpb.MyCommercialProfile, error) {
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, "commercial profile is unavailable")
	}
	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 || !strings.EqualFold(actor.TokenType, "ACCESS") {
		return nil, status.Error(codes.Unauthenticated, "authenticated profile is required")
	}
	projection, err := h.service.GetProfile(ctx, actor.ProfileID)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "commercial profile is unavailable")
	}
	response := &userpb.MyCommercialProfile{ProfileId: actor.ProfileID, HasSubscription: projection.EffectiveSubscription != nil}
	if projection.EffectiveSubscription != nil {
		response.EffectiveSubscription = subscriptionProjectionToProto(*projection.EffectiveSubscription)
	}
	return response, nil
}
