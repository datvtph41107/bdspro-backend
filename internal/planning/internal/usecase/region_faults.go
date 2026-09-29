package usecase

import (
	"fmt"
	"strconv"

	_errors "common/errors"
	"tqd/internal"
)

func regionLayerIDRequired() error {
	return _errors.ReturnError(service.RegionLayerIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "layer_id", Description: "must be greater than zero"}))
}
func regionLabelIDRequired() error {
	return _errors.ReturnError(service.RegionLabelIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "label_id", Description: "must be greater than zero"}))
}
func regionIDRequired() error {
	return _errors.ReturnError(service.RegionRecordIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "id", Description: "must be greater than zero"}))
}
func regionNameRequired() error {
	return _errors.ReturnError(service.RegionNameRequired, _errors.WithViolations(_errors.FieldViolation{Field: "name", Description: "is required"}))
}
func regionUpdateRequired() error {
	return _errors.ReturnError(service.RegionUpdateRequired)
}
func regionNameEmpty() error {
	return _errors.ReturnError(service.RegionNameEmpty, _errors.WithViolations(_errors.FieldViolation{Field: "name", Description: "cannot be empty"}))
}
func regionStatusInvalid(value uint32) error {
	return _errors.ReturnError(
		service.RegionStatusInvalid,
		_errors.WithPublicMessage(fmt.Sprintf("invalid status: %d", value)),
		_errors.WithViolations(_errors.FieldViolation{Field: "status", Description: "unsupported region status"}),
	)
}
func regionGeometryInvalid(cause error) error {
	if cause == nil {
		return _errors.ReturnError(service.RegionGeometryInvalid)
	}
	return _errors.ReturnError(
		service.RegionGeometryInvalid,
		_errors.WithCause(cause),
		_errors.WithPublicMessage(cause.Error()),
	)
}
func regionLayerNotFound(layerID uint64) error {
	return _errors.ReturnError(
		service.RegionLayerNotFound,
		_errors.WithCause(ErrLayerNotFound),
		_errors.WithPublicMessage(fmt.Sprintf("layer %d not found", layerID)),
		_errors.WithMetadata(map[string]string{"layer_id": strconv.FormatUint(layerID, 10)}),
	)
}
func regionNotFound(regionID uint64) error {
	return _errors.ReturnError(
		service.RegionRecordNotFound,
		_errors.WithPublicMessage(fmt.Sprintf("region %d not found", regionID)),
		_errors.WithMetadata(map[string]string{"region_id": strconv.FormatUint(regionID, 10)}),
	)
}
func regionLabelNotFound(labelID uint64) error {
	return _errors.ReturnError(
		service.RegionLabelNotFound,
		_errors.WithPublicMessage(fmt.Sprintf("label %d not found", labelID)),
		_errors.WithMetadata(map[string]string{"label_id": strconv.FormatUint(labelID, 10)}),
	)
}
func regionSyncReferenceNotFound(layerID, labelID uint64) error {
	return _errors.ReturnError(
		service.RegionSyncReferenceNotFound,
		_errors.WithPublicMessage(fmt.Sprintf("legend with landUse not found for layerId=%d labelId=%d", layerID, labelID)),
		_errors.WithMetadata(map[string]string{
			"layer_id": strconv.FormatUint(layerID, 10),
			"label_id": strconv.FormatUint(labelID, 10),
		}),
	)
}
func regionInternal(cause error, code string) error {
	if cause == nil {
		return fmt.Errorf("region operation failed (%s)", code)
	}
	return fmt.Errorf("region operation failed (%s): %w", code, cause)
}
