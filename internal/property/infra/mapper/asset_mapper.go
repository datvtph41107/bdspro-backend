package mapper

import (
	"bdspro/infra/client"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_usecase "common/domain/usecase"
	_models "common/models"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type AssetMapper struct {
	*CommonMapper
	AssetLegalMapper *AssetLegalMapper
	ProductMapper    *ProductMapper
	UserClient       *client.UserClient
	CodeDataUsecase  *_usecase.CodeDataUsecase
}

func NewAssetMapper(
	assetLegalMapper *AssetLegalMapper,
	productMapper *ProductMapper,
	userClient *client.UserClient,
	codeDataUsecase *_usecase.CodeDataUsecase) *AssetMapper {
	result := &AssetMapper{
		CommonMapper:     &CommonMapper{},
		AssetLegalMapper: assetLegalMapper,
		ProductMapper:    productMapper,
		UserClient:       userClient,
		CodeDataUsecase:  codeDataUsecase,
	}
	productMapper.PostMapper.AssetMapper = result
	return result
}

func (m *AssetMapper) MapAssetPb(asset *domain.Asset) *bdspropb.Asset {
	if asset == nil {
		return nil
	}

	result := &bdspropb.Asset{
		Id:               asset.ID,
		Name:             asset.Name,
		Area:             asset.Area,
		ProductId:        asset.ProductID,
		ParentAssetId:    asset.ParentAssetID,
		SplitMergeStatus: asset.SplitMergeStatus,
		PurchasePrice:    asset.PurchasePrice,
		PurchaseDate:     _utils.FormatTimeToString(asset.PurchaseDate),
		LegalStatus:      uint32(asset.LegalStatus),
		Archived:         asset.Archived,
		RentStatus:       uint32(asset.RentStatus),
		Address:          asset.Address,
		Description:      asset.Description,
		PropertyTypeId:   asset.PropertyTypeId,
		OwnerType:        &asset.OwnerOf,
		LegalItems:       m.mapLegalItemsPb(asset.LegalItems),
	}

	if asset.Province != nil {
		result.ProvinceName = asset.Province.Name
		result.ProvinceId = &asset.Province.ID
	}
	// if asset.District != nil {
	// 	result.DistrictName = asset.District.Name
	// 	result.DistrictId = &asset.District.ID
	// }
	if asset.Ward != nil {
		result.WardName = asset.Ward.Name
		result.WardId = &asset.Ward.ID
	}

	return result
}

// MapAssetListPb maps AssetList to protobuf Asset (for list responses)
// Không trả legalItems, chỉ trả imageUrl
func (m *AssetMapper) MapAssetListPb(asset *domain.AssetList) *bdspropb.Asset {
	if asset == nil {
		return nil
	}

	result := &bdspropb.Asset{
		Id:               asset.ID,
		Name:             asset.Name,
		Area:             asset.Area,
		ProductId:        asset.ProductID,
		ParentAssetId:    asset.ParentAssetID,
		SplitMergeStatus: asset.SplitMergeStatus,
		PurchasePrice:    asset.PurchasePrice,
		PurchaseDate:     _utils.FormatTimeToString(asset.PurchaseDate),
		LegalStatus:      uint32(asset.LegalStatus),
		Archived:         asset.Archived,
		RentStatus:       uint32(asset.RentStatus),
		Address:          asset.Address,
		Description:      asset.Description,
		PropertyTypeId:   asset.PropertyTypeId,
		OwnerType:        &asset.OwnerOf,
		PropertyTypeName: asset.PropertyTypeName,
		ProvinceName:     asset.ProvinceName,
		WardName:         asset.WardName,
		ImageUrl:         asset.ImageUrl,
		// Không map LegalItems trong danh sách
		LegalItems: nil,
	}

	if asset.ProvinceID != nil {
		result.ProvinceId = asset.ProvinceID
	}
	if asset.WardID != nil {
		result.WardId = asset.WardID
	}

	return result
}

func (m *AssetMapper) MapAssetListPbList(assets []domain.AssetList) []*bdspropb.Asset {
	if assets == nil {
		return nil
	}

	result := make([]*bdspropb.Asset, len(assets))
	for i, asset := range assets {
		result[i] = m.MapAssetListPb(&asset)
	}
	return result
}

func (m *AssetMapper) AssetPbToDomain(asset *bdspropb.SaveAssetRequest) *domain.Asset {
	result := &domain.Asset{
		BaseEntity: _models.BaseEntity{
			ID: asset.Id,
		},
		Name:   asset.Name,
		Area:   asset.Area,
		WardID: asset.WardId,
		// DistrictID:       asset.DistrictId,
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
		OwnerOf:          enums.EOwnerOf(asset.OwnerType),
		// LegalItems:    asset.LegalStatus,
	}
	if asset.LegalItems != nil {
		result.LegalItems = m.AssetLegalMapper.PbToListAssetLegal(asset.LegalItems)
	}
	return result
}

func (m *AssetMapper) MapAssetPbList(assets []domain.Asset) []*bdspropb.Asset {
	if assets == nil {
		return nil
	}

	result := make([]*bdspropb.Asset, len(assets))
	for i, asset := range assets {
		result[i] = m.MapAssetPb(&asset)
	}
	return result
}

func (m *AssetMapper) AssetSavePbToDTO(req *bdspropb.SaveAssetRequest) *dto.AssetSaveRequest {
	if req == nil {
		return nil
	}

	return &dto.AssetSaveRequest{
		ID:               req.Id,
		Name:             req.Name,
		Area:             req.Area,
		WardID:           req.WardId,
		ProvinceID:       req.ProvinceId,
		ProductID:        req.ProductId,
		ParentAssetID:    req.ParentAssetId,
		SplitMergeStatus: req.SplitMergeStatus,
		PurchasePrice:    req.PurchasePrice,
		PurchaseDate:     req.PurchaseDate,
		LegalStatus:      req.LegalStatus,
		Archived:         req.Archived,
		RentStatus:       req.RentStatus,
		Address:          req.Address,
		Description:      req.Description,
		PropertyTypeID:   req.PropertyTypeId,
		LegalItems:       m.mapLegalItemsDto(req.LegalItems),
		OwnerType:        req.OwnerType,
	}
}

func (m *AssetMapper) mapLegalItemsPb(legalItems []domain.AssetLegal) []*bdspropb.LegalItem {
	if legalItems == nil {
		return nil
	}

	result := make([]*bdspropb.LegalItem, len(legalItems))
	for i, item := range legalItems {
		result[i] = &bdspropb.LegalItem{
			Id:           item.ID,
			AssetId:      item.AssetID,
			DocumentType: item.DocumentType,
			DocumentName: item.DocumentName,
			DocumentUrl:  item.DocumentURL,
			IssuedDate:   _utils.FormatTimeToString(item.IssuedDate),
			ExpiryDate:   _utils.FormatTimeToString(item.ExpiryDate),
			Description:  item.Description,
		}
	}
	return result
}

func (m *AssetMapper) mapLegalItemsDto(legalItems []*bdspropb.LegalItem) []dto.LegalItem {
	if legalItems == nil {
		return nil
	}

	result := make([]dto.LegalItem, len(legalItems))
	for i, item := range legalItems {
		result[i] = dto.LegalItem{
			ID:             item.Id,
			AssetID:        item.AssetId,
			DocumentType:   item.DocumentType,
			DocumentNumber: item.DocumentName,
			DocumentURL:    item.DocumentUrl,
			IssueDate:      item.IssuedDate,
			ExpiryDate:     item.ExpiryDate,
			Status:         0, // Default status
		}
	}
	return result
}
