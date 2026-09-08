package grpcgateway

import (
	"context"
	authpb "pb/types/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterAuthService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		authpb.RegisterAuthServiceHandler,
		authpb.RegisterOAuthServiceHandler,
		authpb.RegisterRoleServiceHandler,
		authpb.RegisterPermissionServiceHandler,
		authpb.RegisterRoleGroupServiceHandler,
		authpb.RegisterUserInfoServiceHandler,
	)
}
