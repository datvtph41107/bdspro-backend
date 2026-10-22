package repository

import (
	"context"
	"time"

	_errors "common/errors"

	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/enums"

	"gorm.io/gorm"
)

// UserCreate persists a new property lineage (parity bdspro PropertyUsecase.UserCreate).
func (r *propertyPostgres) UserCreate(ctx context.Context, originID uint64, bundle *entity.PropertyCreateBundle) (*entity.PropertyLineage, error) {
	if r.db == nil {
		return nil, _errors.InternalServerException("database unavailable")
	}
	if originID == 0 {
		return nil, _errors.UnauthorizedException("origin id is required")
	}
	if bundle == nil {
		bundle = &entity.PropertyCreateBundle{}
	}

	var property *entity.PropertyLineage
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		info := bundle.Info
		if info == nil {
			info = &entity.PropertyInfo{
				SourceType:  uint32(enums.SourceTypeOther),
				LegalStatus: uint32(enums.EHouseCertificateNone),
				Note:        "Note",
				UnitCode:    "UnitCode",
				Identifier:  "Identifier",
				Level:       "Level",
				Scope:       10,
			}
		}
		if err := createSerialEntity(tx, "property_info", info); err != nil {
			return _errors.InternalServerException("Create PropertyInfo Error: ", err.Error())
		}

		var locationID, infoID, landInfoID, buildingID, edvidenceID *uint64
		infoID = &info.ID

		if bundle.Location != nil {
			loc := bundle.Location
			if err := createSerialEntity(tx, "property_location", loc); err != nil {
				return _errors.InternalServerException("Create PropertyLocation Error: ", err.Error())
			}
			locationID = &loc.ID
		}

		if bundle.LandInfo != nil {
			land := bundle.LandInfo
			if err := createSerialEntity(tx, "property_land_info", land); err != nil {
				return _errors.InternalServerException("Create PropertyLandInfo Error: ", err.Error())
			}
			landInfoID = &land.ID
		}

		if bundle.BuildingInfo != nil {
			bi := bundle.BuildingInfo
			if err := createSerialEntity(tx, "property_building_info", bi); err != nil {
				return _errors.InternalServerException("Create PropertyBuildingInfo Error: ", err.Error())
			}
			buildingID = &bi.ID
		}

		if bundle.Edvidence != nil {
			ev := bundle.Edvidence
			if err := createSerialEntity(tx, "property_edvidence", ev); err != nil {
				return _errors.InternalServerException("Create PropertyEdvidence Error: ", err.Error())
			}
			edvidenceID = &ev.ID
		}

		productCount := int64(0)
		assetCount := int64(0)
		listCount := int64(0)
		if bundle.IncludeProduct {
			productCount = 1
		}
		if bundle.IncludeAsset {
			assetCount = 1
		}
		if bundle.IncludePost {
			listCount = 1
		}

		stat := &entity.PropertyStatistic{
			ProductCount: productCount,
			AssetCount:   assetCount,
			ListingCount: listCount,
		}
		if err := createSerialEntity(tx, "property_statistic", stat); err != nil {
			return _errors.InternalServerException("Create PropertyStatistic Error: ", err.Error())
		}
		statisticID := &stat.ID

		property = &entity.PropertyLineage{
			OriginProfileID: &originID,
			LineageStatus:   10,
			PropertyInfoID:  infoID,
			LocationID:      locationID,
			LandInfoID:      landInfoID,
			BuildingInfoID:  buildingID,
			EdvidenceID:     edvidenceID,
			StatisticID:     statisticID,
			NationalID:      bundle.NationalID,
		}
		if err := createSerialEntity(tx, "property_lineage", property); err != nil {
			return _errors.InternalServerException("Create PropertyLineage Error: ", err.Error())
		}

		if bundle.ExternalRef != nil {
			ext := bundle.ExternalRef
			ext.LineageID = property.ID
			if err := createSerialEntity(tx, "property_external_ref", ext); err != nil {
				return _errors.InternalServerException("Create PropertyExternalRef Error: ", err.Error())
			}
		}

		now := time.Now()
		propUser := &entity.PropertyUser{
			OwnerOriginID:     originID,
			PropertyLineageID: property.ID,
			OriginID:          &originID,
			OwneredAt:         &now,
		}
		if err := createSerialEntity(tx, "property_user", propUser); err != nil {
			return _errors.InternalServerException("Create PropertyUser Error: ", err.Error())
		}

		if bundle.Asset != nil {
			if err := createAssetForProperty(tx, property.ID, bundle.Asset); err != nil {
				return err
			}
		}

		if len(bundle.MediaList) > 0 {
			for i := range bundle.MediaList {
				bundle.MediaList[i].LineageID = property.ID
			}
			for i := range bundle.MediaList {
				if err := createSerialEntity(tx, "property_media", &bundle.MediaList[i]); err != nil {
					return _errors.InternalServerException("Create PropertyMedia Error: ", err.Error())
				}
			}
			avatarID := bundle.MediaList[0].ID
			if err := tx.Model(info).Update("avatar_id", avatarID).Error; err != nil {
				return _errors.InternalServerException("Update PropertyInfo avatar Error: ", err.Error())
			}
		}

		for _, tagID := range bundle.AmenityIDs {
			if tagID == 0 {
				continue
			}
			link := entity.PropertyTagLink{PropertyID: property.ID, TagID: tagID}
			if err := tx.Create(&link).Error; err != nil {
				return _errors.InternalServerException("Create PropertyTagLink Error: ", err.Error())
			}
		}
		for _, tagID := range bundle.TagIDs {
			if tagID == 0 {
				continue
			}
			link := entity.PropertyTagLink{PropertyID: property.ID, TagID: tagID}
			if err := tx.Create(&link).Error; err != nil {
				return _errors.InternalServerException("Create PropertyTagLink Error: ", err.Error())
			}
		}

		property.Info = info
		return nil
	})
	if err != nil {
		return nil, err
	}
	return property, nil
}
