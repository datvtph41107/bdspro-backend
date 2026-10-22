package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	"common/httpapi"
	_utils "common/utils"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/enums"
	"mdn4/modules/bdspro/mapper"
	"mdn4/modules/bdspro/services"
	bdspropb "pb/types/bdspro"
)

// DistributeController owns /v2/bdspro/distributions/* (migrated from bdspro-service).
type DistributeController struct {
	svc    services.DistributionService
	mapper *mapper.DistributionMapper
	sync   *_utils.SyncUtil
}

func NewDistributeController(
	svc services.DistributionService,
	mapper *mapper.DistributionMapper,
	sync *_utils.SyncUtil,
) *DistributeController {
	return &DistributeController{svc: svc, mapper: mapper, sync: sync}
}

func (c *DistributeController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.HandleProto(mux, http.MethodPost, "/v2/bdspro/distributions", c.createProto)
	httpapi.HandleProto(mux, http.MethodPut, "/v2/bdspro/distributions/{id}", c.updateProto)
	httpapi.HandleProto(mux, http.MethodGet, "/v2/bdspro/distributions/{id}", c.getProto)
	httpapi.HandleProto(mux, http.MethodGet, "/v2/bdspro/distributions/all/{productId}", c.listProto)
	httpapi.HandleProto(mux, http.MethodPost, "/v2/bdspro/distributions/revoke", c.revokeProto)
	httpapi.HandleProto(mux, http.MethodDelete, "/v2/bdspro/distributions/{distributeId}/users/{originProfileId}", c.revokeUserProto)
}

func (c *DistributeController) createProto(ctx context.Context, req *bdspropb.DistributeRequest) (*bdspropb.DistributeResponse, error) {
	dtoReq := &dto.DistributeDTO{
		ProductID:       req.ProductId,
		PartnerIDs:      req.PartnerIds,
		FromDate:        _utils.ParseStringToTime(req.FromDate),
		ToDate:          _utils.ParseStringToTime(req.ToDate),
		Price:           req.Price,
		CanDeal:         req.CanDeal,
		ChannelPrice:    req.ChannelPrice,
		CommissionType:  enums.CommissionType(req.CommissionType),
		CommissionValue: req.CommissionValue,
		Note:            req.Note,
	}
	result, err := c.svc.CreateDistribute(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	return distributeResultToProto(result), nil
}

func (c *DistributeController) updateProto(ctx context.Context, req *bdspropb.UpdateDistributeRequest) (*bdspropb.DistributeResponse, error) {
	dtoReq := &dto.DistributeDTO{
		PartnerIDs:      req.PartnerIds,
		FromDate:        _utils.ParseStringToTime(req.FromDate),
		ToDate:          _utils.ParseStringToTime(req.ToDate),
		Price:           req.Price,
		CanDeal:         req.CanDeal,
		ChannelPrice:    req.ChannelPrice,
		CommissionType:  enums.CommissionType(req.CommissionType),
		CommissionValue: req.CommissionValue,
		Note:            req.Note,
	}
	result, err := c.svc.UpdateDistribute(ctx, req.Id, dtoReq)
	if err != nil {
		return nil, err
	}
	return distributeResultToProto(result), nil
}

func (c *DistributeController) getProto(ctx context.Context, req *bdspropb.GetDistributionRequest) (*bdspropb.DistributionEntity, error) {
	result, err := c.svc.GetDistributeWithPartners(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return c.mapper.DistributionWithPartnersToProto(ctx, result), nil
}

func (c *DistributeController) listProto(ctx context.Context, req *bdspropb.ListDistributionsRequest) (*bdspropb.ListDistributionsResponse, error) {
	if c.sync != nil {
		key := c.sync.GetKey(ctx, _utils.SyncKeyProductDistrs, req.ProductId)
		updated := c.sync.HasUpdated(ctx, key, req.Timestamp)
		c.sync.PutTimeRequest(ctx, key, req.Timestamp)
		if !updated && req.Page == 0 {
			return &bdspropb.ListDistributionsResponse{}, nil
		}
	}
	filter := &dto.FilterDistributeDTO{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
			Sort: req.Sort,
		},
		Text:         req.Query,
		FromDate:     _utils.ParseStringToTime(req.FromDate),
		ToDate:       _utils.ParseStringToTime(req.ToDate),
		CanDeal:      req.CanDeal,
		ChannelPrice: req.ChannelPrice,
		Sort:         req.Sort,
		Status:       c.mapper.ConvertDistributionStatus(req.Status),
	}
	lists, err := c.svc.GetListDistributesWithPartners(ctx, req.ProductId, filter)
	if err != nil {
		return nil, err
	}
	return &bdspropb.ListDistributionsResponse{
		Data:  c.mapper.DistributionWithPartnersListToProto(ctx, lists),
		Total: int64(len(lists)),
	}, nil
}

func (c *DistributeController) revokeProto(ctx context.Context, req *bdspropb.RevokeDistributeRequest) (*bdspropb.RevokeDistributeResponse, error) {
	revoke := dto.RevokeDTO{
		Reason:      uint32(req.Reason),
		ReasonOther: req.ReasonOther,
	}
	if err := c.svc.RevokeDistribute(ctx, req.DistributeId, revoke); err != nil {
		return &bdspropb.RevokeDistributeResponse{Success: false}, err
	}
	return &bdspropb.RevokeDistributeResponse{Success: true}, nil
}

func (c *DistributeController) revokeUserProto(ctx context.Context, req *bdspropb.RevokeUserFromDistributeRequest) (*bdspropb.RevokeUserFromDistributeResponse, error) {
	if err := c.svc.RevokeUserFromDistribute(ctx, req.DistributeId, req.OriginProfileId); err != nil {
		return &bdspropb.RevokeUserFromDistributeResponse{Success: false}, err
	}
	return &bdspropb.RevokeUserFromDistributeResponse{Success: true}, nil
}

func distributeResultToProto(result *dto.DistributeResultDTO) *bdspropb.DistributeResponse {
	if result == nil {
		return &bdspropb.DistributeResponse{}
	}
	return &bdspropb.DistributeResponse{
		Price:             result.Price,
		CanDeal:           result.CanDeal,
		ChannelPrice:      result.ChannelPrice,
		Status:            bdspropb.DistributionStatus(result.Status),
		CommissionType:    bdspropb.CommissionType(result.CommissionType),
		CommissionValue:   result.CommissionValue,
		Note:              result.Note,
		PartnerCommission: result.PartnerCommission,
		CreatedCount:      int32(result.CreatedCount),
		FromDate:          _utils.FormatTimeToString(result.FromDate),
		ToDate:            _utils.FormatTimeToString(result.ToDate),
	}
}
