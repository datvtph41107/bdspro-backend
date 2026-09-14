package usecase

import (
	"fmt"
	"strconv"

	"common/fault"
)

func qhLegendPayloadRequired() error {
	return fault.Validation(
		"tqd.legend.payload_required",
		"payload is required",
		fault.FieldViolation{
			Field:       "payload",
			Description: "is required",
		},
	)
}

func qhLegendIDRequired() error {
	return fault.Validation(
		"tqd.legend.id_required",
		"id is required",
		fault.FieldViolation{
			Field:       "id",
			Description: "must be greater than zero",
		},
	)
}

func qhLegendLayerIDRequired() error {
	return fault.Validation(
		"tqd.legend.layer_id_required",
		"layerId is required",
		fault.FieldViolation{
			Field:       "layer_id",
			Description: "must be greater than zero",
		},
	)
}

func qhLegendLabelIDRequired() error {
	return fault.Validation(
		"tqd.legend.label_id_required",
		"labelId is required",
		fault.FieldViolation{
			Field:       "label_id",
			Description: "must be greater than zero",
		},
	)
}

func qhLegendLayerNotFound(layerID uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.legend.layer_not_found",
		fmt.Sprintf("layer %d not found", layerID),
	).WithMetadata(map[string]string{
		"layer_id": strconv.FormatUint(layerID, 10),
	})
}

func qhLegendLabelNotFound(labelID uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.legend.label_not_found",
		fmt.Sprintf("label %d not found", labelID),
	).WithMetadata(map[string]string{
		"label_id": strconv.FormatUint(labelID, 10),
	})
}

func qhLegendRecordNotFound(id uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.legend.record_not_found",
		fmt.Sprintf("record %d not found", id),
	).WithMetadata(map[string]string{
		"legend_id": strconv.FormatUint(id, 10),
	})
}

func qhLegendDuplicate(layerID, labelID uint64) error {
	return fault.New(
		fault.KindConflict,
		"tqd.legend.duplicate",
		fmt.Sprintf(
			"layer %d already has legend for label %d",
			layerID,
			labelID,
		),
	).WithMetadata(map[string]string{
		"layer_id": strconv.FormatUint(layerID, 10),
		"label_id": strconv.FormatUint(labelID, 10),
	})
}

func qhLegendInvalidLegendType(value string) error {
	return fault.Validation(
		"tqd.legend.legend_type_invalid",
		fmt.Sprintf("invalid legendType: %s", value),
		fault.FieldViolation{
			Field:       "legend_type",
			Description: "unsupported legend type",
		},
	)
}

func qhLegendInvalidGeometryType(value string) error {
	return fault.Validation(
		"tqd.legend.geometry_type_invalid",
		fmt.Sprintf("invalid geometryType: %s", value),
		fault.FieldViolation{
			Field:       "geometry_type",
			Description: "unsupported geometry type",
		},
	)
}
