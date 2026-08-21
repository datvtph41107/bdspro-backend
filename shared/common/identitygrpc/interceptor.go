package identitygrpc

import (
	"common/rpc"
	"common/transporttrust"

	"google.golang.org/grpc"
)

// Deprecated: use rpc.ServerTransport. This strict facade preserves the old
// identitygrpc contract for callers that have not moved to canonical server
// composition yet.
func UnaryServerInterceptor(config transporttrust.Config) grpc.UnaryServerInterceptor {
	return rpc.StrictIdentityUnaryServerInterceptor(config)
}
