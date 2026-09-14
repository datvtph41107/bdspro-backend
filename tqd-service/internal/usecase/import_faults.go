package usecase

import (
	"fmt"
	"strconv"

	"common/fault"

	qh_domain "tqd/internal/domain/qh"
)

func importLayerIDRequired() error {
	return fault.Validation(
		"tqd.import.layer_id_required",
		"layerId is required",
		fault.FieldViolation{
			Field:       "layer_id",
			Description: "must be greater than zero",
		},
	)
}

func importErrorIDRequired() error {
	return fault.Validation(
		"tqd.import.error_id_required",
		"errorId is required",
		fault.FieldViolation{
			Field:       "error_id",
			Description: "must be greater than zero",
		},
	)
}

func importLayerNotFound(layerID uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.import.layer_not_found",
		fmt.Sprintf("layer %d not found", layerID),
	).WithMetadata(map[string]string{
		"layer_id": strconv.FormatUint(layerID, 10),
	})
}

func importAlreadyInProgress(layerID uint64) error {
	return fault.New(
		fault.KindPrecondition,
		"tqd.import.in_progress",
		fmt.Sprintf("layer %d already has an import in progress", layerID),
	).WithMetadata(map[string]string{
		"layer_id": strconv.FormatUint(layerID, 10),
	})
}

func importErrorNotFound(errorID uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.import.error_not_found",
		fmt.Sprintf("import error %d not found", errorID),
	).WithMetadata(map[string]string{
		"import_error_id": strconv.FormatUint(errorID, 10),
	})
}

func importRetryInProgress(errorID uint64) error {
	return fault.New(
		fault.KindPrecondition,
		"tqd.import.retry_in_progress",
		fmt.Sprintf("import error %d retry is already in progress", errorID),
	).WithMetadata(map[string]string{
		"import_error_id": strconv.FormatUint(errorID, 10),
	})
}

func importRetryNotAllowed(errorID uint64, status qh_domain.ImportErrorStatus) error {
	return fault.New(
		fault.KindPrecondition,
		"tqd.import.retry_not_allowed",
		fmt.Sprintf("import error %d cannot be retried (status=%s)", errorID, status),
	).WithMetadata(map[string]string{
		"import_error_id": strconv.FormatUint(errorID, 10),
		"status":          string(status),
	})
}

func importRetryLockConflict(errorID uint64) error {
	return fault.New(
		fault.KindAborted,
		"tqd.import.retry_lock_conflict",
		fmt.Sprintf("import error %d could not be locked for retry", errorID),
	).WithMetadata(map[string]string{
		"import_error_id": strconv.FormatUint(errorID, 10),
	})
}
