package request

import (
	"errors"
	"strings"
)

const (
	OperationIDHeader      = "X-Operation-ID"
	OperationIDMetadataKey = "x-operation-id"
	MaxOperationIDLength   = 128
)

type OperationIDSource string

const (
	OperationIDSourceCaller = OperationIDSource("caller")

	OperationIDSourceGeneratedFallback = OperationIDSource("generated_fallback")
)

type OperationIDReason string

const (
	OperationIDReasonValid = OperationIDReason("valid")

	OperationIDReasonMissing = OperationIDReason("missing")
)

type OperationIDDecision struct {
	ID string

	Source OperationIDSource

	Reason OperationIDReason
}

var (
	ErrInvalidOperationID = errors.New(
		"invalid operation ID",
	)

	ErrMultipleOperationIDs = errors.New(
		"multiple operation IDs",
	)

	ErrOperationIDGeneration = errors.New(
		"operation ID generation failed",
	)
)

/**
 * NormalizeOperationID normalizes one externally supplied operation ID.
 */
func NormalizeOperationID(
	value string,
) string {
	return strings.TrimSpace(value)
}

/**
 * IsValidOperationID validates one normalized operation ID.
 */
func IsValidOperationID(
	value string,
) bool {
	if value == "" ||
		len(value) > MaxOperationIDLength {

		return false
	}

	for i := 0; i < len(value); i++ {
		c := value[i]

		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '-' ||
			c == '_' ||
			c == '.' ||
			c == ':' {

			continue
		}

		return false
	}

	return true
}

/**
 * AcceptOperationID accepts exactly one valid propagated operation ID.
 */
func AcceptOperationID(
	values []string,
) (
	string,
	bool,
) {
	if len(values) != 1 {
		return "", false
	}

	operationID :=
		NormalizeOperationID(
			values[0],
		)

	if !IsValidOperationID(
		operationID,
	) {
		return "", false
	}

	return operationID, true
}

/**
 * ResolveOperationID resolves caller input at the logical-operation origin boundary.
 */
func ResolveOperationID(
	values []string,
) (
	OperationIDDecision,
	error,
) {
	return resolveOperationID(
		values,
		NewOperationID,
	)
}

/**
 * resolveOperationID classifies caller input and authorizes fallback
 * generation only when the operation ID is missing.
 */
func resolveOperationID(
	values []string,
	generate operationIDGenerator,
) (
	OperationIDDecision,
	error,
) {
	if operationID, ok :=
		AcceptOperationID(
			values,
		); ok {

		return OperationIDDecision{
			ID:     operationID,
			Source: OperationIDSourceCaller,
			Reason: OperationIDReasonValid,
		}, nil
	}

	switch len(values) {
	case 0:
		operationID, err :=
			generate()

		if err != nil {
			return OperationIDDecision{},
				err
		}

		if !IsValidOperationID(
			operationID,
		) {
			return OperationIDDecision{},
				ErrOperationIDGeneration
		}

		return OperationIDDecision{
			ID: operationID,

			Source: OperationIDSourceGeneratedFallback,

			Reason: OperationIDReasonMissing,
		}, nil

	case 1:
		return OperationIDDecision{},
			ErrInvalidOperationID

	default:
		return OperationIDDecision{},
			ErrMultipleOperationIDs
	}
}
