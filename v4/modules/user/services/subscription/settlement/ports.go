package settlement

import "context"

// Store owns the atomic business checkpoint: effect arbitration, catalog
// verification, subscription mutation/event and durable effect evidence commit.
type Store interface {
	ApplySettlement(ctx context.Context, input Acceptance) (Result, error)
}
