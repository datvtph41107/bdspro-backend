package handler

import (
	"context"
	"organization/infrastructure/client"
	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/dto"
	"organization/internal/usecase"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"
)

type InvestmentHandler struct {
	organizationpb.UnimplementedInvestmentServiceServer
	InvestmentUsecase   *usecase.InvestmentUsecase
	InvestmentValidator *validator.InvestmentValidator
	InvestmentMapper    *transformer.InvestmentMapper
	UserClient          *client.UserClient
}

func NewInvestmentHandler(
	investmentUsecase *usecase.InvestmentUsecase,
	investmentValidator *validator.InvestmentValidator,
	investmentMapper *transformer.InvestmentMapper,
	userClient *client.UserClient,
) *InvestmentHandler {
	return &InvestmentHandler{
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
// @Param investment body organizationpb.SubmitInvestRequest true "Investment"
// @Success 200 {object} uint64
// @Router /investment [post]
func (h *InvestmentHandler) CreateInvestment(ctx context.Context, req *organizationpb.SubmitInvestRequest) (*sharepb.IdDTO, error) {
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
// @Param investment body organizationpb.SubmitInvestRequest true "Investment"
// @Success 200 {object} uint64
// @Router /investment/{id} [put]
func (h *InvestmentHandler) UpdateInvestment(ctx context.Context, req *organizationpb.SubmitInvestRequest) (*sharepb.IdDTO, error) {
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
func (h *InvestmentHandler) DeleteInvestment(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
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
// @Success 200 {object} organizationpb.Investment
// @Router /investment/{id} [get]
func (h *InvestmentHandler) DetailInvestment(ctx context.Context, req *sharepb.IdDTO) (*organizationpb.Investment, error) {
	investment, confirmedInvestments, totalInvestment, err := h.InvestmentUsecase.GetInvestmentWithConfirmedHistory(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	result := h.InvestmentMapper.EntityToPbWithConfirmedHistory(investment, confirmedInvestments, totalInvestment)
	h.UserClient.MapToInvestPb(ctx, []*organizationpb.Investment{result})
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
func (h *InvestmentHandler) GetInvestments(ctx context.Context, req *organizationpb.InvestmentListRequest) (*organizationpb.ListInvestment, error) {
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
	result := make([]*organizationpb.Investment, len(investments))
	for i, investment := range investments {
		result[i] = h.InvestmentMapper.EntityToPb(&investment)
	}
	h.UserClient.MapToInvestPb(ctx, result)

	return &organizationpb.ListInvestment{
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
func (h *InvestmentHandler) ApproveInvestment(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
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
func (h *InvestmentHandler) RejectInvestment(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
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
func (h *InvestmentHandler) ConfirmInvestment(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
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
func (h *InvestmentHandler) UnconfirmInvestment(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
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
// @Success 200 {object} organizationpb.GetInvestmentHistoriesResponse
// @Router /investment/histories/{id} [get]
func (h *InvestmentHandler) GetInvestmentHistories(ctx context.Context, req *organizationpb.GetInvestmentHistoriesRequest) (*organizationpb.GetInvestmentHistoriesResponse, error) {
	histories, total, err := h.InvestmentUsecase.GetInvestmentHistories(ctx, req.DealId, req.MemberId)
	if err != nil {
		return nil, err
	}
	return &organizationpb.GetInvestmentHistoriesResponse{
		Data:            h.InvestmentMapper.EntityToPbWithHistories(histories),
		TotalInvestment: total,
	}, nil
}
