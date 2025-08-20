package handler

import (
	"context"
	appointmentpb "pb/types/appointment"

	"github.com/tuandq2112/appointment-service/internal/interface/transformer"
	"github.com/tuandq2112/appointment-service/internal/interface/validator"
	"github.com/tuandq2112/appointment-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppointmentHandler struct {
	appointmentpb.UnimplementedAppointmentServiceServer
	AppointmentUsecase usecase.AppointmentUsecase
	Transformer        transformer.AppointmentTransformer
	Validator          validator.AppointmentValidator
}

func NewAppointmentHandler(
	appointmentUsecase usecase.AppointmentUsecase,
	transformer transformer.AppointmentTransformer,
	validator validator.AppointmentValidator,
) *AppointmentHandler {
	return &AppointmentHandler{
		AppointmentUsecase: appointmentUsecase,
		Transformer:        transformer,
		Validator:          validator,
	}
}

// @Summary Tạo lịch hẹn
// @Description Tạo lịch hẹn
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param appointment body appointmentpb.CreateAppointmentRequest true "Thông tin lịch hẹn"
// @Security BearerAuth
// @Success 200 {object} appointmentpb.Appointment "Thành công"
// @Router /v2/appointment [post]
func (h *AppointmentHandler) CreateAppointment(ctx context.Context, req *appointmentpb.CreateAppointmentRequest) (*appointmentpb.Appointment, error) {
	if err := h.Validator.ValidateCreateAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	appointment := h.Transformer.CreateAppointmentRequestToEntity(req)

	createdAppointment, err := h.AppointmentUsecase.CreateAppointment(ctx, appointment)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToAppointmentResponse(createdAppointment, nil, nil, nil), nil
}

// @Summary Lấy lịch hẹn
// @Description Lấy lịch hẹn
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param id path uint32 true "ID lịch hẹn"
// @Security BearerAuth
// @Success 200 {object} appointmentpb.Appointment "Thành công"
// @Router /v2/appointment/{id} [get]
func (h *AppointmentHandler) GetAppointment(ctx context.Context, req *appointmentpb.GetAppointmentRequest) (*appointmentpb.Appointment, error) {
	if err := h.Validator.ValidateGetAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	appointment, err := h.AppointmentUsecase.GetAppointment(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToAppointmentResponse(appointment.Appointment, appointment.Deal, appointment.Products, appointment.Participants), nil
}

// @Summary Cập nhật lịch hẹn
// @Description Cập nhật lịch hẹn
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param id path uint32 true "ID lịch hẹn"
// @Param appointment body appointmentpb.UpdateAppointmentRequest true "Thông tin lịch hẹn"
// @Security BearerAuth
// @Success 200 {object} appointmentpb.Appointment "Thành công"
// @Router /v2/appointment/{id} [put]
func (h *AppointmentHandler) UpdateAppointment(ctx context.Context, req *appointmentpb.UpdateAppointmentRequest) (*appointmentpb.Appointment, error) {
	if err := h.Validator.ValidateUpdateAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	appointment := h.Transformer.UpdateAppointmentRequestToEntity(req)

	updatedAppointment, err := h.AppointmentUsecase.UpdateAppointment(ctx, appointment)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToAppointmentResponse(updatedAppointment, nil, nil, nil), nil
}

// @Summary Xóa lịch hẹn
// @Description Xóa lịch hẹn
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param id path uint32 true "ID lịch hẹn"
// @Security BearerAuth
// @Success 200 {object} appointmentpb.DeleteAppointmentResponse "Thành công"
// @Router /v2/appointment/{id} [delete]
func (h *AppointmentHandler) DeleteAppointment(ctx context.Context, req *appointmentpb.DeleteAppointmentRequest) (*appointmentpb.DeleteAppointmentResponse, error) {
	if err := h.Validator.ValidateDeleteAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.AppointmentUsecase.DeleteAppointment(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &appointmentpb.DeleteAppointmentResponse{
		Id: req.Id,
	}, nil
}

// @Summary Lấy danh sách lịch hẹn
// @Description Lấy danh sách lịch hẹn
// @Tags Lịch hẹn
// @Accept json
// @Produce json
// @Param page query uint32 false "Trang"
// @Param size query uint32 false "Kích thước trang"
// @Security BearerAuth
// @Success 200 {object} appointmentpb.ListAppointmentsResponse "Thành công"
// @Router /v2/appointment [get]
func (h *AppointmentHandler) ListAppointments(ctx context.Context, req *appointmentpb.ListAppointmentsRequest) (*appointmentpb.ListAppointmentsResponse, error) {
	if err := h.Validator.ValidateListAppointmentsRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	appointments, total, err := h.AppointmentUsecase.ListAppointments(ctx, nil, page, size)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToListAppointmentsResponse(appointments, int32(total)), nil
}
