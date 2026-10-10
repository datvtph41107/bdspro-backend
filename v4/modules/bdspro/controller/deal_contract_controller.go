package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// DealContractController owns transaction/contract/deal routes (v2.1 + v4 aliases).
type DealContractController struct {
	deal     services.DealService
	contract services.ContractDealService
}

func NewDealContractController(deal services.DealService, contract services.ContractDealService) *DealContractController {
	return &DealContractController{deal: deal, contract: contract}
}

func (c *DealContractController) RegisterHTTP(mux *http.ServeMux) {
	registerDealContractRoutes(mux, "/v2.1/transaction/contract/deal", c)
	registerDealContractRoutes(mux, "/v4/transaction/contract/deal", c)
	// Legacy v4 path shape (deal id before resource segment).
	httpapi.GET(mux, "/v4/transaction/contract/deal/{id}/list", c.getContractsLegacy)
	httpapi.GET(mux, "/v4/transaction/contract/deal/{id}/statistic", c.getStatisticsLegacy)
	httpapi.GET(mux, "/v4/transaction/contract/deal/{id}/statistics", c.getStatisticsLegacy)
	httpapi.GET(mux, "/v4/transaction/contract/deal/{id}/process", c.getProcessLegacy)
}

func (c *DealContractController) getProcessLegacy(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	result, _, err := c.deal.GetContractProcess(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}

func (c *DealContractController) getContractsLegacy(ctx context.Context, req dto.DealListByIDRequest) (*httpapi.Result, error) {
	result, err := c.deal.GetContracts(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		return nil, err
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: int64(result.GetTotal())}),
	}
	opts = append(opts, timestampOpt(result.GetTimestamp())...)
	return httpapi.OK(data, opts...), nil
}

func (c *DealContractController) getStatisticsLegacy(ctx context.Context, req dto.DealStatisticsRequest) (*httpapi.Result, error) {
	result, ts, err := c.deal.GetContractStatistics(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result, timestampOpt(ts)...)
}

func (c *DealContractController) processProto(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.ContractProcessResponse, error) {
	result, _, err := c.deal.GetContractProcess(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if result == nil {
		return &bdspropb.ContractProcessResponse{}, nil
	}
	return result, nil
}

func (c *DealContractController) listProto(ctx context.Context, req *bdspropb.GetDealContractRequest) (*bdspropb.DealContractListResponse, error) {
	listReq := dto.DealListByIDRequest{
		DealID:    req.GetDealId(),
		Page:      uint32(req.GetPage()),
		Size:      uint32(req.GetSize()),
		Timestamp: req.GetTimestamp(),
	}
	return c.deal.GetContracts(ctx, listReq)
}

func (c *DealContractController) statisticsProto(ctx context.Context, req *bdspropb.DealCostStatisticsRequestDTO) (*bdspropb.ContractDealStatisticsDTO, error) {
	statReq := dto.DealStatisticsRequest{
		DealID:    req.GetDealId(),
		Payment:   req.GetPayment(),
		Cost:      req.GetCost(),
		Timestamp: req.GetTimestamp(),
	}
	result, _, err := c.deal.GetContractStatistics(ctx, statReq)
	return result, err
}

func (c *DealContractController) detailProto(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.DealContractResponse, error) {
	return c.contract.GetContract(ctx, req.GetId())
}

func (c *DealContractController) deleteProto(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.IdResponse, error) {
	id, err := c.contract.DeleteContract(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &bdspropb.IdResponse{Id: id}, nil
}

func (c *DealContractController) paymentProto(ctx context.Context, req *bdspropb.DealContractPaymentRequest) (*bdspropb.DealContractResponse, error) {
	return c.contract.CreateContractPayment(ctx, req)
}

func (c *DealContractController) paymentActionProto(ctx context.Context, req *bdspropb.PaymentActionRequest) (*bdspropb.ActionResponse, error) {
	return c.contract.CreateActionPayment(ctx, req)
}

func (c *DealContractController) meProto(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.DealContractListResponse, error) {
	return c.contract.GetMyContracts(ctx, uint32(req.GetPage()), uint32(req.GetSize()))
}

func (c *DealContractController) approveProto(context.Context, *bdspropb.ApproveDealContractRequest) (*bdspropb.IdResponse, error) {
	return nil, _errors.BadRequestException("approve contract is not implemented on v4 yet")
}

func (c *DealContractController) summaryProto(context.Context, *sharepb.IdRequest) (*bdspropb.FinancialSummaryResponse, error) {
	return nil, _errors.BadRequestException("contract financial summary is not implemented on v4 yet")
}
