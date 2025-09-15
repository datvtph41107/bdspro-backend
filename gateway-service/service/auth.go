package service

import (
	"context"
	"fmt"
	"log"
	authpb "pb/types/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterAuthService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.auth")
	url := "localhost"
	if env == "release" {
		url = "auth"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register auth service: %v", err)
	}
	err = authpb.RegisterOAuthServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register oauth service: %v", err)
	}
	err = authpb.RegisterRoleServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register role service: %v", err)
	}
	err = authpb.RegisterPermissionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register permission service: %v", err)
	}
	err = authpb.RegisterRoleGroupServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register role group service: %v", err)
	}
	err = authpb.RegisterAuthInternalServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register auth internal service: %v", err)
	}
	err = authpb.RegisterRoleServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register role service: %v", err)
	}
}
