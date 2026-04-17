// service/auth.go (refactored)
package service

import (
	"context"
	authpb "pb/types/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterAuthV2Service(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) {
	_ = authpb.RegisterAuthServiceHandler(ctx, mux, conn)
	_ = authpb.RegisterOAuthServiceHandler(ctx, mux, conn)
	_ = authpb.RegisterRoleServiceHandler(ctx, mux, conn)
	_ = authpb.RegisterPermissionServiceHandler(ctx, mux, conn)
	_ = authpb.RegisterRoleGroupServiceHandler(ctx, mux, conn)
	_ = authpb.RegisterAuthInternalServiceHandler(ctx, mux, conn)
	_ = authpb.RegisterUserInfoServiceHandler(ctx, mux, conn)
}
