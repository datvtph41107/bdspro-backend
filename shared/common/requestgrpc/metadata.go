// Package requestgrpc is the compatibility facade for historical request
// metadata helpers. New transport code should use common/rpc.
package requestgrpc

import (
	"common/rpc"
	"context"

	"google.golang.org/grpc/metadata"
)

// Deprecated: use rpc.AppendRequestFromContext.
func AppendFromContext(ctx context.Context, md metadata.MD) metadata.MD {
	return rpc.AppendRequestFromContext(ctx, md)
}

// Deprecated: use rpc.RestoreRequestFromMetadata.
func RestoreFromMetadata(ctx context.Context, md metadata.MD) (context.Context, error) {
	return rpc.RestoreRequestFromMetadata(ctx, md)
}
