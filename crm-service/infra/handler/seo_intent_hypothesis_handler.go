package handler

import (
	_errors "common/errors"
	"context"

	"crm/internal/dto"
	crmpb "pb/types/crm"
)

func (h *SeoStrategyHandler) CreateSeoIntentHypothesis(
	ctx context.Context,
	req *crmpb.CreateSeoIntentHypothesisRequest,
) (*crmpb.SeoIntentHypothesisResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	result, err := h.intentHypothesisUsecase.Create(
		ctx,
		&dto.CreateSeoIntentHypothesisRequest{
			UserNeedID: req.GetUserNeedId(),
			Code:       req.GetCode(),
			Name:       req.GetName(),
		},
	)
	if err != nil {
		return nil, err
	}
	return h.intentHypothesisMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoIntentHypothesis(
	ctx context.Context,
	req *crmpb.GetSeoIntentHypothesisRequest,
) (*crmpb.SeoIntentHypothesisResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	result, err := h.intentHypothesisUsecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return h.intentHypothesisMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoIntentHypothesisList(
	ctx context.Context,
	req *crmpb.GetSeoIntentHypothesisListRequest,
) (*crmpb.GetSeoIntentHypothesisListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoIntentHypothesisListRequest{}
	}
	items, total, err := h.intentHypothesisUsecase.List(
		ctx,
		&dto.SeoIntentHypothesisListRequest{
			ObjectiveID:  req.GetObjectiveId(),
			UserNeedID:   req.GetUserNeedId(),
			Keyword:      req.GetKeyword(),
			SearchAction: req.GetSearchAction(),
			SubjectType:  req.GetSubjectType(),
			ScopeType:    req.GetScopeType(),
			Status:       req.GetStatus(),
			Confidence:   req.GetConfidence(),
			Page:         req.GetPage(),
			Size:         req.GetSize(),
		},
	)
	if err != nil {
		return nil, err
	}
	response := &crmpb.GetSeoIntentHypothesisListResponse{Total: total}
	for i := range items {
		response.Data = append(
			response.Data,
			h.intentHypothesisMapper.SummaryToPb(&items[i]),
		)
	}
	return response, nil
}

func (h *SeoStrategyHandler) UpdateSeoIntentHypothesis(
	ctx context.Context,
	req *crmpb.UpdateSeoIntentHypothesisRequest,
) (*crmpb.SeoIntentHypothesisResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	reviewDueAt, err := parseOptionalStrategyTime(req.ReviewDueAt)
	if err != nil {
		return nil, _errors.ReturnError(400, "reviewDueAt không hợp lệ: "+err.Error())
	}
	result, err := h.intentHypothesisUsecase.Update(
		ctx,
		&dto.UpdateSeoIntentHypothesisRequest{
			ID:              req.GetId(),
			ExpectedVersion: req.GetExpectedVersion(),
			Name:            req.GetName(),
			GoalStatement:   req.GetGoalStatement(),
			PrimaryQuestion: req.GetPrimaryQuestion(),
			ExpectedResult:  req.GetExpectedResult(),
			SearchAction:    req.GetSearchAction(),
			SubjectType:     req.GetSubjectType(),
			ScopeType:       req.GetScopeType(),
			Confidence:      req.GetConfidence(),
			OwnerCode:       req.GetOwnerCode(),
			ReviewDueAt:     reviewDueAt,
		},
	)
	if err != nil {
		return nil, err
	}
	return h.intentHypothesisMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) LinkSeoIntentEvidence(
	ctx context.Context,
	req *crmpb.LinkSeoIntentEvidenceRequest,
) (*crmpb.SeoIntentHypothesisResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	result, err := h.intentHypothesisUsecase.LinkEvidence(
		ctx,
		&dto.LinkSeoIntentEvidenceRequest{
			IntentHypothesisID: req.GetIntentHypothesisId(),
			QueryEvidenceID:    req.GetQueryEvidenceId(),
			RelationType:       req.GetRelationType(),
			RelevanceNote:      req.GetRelevanceNote(),
		},
	)
	if err != nil {
		return nil, err
	}
	return h.intentHypothesisMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) UnlinkSeoIntentEvidence(
	ctx context.Context,
	req *crmpb.UnlinkSeoIntentEvidenceRequest,
) (*crmpb.SeoIntentHypothesisResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	result, err := h.intentHypothesisUsecase.UnlinkEvidence(
		ctx,
		&dto.UnlinkSeoIntentEvidenceRequest{
			IntentHypothesisID: req.GetIntentHypothesisId(),
			QueryEvidenceID:    req.GetQueryEvidenceId(),
		},
	)
	if err != nil {
		return nil, err
	}
	return h.intentHypothesisMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) EvaluateSeoIntentHypothesisReadiness(
	ctx context.Context,
	req *crmpb.GetSeoIntentHypothesisRequest,
) (*crmpb.SeoIntentHypothesisReadinessResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	result, err := h.intentHypothesisUsecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return h.intentHypothesisMapper.ReadinessToPb(result.Readiness), nil
}

func (h *SeoStrategyHandler) ActivateSeoIntentHypothesis(
	ctx context.Context,
	req *crmpb.SeoIntentHypothesisCommandRequest,
) (*crmpb.SeoIntentHypothesisResponse, error) {
	result, err := h.intentHypothesisUsecase.Activate(
		ctx,
		intentHypothesisCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}
	return h.intentHypothesisMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) PauseSeoIntentHypothesis(
	ctx context.Context,
	req *crmpb.SeoIntentHypothesisCommandRequest,
) (*crmpb.SeoIntentHypothesisResponse, error) {
	result, err := h.intentHypothesisUsecase.Pause(
		ctx,
		intentHypothesisCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}
	return h.intentHypothesisMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) ArchiveSeoIntentHypothesis(
	ctx context.Context,
	req *crmpb.SeoIntentHypothesisCommandRequest,
) (*crmpb.SeoIntentHypothesisResponse, error) {
	result, err := h.intentHypothesisUsecase.Archive(
		ctx,
		intentHypothesisCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}
	return h.intentHypothesisMapper.AggregateToPb(result), nil
}

func intentHypothesisCommandToDTO(
	req *crmpb.SeoIntentHypothesisCommandRequest,
) *dto.SeoIntentHypothesisCommandRequest {
	if req == nil {
		return nil
	}
	return &dto.SeoIntentHypothesisCommandRequest{
		ID:              req.GetId(),
		ExpectedVersion: req.GetExpectedVersion(),
		Reason:          req.GetReason(),
	}
}
