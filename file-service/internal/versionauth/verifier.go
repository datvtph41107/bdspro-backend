package versionauth

import "context"

// Verifier is the narrow capability VersionService consumes from Hub.
// File-service owns this contract; transport and Hub client lifecycle belong to
// the composition adapter, not to the business service constructor.
type Verifier interface {
	Verify(ctx context.Context, apiKey string) error
}
