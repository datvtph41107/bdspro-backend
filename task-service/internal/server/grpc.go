package server

import (
	pb_task "pb/types/task"
)

type GrpcServer struct {
	pb_task.UnimplementedTaskServiceServer
}
