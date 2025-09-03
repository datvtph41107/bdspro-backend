package initial

import (
	"organization/infrastructure/server/grpc"
)

type InitialApp struct {
	GRPCServer *grpc.GRPCServer
}

func NewInitialApp(
	grpcServer *grpc.GRPCServer,
) *InitialApp {
	return &InitialApp{
		GRPCServer: grpcServer,
	}
}
