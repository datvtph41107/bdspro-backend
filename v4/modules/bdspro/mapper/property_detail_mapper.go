package mapper

import (
	"fmt"

	_utils "common/utils"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"mdn4/modules/bdspro/entity"
)

// PropertyLineageDetailToProto maps detail aggregate to PropertyDetailResponse (bdspro EDataModeFull parity).
func PropertyLineageDetailToProto(d *entity.PropertyLineageDetail) *bdspropb.PropertyDetailResponse {
	if d == nil || d.PropertyLineage == nil {
		return nil
	}
	p := d.PropertyLineage
	res := &bdspropb.PropertyDetailResponse{Id: p.ID}
	if p.CreatedAt != nil {
		res.CreatedAt = _utils.FormatTimeToString(p.CreatedAt)
	}
	if p.UpdatedAt != nil {
		res.UpdatedAt = _utils.FormatTimeToString(p.UpdatedAt)
	}
	if p.NationalID != "" || d.VerifiedNationalAt != nil {
		res.Lineage = &bdspropb.PropertyLineageInfo{NationalId: p.NationalID}
		if d.VerifiedNationalAt != nil {
			ms := d.VerifiedNationalAt.UnixMilli()
			res.Lineage.VerifiedNationalAt = &ms
		}
	}

	if p.Identify != nil && p.Identify.ID != 0 {
		res.Identifier = &bdspropb.PropertyIdentify{
			Id:        p.Identify.ID,
			Pid:       p.Identify.PID,
			Version:   p.Identify.Version,
			Code:      fmt.Sprintf("BDS%d", p.Identify.PID),
			LineageId: p.ID,
		}
		if p.Identify.LineageID != nil {
			res.Identifier.LineageId = *p.Identify.LineageID
		}
	}

	if p.Statistic != nil {
		res.Statistic = &bdspropb.PropertyStatistic{
			Id:           p.Statistic.ID,
			ProductCount: &p.Statistic.ProductCount,
			AssetCount:   &p.Statistic.AssetCount,
			ListingCount: &p.Statistic.ListingCount,
		}
	}

	if p.Info != nil {
		pi := p.Info
		info := &bdspropb.PropertyInfo{
			PropertyTypeId: pi.PropertyTypeID,
			ProjectId:      pi.ProjectID,
		}
		if pi.ID != 0 {
			id := pi.ID
			info.PropertyIdentifyId = &id
		}
		if pi.Title != "" {
			t := pi.Title
			info.Title = &t
		}
		if pi.Identifier != "" {
			v := pi.Identifier
			info.Identifier = &v
		}
		if pi.UnitCode != "" {
			v := pi.UnitCode
			info.UnitCode = &v
		}
		if pi.Level != "" {
			v := pi.Level
			info.Level = &v
		}
		if pi.Note != "" {
			v := pi.Note
			info.Note = &v
		}
		if pi.LegalStatus != 0 {
			v := pi.LegalStatus
			info.LegalStatus = &v
		}
		if d.Personalization != nil {
			rs := uint32(d.Personalization.ResolveRecordStatus())
			info.RecordStatus = &rs
		}
		res.Info = info

		if pi.Avatar != nil {
			res.Avatar = mediaItem(pi.Avatar)
		}
		if pi.PropertyType != nil {
			res.PropertyType = &bdspropb.PropertyType{
				Id:   pi.PropertyType.ID,
				Name: pi.PropertyType.Name,
			}
		}
		if pi.Project != nil {
			res.Project = &bdspropb.Project{
				Id:   pi.Project.ID,
				Name: pi.Project.Name,
			}
		}
	}

	if d.Developer != nil {
		dev := d.Developer
		res.Developer = &bdspropb.Developer{
			Id:      dev.ID,
			Name:    dev.Name,
			Slug:    dev.Slug,
			LogoUrl: dev.LogoURL,
			Phone:   dev.Phone,
			Email:   dev.Email,
			Website: dev.Website,
		}
	}

	if d.Personalization != nil {
		pu := d.Personalization
		lineageID := pu.PropertyLineageID
		res.Personalization = &bdspropb.PropertyPersonalization{
			Id:                pu.ID,
			OwnerOriginId:     pu.OwnerOriginID,
			PropertyLineageId: &lineageID,
		}
		if pu.RoleID != 0 {
			roleID := pu.RoleID
			res.Personalization.RoleId = &roleID
		}
		if pu.OriginID != nil {
			res.Personalization.OriginProfileId = pu.OriginID
		}
		if pu.OwneredAt != nil {
			s := _utils.FormatTimeToString(pu.OwneredAt)
			res.Personalization.OwnerAt = &s
		}
		if pu.ArchivedAt != nil {
			s := _utils.FormatTimeToString(pu.ArchivedAt)
			res.Personalization.ArchivedAt = &s
		}
		if pu.HiddenAt != nil {
			s := _utils.FormatTimeToString(pu.HiddenAt)
			res.Personalization.HiddenAt = &s
		}
		rs := uint32(pu.ResolveRecordStatus())
		res.Personalization.RecordStatus = &rs
	}

	if p.Location != nil {
		loc := p.Location
		location := &bdspropb.PropertyLocationInfo{
			Id:         loc.ID,
			Detail:     loc.AddressDetail,
			MapUrl:     loc.MapURL,
			RegionId:   loc.RegionID,
			ProvinceId: loc.ProvinceID,
			DistrictId: loc.DistrictID,
			WardId:     loc.WardID,
			Latitude:   loc.Latitude,
			Longitude:  loc.Longitude,
		}
		if loc.Province != nil {
			location.ProvinceName = loc.Province.Name
		}
		if loc.Ward != nil {
			location.WardName = loc.Ward.Name
		}
		res.Location = location
	}

	if p.LandInfo != nil {
		li := p.LandInfo
		land := &bdspropb.PropertyLandInfoDetail{
			Id:           li.ID,
			DocumentNo:   li.DocumentNo,
			IssuringAuth: li.IssuringAuth,
			Note:         li.Note,
			LandNote:     li.LandNote,
			DocumentType: li.DocumentType,
			AreaTotal:    li.AreaTotal,
			AreaLand:     li.AreaLand,
			AreaPlant:    li.AreaPlant,
			FrontWidth:   li.FrontWidth,
			Depth:        li.Depth,
			StreetWidth:  li.StreetWidth,
		}
		if li.Plot > 0 {
			v := li.Plot
			land.Plot = &v
		}
		if li.Sheet > 0 {
			v := li.Sheet
			land.Sheet = &v
		}
		if li.PurposeUsed != 0 {
			v := li.PurposeUsed
			land.PurposeUsed = &v
		}
		res.LandInfo = land
	}

	if p.BuildingInfo != nil {
		bi := p.BuildingInfo
		res.BuildingInfo = &bdspropb.PropertyBuildingInfo{
			BuildingType:     bi.BuildingType,
			AreaActual:       bi.AreaActual,
			AreaFloor:        bi.AreaFloor,
			AreaConstruction: bi.AreaConstruction,
			Floors:           bi.Floors,
			RoomNumber:       bi.RoomNumber,
			Bedrooms:         bi.Bedrooms,
			Bathrooms:        bi.Bathrooms,
			BuildStatus:      bi.BuildStatus,
			Direction:        bi.Direction,
			BalconyDirection: bi.BalconyDirection,
			Note:             bi.Note,
		}
	}

	if p.Edvidence != nil {
		ev := p.Edvidence
		res.Edvidence = &bdspropb.PropertyEdvidenceInfo{
			Id:          ev.ID,
			Title:       ev.Title,
			FileId:      ev.FileID,
			Description: ev.Description,
		}
	}

	if len(p.MediaList) > 0 {
		res.MediaList = make([]*bdspropb.MediaItem, 0, len(p.MediaList))
		for i := range p.MediaList {
			res.MediaList = append(res.MediaList, mediaItem(&p.MediaList[i]))
		}
		if res.Avatar == nil {
			res.Avatar = res.MediaList[0]
		}
	}

	if len(d.Amenities) > 0 {
		res.Amenities = make([]*sharepb.TagItem, len(d.Amenities))
		for i, a := range d.Amenities {
			res.Amenities[i] = &sharepb.TagItem{Id: a.ID, Name: a.Name}
		}
	}
	if len(d.AreaRegions) > 0 {
		res.AreaRegions = make([]*sharepb.TagItem, len(d.AreaRegions))
		for i, a := range d.AreaRegions {
			res.AreaRegions[i] = &sharepb.TagItem{Id: a.ID, Name: a.Name}
		}
	}

	return res
}

func mediaItem(m *entity.PropertyMedia) *bdspropb.MediaItem {
	if m == nil {
		return nil
	}
	return &bdspropb.MediaItem{
		Id:        m.ID,
		MediaType: m.MediaType,
		MediaUrl:  m.MediaURL,
		Order:     m.SortOrder,
	}
}
