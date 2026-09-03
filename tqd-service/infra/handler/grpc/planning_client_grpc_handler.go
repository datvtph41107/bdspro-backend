package handler_grpc

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	_middleware "common/middleware"
	planningdomain "tqd/internal/domain/planningclient/model"
	planningapp "tqd/internal/usecase/planningclient/application"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
	tqdpb "pb/types/tqd"
)

// PlanningClientGrpcHandler owns planning-project business workflows in TQD.
// HTTP method/path/query/body mapping is generated from planning_client.proto
// and registered only by gateway-service.
type PlanningClientGrpcHandler struct {
	tqdpb.UnimplementedPlanningClientServiceServer
	service *planningapp.Service
}

func NewPlanningClientGrpcHandler(service *planningapp.Service) *PlanningClientGrpcHandler {
	return &PlanningClientGrpcHandler{service: service}
}

func bytesResponse(value any) (*wrapperspb.BytesValue, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "encode response: %v", err)
	}
	return wrapperspb.Bytes(payload), nil
}

func setPlanningHTTPMetadata(ctx context.Context, statusCode int, cacheControl string) {
	pairs := []string{
		"x-http-status", strconv.Itoa(statusCode),
		"content-type", "application/json; charset=utf-8",
		"x-content-type-options", "nosniff",
	}
	if strings.TrimSpace(cacheControl) != "" {
		pairs = append(pairs, "cache-control", cacheControl)
	}
	_ = grpc.SetHeader(ctx, metadata.Pairs(pairs...))
}

func parseOptionalPlanningTime(raw, field string, endOfDay bool) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		value, err := time.Parse(layout, raw)
		if err != nil {
			continue
		}
		if endOfDay && layout == "2006-01-02" {
			value = value.Add(24*time.Hour - time.Nanosecond)
		}
		return &value, nil
	}
	return nil, status.Errorf(codes.InvalidArgument, "%s must be RFC3339 or YYYY-MM-DD", field)
}

func authenticatedProfileID(ctx context.Context) (uint64, error) {
	principal, err := _middleware.PrincipalFromContext(ctx)
	if err != nil || principal == nil || principal.ProfileId == 0 {
		return 0, status.Error(codes.Unauthenticated, "authentication is required")
	}
	return principal.ProfileId, nil
}

func repositoryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return status.Error(codes.NotFound, "planning data was not found")
	}
	return status.Errorf(codes.Internal, "planning repository: %v", err)
}

func positiveOr(value, fallback int32) int {
	if value > 0 {
		return int(value)
	}
	return int(fallback)
}

func requiredID(value uint64, name string) (uint64, error) {
	if value == 0 {
		return 0, status.Errorf(codes.InvalidArgument, "%s is required", name)
	}
	return value, nil
}

func (h *PlanningClientGrpcHandler) ListPlanningProjects(ctx context.Context, req *tqdpb.ListPlanningProjectsRequest) (*wrapperspb.BytesValue, error) {
	if req == nil {
		req = &tqdpb.ListPlanningProjectsRequest{}
	}
	updatedFrom, err := parseOptionalPlanningTime(req.GetUpdatedFrom(), "updatedFrom", false)
	if err != nil {
		return nil, err
	}
	updatedTo, err := parseOptionalPlanningTime(req.GetUpdatedTo(), "updatedTo", true)
	if err != nil {
		return nil, err
	}
	var hasLayer *bool
	if req.GetHasLayer() != nil {
		value := req.GetHasLayer().GetValue()
		hasLayer = &value
	}
	filter := planningdomain.ProjectListFilter{
		Page:           int(req.GetPage()),
		Size:           positiveOr(req.GetSize(), 20),
		Search:         strings.TrimSpace(req.GetSearch()),
		PlanningType:   strings.TrimSpace(req.GetPlanningType()),
		PlanningLevel:  strings.TrimSpace(req.GetPlanningLevel()),
		ValidityStatus: strings.TrimSpace(req.GetValidityStatus()),
		ProcessStatus:  strings.TrimSpace(req.GetProcessStatus()),
		JurisdictionID: req.GetJurisdictionId(),
		Area:           strings.TrimSpace(req.GetArea()),
		LegalStatus:    strings.TrimSpace(req.GetLegalStatus()),
		HasLayer:       hasLayer,
		SortBy:         strings.TrimSpace(req.GetSortBy()),
		SortOrder:      strings.TrimSpace(req.GetSortOrder()),
		Decision:       strings.TrimSpace(req.GetDecision()),
		UpdatedFrom:    updatedFrom,
		UpdatedTo:      updatedTo,
	}
	result, err := h.service.ListProjects(ctx, filter)
	if err != nil {
		return nil, repositoryError(err)
	}
	setPlanningHTTPMetadata(ctx, 200, "public, max-age=15, stale-while-revalidate=60")
	return bytesResponse(result)
}

func (h *PlanningClientGrpcHandler) GetPlanningProject(ctx context.Context, req *tqdpb.PlanningProjectIdRequest) (*wrapperspb.BytesValue, error) {
	id, err := requiredID(req.GetProjectId(), "planning project id")
	if err != nil {
		return nil, err
	}
	result, err := h.service.GetProject(ctx, id)
	if err != nil {
		return nil, repositoryError(err)
	}
	setPlanningHTTPMetadata(ctx, 200, "public, max-age=30, stale-while-revalidate=120")
	return bytesResponse(result)
}

func (h *PlanningClientGrpcHandler) ListPlanningDocuments(ctx context.Context, req *tqdpb.ListPlanningDocumentsRequest) (*wrapperspb.BytesValue, error) {
	id, err := requiredID(req.GetProjectId(), "planning project id")
	if err != nil {
		return nil, err
	}
	result, err := h.service.ListDocuments(ctx, id, int(req.GetPage()), positiveOr(req.GetSize(), 50), strings.TrimSpace(req.GetKeyword()), strings.TrimSpace(req.GetDocumentType()), strings.TrimSpace(req.GetStatus()))
	if err != nil {
		return nil, repositoryError(err)
	}
	setPlanningHTTPMetadata(ctx, 200, "public, max-age=30, stale-while-revalidate=120")
	return bytesResponse(result)
}

func (h *PlanningClientGrpcHandler) GetPlanningDocument(ctx context.Context, req *tqdpb.PlanningDocumentIdRequest) (*wrapperspb.BytesValue, error) {
	id, err := requiredID(req.GetDocumentId(), "planning document id")
	if err != nil {
		return nil, err
	}
	result, err := h.service.GetDocument(ctx, id)
	if err != nil {
		return nil, repositoryError(err)
	}
	setPlanningHTTPMetadata(ctx, 200, "public, max-age=30, stale-while-revalidate=120")
	return bytesResponse(result)
}

func (h *PlanningClientGrpcHandler) ListPlanningEvents(ctx context.Context, req *tqdpb.ListPlanningEventsRequest) (*wrapperspb.BytesValue, error) {
	id, err := requiredID(req.GetProjectId(), "planning project id")
	if err != nil {
		return nil, err
	}
	result, err := h.service.ListEvents(ctx, id, int(req.GetPage()), positiveOr(req.GetSize(), 50))
	if err != nil {
		return nil, repositoryError(err)
	}
	setPlanningHTTPMetadata(ctx, 200, "public, max-age=30, stale-while-revalidate=120")
	return bytesResponse(result)
}

func (h *PlanningClientGrpcHandler) FollowPlanningProject(ctx context.Context, req *tqdpb.FollowPlanningProjectRequest) (*wrapperspb.BytesValue, error) {
	userID, err := authenticatedProfileID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := requiredID(req.GetProjectId(), "planning project id")
	if err != nil {
		return nil, err
	}
	result, err := h.service.Follow(ctx, userID, id, strings.TrimSpace(req.GetNote()))
	if err != nil {
		return nil, repositoryError(err)
	}
	setPlanningHTTPMetadata(ctx, 200, "no-store")
	return bytesResponse(map[string]any{"followed": true, "data": result})
}

func (h *PlanningClientGrpcHandler) UnfollowPlanningProject(ctx context.Context, req *tqdpb.PlanningProjectIdRequest) (*wrapperspb.BytesValue, error) {
	userID, err := authenticatedProfileID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := requiredID(req.GetProjectId(), "planning project id")
	if err != nil {
		return nil, err
	}
	if err := h.service.Unfollow(ctx, userID, id); err != nil {
		return nil, repositoryError(err)
	}
	setPlanningHTTPMetadata(ctx, 200, "no-store")
	return bytesResponse(map[string]any{"followed": false, "planningProjectId": id})
}

func (h *PlanningClientGrpcHandler) GetPlanningProjectFollowStatus(ctx context.Context, req *tqdpb.PlanningProjectIdRequest) (*wrapperspb.BytesValue, error) {
	userID, err := authenticatedProfileID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := requiredID(req.GetProjectId(), "planning project id")
	if err != nil {
		return nil, err
	}
	result, err := h.service.GetFollow(ctx, userID, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		setPlanningHTTPMetadata(ctx, 200, "no-store")
		return bytesResponse(map[string]any{"followed": false, "planningProjectId": id})
	}
	if err != nil {
		return nil, repositoryError(err)
	}
	setPlanningHTTPMetadata(ctx, 200, "no-store")
	return bytesResponse(map[string]any{"followed": true, "data": result})
}

func (h *PlanningClientGrpcHandler) ListFollowedPlanningProjects(ctx context.Context, req *tqdpb.ListFollowedPlanningProjectsRequest) (*wrapperspb.BytesValue, error) {
	userID, err := authenticatedProfileID(ctx)
	if err != nil {
		return nil, err
	}
	result, err := h.service.ListFollowed(ctx, userID, int(req.GetPage()), positiveOr(req.GetSize(), 20), strings.TrimSpace(req.GetSearch()))
	if err != nil {
		return nil, repositoryError(err)
	}
	setPlanningHTTPMetadata(ctx, 200, "no-store")
	return bytesResponse(result)
}

func (h *PlanningClientGrpcHandler) GetParcelQuickView(ctx context.Context, req *tqdpb.ParcelQuickViewRequest) (*wrapperspb.BytesValue, error) {
	id, err := requiredID(req.GetParcelId(), "parcel id")
	if err != nil {
		return nil, err
	}
	result, err := h.service.GetParcelQuickView(ctx, id)
	if err != nil {
		return nil, repositoryError(err)
	}
	setPlanningHTTPMetadata(ctx, 200, "public, max-age=30, stale-while-revalidate=120")
	return bytesResponse(result)
}

// GetPlanningProjectProjection is an internal gRPC contract used by CRM to
// compose the public SEO document. It intentionally has no public HTTP binding.
func (h *PlanningClientGrpcHandler) GetPlanningProjectProjection(ctx context.Context, req *wrapperspb.StringValue) (*wrapperspb.BytesValue, error) {
	identity := strings.TrimSpace(req.GetValue())
	if identity == "" {
		return nil, status.Error(codes.InvalidArgument, "planning project identity is required")
	}
	result, err := h.service.GetPlanningProjectProjection(ctx, identity)
	if err != nil {
		return nil, repositoryError(err)
	}
	return bytesResponse(result)
}

func (h *PlanningClientGrpcHandler) ListPlanningProjectLayers(ctx context.Context, req *tqdpb.PlanningProjectIdRequest) (*wrapperspb.BytesValue, error) {
	id, err := requiredID(req.GetProjectId(), "planning project id")
	if err != nil {
		return nil, err
	}
	result, err := h.service.ListProjectLayers(ctx, id)
	if err != nil {
		return nil, repositoryError(err)
	}
	setPlanningHTTPMetadata(ctx, 200, "public, max-age=30, stale-while-revalidate=120")
	return bytesResponse(map[string]any{"data": result, "total": len(result)})
}
