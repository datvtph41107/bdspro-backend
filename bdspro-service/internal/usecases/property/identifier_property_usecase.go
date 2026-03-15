package property_usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
)

// IdentifierProperty tạo bảng ghi property_identify, gán ID vào lineage và các bảng thuộc tính (info, location, land_info, building_info, evidence, external_ref, media).
func (u *PropertyUsecase) IdentifierProperty(ctx context.Context, lineageID uint64) (identifyID uint64, err error) {
	lineage, err := u.LineageRepo.GetByID(ctx, lineageID)
	if err != nil {
		return 0, err
	}
	if lineage == nil {
		return 0, _errors.NotFoundException("lineage not found")
	}
	if lineage.PropertyIdentifyID != nil {
		identifyID = *lineage.PropertyIdentifyID
		identify, err := u.PropertyIdentifyRepo.GetByID(ctx, identifyID)
		if err == nil && identify != nil && (identify.LineageID == nil || *identify.LineageID != lineageID) {
			_ = u.PropertyIdentifyRepo.UpdateFields(ctx, identifyID, map[string]any{"lineage_id": lineageID})
		}
		return identifyID, nil
	}

	identify := &domain.PropertyIdentify{
		PID:       lineageID,
		Version:   1,
		LineageID: &lineageID,
	}
	if err := u.PropertyIdentifyRepo.Create(ctx, identify); err != nil {
		return 0, fmt.Errorf("create property_identify: %w", err)
	}
	identifyID = identify.ID

	if err := u.LineageRepo.UpdateFields(ctx, lineageID, map[string]any{"property_identify_id": identifyID}); err != nil {
		return 0, fmt.Errorf("update lineage: %w", err)
	}

	if lineage.PropertyInfoID != nil {
		_ = u.InfoRepo.UpdateFields(ctx, *lineage.PropertyInfoID, map[string]any{"property_identify_id": identifyID})
	}
	if lineage.LocationID != nil {
		_ = u.LocationRepo.UpdateFields(ctx, *lineage.LocationID, map[string]any{"property_identify_id": identifyID})
	}
	if lineage.LandInfoID != nil {
		_ = u.LandInfoRepo.UpdateFields(ctx, *lineage.LandInfoID, map[string]any{"property_identify_id": identifyID})
	}
	if lineage.BuildingInfoID != nil {
		_ = u.BuildingInfoRepo.UpdateFields(ctx, *lineage.BuildingInfoID, map[string]any{"property_identify_id": identifyID})
	}
	if lineage.EdvidenceID != nil {
		_ = u.EvidenceRepo.UpdateFields(ctx, *lineage.EdvidenceID, map[string]any{"property_identify_id": identifyID})
	}

	if refs, err := u.ExternalRefRepo.GetByLineageID(ctx, lineageID); err == nil {
		for _, ref := range refs {
			_ = u.ExternalRefRepo.UpdateFields(ctx, ref.ID, map[string]any{"property_identify_id": identifyID})
		}
	}
	if mediaList, err := u.MediaRepo.GetByLineageID(ctx, lineageID); err == nil {
		for _, m := range mediaList {
			_ = u.MediaRepo.UpdateFields(ctx, m.ID, map[string]any{"property_identify_id": identifyID})
		}
	}

	return identifyID, nil
}

// ImportProperty tạo identifier trước, rồi lineage, sau đó update identifier.lineage_id. Dùng cho import data.
func (u *PropertyUsecase) ImportProperty(ctx context.Context, req *dto.CreatePropertyProductRequest) (*domain.PropertyLineage, error) {
	var (
		property    *domain.PropertyLineage
		identifyID  uint64
		locationID  *uint64
		infoID      *uint64
		landInfoID  *uint64
		buildingID  *uint64
		edvidenceID *uint64
	)
	originID := _utils.GetOriginIdFromContext(ctx)

	err := u.Transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Tạo PropertyIdentify trước (chưa có lineage_id)
		identify := &domain.PropertyIdentify{
			PID:     0,
			Version: 1,
		}
		if err := u.PropertyIdentifyRepo.Create(txCtx, identify); err != nil {
			return fmt.Errorf("create property_identify: %w", err)
		}
		identifyID = identify.ID
		identifyIDPtr := &identifyID

		// 2. Tạo Info (có property_identify_id)
		info := &domain.PropertyInfo{
			SourceType:  enums.EPropertySourceType(enums.SourceTypeOther),
			LegalStatus: enums.EHouseCertificateNone,
			Note:        "Note",
			UnitCode:    "UnitCode",
			Identifier:  "Identifier",
			Level:       "Level",
		}
		if req.Info != nil {
			info = req.Info
		}
		info.PropertyIdentifyID = identifyIDPtr
		if err := u.PropertyInfoRepo.Save(txCtx, info); err != nil {
			return _errors.InternalServerException("Create PropertyInfo: %s", err.Error())
		}
		infoID = &info.ID

		// 3. Location
		if req.Location != nil {
			loc := *req.Location
			loc.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyLocationRepo.Create(txCtx, &loc); err != nil {
				return _errors.InternalServerException("Create PropertyLocation: %s", err.Error())
			}
			locationID = &loc.ID
		}

		// 4. LandInfo
		if req.LandInfo != nil {
			land := *req.LandInfo
			land.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyLandInfoRepo.Create(txCtx, &land); err != nil {
				return _errors.InternalServerException("Create PropertyLandInfo: %s", err.Error())
			}
			landInfoID = &land.ID
		}

		// 5. BuildingInfo
		if req.BuildingInfo != nil {
			bld := *req.BuildingInfo
			bld.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyBuildingInfoRepo.Create(txCtx, &bld); err != nil {
				return _errors.InternalServerException("Create PropertyBuildingInfo: %s", err.Error())
			}
			buildingID = &bld.ID
		}

		// 6. Edvidence
		if req.Edvidence != nil {
			edv := *req.Edvidence
			edv.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyEdvidenceRepo.Create(txCtx, &edv); err != nil {
				return _errors.InternalServerException("Create PropertyEdvidence: %s", err.Error())
			}
			edvidenceID = &edv.ID
		}

		// 7. ExternalRef
		if req.ExternalRef != nil {
			ext := *req.ExternalRef
			ext.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyExternalRefRepo.Create(txCtx, &ext); err != nil {
				return _errors.InternalServerException("Create PropertyExternalRef: %s", err.Error())
			}
		}

		// 8. Statistic
		statistic := &domain.PropertyStatistic{
			ProductCount: 0,
			AssetCount:   0,
			ListingCount: 0,
		}
		if err := u.PropertyStatisticRepo.Save(txCtx, statistic); err != nil {
			return _errors.InternalServerException("Create PropertyStatistic: %s", err.Error())
		}
		statisticID := &statistic.ID

		// 9. Tạo Lineage (đã có property_identify_id)
		property = &domain.PropertyLineage{
			PropertyIdentifyID: identifyIDPtr,
			PropertyInfoID:     infoID,
			LocationID:         locationID,
			LandInfoID:         landInfoID,
			BuildingInfoID:     buildingID,
			EdvidenceID:        edvidenceID,
			StatisticID:        statisticID,
		}
		if err := u.PropertyRepo.Create(txCtx, property); err != nil {
			return _errors.InternalServerException("Create PropertyLineage: %s", err.Error())
		}

		// 10. Update identifier: gán lineage_id
		if err := u.PropertyIdentifyRepo.UpdateFields(txCtx, identifyID, map[string]any{"lineage_id": property.ID}); err != nil {
			return fmt.Errorf("update property_identify lineage_id: %w", err)
		}
		// Update PID = lineageID (optional, để đồng bộ)
		_ = u.PropertyIdentifyRepo.UpdateFields(txCtx, identifyID, map[string]any{"pid": property.ID})

		// 11. PropertyUser (owner)
		if originID > 0 {
			_, _ = u.PropertyUserRepo.CreateOwner(txCtx, property.ID, originID)
		}

		// 12. Media (sau khi có lineage)
		if len(req.MediaList) > 0 {
			mediaList := make([]domain.PropertyMedia, len(req.MediaList))
			for i := range req.MediaList {
				mediaList[i] = req.MediaList[i]
				mediaList[i].LineageID = property.ID
				mediaList[i].PropertyIdentifyID = identifyIDPtr
			}
			if err := u.PropertyMediaRepo.CreateBatch(txCtx, mediaList); err != nil {
				return _errors.InternalServerException("Create PropertyMedia: %s", err.Error())
			}
			if len(mediaList) > 0 {
				info.AvatarID = &mediaList[0].ID
				_ = u.PropertyInfoRepo.Save(txCtx, info)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return property, nil
}
