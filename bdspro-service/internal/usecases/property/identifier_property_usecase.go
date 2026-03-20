package property_usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"errors"
	"fmt"
	"log"
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

type DuplicatePropertyError struct {
	IDs []uint64
}

func (e *DuplicatePropertyError) Error() string {
	return fmt.Sprintf("duplicate properties found: %v", e.IDs)
}

// ImportProperty tạo identifier trước, rồi lineage, sau đó update identifier.lineage_id. Dùng cho import data.
func (u *PropertyUsecase) ImportProperty(ctx context.Context, req *dto.CreatePropertyProductRequest) (*domain.PropertyLineage, error) {
	const (
		DuplicateCheckRadiusMeters = 50.0
		AreaTolerancePercent       = 0.1
	)

	log.Printf("[INFO] Starting property import - hasLocation=%v, hasInfo=%v, hasLandInfo=%v, mediaCount=%d",
		req.Location != nil,
		req.Info != nil,
		req.LandInfo != nil,
		len(req.MediaList),
	)

	if err := u.validateImportRequest(req); err != nil {
		log.Printf("[ERROR] Validation failed: %v", err)
		return nil, err
	}

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
		identify := &domain.PropertyIdentify{PID: 0, Version: 1}
		if err := u.PropertyIdentifyRepo.Create(txCtx, identify); err != nil {
			log.Printf("[ERROR] Create property_identify failed: %v", err)
			return fmt.Errorf("create property_identify: %w", err)
		}
		identifyID = identify.ID
		identifyIDPtr := &identifyID
		log.Printf("[DEBUG] Created property_identify id=%d", identifyID)

		if err := u.checkDuplicates(txCtx, req, identifyID); err != nil {
			return err
		}

		info := &domain.PropertyInfo{
			SourceType:  enums.EPropertySourceType(enums.SourceTypeOther),
			LegalStatus: enums.EHouseCertificateNone,
		}
		if req.Info != nil {
			info.Title = req.Info.Title
			info.PropertyTypeID = req.Info.PropertyTypeID
			info.ProjectID = req.Info.ProjectID
			info.Note = req.Info.Note
			info.UnitCode = req.Info.UnitCode
			info.Identifier = req.Info.Identifier
			info.Level = req.Info.Level
		}
		info.PropertyIdentifyID = identifyIDPtr
		if err := u.PropertyInfoRepo.Save(txCtx, info); err != nil {
			log.Printf("[ERROR] Create PropertyInfo failed: %v", err)
			return _errors.InternalServerException("Create PropertyInfo: %s", err.Error())
		}
		infoID = &info.ID
		log.Printf("[DEBUG] Created property_info id=%d", info.ID)

		if req.Location != nil {
			loc := *req.Location
			loc.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyLocationRepo.Create(txCtx, &loc); err != nil {
				log.Printf("[ERROR] Create PropertyLocation failed: %v", err)
				return _errors.InternalServerException("Create PropertyLocation: %s", err.Error())
			}
			idVal := loc.ID
			locationID = &idVal
			log.Printf("[DEBUG] Created property_location id=%d, lat=%v, lng=%v",
				loc.ID, safeDerefFloat(loc.Latitude), safeDerefFloat(loc.Longitude))
		}

		if req.LandInfo != nil {
			land := *req.LandInfo
			land.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyLandInfoRepo.Create(txCtx, &land); err != nil {
				log.Printf("[ERROR] Create PropertyLandInfo failed: %v", err)
				return _errors.InternalServerException("Create PropertyLandInfo: %s", err.Error())
			}
			idVal := land.ID
			landInfoID = &idVal
			log.Printf("[DEBUG] Created property_land_info id=%d, area=%.2f", land.ID, safeDerefFloat(land.AreaTotal))
		}

		if req.BuildingInfo != nil {
			bld := *req.BuildingInfo
			bld.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyBuildingInfoRepo.Create(txCtx, &bld); err != nil {
				log.Printf("[ERROR] Create PropertyBuildingInfo failed: %v", err)
				return _errors.InternalServerException("Create PropertyBuildingInfo: %s", err.Error())
			}
			idVal := bld.ID
			buildingID = &idVal
			log.Printf("[DEBUG] Created property_building_info id=%d, floors=%v", bld.ID, derefUint32(bld.Floors))
		}

		if req.Edvidence != nil {
			edv := *req.Edvidence
			edv.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyEdvidenceRepo.Create(txCtx, &edv); err != nil {
				log.Printf("[ERROR] Create PropertyEdvidence failed: %v", err)
				return _errors.InternalServerException("Create PropertyEdvidence: %s", err.Error())
			}
			idVal := edv.ID
			edvidenceID = &idVal
			log.Printf("[DEBUG] Created property_edvidence id=%d", edv.ID)
		}

		if req.ExternalRef != nil {
			ext := *req.ExternalRef
			ext.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyExternalRefRepo.Create(txCtx, &ext); err != nil {
				log.Printf("[ERROR] Create PropertyExternalRef failed: %v", err)
				return _errors.InternalServerException("Create PropertyExternalRef: %s", err.Error())
			}
			log.Printf("[DEBUG] Created property_external_ref id=%d", ext.ID)
		}

		statistic := &domain.PropertyStatistic{
			ProductCount:       0,
			AssetCount:         0,
			ListingCount:       0,
			PropertyIdentifyID: identifyIDPtr,
		}
		if err := u.PropertyStatisticRepo.Save(txCtx, statistic); err != nil {
			log.Printf("[ERROR] Create PropertyStatistic failed: %v", err)
			return _errors.InternalServerException("Create PropertyStatistic: %s", err.Error())
		}
		statisticID := &statistic.ID
		log.Printf("[DEBUG] Created property_statistic id=%d", statistic.ID)

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
			log.Printf("[ERROR] Create PropertyLineage failed: %v", err)
			return _errors.InternalServerException("Create PropertyLineage: %s", err.Error())
		}
		log.Printf("[INFO] Created property_lineage id=%d", property.ID)

		if err := u.PropertyIdentifyRepo.UpdateFields(txCtx, identifyID, map[string]any{
			"lineage_id": property.ID,
			"pid":        property.ID,
		}); err != nil {
			log.Printf("[ERROR] Update property_identify failed: %v", err)
			return fmt.Errorf("update property_identify: %w", err)
		}
		log.Printf("[DEBUG] Updated property_identify id=%d with lineage_id=%d", identifyID, property.ID)

		if originID > 0 {
			if _, err := u.PropertyUserRepo.CreateOwner(txCtx, property.ID, originID); err != nil {
				log.Printf("[WARN] Failed to create property user for origin %d: %v", originID, err)
			} else {
				log.Printf("[DEBUG] Created property_user for origin %d", originID)
			}
		}

		// 13. Media
		if len(req.MediaList) > 0 {
			mediaList := make([]domain.PropertyMedia, len(req.MediaList))
			for i := range req.MediaList {
				mediaList[i] = req.MediaList[i]
				mediaList[i].LineageID = property.ID
				mediaList[i].PropertyIdentifyID = identifyIDPtr
			}
			if err := u.PropertyMediaRepo.CreateBatch(txCtx, mediaList); err != nil {
				log.Printf("[ERROR] Create PropertyMedia batch failed: %v", err)
				return _errors.InternalServerException("Create PropertyMedia: %s", err.Error())
			}
			log.Printf("[DEBUG] Created %d media items", len(mediaList))

			if len(mediaList) > 0 {
				info.AvatarID = &mediaList[0].ID
				if err := u.PropertyInfoRepo.Save(txCtx, info); err != nil {
					log.Printf("[WARN] Failed to update avatar: %v", err)
				} else {
					log.Printf("[DEBUG] Updated avatar to media id=%d", mediaList[0].ID)
				}
			}
		}

		log.Printf("[INFO] Property imported successfully - property_id=%d, identify_id=%d",
			property.ID, identifyID)

		return nil
	})

	if err != nil {
		var dupErr *DuplicatePropertyError
		if errors.As(err, &dupErr) {
			log.Printf("[WARN] Duplicate properties detected: %v", dupErr.IDs)
		} else {
			log.Printf("[ERROR] Transaction failed: %v", err)
		}
		return nil, err
	}

	return property, nil
}

func (u *PropertyUsecase) validateImportRequest(req *dto.CreatePropertyProductRequest) error {
	if req == nil {
		return _errors.BadRequestException("request cannot be nil")
	}

	// Kiểm tra Location
	if req.Location == nil {
		return _errors.BadRequestException("location is required")
	}
	if req.Location.Latitude == nil {
		return _errors.BadRequestException("latitude is required")
	}
	if req.Location.Longitude == nil {
		return _errors.BadRequestException("longitude is required")
	}
	if *req.Location.Latitude < -90 || *req.Location.Latitude > 90 {
		return _errors.BadRequestException("latitude must be between -90 and 90")
	}
	if *req.Location.Longitude < -180 || *req.Location.Longitude > 180 {
		return _errors.BadRequestException("longitude must be between -180 and 180")
	}

	// Kiểm tra Info
	if req.Info == nil {
		return _errors.BadRequestException("info is required")
	}
	if req.Info.PropertyTypeID == nil {
		return _errors.BadRequestException("property_type_id is required")
	}

	// Kiểm tra diện tích
	if req.LandInfo == nil && req.BuildingInfo == nil {
		return _errors.BadRequestException("either land_info or building_info is required")
	}

	if req.LandInfo != nil && req.LandInfo.AreaTotal != nil {
		if *req.LandInfo.AreaTotal <= 0 {
			return _errors.BadRequestException("area_total must be positive")
		}
	}

	if req.BuildingInfo != nil && req.BuildingInfo.AreaActual != nil {
		if *req.BuildingInfo.AreaActual <= 0 {
			return _errors.BadRequestException("area_actual must be positive")
		}
	}

	return nil
}

func (u *PropertyUsecase) checkDuplicates(ctx context.Context, req *dto.CreatePropertyProductRequest, identifyID uint64) error {
	// Kiểm tra điều kiện cần để check duplicate
	if req.Location == nil || req.Location.Latitude == nil || req.Location.Longitude == nil ||
		req.Info == nil || req.Info.PropertyTypeID == nil {
		return nil // không đủ dữ liệu để check
	}
	const (
		DuplicateCheckRadiusMeters = 50.0
		AreaTolerancePercent       = 0.1
	)
	// Xác định diện tích
	var (
		area       float64
		areaSource string
		hasArea    bool
	)

	if req.LandInfo != nil && req.LandInfo.AreaTotal != nil {
		area = *req.LandInfo.AreaTotal
		areaSource = "land_info"
		hasArea = true
	} else if req.BuildingInfo != nil && req.BuildingInfo.AreaActual != nil {
		area = *req.BuildingInfo.AreaActual
		areaSource = "building_info"
		hasArea = true
	}

	if !hasArea {
		log.Printf("[WARN] Missing area for duplicate check, skipping - property_identify_id=%d", identifyID)
		return nil
	}

	log.Printf("[INFO] Checking duplicates - lat=%.6f, lng=%.6f, area=%.2f, area_source=%s, property_type=%d",
		*req.Location.Latitude, *req.Location.Longitude, area, areaSource, *req.Info.PropertyTypeID)

	duplicates, err := u.PropertyRepo.FindDuplicates(ctx,
		*req.Location.Latitude, *req.Location.Longitude,
		DuplicateCheckRadiusMeters,
		*req.Info.PropertyTypeID,
		area, AreaTolerancePercent)

	if err != nil {
		log.Printf("[ERROR] Duplicate check failed: %v", err)
		return fmt.Errorf("duplicate check failed: %w", err)
	}

	if len(duplicates) > 0 {
		ids := make([]uint64, len(duplicates))
		for i, d := range duplicates {
			ids[i] = d.ID
		}
		log.Printf("[WARN] Found %d duplicate properties: %v", len(ids), ids)
		return &DuplicatePropertyError{IDs: ids}
	}

	log.Printf("[DEBUG] No duplicates found")
	return nil
}
