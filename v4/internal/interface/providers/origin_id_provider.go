package providers

import "context"

// OriginIDProvider resolves CRM origin_id for a profile from the local v4 DSN.
// Used when outbound CRM gRPC is unavailable (e.g. host dev without metadata signing).
type OriginIDProvider interface {
	GetOriginID(ctx context.Context, profileID uint64) (uint64, error)
}
