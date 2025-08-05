package handler

import (
	_dto "common/dto"
	_utils "common/utils"
	"context"
	"errors"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"
	"time"

	"organization/infrastructure/client"
	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/dto"
	"organization/internal/enums"
	"organization/internal/usecase"
)

type GroupDealHandler struct {
	organizationpb.UnimplementedGroupDealServiceServer
	DealUsecase          usecase.DealUsecase
	GroupDealTransformer transformer.GroupDealTransformer
	GroupDealValidator   validator.GroupDealValidator
	BdsproClient         *client.BdsproClient
	UserClient           *client.UserClient
}

func NewGroupDealHandler(
	groupDealUsecase usecase.DealUsecase,
	groupDealTransformer transformer.GroupDealTransformer,
	groupDealValidator validator.GroupDealValidator,
	bdsproClient *client.BdsproClient,
	userClient *client.UserClient,
) *GroupDealHandler {
	return &GroupDealHandler{
		DealUsecase:          groupDealUsecase,
		GroupDealTransformer: groupDealTransformer,
		GroupDealValidator:   groupDealValidator,
		BdsproClient:         bdsproClient,
		UserClient:           userClient,
	}
}

// @Summary Tạo thương vụ
// @Description Tạo thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request body organizationpb.CreateGroupDealRequest true "Thông tin thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.CreateGroupDealResponse "Thành công"
// @Router /deal [post]
func (h *GroupDealHandler) CreateGroupDeal(ctx context.Context, req *organizationpb.CreateGroupDealRequest) (*organizationpb.CreateGroupDealResponse, error) {
	if err := h.GroupDealValidator.ValidateCreateGroupDealRequest(req); err != nil {
		return nil, err
	}

	deal := h.GroupDealTransformer.CreateGroupDealRequestToEntity(req)

	deal, err := h.DealUsecase.CreateDeal(ctx, deal)
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToCreateGroupDealResponse(deal), nil
}

// @Summary Cập nhật thương vụ
// @Description Cập nhật thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request body organizationpb.UpdateGroupDealRequest true "Thông tin thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.UpdateGroupDealResponse "Thành công"
// @Router /deal/{id} [put]
func (h *GroupDealHandler) UpdateGroupDeal(ctx context.Context, req *organizationpb.UpdateGroupDealRequest) (*organizationpb.UpdateGroupDealResponse, error) {
	if err := h.GroupDealValidator.ValidateUpdateGroupDealRequest(req); err != nil {
		return nil, err
	}

	deal := h.GroupDealTransformer.UpdateGroupDealRequestToEntity(req)

	deal, err := h.DealUsecase.UpdateGroupDeal(ctx, deal)
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToUpdateGroupDealResponse(deal), nil
}

// @Summary Xóa thương vụ
// @Description Xóa thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.DeleteGroupDealResponse "Thành công"
// @Router /deal/{id} [delete]
func (h *GroupDealHandler) DeleteGroupDeal(ctx context.Context, req *organizationpb.DeleteGroupDealRequest) (*organizationpb.DeleteGroupDealResponse, error) {
	if err := h.GroupDealValidator.ValidateDeleteGroupDealRequest(req); err != nil {
		return nil, err
	}

	err := h.DealUsecase.DeleteGroupDeal(ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}

	return &organizationpb.DeleteGroupDealResponse{Id: req.Id}, nil
}

// @Summary Lấy thông tin thương vụ
// @Description Lấy thông tin thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.GetGroupDealResponse "Thành công"
// @Router /deal/{id} [get]
func (h *GroupDealHandler) GetGroupDeal(ctx context.Context, req *organizationpb.GetGroupDealRequest) (*organizationpb.GroupDeal, error) {
	if err := h.GroupDealValidator.ValidateGetGroupDealRequest(req); err != nil {
		return nil, err
	}

	deal, err := h.DealUsecase.GetGroupDealByID(ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}

	response := h.GroupDealTransformer.EntityToGetGroupDealResponse(deal)
	h.BdsproClient.MapProductToDealPb(ctx, []*organizationpb.GroupDeal{response})
	h.UserClient.MapMemberToDealPb(ctx, []*organizationpb.GroupDeal{response})
	return response, nil
}

// @Summary Lấy danh sách thương vụ
// @Description Lấy danh sách thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request query organizationpb.GetGroupDealsRequest true "Thông tin lấy danh sách thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.GetGroupDealsResponse "Thành công"
// @Router /deal [get]
func (h *GroupDealHandler) GetGroupDeals(ctx context.Context, req *organizationpb.GetGroupDealsRequest) (*organizationpb.GetGroupDealsResponse, error) {
	if err := h.GroupDealValidator.ValidateGetGroupDealsRequest(req); err != nil {
		return nil, err
	}

	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}
	deals, total, err := h.DealUsecase.GetGroupDealsByGroupID(ctx, uint64(req.GroupId), page, size)
	if err != nil {
		return nil, err
	}

	response := h.GroupDealTransformer.EntityToGetGroupDealsResponse(deals, total)
	h.UserClient.MapMemberToDealPb(ctx, response.Data)
	return response, nil
}

// @Summary Cập nhật trạng thái thương vụ
// @Description Cập nhật trạng thái thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Param request body organizationpb.UpdateGroupDealStatusRequest true "Thông tin cập nhật trạng thái thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.UpdateGroupDealStatusResponse "Thành công"
// @Router /deal/{id}/status [put]
func (h *GroupDealHandler) UpdateGroupDealStatus(ctx context.Context, req *organizationpb.UpdateGroupDealStatusRequest) (*organizationpb.UpdateGroupDealStatusResponse, error) {
	if err := h.GroupDealValidator.ValidateUpdateGroupDealStatusRequest(req); err != nil {
		return nil, err
	}

	err := h.DealUsecase.UpdateGroupDealStatus(ctx, uint64(req.Id), enums.DealStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return &organizationpb.UpdateGroupDealStatusResponse{Id: req.Id}, nil
}

// @Summary Lấy thông tin đầu tư thống kê khai báo thương vụ
// @Description Lấy thông tin đầu tư thống kê khai báo thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.InfoInvestmentResponse "Thành công"
// @Router /deal/{id}/investment [get]
func (h *GroupDealHandler) InfoInvestment(ctx context.Context, req *sharepb.IdDTO) (*organizationpb.InfoInvestmentResponse, error) {
	info, err := h.DealUsecase.InfoInvestment(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &organizationpb.InfoInvestmentResponse{
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
// @Success 200 {object} organizationpb.BankAccountItem "Thành công"
// @Router /deal/{id}/bank-account [get]
func (h *GroupDealHandler) AccountInvestment(ctx context.Context, req *sharepb.IdDTO) (*organizationpb.BankAccountItem, error) {
	account, err := h.DealUsecase.AccountInvestment(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &organizationpb.BankAccountItem{
		// Name:          account.AccountName,
		BankName:      account.BankName,
		AccountNumber: account.AccountNumber,
		AccountName:   account.AccountName,
	}, nil
}

func (h *GroupDealHandler) AddMemberToDeal(ctx context.Context, req *organizationpb.AddMemberToDealRequest) (*organizationpb.AddMemberToDealResponse, error) {
	if err := h.GroupDealValidator.ValidateAddMemberToDealRequest(req); err != nil {
		return nil, err
	}

	err := h.DealUsecase.AddMemberToDeal(ctx, req.DealId, req.MemberId)
	if err != nil {
		return nil, err
	}

	return &organizationpb.AddMemberToDealResponse{
		DealId:   req.DealId,
		MemberId: req.MemberId,
	}, nil
}

// @Summary Tạo thương vụ cho tổ chức
// @Description Tạo thương vụ cho tổ chức
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param organizationId path uint64 true "ID tổ chức"
// @Param request body organizationpb.CreateOrganizationDealRequest true "Thông tin thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.CreateOrganizationDealResponse "Thành công"
// @Router /organization/{organizationId}/deal [post]
func (h *GroupDealHandler) CreateOrganizationDeal(ctx context.Context, req *organizationpb.CreateOrganizationDealRequest) (*organizationpb.CreateOrganizationDealResponse, error) {
	// Lấy organizationId từ path parameter
	organizationId := getOrganizationIdFromContext(ctx)

	deal := h.GroupDealTransformer.CreateOrganizationDealRequestToEntity(req)

	createdDeal, err := h.DealUsecase.CreateOrganizationDeal(ctx, organizationId, deal)
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToCreateOrganizationDealResponse(createdDeal), nil
}

// @Summary Tạo thương vụ cho nhóm
// @Description Tạo thương vụ cho nhóm
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param groupId path uint64 true "ID nhóm"
// @Param request body organizationpb.CreateGroupDealV2Request true "Thông tin thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.CreateGroupDealV2Response "Thành công"
// @Router /group/{groupId}/deal [post]
func (h *GroupDealHandler) CreateGroupDealV2(ctx context.Context, req *organizationpb.CreateGroupDealV2Request) (*organizationpb.CreateGroupDealV2Response, error) {
	// Lấy groupId từ path parameter
	groupId := getGroupIdFromContext(ctx)

	deal := h.GroupDealTransformer.CreateGroupDealV2RequestToEntity(req)

	createdDeal, err := h.DealUsecase.CreateGroupDeal(ctx, groupId, deal)
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToCreateGroupDealV2Response(createdDeal), nil
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
// @Success 200 {object} organizationpb.GetOrganizationDealsResponse "Thành công"
// @Router /organization/{organizationId}/deals [get]
func (h *GroupDealHandler) GetOrganizationDeals(ctx context.Context, req *organizationpb.GetOrganizationDealsRequest) (*organizationpb.GetOrganizationDealsResponse, error) {
	organizationId := req.OrganizationId
	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	deals, total, err := h.DealUsecase.GetOrganizationDealsByOrganizationID(ctx, organizationId, page, size)
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToGetOrganizationDealsResponse(deals, total), nil
}

// @Summary Lấy danh sách thương vụ của tổ chức hiện tại mà người dùng request đến
// @Description Lấy danh sách thương vụ của tổ chức hiện tại
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request query sharepb.IdRequest true "Thông tin lấy danh sách thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.GetOrganizationDealsResponse "Thành công"
// @Router /organization/deals/current [get]
func (h *GroupDealHandler) GetOrganizationDealsCurrent(ctx context.Context, req *sharepb.IdRequest) (*organizationpb.GetOrganizationDealsResponse, error) {
	organizationId := getOrganizationIdFromContext(ctx)
	if organizationId == 0 {
		return nil, errors.New("organizationId is required")
	}

	pagable := _dto.Pagable{
		Page: int(req.GetPage()),
		Size: int(req.GetSize()),
	}

	deals, total, err := h.DealUsecase.GetOrganizationDealsByOrganizationID(ctx, organizationId, pagable.GetPage(), pagable.GetSize())
	if err != nil {
		return nil, err
	}

	response := h.GroupDealTransformer.EntityToGetOrganizationDealsResponse(deals, total)
	h.UserClient.MapMemberToDealPb(ctx, response.Data)
	return response, nil
}

// @Summary Lấy danh sách thương vụ của nhóm
// @Description Lấy danh sách thương vụ của nhóm
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param groupId path uint64 true "ID nhóm"
// @Param page query uint32 false "Trang" default(1)
// @Param size query uint32 false "Kích thước" default(10)
// @Security BearerAuth
// @Success 200 {object} organizationpb.GetGroupDealsV2Response "Thành công"
// @Router /group/{groupId}/deals [get]
func (h *GroupDealHandler) GetGroupDealsV2(ctx context.Context, req *organizationpb.GetGroupDealsV2Request) (*organizationpb.GetGroupDealsV2Response, error) {
	groupId := req.GroupId
	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	deals, total, err := h.DealUsecase.GetGroupDealsByGroupID(ctx, groupId, page, size)
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToGetGroupDealsV2Response(deals, total), nil
}

// @Summary Cập nhật trạng thái thương vụ cho tổ chức
// @Description Cập nhật trạng thái thương vụ cho tổ chức
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param organizationId path uint64 true "ID tổ chức"
// @Param dealId path uint64 true "ID thương vụ"
// @Param request body organizationpb.UpdateOrganizationDealStatusRequest true "Thông tin cập nhật"
// @Security BearerAuth
// @Success 200 {object} organizationpb.UpdateOrganizationDealStatusResponse "Thành công"
// @Router /organization/{organizationId}/deal/{dealId}/status [put]
func (h *GroupDealHandler) UpdateOrganizationDealStatus(ctx context.Context, req *organizationpb.UpdateOrganizationDealStatusRequest) (*organizationpb.UpdateOrganizationDealStatusResponse, error) {
	err := h.DealUsecase.UpdateOrganizationDealStatus(ctx, req.OrganizationId, req.DealId, enums.DealStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return &organizationpb.UpdateOrganizationDealStatusResponse{DealId: req.DealId}, nil
}

// @Summary Cập nhật trạng thái thương vụ cho nhóm
// @Description Cập nhật trạng thái thương vụ cho nhóm
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param groupId path uint64 true "ID nhóm"
// @Param dealId path uint64 true "ID thương vụ"
// @Param request body organizationpb.UpdateGroupDealStatusV2Request true "Thông tin cập nhật"
// @Security BearerAuth
// @Success 200 {object} organizationpb.UpdateGroupDealStatusV2Response "Thành công"
// @Router /group/{groupId}/deal/{dealId}/status [put]
func (h *GroupDealHandler) UpdateGroupDealStatusV2(ctx context.Context, req *organizationpb.UpdateGroupDealStatusV2Request) (*organizationpb.UpdateGroupDealStatusV2Response, error) {
	err := h.DealUsecase.UpdateGroupDealStatusV2(ctx, req.GroupId, req.DealId, enums.DealStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return &organizationpb.UpdateGroupDealStatusV2Response{DealId: req.DealId}, nil
}

// @Summary Hủy thương vụ
// @Description Hủy thương vụ với lý do
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path uint32 true "ID thương vụ"
// @Param request body organizationpb.CancelDealRequest true "Lý do hủy"
// @Security BearerAuth
// @Success 200 {object} organizationpb.CancelDealResponse "Thành công"
// @Router /deal/{id}/cancel [put]
func (h *GroupDealHandler) CancelDeal(ctx context.Context, req *organizationpb.CancelDealRequest) (*organizationpb.CancelDealResponse, error) {
	err := h.DealUsecase.CancelDeal(ctx, uint64(req.Id), req.Reason)
	if err != nil {
		return nil, err
	}

	return &organizationpb.CancelDealResponse{Id: req.Id}, nil
}

func (h *GroupDealHandler) UpdateDealSetting(ctx context.Context, req *organizationpb.UpdateDealSettingRequest) (*organizationpb.UpdateDealSettingResponse, error) {
	fromDate, err := time.Parse(time.DateOnly, req.FromDate)
	if err != nil {
		return nil, err
	}
	toDate, err := time.Parse(time.DateOnly, req.ToDate)
	if err != nil {
		return nil, err
	}
	err = h.DealUsecase.UpdateDealSetting(ctx, uint64(req.DealId), &dto.DealSetting{
		FromDate:                fromDate,
		ToDate:                  toDate,
		AllowSharing:            req.AllowSharing,
		MemberCanAddTransaction: req.MemberCanAddTransaction,
		OnlyOwnerGetCommission:  req.OnlyOwnerGetCommission,
		InternalNote:            req.InternalNote,
	})
	if err != nil {
		return nil, err
	}

	return &organizationpb.UpdateDealSettingResponse{DealId: req.DealId}, nil
}

func (h *GroupDealHandler) Members(ctx context.Context, req *sharepb.IdDTO) (*organizationpb.GetDealMembersResponse, error) {
	members, err := h.DealUsecase.GetDealMembers(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	response := make([]*organizationpb.DealMember, len(members))
	for i, member := range members {
		response[i] = &organizationpb.DealMember{
			DealId:          member.DealID,
			MemberId:        member.MemberID,
			MemberType:      uint32(member.MemberType),
			AmountCommit:    member.AmountCommit,
			CommissionValue: member.CommissionValue,
			CommissionType:  uint32(member.CommissionType),
			Note:            member.Note,
		}
	}
	return &organizationpb.GetDealMembersResponse{Data: response}, nil
}

// @Summary Thêm sản phẩm thương vụ
// @Description Thêm sản phẩm thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Param productIds body []uint64 true "Danh sách ID sản phẩm"
// @Security BearerAuth
// @Success 200 {object} organizationpb.SaveDealProductResponse "Thành công"
// @Router /deal/{dealId}/products [post]
func (h *GroupDealHandler) SaveDealProduct(ctx context.Context, req *organizationpb.SaveDealProductRequest) (*organizationpb.SaveDealProductResponse, error) {
	_, err := h.DealUsecase.SaveDealProduct(ctx, req.DealId, req.ProductIds)
	if err != nil {
		return nil, err
	}
	return &organizationpb.SaveDealProductResponse{Status: 1}, nil
}

// @Summary Lấy danh sách sản phẩm thương vụ
// @Description Lấy danh sách sản phẩm thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.GetDealProductsResponse "Thành công"
// @Router /deal/{dealId}/products [get]
func (h *GroupDealHandler) GetDealProducts(ctx context.Context, req *sharepb.IdRequest) (*organizationpb.GetDealProductsResponse, error) {
	pagable := _dto.Pagable{
		Page: int(req.GetPage()),
		Size: int(req.GetSize()),
	}
	products, err := h.DealUsecase.GetDealProducts(ctx, req.Id, pagable)
	if err != nil {
		return nil, err
	}
	productIds := make([]uint64, len(products))
	for i, product := range products {
		productIds[i] = product.ProductID
	}
	productAttachments, err := h.BdsproClient.GetProductAttachmentByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}
	return &organizationpb.GetDealProductsResponse{Data: productAttachments}, nil
}

// @Summary Tổng quan thương vụ
// @Description Tổng quan thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Param query query organizationpb.SummaryRequest true "Thông tin tổng quan"
// @Security BearerAuth
// @Success 200 {object} organizationpb.SummaryResponse "Thành công"
// @Router /deal/{dealId}/summary [get]
func (h *GroupDealHandler) Summary(ctx context.Context, req *organizationpb.SummaryRequest) (*organizationpb.SummaryResponse, error) {
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

	return &organizationpb.SummaryResponse{
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
		AmountProfit:    5000000000,
		AmountTotal:     10000000000,
		NumTransaction:  6,
	}, nil
}

// Helper functions để lấy ID từ context
func getOrganizationIdFromContext(ctx context.Context) uint64 {
	// Implementation sẽ được thêm sau khi có middleware
	return _utils.GetOrganizationIdFromContext(ctx)
}

func getGroupIdFromContext(ctx context.Context) uint64 {
	// Implementation sẽ được thêm sau khi có middleware
	return 0
}
