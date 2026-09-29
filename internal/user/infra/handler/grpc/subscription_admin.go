package grpc

import (
	"context"
	"errors"

	userpb "pb/types/user"
	subscriptiondomain "user/internal/domain/subscription"
	subscriptionusecase "user/internal/usecase/subscription"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SubscriptionAdminHandler struct {
	userpb.UnimplementedAdminCommercialServiceServer
	service    *subscriptionusecase.AdminService
	authorizer subscriptionusecase.PermissionAuthorizer
}

func NewSubscriptionAdminHandler(service *subscriptionusecase.AdminService, authorizer subscriptionusecase.PermissionAuthorizer) *SubscriptionAdminHandler {
	return &SubscriptionAdminHandler{service: service, authorizer: authorizer}
}

func (h *SubscriptionAdminHandler) ListSubscriptions(ctx context.Context, req *userpb.ListAdminSubscriptionsRequest) (*userpb.ListAdminSubscriptionsResponse, error) {
	if err := h.authorize(ctx); err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, "subscription projection is unavailable")
	}
	if req == nil {
		req = &userpb.ListAdminSubscriptionsRequest{}
	}
	page, err := h.service.List(ctx, subscriptionusecase.AdminQuery{
		Page: req.GetPage(), PageSize: req.GetPageSize(), ProfileID: req.GetProfileId(),
		Status: req.GetStatus(), ProductCode: req.GetProductCode(), PlanCode: req.GetPlanCode(),
	})
	if err != nil {
		return nil, subscriptionProjectionError(err)
	}
	response := &userpb.ListAdminSubscriptionsResponse{Total: page.Total, Page: page.Page, PageSize: page.PageSize}
	for i := range page.Subscriptions {
		response.Subscriptions = append(response.Subscriptions, subscriptionProjectionToProto(page.Subscriptions[i]))
	}
	return response, nil
}

func (h *SubscriptionAdminHandler) GetUserProjection(ctx context.Context, req *userpb.GetAdminUserCommercialProjectionRequest) (*userpb.AdminUserCommercialProjection, error) {
	if err := h.authorize(ctx); err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, "subscription projection is unavailable")
	}
	if req == nil || req.GetProfileId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "profile id is required")
	}
	projection, err := h.service.GetUser(ctx, req.GetProfileId())
	if err != nil {
		return nil, subscriptionProjectionError(err)
	}
	return subscriptionUserProjectionToProto(projection), nil
}

func (h *SubscriptionAdminHandler) ListUserProjections(ctx context.Context, req *userpb.ListAdminUserCommercialProjectionsRequest) (*userpb.ListAdminUserCommercialProjectionsResponse, error) {
	if err := h.authorize(ctx); err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, "subscription projection is unavailable")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "profile ids are required")
	}
	projections, err := h.service.ListUsers(ctx, req.GetProfileIds())
	if err != nil {
		return nil, subscriptionProjectionError(err)
	}
	response := &userpb.ListAdminUserCommercialProjectionsResponse{Projections: make([]*userpb.AdminUserCommercialProjection, 0, len(projections))}
	for _, projection := range projections {
		response.Projections = append(response.Projections, subscriptionUserProjectionToProto(projection))
	}
	return response, nil
}

func (h *SubscriptionAdminHandler) authorize(ctx context.Context) error {
	var authorizer subscriptionusecase.PermissionAuthorizer
	if h != nil {
		authorizer = h.authorizer
	}
	_, err := subscriptionusecase.RequireViewPermission(ctx, authorizer)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, subscriptionusecase.ErrAdminActorRequired):
		return status.Error(codes.Unauthenticated, "trusted admin actor identity is required")
	case errors.Is(err, subscriptionusecase.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, "COMMERCIAL_SUBSCRIPTION_VIEW permission is required")
	default:
		return status.Error(codes.Unavailable, "commercial permission authority failed")
	}
}

func subscriptionProjectionError(err error) error {
	switch err.Error() {
	case "profile id is required", "between 1 and 100 profile ids are required", "page size must be between 1 and 100":
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "subscription projection failed")
	}
}

func subscriptionUserProjectionToProto(projection subscriptiondomain.AdminUserProjection) *userpb.AdminUserCommercialProjection {
	response := &userpb.AdminUserCommercialProjection{ProfileId: projection.ProfileID, HasSubscription: projection.EffectiveSubscription != nil}
	if projection.EffectiveSubscription != nil {
		response.EffectiveSubscription = subscriptionProjectionToProto(*projection.EffectiveSubscription)
	}
	return response
}

func subscriptionProjectionToProto(item subscriptiondomain.AdminProjection) *userpb.AdminSubscriptionProjection {
	response := &userpb.AdminSubscriptionProjection{
		SubscriptionId: item.SubscriptionID, SubjectKind: item.SubjectKind, SubjectId: item.SubjectID,
		ProductCode: item.ProductCode, ProductDisplayName: item.ProductDisplayName,
		PlanCode: item.PlanCode, PlanDisplayName: item.PlanDisplayName,
		PlanVersionId: item.PlanVersionID, PlanVersion: item.PlanVersion, Status: item.Status,
		OrderReference: item.OrderReference,
	}
	if !item.StartedAt.IsZero() {
		response.StartedAt = timestamppb.New(item.StartedAt)
	}
	if !item.CurrentPeriodStart.IsZero() {
		response.CurrentPeriodStart = timestamppb.New(item.CurrentPeriodStart)
	}
	if !item.CurrentPeriodEnd.IsZero() {
		response.CurrentPeriodEnd = timestamppb.New(item.CurrentPeriodEnd)
	}
	if item.AccessUntil != nil {
		response.AccessUntil = timestamppb.New(*item.AccessUntil)
	}
	if item.PendingPlanVersionID != nil {
		response.PendingPlanVersionId = *item.PendingPlanVersionID
	}
	if item.PendingEffectiveAt != nil {
		response.PendingEffectiveAt = timestamppb.New(*item.PendingEffectiveAt)
	}
	if !item.UpdatedAt.IsZero() {
		response.UpdatedAt = timestamppb.New(item.UpdatedAt)
	}
	return response
}
