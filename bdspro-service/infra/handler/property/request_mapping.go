package property_handler

import (
	"bdspro/internal/dto"
	"common/pkg/fieldmask"
	bdspropb "pb/types/bdspro"
)

func MapUpdateRequestToCmd(req *bdspropb.UpdatePropertyRequest) *dto.UpdatePropertyDTO {
	cmd := &dto.UpdatePropertyDTO{
		LineageID: req.Id,
		Mask:      fieldmask.FromPaths(req.GetFieldMask().GetPaths()),
	}

	// Impact assessment
	if req.ImpactAssessment != nil {
		cmd.ImpactAssessment = &dto.ImpactAssessmentDTO{
			PreviewOnly: req.ImpactAssessment.PreviewOnly,
			Confirmed:   req.ImpactAssessment.Confirmed,
			SessionID:   req.ImpactAssessment.SessionId,
		}
	}

	// Sub-domains
	cmd.Info = mapInfoPayload(req.Info)
	cmd.Location = mapLocationPayload(req.Location)
	cmd.LandInfo = mapLandInfoPayload(req.LandInfo)
	cmd.BuildingInfo = mapBuildingPayload(req.BuildingInfo)
	cmd.Evidence = mapEvidencePayload(req.Evidence)

	// Side blocks
	cmd.Media = mapMediaPayload(req.Media)
	cmd.Amenities = mapAmenityPayload(req.Amenities)
	cmd.AreaRegions = mapAreaRegionPayload(req.AreaRegions)
	cmd.ExternalRefs = mapExternalRefPayload(req.ExternalRefs)
	cmd.LineageMeta = mapLineageMetaPayload(req.LineageMeta)
	cmd.Personalization = mapPersonalizationPayload(req.Personalization)

	return cmd
}

// MapPreviewRequestToCmd chuyển đổi PreviewPropertyUpdateRequest sang UpdatePropertyDTO
func MapPreviewRequestToCmd(req *bdspropb.PreviewPropertyUpdateRequest) *dto.UpdatePropertyDTO {
	cmd := &dto.UpdatePropertyDTO{
		LineageID: req.Id,
		Mask:      fieldmask.FromPaths(req.GetFieldMask().GetPaths()),
		// Preview mode will be set by handler
	}

	// Map only fields needed for preview (can reuse full mapping or subset)
	cmd.Info = mapInfoPayload(req.Info)
	cmd.Location = mapLocationPayload(req.Location)
	cmd.LandInfo = mapLandInfoPayload(req.LandInfo)
	cmd.BuildingInfo = mapBuildingPayload(req.BuildingInfo)
	cmd.Evidence = mapEvidencePayload(req.Evidence)
	cmd.Media = mapMediaPayload(req.Media)
	cmd.Amenities = mapAmenityPayload(req.Amenities)
	cmd.AreaRegions = mapAreaRegionPayload(req.AreaRegions)
	cmd.ExternalRefs = mapExternalRefPayload(req.ExternalRefs)
	cmd.LineageMeta = mapLineageMetaPayload(req.LineageMeta)
	cmd.Personalization = mapPersonalizationPayload(req.Personalization)

	return cmd
}

// MapReportRequestToCmd chuyển đổi SendReportPropertyRequest sang SubmitReportDTO
func MapReportRequestToCmd(req *bdspropb.SendReportPropertyRequest) *dto.SubmitReportDTO {
	cmd := &dto.SubmitReportDTO{
		LineageID:          req.Id,
		ReportType:         req.Type,
		SubIssues:          req.SubIssues,
		Note:               req.Note,
		RelatedLineageID:   req.RelatedLineageId,
		AttachmentMediaIDs: req.AttachmentMediaIds,
		AttachmentFileIDs:  req.AttachmentFileIds,
	}
	if req.ContributeData != nil {
		cmd.ContributeData = MapUpdateRequestToCmd(req.ContributeData)
	}
	return cmd
}

// =============================================================================
// PAYLOAD MAPPERS (reusable)
// =============================================================================

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
		WardID:        p.WardId,
	}
}

func mapLandInfoPayload(p *bdspropb.UpdateLandInfoPayload) *dto.UpdateLandInfoDTO {
	if p == nil {
		return nil
	}
	return &dto.UpdateLandInfoDTO{
		DocumentNo:   p.DocumentNo,
		IssuringAuth: p.IssuringAuth,
		DocumentType: p.DocumentType,
		Plot:         p.Plot,
		Sheet:        p.Sheet,
		AreaTotal:    p.AreaTotal,
		AreaLand:     p.AreaLand,
		AreaPlant:    p.AreaPlant,
		FrontWidth:   p.FrontWidth,
		Depth:        p.Depth,
		StreetWidth:  p.StreetWidth,
		ExpiredLand:  p.ExpiredLand,
		ExpiredPlant: p.ExpiredPlant,
		Note:         p.Note,
		LandNote:     p.LandNote,
		PurposeUsed:  p.PurposeUsed,
	}
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

func mapPersonalizationPayload(p *bdspropb.UpdatePersonalizationPayload) *dto.UpdatePersonalizationDTO {
	if p == nil {
		return nil
	}
	return &dto.UpdatePersonalizationDTO{
		ArchivedAt: p.ArchivedAt,
		HiddenAt:   p.HiddenAt,
	}
}

// =============================================================================
// RESPONSE MAPPERS
// =============================================================================

// MapUpdatePropertyResultToProto chuyển đổi UpdatePropertyResult sang UpdatePropertyResponse
func MapUpdatePropertyResultToProto(result *dto.UpdatePropertyResult) *bdspropb.UpdatePropertyResponse {
	if result == nil {
		return nil
	}
	resp := &bdspropb.UpdatePropertyResponse{
		LineageId:       result.LineageID,
		UpdatedAt:       result.UpdatedAt,
		RequiresConfirm: result.RequiresConfirm,
		Message:         "updated successfully",
	}
	if result.ImpactEvaluation != nil {
		resp.ImpactEvaluation = MapImpactEvaluationToProto(result.ImpactEvaluation)
	}
	return resp
}

// MapPreviewResultToProto chuyển đổi UpdatePropertyResult sang PreviewPropertyUpdateResponse
func MapPreviewResultToProto(result *dto.UpdatePropertyResult) *bdspropb.PreviewPropertyUpdateResponse {
	if result == nil {
		return nil
	}
	return &bdspropb.PreviewPropertyUpdateResponse{
		LineageId:        result.LineageID,
		RequiresConfirm:  result.RequiresConfirm,
		ImpactEvaluation: MapImpactEvaluationToProto(result.ImpactEvaluation),
	}
}

// MapReportResultToProto chuyển đổi SubmitReportResult sang ReportPropertyResponse
func MapReportResultToProto(result *dto.SubmitReportResult) *bdspropb.ReportPropertyResponse {
	if result == nil {
		return nil
	}
	return &bdspropb.ReportPropertyResponse{
		ReportId:         result.ReportID,
		LineageId:        result.LineageID,
		ProposeLineageId: derefUint64(result.ProposeLineageID),
		Type:             result.ReportType,
		Note:             result.Note,
		CreatedAt:        result.CreatedAt,
		Message:          "report submitted successfully",
	}
}

// MapImpactEvaluationToProto chuyển đổi ImpactEvaluationResponse sang proto ImpactEvaluation
func MapImpactEvaluationToProto(eval *dto.ImpactEvaluationResponse) *bdspropb.ImpactEvaluation {
	if eval == nil {
		return nil
	}
	proto := &bdspropb.ImpactEvaluation{
		ImpactLevel:     eval.ImpactLevel,
		ActionPolicy:    eval.ActionPolicy,
		RequiresConfirm: eval.RequiresConfirm,
		BlockReason:     eval.BlockReason,
		Warnings:        eval.Warnings,
		Suggestions:     eval.Suggestions,
	}
	if eval.AffectedSummary != nil {
		proto.AffectedSummary = &bdspropb.AffectedSummary{
			Products: int32(eval.AffectedSummary.Products),
			Listings: int32(eval.AffectedSummary.Listings),
			Assets:   int32(eval.AffectedSummary.Assets),
			Deals:    int32(eval.AffectedSummary.Deals),
			CrmNotes: int32(eval.AffectedSummary.CrmNotes),
			Total:    int32(eval.AffectedSummary.Total),
		}
	}
	if eval.ProjectedStates != nil {
		proto.ProjectedStates = &bdspropb.ProjectedStates{
			Active:     int32(eval.ProjectedStates.Active),
			Restricted: int32(eval.ProjectedStates.Restricted),
			Frozen:     int32(eval.ProjectedStates.Frozen),
			Archived:   int32(eval.ProjectedStates.Archived),
		}
	}
	if len(eval.ChangedFields) > 0 {
		proto.ChangedFields = make([]*bdspropb.FieldChange, len(eval.ChangedFields))
		for i, f := range eval.ChangedFields {
			proto.ChangedFields[i] = &bdspropb.FieldChange{
				FieldPath:  f.FieldPath,
				OldValue:   f.OldValue,
				NewValue:   f.NewValue,
				ChangeType: f.ChangeType,
			}
		}
	}
	return proto
}

// Helper
func derefUint64(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}
