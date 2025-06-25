package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
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

	if len(product.MediaList) > 0 {
		result.MediaList = make([]*pb_bdspro.MediaItem, len(product.MediaList))
		for i, media := range product.MediaList {
			result.MediaList[i] = &pb_bdspro.MediaItem{
				Id:        media.ID,
				MediaUrl:  media.MediaURL,
				MediaType: media.MediaType,
			}
		}
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

func MapProductAttachmentPbList(products []domain.Product) []*pb_bdspro.ProductAttachment {
	productList := make([]*pb_bdspro.ProductAttachment, len(products))
	for i, product := range products {
		productList[i] = &pb_bdspro.ProductAttachment{
			Id:             product.ID,
			Name:           product.Name,
			Code:           product.Code,
			Area:           product.Area,
			AssetId:        product.AssetId,
			DepositeId:     product.DepositeID,
			RentStatus:     uint32(product.RentStatus),
			SaleStatus:     uint32(product.SaleStatus),
			SourceType:     int32(product.SourceType),
			SourceTypeName: enums.SourceTypeMap[product.SourceType],
		}
		if product.Price != nil {
			productList[i].Price = *product.Price.SalePrice
		}
		if product.MediaList != nil && len(product.MediaList) > 0 {
			productList[i].Image = product.MediaList[0].MediaURL
		}
	}
	return productList
}
