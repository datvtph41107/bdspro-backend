package usecase

import (
	"fmt"
	"strconv"

	"common/fault"
)

func qhLandUsePayloadRequired() error {
	return fault.Validation(
		"tqd.land_use.payload_required",
		"payload is required",
		fault.FieldViolation{
			Field:       "payload",
			Description: "is required",
		},
	)
}

func qhLandUseIDRequired() error {
	return fault.Validation(
		"tqd.land_use.id_required",
		"id is required",
		fault.FieldViolation{
			Field:       "id",
			Description: "must be greater than zero",
		},
	)
}

func qhLandUseLayerIDRequired() error {
	return fault.Validation(
		"tqd.land_use.layer_id_required",
		"layerId is required",
		fault.FieldViolation{
			Field:       "layer_id",
			Description: "must be greater than zero",
		},
	)
}

func qhLandUseIDValueRequired() error {
	return fault.Validation(
		"tqd.land_use.land_use_id_required",
		"landUseId is required",
		fault.FieldViolation{
			Field:       "land_use_id",
			Description: "must be greater than zero",
		},
	)
}

func qhLandUseGroupIDRequired() error {
	return fault.Validation(
		"tqd.land_use.group_id_required",
		"landUseGroupId is required",
		fault.FieldViolation{
			Field:       "land_use_group_id",
			Description: "must be greater than zero",
		},
	)
}

func qhLandUseNameRequired() error {
	return fault.Validation(
		"tqd.land_use.name_required",
		"name is required",
		fault.FieldViolation{
			Field:       "name",
			Description: "is required",
		},
	)
}

func qhLandUseLayerNotFound(layerID uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.land_use.layer_not_found",
		fmt.Sprintf("layer %d not found", layerID),
	).WithMetadata(map[string]string{
		"layer_id": strconv.FormatUint(layerID, 10),
	})
}

func qhLandUseNotFound(landUseID uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.land_use.not_found",
		fmt.Sprintf("land use %d not found", landUseID),
	).WithMetadata(map[string]string{
		"land_use_id": strconv.FormatUint(landUseID, 10),
	})
}

func qhLandUseRecordNotFound(id uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.land_use.record_not_found",
		fmt.Sprintf("record %d not found", id),
	).WithMetadata(map[string]string{
		"land_use_id": strconv.FormatUint(id, 10),
	})
}

func qhLandUseDuplicate(layerID, landUseID uint64) error {
	return fault.New(
		fault.KindConflict,
		"tqd.land_use.duplicate",
		fmt.Sprintf(
			"layer %d already has land use %d",
			layerID,
			landUseID,
		),
	).WithMetadata(map[string]string{
		"layer_id":    strconv.FormatUint(layerID, 10),
		"land_use_id": strconv.FormatUint(landUseID, 10),
	})
}
