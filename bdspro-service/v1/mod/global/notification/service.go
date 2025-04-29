package bdspro_noti

import (
	pb_noti "pb/types/notification"

	"google.golang.org/grpc"
)

type NotiService struct {
	NotiClient pb_noti.NotificationServiceClient
}

func NewNotiService() *NotiService {
	conn, err := grpc.Dial("localhost:8004", grpc.WithInsecure())
	if err != nil {
		return &NotiService{}
	}
	client := pb_noti.NewNotificationServiceClient(conn)
	return &NotiService{
		NotiClient: client,
	}
}
