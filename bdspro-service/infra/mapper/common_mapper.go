package mapper

import (
	"bdspro/internal/domain"
)

type CommonMapper struct {
}

func NewCommonMapper() *CommonMapper {
	return &CommonMapper{}
}

func (m *CommonMapper) GetProvinceName(province *domain.ProvinceV2) string {
	if province == nil {
		return ""
	}
	return province.Name
}

func (m *CommonMapper) GetDistrictName(district *domain.District) string {
	if district == nil {
		return ""
	}
	return district.Name
}

func (m *CommonMapper) GetWardName(ward *domain.WardV2) string {
	if ward == nil {
		return ""
	}
	return ward.Name
}

func (m *CommonMapper) GetSalePrice(saleStatus *domain.ProductPrice) *float64 {
	if saleStatus == nil {
		return nil
	}
	return saleStatus.SalePrice
}

func (m *CommonMapper) GetRentPrice(price *domain.ProductPrice) *float64 {
	if price == nil {
		return nil
	}
	return price.RentPrice
}

func (m *CommonMapper) GetPropertyTypeName(propertyType *domain.PropertyTypeItem) string {
	if propertyType == nil {
		return ""
	}
	return propertyType.Name
}
