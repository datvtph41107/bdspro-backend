package handler

import (
	_errors "common/errors"
	"context"

	"crm/internal/dto"
	crmpb "pb/types/crm"
)

func (h *SeoStrategyHandler) CreateSeoAudienceSegment(
	ctx context.Context,
	req *crmpb.CreateSeoAudienceSegmentRequest,
) (*crmpb.SeoAudienceSegmentResponse, error) {
	result, err := h.audienceUsecase.Create(
		ctx,
		&dto.CreateSeoAudienceRequest{
			Code: req.GetCode(),
			Name: req.GetName(),
		},
	)
	if err != nil {
		return nil, err
	}

	return h.audienceMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoAudienceSegment(
	ctx context.Context,
	req *crmpb.GetSeoAudienceSegmentRequest,
) (*crmpb.SeoAudienceSegmentResponse, error) {
	result, err := h.audienceUsecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return h.audienceMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoAudienceSegmentList(
	ctx context.Context,
	req *crmpb.GetSeoAudienceSegmentListRequest,
) (*crmpb.GetSeoAudienceSegmentListResponse, error) {
	items, total, err := h.audienceUsecase.List(
		ctx,
		&dto.SeoAudienceListRequest{
			Keyword:      req.GetKeyword(),
			Status:       req.GetStatus(),
			DecisionRole: req.GetDecisionRole(),
			Page:         req.GetPage(),
			Size:         req.GetSize(),
		},
	)
	if err != nil {
		return nil, err
	}

	response := &crmpb.GetSeoAudienceSegmentListResponse{
		Total: total,
	}

	for i := range items {
		response.Data = append(
			response.Data,
			h.audienceMapper.SummaryToPb(&items[i]),
		)
	}

	return response, nil
}

func (h *SeoStrategyHandler) UpdateSeoAudienceSegment(
	ctx context.Context,
	req *crmpb.UpdateSeoAudienceSegmentRequest,
) (*crmpb.SeoAudienceSegmentResponse, error) {
	reviewDueAt, err := parseOptionalStrategyTime(req.ReviewDueAt)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

	result, err := h.audienceUsecase.Update(
		ctx,
		&dto.UpdateSeoAudienceRequest{
			ID:              req.GetId(),
			ExpectedVersion: req.GetExpectedVersion(),
			Name:            req.GetName(),
			Description:     req.GetDescription(),
			DecisionRole:    req.GetDecisionRole(),
			KnowledgeLevel:  req.GetKnowledgeLevel(),
			OwnerCode:       req.GetOwnerCode(),
			ReviewDueAt:     reviewDueAt,
		},
	)
	if err != nil {
		return nil, err
	}

	return h.audienceMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) ActivateSeoAudienceSegment(
	ctx context.Context,
	req *crmpb.SeoAudienceCommandRequest,
) (*crmpb.SeoAudienceSegmentResponse, error) {
	result, err := h.audienceUsecase.Activate(
		ctx,
		audienceCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}

	return h.audienceMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) PauseSeoAudienceSegment(
	ctx context.Context,
	req *crmpb.SeoAudienceCommandRequest,
) (*crmpb.SeoAudienceSegmentResponse, error) {
	result, err := h.audienceUsecase.Pause(
		ctx,
		audienceCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}

	return h.audienceMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) ArchiveSeoAudienceSegment(
	ctx context.Context,
	req *crmpb.SeoAudienceCommandRequest,
) (*crmpb.SeoAudienceSegmentResponse, error) {
	result, err := h.audienceUsecase.Archive(
		ctx,
		audienceCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}

	return h.audienceMapper.AggregateToPb(result), nil
}

func audienceCommandToDTO(
	req *crmpb.SeoAudienceCommandRequest,
) *dto.SeoAudienceCommandRequest {
	return &dto.SeoAudienceCommandRequest{
		ID:              req.GetId(),
		ExpectedVersion: req.GetExpectedVersion(),
		Reason:          req.GetReason(),
	}
}
