package usecase

import (
	"fmt"
	"strconv"

	"common/fault"
)

func regionLayerIDRequired() error {
	return fault.Validation(
		"tqd.region.layer_id_required",
		"layerId is required",
		fault.FieldViolation{
			Field:       "layer_id",
			Description: "must be greater than zero",
		},
	)
}

func regionLabelIDRequired() error {
	return fault.Validation(
		"tqd.region.label_id_required",
		"labelId is required",
		fault.FieldViolation{
			Field:       "label_id",
			Description: "must be greater than zero",
		},
	)
}

func regionIDRequired() error {
	return fault.Validation(
		"tqd.region.id_required",
		"id is required",
		fault.FieldViolation{
			Field:       "id",
			Description: "must be greater than zero",
		},
	)
}

func regionNameRequired() error {
	return fault.Validation(
		"tqd.region.name_required",
		"name is required",
		fault.FieldViolation{
			Field:       "name",
			Description: "is required",
		},
	)
}

func regionUpdateRequired() error {
	return fault.Validation(
		"tqd.region.update_required",
		"at least one field to update is required",
	)
}

func regionNameEmpty() error {
	return fault.Validation(
		"tqd.region.name_empty",
		"name cannot be empty",
		fault.FieldViolation{
			Field:       "name",
			Description: "cannot be empty",
		},
	)
}

func regionStatusInvalid(value uint32) error {
	return fault.Validation(
		"tqd.region.status_invalid",
		fmt.Sprintf("invalid status: %d", value),
		fault.FieldViolation{
			Field:       "status",
			Description: "unsupported region status",
		},
	)
}

func regionGeometryInvalid(cause error) error {
	message := "invalid geometry"
	if cause != nil && cause.Error() != "" {
		message = cause.Error()
	}

	return fault.Wrap(
		cause,
		fault.KindValidation,
		"tqd.region.geometry_invalid",
		message,
	)
}

func regionLayerNotFound(layerID uint64) error {
	return fault.Wrap(
		ErrLayerNotFound,
		fault.KindNotFound,
		"tqd.region.layer_not_found",
		fmt.Sprintf("layer %d not found", layerID),
	).WithMetadata(map[string]string{
		"layer_id": strconv.FormatUint(layerID, 10),
	})
}

func regionNotFound(regionID uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.region.not_found",
		fmt.Sprintf("region %d not found", regionID),
	).WithMetadata(map[string]string{
		"region_id": strconv.FormatUint(regionID, 10),
	})
}

func regionLabelNotFound(labelID uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.region.label_not_found",
		fmt.Sprintf("label %d not found", labelID),
	).WithMetadata(map[string]string{
		"label_id": strconv.FormatUint(labelID, 10),
	})
}

func regionSyncReferenceNotFound(
	layerID uint64,
	labelID uint64,
) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.region.sync_reference_not_found",
		fmt.Sprintf(
			"legend with landUse not found for layerId=%d labelId=%d",
			layerID,
			labelID,
		),
	).WithMetadata(map[string]string{
		"layer_id": strconv.FormatUint(layerID, 10),
		"label_id": strconv.FormatUint(labelID, 10),
	})
}

func regionInternal(cause error, code string) error {
	return fault.Wrap(
		cause,
		fault.KindInternal,
		code,
		"region operation failed",
	)
}
