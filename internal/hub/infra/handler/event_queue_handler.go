package handler

import (
	_dto "common/domain/dto"
	"context"
	"hub/infra/mapper"
	"hub/internal/dto"
	"hub/internal/enums"
	"hub/internal/usecase"
	hubpb "pb/types/hub"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventQueueService struct {
	hubpb.UnimplementedEventQueueServiceServer
	UC     *usecase.EventQueueUsecase
	Mapper *mapper.EventQueueMapper
}

func NewEventQueueService(
	uc *usecase.EventQueueUsecase,
	mapper *mapper.EventQueueMapper,
) *EventQueueService {
	return &EventQueueService{
		UC:     uc,
		Mapper: mapper,
	}
}

// @Summary Lấy danh sách event queue
// @Description Lấy danh sách event queue với phân trang và filter
// @Tags Event Queue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param text query string false "Tìm kiếm theo tên hoặc loại event"
// @Param profileId query int false "Profile ID"
// @Param organizationId query int false "Organization ID"
// @Param status query []int false "Trạng thái (10: Pending, 20: Processing, 30: Completed, 40: Failed, 50: Cancelled)"
// @Param eventType query []string false "Loại event"
// @Param fromDate query string false "Từ ngày (YYYY-MM-DD)"
// @Param toDate query string false "Đến ngày (YYYY-MM-DD)"
// @Router /event-queue/list [get]
func (s *EventQueueService) List(ctx context.Context, req *hubpb.EventQueueListRequest) (*hubpb.EventQueueListResponse, error) {
	searchDTO := dto.EventQueueSearchDTO{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
	}

	if req.Text != nil {
		searchDTO.Text = req.Text
	}
	if req.ProfileId != nil {
		searchDTO.ProfileID = req.ProfileId
	}
	if req.OrganizationId != nil {
		searchDTO.OrganizationID = req.OrganizationId
	}
	if len(req.Status) > 0 {
		statuses := make([]enums.EEventQueueStatus, len(req.Status))
		for i, s := range req.Status {
			statuses[i] = enums.EEventQueueStatus(s)
		}
		searchDTO.Status = statuses
	}
	if len(req.EventType) > 0 {
		eventTypes := make([]string, len(req.EventType))
		copy(eventTypes, req.EventType)
		searchDTO.EventType = eventTypes
	}
	if req.FromDate != nil {
		searchDTO.FromDate = req.FromDate
	}
	if req.ToDate != nil {
		searchDTO.ToDate = req.ToDate
	}

	events, total, err := s.UC.Search(ctx, searchDTO)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	eventsPb := s.Mapper.ListEntityToPb(events)

	return &hubpb.EventQueueListResponse{
		Data:  eventsPb,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy chi tiết event queue
// @Description Lấy thông tin chi tiết của một event queue
// @Tags Event Queue
// @Accept json
// @Produce json
// @Param id path int true "ID của event queue"
// @Security BearerAuth
// @Router /event-queue/detail/{id} [get]
func (s *EventQueueService) Detail(ctx context.Context, req *sharepb.IdRequest) (*hubpb.EventQueueDTO, error) {
	event, err := s.UC.Detail(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return s.Mapper.EntityToPb(event), nil
}

// @Summary Tạo mới event queue
// @Description Tạo mới một event queue
// @Tags Event Queue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body hubpb.EventQueueSaveDTO true "Thông tin event queue"
// @Router /event-queue [post]
func (s *EventQueueService) Create(ctx context.Context, req *hubpb.EventQueueSaveDTO) (*hubpb.EventQueueDTO, error) {
	saveDTO := dto.EventQueueSaveDTO{
		EventType: req.EventType,
		EventName: req.EventName,
	}

	if req.EventData != nil {
		saveDTO.EventData = req.EventData
	}
	if req.Metadata != nil {
		saveDTO.Metadata = req.Metadata
	}
	if req.Status != nil {
		status := enums.EEventQueueStatus(*req.Status)
		saveDTO.Status = &status
	}
	if req.ProfileId != nil {
		saveDTO.ProfileID = req.ProfileId
	}
	if req.OrganizationId != nil {
		saveDTO.OrganizationID = req.OrganizationId
	}
	if req.ScheduledAt != nil {
		saveDTO.ScheduledAt = req.ScheduledAt
	}
	if req.RetryCount != nil {
		saveDTO.RetryCount = req.RetryCount
	}
	if req.MaxRetries != nil {
		saveDTO.MaxRetries = req.MaxRetries
	}
	if req.ErrorMessage != nil {
		saveDTO.ErrorMessage = req.ErrorMessage
	}

	event, err := s.UC.Create(ctx, saveDTO)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return s.Mapper.EntityToPb(event), nil
}

// @Summary Cập nhật event queue
// @Description Cập nhật thông tin event queue
// @Tags Event Queue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của event queue"
// @Param body body hubpb.EventQueueSaveDTO true "Thông tin event queue"
// @Router /event-queue/edit/{id} [put]
func (s *EventQueueService) Update(ctx context.Context, req *hubpb.EventQueueSaveDTO) (*hubpb.EventQueueDTO, error) {
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "ID không được để trống")
	}

	saveDTO := dto.EventQueueSaveDTO{
		EventType: req.EventType,
		EventName: req.EventName,
	}

	if req.EventData != nil {
		saveDTO.EventData = req.EventData
	}
	if req.Metadata != nil {
		saveDTO.Metadata = req.Metadata
	}
	if req.Status != nil {
		status := enums.EEventQueueStatus(*req.Status)
		saveDTO.Status = &status
	}
	if req.ProfileId != nil {
		saveDTO.ProfileID = req.ProfileId
	}
	if req.OrganizationId != nil {
		saveDTO.OrganizationID = req.OrganizationId
	}
	if req.ScheduledAt != nil {
		saveDTO.ScheduledAt = req.ScheduledAt
	}
	if req.RetryCount != nil {
		saveDTO.RetryCount = req.RetryCount
	}
	if req.MaxRetries != nil {
		saveDTO.MaxRetries = req.MaxRetries
	}
	if req.ErrorMessage != nil {
		saveDTO.ErrorMessage = req.ErrorMessage
	}

	event, err := s.UC.Update(ctx, req.Id, saveDTO)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return s.Mapper.EntityToPb(event), nil
}

// @Summary Xóa event queue
// @Description Xóa mềm event queue
// @Tags Event Queue
// @Accept json
// @Produce json
// @Param id path int true "ID của event queue"
// @Security BearerAuth
// @Router /event-queue/{id} [delete]
func (s *EventQueueService) Delete(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	err := s.UC.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return &sharepb.SubmitResponse{
		Message: "Xóa event queue thành công",
	}, nil
}

// @Summary Cập nhật trạng thái event queue
// @Description Cập nhật trạng thái của event queue
// @Tags Event Queue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của event queue"
// @Param body body hubpb.EventQueueStatusDTO true "Thông tin trạng thái"
// @Router /event-queue/status/{id} [put]
func (s *EventQueueService) UpdateStatus(ctx context.Context, req *hubpb.EventQueueStatusDTO) (*hubpb.EventQueueDTO, error) {
	statusDTO := dto.EventQueueStatusDTO{
		ID:     req.Id,
		Status: enums.EEventQueueStatus(req.Status),
	}

	if req.ErrorMessage != nil {
		statusDTO.ErrorMessage = req.ErrorMessage
	}

	event, err := s.UC.UpdateStatus(ctx, statusDTO)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return s.Mapper.EntityToPb(event), nil
}
