package versionauth

import (
	"context"
	"errors"
)

var (
	// ErrInvalid means the caller did not provide acceptable API-key evidence.
	ErrInvalid = errors.New("invalid version API key")

	// ErrUnavailable means the verifier could not produce an authorization
	// decision because its authority/dependency was unavailable.
	ErrUnavailable = errors.New("version API key verifier unavailable")
)

// Verifier is the narrow API-key authority consumed by File version upload.
// File-service owns this contract; transport and Hub client lifecycle belong to
// the composition adapter, not to the business service constructor.
type Verifier interface {
	Verify(ctx context.Context, apiKey string) error
}
