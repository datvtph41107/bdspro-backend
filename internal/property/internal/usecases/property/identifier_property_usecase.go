package property_usecases

import (
	"bdspro/internal"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_errors "common/errors"
	"common/logging"
	_utils "common/utils"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

// IdentifierProperty tạo bảng ghi property_identify, gán ID vào lineage và các bảng thuộc tính (info, location, land_info, building_info, evidence, external_ref, media).
func (u *PropertyUsecase) IdentifierProperty(ctx context.Context, lineageID uint64) (identifyID uint64, err error) {
	lineage, err := u.LineageRepo.GetByID(ctx, lineageID)
	if err != nil {
		return 0, err
	}
	if lineage == nil {
		return 0, _errors.ReturnError(service.PropertyLineageNotFound)
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

	logger := logging.FromContext(ctx)
	logger.Info(
		"starting property import",
		slog.Bool("has_location", req.Location != nil),
		slog.Bool("has_info", req.Info != nil),
		slog.Bool("has_land_info", req.LandInfo != nil),
		slog.Int("media_count", len(req.MediaList)),
	)

	if err := u.validateImportRequest(req); err != nil {
		logger.Error(
			"property import validation failed",
			slog.Any("error", err),
		)
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
		txLogger := logging.FromContext(txCtx)

		identify := &domain.PropertyIdentify{PID: 0, Version: 1}
		if err := u.PropertyIdentifyRepo.Create(txCtx, identify); err != nil {
			txLogger.Error(
				"create property identify failed",
				slog.Any("error", err),
			)
			return fmt.Errorf("create property_identify: %w", err)
		}
		identifyID = identify.ID
		identifyIDPtr := &identifyID
		txLogger.Debug(
			"created property identify",
			slog.Uint64("property_identify_id", identifyID),
		)

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
			txLogger.Error(
				"create property info failed",
				slog.Any("error", err),
			)
			return fmt.Errorf("create property info: %w", err)
		}
		infoID = &info.ID
		txLogger.Debug(
			"created property info",
			slog.Any("property_info_id", info.ID),
		)

		if req.Location != nil {
			loc := *req.Location
			loc.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyLocationRepo.Create(txCtx, &loc); err != nil {
				txLogger.Error(
					"create property location failed",
					slog.Any("error", err),
				)
				return fmt.Errorf("create property location: %w", err)
			}
			idVal := loc.ID
			locationID = &idVal
			txLogger.Debug(
				"created property location",
				slog.Any("property_location_id", loc.ID),
				slog.Float64("latitude", safeDerefFloat(loc.Latitude)),
				slog.Float64("longitude", safeDerefFloat(loc.Longitude)),
			)
		}

		if req.LandInfo != nil {
			land := *req.LandInfo
			land.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyLandInfoRepo.Create(txCtx, &land); err != nil {
				txLogger.Error(
					"create property land info failed",
					slog.Any("error", err),
				)
				return fmt.Errorf("create property land info: %w", err)
			}
			idVal := land.ID
			landInfoID = &idVal
			txLogger.Debug(
				"created property land info",
				slog.Any("property_land_info_id", land.ID),
				slog.Float64("area", safeDerefFloat(land.AreaTotal)),
			)
		}

		if req.BuildingInfo != nil {
			bld := *req.BuildingInfo
			bld.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyBuildingInfoRepo.Create(txCtx, &bld); err != nil {
				txLogger.Error(
					"create property building info failed",
					slog.Any("error", err),
				)
				return fmt.Errorf("create property building info: %w", err)
			}
			idVal := bld.ID
			buildingID = &idVal
			txLogger.Debug(
				"created property building info",
				slog.Any("property_building_info_id", bld.ID),
				slog.Any("floors", derefUint32(bld.Floors)),
			)
		}

		if req.Edvidence != nil {
			edv := *req.Edvidence
			edv.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyEdvidenceRepo.Create(txCtx, &edv); err != nil {
				txLogger.Error(
					"create property evidence failed",
					slog.Any("error", err),
				)
				return fmt.Errorf("create property evidence: %w", err)
			}
			idVal := edv.ID
			edvidenceID = &idVal
			txLogger.Debug(
				"created property evidence",
				slog.Any("property_evidence_id", edv.ID),
			)
		}

		if req.ExternalRef != nil {
			ext := *req.ExternalRef
			ext.PropertyIdentifyID = identifyIDPtr
			if err := u.PropertyExternalRefRepo.Create(txCtx, &ext); err != nil {
				txLogger.Error(
					"create property external reference failed",
					slog.Any("error", err),
				)
				return fmt.Errorf("create property external ref: %w", err)
			}
			txLogger.Debug(
				"created property external reference",
				slog.Any("property_external_ref_id", ext.ID),
			)
		}

		statistic := &domain.PropertyStatistic{
			ProductCount:       0,
			AssetCount:         0,
			ListingCount:       0,
			PropertyIdentifyID: identifyIDPtr,
		}
		if err := u.PropertyStatisticRepo.Save(txCtx, statistic); err != nil {
			txLogger.Error(
				"create property statistic failed",
				slog.Any("error", err),
			)
			return fmt.Errorf("create property statistic: %w", err)
		}
		statisticID := &statistic.ID
		txLogger.Debug(
			"created property statistic",
			slog.Any("property_statistic_id", statistic.ID),
		)

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
			txLogger.Error(
				"create property lineage failed",
				slog.Any("error", err),
			)
			return fmt.Errorf("create property lineage: %w", err)
		}
		txLogger.Info(
			"created property lineage",
			slog.Any("property_lineage_id", property.ID),
		)

		if err := u.PropertyIdentifyRepo.UpdateFields(txCtx, identifyID, map[string]any{
			"lineage_id": property.ID,
			"pid":        property.ID,
		}); err != nil {
			txLogger.Error(
				"update property identify failed",
				slog.Any("error", err),
			)
			return fmt.Errorf("update property_identify: %w", err)
		}
		txLogger.Debug(
			"updated property identify lineage",
			slog.Uint64("property_identify_id", identifyID),
			slog.Any("lineage_id", property.ID),
		)

		if originID > 0 {
			if _, err := u.PropertyUserRepo.CreateOwner(txCtx, property.ID, originID); err != nil {
				txLogger.Warn(
					"create property user failed",
					slog.Any("origin_id", originID),
					slog.Any("error", err),
				)
			} else {
				txLogger.Debug(
					"created property user",
					slog.Any("origin_id", originID),
				)
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
				txLogger.Error(
					"create property media batch failed",
					slog.Any("error", err),
				)
				return fmt.Errorf("create property media: %w", err)
			}
			txLogger.Debug(
				"created property media items",
				slog.Int("media_count", len(mediaList)),
			)

			if len(mediaList) > 0 {
				info.AvatarID = &mediaList[0].ID
				if err := u.PropertyInfoRepo.Save(txCtx, info); err != nil {
					txLogger.Warn(
						"update property avatar failed",
						slog.Any("error", err),
					)
				} else {
					txLogger.Debug(
						"updated property avatar",
						slog.Any("media_id", mediaList[0].ID),
					)
				}
			}
		}

		txLogger.Info(
			"property imported successfully",
			slog.Any("property_id", property.ID),
			slog.Uint64("property_identify_id", identifyID),
		)

		return nil
	})

	if err != nil {
		var dupErr *DuplicatePropertyError
		if errors.As(err, &dupErr) {
			logger.Warn(
				"duplicate properties detected",
				slog.Any("duplicate_property_ids", dupErr.IDs),
			)
		} else {
			logger.Error(
				"property import transaction failed",
				slog.Any("error", err),
			)
		}
		return nil, err
	}

	return property, nil
}

func (u *PropertyUsecase) validateImportRequest(req *dto.CreatePropertyProductRequest) error {
	if req == nil {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("request cannot be nil"))
	}

	// Kiểm tra Location
	if req.Location == nil {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("location is required"))
	}
	if req.Location.Latitude == nil {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("latitude is required"))
	}
	if req.Location.Longitude == nil {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("longitude is required"))
	}
	if *req.Location.Latitude < -90 || *req.Location.Latitude > 90 {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("latitude must be between -90 and 90"))
	}
	if *req.Location.Longitude < -180 || *req.Location.Longitude > 180 {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("longitude must be between -180 and 180"))
	}

	// Kiểm tra Info
	if req.Info == nil {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("info is required"))
	}
	if req.Info.PropertyTypeID == nil {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("property_type_id is required"))
	}

	// Kiểm tra diện tích
	if req.LandInfo == nil && req.BuildingInfo == nil {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("either land_info or building_info is required"))
	}

	if req.LandInfo != nil && req.LandInfo.AreaTotal != nil {
		if *req.LandInfo.AreaTotal <= 0 {
			return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("area_total must be positive"))
		}
	}

	if req.BuildingInfo != nil && req.BuildingInfo.AreaActual != nil {
		if *req.BuildingInfo.AreaActual <= 0 {
			return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("area_actual must be positive"))
		}
	}

	return nil
}

func (u *PropertyUsecase) checkDuplicates(ctx context.Context, req *dto.CreatePropertyProductRequest, identifyID uint64) error {
	logger := logging.FromContext(ctx)

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
		logger.Warn(
			"skipping duplicate check because area is missing",
			slog.Uint64("property_identify_id", identifyID),
		)
		return nil
	}

	logger.Info(
		"checking property duplicates",
		slog.Float64("latitude", *req.Location.Latitude),
		slog.Float64("longitude", *req.Location.Longitude),
		slog.Float64("area", area),
		slog.String("area_source", areaSource),
		slog.Any("property_type_id", *req.Info.PropertyTypeID),
	)

	duplicates, err := u.PropertyRepo.FindDuplicates(ctx,
		*req.Location.Latitude, *req.Location.Longitude,
		DuplicateCheckRadiusMeters,
		*req.Info.PropertyTypeID,
		area, AreaTolerancePercent)

	if err != nil {
		logger.Error(
			"property duplicate check failed",
			slog.Any("error", err),
		)
		return fmt.Errorf("duplicate check failed: %w", err)
	}

	if len(duplicates) > 0 {
		ids := make([]uint64, len(duplicates))
		for i, d := range duplicates {
			ids[i] = d.ID
		}
		logger.Warn(
			"duplicate properties found",
			slog.Int("duplicate_count", len(ids)),
			slog.Any("duplicate_property_ids", ids),
		)
		return &DuplicatePropertyError{IDs: ids}
	}

	logger.Debug("no duplicate properties found")
	return nil
}
