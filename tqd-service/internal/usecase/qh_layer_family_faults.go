package usecase

import (
	"fmt"
	"strconv"

	_errors "common/errors"
	"tqd/internal"
)

var (
	ErrQHLayerFamilyPayloadRequired = _errors.ReturnError(service.LayerFamilyPayloadRequired)
	ErrQHLayerFamilyNameRequired    = _errors.ReturnError(service.LayerFamilyNameRequired, _errors.WithViolations(_errors.FieldViolation{Field: "name", Description: "is required"}))
	ErrQHLayerFamilyIDRequired      = _errors.ReturnError(service.LayerFamilyIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "id", Description: "must be positive"}))
)

func qhLayerFamilyNotFound(id uint64) error {
	return _errors.ReturnError(service.LayerFamilyNotFound, _errors.WithPublicMessage(fmt.Sprintf("layer family %d was not found", id)), _errors.WithMetadata(map[string]string{"id": strconv.FormatUint(id, 10)}))
}
