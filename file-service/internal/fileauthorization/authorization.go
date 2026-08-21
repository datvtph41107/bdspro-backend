package fileauthorization

import (
	"context"
	"errors"
)

// Capability names the exact current-authority decision consumed by File-service.
type Capability string

const (
	// PrivateReadAny allows an operator to read a private file without a per-file
	// file_access row. It is intentionally stronger than ordinary signed access.
	PrivateReadAny Capability = "FILE_PRIVATE_READ_ANY"
)

var (
	ErrDenied      = errors.New("file permission denied")
	ErrUnavailable = errors.New("file permission authority unavailable")
)

// Authorizer is owned by File-service as the consumer contract. Implementations
// may use Auth gRPC, but callers do not know transport, role names, or numeric IDs.
type Authorizer interface {
	Require(ctx context.Context, capability Capability) error
}

// UnavailableAuthorizer is the safe fallback when no current permission authority
// was composed. Missing authority never turns into an implicit admin bypass.
type UnavailableAuthorizer struct{}

func (UnavailableAuthorizer) Require(context.Context, Capability) error {
	return ErrUnavailable
}
