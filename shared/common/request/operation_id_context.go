package request

import (
	"context"
	"errors"
)

type operationIDContextKey struct{}

var ErrOperationIDContextConflict = errors.New(
	"operation ID context conflict",
)

/**
 * BindOperationID binds one canonical logical-operation identity.
 *
 * The first valid binding wins. Rebinding the same identity is
 * idempotent; binding a different identity returns a conflict.
 */
func BindOperationID(
	ctx context.Context,
	operationID string,
) (
	context.Context,
	error,
) {
	if ctx == nil {
		panic(
			"request: nil context",
		)
	}

	if !IsValidOperationID(
		operationID,
	) {
		return ctx, ErrInvalidOperationID
	}

	if existing, ok :=
		OperationIDFromContext(
			ctx,
		); ok {

		if existing == operationID {
			return ctx, nil
		}

		return ctx, ErrOperationIDContextConflict
	}

	return context.WithValue(
			ctx,
			operationIDContextKey{},
			operationID,
		),
		nil
}

/**
 * OperationIDFromContext reads the canonical logical-operation identity.
 */
func OperationIDFromContext(
	ctx context.Context,
) (
	string,
	bool,
) {
	if ctx == nil {
		return "", false
	}

	operationID, ok := ctx.Value(
		operationIDContextKey{},
	).(string)

	if !ok ||
		!IsValidOperationID(
			operationID,
		) {

		return "", false
	}

	return operationID, true
}
