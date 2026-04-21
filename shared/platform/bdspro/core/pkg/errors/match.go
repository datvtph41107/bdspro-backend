package apperr

import (
	stdErrors "errors"
)

// As extracts *AppError từ error chain.
// Dùng thay cho type assertion trực tiếp — handle wrapped errors.
//
// Ví dụ:
//
//	ae, ok := apperr.As(err)
//	if !ok { /* not an AppError */ }

func As(err error) (*AppError, bool) {
	var ae *AppError
	if stdErrors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

func IsKind(err error, kind Kind) bool {
	ae, ok := As(err)
	return ok && ae.kind == kind
}

func IsCode(err error, code Code) bool {
	ae, ok := As(err)
	return ok && ae.code == code
}

func IsRetryable(err error) bool {
	ae, ok := As(err)
	return ok && ae.retryable
}

func IsClientFault(err error) bool {
	ae, ok := As(err)
	return ok && ae.kind.IsClientFault()
}

func IsSystemFault(err error) bool {
	ae, ok := As(err)
	return ok && ae.kind.IsSystemFault()
}
