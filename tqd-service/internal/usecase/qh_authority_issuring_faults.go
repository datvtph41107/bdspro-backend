package usecase

import (
	"fmt"
	"strconv"

	"common/fault"
)

var (
	ErrQHAuthorityIssuringPayloadRequired = fault.Validation(
		"tqd.qh_authority_issuring.payload_required",
		"payload is required",
	)
	ErrQHAuthorityIssuringNameRequired = fault.Validation(
		"tqd.qh_authority_issuring.name_required",
		"name is required",
		fault.FieldViolation{Field: "name", Description: "is required"},
	)
	ErrQHAuthorityIssuringCodeRequired = fault.Validation(
		"tqd.qh_authority_issuring.code_required",
		"code is required",
		fault.FieldViolation{Field: "code", Description: "is required"},
	)
	ErrQHAuthorityIssuringIDRequired = fault.Validation(
		"tqd.qh_authority_issuring.id_required",
		"id is required",
		fault.FieldViolation{Field: "id", Description: "must be positive"},
	)
)

func qhAuthorityIssuringNotFound(id uint64) error {
	return fault.New(
		fault.KindNotFound,
		"tqd.qh_authority_issuring.not_found",
		fmt.Sprintf("authority issuring %d was not found", id),
	).WithMetadata(map[string]string{
		"id": strconv.FormatUint(id, 10),
	})
}

func qhAuthorityIssuringCodeConflict(code string) error {
	return fault.New(
		fault.KindConflict,
		"tqd.qh_authority_issuring.code_conflict",
		fmt.Sprintf("authority issuring code %q already exists", code),
	).WithMetadata(map[string]string{
		"code": code,
	})
}
