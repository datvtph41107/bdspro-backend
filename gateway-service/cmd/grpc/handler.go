package cmd

import pb_noti "pb/types/notification"

type ServerGrpc struct {
	pb_noti.UnimplementedNotificationServiceServer
}
