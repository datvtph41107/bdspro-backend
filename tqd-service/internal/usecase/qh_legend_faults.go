package usecase

import (
	"fmt"
	"strconv"

	_errors "common/errors"
	"tqd/internal"
)

func qhLegendPayloadRequired() error {
	return _errors.ReturnError(service.LegendPayloadRequired, _errors.WithViolations(_errors.FieldViolation{Field: "payload", Description: "is required"}))
}
func qhLegendIDRequired() error {
	return _errors.ReturnError(service.LegendIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "id", Description: "must be greater than zero"}))
}
func qhLegendLayerIDRequired() error {
	return _errors.ReturnError(service.LegendLayerIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "layer_id", Description: "must be greater than zero"}))
}
func qhLegendLabelIDRequired() error {
	return _errors.ReturnError(service.LegendLabelIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "label_id", Description: "must be greater than zero"}))
}
func qhLegendLayerNotFound(layerID uint64) error {
	return _errors.ReturnError(service.LegendLayerNotFound, _errors.WithPublicMessage(fmt.Sprintf("layer %d not found", layerID)), _errors.WithMetadata(map[string]string{"layer_id": strconv.FormatUint(layerID, 10)}))
}
func qhLegendLabelNotFound(labelID uint64) error {
	return _errors.ReturnError(service.LegendLabelNotFound, _errors.WithPublicMessage(fmt.Sprintf("label %d not found", labelID)), _errors.WithMetadata(map[string]string{"label_id": strconv.FormatUint(labelID, 10)}))
}
func qhLegendRecordNotFound(id uint64) error {
	return _errors.ReturnError(service.LegendRecordNotFound, _errors.WithPublicMessage(fmt.Sprintf("record %d not found", id)), _errors.WithMetadata(map[string]string{"legend_id": strconv.FormatUint(id, 10)}))
}
func qhLegendDuplicate(layerID, labelID uint64) error {
	return _errors.ReturnError(
		service.LegendDuplicate,
		_errors.WithPublicMessage(fmt.Sprintf("layer %d already has legend for label %d", layerID, labelID)),
		_errors.WithMetadata(map[string]string{"layer_id": strconv.FormatUint(layerID, 10), "label_id": strconv.FormatUint(labelID, 10)}),
	)
}
func qhLegendInvalidLegendType(value string) error {
	return _errors.ReturnError(service.LegendTypeInvalid, _errors.WithPublicMessage(fmt.Sprintf("invalid legendType: %s", value)), _errors.WithViolations(_errors.FieldViolation{Field: "legend_type", Description: "unsupported legend type"}))
}
func qhLegendInvalidGeometryType(value string) error {
	return _errors.ReturnError(service.LegendGeometryTypeInvalid, _errors.WithPublicMessage(fmt.Sprintf("invalid geometryType: %s", value)), _errors.WithViolations(_errors.FieldViolation{Field: "geometry_type", Description: "unsupported geometry type"}))
}
