// Package identitygrpc is the compatibility facade for historical Actor/Caller
// metadata helpers. New transport code should use common/rpc.
package identitygrpc

import (
	"common/rpc"
	"context"

	"google.golang.org/grpc/metadata"
)

// Deprecated: use rpc.AppendIdentityFromContext.
func AppendFromContext(ctx context.Context, md metadata.MD) metadata.MD {
	return rpc.AppendIdentityFromContext(ctx, md)
}

// Deprecated: use rpc.RestoreIdentityFromMetadata with an already verified
// transport context.
func RestoreFromMetadata(ctx context.Context, md metadata.MD) (context.Context, error) {
	return rpc.RestoreIdentityFromMetadata(ctx, md, true)
}
