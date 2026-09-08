package handler

import (
	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/infra/validator"
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type DealInvestmentHandler struct {
	bdspropb.UnimplementedInvestmentServiceServer
	InvestmentUsecase   *usecases.InvestmentUsecase
	InvestmentValidator *validator.InvestmentValidator
	InvestmentMapper    *mapper.InvestmentMapper
	UserClient          *client.UserClient
}

func NewDealInvestmentHandler(
	investmentUsecase *usecases.InvestmentUsecase,
	investmentValidator *validator.InvestmentValidator,
	investmentMapper *mapper.InvestmentMapper,
	userClient *client.UserClient,
) *DealInvestmentHandler {
	return &DealInvestmentHandler{
		InvestmentUsecase:   investmentUsecase,
		InvestmentValidator: investmentValidator,
		InvestmentMapper:    investmentMapper,
		UserClient:          userClient,
	}
}

// @Summary submit form khai báo góp vốn
// @Description submit form khai báo góp vốn
// @Tags Investment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param investment body bdspropb.SubmitInvestRequest true "Investment"
// @Success 200 {object} uint64
// @Router /investment [post]
func (h *DealInvestmentHandler) CreateInvestment(ctx context.Context, req *bdspropb.SubmitInvestRequest) (*sharepb.IdDTO, error) {
	investment := h.InvestmentMapper.SubmitInvestRequestToEntity(req)
	if err := h.InvestmentValidator.ValidateCreateInvestment(investment); err != nil {
		return nil, err
	}
	id, err := h.InvestmentUsecase.CreateInvestment(ctx, investment)
	if err != nil {
		return nil, err
	}
	return &sharepb.IdDTO{Id: *id}, nil
}

// @Summary update form khai báo góp vốn
// @Description update form khai báo góp vốn
// @Tags Investment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param investment body bdspropb.SubmitInvestRequest true "Investment"
// @Success 200 {object} uint64
// @Router /investment/{id} [put]
func (h *DealInvestmentHandler) UpdateInvestment(ctx context.Context, req *bdspropb.SubmitInvestRequest) (*sharepb.IdDTO, error) {
	investment := h.InvestmentMapper.SubmitInvestRequestToEntity(req)
	if err := h.InvestmentValidator.ValidateUpdateInvestment(investment); err != nil {
		return nil, err
	}
	id, err := h.InvestmentUsecase.UpdateInvestment(ctx, investment)
	if err != nil {
		return nil, err
	}
	return &sharepb.IdDTO{Id: *id}, nil
}

// @Summary delete form khai báo góp vốn
// @Description delete form khai báo góp vốn
// @Tags Investment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Success 200 {object} uint64
// @Router /investment/{id} [delete]
func (h *DealInvestmentHandler) DeleteInvestment(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
	err := h.InvestmentUsecase.DeleteInvestment(ctx, &req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.IdDTO{Id: req.Id}, nil
}

// @Summary get form khai báo góp vốn
// @Description get form khai báo góp vốn
// @Tags Investment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Success 200 {object} bdspropb.Investment
// @Router /investment/{id} [get]
func (h *DealInvestmentHandler) DetailInvestment(ctx context.Context, req *sharepb.IdDTO) (*bdspropb.Investment, error) {
	investment, confirmedInvestments, totalInvestment, err := h.InvestmentUsecase.GetInvestmentWithConfirmedHistory(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	result := h.InvestmentMapper.EntityToPbWithConfirmedHistory(investment, confirmedInvestments, totalInvestment)
	h.UserClient.MapToInvestPb(ctx, []*bdspropb.Investment{result})
	return result, nil
}

// @Summary get form khai báo góp vốn
// @Description get form khai báo góp vốn
// @Tags Investment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entity.Investment
// @Param dealId path uint64 true "Deal ID"
// @Router /investments/{dealId} [get]
func (h *DealInvestmentHandler) GetInvestments(ctx context.Context, req *bdspropb.InvestmentListRequest) (*bdspropb.ListInvestment, error) {
	dto := &dto.InvestmentDTO{
		DealID: req.DealId,
	}

	// Handle nullable member ID
	if req.MemberId != nil {
		dto.MemberID = *req.MemberId
	}

	investments, total, err := h.InvestmentUsecase.GetInvestments(ctx, dto)
	if err != nil {
		return nil, err
	}
	result := make([]*bdspropb.Investment, len(investments))
	for i, investment := range investments {
		result[i] = h.InvestmentMapper.EntityToPb(&investment)
	}
	h.UserClient.MapToInvestPb(ctx, result)

	return &bdspropb.ListInvestment{
		Data:  result,
		Total: int32(total),
	}, nil
}

// @Summary approve form khai báo góp vốn
// @Description approve form khai báo góp vốn
// @Tags Investment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Success 200 {object} uint64
// @Router /investment/{id}/approve [post]
func (h *DealInvestmentHandler) ApproveInvestment(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
	id, err := h.InvestmentUsecase.ApproveInvestment(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.IdDTO{
		Id: id,
	}, nil
}

// @Summary reject form khai báo góp vốn
// @Description reject form khai báo góp vốn
// @Tags Investment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Success 200 {object} uint64
// @Router /investment/{id}/reject [post]
func (h *DealInvestmentHandler) RejectInvestment(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
	id, err := h.InvestmentUsecase.RejectInvestment(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.IdDTO{
		Id: id,
	}, nil
}

// @Summary confirm form khai báo góp vốn
// @Description confirm form khai báo góp vốn
// @Tags Investment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Success 200 {object} uint64
// @Router /investment/{id}/confirm [post]
func (h *DealInvestmentHandler) ConfirmInvestment(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
	id, err := h.InvestmentUsecase.ChangeConfirmation(ctx, req.Id, true)
	if err != nil {
		return nil, err
	}

	return &sharepb.IdDTO{
		Id: id,
	}, nil
}

// @Summary unconfirm form khai báo góp vốn
// @Description unconfirm form khai báo góp vốn
// @Tags Investment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Success 200 {object} uint64
// @Router /investment/{id}/unconfirm [post]
func (h *DealInvestmentHandler) UnconfirmInvestment(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
	id, err := h.InvestmentUsecase.ChangeConfirmation(ctx, req.Id, false)
	if err != nil {
		return nil, err
	}

	return &sharepb.IdDTO{
		Id: id,
	}, nil
}

// @Summary get investment histories
// @Description get investment histories
// @Tags Investment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Success 200 {object} bdspropb.GetInvestmentHistoriesResponse
// @Router /investment/histories/{id} [get]
func (h *DealInvestmentHandler) GetInvestmentHistories(ctx context.Context, req *bdspropb.GetInvestmentHistoriesRequest) (*bdspropb.GetInvestmentHistoriesResponse, error) {
	histories, total, err := h.InvestmentUsecase.GetInvestmentHistories(ctx, req.DealId, req.MemberId)
	if err != nil {
		return nil, err
	}
	return &bdspropb.GetInvestmentHistoriesResponse{
		Data:            h.InvestmentMapper.EntityToPbWithHistories(histories),
		TotalInvestment: total,
	}, nil
}
