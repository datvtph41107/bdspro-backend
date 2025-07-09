package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_models "common/models"
	_utils "common/utils"
	pb_bdspro "pb/types/bdspro"
)

type AssetMapper struct {
	AssetLegalMapper *AssetLegalMapper
	ProductMapper    *ProductMapper
}

func NewAssetMapper(assetLegalMapper *AssetLegalMapper, productMapper *ProductMapper) *AssetMapper {
	result := &AssetMapper{
		AssetLegalMapper: assetLegalMapper,
		ProductMapper:    productMapper,
	}
	productMapper.PostMapper.AssetMapper = result
	return result
}

func (m *AssetMapper) AssetToPb(asset *domain.Asset) *pb_bdspro.Asset {
	result := &pb_bdspro.Asset{
		Id:             asset.ID,
		Name:           asset.Name,
		Area:           asset.Area,
		Address:        asset.Address,
		WardId:         asset.WardID,
		DistrictId:     asset.DistrictID,
		ProvinceId:     asset.ProvinceID,
		PurchasePrice:  asset.PurchasePrice,
		PurchaseDate:   _utils.FormatTimeToString(asset.PurchaseDate),
		LegalStatus:    uint32(asset.LegalStatus),
		Archived:       asset.Archived,
		RentStatus:     uint32(asset.RentStatus),
		Description:    asset.Description,
		PropertyTypeId: asset.PropertyTypeId,
		ProductId:      asset.ProductID,
		// DocStatus:     uint32(asset.LegalStatus),
		// MediaList:     m.LegalItemsToPb(asset.LegalItems),
	}
	if asset.PropertyType != nil {
		result.PropertyTypeName = asset.PropertyType.Name
	}
	if asset.Ward != nil {
		result.WardName = asset.Ward.Name
	}
	if asset.District != nil {
		result.DistrictName = asset.District.Name
	}
	if asset.Province != nil {
		result.ProvinceName = asset.Province.Name
	}
	if asset.LegalItems != nil {
		result.LegalItems = m.LegalItemsToPb(*asset.LegalItems)
	}
	if asset.Product != nil {
		result.Product = m.ProductMapper.MapProductPb(asset.Product)
	}
	return result
}

func (m *AssetMapper) AssetPbToDomain(asset *pb_bdspro.SaveAssetRequest) *domain.Asset {
	result := &domain.Asset{
		BaseEntity: _models.BaseEntity{
			ID: asset.Id,
		},
		Name:             asset.Name,
		Area:             asset.Area,
		WardID:           asset.WardId,
		DistrictID:       asset.DistrictId,
		ProvinceID:       asset.ProvinceId,
		ProductID:        asset.ProductId,
		ParentAssetID:    asset.ParentAssetId,
		SplitMergeStatus: asset.SplitMergeStatus,
		PurchasePrice:    asset.PurchasePrice,
		PurchaseDate:     _utils.ParseStringToTime(asset.PurchaseDate),
		LegalStatus:      enums.EDocType(asset.LegalStatus),
		Archived:         asset.Archived,
		RentStatus:       enums.EAssetStatus(asset.RentStatus),
		Address:          asset.Address,
		Description:      asset.Description,
		PropertyTypeId:   asset.PropertyTypeId,
		OwnerType:        enums.EOwnerType(asset.OwnerType),
		// LegalItems:    asset.LegalStatus,
	}
	if asset.LegalItems != nil {
		result.LegalItems = m.AssetLegalMapper.PbToListAssetLegal(asset.LegalItems)
	}
	return result
}

func (m *AssetMapper) LegalItemsToPb(legalItems []domain.AssetLegal) []*pb_bdspro.LegalItem {
	pbLegalItems := make([]*pb_bdspro.LegalItem, len(legalItems))
	for i, legalItem := range legalItems {
		pbLegalItems[i] = &pb_bdspro.LegalItem{
			Id:           legalItem.ID,
			AssetId:      legalItem.AssetID,
			DocumentName: legalItem.DocumentName,
			DocumentUrl:  legalItem.DocumentURL,
			DocumentType: legalItem.DocumentType,
			IssuedDate:   legalItem.IssuedDate,
			ExpiryDate:   legalItem.ExpiryDate,
			Description:  legalItem.Description,
		}
	}
	return pbLegalItems
}
