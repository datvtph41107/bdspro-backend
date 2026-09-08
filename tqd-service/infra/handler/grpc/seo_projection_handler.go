package handler_grpc

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	tqdpb "pb/types/tqd"
	publiccontent_postgres "tqd/infra/postgres/publiccontent"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const seoProjectionTimeout = 15 * time.Second

// SeoProjectionGrpcHandler exposes authoritative TQD source projections to
// internal backend consumers. Public/browser traffic is owned by Gateway;
// CRM and Gateway consume this projection exclusively over gRPC.
type SeoProjectionGrpcHandler struct {
	tqdpb.UnimplementedTqdSeoProjectionServiceServer
	repository *publiccontent_postgres.Repository
}

func NewSeoProjectionGrpcHandler(
	repository *publiccontent_postgres.Repository,
) *SeoProjectionGrpcHandler {
	return &SeoProjectionGrpcHandler{repository: repository}
}

func (h *SeoProjectionGrpcHandler) GetAdministrativeUnitProjection(
	ctx context.Context,
	req *wrapperspb.StringValue,
) (*wrapperspb.BytesValue, error) {
	if h == nil || h.repository == nil {
		return nil, status.Error(codes.Unavailable, "administrative unit projection repository is unavailable")
	}

	identity := ""
	if req != nil {
		identity = strings.TrimSpace(req.Value)
	}
	if identity == "" {
		return nil, status.Error(codes.InvalidArgument, "administrative unit identity is required")
	}

	queryCtx := ctx
	cancel := func() {}
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		queryCtx, cancel = context.WithTimeout(ctx, seoProjectionTimeout)
	}
	defer cancel()

	projection, err := h.repository.GetAdministrativeUnitProjection(queryCtx, identity)
	if err != nil {
		if publiccontent_postgres.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "administrative unit was not found")
		}
		if queryCtx.Err() != nil {
			return nil, status.Error(codes.DeadlineExceeded, "administrative unit projection timed out")
		}
		return nil, status.Errorf(codes.Internal, "load administrative unit projection: %v", err)
	}

	payload, err := json.Marshal(projection)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "encode administrative unit projection: %v", err)
	}
	setPlanningHTTPMetadata(ctx, 200, "public, max-age=30, stale-while-revalidate=120")
	return wrapperspb.Bytes(payload), nil
}

func (h *SeoProjectionGrpcHandler) GetPlanningProjectProjection(
	ctx context.Context,
	req *wrapperspb.StringValue,
) (*wrapperspb.BytesValue, error) {
	identity := ""
	if req != nil {
		identity = strings.TrimSpace(req.Value)
	}
	payload, err := h.loadPlanningProjectionJSON(ctx, identity)
	if err != nil {
		return nil, err
	}
	return wrapperspb.Bytes(payload), nil
}

func (h *SeoProjectionGrpcHandler) GetParcelQuickViewPublic(
	ctx context.Context,
	req *tqdpb.GetParcelQuickViewRequest,
) (*httpbody.HttpBody, error) {
	if h == nil || h.repository == nil {
		return nil, status.Error(codes.Unavailable, "parcel quick-view repository is unavailable")
	}
	parcelID := uint64(0)
	if req != nil {
		parcelID = req.ParcelId
	}
	if parcelID == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel id must be a positive integer")
	}

	result, err := h.repository.GetParcelQuickView(ctx, parcelID)
	if err != nil {
		if publiccontent_postgres.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "parcel was not found")
		}
		return nil, status.Errorf(codes.Internal, "load parcel quick-view: %v", err)
	}
	payload, err := json.Marshal(result)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "encode parcel quick-view: %v", err)
	}
	return &httpbody.HttpBody{
		ContentType: "application/json; charset=utf-8",
		Data:        payload,
	}, nil
}

func (h *SeoProjectionGrpcHandler) GetPlanningProjectProjectionPublic(
	ctx context.Context,
	req *tqdpb.GetPlanningProjectProjectionRequest,
) (*httpbody.HttpBody, error) {
	identity := ""
	if req != nil {
		identity = strings.TrimSpace(req.Identity)
	}
	payload, err := h.loadPlanningProjectionJSON(ctx, identity)
	if err != nil {
		return nil, err
	}
	return &httpbody.HttpBody{
		ContentType: "application/json; charset=utf-8",
		Data:        payload,
	}, nil
}

func (h *SeoProjectionGrpcHandler) loadPlanningProjectionJSON(ctx context.Context, identity string) ([]byte, error) {
	if h == nil || h.repository == nil {
		return nil, status.Error(codes.Unavailable, "planning project projection repository is unavailable")
	}
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return nil, status.Error(codes.InvalidArgument, "planning project identity is required")
	}

	queryCtx := ctx
	cancel := func() {}
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		queryCtx, cancel = context.WithTimeout(ctx, seoProjectionTimeout)
	}
	defer cancel()

	projection, err := h.repository.GetPlanningProjectProjection(queryCtx, identity)
	if err != nil {
		if publiccontent_postgres.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "planning project was not found")
		}
		if queryCtx.Err() != nil {
			return nil, status.Error(codes.DeadlineExceeded, "planning project projection timed out")
		}
		return nil, status.Errorf(codes.Internal, "load planning project projection: %v", err)
	}

	payload, err := json.Marshal(projection)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "encode planning project projection: %v", err)
	}
	setPlanningHTTPMetadata(ctx, 200, "public, max-age=30, stale-while-revalidate=120")
	return payload, nil
}
