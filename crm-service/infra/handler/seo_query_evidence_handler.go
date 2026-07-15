package handler

import (
	_errors "common/errors"
	"context"
	"encoding/json"

	"crm/internal/dto"
	crmpb "pb/types/crm"
)

func (h *SeoStrategyHandler) UpsertSeoQueryEvidence(
	ctx context.Context,
	req *crmpb.UpsertSeoQueryEvidenceRequest,
) (*crmpb.SeoQueryEvidenceResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}

	observedFrom, err := parseOptionalStrategyTime(req.ObservedFrom)
	if err != nil {
		return nil, _errors.ReturnError(
			400,
			"observedFrom không hợp lệ: "+err.Error(),
		)
	}
	observedTo, err := parseOptionalStrategyTime(req.ObservedTo)
	if err != nil {
		return nil, _errors.ReturnError(
			400,
			"observedTo không hợp lệ: "+err.Error(),
		)
	}
	validUntil, err := parseOptionalStrategyTime(req.ValidUntil)
	if err != nil {
		return nil, _errors.ReturnError(
			400,
			"validUntil không hợp lệ: "+err.Error(),
		)
	}

	result, err := h.queryEvidenceUsecase.Upsert(
		ctx,
		&dto.UpsertSeoQueryEvidenceRequest{
			UserNeedID:      req.GetUserNeedId(),
			RawQuery:        req.GetRawQuery(),
			SourceSystem:    req.GetSourceSystem(),
			SourceKind:      req.GetSourceKind(),
			ObservationType: req.GetObservationType(),
			SourceRecordKey: req.GetSourceRecordKey(),
			SourceReference: req.GetSourceReference(),
			Locale:          req.GetLocale(),
			CountryCode:     req.GetCountryCode(),
			DeviceCategory:  req.GetDeviceCategory(),
			ObservedFrom:    observedFrom,
			ObservedTo:      observedTo,
			ValidUntil:      validUntil,
			Impressions:     req.Impressions,
			Clicks:          req.Clicks,
			CTR:             req.Ctr,
			AvgPosition:     req.AvgPosition,
			SearchVolume:    req.SearchVolume,
			Confidence:      req.GetConfidence(),
			Metadata:        json.RawMessage(req.GetMetadataJson()),
		},
	)
	if err != nil {
		return nil, err
	}

	return h.queryEvidenceMapper.ToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoQueryEvidence(
	ctx context.Context,
	req *crmpb.GetSeoQueryEvidenceRequest,
) (*crmpb.SeoQueryEvidenceResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}

	result, err := h.queryEvidenceUsecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return h.queryEvidenceMapper.ToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoQueryEvidenceList(
	ctx context.Context,
	req *crmpb.GetSeoQueryEvidenceListRequest,
) (*crmpb.GetSeoQueryEvidenceListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoQueryEvidenceListRequest{}
	}

	items, total, err := h.queryEvidenceUsecase.List(
		ctx,
		&dto.SeoQueryEvidenceListRequest{
			ObjectiveID:     req.GetObjectiveId(),
			UserNeedID:      req.GetUserNeedId(),
			Keyword:         req.GetKeyword(),
			SourceSystem:    req.GetSourceSystem(),
			SourceKind:      req.GetSourceKind(),
			ObservationType: req.GetObservationType(),
			DeviceCategory:  req.GetDeviceCategory(),
			Confidence:      req.GetConfidence(),
			Status:          req.GetStatus(),
			CurrentOnly:     req.GetCurrentOnly(),
			Page:            req.GetPage(),
			Size:            req.GetSize(),
		},
	)
	if err != nil {
		return nil, err
	}

	response := &crmpb.GetSeoQueryEvidenceListResponse{Total: total}
	for i := range items {
		response.Data = append(
			response.Data,
			h.queryEvidenceMapper.ToPb(&items[i]),
		)
	}
	return response, nil
}

func (h *SeoStrategyHandler) SupersedeSeoQueryEvidence(
	ctx context.Context,
	req *crmpb.SeoQueryEvidenceStatusCommandRequest,
) (*crmpb.SeoQueryEvidenceResponse, error) {
	command, err := queryEvidenceStatusCommandToDTO(req)
	if err != nil {
		return nil, err
	}
	result, err := h.queryEvidenceUsecase.Supersede(ctx, command)
	if err != nil {
		return nil, err
	}
	return h.queryEvidenceMapper.ToPb(result), nil
}

func (h *SeoStrategyHandler) RejectSeoQueryEvidence(
	ctx context.Context,
	req *crmpb.SeoQueryEvidenceStatusCommandRequest,
) (*crmpb.SeoQueryEvidenceResponse, error) {
	command, err := queryEvidenceStatusCommandToDTO(req)
	if err != nil {
		return nil, err
	}
	result, err := h.queryEvidenceUsecase.Reject(ctx, command)
	if err != nil {
		return nil, err
	}
	return h.queryEvidenceMapper.ToPb(result), nil
}

func queryEvidenceStatusCommandToDTO(
	req *crmpb.SeoQueryEvidenceStatusCommandRequest,
) (*dto.SeoQueryEvidenceStatusCommandRequest, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	return &dto.SeoQueryEvidenceStatusCommandRequest{
		ID:             req.GetId(),
		ExpectedStatus: req.GetExpectedStatus(),
	}, nil
}
