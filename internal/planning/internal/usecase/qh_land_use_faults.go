package usecase

import (
	"fmt"
	"strconv"

	_errors "common/errors"
	"tqd/internal"
)

func qhLandUsePayloadRequired() error {
	return _errors.ReturnError(service.LandUsePayloadRequired, _errors.WithViolations(_errors.FieldViolation{Field: "payload", Description: "is required"}))
}
func qhLandUseIDRequired() error {
	return _errors.ReturnError(service.LandUseIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "id", Description: "must be greater than zero"}))
}
func qhLandUseLayerIDRequired() error {
	return _errors.ReturnError(service.LandUseLayerIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "layer_id", Description: "must be greater than zero"}))
}
func qhLandUseIDValueRequired() error {
	return _errors.ReturnError(service.LandUseValueIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "land_use_id", Description: "must be greater than zero"}))
}
func qhLandUseGroupIDRequired() error {
	return _errors.ReturnError(service.LandUseGroupIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "land_use_group_id", Description: "must be greater than zero"}))
}
func qhLandUseNameRequired() error {
	return _errors.ReturnError(service.LandUseNameRequired, _errors.WithViolations(_errors.FieldViolation{Field: "name", Description: "is required"}))
}
func qhLandUseLayerNotFound(layerID uint64) error {
	return _errors.ReturnError(service.LandUseLayerNotFound, _errors.WithPublicMessage(fmt.Sprintf("layer %d not found", layerID)), _errors.WithMetadata(map[string]string{"layer_id": strconv.FormatUint(layerID, 10)}))
}
func qhLandUseNotFound(landUseID uint64) error {
	return _errors.ReturnError(service.LandUseNotFound, _errors.WithPublicMessage(fmt.Sprintf("land use %d not found", landUseID)), _errors.WithMetadata(map[string]string{"land_use_id": strconv.FormatUint(landUseID, 10)}))
}
func qhLandUseRecordNotFound(id uint64) error {
	return _errors.ReturnError(service.LandUseRecordNotFound, _errors.WithPublicMessage(fmt.Sprintf("record %d not found", id)), _errors.WithMetadata(map[string]string{"land_use_id": strconv.FormatUint(id, 10)}))
}
func qhLandUseDuplicate(layerID, landUseID uint64) error {
	return _errors.ReturnError(
		service.LandUseDuplicate,
		_errors.WithPublicMessage(fmt.Sprintf("layer %d already has land use %d", layerID, landUseID)),
		_errors.WithMetadata(map[string]string{"layer_id": strconv.FormatUint(layerID, 10), "land_use_id": strconv.FormatUint(landUseID, 10)}),
	)
}
