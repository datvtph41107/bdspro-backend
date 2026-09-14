package grpc

import (
	_fault "common/fault"
	"context"
	"errors"
	"time"

	userpb "pb/types/user"
	plan "user/internal/domain/plan"
	"user/internal/usecase/plan/admin"
	"user/internal/usecase/plan/publish"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/** Handler translates the trusted Admin catalog API to the domain service. */
type PlanVersionAdminHandler struct {
	userpb.UnimplementedAdminCatalogServiceServer
	service    *admin.Service
	authorizer admin.PermissionAuthorizer
}

/** NewHandler creates an Admin catalog gRPC handler. */
func NewPlanVersionAdminHandler(service *admin.Service, authorizer admin.PermissionAuthorizer) *PlanVersionAdminHandler {
	return &PlanVersionAdminHandler{service: service, authorizer: authorizer}
}

func (h *PlanVersionAdminHandler) CreatePlanVersionDraft(ctx context.Context, req *userpb.CreateCatalogPlanVersionDraftRequest) (*userpb.CatalogPlanVersion, error) {
	actorID, err := h.authorize(ctx, admin.PermissionPlanManage)
	if err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "plan version administration is not configured")
	}
	command, err := planDraftCommand(actorID, req.GetDraft())
	if err != nil {
		return nil, err
	}
	aggregate, err := h.service.CreateDraft(ctx, command)
	if err != nil {
		return nil, planServiceError(err)
	}
	return planVersionToProto(aggregate), nil
}

/** ListPlanVersions returns terms and quota evidence; it never changes policy. */
func (h *PlanVersionAdminHandler) ListPlanVersions(
	ctx context.Context,
	req *userpb.ListCatalogPlanVersionsRequest,
) (*userpb.ListCatalogPlanVersionsResponse, error) {
	if _, err := h.authorize(ctx, admin.PermissionPlanView); err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "plan version administration is not configured")
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
		return nil, planServiceError(err)
	}
	response := &userpb.ListCatalogPlanVersionsResponse{
		PlanVersions: make([]*userpb.CatalogPlanVersion, 0, len(page.PlanVersions)),
		Total:        page.Total,
	}
	for _, aggregate := range page.PlanVersions {
		response.PlanVersions = append(response.PlanVersions, planVersionToProto(aggregate))
	}
	return response, nil
}

/** GetPlanVersion returns one immutable plan-version contract. */
func (h *PlanVersionAdminHandler) GetPlanVersion(
	ctx context.Context,
	req *userpb.GetCatalogPlanVersionRequest,
) (*userpb.CatalogPlanVersion, error) {
	if _, err := h.authorize(ctx, admin.PermissionPlanView); err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "plan version administration is not configured")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	aggregate, err := h.service.GetPlanVersion(ctx, req.GetPlanVersionId())
	if err != nil {
		return nil, planServiceError(err)
	}
	return planVersionToProto(aggregate), nil
}

func (h *PlanVersionAdminHandler) UpdatePlanVersionDraft(ctx context.Context, req *userpb.UpdateCatalogPlanVersionDraftRequest) (*userpb.CatalogPlanVersion, error) {
	actorID, err := h.authorize(ctx, admin.PermissionPlanManage)
	if err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "plan version administration is not configured")
	}
	if req == nil || req.GetPlanVersionId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "plan version id is required")
	}
	command, err := planDraftCommand(actorID, req.GetDraft())
	if err != nil {
		return nil, err
	}
	aggregate, err := h.service.UpdateDraft(ctx, req.GetPlanVersionId(), command)
	if err != nil {
		return nil, planServiceError(err)
	}
	return planVersionToProto(aggregate), nil
}

func (h *PlanVersionAdminHandler) ValidatePlanVersionDraft(ctx context.Context, req *userpb.ValidateCatalogPlanVersionDraftRequest) (*userpb.ValidateCatalogPlanVersionDraftResponse, error) {
	if _, err := h.authorize(ctx, admin.PermissionPlanManage); err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "plan version administration is not configured")
	}
	if req == nil || req.GetPlanVersionId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "plan version id is required")
	}
	checksum, err := h.service.ValidateDraft(ctx, req.GetPlanVersionId())
	if err != nil {
		return nil, planServiceError(err)
	}
	return &userpb.ValidateCatalogPlanVersionDraftResponse{Valid: true, TermsChecksum: checksum}, nil
}

func (h *PlanVersionAdminHandler) DeletePlanVersionDraft(ctx context.Context, req *userpb.DeleteCatalogPlanVersionDraftRequest) (*emptypb.Empty, error) {
	actorID, err := h.authorize(ctx, admin.PermissionPlanManage)
	if err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "plan version administration is not configured")
	}
	if req == nil || req.GetPlanVersionId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "plan version id is required")
	}
	if err := h.service.DeleteDraft(ctx, req.GetPlanVersionId(), actorID); err != nil {
		return nil, planServiceError(err)
	}
	return &emptypb.Empty{}, nil
}

/** PublishPlanVersion executes the only allowed state transition: draft to active. */
func (h *PlanVersionAdminHandler) PublishPlanVersion(
	ctx context.Context,
	req *userpb.PublishCatalogPlanVersionRequest,
) (*userpb.PublishCatalogPlanVersionResponse, error) {
	actorID, err := h.authorize(ctx, admin.PermissionPlanPublish)
	if err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "plan version administration is not configured")
	}
	if req == nil || req.GetPlanVersionId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "plan version id is required")
	}
	// Publishing immediately is a server-owned clock decision. Requiring the
	// browser to submit "now" makes a valid request stale while it is in flight.
	// An explicit value is still accepted for scheduled publication.
	publishedAt := time.Now().UTC()
	effectiveFrom := publishedAt
	if req.EffectiveFrom != nil {
		effectiveFrom, err = planRequiredTimestamp(req.EffectiveFrom, "effective_from")
		if err != nil {
			return nil, err
		}
	}
	effectiveUntil, err := planOptionalTimestamp(req.EffectiveUntil, "effective_until")
	if err != nil {
		return nil, err
	}
	result, err := h.service.Publish(ctx, publish.Command{
		PlanVersionID:  req.GetPlanVersionId(),
		ActorID:        actorID,
		EffectiveFrom:  effectiveFrom,
		EffectiveUntil: effectiveUntil,
		PublishedAt:    publishedAt,
	})
	if err != nil {
		return nil, planServiceError(err)
	}
	aggregate, err := h.service.GetPlanVersion(ctx, result.PlanVersionID)
	if err != nil {
		return nil, planServiceError(err)
	}
	return &userpb.PublishCatalogPlanVersionResponse{
		PlanVersion:   planVersionToProto(aggregate),
		TermsChecksum: result.TermsChecksum,
		PublishedAt:   timestamppb.New(result.PublishedAt),
	}, nil
}

// RetirePlanVersion stops future sales while preserving historical terms.
func (h *PlanVersionAdminHandler) RetirePlanVersion(ctx context.Context, req *userpb.RetireCatalogPlanVersionRequest) (*userpb.CatalogPlanVersion, error) {
	actorID, err := h.authorize(ctx, admin.PermissionPlanPublish)
	if err != nil {
		return nil, err
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "plan version administration is not configured")
	}
	if req == nil || req.GetPlanVersionId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "plan version id is required")
	}
	aggregate, err := h.service.Retire(ctx, req.GetPlanVersionId(), actorID)
	if err != nil {
		return nil, planServiceError(err)
	}
	return planVersionToProto(aggregate), nil
}

func (h *PlanVersionAdminHandler) authorize(ctx context.Context, permissionCode string) (uint64, error) {
	var authorizer admin.PermissionAuthorizer
	if h != nil {
		authorizer = h.authorizer
	}
	actorID, err := admin.RequirePermission(ctx, authorizer, permissionCode)
	switch {
	case err == nil:
		return actorID, nil
	case errors.Is(err, admin.ErrUnauthorized):
		return 0, planAuthorizationError(err)
	case errors.Is(err, admin.ErrPermissionDenied):
		return 0, status.Errorf(codes.PermissionDenied, "permission %s is required", permissionCode)
	default:
		return 0, status.Error(codes.Unavailable, "plan permission authority failed")
	}
}

func planRequiredTimestamp(value *timestamppb.Timestamp, name string) (time.Time, error) {
	if value == nil {
		return time.Time{}, status.Errorf(codes.InvalidArgument, "%s is required", name)
	}
	if err := value.CheckValid(); err != nil {
		return time.Time{}, status.Errorf(codes.InvalidArgument, "%s is invalid", name)
	}
	return value.AsTime().UTC(), nil
}

func planOptionalTimestamp(value *timestamppb.Timestamp, name string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	if err := value.CheckValid(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s is invalid", name)
	}
	result := value.AsTime().UTC()
	return &result, nil
}

func planAuthorizationError(err error) error {
	if errors.Is(err, admin.ErrUnauthorized) {
		return status.Error(codes.PermissionDenied, "trusted plan actor identity is required")
	}
	return status.Error(codes.Internal, "plan authorization failed")
}

func planServiceError(err error) error {
	switch {
	case errors.Is(err, admin.ErrPlanVersionNotFound),
		errors.Is(err, publish.ErrPlanVersionNotFound):
		return status.Error(codes.NotFound, "plan version was not found")
	case errors.Is(err, admin.ErrPlanVersionConflict):
		return status.Error(codes.AlreadyExists, "plan version already exists")
	case errors.Is(err, admin.ErrPlanVersionNotDraft):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, admin.ErrPlanVersionNotActive),
		errors.Is(err, admin.ErrStablePlanFieldsImmutable):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, publish.ErrPlanVersionNotDraft),
		errors.Is(err, publish.ErrParentRetired):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, publish.ErrConcurrentPublish):
		return status.Error(codes.Aborted, "plan version changed; refresh and try again")
	}
	if _, ok := _fault.As(err); ok {
		return _fault.ToGRPC(err)
	}
	if statusError, ok := status.FromError(err); ok {
		return statusError.Err()
	}
	return status.Error(codes.Internal, "plan version administration failed")
}

func planDraftCommand(actorID uint64, draft *userpb.CatalogPlanVersionDraft) (admin.DraftCommand, error) {
	if draft == nil {
		return admin.DraftCommand{}, status.Error(codes.InvalidArgument, "draft is required")
	}
	command := admin.DraftCommand{
		ActorID:            actorID,
		ProductDisplayName: draft.GetProductDisplayName(),
		TierRank:           draft.GetTierRank(),
		Terms: plan.PlanVersion{
			ProductCode:          draft.GetProductCode(),
			PlanCode:             draft.GetPlanCode(),
			Version:              draft.GetVersion(),
			DisplayName:          draft.GetDisplayName(),
			Status:               plan.StatusDraft,
			SubjectScope:         plan.SubjectScope(draft.GetSubjectScope()),
			SubscriptionTermDays: draft.GetSubscriptionTermDays(),
			Entitlements:         make([]plan.Entitlement, 0, len(draft.GetEntitlements())),
			Prices:               make([]plan.PriceItem, 0, len(draft.GetPrices())),
			Operations:           make([]plan.OperationBinding, 0, len(draft.GetOperations())),
		},
	}
	for _, item := range draft.GetEntitlements() {
		command.Terms.Entitlements = append(command.Terms.Entitlements, plan.Entitlement{Code: item.GetCode(), Kind: plan.EntitlementKind(item.GetKind()), FeatureCode: item.GetFeatureCode(), MeterCode: item.GetMeterCode(), Amount: item.GetAmount(), Unlimited: item.GetUnlimited(), Period: plan.PeriodKind(item.GetPeriod())})
	}
	for _, item := range draft.GetPrices() {
		command.Terms.Prices = append(command.Terms.Prices, plan.PriceItem{Code: item.GetCode(), Kind: plan.PriceKind(item.GetKind()), Currency: item.GetCurrency(), AmountMinor: item.GetAmountMinor(), BillingUnit: item.GetBillingUnit(), MeterCode: item.GetMeterCode(), Quantity: item.GetQuantity()})
	}
	for _, item := range draft.GetOperations() {
		command.Terms.Operations = append(command.Terms.Operations, plan.OperationBinding{Code: item.GetCode(), FeatureCode: item.GetFeatureCode(), MeterCode: item.GetMeterCode(), UnitsPerAction: item.GetUnitsPerAction()})
	}
	return command, nil
}

func planVersionToProto(aggregate plan.PlanVersionAggregate) *userpb.CatalogPlanVersion {
	response := &userpb.CatalogPlanVersion{
		Id:                   aggregate.ID,
		ProductCode:          aggregate.ProductCode,
		ProductDisplayName:   aggregate.ProductDisplayName,
		PlanCode:             aggregate.PlanCode,
		Version:              aggregate.Version,
		DisplayName:          aggregate.DisplayName,
		Status:               string(aggregate.Status),
		SubjectScope:         string(aggregate.SubjectScope),
		TermsChecksum:        aggregate.TermsChecksum,
		TierRank:             aggregate.TierRank,
		SubscriptionTermDays: aggregate.SubscriptionTermDays,
		Entitlements:         make([]*userpb.CatalogEntitlement, 0, len(aggregate.Entitlements)),
		Operations:           make([]*userpb.CatalogOperationPolicy, 0, len(aggregate.Operations)),
		Prices:               make([]*userpb.CatalogPriceItem, 0, len(aggregate.Prices)),
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
