package rpc

import (
	"fmt"
	"strings"
	"time"

	qhprorpc "common/rpc"
	"common/rpcenv"
	"pb/clients"
	assistantpb "pb/types/assistant"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const assistantGrpcMaxMsgBytes = 32 << 20

func openBound[T any](target string, bind func(grpc.ClientConnInterface) T, options ...grpc.CallOption) (T, func(), error) {
	if strings.TrimSpace(target) == "" {
		var zero T
		return zero, nil, fmt.Errorf("gRPC target is required")
	}
	return qhprorpc.NewBoundClient(qhprorpc.ClientConfig{
		Target:             strings.TrimSpace(target),
		BackoffMaxDelay:    5 * time.Second,
		Credentials:        insecure.NewCredentials(),
		Transport:          rpcenv.LoadTransportConfig(),
		DefaultCallOptions: options,
	}, bind)
}

func OpenUserRPCClient(target string) (*clients.UserGrpcClient, func(), error) {
	return openBound(target, clients.BindUserGrpcClient)
}

func OpenAuthRPCClient(target string) (*clients.AuthGrpcClient, func(), error) {
	return openBound(target, func(conn grpc.ClientConnInterface) *clients.AuthGrpcClient {
		return clients.BindAuthGrpcClient(conn, nil)
	})
}

func OpenAssistantRPCClient(target string) (assistantpb.AssistantServiceClient, func(), error) {
	return openBound(
		target,
		func(conn grpc.ClientConnInterface) assistantpb.AssistantServiceClient {
			return assistantpb.NewAssistantServiceClient(conn)
		},
		grpc.MaxCallRecvMsgSize(assistantGrpcMaxMsgBytes),
		grpc.MaxCallSendMsgSize(assistantGrpcMaxMsgBytes),
	)
}
