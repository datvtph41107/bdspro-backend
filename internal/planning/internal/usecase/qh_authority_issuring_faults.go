package usecase

import (
	"fmt"
	"strconv"

	_errors "common/errors"
	"tqd/internal"
)

var (
	ErrQHAuthorityIssuringPayloadRequired = _errors.ReturnError(service.AuthorityPayloadRequired)
	ErrQHAuthorityIssuringNameRequired    = _errors.ReturnError(service.AuthorityNameRequired, _errors.WithViolations(_errors.FieldViolation{Field: "name", Description: "is required"}))
	ErrQHAuthorityIssuringCodeRequired    = _errors.ReturnError(service.AuthorityCodeRequired, _errors.WithViolations(_errors.FieldViolation{Field: "code", Description: "is required"}))
	ErrQHAuthorityIssuringIDRequired      = _errors.ReturnError(service.AuthorityIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "id", Description: "must be positive"}))
)

func qhAuthorityIssuringNotFound(id uint64) error {
	return _errors.ReturnError(service.AuthorityNotFound, _errors.WithPublicMessage(fmt.Sprintf("authority issuring %d was not found", id)), _errors.WithMetadata(map[string]string{"id": strconv.FormatUint(id, 10)}))
}
func qhAuthorityIssuringCodeConflict(code string) error {
	return _errors.ReturnError(service.AuthorityCodeConflict, _errors.WithPublicMessage(fmt.Sprintf("authority issuring code %q already exists", code)), _errors.WithMetadata(map[string]string{"code": code}))
}
