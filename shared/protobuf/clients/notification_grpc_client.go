package clients

import (
	"common"
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	"time"

	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"

	"google.golang.org/grpc"
)

type NotificationClient struct {
	Client        notificationpb.InternalNotificationServiceClient
	historyClient notificationpb.HistoryServiceClient
}

func MakeNotificationClient(env_runtime string) *NotificationClient {
	notificationClient := &NotificationClient{}
	go func() {
		name := "notification"
		url := common.GetPrefixProtobufUrl(env_runtime, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:8204", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		client := notificationpb.NewInternalNotificationServiceClient(conn)
		historyClient := notificationpb.NewHistoryServiceClient(conn)
		notificationClient.Client = client
		notificationClient.historyClient = historyClient
	}()
	return notificationClient
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
func (s *NotificationClient) CreateNotification(ctx context.Context,
	title string,
	message []string,
	notificationType int32,
	link string,
	ownerID uint64,
	attachData []string,
) error {
	_, err := s.Client.Create(ctx, &sharepb.NotificationRequest{
		Title:      title,
		Message:    message,
		Type:       notificationType,
		Link:       link,
		OwnerId:    ownerID,
		AttachData: attachData,
	})
	return err
}
