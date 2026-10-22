package repository

import (
	_errors "common/errors"

	"mdn4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// createAssetForProperty persists an asset + legal items and links it to the property lineage
// (parity bdspro PropertyUsecase.CreateAssetFromProperty + PropertyRelation asset + AssetLegal).
func createAssetForProperty(tx *gorm.DB, propertyID uint64, asset *entity.Asset) error {
	legalItems := asset.LegalItems
	asset.LegalItems = nil

	if err := createSerialEntity(tx, "assets", asset); err != nil {
		return _errors.InternalServerException("Create Asset Error: ", err.Error())
	}

	for i := range legalItems {
		legalItems[i].AssetID = asset.ID
		if err := createSerialEntity(tx, "asset_legals", &legalItems[i]); err != nil {
			return _errors.InternalServerException("Create AssetLegal Error: ", err.Error())
		}
	}
	if len(legalItems) > 0 {
		imageID := legalItems[0].ID
		if err := tx.Model(&entity.Asset{}).Where("id = ?", asset.ID).Update("image_id", imageID).Error; err != nil {
			return _errors.InternalServerException("Update Asset image Error: ", err.Error())
		}
		asset.ImageID = &imageID
	}
	asset.LegalItems = legalItems

	assetID := asset.ID
	rel := &entity.PropertyRelation{
		PropertyID:   propertyID,
		RelationID:   &assetID,
		RelationType: entity.RelationTypeAsset,
	}
	if err := tx.Create(rel).Error; err != nil {
		return _errors.InternalServerException("Create PropertyRelation asset Error: ", err.Error())
	}
	return nil
}
