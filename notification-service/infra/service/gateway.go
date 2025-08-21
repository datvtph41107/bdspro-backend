package service

import (
	_dto "common/dto"
	_utils "common/utils"
	"context"
	shared_enum "pb/enums"
	notipb "pb/types/notification"
	"time"

	"noti/infra/mapper"
	"noti/internal/domain"
	"noti/internal/dto"
	"noti/internal/enums"
	"noti/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GatewayNotiService struct {
	Usecase             *usecase.NotificationUsecase
	AdminHistoryUsecase *usecase.AdminHistoryUsecase
	AdminHistoryMapper  *mapper.AdminHistoryMapper
	notipb.UnimplementedGatewayNotificationServiceServer
}

func NewGatewayNotiService(
	usecase *usecase.NotificationUsecase,
	adminHistoryUsecase *usecase.AdminHistoryUsecase,
	adminHistoryMapper *mapper.AdminHistoryMapper,
) *GatewayNotiService {
	return &GatewayNotiService{
		Usecase:             usecase,
		AdminHistoryUsecase: adminHistoryUsecase,
		AdminHistoryMapper:  adminHistoryMapper,
	}
}

// @Summary Tạo thông báo mới
// @Description API này không bắn thông báo tới thiết bị của người dùng, chỉ tạo thông báo
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body notipb.GatewayNotiNewRequest true "Chi tiết thông báo"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /new [post]
func (s *GatewayNotiService) Create(ctx context.Context, req *notipb.GatewayNotiNewRequest) (*notipb.NotificationDTO, error) {
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
func (s *GatewayNotiService) Count(ctx context.Context, req *notipb.EmptyRequest) (*notipb.CountResponse, error) {
	count := s.Usecase.Count(ctx)
	return &notipb.CountResponse{
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
func (s *GatewayNotiService) Read(ctx context.Context, req *notipb.IdRequest) (*notipb.StatusResponse, error) {
	_, err := s.Usecase.Read(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method Read not implemented")
	}
	return &notipb.StatusResponse{
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
func (s *GatewayNotiService) ReadAll(ctx context.Context, req *notipb.EmptyRequest) (*notipb.StatusResponse, error) {
	_, err := s.Usecase.ReadAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method ReadAll not implemented")
	}
	return &notipb.StatusResponse{
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
func (s *GatewayNotiService) Remove(ctx context.Context, req *notipb.IdRequest) (*notipb.StatusResponse, error) {
	_, err := s.Usecase.Remove(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method Remove not implemented")
	}
	return &notipb.StatusResponse{
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
// @Success 200 {object} []notipb.NotificationDTO "Danh sách thông báo"
// @Failure 500 {object} map[string]string "error"
// @Router /list [get]
func (s *GatewayNotiService) List(ctx context.Context, req *notipb.ListRequest) (*notipb.ListResponse, error) {
	query := dto.SearchNotiRequest{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		IsRead: &req.IsRead,
	}
	notifications, total, err := s.Usecase.Search(ctx, &query)
	var responses []*notipb.NotificationDTO
	for _, noti := range notifications {
		responses = append(responses, s.mapToPB(&noti))
	}

	return &notipb.ListResponse{
		Data:  responses,
		Total: int32(total),
	}, err
}

// todo: implement
func (s *GatewayNotiService) PushByToken(ctx context.Context, req *notipb.PushTokenRequest) (*notipb.PushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushByToken not implemented")
}

func (s *GatewayNotiService) PushByTopic(ctx context.Context, req *notipb.PushTopicRequest) (*notipb.PushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushByTopic not implemented")
}

// @Summary Tìm kiếm lịch sử admin
// @Description Tìm kiếm lịch sử các hành động của admin
// @Tags Admin History
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang cần lấy (default = 0)"
// @Param size query int false "Số lượng trên mỗi trang (default = 20)"
// @Param adminId query int false "ID của admin"
// @Param ownerId query int false "ID của owner"
// @Param ownerType query int false "Loại owner"
// @Param targetId query int false "ID của target"
// @Param targetType query int false "Loại target"
// @Param actionType query int false "Loại hành động"
// @Param adminRole query string false "Vai trò admin"
// @Param fromDate query string false "Từ ngày (RFC3339 format)"
// @Param toDate query string false "Đến ngày (RFC3339 format)"
// @Param isInternal query bool false "Có phải internal không"
// @Success 200 {object} notipb.AdminHistoryListResponse "Danh sách lịch sử admin"
// @Failure 500 {object} map[string]string "error"
// @Router /admin-history/search [get]
func (s *GatewayNotiService) SearchAdminHistory(ctx context.Context, request *notipb.AdminHistorySearchRequest) (*notipb.AdminHistoryListResponse, error) {
	// Convert protobuf request to DTO
	searchDTO := &dto.AdminHistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(request.Page),
			Size: uint32(request.Size),
		},
	}

	if request.AdminId != nil {
		searchDTO.AdminID = request.AdminId
	}
	if request.OwnerId != nil {
		searchDTO.OwnerID = request.OwnerId
	}
	if request.OwnerType != nil {
		ownerType := shared_enum.EOwnerType(*request.OwnerType)
		searchDTO.OwnerType = &ownerType
	}
	if request.TargetId != nil {
		searchDTO.TargetID = request.TargetId
	}
	if request.TargetType != nil {
		targetType := shared_enum.ETargetHistory(*request.TargetType)
		searchDTO.TargetType = &targetType
	}
	if request.ActionType != nil {
		actionType := shared_enum.EHistory(*request.ActionType)
		searchDTO.ActionType = &actionType
	}
	if request.AdminRole != nil {
		searchDTO.AdminRole = request.AdminRole
	}
	if request.FromDate != nil {
		if fromDate, err := time.Parse(time.RFC3339, *request.FromDate); err == nil {
			searchDTO.FromDate = &fromDate
		}
	}
	if request.ToDate != nil {
		if toDate, err := time.Parse(time.RFC3339, *request.ToDate); err == nil {
			searchDTO.ToDate = &toDate
		}
	}
	if request.IsInternal != nil {
		searchDTO.IsInternal = request.IsInternal
	}

	histories, total, err := s.AdminHistoryUsecase.SearchInternal(ctx, *searchDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to search admin history: %v", err)
	}

	dtos := s.AdminHistoryMapper.EntitiesToDTOs(histories)

	// Convert DTOs to protobuf DTOs
	pbDtos := make([]*notipb.AdminHistoryDTO, len(dtos))
	for i, dto := range dtos {
		pbDtos[i] = s.convertToPbDTO(&dto)
	}

	return &notipb.AdminHistoryListResponse{
		Data:  pbDtos,
		Total: int32(total),
	}, nil
}

// convertToPbDTO chuyển đổi AdminHistoryDTO thành protobuf AdminHistoryDTO
func (s *GatewayNotiService) convertToPbDTO(dto *dto.AdminHistoryDTO) *notipb.AdminHistoryDTO {
	return &notipb.AdminHistoryDTO{
		Id:         dto.ID,
		TargetId:   dto.TargetId,
		TargetType: int32(dto.TargetType),
		ActionType: int32(dto.ActionType),
		Title:      dto.Title,
		Note:       dto.Note,
		PreStage:   dto.PreStage,
		AfterStage: dto.AfterStage,
		AdminId:    dto.AdminID,
		OwnerId:    dto.OwnerID,
		OwnerType:  int32(dto.OwnerType),
		IsInternal: dto.IsInternal,
		AdminRole:  dto.AdminRole,
		IpAddress:  dto.IPAddress,
		UserAgent:  dto.UserAgent,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
	}
}

// -- mapping ---
func (s *GatewayNotiService) mapToPB(noti *domain.NotificationEntity) *notipb.NotificationDTO {
	return &notipb.NotificationDTO{
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
