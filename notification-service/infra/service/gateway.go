package service

import (
	_dto "common/dto"
	_utils "common/utils"
	"context"
	"noti/internal/domain"
	"noti/internal/dto"
	"noti/internal/enums"
	"noti/internal/usecase"
	pb_noti "pb/types/notification"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GatewayNotiService struct {
	Usecase *usecase.NotificationUsecase
	pb_noti.UnimplementedGatewayNotificationServiceServer
}

func NewGatewayNotiService(usecase *usecase.NotificationUsecase) *GatewayNotiService {
	return &GatewayNotiService{
		Usecase: usecase,
	}
}

// @Summary Tạo thông báo mới
// @Description API này không bắn thông báo tới thiết bị của người dùng, chỉ tạo thông báo
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body pb_noti.GatewayNotiNewRequest true "Chi tiết thông báo"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /new [post]
func (s *GatewayNotiService) Create(ctx context.Context, req *pb_noti.GatewayNotiNewRequest) (*pb_noti.NotificationDTO, error) {
	dto := &dto.NotiNewRequest{
		Title:   req.Title,
		Message: []string{req.Message},
		Link:    req.Link,
		Type:    enums.TypeEnum(req.Type),
	}
	notification, err := s.Usecase.CreateNotification(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method Send not implemented")
	}

	response := s.mapToPB(notification)
	return response, nil
}

// @Summary Số thông báo mới chưa đọc
// @Description Hiển thị số lượng thông báo chưa đọc của người dùng
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]int "count"
// @Failure 500 {object} map[string]string "error"
// @Router /count [get]
func (s *GatewayNotiService) Count(ctx context.Context, req *pb_noti.EmptyRequest) (*pb_noti.CountResponse, error) {
	count := s.Usecase.Count(ctx)
	return &pb_noti.CountResponse{
		Count: count,
	}, nil
}

// @Summary Đánh dấu thông báo đã đọc
// @Description Đánh dấu thông báo với ID cụ thể là đã đọc
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của thông báo"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /read/{id} [put]
func (s *GatewayNotiService) Read(ctx context.Context, req *pb_noti.IdRequest) (*pb_noti.StatusResponse, error) {
	_, err := s.Usecase.Read(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method Read not implemented")
	}
	return &pb_noti.StatusResponse{
		Status:  true,
		Message: "Thông báo đã được đánh dấu là đã đọc",
	}, nil
}

// @Summary Đánh dấu tất cả thông báo là đã đọc
// @Description Đánh dấu tất cả thông báo là đã đọc
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /read-all [put]
func (s *GatewayNotiService) ReadAll(ctx context.Context, req *pb_noti.EmptyRequest) (*pb_noti.StatusResponse, error) {
	_, err := s.Usecase.ReadAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method ReadAll not implemented")
	}
	return &pb_noti.StatusResponse{
		Status:  true,
		Message: "Thông báo đã được đánh dấu là đã đọc",
	}, nil
}

// @Summary Xóa thông báo
// @Description Xóa thông báo với ID cụ thể
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của thông báo"
// @Router /remove/{id} [delete]
func (s *GatewayNotiService) Remove(ctx context.Context, req *pb_noti.IdRequest) (*pb_noti.StatusResponse, error) {
	_, err := s.Usecase.Remove(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method Remove not implemented")
	}
	return &pb_noti.StatusResponse{
		Status:  true,
		Message: "Thông báo đã được xóa",
	}, nil
}

// @Summary Lấy danh sách thông báo
// @Description Trả về danh sách thông báo của người dùng
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang cần lấy (default = 0)"
// @Param size query int false "Số lượng trên mỗi trang (default = 200)"
// @Success 200 {object} []pb_noti.NotificationDTO "Danh sách thông báo"
// @Failure 500 {object} map[string]string "error"
// @Router /list [get]
func (s *GatewayNotiService) List(ctx context.Context, req *pb_noti.ListRequest) (*pb_noti.ListResponse, error) {
	query := dto.SearchNotiRequest{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		IsRead: &req.IsRead,
	}
	notifications, total, err := s.Usecase.Search(ctx, &query)
	var responses []*pb_noti.NotificationDTO
	for _, noti := range notifications {
		responses = append(responses, s.mapToPB(&noti))
	}

	return &pb_noti.ListResponse{
		Data:  responses,
		Total: int32(total),
	}, err
}

// todo: implement
func (s *GatewayNotiService) PushByToken(ctx context.Context, req *pb_noti.PushTokenRequest) (*pb_noti.PushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushByToken not implemented")
}
func (s *GatewayNotiService) PushByTopic(ctx context.Context, req *pb_noti.PushTopicRequest) (*pb_noti.PushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushByTopic not implemented")
}

// -- mapping ---
func (s *GatewayNotiService) mapToPB(noti *domain.NotificationEntity) *pb_noti.NotificationDTO {
	return &pb_noti.NotificationDTO{
		Id:        noti.ID,
		UserId:    noti.OwnerID,
		Title:     noti.Title,
		Message:   noti.Message,
		Link:      noti.Link,
		Type:      int32(noti.Type),
		Timestamp: _utils.FormatTimeToString(noti.CreatedAt),
		IsRead:    noti.IsRead,
		TargetId:  &noti.TargetID,
	}
}
