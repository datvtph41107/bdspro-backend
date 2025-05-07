package transmit

import (
	pb_noti "pb/types/notification"
	pb_task "pb/types/task"

	"google.golang.org/grpc"
)

type TransmitService struct {
	NotiClient pb_noti.NotificationServiceClient
	TaskClient pb_task.TaskServiceClient
}

func NewNotiService() *TransmitService {
	conn, err := grpc.Dial("localhost:8004", grpc.WithInsecure())
	if err != nil {
		return &TransmitService{}
	}
	client := pb_noti.NewNotificationServiceClient(conn)

	connTask, err := grpc.Dial("localhost:8008", grpc.WithInsecure())
	if err != nil {
		return &TransmitService{}
	}
	clientTask := pb_task.NewTaskServiceClient(connTask)
	return &TransmitService{
		NotiClient: client,
		TaskClient: clientTask,
	}
}
