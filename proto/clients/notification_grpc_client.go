package clients

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"context"
	"fmt"
	"log"

	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"

	"google.golang.org/grpc"
)

type NotificationClient struct {
	Client        notificationpb.InternalNotificationServiceClient
	historyClient notificationpb.HistoryServiceClient
}

func BindNotificationClient(conn grpc.ClientConnInterface) *NotificationClient {
	return &NotificationClient{
		Client:        notificationpb.NewInternalNotificationServiceClient(conn),
		historyClient: notificationpb.NewHistoryServiceClient(conn),
	}
}

func (s *NotificationClient) NotifyToUser(
	ctx context.Context,
	typeNoti int32,
	userId *uint64,
	attachData []string,
	targetID uint64,
	isMerge bool,
) error {
	_, err := s.Client.Create(ctx, &sharepb.NotificationRequest{
		AttachData: attachData,
		Link:       fmt.Sprintf("/news-feed/%d", targetID),
		Type:       typeNoti,
		OwnerId:    *userId,
		TargetId:   targetID,
		TargetType: typeNoti,
		IsMerge:    isMerge,
	})
	return err
}

func (s *NotificationClient) SendBatch(
	ctx context.Context,
	notifications []*sharepb.NotificationRequest,
) error {
	_, err := s.Client.SendBatch(ctx, &notificationpb.NotiBatchRequest{
		Datas:        notifications,
		SendToDevice: false,
	})
	return err
}

func (s *NotificationClient) BatchPersonConfig(ctx context.Context, configs []*sharepb.PersonConfigV3Proto) error {
	_, err := s.Client.BatchPersonConfig(ctx, &notificationpb.BatchPersonConfigRequest{
		Configs: configs,
	})
	return err
}

func (s *NotificationClient) CreateHistory(ctx context.Context, history *_dto.HistoryDTO) error {
	_, err := s.historyClient.Create(ctx, &notificationpb.HistoryDTO{
		ActionType: int32(history.ActionType),
		TargetId:   history.TargetId,
		OwnerId:    history.OwnerID,
		Note:       history.Note,
		Title:      history.Title,
	})
	return err
}

func (s *NotificationClient) SendNotification(ctx context.Context, noti *_dto.NotificationDTO) error {
	_, err := s.Client.Create(ctx, &sharepb.NotificationRequest{
		Title:   noti.Title,
		Message: noti.Message,
		Type:    int32(noti.NotificationType),
	})
	return err
}

func (s *NotificationClient) SendBatchNotification(ctx context.Context, noti *_dto.NotificationBatchDTO) error {
	notifications := make([]*sharepb.NotificationRequest, 0)
	for _, notification := range noti.Datas {
		notifications = append(notifications, &sharepb.NotificationRequest{
			Title:   notification.Title,
			Message: notification.Message,
			Type:    int32(notification.NotificationType),
		})
	}
	_, err := s.Client.SendBatch(ctx, &notificationpb.NotiBatchRequest{
		Datas:        notifications,
		SendToDevice: false,
	})
	return err
}

func (s *NotificationClient) CreateCRMHistory(ctx context.Context, history *_dto.HistoryDTO) error {
	_, err := s.historyClient.CreateInternal(ctx, &notificationpb.HistoryDTO{
		ActionType: int32(history.ActionType),
		TargetId:   history.TargetId,
		OwnerId:    history.OwnerID,
		Note:       history.Note,
		Title:      history.Title,
	})
	return err
}

// CreateAdminHistory tạo lịch sử admin
func (s *NotificationClient) CreateAdminHistory(ctx context.Context,
	adminID uint64,
	targetID uint64,
	targetType int32,
	actionType int32,
	title string,
	notes []string,
	preStage string,
	afterStage string,
	ownerID *uint64,
	ownerType int32,
	adminRole string,
	ipAddress string,
	userAgent string,
) error {
	log.Printf("Creating admin history: adminID=%d, targetID=%d, actionType=%d, title=%s", adminID, targetID, actionType, title)

	response, err := s.Client.CreateAdminHistory(ctx, &notificationpb.AdminHistoryCreateRequest{
		TargetId:   targetID,
		TargetType: targetType,
		ActionType: actionType,
		Title:      title,
		Note:       notes,
		PreStage:   preStage,
		AfterStage: afterStage,
		AdminId:    adminID,
		OwnerId:    ownerID,
		OwnerType:  ownerType,
		IsInternal: false,
		AdminRole:  adminRole,
		IpAddress:  ipAddress,
		UserAgent:  userAgent,
	})
	if err != nil {
		log.Printf("Error creating admin history: %v", err)
		return err
	}

	log.Printf("Admin history created successfully: ID=%d", response.Id)
	return nil
}

// CreateAdminHistoryBatch tạo hàng loạt lịch sử admin
func (s *NotificationClient) CreateAdminHistoryBatch(ctx context.Context, requests []*notificationpb.AdminHistoryCreateRequest) error {
	_, err := s.Client.CreateAdminHistoryBatch(ctx, &notificationpb.AdminHistoryBatchRequest{
		Datas: requests,
	})
	return err
}

func (s *NotificationClient) CreateHistoryAuth(ctx context.Context, req *notificationpb.HistoryAuthCreateRequest) (*notificationpb.HistoryAuthDTO, error) {
	return s.Client.CreateHistoryAuth(ctx, req)
}

// LogAdminLogin ghi log khi admin đăng nhập thành công
func (s *NotificationClient) LogAdminLogin(ctx context.Context, adminID uint64, adminRole string, ipAddress string, userAgent string) error {
	return s.CreateAdminHistory(ctx,
		adminID,                        // adminID
		adminID,                        // targetID (chính admin đó)
		int32(_enum.TargetHistoryLead), // targetType
		int32(_enum.HistoryCreateStep), // actionType
		"Admin đăng nhập hệ thống",     // title
		[]string{"Đăng nhập thành công vào hệ thống"}, // notes
		"Chưa đăng nhập",            // preStage
		"Đã đăng nhập",              // afterStage
		&adminID,                    // ownerID
		int32(_enum.EOwnerOfMember), // ownerType
		adminRole,                   // adminRole
		ipAddress,                   // ipAddress
		userAgent,                   // userAgent
	)
}

// CreateNotification tạo thông báo cho user với đầy đủ thông tin
func (c *NotificationClient) CreateNotification(
	ctx context.Context,
	avatar string,
	title string,
	message []string,
	notificationType _enum.ENotificationType,
	targetId *uint64,
	ownerId uint64,
	ownerOf _enum.EOwnerOf,
	attachData []string,
) error {
	log.Println("avatar", avatar)
	_, err := c.Client.CreateToOwner(ctx, &sharepb.NotificationRequest{
		Avatar:           avatar,
		Title:            title,
		Message:          message,
		Type:             int32(notificationType),
		OwnerId:          ownerId,
		OwnerOf:          uint32(ownerOf),
		NotificationType: uint32(notificationType),
		TargetId:         *targetId,
		AttachData:       attachData,
		SendToDevice:     false,
	})
	return err
}

func (c *NotificationClient) RemoveNotification(
	ctx context.Context,
	ownerOf _enum.EOwnerOf,
	ownerId uint64,
	targetId uint64,
	notificationType _enum.ENotificationType,
) error {
	_, err := c.Client.RemoveToOwner(ctx, &sharepb.NotificationRequest{
		OwnerOf:          uint32(ownerOf),
		OwnerId:          ownerId,
		TargetId:         targetId,
		NotificationType: uint32(notificationType),
	})
	return err
}

func (s *NotificationClient) PushByUserId(
	ctx context.Context,
	userId uint64,
	title string,
	body string,
	data map[string]string,
) error {
	_, err := s.Client.PushByUserId(ctx, &notificationpb.PushByUserIdRequest{
		UserId: userId,
		Title:  title,
		Body:   body,
		Data:   data,
	})
	return err
}

func (s *NotificationClient) CreatePropertyHistoryActivity(
	ctx context.Context,
	subjectId uint64,
	subjectType uint32,
	action uint32,
	description string,
	actorId uint64,
) error {
	_, err := s.Client.CreatePropertyHistoryActivity(ctx, &notificationpb.CreatePropertyHistoryRequest{
		SubjectId:   subjectId,
		Action:      action,
		Description: description,
		ActorId:     actorId,
	})
	return err
}

func (s *NotificationClient) GetPersonConfigs(
	ctx context.Context,
	userID uint64,
) ([]*sharepb.PersonConfigV3Proto, error) {
	if userID == 0 {
		return []*sharepb.PersonConfigV3Proto{}, nil
	}
	if s == nil || s.Client == nil {
		return nil, fmt.Errorf("notification grpc client chưa được khởi tạo")
	}

	resp, err := s.Client.GetPersonConfigs(
		ctx,
		&notificationpb.GetPersonConfigsRequest{
			UserId: userID,
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.GetData(), nil
}
