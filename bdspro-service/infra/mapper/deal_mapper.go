package mapper

import (
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/utils"
)

type DealMapper interface {
	CreateGroupDealRequestToEntity(request *bdspropb.CreateGroupDealRequest) *domain.Deal
	CreateOrganizationDealRequestToEntity(request *bdspropb.CreateOrganizationDealRequest) *domain.Deal
	CreateGroupDealV2RequestToEntity(request *bdspropb.CreateGroupDealV2Request) *domain.Deal
	UpdateGroupDealRequestToEntity(request *bdspropb.UpdateGroupDealRequest) *domain.Deal
	EntityToCreateGroupDealResponse(deal *domain.Deal) *bdspropb.CreateGroupDealResponse
	EntityToCreateOrganizationDealResponse(deal *domain.Deal) *bdspropb.CreateOrganizationDealResponse
	EntityToCreateGroupDealV2Response(deal *domain.Deal) *bdspropb.CreateGroupDealV2Response
	EntityToUpdateGroupDealResponse(deal *domain.Deal) *bdspropb.UpdateGroupDealResponse
	EntityToGetGroupDealResponse(deal *domain.Deal) *bdspropb.GroupDeal
	EntityToDealDetailResponse(deal *domain.Deal) *bdspropb.GroupDeal
	EntityToGetOrganizationDealsResponse(deals []*domain.Deal, total uint32) *bdspropb.GetOrganizationDealsResponse
	EntityToGetDealsResponse(deals []*domain.Deal, total uint32) *bdspropb.GetDealsResponse
	EntityToDealMemberProto(member *domain.DealMember) *bdspropb.GetDealMemberByDealAndMemberResponse
	EntityToGetGroupDealsV2Response(deals []*domain.Deal, total uint32) *bdspropb.GetGroupDealsV2Response

	// Admin API
	AdminDealSearchRequestToDTO(req *bdspropb.AdminDealSearchRequest) *dto.AdminDealSearchRequest
	EntityToAdminDealItem(deal *domain.Deal) *bdspropb.AdminDealItem
	EntityToAdminDealSearchResponse(deals []*domain.Deal, total int64) *bdspropb.AdminDealSearchResponse

	DealMemberToPb(dealMember *domain.DealMember) *sharepb.DealMember
	DealMembersToPb(dealMembers []*domain.DealMember) []*sharepb.DealMember
}

type groupDealTransformer struct {
	DocumentMapper    *DocumentMapper
	BankAccountMapper *BankAccountMapper
}

func NewGroupDealTransformer(
	documentMapper *DocumentMapper,
	bankAccountMapper *BankAccountMapper,
) DealMapper {
	return &groupDealTransformer{
		DocumentMapper:    documentMapper,
		BankAccountMapper: bankAccountMapper,
	}
}

func (t *groupDealTransformer) CreateGroupDealRequestToEntity(request *bdspropb.CreateGroupDealRequest) *domain.Deal {
	result := &domain.Deal{
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
		OwnerType:      enums.EOwnerOf(request.OwnerType),
		ChargePersonID: &request.ChargePersonId,
		IsUnilateral:   request.IsUnilateral,
		IsManual:       request.IsManual,
		FromDate:       _utils.ParseStringToTime(request.FromDate),
		ToDate:         _utils.ParseStringToTime(request.ToDate),
		Description:    request.Description,

		ShareVisibility: enums.EVisibility(request.ShareVisibility),
	}

	if request.BankAccount != nil {
		result.BankAccount = &domain.BankAccount{
			BankName:      request.BankAccount.BankName,
			AccountNumber: request.BankAccount.AccountNumber,
			AccountName:   request.BankAccount.AccountName,
		}
	}

	if request.Documents != nil {
		result.Documents = make([]domain.AttachDocument, len(request.Documents))
		for i, document := range request.Documents {
			result.Documents[i] = domain.AttachDocument{
				DocName: document.Name,
				DocPath: document.Path,
				DocType: document.Type,
				DocSize: document.Size,
			}
		}
	}

	return result
}

func (t *groupDealTransformer) EntityToDealMemberProto(
	member *domain.DealMember,
) *bdspropb.GetDealMemberByDealAndMemberResponse {

	if member == nil {
		return nil
	}

	pb := &bdspropb.GetDealMemberByDealAndMemberResponse{
		DealId:   member.DealID,
		MemberId: member.MemberID,

		RoleId:          member.RoleID,
		RoleKey:         uint32(member.RoleKey),
		Status:          uint32(member.Status),
		AmountCommit:    member.AmountCommit,
		CommissionValue: member.CommissionValue,
		CommissionType:  uint32(member.CommissionType),
		Note:            member.Note,
		Message:         member.Message,

		InviterId:      member.InviterID,
		DoneInvestment: member.DoneInvestment,
		IsOwner:        member.IsOwner,
	}

	if member.InvitedAt.IsZero() == false {
		pb.InvitedAt = _utils.FormatTimeToString(&member.InvitedAt)
	}

	if member.RespondedAt != nil {
		respondedAt := _utils.FormatTimeToString(member.RespondedAt)
		pb.RespondedAt = respondedAt
	}

	if member.WithdrawnAt != nil {
		withdrawnAt := _utils.FormatTimeToString(member.WithdrawnAt)
		pb.WithdrawnAt = withdrawnAt
	}

	return pb
}

func (t *groupDealTransformer) CreateOrganizationDealRequestToEntity(request *bdspropb.CreateOrganizationDealRequest) *domain.Deal {
	result := &domain.Deal{
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
		result.BankAccount = &domain.BankAccount{
			BankName:      request.BankAccount.BankName,
			AccountNumber: request.BankAccount.AccountNumber,
			AccountName:   request.BankAccount.AccountName,
		}
	}

	if request.Documents != nil {
		result.Documents = make([]domain.AttachDocument, len(request.Documents))
		for i, document := range request.Documents {
			result.Documents[i] = domain.AttachDocument{
				DocName: document.Name,
				DocPath: document.Path,
				DocType: document.Type,
				DocSize: document.Size,
			}
		}
	}

	return result
}

func (t *groupDealTransformer) CreateGroupDealV2RequestToEntity(request *bdspropb.CreateGroupDealV2Request) *domain.Deal {
	result := &domain.Deal{
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

		ShareVisibility: enums.EVisibility(request.ShareVisibility),
	}

	if request.BankAccount != nil {
		result.BankAccount = &domain.BankAccount{
			BankName:      request.BankAccount.BankName,
			AccountNumber: request.BankAccount.AccountNumber,
			AccountName:   request.BankAccount.AccountName,
		}
	}

	if request.Documents != nil {
		result.Documents = make([]domain.AttachDocument, len(request.Documents))
		for i, document := range request.Documents {
			result.Documents[i] = domain.AttachDocument{
				DocName: document.Name,
				DocPath: document.Path,
				DocType: document.Type,
				DocSize: document.Size,
			}
		}
	}

	return result
}

func (t *groupDealTransformer) UpdateGroupDealRequestToEntity(request *bdspropb.UpdateGroupDealRequest) *domain.Deal {
	return &domain.Deal{
		// ID:     uint32(request.Id),
		Name: request.Name,
		// Amount: request.Amount,
		Status: enums.DealStatus(request.Status),
	}
}

func (t *groupDealTransformer) EntityToCreateGroupDealResponse(deal *domain.Deal) *bdspropb.CreateGroupDealResponse {
	return &bdspropb.CreateGroupDealResponse{
		Id: uint32(deal.ID),
	}
}

func (t *groupDealTransformer) EntityToCreateOrganizationDealResponse(deal *domain.Deal) *bdspropb.CreateOrganizationDealResponse {
	return &bdspropb.CreateOrganizationDealResponse{
		Id: uint32(deal.ID),
	}
}

func (t *groupDealTransformer) EntityToCreateGroupDealV2Response(deal *domain.Deal) *bdspropb.CreateGroupDealV2Response {
	return &bdspropb.CreateGroupDealV2Response{
		Id: uint32(deal.ID),
	}
}

func (t *groupDealTransformer) EntityToUpdateGroupDealResponse(deal *domain.Deal) *bdspropb.UpdateGroupDealResponse {
	return &bdspropb.UpdateGroupDealResponse{
		Id: uint32(deal.ID),
	}
}

func (t *groupDealTransformer) EntityToPb(deal *domain.Deal) *bdspropb.GroupDeal {
	result := &bdspropb.GroupDeal{
		Id:   uint32(deal.ID),
		Name: deal.Name,
		// Amount: deal.Amount,
		Status:       uint32(deal.Status),
		OwnerId:      uint32(deal.OwnerId),
		OwnerType:    uint32(deal.OwnerType),
		TargetProfit: deal.TargetProfit,
		DealType:     uint32(deal.DealType),
		DealTypeName: enums.DealTypeNames[deal.DealType],
		Note:         deal.Note,
		NumMember:    deal.NumMember,
		IsUnilateral: deal.IsUnilateral,
		IsManual:     deal.IsManual,
	}

	if deal.Documents != nil {
		result.Documents = make([]*bdspropb.DocumentItem, len(deal.Documents))
		for i, document := range deal.Documents {
			result.Documents[i] = t.DocumentMapper.EntityToPb(&document)
		}
	}

	if deal.BankAccount != nil {
		result.BankAccount = t.BankAccountMapper.EntityToPb(deal.BankAccount)
	}

	return result
}

func (t *groupDealTransformer) EntityToGetGroupDealResponse(deal *domain.Deal) *bdspropb.GroupDeal {
	response := &bdspropb.GroupDeal{
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

func (t *groupDealTransformer) EntityToDealDetailResponse(deal *domain.Deal) *bdspropb.GroupDeal {
	response := &bdspropb.GroupDeal{
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
		Description:             deal.Description,
		ShareVisibility:         uint32(deal.ShareVisibility),
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

func (t *groupDealTransformer) EntityToGetOrganizationDealsResponse(deals []*domain.Deal, total uint32) *bdspropb.GetOrganizationDealsResponse {
	response := &bdspropb.GetOrganizationDealsResponse{
		Total: total,
		Data:  make([]*bdspropb.GroupDeal, len(deals)),
	}

	for i, deal := range deals {
		response.Data[i] = t.EntityToGetGroupDealResponse(deal)
	}

	return response
}

func (t *groupDealTransformer) EntityToGetDealsResponse(deals []*domain.Deal, total uint32) *bdspropb.GetDealsResponse {

	response := &bdspropb.GetDealsResponse{
		Total: total,
		Data:  make([]*bdspropb.GroupDeal, len(deals)),
	}

	for i, deal := range deals {
		response.Data[i] = t.EntityToGetGroupDealResponse(deal)
	}

	return response
}

// func (t *groupDealTransformer) EntityToDealMemberProto(member *domain.DealMember) *sharepb.DealMember {
// 	if member == nil {
// 		return nil
// 	}

// 	proto := &sharepb.DealMember{
// 		DealId:   member.DealID,
// 		MemberId: member.MemberID,
// 		// MemberType:      uint32(member.MemberType),
// 		AmountCommit:    member.AmountCommit,
// 		CommissionValue: member.CommissionValue,
// 		CommissionType:  uint32(member.CommissionType),
// 		Note:            member.Note,
// 		RoleId:          &member.RoleID,
// 		Status:          uint32(member.Status),
// 		Message:         member.Message,
// 		InvitedAt:       _utils.FormatTimeToString(&member.InvitedAt),
// 		InviterId:       member.InviterID,
// 		// colorId:         &member.ColorId,
// 		// IsUnilateral:    member.IsUnilateral,
// 		DoneInvestment: member.DoneInvestment,
// 		IsOwner:        member.IsOwner,
// 	}

// 	// Set responded_at if not nil
// 	if member.RespondedAt != nil {
// 		respondedAt := _utils.FormatTimeToString(member.RespondedAt)
// 		proto.RespondedAt = respondedAt
// 	}

// 	// Set withdrawn_at if not nil
// 	if member.WithdrawnAt != nil {
// 		withdrawnAt := _utils.FormatTimeToString(member.WithdrawnAt)
// 		proto.WithdrawnAt = withdrawnAt
// 	}

// 	return proto
// }

func (t *groupDealTransformer) EntityToGetGroupDealsV2Response(deals []*domain.Deal, total uint32) *bdspropb.GetGroupDealsV2Response {
	response := &bdspropb.GetGroupDealsV2Response{
		Total: total,
		Data:  make([]*bdspropb.GroupDeal, len(deals)),
	}

	for i, deal := range deals {
		response.Data[i] = t.EntityToGetGroupDealResponse(deal)
	}

	return response
}

// AdminDealSearchRequestToDTO converts proto request to DTO
func (t *groupDealTransformer) AdminDealSearchRequestToDTO(req *bdspropb.AdminDealSearchRequest) *dto.AdminDealSearchRequest {
	return &dto.AdminDealSearchRequest{
		Page:      int(req.Page),
		Size:      int(req.Size),
		Keyword:   req.Keyword,
		Status:    req.Status,
		DealType:  req.DealType,
		GroupID:   req.GroupId,
		OwnerID:   req.OwnerId,
		OwnerType: req.OwnerType,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		MinAmount: req.MinAmount,
		MaxAmount: req.MaxAmount,
		// OrganizationID: req.OrganizationId,
	}
}

// EntityToAdminDealItem converts entity to proto admin deal item
func (t *groupDealTransformer) EntityToAdminDealItem(deal *domain.Deal) *bdspropb.AdminDealItem {
	item := &bdspropb.AdminDealItem{
		Id:           deal.ID,
		Name:         deal.Name,
		TargetProfit: deal.TargetProfit,
		TotalAmount:  0, // Will be filled from transaction service
		ProfitAmount: 0, // Will be filled from transaction service
		Status:       uint32(deal.Status),
		StatusName:   enums.DealStatusMap[deal.Status],
		DealType:     uint32(deal.DealType),
		DealTypeName: enums.DealTypeNames[deal.DealType],
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
func (t *groupDealTransformer) EntityToAdminDealSearchResponse(deals []*domain.Deal, total int64) *bdspropb.AdminDealSearchResponse {
	response := &bdspropb.AdminDealSearchResponse{
		Total: total,
		Data:  make([]*bdspropb.AdminDealItem, len(deals)),
	}

	for i, deal := range deals {
		response.Data[i] = t.EntityToAdminDealItem(deal)
	}

	return response
}

func (t *groupDealTransformer) MapProductToDealPb(ctx context.Context, deals []*bdspropb.GroupDeal) {
	// ids := make([]uint64, len(deals))
	// for _, deal := range deals {
	// 	if deal.Products != nil {
	// 		for _, product := range deal.Products {
	// 			ids = append(ids, product.Id)
	// 		}
	// 	}
	// }

	// products, err := c.GetProductAttachmentByIds(ctx, ids)
	// if err != nil {
	// 	return
	// }

	// mapProduct := make(map[uint64]*bdspropb.ProductAttachment)
	// for _, product := range products.Data {
	// 	mapProduct[product.Id] = product
	// }

	// for _, deal := range deals {
	// 	if deal.Products != nil {
	// 		for i, product := range deal.Products {
	// 			deal.Products[i] = mapProduct[product.Id]
	// 		}
	// 	}
	// }
}

func (m *groupDealTransformer) DealMemberToPb(dealMember *domain.DealMember) *sharepb.DealMember {
	return &sharepb.DealMember{
		DealId:   uint64(dealMember.DealID),
		MemberId: uint64(dealMember.MemberID),
		// MemberType: uint32(dealMember.MemberType),
		RoleId:  &dealMember.RoleID,
		RoleKey: uint32(dealMember.RoleKey),
	}
}

func (t *groupDealTransformer) DealMembersToPb(dealMembers []*domain.DealMember) []*sharepb.DealMember {
	pbDealMembers := make([]*sharepb.DealMember, len(dealMembers))
	for i, dealMember := range dealMembers {
		pbDealMembers[i] = t.DealMemberToPb(dealMember)
	}
	return pbDealMembers
}
