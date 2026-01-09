package handler

import (
	"context"
	"fmt"
	"time"

	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	bdspropb "pb/types/bdspro"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DistributionHandler struct {
	bdspropb.UnimplementedDistributionServiceServer
	usecase *usecases.DistributionUsecase
	mapper  *mapper.DistributionMapper
}

func NewDistributionHandler(usecase *usecases.DistributionUsecase, mapper *mapper.DistributionMapper) *DistributionHandler {
	return &DistributionHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

func (h *DistributionHandler) CreateDistribute(ctx context.Context, req *bdspropb.DistributeRequest) (*bdspropb.DistributeResponse, error) {
	dtoReq := &dto.DistributeDTO{
		ProductID:  req.ProductId,
		PartnerIDs: req.PartnerIds,
		Policy:     h.mapper.ProtoToDistributionPolicy(req.Policy),
		FromDate:   timestamppbToTimePtr(req.FromDate),
		ToDate:     timestamppbToTimePtr(req.ToDate),
	}

	result, err := h.usecase.CreateDistribute(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create distribute: %v", err)
	}

	return &bdspropb.DistributeResponse{
		CreatedCount: int32(result.CreatedCount),
		FromDate:     req.FromDate,
		ToDate:       req.ToDate,
	}, nil
}

func (h *DistributionHandler) UpdateDistribute(ctx context.Context, req *bdspropb.UpdateDistributeRequest) (*bdspropb.DistributeResponse, error) {
	dtoReq := &dto.DistributeDTO{
		ProductID:  req.Id,
		PartnerIDs: req.PartnerIds,
		Policy:     h.mapper.ProtoToDistributionPolicy(req.Policy),
		FromDate:   timestamppbToTimePtr(req.FromDate),
		ToDate:     timestamppbToTimePtr(req.ToDate),
	}

	result, err := h.usecase.UpdateDistribute(ctx, req.Id, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update distribute: %v", err)
	}

	return &bdspropb.DistributeResponse{
		CreatedCount: int32(result.CreatedCount),
		FromDate:     req.FromDate,
		ToDate:       req.ToDate,
	}, nil
}

func (h *DistributionHandler) GetDistribution(ctx context.Context, req *bdspropb.GetDistributionRequest) (*bdspropb.DistributionEntity, error) {
	entity, err := h.usecase.GetDistributeWithPartners(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "distribution not found: %v", err)
	}
	return h.mapper.DistributionEntityToProto(entity), nil
}

// ---------------- List ----------------
func (h *DistributionHandler) ListDistributions(ctx context.Context, req *bdspropb.ListDistributionsRequest) (*bdspropb.ListDistributionsResponse, error) {
	entities, err := h.usecase.GetListDistributes(ctx, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list distributions: %v", err)
	}

	return &bdspropb.ListDistributionsResponse{
		Distributions: h.mapper.DistributionListToProto(entities),
	}, nil
}

func (h *DistributionHandler) RevokeDistribute(ctx context.Context, req *bdspropb.RevokeDistributeRequest) (*bdspropb.RevokeDistributeResponse, error) {
	revoke := dto.RevokeDTO{
		Reason:      uint32(req.Reason),
		ReasonOther: req.ReasonOther,
	}

	if err := h.usecase.RevokeDistribute(ctx, req.DistributeId, revoke); err != nil {
		fmt.Printf("Failed to revoke distribute %d: %v\n", req.DistributeId, err)
		return &bdspropb.RevokeDistributeResponse{Success: false}, status.Errorf(codes.Internal, "failed to revoke distribute: %v", err)
	}

	return &bdspropb.RevokeDistributeResponse{Success: true}, nil
}

func timestamppbToTimePtr(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}
