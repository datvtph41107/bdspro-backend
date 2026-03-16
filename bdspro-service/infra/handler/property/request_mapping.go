package property_handler

import (
	"bdspro/internal/dto"
	"common/pkg/fieldmask"
	bdspropb "pb/types/bdspro"
)

func mapUpdateRequestToCmd(req *bdspropb.UpdatePropertyRequest) *dto.UpdatePropertyDTO {
	return &dto.UpdatePropertyDTO{
		LineageID: req.Id,
		// FieldMask: khởi tạo 1 lần — value object truyền xuống usecase qua DTO
		Mask:         fieldmask.FromPaths(req.GetFieldMask().GetPaths()),
		Info:         mapInfoPayload(req.Info),
		Location:     mapLocationPayload(req.Location),
		LandInfo:     mapLandInfoPayload(req.LandInfo),
		BuildingInfo: mapBuildingPayload(req.BuildingInfo),
		Evidence:     mapEvidencePayload(req.Evidence),
		Media:        mapMediaPayload(req.Media),
		Amenities:    mapAmenityPayload(req.Amenities),
		ExternalRefs: mapExternalRefPayload(req.ExternalRefs),
		AreaRegions:  mapAreaRegionPayload(req.AreaRegions),
		LineageMeta:  mapLineageMetaPayload(req.LineageMeta),
	}
}

func mapInfoPayload(p *bdspropb.UpdateInfoPayload) *dto.UpdateInfoDTO {
	if p == nil {
		return nil
	}
	return &dto.UpdateInfoDTO{
		Title:           p.Title,
		Note:            p.Note,
		UnitCode:        p.UnitCode,
		Identifier:      p.Identifier,
		Level:           p.Level,
		Scope:           p.Scope,
		SourceType:      p.SourceType,
		LegalStatus:     p.LegalStatus,
		RecordStatus:    p.RecordStatus,
		AvatarID:        p.AvatarId,
		PropertyTypeID:  p.PropertyTypeId,
		ProjectID:       p.ProjectId,
		OriginProfileID: p.OriginProfileId,
	}
}

func mapLocationPayload(p *bdspropb.UpdateLocationPayload) *dto.UpdateLocationDTO {
	if p == nil {
		return nil
	}
	return &dto.UpdateLocationDTO{
		AddressDetail: p.AddressDetail,
		MapURL:        p.MapUrl,
		Latitude:      p.Latitude,
		Longitude:     p.Longitude,
		RegionID:      p.RegionId,
		ProvinceID:    p.ProvinceId,
		DistrictID:    p.DistrictId,
		WardID:        p.WardId,
	}
}

func mapLandInfoPayload(p *bdspropb.UpdateLandInfoPayload) *dto.UpdateLandInfoDTO {
	if p == nil {
		return nil
	}
	c := &dto.UpdateLandInfoDTO{
		Note:         p.Note,
		LandNote:     p.LandNote,
		PurposeUsed:  p.PurposeUsed,
		ExpiredLand:  p.ExpiredLand,
		ExpiredPlant: p.ExpiredPlant,
		AreaTotal:    p.AreaTotal,
		AreaLand:     p.AreaLand,
		AreaPlant:    p.AreaPlant,
		FrontWidth:   p.FrontWidth,
		Depth:        p.Depth,
		StreetWidth:  p.StreetWidth,
		DocumentNo:   p.DocumentNo,
		IssuringAuth: p.IssuringAuth,
		DocumentType: p.DocumentType,
		Plot:         p.Plot,
		Sheet:        p.Sheet,
	}

	return c
}

func mapBuildingPayload(p *bdspropb.UpdateBuildingPayload) *dto.UpdateBuildingDTO {
	if p == nil {
		return nil
	}
	return &dto.UpdateBuildingDTO{
		Note:             p.Note,
		BuildStatus:      p.BuildStatus,
		BuildingType:     p.BuildingType,
		Direction:        p.Direction,
		BalconyDir:       p.BalconyDirection,
		AreaActual:       p.AreaActual,
		AreaFloor:        p.AreaFloor,
		AreaConstruction: p.AreaConstruction,
		Floors:           p.Floors,
		RoomNumber:       p.RoomNumber,
		Bedrooms:         p.Bedrooms,
		Bathrooms:        p.Bathrooms,
	}
}

func mapEvidencePayload(p *bdspropb.UpdateEvidencePayload) *dto.UpdateEvidenceDTO {
	if p == nil {
		return nil
	}
	return &dto.UpdateEvidenceDTO{
		Title:       p.Title,
		Description: p.Description,
		FileID:      p.FileId,
	}
}

func mapMediaPayload(p *bdspropb.UpdateMediaPayload) *dto.UpdateMediaDTO {
	if p == nil {
		return nil
	}
	items := make([]dto.UpdateMediaItemDTO, 0, len(p.Items))
	for _, item := range p.Items {
		items = append(items, dto.UpdateMediaItemDTO{
			ID:        item.Id,
			MediaType: item.MediaType,
			MediaURL:  item.MediaUrl,
			ThumbURL:  item.ThumbUrl,
			SortOrder: item.SortOrder,
			IsCover:   item.IsCover,
		})
	}
	return &dto.UpdateMediaDTO{Items: items}
}

func mapAmenityPayload(p *bdspropb.UpdateAmenityPayload) *dto.UpdateAmenityDTO {
	if p == nil {
		return nil
	}
	return &dto.UpdateAmenityDTO{AmenityIDs: p.AmenityIds}
}

func mapAreaRegionPayload(p *bdspropb.UpdateAreaRegionPayload) *dto.UpdateAreaRegionsDTO {
	if p == nil {
		return nil
	}
	return &dto.UpdateAreaRegionsDTO{AreaRegionIDs: p.AreaRegionIds}
}

func mapExternalRefPayload(p *bdspropb.UpdateExternalRefPayload) *dto.UpdateExternalRefDTO {
	if p == nil {
		return nil
	}
	items := make([]dto.UpdateExternalRefItemDTO, 0, len(p.Items))
	for _, item := range p.Items {
		// Scalar copy — external_ref upsert toàn bộ field, không partial
		items = append(items, dto.UpdateExternalRefItemDTO{
			ID:            item.Id,
			ExternalRefID: item.ExternalRefId,
			SourceSystem:  item.SourceSystem,
			SourceCode:    item.SourceCode,
			Confidence:    item.Confidence,
			SyncStatus:    item.SyncStatus,
			Note:          item.Note,
		})
	}
	return &dto.UpdateExternalRefDTO{Items: items}
}

func mapLineageMetaPayload(p *bdspropb.UpdateLineageMetaPayload) *dto.UpdateLineageMetaDTO {
	if p == nil {
		return nil
	}
	return &dto.UpdateLineageMetaDTO{
		NationalID:         p.NationalId,
		LineageStatus:      p.LineageStatus,
		VerifiedNationalAt: p.VerifiedNationalAt,
	}
}

func mapReportRequestToCmd(req *bdspropb.SendReportPropertyRequest) *dto.SubmitReportDTO {
	cmd := &dto.SubmitReportDTO{
		LineageID:          req.Id,
		ReportType:         int32(req.Type),
		SubIssues:          req.SubIssues,
		AttachmentMediaIDs: req.AttachmentMediaIds,
		AttachmentFileIDs:  req.AttachmentFileIds,
	}

	if req.Note != nil {
		cmd.Note = req.Note
	}

	if req.RelatedLineageId != nil {
		id := *req.RelatedLineageId
		cmd.RelatedLineageID = &id
	}

	if req.ContributeData != nil {
		updateDTO := mapUpdateRequestToCmd(req.ContributeData)
		updateDTO.LineageID = req.Id
		cmd.ContributeData = updateDTO
	}

	return cmd
}
