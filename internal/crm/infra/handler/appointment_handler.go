package handler

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
	"time"

	"crm/infra/job"
	"crm/infra/mapper"
	"crm/infra/validator"
	"crm/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppointmentHandler struct {
	crmpb.UnimplementedAppointmentServiceServer
	AppointmentUsecase usecase.AppointmentUsecase
	Transformer        mapper.AppointmentTransformer
	Validator          validator.AppointmentValidator
	DashboardStatsJob  *job.DashboardStatsJob
	SyncProvider       *_utils.SyncUtil
}

func NewAppointmentHandler(
	appointmentUsecase usecase.AppointmentUsecase,
	transformer mapper.AppointmentTransformer,
	validator validator.AppointmentValidator,
	dashboardStatsJob *job.DashboardStatsJob,
	SyncProvider *_utils.SyncUtil,
) *AppointmentHandler {
	return &AppointmentHandler{
		AppointmentUsecase: appointmentUsecase,
		Transformer:        transformer,
		Validator:          validator,
		DashboardStatsJob:  dashboardStatsJob,
		SyncProvider:       SyncProvider,
	}
}

// @Summary Tạo lịch hẹn
// @Description Tạo lịch hẹn
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param appointment body crmpb.CreateAppointmentRequest true "Thông tin lịch hẹn"
// @Security BearerAuth
// @Success 200 {object} crmpb.Appointment "Thành công"
// @Router /v2/appointment [post]
func (h *AppointmentHandler) CreateAppointment(ctx context.Context, req *crmpb.CreateAppointmentRequest) (*crmpb.Appointment, error) {
	if err := h.Validator.ValidateCreateAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	appointment := h.Transformer.CreateAppointmentRequestToEntity(req)

	createdAppointment, err := h.AppointmentUsecase.CreateAppointment(ctx, appointment)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToAppointmentResponse(createdAppointment, nil, nil, nil, nil), nil
}

// @Summary Lấy lịch hẹn
// @Description Lấy lịch hẹn
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param id path uint32 true "ID lịch hẹn"
// @Security BearerAuth
// @Success 200 {object} crmpb.Appointment "Thành công"
// @Router /v2/appointment/{id} [get]
func (h *AppointmentHandler) GetAppointment(ctx context.Context, req *crmpb.GetAppointmentRequest) (*crmpb.Appointment, error) {
	if err := h.Validator.ValidateGetAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	appointment, err := h.AppointmentUsecase.GetAppointment(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToAppointmentResponse(appointment.Appointment, appointment.Deal, appointment.Products, appointment.Participants, appointment.DealContract), nil
}

// @Summary Cập nhật lịch hẹn
// @Description Cập nhật lịch hẹn
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param id path uint32 true "ID lịch hẹn"
// @Param appointment body crmpb.UpdateAppointmentRequest true "Thông tin lịch hẹn"
// @Security BearerAuth
// @Success 200 {object} crmpb.Appointment "Thành công"
// @Router /v2/appointment/{id} [put]
func (h *AppointmentHandler) UpdateAppointment(ctx context.Context, req *crmpb.UpdateAppointmentRequest) (*crmpb.Appointment, error) {
	if err := h.Validator.ValidateUpdateAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	appointment := h.Transformer.UpdateAppointmentRequestToEntity(req)

	updatedAppointment, err := h.AppointmentUsecase.UpdateAppointment(ctx, appointment)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToAppointmentResponse(updatedAppointment, nil, nil, nil, nil), nil
}

// @Summary Xóa lịch hẹn
// @Description Xóa lịch hẹn
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param id path uint32 true "ID lịch hẹn"
// @Security BearerAuth
// @Success 200 {object} crmpb.DeleteAppointmentResponse "Thành công"
// @Router /v2/appointment/{id} [delete]
func (h *AppointmentHandler) DeleteAppointment(ctx context.Context, req *crmpb.DeleteAppointmentRequest) (*crmpb.DeleteAppointmentResponse, error) {
	if err := h.Validator.ValidateDeleteAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.AppointmentUsecase.DeleteAppointment(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &crmpb.DeleteAppointmentResponse{
		Id: req.Id,
	}, nil
}

// @Summary Lấy danh sách lịch hẹn
// @Description Lấy danh sách lịch hẹn
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param title query string false "Tiêu đề"
// @Param status query uint32 false "Trạng thái"
// @Param startTime query string false "Thời gian bắt đầu"
// @Param endTime query string false "Thời gian kết thúc"
// @Param dealIds query []uint32 false "Danh sách ID deal"
// @Param transactionId query uint64 false "ID giao dịch"
// @Param page query int32 false "Trang"
// @Param size query int32 false "Kích thước trang"
// @Param sort query string false "Sắp xếp: 'startTime,asc' (tăng dần) hoặc 'startTime,desc' (giảm dần, mặc định)"
// @Param participantIds query []uint32 false "Danh sách user IDs của participants"
// @Param fromDate query string false "Lọc từ ngày (format: YYYY-MM-DD)"
// @Param toDate query string false "Lọc đến ngày (format: YYYY-MM-DD)"
// @Param containProduct query bool false "Lọc theo có product (true = có, false = không có)"
// @Param containDeal query bool false "Lọc theo có deal (true = có, false = không có)"
// @Param containTransaction query bool false "Lọc theo có transaction (true = có, false = không có)"
// @Param statusFilter query string false "Lọc theo status (format: '10,20' - parse thành []int32)"
// @Security BearerAuth
// @Success 200 {object} crmpb.ListAppointmentsResponse "Thành công"
// @Router /v2.1/appointment [get]
func (h *AppointmentHandler) ListAppointments(ctx context.Context, req *crmpb.ListAppointmentsRequest) (*crmpb.ListAppointmentsResponse, error) {
	if err := h.Validator.ValidateListAppointmentsRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pagable := _dto.Pagable{}
	if req.Page != nil {
		pagable.Page = uint32(*req.Page)
	}
	if req.Size != nil {
		pagable.Size = uint32(*req.Size)
	}
	if req.Sort != nil && *req.Sort != "" {
		pagable.Sort = *req.Sort
	}

	// Build filter từ request
	filter := make(map[string]any)
	if req.Title != "" {
		filter["title"] = req.Title
	}
	if req.Status != 0 {
		filter["status"] = req.Status
	}
	if len(req.DealIds) > 0 {
		// Convert []uint32 to []uint32 for dealIds filter
		dealIds := make([]uint32, len(req.DealIds))
		copy(dealIds, req.DealIds)
		filter["deal_ids"] = dealIds
	}
	if req.TransactionId != nil && *req.TransactionId > 0 {
		filter["transaction_id"] = *req.TransactionId
	}
	if len(req.ParticipantIds) > 0 {
		filter["participant_ids"] = _utils.ParseInt64Array(req.ParticipantIds)
	}
	if req.FromDate != nil && *req.FromDate != "" {
		// Giữ nguyên string format YYYY-MM-DD để query DATE trực tiếp
		filter["from_date"] = *req.FromDate
	}
	if req.ToDate != nil && *req.ToDate != "" {
		// Giữ nguyên string format YYYY-MM-DD để query DATE trực tiếp
		filter["to_date"] = *req.ToDate
	}
	if req.ContainProduct != nil {
		filter["contain_product"] = *req.ContainProduct
	}
	if req.ContainDeal != nil {
		filter["contain_deal"] = *req.ContainDeal
	}
	if req.ContainTransaction != nil {
		filter["contain_transaction"] = *req.ContainTransaction
	}
	if req.StatusFilters != nil && *req.StatusFilters != "" {
		statusArray := _utils.ParseInt64Array(*req.StatusFilters)
		if len(statusArray) > 0 {
			filter["status_filters"] = statusArray
		}
	}

	appointments, total, err := h.AppointmentUsecase.ListAppointments(ctx, filter, pagable)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToListAppointmentsResponse(appointments, int32(total)), nil
}

// @Summary [Admin] Lấy danh sách tất cả lịch hẹn trong hệ thống
// @Description Lấy danh sách tất cả lịch hẹn trong hệ thống cho admin
// @Tags Admin - Lịch hẹn
// @Accept json
// @Produce json
// @Param title query string false "Tìm kiếm theo tiêu đề"
// @Param status query uint32 false "Lọc theo trạng thái"
// @Param startTime query string false "Lọc theo thời gian bắt đầu"
// @Param endTime query string false "Lọc theo thời gian kết thúc"
// @Param createdBy query uint32 false "Lọc theo người tạo"
// @Param page query int32 false "Trang"
// @Param size query int32 false "Kích thước trang"
// @Param sort query string false "Sắp xếp: 'startTime,asc' (tăng dần) hoặc 'startTime,desc' (giảm dần, mặc định)"
// @Security BearerAuth
// @Success 200 {object} crmpb.AdminListAppointmentsResponse "Thành công"
// @Router /v2/admin/appointment [get]
func (h *AppointmentHandler) AdminListAppointments(ctx context.Context, req *crmpb.AdminListAppointmentsRequest) (*crmpb.AdminListAppointmentsResponse, error) {
	if err := h.Validator.ValidateAdminListAppointmentsRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var pagable _dto.Pagable
	if req.Page != nil && req.Size != nil {
		pagable = _dto.Pagable{
			Page: uint32(*req.Page),
			Size: uint32(*req.Size),
		}
	}
	if req.Sort != nil && *req.Sort != "" {
		pagable.Sort = *req.Sort
	}

	// Build filter từ request
	filter := make(map[string]any)
	if req.Title != nil && *req.Title != "" {
		filter["title"] = *req.Title
	}
	if req.Status != nil {
		filter["status"] = *req.Status
	}
	if req.CreatedBy != nil {
		filter["created_by"] = *req.CreatedBy
	}

	appointments, total, err := h.AppointmentUsecase.ListAppointments(ctx, filter, pagable)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return h.Transformer.EntityToAdminListAppointmentsResponse(appointments, int32(total)), nil
}

func (h *AppointmentHandler) GetAppointmentStats(ctx context.Context, req *sharepb.GetStatsRequest) (*sharepb.GetStatsResponse, error) {

	firstDate, err := time.Parse(time.DateOnly, req.FirstDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	secondDate, err := time.Parse(time.DateOnly, req.SecondDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	stats, stats2, err := h.DashboardStatsJob.GetStats(ctx, firstDate, secondDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sharepb.GetStatsResponse{
		FirstDate:       stats.CalculateTime.Format(time.DateOnly),
		SecondDate:      stats2.CalculateTime.Format(time.DateOnly),
		FirstDateCount:  uint32(stats.Count),
		SecondDateCount: uint32(stats2.Count),
	}, nil
}

// @Summary Lấy danh sách lịch hẹn theo sản phẩm
// @Description Lấy danh sách lịch hẹn theo ID sản phẩm
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param id path uint64 true "ID sản phẩm"
// @Security BearerAuth
// @Success 200 {object} crmpb.ListAppointmentsResponse "Thành công"
// @Router /v2.1/appointment/list/by-product/{id} [get]
func (h *AppointmentHandler) GetAppointmentsByProduct(ctx context.Context, req *sharepb.SyncRequest) (*crmpb.ListAppointmentsResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyProductAppointments, req.Id)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	// client mặc định sẽ truyền page = 0 GET /v2/deals?page=0
	if !updated && req.Page == 0 {
		return &crmpb.ListAppointmentsResponse{}, nil
	}
	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}
	appointments, total, err := h.AppointmentUsecase.GetAppointmentsByProduct(ctx, req.Id, pagable)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// if req.Page == 0 {
	// 	time := time.Now().UnixMilli() - 24*60*60*1000
	// 	if len(appointments) > 0 && appointments[0].ProductUpdatedAt != nil {
	// 		time = appointments[0].ProductUpdatedAt.UnixMilli()
	// 	}
	// 	h.SyncProvider.PutTimestamp(ctx, key, time)
	// }
	return h.Transformer.EntityToListAppointmentsResponse(appointments, int32(total)), nil
}

// @Summary Lấy danh sách lịch hẹn theo deal
// @Description Lấy danh sách lịch hẹn theo ID deal
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param id path uint64 true "ID deal"
// @Security BearerAuth
// @Success 200 {object} crmpb.ListAppointmentsResponse "Thành công"
// @Router /v2.1/appointment/list/by-deal/{id} [get]
func (h *AppointmentHandler) GetAppointmentsByDeal(ctx context.Context, req *sharepb.IdRequest) (*crmpb.ListAppointmentsResponse, error) {
	appointments, total, err := h.AppointmentUsecase.GetAppointmentsByDeal(ctx, uint32(req.Id), 0, 100)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return h.Transformer.EntityToListAppointmentsResponse(appointments, int32(total)), nil
}

// @Summary Cập nhật trạng thái lịch hẹn
// @Description Cập nhật trạng thái lịch hẹn (approve: Đã xác nhận, cancel: Đã hủy)
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param id path uint32 true "ID lịch hẹn"
// @Param request body crmpb.UpdateAppointmentStatusRequest true "Thông tin cập nhật trạng thái"
// @Security BearerAuth
// @Success 200 {object} crmpb.Appointment "Thành công"
// @Router /v2.1/appointment/{id}/status [patch]
func (h *AppointmentHandler) UpdateAppointmentStatus(ctx context.Context, req *crmpb.UpdateAppointmentStatusRequest) (*crmpb.Appointment, error) {
	if err := h.Validator.ValidateUpdateAppointmentStatusRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	updatedAppointment, err := h.AppointmentUsecase.UpdateAppointmentStatus(ctx, req.Id, req.Action)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Lấy thông tin đầy đủ để trả về
	appointmentDTO, err := h.AppointmentUsecase.GetAppointment(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToAppointmentResponse(updatedAppointment, appointmentDTO.Deal, appointmentDTO.Products, appointmentDTO.Participants, appointmentDTO.DealContract), nil
}
