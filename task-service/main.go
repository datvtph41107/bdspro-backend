package main

import (
	"common"
	"task/config"
	"task/internal/server"

	"google.golang.org/grpc"

	pb_task "pb/types/task"
)

func main() {
	parameters := common.SetupParameters{
		ConfigFile:    config.ConfigFile,
		Engine:        nil,
		PublicRoutes:  nil,
		Prefix:        "/task",
		IgnoreSwagger: true,
	}

	common.SetupGRPCServer(parameters, func(s *grpc.Server) {
		// app, _ := InitializeApp()
		pb_task.RegisterTaskServiceServer(s, server.GrpcServer{})
	})
}
