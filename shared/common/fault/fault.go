// Package fault owns the backend error identity contract.
//
// Business and application code create typed faults here. Transport layers map
// those faults exactly once; callers must never classify failures by parsing an
// error message.
package fault

import (
	"errors"
	"strings"
)

type Kind string

const (
	KindValidation        Kind = "validation"
	KindUnauthenticated   Kind = "unauthenticated"
	KindPermissionDenied  Kind = "permission_denied"
	KindNotFound          Kind = "not_found"
	KindConflict          Kind = "conflict"
	KindPrecondition      Kind = "precondition"
	KindAborted           Kind = "aborted"
	KindResourceExhausted Kind = "resource_exhausted"
	KindUnavailable       Kind = "unavailable"
	KindInternal          Kind = "internal"
)

type FieldViolation struct {
	Field       string
	Description string
}

type Error struct {
	kind       Kind
	code       string
	message    string
	cause      error
	violations []FieldViolation
	metadata   map[string]string
}

func New(kind Kind, code, message string) *Error {
	return &Error{kind: kind, code: strings.TrimSpace(code), message: strings.TrimSpace(message)}
}

func Wrap(cause error, kind Kind, code, message string) *Error {
	if cause == nil {
		return New(kind, code, message)
	}
	return &Error{kind: kind, code: strings.TrimSpace(code), message: strings.TrimSpace(message), cause: cause}
}

func Validation(code, message string, violations ...FieldViolation) *Error {
	return New(KindValidation, code, message).WithViolations(violations...)
}

func (e *Error) WithViolations(violations ...FieldViolation) *Error {
	if e == nil {
		return e
	}
	e.violations = append(e.violations, violations...)
	return e
}

func (e *Error) WithMetadata(metadata map[string]string) *Error {
	if e == nil || len(metadata) == 0 {
		return e
	}
	if e.metadata == nil {
		e.metadata = make(map[string]string, len(metadata))
	}
	for key, value := range metadata {
		if key = strings.TrimSpace(key); key != "" {
			e.metadata[key] = value
		}
	}
	return e
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.message != "" {
		return e.message
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return string(e.kind)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *Error) Kind() Kind {
	if e == nil {
		return ""
	}
	return e.kind
}

func (e *Error) Code() string {
	if e == nil {
		return ""
	}
	return e.code
}

func (e *Error) PublicMessage() string {
	if e == nil {
		return ""
	}
	return e.message
}

func (e *Error) Violations() []FieldViolation {
	if e == nil || len(e.violations) == 0 {
		return nil
	}
	result := make([]FieldViolation, len(e.violations))
	copy(result, e.violations)
	return result
}

func (e *Error) Metadata() map[string]string {
	if e == nil || len(e.metadata) == 0 {
		return nil
	}
	result := make(map[string]string, len(e.metadata))
	for key, value := range e.metadata {
		result[key] = value
	}
	return result
}

func As(err error) (*Error, bool) {
	var target *Error
	if !errors.As(err, &target) {
		return nil, false
	}
	return target, true
}

func IsKind(err error, kind Kind) bool {
	target, ok := As(err)
	return ok && target.Kind() == kind
}
