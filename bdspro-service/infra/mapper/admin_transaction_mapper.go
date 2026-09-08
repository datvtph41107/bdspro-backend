package mapper

import (
	"bdspro/infra/client"
	"bdspro/internal/domain"
	admin_repo "bdspro/internal/repo/admin"
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
)

type AdminTransactionMapper struct {
	UserClient *client.UserClient
}

func NewAdminTransactionMapper(userClient *client.UserClient) *AdminTransactionMapper {
	return &AdminTransactionMapper{
		UserClient: userClient,
	}
}

// Map request proto to filter
func (m *AdminTransactionMapper) RequestToFilter(req *bdspropb.AdminTransactionSearchRequest) *admin_repo.TransactionFilter {
	filter := &admin_repo.TransactionFilter{
		Page: int(req.Page),
		Size: int(req.Size),
	}

	if req.Keyword != nil {
		filter.Keyword = req.Keyword
	}
	if req.TransactionType != nil {
		filter.TransactionType = req.TransactionType
	}
	if req.TransactionStatus != nil {
		filter.TransactionStatus = req.TransactionStatus
	}
	if req.ApprovalStatus != nil {
		filter.ApprovalStatus = req.ApprovalStatus
	}
	if req.StartDate != nil {
		filter.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		filter.EndDate = req.EndDate
	}
	if req.ProductId != nil {
		filter.ProductId = req.ProductId
	}
	if req.ContactId != nil {
		filter.ContactId = req.ContactId
	}
	if req.OwnerId != nil {
		filter.OwnerId = req.OwnerId
	}
	if req.OwnerType != nil {
		filter.OwnerType = req.OwnerType
	}

	return filter
}

// Map entity to list item proto
func (m *AdminTransactionMapper) EntityToListItem(ctx context.Context, entity *domain.Transaction) *bdspropb.AdminTransactionItem {
	if entity == nil {
		return nil
	}

	item := &bdspropb.AdminTransactionItem{
		Id:                    entity.ID,
		TransactionName:       entity.TransactionName,
		Amount:                entity.Amount,
		Currency:              entity.Currency,
		TransactionType:       uint32(entity.TransactionType),
		TransactionTypeName:   entity.GetTransactionTypeName(),
		TransactionStatus:     uint32(entity.TransactionStatus),
		TransactionStatusName: entity.GetTransactionStatusName(),
		ApprovalStatus:        string(entity.ApprovalStatus),
		TransactionDate:       _utils.FormatTimeToString(&entity.TransactionDate),
		OwnerId:               entity.OwnerId,
		OwnerType:             uint32(entity.OwnerType),
		CreatedAt:             _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:             _utils.FormatTimeToString(entity.UpdatedAt),
	}

	// Deposite amount
	if entity.DepositeAmount != nil {
		item.DepositeAmount = entity.DepositeAmount
	}

	// Product info
	if entity.ProductId != nil {
		item.ProductId = entity.ProductId
		if entity.Product != nil {
			item.ProductName = &entity.Product.Name
			item.ProductCode = &entity.Product.Code
		}
	}

	// Contact info
	if entity.ContactId != nil {
		item.ContactId = entity.ContactId
		// TODO: Load contact info from CRM service if needed
	}

	// Owner info
	if entity.OwnerId > 0 {
		profile := m.UserClient.GetProfileById(ctx, entity.OwnerId)
		if profile != nil {
			item.Owner = profile
		}
	}

	// Approver info
	if entity.ApprovedBy != nil && *entity.ApprovedBy > 0 {
		approvedBy := uint64(*entity.ApprovedBy)
		item.ApprovedBy = &approvedBy

		profile := m.UserClient.GetProfileById(ctx, approvedBy)
		if profile != nil {
			item.Approver = profile
		}

		if entity.ApprovalDate != nil {
			approvalDate := _utils.FormatTimeToString(entity.ApprovalDate)
			item.ApprovalDate = &approvalDate
		}
	}

	return item
}

// Map entities to list response
func (m *AdminTransactionMapper) EntitiesToListResponse(ctx context.Context, entities []*domain.Transaction, total int64) *bdspropb.AdminTransactionSearchResponse {
	items := make([]*bdspropb.AdminTransactionItem, 0, len(entities))

	for _, entity := range entities {
		item := m.EntityToListItem(ctx, entity)
		if item != nil {
			items = append(items, item)
		}
	}

	return &bdspropb.AdminTransactionSearchResponse{
		Data:  items,
		Total: total,
	}
}

// Map entity to detail proto
func (m *AdminTransactionMapper) EntityToDetail(ctx context.Context, entity *domain.Transaction) *bdspropb.AdminTransactionDetail {
	if entity == nil {
		return nil
	}

	detail := &bdspropb.AdminTransactionDetail{
		Id:                    entity.ID,
		TransactionName:       entity.TransactionName,
		Amount:                entity.Amount,
		Currency:              entity.Currency,
		Description:           entity.Description,
		TransactionType:       uint32(entity.TransactionType),
		TransactionTypeName:   entity.GetTransactionTypeName(),
		TransactionStatus:     uint32(entity.TransactionStatus),
		TransactionStatusName: entity.GetTransactionStatusName(),
		ApprovalStatus:        string(entity.ApprovalStatus),
		TransactionDate:       _utils.FormatTimeToString(&entity.TransactionDate),
		OwnerId:               entity.OwnerId,
		OwnerType:             uint32(entity.OwnerType),
		CategoryId:            entity.CategoryId,
		PaymentMethodId:       entity.PaymentMethodId,
		CreatedAt:             _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:             _utils.FormatTimeToString(entity.UpdatedAt),
	}

	// Deposite info
	if entity.DepositeAmount != nil {
		detail.DepositeAmount = entity.DepositeAmount
	}
	if entity.DepositeNote != nil {
		detail.DepositeNote = entity.DepositeNote
	}

	// Related IDs
	if entity.RelatedDealId != nil {
		detail.RelatedDealId = entity.RelatedDealId
	}
	if entity.RelatedTransactionId != nil {
		detail.RelatedTransactionId = entity.RelatedTransactionId
	}

	// Product info
	if entity.ProductId != nil {
		detail.ProductId = entity.ProductId
		if entity.Product != nil {
			detail.Product = &bdspropb.ProductInfo{
				Id:   entity.Product.ID,
				Name: entity.Product.Name,
				Code: entity.Product.Code,
			}

			if entity.Product.Area > 0 {
				area := entity.Product.Area
				detail.Product.Area = area
			}
			// TODO: Add region names if needed
		}
	}

	// Contact info
	if entity.ContactId != nil {
		detail.ContactId = entity.ContactId
		// TODO: Load contact info from CRM service if needed
	}

	// Owner info
	if entity.OwnerId > 0 {
		profile := m.UserClient.GetProfileById(ctx, entity.OwnerId)
		if profile != nil {
			detail.Owner = profile
		}
	}

	// Approver info
	if entity.ApprovedBy != nil && *entity.ApprovedBy > 0 {
		approvedBy := uint64(*entity.ApprovedBy)
		detail.ApprovedBy = &approvedBy

		profile := m.UserClient.GetProfileById(ctx, approvedBy)
		if profile != nil {
			detail.Approver = profile
		}

		if entity.ApprovalDate != nil {
			approvalDate := _utils.FormatTimeToString(entity.ApprovalDate)
			detail.ApprovalDate = &approvalDate
		}
	}

	// Creator and Updater
	if entity.CreatedBy != nil && *entity.CreatedBy > 0 {
		createdBy := *entity.CreatedBy
		detail.CreatedBy = &createdBy

		profile := m.UserClient.GetProfileById(ctx, createdBy)
		if profile != nil {
			detail.Creator = profile
		}
	}

	if entity.UpdatedBy != nil && *entity.UpdatedBy > 0 {
		updatedBy := *entity.UpdatedBy
		detail.UpdatedBy = &updatedBy

		profile := m.UserClient.GetProfileById(ctx, updatedBy)
		if profile != nil {
			detail.Updater = profile
		}
	}

	return detail
}

// Map product info from entity
func (m *AdminTransactionMapper) mapProductInfo(product *domain.Product) *bdspropb.ProductInfo {
	if product == nil {
		return nil
	}

	info := &bdspropb.ProductInfo{
		Id:   product.ID,
		Name: product.Name,
		Code: product.Code,
	}

	if product.Area > 0 {
		info.Area = product.Area
	}

	return info
}

// Map contact info (placeholder - implement when CRM integration is ready)
func (m *AdminTransactionMapper) mapContactInfo(contactId uint64) *bdspropb.ContactInfo {
	// TODO: Implement CRM client to fetch contact info
	return &bdspropb.ContactInfo{
		Id: contactId,
	}
}
