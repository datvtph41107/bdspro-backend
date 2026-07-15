package handler

import (
	_errors "common/errors"
	"context"

	"crm/internal/dto"
	crmpb "pb/types/crm"
)

func (h *SeoStrategyHandler) CreateSeoQueryCluster(
	ctx context.Context,
	req *crmpb.CreateSeoQueryClusterRequest,
) (*crmpb.SeoQueryClusterResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	result, err := h.queryClusterUsecase.Create(ctx, &dto.CreateSeoQueryClusterRequest{
		IntentHypothesisID: req.GetIntentHypothesisId(),
		Code:               req.GetCode(),
		Name:               req.GetName(),
	})
	if err != nil {
		return nil, err
	}
	return h.queryClusterMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoQueryCluster(
	ctx context.Context,
	req *crmpb.GetSeoQueryClusterRequest,
) (*crmpb.SeoQueryClusterResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	result, err := h.queryClusterUsecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return h.queryClusterMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoQueryClusterList(
	ctx context.Context,
	req *crmpb.GetSeoQueryClusterListRequest,
) (*crmpb.GetSeoQueryClusterListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoQueryClusterListRequest{}
	}
	items, total, err := h.queryClusterUsecase.List(ctx, &dto.SeoQueryClusterListRequest{
		ObjectiveID:        req.GetObjectiveId(),
		UserNeedID:         req.GetUserNeedId(),
		IntentHypothesisID: req.GetIntentHypothesisId(),
		Keyword:            req.GetKeyword(),
		Status:             req.GetStatus(),
		Confidence:         req.GetConfidence(),
		Page:               req.GetPage(),
		Size:               req.GetSize(),
	})
	if err != nil {
		return nil, err
	}
	response := &crmpb.GetSeoQueryClusterListResponse{Total: total}
	for i := range items {
		response.Data = append(response.Data, h.queryClusterMapper.SummaryToPb(&items[i]))
	}
	return response, nil
}

func (h *SeoStrategyHandler) UpdateSeoQueryCluster(
	ctx context.Context,
	req *crmpb.UpdateSeoQueryClusterRequest,
) (*crmpb.SeoQueryClusterResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	reviewDueAt, err := parseOptionalStrategyTime(req.ReviewDueAt)
	if err != nil {
		return nil, _errors.ReturnError(400, "reviewDueAt không hợp lệ: "+err.Error())
	}
	result, err := h.queryClusterUsecase.Update(ctx, &dto.UpdateSeoQueryClusterRequest{
		ID:              req.GetId(),
		ExpectedVersion: req.GetExpectedVersion(),
		Name:            req.GetName(),
		SearchPromise:   req.GetSearchPromise(),
		InclusionRule:   req.GetInclusionRule(),
		ExclusionRule:   req.GetExclusionRule(),
		Confidence:      req.GetConfidence(),
		OwnerCode:       req.GetOwnerCode(),
		ReviewDueAt:     reviewDueAt,
	})
	if err != nil {
		return nil, err
	}
	return h.queryClusterMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) UpsertSeoQueryClusterMember(
	ctx context.Context,
	req *crmpb.UpsertSeoQueryClusterMemberRequest,
) (*crmpb.SeoQueryClusterResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	result, err := h.queryClusterUsecase.UpsertMember(ctx, &dto.UpsertSeoQueryClusterMemberRequest{
		ClusterID:       req.GetClusterId(),
		QueryEvidenceID: req.GetQueryEvidenceId(),
		ExpectedVersion: req.GetExpectedVersion(),
		MembershipType:  req.GetMembershipType(),
		Rationale:       req.GetRationale(),
	})
	if err != nil {
		return nil, err
	}
	return h.queryClusterMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) RemoveSeoQueryClusterMember(
	ctx context.Context,
	req *crmpb.RemoveSeoQueryClusterMemberRequest,
) (*crmpb.SeoQueryClusterResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	result, err := h.queryClusterUsecase.RemoveMember(ctx, &dto.RemoveSeoQueryClusterMemberRequest{
		ClusterID:       req.GetClusterId(),
		QueryEvidenceID: req.GetQueryEvidenceId(),
		ExpectedVersion: req.GetExpectedVersion(),
	})
	if err != nil {
		return nil, err
	}
	return h.queryClusterMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) EvaluateSeoQueryClusterReadiness(
	ctx context.Context,
	req *crmpb.GetSeoQueryClusterRequest,
) (*crmpb.SeoQueryClusterReadinessResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	result, err := h.queryClusterUsecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return h.queryClusterMapper.ReadinessToPb(result.Readiness), nil
}

func (h *SeoStrategyHandler) ActivateSeoQueryCluster(
	ctx context.Context,
	req *crmpb.SeoQueryClusterCommandRequest,
) (*crmpb.SeoQueryClusterResponse, error) {
	result, err := h.queryClusterUsecase.Activate(ctx, queryClusterCommandToDTO(req))
	if err != nil {
		return nil, err
	}
	return h.queryClusterMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) PauseSeoQueryCluster(
	ctx context.Context,
	req *crmpb.SeoQueryClusterCommandRequest,
) (*crmpb.SeoQueryClusterResponse, error) {
	result, err := h.queryClusterUsecase.Pause(ctx, queryClusterCommandToDTO(req))
	if err != nil {
		return nil, err
	}
	return h.queryClusterMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) ArchiveSeoQueryCluster(
	ctx context.Context,
	req *crmpb.SeoQueryClusterCommandRequest,
) (*crmpb.SeoQueryClusterResponse, error) {
	result, err := h.queryClusterUsecase.Archive(ctx, queryClusterCommandToDTO(req))
	if err != nil {
		return nil, err
	}
	return h.queryClusterMapper.AggregateToPb(result), nil
}

func queryClusterCommandToDTO(req *crmpb.SeoQueryClusterCommandRequest) *dto.SeoQueryClusterCommandRequest {
	if req == nil {
		return nil
	}
	return &dto.SeoQueryClusterCommandRequest{
		ID:              req.GetId(),
		ExpectedVersion: req.GetExpectedVersion(),
		Reason:          req.GetReason(),
	}
}
