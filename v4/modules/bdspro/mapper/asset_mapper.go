package mapper

import (
	"time"

	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/enums"
)

// CreatePropertyAssetToEntity maps CreatePropertyRequest.asset (+ legalInfo) to an Asset.
// Name / type / location / area are taken from the property bundle because AssetCreateWithProduct does not carry them.
func CreatePropertyAssetToEntity(req *bdspropb.CreatePropertyRequest, bundle *entity.PropertyCreateBundle, ownerProfileID uint64) *entity.Asset {
	a := req.GetAsset()
	if a == nil {
		return nil
	}
	asset := &entity.Asset{
		PurchasePrice: a.GetPurchasePrice(),
		LegalStatus:   uint32(enums.ERedBook),
		RentStatus:    uint32(enums.EAssetStatusOwning),
		Description:   a.GetDescription(),
		OwnerOf:       uint32(enums.EOwnerOfMember),
	}
	if a.GetLegalStatus() > 0 {
		asset.LegalStatus = uint32(a.GetLegalStatus())
	}
	if a.GetRentStatus() > 0 {
		asset.RentStatus = uint32(a.GetRentStatus())
	}
	if a.GetPurchaseDate() > 0 {
		t := time.Unix(int64(a.GetPurchaseDate()), 0)
		asset.PurchaseDate = &t
	}
	if ownerProfileID > 0 {
		owner := ownerProfileID
		asset.OwnerID = &owner
	}

	if bundle != nil {
		if info := bundle.Info; info != nil {
			asset.Name = info.Title
			asset.PropertyTypeID = info.PropertyTypeID
		}
		if loc := bundle.Location; loc != nil {
			asset.ProvinceID = loc.ProvinceID
			asset.WardID = loc.WardID
			asset.Address = loc.AddressDetail
		}
		if land := bundle.LandInfo; land != nil && land.AreaTotal != nil {
			asset.Area = *land.AreaTotal
		}
	}

	for _, item := range a.GetLegalItems() {
		if item == nil {
			continue
		}
		asset.LegalItems = append(asset.LegalItems, entity.AssetLegal{
			DocumentName: item.GetDocumentName(),
			DocumentURL:  item.GetDocumentUrl(),
			DocumentType: item.GetDocumentType(),
			IssuedDate:   _utils.ParseStringToTime(item.GetIssuedDate()),
			ExpiryDate:   _utils.ParseStringToTime(item.GetExpiryDate()),
			Description:  item.GetDescription(),
		})
	}
	for _, item := range req.GetLegalInfo() {
		if item == nil {
			continue
		}
		legal := entity.AssetLegal{
			DocumentName:        item.GetDocumentName(),
			DocumentURL:         item.GetDocumentUrl(),
			DocumentType:        item.GetDocumentType(),
			Description:         item.GetDescription(),
			RelatedSplitMergeID: item.RelatedSplitMergeId,
		}
		if item.IssuedDate != nil {
			t := time.UnixMilli(*item.IssuedDate)
			legal.IssuedDate = &t
		}
		if item.ExpiryDate != nil {
			t := time.UnixMilli(*item.ExpiryDate)
			legal.ExpiryDate = &t
		}
		asset.LegalItems = append(asset.LegalItems, legal)
	}
	return asset
}
