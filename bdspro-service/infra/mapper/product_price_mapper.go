package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type ProductPriceMapper struct{}

func NewProductPriceMapper() *ProductPriceMapper {
	return &ProductPriceMapper{}
}

// @bind: bdspro/infra/mapper.ProductPriceMapper
func (m *ProductPriceMapper) EntityToDTO(entity *domain.ProductPrice) *dto.PriceHistoryResponse {
	if entity == nil {
		return nil
	}

	return &dto.PriceHistoryResponse{
		ID:                 entity.ID,
		Currency:           entity.Currency,
		SalePrice:          entity.SalePrice,
		SaleCommission:     entity.SaleCommission,
		SaleCommissionType: entity.SaleCommissionType,
		Deposite:           entity.Deposite,
		RentPrice:          entity.RentPrice,
		RentCommission:     entity.RentCommission,
		RentCommissionType: entity.RentCommissionType,
		RentPaymentCycle:   entity.RentPaymentCycle,
		CreatedAt:          _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:          _utils.FormatTimeToString(entity.UpdatedAt),
	}
}

func (m *ProductPriceMapper) EntityListToDTO(entities []*domain.ProductPrice) []*dto.PriceHistoryResponse {
	if entities == nil {
		return []*dto.PriceHistoryResponse{}
	}

	result := make([]*dto.PriceHistoryResponse, 0, len(entities))
	for _, entity := range entities {
		result = append(result, m.EntityToDTO(entity))
	}
	return result
}

func (m *ProductPriceMapper) DTOToPb(dto *dto.PriceHistoryResponse) *bdspropb.PriceHistoryItem {
	if dto == nil {
		return nil
	}

	pb := &bdspropb.PriceHistoryItem{
		Id:        dto.ID,
		Currency:  dto.Currency,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}

	if dto.SalePrice != nil {
		pb.SalePrice = *dto.SalePrice
	}
	if dto.SaleCommission != nil {
		pb.SaleCommission = *dto.SaleCommission
	}
	if dto.SaleCommissionType != nil {
		pb.SaleCommissionType = *dto.SaleCommissionType
	}
	if dto.Deposite != nil {
		pb.Deposite = *dto.Deposite
	}
	if dto.RentPrice != nil {
		pb.RentPrice = *dto.RentPrice
	}
	if dto.RentCommission != nil {
		pb.RentCommission = *dto.RentCommission
	}
	if dto.RentCommissionType != nil {
		pb.RentCommissionType = *dto.RentCommissionType
	}
	if dto.RentPaymentCycle != nil {
		pb.RentPaymentCycle = *dto.RentPaymentCycle
	}

	return pb
}

func (m *ProductPriceMapper) DTOListToPb(dtos []*dto.PriceHistoryResponse) []*bdspropb.PriceHistoryItem {
	if dtos == nil {
		return []*bdspropb.PriceHistoryItem{}
	}

	result := make([]*bdspropb.PriceHistoryItem, 0, len(dtos))
	for _, dto := range dtos {
		result = append(result, m.DTOToPb(dto))
	}
	return result
}
