package mapper

import (
	"bdspro/internal/domain"
	pb_bdspro "pb/types/bdspro"
)

func MapProductPb(product *domain.Product) *pb_bdspro.ProductResponse {
	result := &pb_bdspro.ProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Area:        product.Area,
		Address:     product.Address,
		SaleStatus:  uint32(product.SaleStatus),
		RentStatus:  uint32(product.RentStatus),
	}

	if product.District != nil {
		result.DistrictName = product.District.Name
	}
	if product.Province != nil {
		result.ProvinceName = product.Province.Name
	}
	if product.Ward != nil {
		result.WardName = product.Ward.Name
	}

	if product.Price != nil {
		result.Price = &pb_bdspro.ProductPrice{
			// PriceOwner:     product.Price.PriceOwner,
			Id:             product.Price.ID,
			Currency:       product.Price.Currency,
			SalePrice:      product.Price.SalePrice,
			SaleCommission: &product.Price.SaleCommission,
			RentPrice:      product.Price.RentPrice,
			RentCommission: product.Price.RentCommission,
			Deposite:       product.Price.Deposite,
		}

	}

	return result
}

func MapProductPbList(products []domain.Product) []*pb_bdspro.ProductResponse {
	productList := make([]*pb_bdspro.ProductResponse, len(products))
	for i, product := range products {
		productList[i] = MapProductPb(&product)
	}
	return productList
}
