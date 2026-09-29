package _errors

import (
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/grpc/codes"
)

type Code int32
type Key string

type Spec struct {
	code              Code
	key               Key
	message           string
	rpc               codes.Code
	legacyHTTP200     bool
	legacyCode        int32
	legacyProblemCode string
}

type SpecOption func(*Spec)

// LegacyHTTP200 preserves the historical Gateway/direct-gRPC envelope.
// It is migration-only and must be paired with LegacyCode when the historical
// public code differs from the canonical application Code.
func LegacyHTTP200() SpecOption {
	return func(spec *Spec) {
		spec.legacyHTTP200 = true
	}
}

// LegacyCode preserves the currently consumed numeric body/detail code while
// Spec.Code remains the globally stable canonical application identity.
func LegacyCode(code int32) SpecOption {
	return func(spec *Spec) {
		spec.legacyCode = code
	}
}

// LegacyProblemCode preserves an existing public string code while the
// canonical numeric Code + semantic Key become the internal source of truth.
func LegacyProblemCode(code string) SpecOption {
	return func(spec *Spec) {
		spec.legacyProblemCode = strings.TrimSpace(code)
	}
}

var keyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

func MustSpec(code Code, key Key, message string, rpc codes.Code, options ...SpecOption) Spec {
	spec := Spec{
		code:    code,
		key:     Key(strings.TrimSpace(string(key))),
		message: strings.TrimSpace(message),
		rpc:     rpc,
	}
	for _, option := range options {
		if option != nil {
			option(&spec)
		}
	}
	if err := spec.Validate(); err != nil {
		panic(err)
	}
	return spec
}

func (s Spec) Validate() error {
	if s.code <= 0 {
		return fmt.Errorf("error spec code must be positive")
	}
	if !keyPattern.MatchString(string(s.key)) {
		return fmt.Errorf("error spec key %q must be UPPER_SNAKE_CASE", s.key)
	}
	if s.message == "" {
		return fmt.Errorf("error spec %s message must not be empty", s.key)
	}
	if s.rpc == codes.OK {
		return fmt.Errorf("error spec %s cannot use gRPC OK", s.key)
	}
	return nil
}

func (s Spec) Code() Code                { return s.code }
func (s Spec) Key() Key                  { return s.key }
func (s Spec) Message() string           { return s.message }
func (s Spec) RPCCode() codes.Code       { return s.rpc }
func (s Spec) LegacyHTTP200() bool       { return s.legacyHTTP200 }
func (s Spec) LegacyCode() int32         { return s.legacyCode }
func (s Spec) LegacyProblemCode() string { return s.legacyProblemCode }
