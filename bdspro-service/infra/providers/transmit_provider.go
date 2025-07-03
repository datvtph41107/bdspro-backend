package providers

import (
	notificationpb "pb/types/notification"
	pb_task "pb/types/task"

	"google.golang.org/grpc"
)

type TransmitProvider struct {
	NotiClient notificationpb.NotificationServiceClient
	TaskClient pb_task.TaskServiceClient
}

func NewNotiService() *TransmitProvider {
	conn, err := grpc.Dial("localhost:8204", grpc.WithInsecure())
	if err != nil {
		return &TransmitProvider{}
	}
	client := notificationpb.NewNotificationServiceClient(conn)

	connTask, err := grpc.Dial("localhost:8008", grpc.WithInsecure())
	if err != nil {
		return &TransmitProvider{}
	}
	clientTask := pb_task.NewTaskServiceClient(connTask)
	return &TransmitProvider{
		NotiClient: client,
		TaskClient: clientTask,
	}
}
