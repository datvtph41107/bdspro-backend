package services

import (
	_dto "common/domain/dto"
	assistantpb "pb/types/assistant"
	sharepb "pb/types/shared"
)

func mapSuggestResponse(product *_dto.ProductV3DTO, modelName string, processingTime int64) *assistantpb.SuggestResponse {
	if product == nil {
		return &assistantpb.SuggestResponse{
			ModelName:      modelName,
			ProcessingTime: processingTime,
		}
	}

	productInfo := &sharepb.ProductV3Proto{
		Name:        product.Name,
		Description: product.Description,
		Area:        product.Area,
	}

	if product.SourceType != nil {
		productInfo.SourceType = &sharepb.ItemV3Proto{
			Id:   product.SourceType.ID,
			Name: product.SourceType.Name,
		}
	}

	if product.HouseInfo != nil {
		productInfo.HouseInfo = &sharepb.HouseInfoV3Proto{
			NumBedroom:  product.HouseInfo.NumBedroom,
			NumBathroom: product.HouseInfo.NumBathroom,
			NumFloor:    product.HouseInfo.NumFloor,
			NumFront:    product.HouseInfo.NumFront,
			NumCarPark:  product.HouseInfo.NumCarPark,
			Orientation: product.HouseInfo.Orientation,
			Furniture:   product.HouseInfo.Furniture,
		}
	}

	if product.PriceData != nil {
		productInfo.PriceData = &sharepb.PriceV3Proto{
			Currency:           product.PriceData.Currency,
			SalePrice:          product.PriceData.SalePrice,
			SaleCommission:     product.PriceData.SaleCommission,
			SaleCommissionType: product.PriceData.SaleCommissionType,
			RentPrice:          product.PriceData.RentPrice,
		}
	}

	if product.TransactionType != nil {
		productInfo.TransactionType = &sharepb.ItemV3Proto{
			Id:   product.TransactionType.ID,
			Name: product.TransactionType.Name,
		}
	}

	if product.Address != nil {
		productInfo.Address = &sharepb.AddressV3Proto{
			Detail:       product.Address.Detail,
			ProvinceId:   product.Address.ProvinceID,
			DistrictId:   product.Address.DistrictID,
			WardId:       product.Address.WardID,
			ProvinceName: product.Address.ProvinceName,
			DistrictName: product.Address.DistrictName,
			WardName:     product.Address.WardName,
		}
	}
	if product.PropertyType != nil {
		productInfo.PropertyType = &sharepb.ItemV3Proto{
			Id:   product.PropertyType.ID,
			Name: product.PropertyType.Name,
		}
	}
	if product.DocType != nil {
		productInfo.DocType = &sharepb.ItemV3Proto{
			Id:   product.DocType.ID,
			Name: product.DocType.Name,
		}
	}
	if len(product.Amenities) > 0 {
		productInfo.Amenities = make([]*sharepb.ItemV3Proto, len(product.Amenities))
		for i, amenity := range product.Amenities {
			productInfo.Amenities[i] = &sharepb.ItemV3Proto{
				Id:   amenity.ID,
				Name: amenity.Name,
			}
		}
	}

	return &assistantpb.SuggestResponse{
		ModelName:      modelName,
		ProcessingTime: processingTime,
		Product:        productInfo,
	}
}
