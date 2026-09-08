package authclient

import (
	"fmt"
	"strings"
	"time"

	qhprorpc "common/rpc"
	"common/rpcenv"
	"pb/clients"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewRPCClientAt tạo đúng một client Auth nội bộ cho process Payment. Auth
// vẫn là owner permission; Payment chỉ gửi yêu cầu kiểm tra theo actor context.
func NewRPCClientAt(target string) (*clients.AuthGrpcClient, func(), error) {
	if strings.TrimSpace(target) == "" {
		return nil, nil, fmt.Errorf("auth gRPC address is required")
	}
	return qhprorpc.NewBoundClient(qhprorpc.ClientConfig{
		Target: target, BackoffMaxDelay: 5 * time.Second,
		Credentials: insecure.NewCredentials(), Transport: rpcenv.LoadTransportConfig(),
	}, func(conn grpc.ClientConnInterface) *clients.AuthGrpcClient {
		return clients.BindAuthGrpcClient(conn, nil)
	})
}
