package handler

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"errors"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/infra/validator"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"bdspro/internal/usecases"
)

type DealHandler struct {
	bdspropb.UnimplementedDealServiceServer
	DealUsecase        usecases.DealUsecase
	DealTransformer    mapper.DealMapper
	DealValidator      validator.DealValidator
	OrganizationClient *client.OrganizationClient
	UserClient         *client.UserClient
	PaymentClient      *client.PaymentClient
	ProductRepo        repo.ProductRepo
}

func NewDealHandler(
	dealUsecase usecases.DealUsecase,
	dealTransformer mapper.DealMapper,
	dealValidator validator.DealValidator,
	organizationClient *client.OrganizationClient,
	userClient *client.UserClient,
	paymentClient *client.PaymentClient,
	productRepo repo.ProductRepo,
) *DealHandler {
	return &DealHandler{
		DealUsecase:        dealUsecase,
		DealTransformer:    dealTransformer,
		DealValidator:      dealValidator,
		OrganizationClient: organizationClient,
		UserClient:         userClient,
		PaymentClient:      paymentClient,
		ProductRepo:        productRepo,
	}
}

func (h *DealHandler) GetDealsWithoutProduct(
	ctx context.Context,
	req *sharepb.IdRequest,
) (*bdspropb.GetDealsResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	productID := req.Id
	deals, total, err := h.DealUsecase.GetDealsWithoutProduct(
		ctx,
		userID,
		productID,
		_dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
	)
	if err != nil {
		return nil, err
	}

	return h.DealTransformer.EntityToGetDealsResponse(
		deals,
		uint32(total),
	), nil
}

func (h *DealHandler) AddProductToDeal(ctx context.Context, req *bdspropb.AddProductToDealRequest) (*bdspropb.AddProductToDealResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	err := h.DealUsecase.AddProductToDeal(
		ctx,
		req.DealId,
		userID,
		req.ProductId,
	)
	if err != nil {
		return nil, err
	}

	return &bdspropb.AddProductToDealResponse{
		Success: true,
	}, nil
}

func (h *DealHandler) GetByDealAndMember(
	ctx context.Context,
	req *bdspropb.GetDealMemberByDealAndMemberRequest,
) (*bdspropb.GetDealMemberByDealAndMemberResponse, error) {

	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, errors.New("unauthorized")
	}

	dealMember, err := h.DealUsecase.GetByDealAndMember(
		ctx,
		req.DealId,
		userID,
	)
	if err != nil {
		return nil, err
	}

	if dealMember == nil {
		return &bdspropb.GetDealMemberByDealAndMemberResponse{}, nil
	}

	return h.DealTransformer.EntityToDealMemberProto(dealMember), nil
}

// @Summary Tạo thương vụ
// @Description Tạo thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.CreateGroupDealRequest true "Thông tin thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.CreateGroupDealResponse "Thành công"
// @Router /deal [post]
func (h *DealHandler) CreateGroupDeal(ctx context.Context, req *bdspropb.CreateGroupDealRequest) (*bdspropb.CreateGroupDealResponse, error) {
	if err := h.DealValidator.ValidateCreateGroupDealRequest(req); err != nil {
		return nil, err
	}

	deal := h.DealTransformer.CreateGroupDealRequestToEntity(req)

	deal, err := h.DealUsecase.CreateDeal(ctx, deal)
	if err != nil {
		return nil, err
	}

	return h.DealTransformer.EntityToCreateGroupDealResponse(deal), nil
}

// @Summary Cập nhật thương vụ
// @Description Cập nhật thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.UpdateGroupDealRequest true "Thông tin thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.UpdateGroupDealResponse "Thành công"
// @Router /deal/{id} [put]
func (h *DealHandler) UpdateGroupDeal(ctx context.Context, req *bdspropb.UpdateGroupDealRequest) (*bdspropb.UpdateGroupDealResponse, error) {
	if err := h.DealValidator.ValidateUpdateGroupDealRequest(req); err != nil {
		return nil, err
	}

	deal := h.DealTransformer.UpdateGroupDealRequestToEntity(req)

	deal, err := h.DealUsecase.UpdateGroupDeal(ctx, deal)
	if err != nil {
		return nil, err
	}

	return h.DealTransformer.EntityToUpdateGroupDealResponse(deal), nil
}

// @Summary Xóa thương vụ
// @Description Xóa thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.DeleteGroupDealResponse "Thành công"
// @Router /deal/{id} [delete]
func (h *DealHandler) DeleteGroupDeal(ctx context.Context, req *bdspropb.DeleteGroupDealRequest) (*bdspropb.DeleteGroupDealResponse, error) {
	if err := h.DealValidator.ValidateDeleteGroupDealRequest(req); err != nil {
		return nil, err
	}

	err := h.DealUsecase.DeleteGroupDeal(ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}

	return &bdspropb.DeleteGroupDealResponse{Id: req.Id}, nil
}

// @Summary Lấy thông tin thương vụ
// @Description Lấy thông tin thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetGroupDealResponse "Thành công"
// @Router /deal/{id} [get]
func (h *DealHandler) GetGroupDeal(ctx context.Context, req *bdspropb.GetGroupDealRequest) (*bdspropb.GroupDeal, error) {
	if err := h.DealValidator.ValidateGetGroupDealRequest(req); err != nil {
		return nil, err
	}

	deal, err := h.DealUsecase.GetGroupDealByID(ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}

	response := h.DealTransformer.EntityToDealDetailResponse(deal)
	// h.DealTransformer.MapProductToDealPb(ctx, []*bdspropb.GroupDeal{response})
	h.UserClient.MapMemberToDealPb(ctx, []*bdspropb.GroupDeal{response})
	if response.BankAccount != nil {
		bank, _ := h.PaymentClient.GetBankById(ctx, response.BankAccount.BankId)
		response.BankAccount.Bank = bank
	}

	return response, nil
}

// @Summary Cập nhật trạng thái thương vụ
// @Description Cập nhật trạng thái thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Param request body bdspropb.UpdateGroupDealStatusRequest true "Thông tin cập nhật trạng thái thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.UpdateGroupDealStatusResponse "Thành công"
// @Router /deal/{id}/status [put]
func (h *DealHandler) UpdateGroupDealStatus(ctx context.Context, req *bdspropb.UpdateGroupDealStatusRequest) (*bdspropb.UpdateGroupDealStatusResponse, error) {
	if err := h.DealValidator.ValidateUpdateGroupDealStatusRequest(req); err != nil {
		return nil, err
	}

	err := h.DealUsecase.UpdateGroupDealStatus(ctx, uint64(req.Id), enums.DealStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return &bdspropb.UpdateGroupDealStatusResponse{Id: req.Id}, nil
}

// @Summary Lấy thông tin đầu tư thống kê khai báo thương vụ
// @Description Lấy thông tin đầu tư thống kê khai báo thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.InfoInvestmentResponse "Thành công"
// @Router /deal/{id}/investment [get]
func (h *DealHandler) InfoInvestment(ctx context.Context, req *sharepb.IdDTO) (*bdspropb.InfoInvestmentResponse, error) {
	info, err := h.DealUsecase.InfoInvestment(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &bdspropb.InfoInvestmentResponse{
		AmountPending:   info.AmountPending,
		AmountApproved:  info.AmountApproved,
		AmountRejected:  info.AmountRejected,
		AmountRest:      info.AmountRest,
		NumPending:      info.NumPending,
		NumApproved:     info.NumApproved,
		NumRejected:     info.NumRejected,
		NumInvestment:   info.NumInvestment,
		NumMemberSubmit: info.NumMemberSubmit,
		AmountTarget:    info.AmountTarget,
		PercentDone:     info.PercentDone,
	}, nil
}

// @Summary Lấy thông tin tài khoản đầu tư
// @Description Lấy thông tin tài khoản đầu tư
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.BankAccountItem "Thành công"
// @Router /deal/{id}/bank-account [get]
func (h *DealHandler) AccountInvestment(ctx context.Context, req *sharepb.IdDTO) (*bdspropb.BankAccountItem, error) {
	account, err := h.DealUsecase.AccountInvestment(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	result := &bdspropb.BankAccountItem{
		// Name:          account.AccountName,
		BankName:      account.BankName,
		AccountNumber: account.AccountNumber,
		AccountName:   account.AccountName,
		BankId:        account.BankId,
	}

	bank, _ := h.PaymentClient.GetBankById(ctx, account.BankId)
	if bank != nil {
		result.BankName = bank.Name
	}

	return result, nil
}

func (h *DealHandler) AddMemberToDeal(ctx context.Context, req *bdspropb.AddMemberToDealRequest) (*bdspropb.AddMemberToDealResponse, error) {
	if err := h.DealValidator.ValidateAddMemberToDealRequest(req); err != nil {
		return nil, err
	}

	err := h.DealUsecase.AddMemberToDeal(ctx, req.DealId, req.MemberId)
	if err != nil {
		return nil, err
	}

	return &bdspropb.AddMemberToDealResponse{
		DealId:   req.DealId,
		MemberId: req.MemberId,
	}, nil
}

// @Summary Tạo thương vụ cho tổ chức
// @Description Tạo thương vụ cho tổ chức
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param bdsproId path uint64 true "ID tổ chức"
// @Param request body bdspropb.CreateOrganizationDealRequest true "Thông tin thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.CreateOrganizationDealResponse "Thành công"
// @Router /organization/{organizationId}/deal [post]
func (h *DealHandler) CreateOrganizationDeal(ctx context.Context, req *bdspropb.CreateOrganizationDealRequest) (*bdspropb.CreateOrganizationDealResponse, error) {
	// Lấy organizationId từ path parameter
	organizationId := _utils.GetOrganizationIdFromContext(ctx)

	deal := h.DealTransformer.CreateOrganizationDealRequestToEntity(req)

	createdDeal, err := h.DealUsecase.CreateOrganizationDeal(ctx, organizationId, deal)
	if err != nil {
		return nil, err
	}

	return h.DealTransformer.EntityToCreateOrganizationDealResponse(createdDeal), nil
}

// @Summary Tạo thương vụ cho nhóm
// @Description Tạo thương vụ cho nhóm
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param groupId path uint64 true "ID nhóm"
// @Param request body bdspropb.CreateGroupDealV2Request true "Thông tin thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.CreateGroupDealV2Response "Thành công"
// @Router /group/{groupId}/deal [post]
func (h *DealHandler) CreateGroupDealV2(ctx context.Context, req *bdspropb.CreateGroupDealV2Request) (*bdspropb.CreateGroupDealV2Response, error) {
	// Lấy groupId từ path parameter
	groupId := getGroupIdFromContext(ctx)

	deal := h.DealTransformer.CreateGroupDealV2RequestToEntity(req)

	createdDeal, err := h.DealUsecase.CreateGroupDeal(ctx, groupId, deal)
	if err != nil {
		return nil, err
	}

	return h.DealTransformer.EntityToCreateGroupDealV2Response(createdDeal), nil
}

// @Summary Lấy danh sách thương vụ của tổ chức
// @Description Lấy danh sách thương vụ của tổ chức
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param organizationId path uint64 true "ID tổ chức"
// @Param page query uint32 false "Trang" default(1)
// @Param size query uint32 false "Kích thước" default(10)
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetOrganizationDealsResponse "Thành công"
// @Router /organization/{organizationId}/deals [get]
func (h *DealHandler) GetOrganizationDeals(ctx context.Context, req *bdspropb.GetOrganizationDealsRequest) (*bdspropb.GetOrganizationDealsResponse, error) {
	organizationId := req.OrganizationId
	page := uint32(0)
	size := uint32(10)
	if req.Page != nil {
		page = *req.Page
	}
	if req.Size != nil {
		size = *req.Size
	}

	searchRequest := &dto.DealSearchRequest{
		Pagable: _dto.Pagable{
			Page: page,
			Size: size,
		},
		OrganizationID: &organizationId,
	}

	deals, total, err := h.DealUsecase.GetOrganizationDeals(ctx, searchRequest)
	if err != nil {
		return nil, err
	}

	return h.DealTransformer.EntityToGetOrganizationDealsResponse(deals, total), nil
}

// @Summary Lấy danh sách thương vụ của tổ chức hiện tại mà người dùng request đến
// @Description Lấy danh sách thương vụ của tổ chức hiện tại
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request query sharepb.IdRequest true "Thông tin lấy danh sách thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetOrganizationDealsResponse "Thành công"
// @Router /bdspro/deals/current [get]
func (h *DealHandler) GetOrganizationDealsCurrent(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.GetOrganizationDealsResponse, error) {
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	if organizationId == 0 {
		return nil, errors.New("organizationId is required")
	}

	searchRequest := &dto.DealSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.GetPage(),
			Size: req.GetSize(),
		},
		OrganizationID: &organizationId,
	}

	deals, total, err := h.DealUsecase.GetOrganizationDeals(ctx, searchRequest)
	if err != nil {
		return nil, err
	}

	response := h.DealTransformer.EntityToGetOrganizationDealsResponse(deals, total)
	// h.OrganizationClient.MapMemberToDealPb(ctx, response.Data)
	return response, nil
}

// @Summary Tìm kiếm deals chung
// @Description Tìm kiếm deals với các điều kiện lọc theo organizationId, groupId và text
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.DealSearchRequest true "Thông tin tìm kiếm deals"
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetDealsResponse "Thành công"
// @Router /v2/org/deals/user [get]
func (h *DealHandler) GetDealsUser(ctx context.Context, req *bdspropb.DealSearchRequest) (*bdspropb.GetDealsResponse, error) {
	searchRequest := &dto.DealSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.GetPage(),
			Size: req.GetSize(),
		},
		Text:           req.GetText(),
		OrganizationID: req.OrganizationId,
		GroupID:        req.GroupId,
		// OwnerID:        req.OwnerId,
	}
	// if req.OwnerOf != nil {
	// 	searchRequest.OwnerType = enums.TxOwnerType(*req.OwnerOf)
	// }

	deals, total, err := h.DealUsecase.GetDeals(ctx, searchRequest)
	if err != nil {
		return nil, err
	}

	response := h.DealTransformer.EntityToGetDealsResponse(deals, total)
	h.UserClient.MapMemberToDealPb(ctx, response.Data)
	return response, nil
}

// @Summary Tìm kiếm deals chung
// @Description Tìm kiếm deals với các điều kiện lọc theo organizationId, groupId và text
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.DealSearchRequest true "Thông tin tìm kiếm deals"
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetDealsResponse "Thành công"
// @Router /v2/org/deals/user/bdspro [get]
func (h *DealHandler) GetDealsbdspro(ctx context.Context, req *bdspropb.DealSearchRequest) (*bdspropb.GetDealsResponse, error) {
	searchRequest := &dto.DealSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.GetPage(),
			Size: req.GetSize(),
		},
		Text:           req.GetText(),
		OrganizationID: req.OrganizationId,
		GroupID:        req.GroupId,
	}

	deals, total, err := h.DealUsecase.GetDeals(ctx, searchRequest)
	if err != nil {
		return nil, err
	}

	response := h.DealTransformer.EntityToGetDealsResponse(deals, total)
	// h.OrganizationClient.MapMemberToDealPb(ctx, response.Data)
	return response, nil
}

// @Summary Cập nhật trạng thái thương vụ cho tổ chức
// @Description Cập nhật trạng thái thương vụ cho tổ chức
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param organizationId path uint64 true "ID tổ chức"
// @Param dealId path uint64 true "ID thương vụ"
// @Param request body bdspropb.UpdateOrganizationDealStatusRequest true "Thông tin cập nhật"
// @Security BearerAuth
// @Success 200 {object} bdspropb.UpdateOrganizationDealStatusResponse "Thành công"
// @Router /organization/{organizationId}/deal/{dealId}/status [put]
func (h *DealHandler) UpdateOrganizationDealStatus(ctx context.Context, req *bdspropb.UpdateOrganizationDealStatusRequest) (*bdspropb.UpdateOrganizationDealStatusResponse, error) {
	err := h.DealUsecase.UpdateOrganizationDealStatus(ctx, req.OrganizationId, req.DealId, enums.DealStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return &bdspropb.UpdateOrganizationDealStatusResponse{DealId: req.DealId}, nil
}

// @Summary Cập nhật trạng thái thương vụ cho nhóm
// @Description Cập nhật trạng thái thương vụ cho nhóm
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param groupId path uint64 true "ID nhóm"
// @Param dealId path uint64 true "ID thương vụ"
// @Param request body bdspropb.UpdateGroupDealStatusV2Request true "Thông tin cập nhật"
// @Security BearerAuth
// @Success 200 {object} bdspropb.UpdateGroupDealStatusV2Response "Thành công"
// @Router /group/{groupId}/deal/{dealId}/status [put]
func (h *DealHandler) UpdateGroupDealStatusV2(ctx context.Context, req *bdspropb.UpdateGroupDealStatusV2Request) (*bdspropb.UpdateGroupDealStatusV2Response, error) {
	err := h.DealUsecase.UpdateGroupDealStatusV2(ctx, req.GroupId, req.DealId, enums.DealStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return &bdspropb.UpdateGroupDealStatusV2Response{DealId: req.DealId}, nil
}

// @Summary Hủy thương vụ
// @Description Hủy thương vụ với lý do
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path uint32 true "ID thương vụ"
// @Param request body bdspropb.CancelDealRequest true "Lý do hủy"
// @Security BearerAuth
// @Success 200 {object} bdspropb.CancelDealResponse "Thành công"
// @Router /deal/{id}/cancel [put]
func (h *DealHandler) CancelDeal(ctx context.Context, req *bdspropb.CancelDealRequest) (*bdspropb.CancelDealResponse, error) {
	err := h.DealUsecase.CancelDeal(ctx, uint64(req.Id), req.Reason)
	if err != nil {
		return nil, err
	}

	return &bdspropb.CancelDealResponse{Id: req.Id}, nil
}

// @Summary Cập nhật thông tin thương vụ
// @Description Cập nhật thông tin thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path uint64 true "ID thương vụ"
// @Param request body bdspropb.UpdateDealSettingRequest true "Thông tin cập nhật"
// @Security BearerAuth
// @Success 200 {object} bdspropb.UpdateDealSettingResponse "Thành công"
// @Router /deal/{id}/setting [put]
func (h *DealHandler) UpdateDealSetting(ctx context.Context, req *bdspropb.UpdateDealSettingRequest) (*bdspropb.UpdateDealSettingResponse, error) {
	dto := &dto.DealSetting{
		FromDate:                _utils.ParseStringToTime(req.FromDate),
		ToDate:                  _utils.ParseStringToTime(req.ToDate),
		AllowSharing:            req.AllowSharing,
		MemberCanAddTransaction: req.MemberCanAddTransaction,
		OnlyOwnerGetCommission:  req.OnlyOwnerGetCommission,
		InternalNote:            req.InternalNote,
		OwnerId:                 req.OwnerId,
		Status:                  enums.DealStatus(req.Status),
		AllowManualInput:        req.AllowManualInput,
		AccountNumber:           req.AccountNumber,
		AccountName:             req.AccountName,
		BankId:                  req.BankId,
		TargetProfit:            req.TargetProfit,
		Name:                    req.Name,
		ShareVisibility:         enums.EVisibility(req.ShareVisibility),
		Description:             req.Description,
	}
	err := h.DealUsecase.UpdateDealSetting(ctx, uint64(req.DealId), dto)
	if err != nil {
		return nil, err
	}

	return &bdspropb.UpdateDealSettingResponse{DealId: req.DealId}, nil
}

// @Summary Cập nhật cho phép nhập thủ công
// @Description Cập nhật cho phép nhập thủ công
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path uint64 true "ID thương vụ"
// @Param request body bdspropb.UpdateAllowManualInputRequest true "Thông tin cập nhật"
// @Security BearerAuth
// @Success 200 {object} bdspropb.UpdateAllowManualInputResponse "Thành công"
// @Router /deal/{id}/allow-manual-input [put]
func (h *DealHandler) UpdateAllowManualInput(ctx context.Context, req *bdspropb.UpdateAllowManualInputRequest) (*bdspropb.UpdateAllowManualInputResponse, error) {
	err := h.DealUsecase.UpdateAllowManualInput(ctx, req.DealId, req.AllowManualInput)
	if err != nil {
		return nil, err
	}

	return &bdspropb.UpdateAllowManualInputResponse{DealId: req.DealId}, nil
}

// @Summary Lấy danh sách thành viên thương vụ
// @Description Lấy danh sách thành viên thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path uint64 true "ID thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetDealMembersResponse "Thành công"
// @Router /deal/{id}/members [get]
func (h *DealHandler) Members(ctx context.Context, req *bdspropb.SearchDealMembersRequest) (*bdspropb.GetDealMembersResponse, error) {
	members, err := h.DealUsecase.GetDealMembers(ctx, dto.SearchMembersRequest{
		DealID:         req.DealId,
		Keyword:        req.Keyword,
		Page:           int(req.Page),
		Size:           int(req.Size),
		DoneInvestment: req.DoneInvestment,
	})
	if err != nil {
		return nil, err
	}

	response := make([]*sharepb.DealMember, len(members))
	for i, member := range members {
		response[i] = &sharepb.DealMember{
			DealId:   member.DealID,
			MemberId: member.MemberID,
			// MemberType:      uint32(member.MemberType),
			AmountCommit:    member.AmountCommit,
			CommissionValue: member.CommissionValue,
			CommissionType:  uint32(member.CommissionType),
			Note:            member.Note,
			RoleId:          &member.RoleID,
			IsOwner:         member.IsOwner,
			DoneInvestment:  member.DoneInvestment,
		}
	}
	h.UserClient.MapToDealInvitationPb(ctx, response)
	return &bdspropb.GetDealMembersResponse{Data: response}, nil
}

// @Summary Thêm sản phẩm thương vụ
// @Description Thêm sản phẩm thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Param productIds body []uint64 true "Danh sách ID sản phẩm"
// @Security BearerAuth
// @Success 200 {object} bdspropb.SaveDealProductResponse "Thành công"
// @Router /deal/{dealId}/products [post]
func (h *DealHandler) SaveDealProduct(ctx context.Context, req *bdspropb.SaveDealProductRequest) (*bdspropb.SaveDealProductResponse, error) {
	_, err := h.DealUsecase.SaveDealProduct(ctx, req.DealId, req.ProductIds)
	if err != nil {
		return nil, err
	}
	return &bdspropb.SaveDealProductResponse{Status: 1}, nil
}

// @Summary Lấy danh sách sản phẩm thương vụ
// @Description Lấy danh sách sản phẩm thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetDealProductsResponse "Thành công"
// @Router /deal/{dealId}/products [get]
func (h *DealHandler) GetDealProducts(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.GetDealProductsResponse, error) {
	pagable := _dto.Pagable{
		Page: req.GetPage(),
		Size: req.GetSize(),
	}
	products, err := h.DealUsecase.GetDealProducts(ctx, req.Id, pagable)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return &bdspropb.GetDealProductsResponse{
			Data: []*bdspropb.ProductAttachment{},
		}, nil
	}

	productIds := make([]uint64, len(products))
	for i, product := range products {
		productIds[i] = product.ProductID
	}
	// productAttachments, err := h.BdsproClient.GetProductAttachmentByIds(ctx, productIds)
	productDetails, err := h.ProductRepo.GetDetailByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	productAttachments := make([]*bdspropb.ProductAttachment, len(productDetails))
	for i, product := range productDetails {
		productAttachments[i] = &bdspropb.ProductAttachment{
			Id:             product.ID,
			Name:           product.Name,
			Code:           product.Code,
			Area:           product.Area,
			AssetId:        product.AssetId,
			RentStatus:     uint32(product.RentStatus),
			SaleStatus:     uint32(product.SaleStatus),
			SourceType:     int32(product.SourceType),
			SourceTypeName: enums.SourceTypeMap[product.SourceType],
		}
		// Map price
		if product.Price != nil && product.Price.SalePrice != nil {
			productAttachments[i].Price = *product.Price.SalePrice
		}
		// Map image from MediaList
		if len(product.MediaList) > 0 {
			productAttachments[i].Image = product.MediaList[0].MediaURL
		}
	}
	return &bdspropb.GetDealProductsResponse{
		Data: productAttachments,
	}, nil
}

// @Summary Tổng quan thương vụ
// @Description Tổng quan thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Param query query bdspropb.SummaryRequest true "Thông tin tổng quan"
// @Security BearerAuth
// @Success 200 {object} bdspropb.SummaryResponse "Thành công"
// @Router /deal/{dealId}/summary [get]
func (h *DealHandler) Summary(ctx context.Context, req *bdspropb.SummaryRequest) (*bdspropb.SummaryResponse, error) {
	summary, err := h.DealUsecase.Summary(ctx, &dto.SummaryRequest{
		DealId:          req.DealId,
		AmountPending:   req.AmountPending,
		AmountApproved:  req.AmountApproved,
		AmountRejected:  req.AmountRejected,
		NumPending:      req.NumPending,
		NumApproved:     req.NumApproved,
		NumRejected:     req.NumRejected,
		NumInvestment:   req.NumInvestment,
		NumMemberSubmit: req.NumMemberSubmit,
		AmountTarget:    req.AmountTarget,
		PercentDone:     req.PercentDone,
		AmountCost:      req.AmountCost,
		AmountProfit:    req.AmountProfit,
		AmountTotal:     req.AmountTotal,
		NumTransaction:  req.NumTransaction,
		All:             req.All,
	})
	if err != nil {
		return nil, err
	}

	return &bdspropb.SummaryResponse{
		AmountPending:   summary.AmountPending,
		AmountApproved:  summary.AmountApproved,
		AmountRejected:  summary.AmountRejected,
		AmountRest:      summary.AmountRest,
		NumPending:      summary.NumPending,
		NumApproved:     summary.NumApproved,
		NumRejected:     summary.NumRejected,
		NumInvestment:   summary.NumInvestment,
		NumMemberSubmit: summary.NumMemberSubmit,
		AmountTarget:    summary.AmountTarget,
		PercentDone:     summary.PercentDone,
		AmountCost:      summary.AmountCost,
		AmountProfit:    summary.AmountApproved,
		AmountTotal:     10000000000,
		NumTransaction:  summary.NumTransaction,
		TotalPayment:    summary.TotalPayment,
	}, nil
}

// @Summary Lấy thông tin process của deal
// @Description Lấy thông tin process của deal bao gồm statusName, timestamp, code, transactionId, color
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetDealProcessResponse "Thành công"
// @Router /deal/{dealId}/process [get]
func (h *DealHandler) GetDealProcess(ctx context.Context, req *bdspropb.GetDealProcessRequest) (*bdspropb.GetDealProcessResponse, error) {
	// Lấy profileId từ context
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, errors.New("unauthorized")
	}

	// Lấy bdsproId từ context
	bdsproID := _utils.GetOrganizationIdFromContext(ctx)
	if bdsproID == 0 {
		return nil, errors.New("bdspro not found")
	}

	// Gọi usecase để lấy thông tin process
	processItems, err := h.DealUsecase.GetDealProcess(ctx, req.DealId, bdsproID)
	if err != nil {
		return nil, err
	}

	// Chuyển đổi sang response
	response := &bdspropb.GetDealProcessResponse{
		Data:  make([]*bdspropb.DealProcessItem, len(processItems)),
		Total: uint32(len(processItems)),
	}

	for i, item := range processItems {
		response.Data[i] = &bdspropb.DealProcessItem{
			StatusName:    item.StatusName,
			Timestamp:     item.Timestamp,
			Code:          item.Code,
			TransactionId: item.TransactionId,
			Color:         item.Color,
			BgColor:       item.BgColor,
			BorderColor:   item.BorderColor,
		}
	}

	return response, nil
}

// @Summary Lấy danh sách thương vụ theo chi nhánh
// @Description Lấy danh sách thương vụ theo chi nhánh
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param branchId path uint64 true "ID chi nhánh"
// @Param page query uint32 false "Trang" default(1)
// @Param size query uint32 false "Kích thước" default(10)
// @Security BearerAuth
// @Router /bdspro/branch/{branchId}/deals [get]
func (h *DealHandler) GetBranchDeals(ctx context.Context, req *bdspropb.GetBranchDealsRequest) (*bdspropb.GetGroupDealsV2Response, error) {
	if err := h.DealValidator.ValidateGetBranchDealsRequest(req); err != nil {
		return nil, err
	}

	branchId := req.BranchId
	page := 1
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	deals, total, err := h.DealUsecase.GetBranchDealsByBranchID(ctx, branchId, page, size)
	if err != nil {
		return nil, err
	}

	return h.DealTransformer.EntityToGetGroupDealsV2Response(deals, total), nil
}

// Helper functions để lấy ID từ context
func GetOrganizationIdFromContext(ctx context.Context) uint64 {
	// Implementation sẽ được thêm sau khi có middleware
	return _utils.GetOrganizationIdFromContext(ctx)
}

func getGroupIdFromContext(ctx context.Context) uint64 {
	// Implementation sẽ được thêm sau khi có middleware
	return 0
}

// @Summary Lấy tổng quan thương vụ theo tổ chức
// @Description Lấy thống kê tổng quan thương vụ theo tổ chức
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param organizationId path uint64 true "ID tổ chức"
// @Security BearerAuth
// @Router /organization/{organizationId}/deals/overview [get]
func (h *DealHandler) GetOrganizationDealOverview(ctx context.Context, req *bdspropb.GetOrganizationDealOverviewRequest) (*bdspropb.GetOrganizationDealOverviewResponse, error) {
	organizationID := req.OrganizationId

	// Lấy thống kê từ usecase
	overview, err := h.DealUsecase.GetOrganizationDealOverview(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	total := overview.TotalDealProcessing + overview.TotalDealCompleted + overview.TotalDealCanceled

	// Map sang response
	response := &bdspropb.GetOrganizationDealOverviewResponse{
		OrganizationId:      organizationID,
		TotalDeals:          total,
		TotalCapital:        overview.TotalCapital,
		EstimatedProfit:     overview.EstimatedProfit,
		TotalDealProcessing: overview.TotalDealProcessing,
		TotalDealCompleted:  overview.TotalDealCompleted,
		TotalDealCanceled:   overview.TotalDealCanceled,
	}

	return response, nil
}

// @Summary Lấy tổng quan thương vụ theo nhóm
// @Description Lấy thống kê tổng quan thương vụ theo nhóm
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param groupId path uint64 true "ID nhóm"
// @Security BearerAuth
// @Router /group/{groupId}/deals/overview [get]
func (h *DealHandler) GetGroupDealOverview(ctx context.Context, req *bdspropb.GetGroupDealOverviewRequest) (*bdspropb.GetGroupDealOverviewResponse, error) {
	groupID := req.GroupId

	// Lấy thống kê từ usecase
	overview, err := h.DealUsecase.GetGroupDealOverview(ctx, groupID)
	if err != nil {
		return nil, err
	}

	total := overview.TotalDealProcessing + overview.TotalDealCompleted + overview.TotalDealCanceled

	// Map sang response
	response := &bdspropb.GetGroupDealOverviewResponse{
		GroupId:         groupID,
		TotalDeals:      total,
		TotalCapital:    overview.TotalCapital,
		EstimatedProfit: overview.EstimatedProfit,
		ProcessingDeals: overview.ProcessingDeals,
		CompletedDeals:  overview.TotalDealCompleted,
	}

	return response, nil
}

// @Summary Cập nhật vai trò thành viên thương vụ
// @Description Cập nhật vai trò thành viên thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.UpdateDealMemberRoleRequest true "Thông tin cập nhật"
// @Security BearerAuth
// @Router /deal/member/role [put]
func (h *DealHandler) UpdateDealMemberRole(ctx context.Context, req *bdspropb.UpdateDealMemberRoleRequest) (*bdspropb.UpdateDealMemberRoleResponse, error) {
	_, err := h.DealUsecase.UpdateDealMemberRole(ctx, &dto.UpdateDealMemberRoleRequest{
		DealID:   req.DealId,
		MemberID: req.MemberId,
		RoleKey:  req.RoleKey,
	})
	if err != nil {
		return nil, err
	}
	return &bdspropb.UpdateDealMemberRoleResponse{
		Success: true,
		Message: "Update deal member role successfully",
	}, nil
}

// @Summary Lấy toàn bộ danh sách thương vụ (Admin)
// @Description Lấy toàn bộ danh sách thương vụ trong hệ thống (dành cho admin)
// @Tags Admin - Thương vụ
// @Accept json
// @Produce json
// @Param page query int true "Trang hiện tại"
// @Param size query int true "Số lượng items mỗi trang"
// @Param keyword query string false "Từ khóa tìm kiếm"
// @Param status query int false "Trạng thái deal"
// @Param dealType query int false "Loại deal"
// @Param bdsproId query int false "ID tổ chức"
// @Param groupId query int false "ID nhóm"
// @Param ownerId query int false "ID chủ sở hữu"
// @Param ownerType query int false "Loại chủ sở hữu"
// @Param startDate query string false "Ngày bắt đầu"
// @Param endDate query string false "Ngày kết thúc"
// @Param minAmount query float64 false "Số tiền tối thiểu"
// @Param maxAmount query float64 false "Số tiền tối đa"
// @Security BearerAuth
// @Success 200 {object} bdspropb.AdminDealSearchResponse "Thành công"
// @Router /v2/org/admin/deal/list [get]
func (h *DealHandler) GetAllDeals(ctx context.Context, req *bdspropb.AdminDealSearchRequest) (*bdspropb.AdminDealSearchResponse, error) {
	// Convert request to DTO
	searchDTO := h.DealTransformer.AdminDealSearchRequestToDTO(req)

	// Call usecase
	deals, total, err := h.DealUsecase.GetAllDeals(ctx, searchDTO)
	if err != nil {
		return nil, err
	}

	// Map to response
	return h.DealTransformer.EntityToAdminDealSearchResponse(deals, total), nil
}
