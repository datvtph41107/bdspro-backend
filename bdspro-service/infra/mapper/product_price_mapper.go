package mapper

import (
	"bdspro/internal/dto"
	bdspropb "pb/types/bdspro"
)

type ProductPriceMapper struct{}

func NewProductPriceMapper() *ProductPriceMapper {
	return &ProductPriceMapper{}
}

func (m *ProductPriceMapper) DTOToPbDetail(
	dto *dto.ProductPriceDetailDTO,
) *bdspropb.ProductPriceDetail {
	if dto == nil {
		return nil
	}

	return &bdspropb.ProductPriceDetail{
		Price:        dto.Price,
		PricePerM2:   dto.PricePerM2,
		Currency:     dto.Currency,
		Cost:         dto.Cost,
		Commission:   dto.Commission,
		Channel:      dto.Channel,
		ChannelValue: dto.ChannelValue,
		Profit:       dto.Profit,
		Note:         dto.Note,
	}
}

func (m *ProductPriceMapper) DTOToPb(item *dto.PriceHistoryItemDTO) *bdspropb.PriceHistoryItem {
	if item == nil {
		return nil
	}

	return &bdspropb.PriceHistoryItem{
		Id:              item.ID,
		ProductId:       item.ProductID,
		TransactionType: bdspropb.ProductTransactionType(item.TransactionType),
		Currency:        item.Currency,
		Old:             m.mapPrice(item.Old),
		New:             m.mapPrice(item.New),
		IsCurrentPrice:  item.IsCurrent,
		CreatedAt:       item.CreatedAt,
		ChangeNote:      item.ChangeNote,
		DistributeId:    item.DistributeID,
	}
}

func (m *ProductPriceMapper) mapPrice(price *dto.PriceDTO) *bdspropb.Price {
	if price == nil {
		return nil
	}

	return &bdspropb.Price{
		Price:          value(price.Price),
		Commission:     value(price.Commission),
		CommissionType: value(price.CommissionType),
		CreatedBy:      m.mapUser(price.CreatedBy),
	}
}

func (m *ProductPriceMapper) mapUser(
	user *dto.UserInfoDTO,
) *bdspropb.UserInfo {

	if user == nil {
		return nil
	}

	return &bdspropb.UserInfo{
		Id:     user.ID,
		Name:   user.Name,
		Avatar: user.Avatar,
		Role:   user.Role,
	}
}

func value[T any](v *T) T {
	var zero T
	if v == nil {
		return zero
	}
	return *v
}

func (m *ProductPriceMapper) DTOListToPb(items []*dto.PriceHistoryItemDTO) []*bdspropb.PriceHistoryItem {
	if items == nil {
		return []*bdspropb.PriceHistoryItem{}
	}

	result := make([]*bdspropb.PriceHistoryItem, 0, len(items))
	for _, item := range items {
		result = append(result, m.DTOToPb(item))
	}
	return result
}
