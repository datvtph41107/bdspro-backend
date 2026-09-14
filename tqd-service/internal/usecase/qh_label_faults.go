package usecase

import (
	"fmt"
	"strconv"

	"common/fault"
)

func qhLabelSourceNotFound(id uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.qh_label.source_not_found",
		fmt.Sprintf("source label %d was not found", id),
	).WithMetadata(map[string]string{
		"source_label_id": strconv.FormatUint(id, 10),
	})
}

func qhLabelSourceLayerMismatch(id, layerID uint64) error {
	return fault.Validation(
		"tqd.qh_label.source_layer_mismatch",
		fmt.Sprintf("source label %d belongs to another layer", id),
		fault.FieldViolation{
			Field:       "source_label_ids",
			Description: "contains a label from another layer",
		},
	).WithMetadata(map[string]string{
		"source_label_id": strconv.FormatUint(id, 10),
		"layer_id":        strconv.FormatUint(layerID, 10),
	})
}

func qhLabelNameConflict(layerID uint64) error {
	return fault.New(
		fault.KindConflict,
		"tqd.qh_label.name_conflict",
		"label name already exists in layer",
	).WithMetadata(map[string]string{
		"layer_id": strconv.FormatUint(layerID, 10),
	})
}
