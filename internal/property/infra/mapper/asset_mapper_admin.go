package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
)

// MapAdminAssetList maps domain.Asset list to AdminAssetItem list
func (m *AssetMapper) MapAdminAssetList(ctx context.Context, assets []domain.Asset) []*bdspropb.AdminAssetItem {
	if assets == nil {
		return nil
	}

	result := make([]*bdspropb.AdminAssetItem, len(assets))
	ownerIds := make(map[uint64]struct{})
	for i, asset := range assets {
		result[i] = m.MapAdminAssetItem(ctx, &asset)
		if asset.OwnerID != nil {
			ownerIds[*asset.OwnerID] = struct{}{}
		}
	}

	ownerMap, err := m.UserClient.GetMapByIDs(ctx, ownerIds)
	if err != nil {
		return nil
	}

	for i, asset := range assets {
		if asset.OwnerID != nil {
			if owner, ok := ownerMap[*asset.OwnerID]; ok {
				result[i].Owner = owner
			}
		}
	}

	return result
}

// MapAdminAssetItem maps domain.Asset to AdminAssetItem
func (m *AssetMapper) MapAdminAssetItem(ctx context.Context, asset *domain.Asset) *bdspropb.AdminAssetItem {
	if asset == nil {
		return nil
	}

	result := &bdspropb.AdminAssetItem{
		Id:               asset.ID,
		Code:             m.CodeDataUsecase.GetAssetCode(ctx, asset.ID), // Generate code based on ID
		Name:             asset.Name,
		Status:           uint32(asset.RentStatus), // Using RentStatus as main status
		StatusName:       m.getAssetStatusName(asset.RentStatus),
		Area:             asset.Area,
		Address:          asset.Address,
		PurchasePrice:    asset.PurchasePrice,
		PurchaseDate:     _utils.FormatTimeToString(asset.PurchaseDate),
		LegalStatus:      uint32(asset.LegalStatus),
		LegalStatusName:  m.getLegalStatusName(asset.LegalStatus),
		CreatedAt:        _utils.FormatTimeToString(asset.CreatedAt),
		UpdatedAt:        _utils.FormatTimeToString(asset.UpdatedAt),
		ProvinceName:     m.CommonMapper.GetProvinceName(asset.Province),
		WardName:         m.CommonMapper.GetWardName(asset.Ward),
		PropertyTypeName: m.CommonMapper.GetPropertyTypeName(asset.PropertyType),
		ProductId:        asset.ProductID,
		ProductCode:      m.CodeDataUsecase.GetProductCodePtr(ctx, asset.ProductID),
		// DistrictName:     m.CommonMapper.GetDistrictName(asset.District),
	}

	// Set contract information (get the latest contract)
	if len(asset.AssetExploitations) > 0 {
		// Get the latest contract by creation date
		latestContract := asset.AssetExploitations[0]
		for _, contract := range asset.AssetExploitations {
			if contract.CreatedAt != nil && latestContract.CreatedAt != nil {
				if contract.CreatedAt.After(*latestContract.CreatedAt) {
					latestContract = contract
				}
			}
		}
		result.ContractEndDate = _utils.FormatTimeToString(latestContract.EndDate)
		result.ContractStatus = uint32(latestContract.Status)
		result.ContractStatusName = m.getContractStatusName(latestContract.Status)
	}

	return result
}

// Helper function to get asset status name
func (m *AssetMapper) getAssetStatusName(status enums.EAssetStatus) string {
	statusNames := map[enums.EAssetStatus]string{
		enums.EAssetStatusOwning:  "Đang sở hữu",
		enums.EAssetStatusRenting: "Đang cho thuê",
		enums.EAssetStatusSold:    "Đã bán",
		enums.EAssetStatusSelling: "Đang bán",
		enums.EAssetStatusNotSold: "Chưa bán",
	}
	if name, exists := statusNames[status]; exists {
		return name
	}
	return "Không xác định"
}

// Helper function to get legal status name
func (m *AssetMapper) getLegalStatusName(status enums.EDocType) string {
	statusNames := map[enums.EDocType]string{
		enums.ERedBook:         "Sổ đỏ",
		enums.EPinkBook:        "Sổ hồng",
		enums.ENotarizedRecord: "Vi bằng",
	}
	if name, exists := statusNames[status]; exists {
		return name
	}
	return "Không xác định"
}

// Helper function to get contract status name
func (m *AssetMapper) getContractStatusName(status uint) string {
	statusNames := map[uint]string{
		10: "Chưa kích hoạt",
		20: "Đang hoạt động",
		30: "Tạm dừng",
		40: "Kết thúc",
		50: "Hủy bỏ",
	}
	if name, exists := statusNames[status]; exists {
		return name
	}
	return "Không xác định"
}
