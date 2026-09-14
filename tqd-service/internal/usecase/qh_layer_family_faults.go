package usecase

import (
	"fmt"
	"strconv"

	"common/fault"
)

var (
	ErrQHLayerFamilyPayloadRequired = fault.Validation(
		"tqd.qh_layer_family.payload_required",
		"payload is required",
	)
	ErrQHLayerFamilyNameRequired = fault.Validation(
		"tqd.qh_layer_family.name_required",
		"name is required",
		fault.FieldViolation{Field: "name", Description: "is required"},
	)
	ErrQHLayerFamilyIDRequired = fault.Validation(
		"tqd.qh_layer_family.id_required",
		"id is required",
		fault.FieldViolation{Field: "id", Description: "must be positive"},
	)
)

func qhLayerFamilyNotFound(id uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.qh_layer_family.not_found",
		fmt.Sprintf("layer family %d was not found", id),
	).WithMetadata(map[string]string{
		"id": strconv.FormatUint(id, 10),
	})
}
