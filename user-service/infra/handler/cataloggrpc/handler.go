package cataloggrpc

import (
	"context"
	"errors"
	"strings"
	"time"

	userpb "pb/types/user"
	"user/internal/modules/catalog/admin"
	"user/internal/modules/catalog/domain"
	"user/internal/modules/catalog/publication"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/** Handler translates the trusted Admin catalog API to the domain service. */
type Handler struct {
	userpb.UnimplementedAdminCatalogServiceServer
	service    *admin.Service
	authorizer admin.PermissionAuthorizer
}

/** NewHandler creates an Admin catalog gRPC handler. */
func NewHandler(service *admin.Service, authorizer admin.PermissionAuthorizer) *Handler {
	return &Handler{service: service, authorizer: authorizer}
}

/** ListPlanVersions returns terms and quota evidence; it never changes policy. */
func (h *Handler) ListPlanVersions(
	ctx context.Context,
	req *userpb.ListCatalogPlanVersionsRequest,
) (*userpb.ListCatalogPlanVersionsResponse, error) {
	if _, err := h.authorize(ctx, admin.PermissionPlanView); err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "catalog control plane is not configured")
	}
	if req == nil {
		req = &userpb.ListCatalogPlanVersionsRequest{}
	}
	page, err := h.service.ListPlanVersions(ctx, admin.Query{
		Page:        req.GetPage(),
		PageSize:    req.GetPageSize(),
		ProductCode: req.GetProductCode(),
		PlanCode:    req.GetPlanCode(),
		Status:      req.GetStatus(),
	})
	if err != nil {
		return nil, serviceError(err)
	}
	response := &userpb.ListCatalogPlanVersionsResponse{
		PlanVersions: make([]*userpb.CatalogPlanVersion, 0, len(page.PlanVersions)),
		Total:        page.Total,
	}
	for _, aggregate := range page.PlanVersions {
		response.PlanVersions = append(response.PlanVersions, toProto(aggregate))
	}
	return response, nil
}

/** GetPlanVersion returns one immutable plan-version contract. */
func (h *Handler) GetPlanVersion(
	ctx context.Context,
	req *userpb.GetCatalogPlanVersionRequest,
) (*userpb.CatalogPlanVersion, error) {
	if _, err := h.authorize(ctx, admin.PermissionPlanView); err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "catalog control plane is not configured")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	aggregate, err := h.service.GetPlanVersion(ctx, req.GetPlanVersionId())
	if err != nil {
		return nil, serviceError(err)
	}
	return toProto(aggregate), nil
}

/** PublishPlanVersion executes the only allowed state transition: draft to active. */
func (h *Handler) PublishPlanVersion(
	ctx context.Context,
	req *userpb.PublishCatalogPlanVersionRequest,
) (*userpb.PublishCatalogPlanVersionResponse, error) {
	actorID, err := h.authorize(ctx, admin.PermissionPlanPublish)
	if err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "catalog control plane is not configured")
	}
	if req == nil || req.GetPlanVersionId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "plan version id is required")
	}
	effectiveFrom, err := requiredTimestamp(req.EffectiveFrom, "effective_from")
	if err != nil {
		return nil, err
	}
	effectiveUntil, err := optionalTimestamp(req.EffectiveUntil, "effective_until")
	if err != nil {
		return nil, err
	}
	result, err := h.service.Publish(ctx, publication.Command{
		PlanVersionID:  req.GetPlanVersionId(),
		ActorID:        actorID,
		EffectiveFrom:  effectiveFrom,
		EffectiveUntil: effectiveUntil,
	})
	if err != nil {
		return nil, serviceError(err)
	}
	aggregate, err := h.service.GetPlanVersion(ctx, result.PlanVersionID)
	if err != nil {
		return nil, serviceError(err)
	}
	return &userpb.PublishCatalogPlanVersionResponse{
		PlanVersion:   toProto(aggregate),
		TermsChecksum: result.TermsChecksum,
		PublishedAt:   timestamppb.New(result.PublishedAt),
	}, nil
}

func (h *Handler) authorize(ctx context.Context, permissionCode string) (uint64, error) {
	var authorizer admin.PermissionAuthorizer
	if h != nil {
		authorizer = h.authorizer
	}
	actorID, err := admin.RequirePermission(ctx, authorizer, permissionCode)
	switch {
	case err == nil:
		return actorID, nil
	case errors.Is(err, admin.ErrUnauthorized):
		return 0, authorizationError(err)
	case errors.Is(err, admin.ErrPermissionDenied):
		return 0, status.Errorf(codes.PermissionDenied, "permission %s is required", permissionCode)
	default:
		return 0, status.Error(codes.Unavailable, "catalog permission authority failed")
	}
}

func requiredTimestamp(value *timestamppb.Timestamp, name string) (time.Time, error) {
	if value == nil {
		return time.Time{}, status.Errorf(codes.InvalidArgument, "%s is required", name)
	}
	if err := value.CheckValid(); err != nil {
		return time.Time{}, status.Errorf(codes.InvalidArgument, "%s is invalid", name)
	}
	return value.AsTime().UTC(), nil
}

func optionalTimestamp(value *timestamppb.Timestamp, name string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	if err := value.CheckValid(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s is invalid", name)
	}
	result := value.AsTime().UTC()
	return &result, nil
}

func authorizationError(err error) error {
	if errors.Is(err, admin.ErrUnauthorized) {
		return status.Error(codes.PermissionDenied, "trusted catalog actor identity is required")
	}
	return status.Error(codes.Internal, "catalog authorization failed")
}

func serviceError(err error) error {
	switch {
	case errors.Is(err, admin.ErrPlanVersionNotFound),
		errors.Is(err, publication.ErrPlanVersionNotFound):
		return status.Error(codes.NotFound, "catalog plan version was not found")
	case errors.Is(err, publication.ErrPlanVersionNotDraft),
		errors.Is(err, publication.ErrParentRetired):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, publication.ErrConcurrentPublish):
		return status.Error(codes.Aborted, "catalog plan version changed; refresh and try again")
	}
	if statusError, ok := status.FromError(err); ok {
		return statusError.Err()
	}
	if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "invalid") ||
		strings.Contains(err.Error(), "page size") {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return status.Error(codes.Internal, "catalog control plane failed")
}

func toProto(aggregate domain.PlanVersionAggregate) *userpb.CatalogPlanVersion {
	response := &userpb.CatalogPlanVersion{
		Id:                 aggregate.ID,
		ProductCode:        aggregate.ProductCode,
		ProductDisplayName: aggregate.ProductDisplayName,
		PlanCode:           aggregate.PlanCode,
		Version:            aggregate.Version,
		DisplayName:        aggregate.DisplayName,
		Status:             string(aggregate.Status),
		SubjectScope:       string(aggregate.SubjectScope),
		TermsChecksum:      aggregate.TermsChecksum,
		Entitlements:       make([]*userpb.CatalogEntitlement, 0, len(aggregate.Entitlements)),
		Operations:         make([]*userpb.CatalogOperationPolicy, 0, len(aggregate.Operations)),
		Prices:             make([]*userpb.CatalogPriceItem, 0, len(aggregate.Prices)),
	}
	if aggregate.EffectiveFrom != nil {
		response.EffectiveFrom = timestamppb.New(*aggregate.EffectiveFrom)
	}
	if aggregate.EffectiveUntil != nil {
		response.EffectiveUntil = timestamppb.New(*aggregate.EffectiveUntil)
	}
	if aggregate.PublishedAt != nil {
		response.PublishedAt = timestamppb.New(*aggregate.PublishedAt)
	}
	for _, entitlement := range aggregate.Entitlements {
		response.Entitlements = append(response.Entitlements, &userpb.CatalogEntitlement{
			Code:        entitlement.Code,
			Kind:        string(entitlement.Kind),
			FeatureCode: entitlement.FeatureCode,
			MeterCode:   entitlement.MeterCode,
			Amount:      entitlement.Amount,
			Unlimited:   entitlement.Unlimited,
			Period:      string(entitlement.Period),
		})
	}
	for _, operation := range aggregate.Operations {
		response.Operations = append(response.Operations, &userpb.CatalogOperationPolicy{
			Code:           operation.Code,
			FeatureCode:    operation.FeatureCode,
			MeterCode:      operation.MeterCode,
			UnitsPerAction: operation.UnitsPerAction,
		})
	}
	for _, price := range aggregate.Prices {
		response.Prices = append(response.Prices, &userpb.CatalogPriceItem{
			Code:        price.Code,
			Kind:        string(price.Kind),
			Currency:    price.Currency,
			AmountMinor: price.AmountMinor,
			BillingUnit: price.BillingUnit,
			MeterCode:   price.MeterCode,
			Quantity:    price.Quantity,
		})
	}
	return response
}
