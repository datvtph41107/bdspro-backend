package transformer

import (
	"organization/internal/domain/entity"
	"organization/internal/enums"
	bdspropb "pb/types/bdspro"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"
)

type GroupDealTransformer interface {
	CreateGroupDealRequestToEntity(request *organizationpb.CreateGroupDealRequest) *entity.Deal
	UpdateGroupDealRequestToEntity(request *organizationpb.UpdateGroupDealRequest) *entity.Deal
	EntityToCreateGroupDealResponse(deal *entity.Deal) *organizationpb.CreateGroupDealResponse
	EntityToUpdateGroupDealResponse(deal *entity.Deal) *organizationpb.UpdateGroupDealResponse
	EntityToGetGroupDealResponse(deal *entity.Deal) *organizationpb.GroupDeal
	EntityToGetGroupDealsResponse(deals []*entity.Deal, total uint32) *organizationpb.GetGroupDealsResponse
}

type groupDealTransformer struct {
	DocumentMapper    *DocumentMapper
	BankAccountMapper *BankAccountMapper
}

func NewGroupDealTransformer(
	documentMapper *DocumentMapper,
	bankAccountMapper *BankAccountMapper,
) GroupDealTransformer {
	return &groupDealTransformer{
		DocumentMapper:    documentMapper,
		BankAccountMapper: bankAccountMapper,
	}
}

func (t *groupDealTransformer) CreateGroupDealRequestToEntity(request *organizationpb.CreateGroupDealRequest) *entity.Deal {
	result := &entity.Deal{
		Name: request.Name,
		// Amount:  request.Amount,
		Status:        enums.DealStatus(request.Status),
		TargetProfit:  request.TargetProfit,
		DealType:      enums.DealType(request.DealType),
		Note:          request.Note,
		BankAccountID: request.BankAccountId,
		CustomerIds:   request.CustomerIds,
		PartnerIds:    request.PartnerIds,
		MemberIds:     request.MemberIds,
		ProductIds:    request.ProductIds,
		OwnerId:       request.OwnerId,
		// OwnerType:     enums.OwnerTypeGroup,
		ChargePersonID: &request.ChargePersonId,
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
	result := t.EntityToPb(deal)

	// if deal.Products != nil {
	// 	result.Products = make([]*organizationpb.ProductItem, len(deal.Products))
	// 	for i, product := range deal.Products {
	// 		result.Products[i] = &organizationpb.ProductItem{
	// 			Id: uint32(product.ProductID),
	// 		}
	// 	}
	// }

	if deal.BankAccount != nil {
		result.BankAccount = &organizationpb.BankAccountItem{
			Name:          deal.BankAccount.BankName,
			BankName:      deal.BankAccount.BankName,
			AccountNumber: deal.BankAccount.AccountNumber,
			AccountName:   deal.BankAccount.AccountName,
		}
	}

	if deal.Documents != nil {
		result.Documents = make([]*organizationpb.DocumentItem, len(deal.Documents))
		for i, document := range deal.Documents {
			result.Documents[i] = &organizationpb.DocumentItem{
				Name: document.DocName,
				Path: document.DocPath,
				Type: document.DocType,
				Size: document.DocSize,
			}
		}
	}

	if deal.Products != nil {
		result.Products = make([]*bdspropb.ProductAttachment, len(deal.Products))
		for i, product := range deal.Products {
			result.Products[i] = &bdspropb.ProductAttachment{
				Id: uint64(product.ProductID),
			}
		}
	}

	if deal.Members != nil {
		result.Members = make([]*sharepb.ProfileItem, len(deal.Members))
		for i, member := range deal.Members {
			result.Members[i] = &sharepb.ProfileItem{
				Id: uint64(member.MemberID),
			}
		}
	}

	if deal.Customers != nil {
		result.Customers = make([]*sharepb.ProfileItem, len(deal.Customers))
		for i, customer := range deal.Customers {
			result.Customers[i] = &sharepb.ProfileItem{
				Id: uint64(customer.MemberID),
			}
		}
	}

	if deal.Partners != nil {
		result.Partners = make([]*sharepb.ProfileItem, len(deal.Partners))
		for i, partner := range deal.Partners {
			result.Partners[i] = &sharepb.ProfileItem{
				Id: uint64(partner.MemberID),
			}
		}
	}

	return result
}

func (t *groupDealTransformer) EntityToGetGroupDealsResponse(deals []*entity.Deal, total uint32) *organizationpb.GetGroupDealsResponse {
	dealsResponse := make([]*organizationpb.GroupDeal, len(deals))
	for i, deal := range deals {
		dealsResponse[i] = &organizationpb.GroupDeal{
			Id:   uint32(deal.ID),
			Name: deal.Name,
			// Amount:       deal.Amount,
			DealType:     uint32(deal.DealType),
			DealTypeName: enums.DealTypeMap[deal.DealType],
			// CreatedAt:    _utils.FormatTimeToString(&deal.CreatedAt),
			// UpdatedAt:    _utils.FormatTimeToString(&deal.UpdatedAt),
		}
	}
	return &organizationpb.GetGroupDealsResponse{
		Data:  dealsResponse,
		Total: total,
	}
}
