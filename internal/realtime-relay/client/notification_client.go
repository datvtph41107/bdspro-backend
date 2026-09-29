package client

import (
	clients "pb/clients"
)

type NotificationClient struct {
	*clients.NotificationClient
}

func NewNotificationClient(rpcClient *clients.NotificationClient) *NotificationClient {
	return &NotificationClient{NotificationClient: rpcClient}
}

// func (c *NotificationClient) NotifyToUser(
// 	ctx context.Context,
// 	typeNoti int32,
// 	userId uint64,
// 	attachData []string,
// 	targetID uint64,
// 	isMerge bool,
// ) error {
// 	if c.NotificationClient == nil {
// 		return fmt.Errorf("notification client not available")
// 	}

// 	// Tạo title và body từ attachData
// 	title := "Tin nhắn mới"
// 	body := ""
// 	if len(attachData) > 0 {
// 		body = attachData[0]
// 	}

// 	// Tạo data payload
// 	data := map[string]string{
// 		"type":       fmt.Sprintf("%d", typeNoti),
// 		"targetId":   fmt.Sprintf("%d", targetID),
// 		"link":       fmt.Sprintf("/chat/%d", targetID),
// 		"targetType": fmt.Sprintf("%d", typeNoti),
// 	}

// 	// Gọi PushByUserId để chỉ push notification, không lưu vào DB
// 	_, err := c.Client.PushByUserId(ctx, &notificationpb.PushByUserIdRequest{
// 		UserId: userId,
// 		Title:  title,
// 		Body:   body,
// 		Data:   data,
// 	})
// 	return err
// }
