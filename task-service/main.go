package main

import (
	"common"
	"fmt"
	"task/config"
	"task/internal/server"
	"time"

	"github.com/go-co-op/gocron"
	"google.golang.org/grpc"

	pb_task "pb/types/task"
)

func main() {
	s := gocron.NewScheduler(time.Local)
	s.Every(5).Second().At("00:00").Do(func() {
		fmt.Println("Hello, World!", time.Now().Format("2006-01-02 15:04:05"))
	})
	s.StartAsync()

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
