package usecase

import (
	"fmt"
	"strconv"

	_errors "common/errors"
	"tqd/internal"
)

func qhLabelSourceNotFound(id uint64) error {
	return _errors.ReturnError(service.LabelSourceNotFound, _errors.WithPublicMessage(fmt.Sprintf("source label %d was not found", id)), _errors.WithMetadata(map[string]string{"source_label_id": strconv.FormatUint(id, 10)}))
}
func qhLabelSourceLayerMismatch(id, layerID uint64) error {
	return _errors.ReturnError(
		service.LabelSourceLayerMismatch,
		_errors.WithPublicMessage(fmt.Sprintf("source label %d belongs to another layer", id)),
		_errors.WithViolations(_errors.FieldViolation{Field: "source_label_ids", Description: "contains a label from another layer"}),
		_errors.WithMetadata(map[string]string{"source_label_id": strconv.FormatUint(id, 10), "layer_id": strconv.FormatUint(layerID, 10)}),
	)
}
func qhLabelNameConflict(layerID uint64) error {
	return _errors.ReturnError(service.LabelNameConflict, _errors.WithMetadata(map[string]string{"layer_id": strconv.FormatUint(layerID, 10)}))
}
