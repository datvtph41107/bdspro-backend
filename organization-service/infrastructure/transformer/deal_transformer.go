package transformer

import (
	_utils "common/utils"
	organizationpb "pb/types/organization"

	"organization/infrastructure/utils"
	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/enums"
)

type DealTransformer interface {
	CreateGroupDealRequestToEntity(request *organizationpb.CreateGroupDealRequest) *entity.Deal
	CreateOrganizationDealRequestToEntity(request *organizationpb.CreateOrganizationDealRequest) *entity.Deal
	CreateGroupDealV2RequestToEntity(request *organizationpb.CreateGroupDealV2Request) *entity.Deal
	UpdateGroupDealRequestToEntity(request *organizationpb.UpdateGroupDealRequest) *entity.Deal
	EntityToCreateGroupDealResponse(deal *entity.Deal) *organizationpb.CreateGroupDealResponse
	EntityToCreateOrganizationDealResponse(deal *entity.Deal) *organizationpb.CreateOrganizationDealResponse
	EntityToCreateGroupDealV2Response(deal *entity.Deal) *organizationpb.CreateGroupDealV2Response
	EntityToUpdateGroupDealResponse(deal *entity.Deal) *organizationpb.UpdateGroupDealResponse
	EntityToGetGroupDealResponse(deal *entity.Deal) *organizationpb.GroupDeal
	EntityToDealDetailResponse(deal *entity.Deal) *organizationpb.GroupDeal
	EntityToGetOrganizationDealsResponse(deals []*entity.Deal, total uint32) *organizationpb.GetOrganizationDealsResponse
	EntityToGetDealsResponse(deals []*entity.Deal, total uint32) *organizationpb.GetDealsResponse
	EntityToDealMemberProto(member *entity.DealMember) *organizationpb.DealMember
	EntityToGetGroupDealsV2Response(deals []*entity.Deal, total uint32) *organizationpb.GetGroupDealsV2Response

	// Admin API
	AdminDealSearchRequestToDTO(req *organizationpb.AdminDealSearchRequest) *dto.AdminDealSearchRequest
	EntityToAdminDealItem(deal *entity.Deal) *organizationpb.AdminDealItem
	EntityToAdminDealSearchResponse(deals []*entity.Deal, total int64) *organizationpb.AdminDealSearchResponse
}

type groupDealTransformer struct {
	DocumentMapper    *DocumentMapper
	BankAccountMapper *BankAccountMapper
}

func NewGroupDealTransformer(
	documentMapper *DocumentMapper,
	bankAccountMapper *BankAccountMapper,
) DealTransformer {
	return &groupDealTransformer{
		DocumentMapper:    documentMapper,
		BankAccountMapper: bankAccountMapper,
	}
}

func (t *groupDealTransformer) CreateGroupDealRequestToEntity(request *organizationpb.CreateGroupDealRequest) *entity.Deal {
	result := &entity.Deal{
		Name: request.Name,
		// Amount:  request.Amount,
		Status:         enums.DealStatus(request.Status),
		TargetProfit:   request.TargetProfit,
		DealType:       enums.DealType(request.DealType),
		Note:           request.Note,
		BankAccountID:  request.BankAccountId,
		CustomerIds:    request.CustomerIds,
		PartnerIds:     request.PartnerIds,
		MemberIds:      request.MemberIds,
		ProductIds:     request.ProductIds,
		OwnerId:        request.OwnerId,
		OwnerType:      enums.OwnerOf(request.OwnerType),
		ChargePersonID: &request.ChargePersonId,
		IsUnilateral:   request.IsUnilateral,
		IsManual:       request.IsManual,
		FromDate:       _utils.ParseStringToTime(request.FromDate),
		ToDate:         _utils.ParseStringToTime(request.ToDate),
	}

	if request.BankAccount != nil {
		result.BankAccount = &entity.BankAccount{
			BankName:      request.BankAccount.BankName,
			AccountNumber: request.BankAccount.AccountNumber,
			AccountName:   request.BankAccount.AccountName,
		}
	}

	if request.Documents != nil {
		result.Documents = make([]entity.AttachDocument, len(request.Documents))
		for i, document := range request.Documents {
			result.Documents[i] = entity.AttachDocument{
				DocName: document.Name,
				DocPath: document.Path,
				DocType: document.Type,
				DocSize: document.Size,
			}
		}
	}

	return result
}

func (t *groupDealTransformer) CreateOrganizationDealRequestToEntity(request *organizationpb.CreateOrganizationDealRequest) *entity.Deal {
	result := &entity.Deal{
		Name:           request.Name,
		TargetProfit:   request.TargetProfit,
		DealType:       enums.DealType(request.DealType),
		Note:           request.Note,
		BankAccountID:  request.BankAccountId,
		CustomerIds:    request.CustomerIds,
		PartnerIds:     request.PartnerIds,
		MemberIds:      request.MemberIds,
		ProductIds:     request.ProductIds,
		ChargePersonID: &request.ChargePersonId,
		IsUnilateral:   request.IsUnilateral,
		IsManual:       request.IsManual,
	}

	if request.BankAccount != nil {
		result.BankAccount = &entity.BankAccount{
			BankName:      request.BankAccount.BankName,
			AccountNumber: request.BankAccount.AccountNumber,
			AccountName:   request.BankAccount.AccountName,
		}
	}

	if request.Documents != nil {
		result.Documents = make([]entity.AttachDocument, len(request.Documents))
		for i, document := range request.Documents {
			result.Documents[i] = entity.AttachDocument{
				DocName: document.Name,
				DocPath: document.Path,
				DocType: document.Type,
				DocSize: document.Size,
			}
		}
	}

	return result
}

func (t *groupDealTransformer) CreateGroupDealV2RequestToEntity(request *organizationpb.CreateGroupDealV2Request) *entity.Deal {
	result := &entity.Deal{
		Name:           request.Name,
		TargetProfit:   request.TargetProfit,
		DealType:       enums.DealType(request.DealType),
		Note:           request.Note,
		BankAccountID:  request.BankAccountId,
		CustomerIds:    request.CustomerIds,
		PartnerIds:     request.PartnerIds,
		MemberIds:      request.MemberIds,
		ProductIds:     request.ProductIds,
		ChargePersonID: &request.ChargePersonId,
		IsUnilateral:   request.IsUnilateral,
		IsManual:       request.IsManual,
	}

	if request.BankAccount != nil {
		result.BankAccount = &entity.BankAccount{
			BankName:      request.BankAccount.BankName,
			AccountNumber: request.BankAccount.AccountNumber,
			AccountName:   request.BankAccount.AccountName,
		}
	}

	if request.Documents != nil {
		result.Documents = make([]entity.AttachDocument, len(request.Documents))
		for i, document := range request.Documents {
			result.Documents[i] = entity.AttachDocument{
				DocName: document.Name,
				DocPath: document.Path,
				DocType: document.Type,
				DocSize: document.Size,
			}
		}
	}

	return result
}

func (t *groupDealTransformer) UpdateGroupDealRequestToEntity(request *organizationpb.UpdateGroupDealRequest) *entity.Deal {
	return &entity.Deal{
		// ID:     uint32(request.Id),
		Name: request.Name,
		// Amount: request.Amount,
		Status: enums.DealStatus(request.Status),
	}
}

func (t *groupDealTransformer) EntityToCreateGroupDealResponse(deal *entity.Deal) *organizationpb.CreateGroupDealResponse {
	return &organizationpb.CreateGroupDealResponse{
		Id: uint32(deal.ID),
	}
}

func (t *groupDealTransformer) EntityToCreateOrganizationDealResponse(deal *entity.Deal) *organizationpb.CreateOrganizationDealResponse {
	return &organizationpb.CreateOrganizationDealResponse{
		Id: uint32(deal.ID),
	}
}

func (t *groupDealTransformer) EntityToCreateGroupDealV2Response(deal *entity.Deal) *organizationpb.CreateGroupDealV2Response {
	return &organizationpb.CreateGroupDealV2Response{
		Id: uint32(deal.ID),
	}
}

func (t *groupDealTransformer) EntityToUpdateGroupDealResponse(deal *entity.Deal) *organizationpb.UpdateGroupDealResponse {
	return &organizationpb.UpdateGroupDealResponse{
		Id: uint32(deal.ID),
	}
}

func (t *groupDealTransformer) EntityToPb(deal *entity.Deal) *organizationpb.GroupDeal {
	result := &organizationpb.GroupDeal{
		Id:   uint32(deal.ID),
		Name: deal.Name,
		// Amount: deal.Amount,
		Status:       uint32(deal.Status),
		OwnerId:      uint32(deal.OwnerId),
		OwnerType:    uint32(deal.OwnerType),
		TargetProfit: deal.TargetProfit,
		DealType:     uint32(deal.DealType),
		DealTypeName: enums.DealTypeMap[deal.DealType],
		Note:         deal.Note,
		NumMember:    deal.NumMember,
		IsUnilateral: deal.IsUnilateral,
		IsManual:     deal.IsManual,
	}

	if deal.Documents != nil {
		result.Documents = make([]*organizationpb.DocumentItem, len(deal.Documents))
		for i, document := range deal.Documents {
			result.Documents[i] = t.DocumentMapper.EntityToPb(&document)
		}
	}

	if deal.BankAccount != nil {
		result.BankAccount = t.BankAccountMapper.EntityToPb(deal.BankAccount)
	}

	return result
}

func (t *groupDealTransformer) EntityToGetGroupDealResponse(deal *entity.Deal) *organizationpb.GroupDeal {
	response := &organizationpb.GroupDeal{
		Id:                      uint32(deal.ID),
		Name:                    deal.Name,
		TargetProfit:            deal.TargetProfit,
		Status:                  uint32(deal.Status),
		StatusName:              enums.DealStatusMap[deal.Status],
		DealType:                uint32(deal.DealType),
		OwnerId:                 uint32(deal.OwnerId),
		OwnerType:               uint32(deal.OwnerType),
		Note:                    deal.Note,
		CancelReason:            deal.CancelReason,
		FromDate:                _utils.FormatTimeToString(deal.FromDate),
		ToDate:                  _utils.FormatTimeToString(deal.ToDate),
		AllowSharing:            deal.AllowSharing,
		MemberCanAddTransaction: deal.MemberCanAddTransaction,
		OnlyOwnerGetCommission:  deal.OnlyOwnerGetCommission,
		InternalNote:            deal.InternalNote,
		Code:                    utils.GenerateDealCode(deal.ID),
	}

	if deal.ChargePersonID != nil {
		response.ChargePersonId = *deal.ChargePersonID
	}

	if deal.BankAccount != nil {
		response.BankAccount = t.BankAccountMapper.EntityToPb(deal.BankAccount)
	}

	// Add other fields as needed
	return response
}

func (t *groupDealTransformer) EntityToDealDetailResponse(deal *entity.Deal) *organizationpb.GroupDeal {
	response := &organizationpb.GroupDeal{
		Id:                      uint32(deal.ID),
		Name:                    deal.Name,
		TargetProfit:            deal.TargetProfit,
		Status:                  uint32(deal.Status),
		StatusName:              enums.DealStatusMap[deal.Status],
		DealType:                uint32(deal.DealType),
		OwnerId:                 uint32(deal.OwnerId),
		OwnerType:               uint32(deal.OwnerType),
		Note:                    deal.Note,
		CancelReason:            deal.CancelReason,
		FromDate:                _utils.FormatTimeToString(deal.FromDate),
		ToDate:                  _utils.FormatTimeToString(deal.ToDate),
		AllowSharing:            deal.AllowSharing,
		MemberCanAddTransaction: deal.MemberCanAddTransaction,
		OnlyOwnerGetCommission:  deal.OnlyOwnerGetCommission,
		InternalNote:            deal.InternalNote,
		Code:                    utils.GenerateDealCode(deal.ID),
		CreatedAt:               _utils.FormatTimeToString(deal.CreatedAt),
	}

	if deal.ChargePersonID != nil {
		response.ChargePersonId = *deal.ChargePersonID
	}

	if deal.BankAccount != nil {
		response.BankAccount = t.BankAccountMapper.EntityToPb(deal.BankAccount)
	}

	// Add other fields as needed
	return response
}

func (t *groupDealTransformer) EntityToGetOrganizationDealsResponse(deals []*entity.Deal, total uint32) *organizationpb.GetOrganizationDealsResponse {
	response := &organizationpb.GetOrganizationDealsResponse{
		Total: total,
		Data:  make([]*organizationpb.GroupDeal, len(deals)),
	}

	for i, deal := range deals {
		response.Data[i] = t.EntityToGetGroupDealResponse(deal)
	}

	return response
}

func (t *groupDealTransformer) EntityToGetDealsResponse(deals []*entity.Deal, total uint32) *organizationpb.GetDealsResponse {

	response := &organizationpb.GetDealsResponse{
		Total: total,
		Data:  make([]*organizationpb.GroupDeal, len(deals)),
	}

	for i, deal := range deals {
		response.Data[i] = t.EntityToGetGroupDealResponse(deal)
	}

	return response
}

func (t *groupDealTransformer) EntityToDealMemberProto(member *entity.DealMember) *organizationpb.DealMember {
	if member == nil {
		return nil
	}

	proto := &organizationpb.DealMember{
		DealId:   member.DealID,
		MemberId: member.MemberID,
		// MemberType:      uint32(member.MemberType),
		AmountCommit:    member.AmountCommit,
		CommissionValue: member.CommissionValue,
		CommissionType:  uint32(member.CommissionType),
		Note:            member.Note,
		RoleId:          &member.RoleID,
		Status:          uint32(member.Status),
		Message:         member.Message,
		InvitedAt:       _utils.FormatTimeToString(&member.InvitedAt),
		InviterId:       member.InviterID,
		// colorId:         &member.ColorId,
		// IsUnilateral:    member.IsUnilateral,
		DoneInvestment: member.DoneInvestment,
		IsOwner:        member.IsOwner,
	}

	// Set responded_at if not nil
	if member.RespondedAt != nil {
		respondedAt := _utils.FormatTimeToString(member.RespondedAt)
		proto.RespondedAt = respondedAt
	}

	// Set withdrawn_at if not nil
	if member.WithdrawnAt != nil {
		withdrawnAt := _utils.FormatTimeToString(member.WithdrawnAt)
		proto.WithdrawnAt = withdrawnAt
	}

	return proto
}

func (t *groupDealTransformer) EntityToGetGroupDealsV2Response(deals []*entity.Deal, total uint32) *organizationpb.GetGroupDealsV2Response {
	response := &organizationpb.GetGroupDealsV2Response{
		Total: total,
		Data:  make([]*organizationpb.GroupDeal, len(deals)),
	}

	for i, deal := range deals {
		response.Data[i] = t.EntityToGetGroupDealResponse(deal)
	}

	return response
}

// AdminDealSearchRequestToDTO converts proto request to DTO
func (t *groupDealTransformer) AdminDealSearchRequestToDTO(req *organizationpb.AdminDealSearchRequest) *dto.AdminDealSearchRequest {
	return &dto.AdminDealSearchRequest{
		Page:           int(req.Page),
		Size:           int(req.Size),
		Keyword:        req.Keyword,
		Status:         req.Status,
		DealType:       req.DealType,
		OrganizationID: req.OrganizationId,
		GroupID:        req.GroupId,
		OwnerID:        req.OwnerId,
		OwnerType:      req.OwnerType,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		MinAmount:      req.MinAmount,
		MaxAmount:      req.MaxAmount,
	}
}

// EntityToAdminDealItem converts entity to proto admin deal item
func (t *groupDealTransformer) EntityToAdminDealItem(deal *entity.Deal) *organizationpb.AdminDealItem {
	item := &organizationpb.AdminDealItem{
		Id:           deal.ID,
		Name:         deal.Name,
		TargetProfit: deal.TargetProfit,
		TotalAmount:  0, // Will be filled from transaction service
		ProfitAmount: 0, // Will be filled from transaction service
		Status:       uint32(deal.Status),
		StatusName:   enums.DealStatusMap[deal.Status],
		DealType:     uint32(deal.DealType),
		DealTypeName: enums.DealTypeMap[deal.DealType],
		Code:         utils.GenerateDealCode(deal.ID),
		OwnerId:      deal.OwnerId,
		OwnerType:    uint32(deal.OwnerType),
		NumMember:    deal.NumMember,
		CreatedAt:    _utils.FormatTimeToString(deal.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(deal.UpdatedAt),
		IsUnilateral: deal.IsUnilateral,
		IsManual:     deal.IsManual,
		Note:         deal.Note,
	}

	if deal.FromDate != nil {
		fromDate := _utils.FormatTimeToString(deal.FromDate)
		item.FromDate = &fromDate
	}

	if deal.ToDate != nil {
		toDate := _utils.FormatTimeToString(deal.ToDate)
		item.ToDate = &toDate
	}

	if deal.CancelReason != nil {
		item.CancelReason = deal.CancelReason
	}

	return item
}

// EntityToAdminDealSearchResponse converts entities to proto response
func (t *groupDealTransformer) EntityToAdminDealSearchResponse(deals []*entity.Deal, total int64) *organizationpb.AdminDealSearchResponse {
	response := &organizationpb.AdminDealSearchResponse{
		Total: total,
		Data:  make([]*organizationpb.AdminDealItem, len(deals)),
	}

	for i, deal := range deals {
		response.Data[i] = t.EntityToAdminDealItem(deal)
	}

	return response
}
