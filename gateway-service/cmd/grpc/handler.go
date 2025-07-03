package cmd

import notificationpb "pb/types/notification"

type ServerGrpc struct {
	notificationpb.UnimplementedNotificationServiceServer
}
