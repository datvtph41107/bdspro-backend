// Package grpcgateway owns generated grpc-gateway handler registration for
// the Gateway's explicitly opened upstream ClientConns.
package grpcgateway

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type registrar func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

func registerAll(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn, registrars ...registrar) error {
	for _, register := range registrars {
		if err := register(ctx, mux, conn); err != nil {
			return err
		}
	}
	return nil
}
