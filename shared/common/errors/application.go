package _errors

import (
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
)

type FieldViolation struct {
	Field       string
	Description string
}

type Error struct {
	spec       Spec
	message    string
	cause      error
	violations []FieldViolation
	metadata   map[string]string
	legacyCode *int32
}

type Option func(*Error)

func WithCause(cause error) Option {
	return func(target *Error) {
		target.cause = cause
	}
}

// WithPublicMessage is intentionally an occurrence-level presentation override.
// Prefer the canonical Spec message. Keep this only for values that are truly
// dynamic and cannot be represented as typed details yet.
func WithPublicMessage(message string) Option {
	return func(target *Error) {
		target.message = strings.TrimSpace(message)
	}
}

func WithViolations(violations ...FieldViolation) Option {
	return func(target *Error) {
		target.violations = append(target.violations, violations...)
	}
}

// WithMetadata carries bounded structured occurrence context into ErrorInfo.
// It is not part of Spec identity and must never contain secrets or root-cause
// text. Prefer typed details where a stable consumer exists.
func WithLegacyCode(code int32) Option {
	return func(target *Error) {
		target.legacyCode = &code
	}
}

func WithMetadata(metadata map[string]string) Option {
	return func(target *Error) {
		if len(metadata) == 0 {
			return
		}
		if target.metadata == nil {
			target.metadata = make(map[string]string, len(metadata))
		}
		for key, value := range metadata {
			key = strings.TrimSpace(key)
			if key != "" {
				target.metadata[key] = value
			}
		}
	}
}

func newError(spec Spec, options ...Option) *Error {
	target := &Error{spec: spec}
	for _, option := range options {
		if option != nil {
			option(target)
		}
	}
	return target
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.PublicMessage()
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *Error) Spec() Spec {
	if e == nil {
		return Spec{}
	}
	return e.spec
}

func (e *Error) Code() Code {
	if e == nil {
		return 0
	}
	return e.spec.Code()
}

func (e *Error) Key() Key {
	if e == nil {
		return ""
	}
	return e.spec.Key()
}

func (e *Error) PublicMessage() string {
	if e == nil {
		return ""
	}
	if e.message != "" {
		return e.message
	}
	return e.spec.Message()
}

func (e *Error) RPCCode() codes.Code {
	if e == nil {
		return codes.Internal
	}
	return e.spec.RPCCode()
}

func (e *Error) LegacyCode() (int32, bool) {
	if e == nil || e.legacyCode == nil {
		return 0, false
	}
	return *e.legacyCode, true
}

func (e *Error) Violations() []FieldViolation {
	if e == nil || len(e.violations) == 0 {
		return nil
	}
	out := make([]FieldViolation, len(e.violations))
	copy(out, e.violations)
	return out
}

func (e *Error) Metadata() map[string]string {
	if e == nil || len(e.metadata) == 0 {
		return nil
	}
	out := make(map[string]string, len(e.metadata))
	for key, value := range e.metadata {
		out[key] = value
	}
	return out
}

func As(err error) (*Error, bool) {
	var target *Error
	if !errors.As(err, &target) {
		return nil, false
	}
	return target, true
}
