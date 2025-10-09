package initial

import "appointment/internal/server/grpc"

type InitialApp struct {
	GRPCServer *grpc.GRPCServer
}

func NewInitialApp(GRPCServer *grpc.GRPCServer) *InitialApp {
	return &InitialApp{
		GRPCServer: GRPCServer,
	}
}
