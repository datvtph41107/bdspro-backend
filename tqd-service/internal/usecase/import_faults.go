package usecase

import (
	"fmt"
	"strconv"

	_errors "common/errors"
	"tqd/internal"
	qh_domain "tqd/internal/domain/qh"
)

func importLayerIDRequired() error {
	return _errors.ReturnError(
		service.ImportLayerIDRequired,
		_errors.WithViolations(_errors.FieldViolation{Field: "layer_id", Description: "must be greater than zero"}),
	)
}

func importErrorIDRequired() error {
	return _errors.ReturnError(
		service.ImportErrorIDRequired,
		_errors.WithViolations(_errors.FieldViolation{Field: "error_id", Description: "must be greater than zero"}),
	)
}

func importLayerNotFound(layerID uint64) error {
	return _errors.ReturnError(
		service.ImportLayerNotFound,
		_errors.WithPublicMessage(fmt.Sprintf("layer %d not found", layerID)),
		_errors.WithMetadata(map[string]string{"layer_id": strconv.FormatUint(layerID, 10)}),
	)
}

func importAlreadyInProgress(layerID uint64) error {
	return _errors.ReturnError(
		service.ImportInProgress,
		_errors.WithPublicMessage(fmt.Sprintf("layer %d already has an import in progress", layerID)),
		_errors.WithMetadata(map[string]string{"layer_id": strconv.FormatUint(layerID, 10)}),
	)
}

func importErrorNotFound(errorID uint64) error {
	return _errors.ReturnError(
		service.ImportErrorNotFound,
		_errors.WithPublicMessage(fmt.Sprintf("import error %d not found", errorID)),
		_errors.WithMetadata(map[string]string{"import_error_id": strconv.FormatUint(errorID, 10)}),
	)
}

func importRetryInProgress(errorID uint64) error {
	return _errors.ReturnError(
		service.ImportRetryInProgress,
		_errors.WithPublicMessage(fmt.Sprintf("import error %d retry is already in progress", errorID)),
		_errors.WithMetadata(map[string]string{"import_error_id": strconv.FormatUint(errorID, 10)}),
	)
}

func importRetryNotAllowed(errorID uint64, status qh_domain.ImportErrorStatus) error {
	return _errors.ReturnError(
		service.ImportRetryNotAllowed,
		_errors.WithPublicMessage(fmt.Sprintf("import error %d cannot be retried (status=%s)", errorID, status)),
		_errors.WithMetadata(map[string]string{
			"import_error_id": strconv.FormatUint(errorID, 10),
			"status":          string(status),
		}),
	)
}

func importRetryLockConflict(errorID uint64) error {
	return _errors.ReturnError(
		service.ImportRetryLockConflict,
		_errors.WithPublicMessage(fmt.Sprintf("import error %d could not be locked for retry", errorID)),
		_errors.WithMetadata(map[string]string{"import_error_id": strconv.FormatUint(errorID, 10)}),
	)
}
