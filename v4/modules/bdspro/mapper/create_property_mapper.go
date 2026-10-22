package mapper

import (
	"time"

	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/enums"
)

// CreatePropertyRequestToBundle maps CreatePropertyRequest (pb) for UserCreate.
func CreatePropertyRequestToBundle(req *bdspropb.CreatePropertyRequest, profileID uint64) *entity.PropertyCreateBundle {
	if req == nil {
		return &entity.PropertyCreateBundle{}
	}
	out := &entity.PropertyCreateBundle{
		AmenityIDs: append([]uint64(nil), req.GetAmenityIds()...),
		TagIDs:     append([]uint64(nil), req.GetTagIds()...),
		Note:       req.GetNote(),
	}
	out.IncludeProduct = req.GetProduct() != nil
	out.IncludeAsset = req.GetAsset() != nil
	out.IncludePost = req.GetListing() != nil

	if bi := req.GetInfo(); bi != nil {
		info := &entity.PropertyInfo{
			SourceType:  uint32(enums.SourceTypeOther),
			LegalStatus: uint32(enums.EHouseCertificateNone),
			Note:        "Note",
			UnitCode:    "UnitCode",
			Identifier:  "Identifier",
			Level:       "Level",
			Scope:       10,
		}
		if bi.SourceType != nil {
			info.SourceType = *bi.SourceType
		}
		if bi.PropertyTypeId != nil {
			info.PropertyTypeID = bi.PropertyTypeId
		}
		if bi.ProjectId != nil {
			info.ProjectID = bi.ProjectId
		}
		if bi.LegalStatus != nil {
			info.LegalStatus = *bi.LegalStatus
		}
		if bi.Note != nil {
			info.Note = *bi.Note
		}
		if bi.UnitCode != nil {
			info.UnitCode = *bi.UnitCode
		}
		if bi.Identifier != nil {
			info.Identifier = *bi.Identifier
		}
		if bi.Level != nil {
			info.Level = *bi.Level
		}
		if bi.Title != nil {
			info.Title = *bi.Title
		}
		out.Info = info
	}

	if loc := req.GetLocation(); loc != nil {
		pl := &entity.PropertyLocation{
			AddressDetail: loc.GetAddressDetail(),
			MapURL:        loc.GetMapUrl(),
			PositionType:  uint32(enums.EPositionTypeRelative),
		}
		pl.RegionID = loc.RegionId
		pl.ProvinceID = loc.ProvinceId
		pl.DistrictID = loc.DistrictId
		pl.WardID = loc.WardId
		pl.Latitude = loc.Latitude
		pl.Longitude = loc.Longitude
		out.Location = pl
	}

	if li := req.GetLandInfo(); li != nil {
		land := &entity.PropertyLandInfo{
			DocumentNo:   li.GetDocumentNo(),
			IssuringAuth: li.GetIssuringAuth(),
			Note:         li.GetNote(),
			LandNote:     li.GetLandNote(),
		}
		if li.DocumentType != nil {
			land.DocumentType = uint32(enums.EDocType(*li.DocumentType))
		}
		if li.Plot != nil {
			land.Plot = *li.Plot
		}
		if li.Sheet != nil {
			land.Sheet = *li.Sheet
		}
		land.AreaTotal = li.AreaTotal
		land.AreaLand = li.AreaLand
		land.AreaPlant = li.AreaPlant
		land.FrontWidth = li.FrontWidth
		land.Depth = li.Depth
		land.StreetWidth = li.StreetWidth
		if li.PurposeUsed != nil {
			land.PurposeUsed = *li.PurposeUsed
		}
		if li.ExpiredLand != nil {
			if t := _utils.ParseStringToTime(*li.ExpiredLand); t != nil {
				// entity has no time fields on land in v4 — skip
				_ = t
			}
		}
		out.LandInfo = land
	}

	if bi := req.GetBuildingInfo(); bi != nil {
		buildingType := uint32(enums.EBuildingTypeResidential)
		if bi.BuildingType != nil && *bi.BuildingType != 0 {
			buildingType = *bi.BuildingType
		}
		building := &entity.PropertyBuildingInfo{
			BuildingType: buildingType,
			Note:         bi.GetNote(),
			BuildStatus:  uint32(enums.EBuildStatusCompleted),
		}
		building.AreaActual = bi.AreaActual
		building.AreaFloor = bi.AreaFloor
		building.AreaConstruction = bi.AreaConstruction
		building.Floors = bi.Floors
		building.RoomNumber = bi.RoomNumber
		building.Bedrooms = bi.Bedrooms
		building.Bathrooms = bi.Bathrooms
		if bi.BuildStatus != nil {
			building.BuildStatus = *bi.BuildStatus
		}
		if bi.Direction != nil && *bi.Direction != 0 {
			building.Direction = *bi.Direction
		}
		if bi.BalconyDirection != nil && *bi.BalconyDirection != 0 {
			building.BalconyDirection = *bi.BalconyDirection
		}
		out.BuildingInfo = building
	}

	if er := req.GetExternalRef(); er != nil {
		ext := &entity.PropertyExternalRef{
			SourceCode: er.GetSourceCode(),
			Note:       er.GetNote(),
		}
		if er.ExternalRefId != nil {
			ext.ExternalRefID = *er.ExternalRefId
		}
		if er.SourceSystem != nil {
			ext.SourceSystem = *er.SourceSystem
		}
		if er.Confidence != nil {
			ext.Confidence = *er.Confidence
		}
		if er.SyncStatus != nil {
			ext.SyncStatus = *er.SyncStatus
		}
		out.ExternalRef = ext
	}

	if ev := req.GetEdvidence(); ev != nil {
		out.Edvidence = &entity.PropertyEdvidence{
			Title:       ev.GetTitle(),
			Description: ev.GetDescription(),
			FileID:      ev.FileId,
		}
	}

	if ln := req.GetLineage(); ln != nil {
		out.NationalID = ln.GetNationalId()
		if ln.VerifiedNationalAt != nil {
			t := time.UnixMilli(*ln.VerifiedNationalAt)
			_ = t
		}
	}

	for _, media := range req.GetMediaList() {
		if media == nil {
			continue
		}
		out.MediaList = append(out.MediaList, entity.PropertyMedia{
			MediaType: media.GetMediaType(),
			MediaURL:  media.GetUrl(),
			ThumbURL:  media.GetThumbUrl(),
			SortOrder: media.GetSortOrder(),
		})
	}

	if req.GetProduct() != nil {
		out.Product = ProductSaveToEntity(req.GetProduct(), profileID)
	}

	if post := req.GetListing(); post != nil {
		out.Listing = PostSaveWithProductToListing(post, profileID)
	}

	out.Asset = CreatePropertyAssetToEntity(req, out, profileID)

	return out
}
