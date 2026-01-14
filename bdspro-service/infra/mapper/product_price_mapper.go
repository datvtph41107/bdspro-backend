package mapper

import (
	"bdspro/internal/dto"
	bdspropb "pb/types/bdspro"
)

type ProductPriceMapper struct{}

func NewProductPriceMapper() *ProductPriceMapper {
	return &ProductPriceMapper{}
}

func (m *ProductPriceMapper) DTOToPb(
	item *dto.PriceHistoryItemDTO,
) *bdspropb.PriceHistoryItem {
	if item == nil {
		return nil
	}

	pb := &bdspropb.PriceHistoryItem{
		Id:              item.ID,
		ProductId:       item.ProductID,
		TransactionType: bdspropb.ProductTransactionType(item.TransactionType),
		Currency:        item.Currency,
		CreatedAt:       item.CreatedAt,
	}

	if item.Old != nil {
		pb.Old = m.mapPriceSnapshot(item.Old)
	}

	if item.New != nil {
		pb.New = m.mapPriceSnapshot(item.New)
	}

	return pb
}

func (m *ProductPriceMapper) mapPriceSnapshot(
	dto *dto.PriceDTO,
) *bdspropb.Price {

	if dto == nil {
		return nil
	}

	snapshot := &bdspropb.Price{}

	if dto.Sale != nil {
		snapshot.Sale = m.mapSalePrice(dto.Sale)
	}

	if dto.Rent != nil {
		snapshot.Rent = m.mapRentPrice(dto.Rent)
	}

	return snapshot
}

func (m *ProductPriceMapper) mapSalePrice(
	sale *dto.SalePriceDTO,
) *bdspropb.SalePrice {

	if sale == nil {
		return nil
	}

	return &bdspropb.SalePrice{
		Price:          value(sale.Price),
		Commission:     value(sale.Commission),
		CommissionType: value(sale.CommissionType),
		Deposit:        value(sale.Deposite),
	}
}

func (m *ProductPriceMapper) mapRentPrice(
	rent *dto.RentPriceDTO,
) *bdspropb.RentPrice {

	if rent == nil {
		return nil
	}

	return &bdspropb.RentPrice{
		Price:          value(rent.Price),
		Commission:     value(rent.Commission),
		CommissionType: value(rent.CommissionType),
		PaymentCycle:   value(rent.PaymentCycle),
	}
}

func value[T any](v *T) T {
	var zero T
	if v == nil {
		return zero
	}
	return *v
}

func (m *ProductPriceMapper) DTOListToPb(
	items []*dto.PriceHistoryItemDTO,
) []*bdspropb.PriceHistoryItem {

	if items == nil {
		return []*bdspropb.PriceHistoryItem{}
	}

	result := make([]*bdspropb.PriceHistoryItem, 0, len(items))
	for _, item := range items {
		result = append(result, m.DTOToPb(item))
	}
	return result
}
