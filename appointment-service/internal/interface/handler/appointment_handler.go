package handler

import (
	"context"
	appointmentpb "pb/types/appointment"

	"github.com/tuandq2112/appointment/appconfig"
	"github.com/tuandq2112/appointment/internal/domain/usecase"
	"github.com/tuandq2112/appointment/internal/interface/transformer"
	"github.com/tuandq2112/appointment/internal/interface/validator"
	"github.com/tuandq2112/appointment/pkg/utils"
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

func (h *AppointmentHandler) CreateAppointment(ctx context.Context, req *appointmentpb.CreateAppointmentRequest) (*appointmentpb.Appointment, error) {
	if err := h.Validator.ValidateCreateAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID := ctx.Value(appconfig.USER_CONTEXT).(string)
	appointment := h.Transformer.CreateAppointmentRequestToEntity(req)
	appointment.UserID = userID
	appointment.CreatedBy = userID

	createdAppointment, err := h.AppointmentUsecase.CreateAppointment(ctx, appointment)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToAppointmentResponse(createdAppointment), nil
}

func (h *AppointmentHandler) GetAppointment(ctx context.Context, req *appointmentpb.GetAppointmentRequest) (*appointmentpb.Appointment, error) {
	if err := h.Validator.ValidateGetAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	appointment, err := h.AppointmentUsecase.GetAppointment(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToAppointmentResponse(appointment), nil
}

func (h *AppointmentHandler) UpdateAppointment(ctx context.Context, req *appointmentpb.UpdateAppointmentRequest) (*appointmentpb.Appointment, error) {
	if err := h.Validator.ValidateUpdateAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID := ctx.Value(appconfig.USER_CONTEXT).(string)
	appointment := h.Transformer.UpdateAppointmentRequestToEntity(req)
	appointment.UserID = userID
	appointment.UpdatedBy = userID

	updatedAppointment, err := h.AppointmentUsecase.UpdateAppointment(ctx, appointment)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToAppointmentResponse(updatedAppointment), nil
}

func (h *AppointmentHandler) DeleteAppointment(ctx context.Context, req *appointmentpb.DeleteAppointmentRequest) (*appointmentpb.DeleteAppointmentResponse, error) {
	if err := h.Validator.ValidateDeleteAppointmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.AppointmentUsecase.DeleteAppointment(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &appointmentpb.DeleteAppointmentResponse{
		Success: true,
		Message: "Appointment deleted successfully",
	}, nil
}

func (h *AppointmentHandler) ListAppointments(ctx context.Context, req *appointmentpb.ListAppointmentsRequest) (*appointmentpb.ListAppointmentsResponse, error) {
	if err := h.Validator.ValidateListAppointmentsRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	page := utils.GetPageFromRequest(req.Page)
	size := utils.GetSizeFromRequest(req.Size)

	appointments, total, err := h.AppointmentUsecase.ListAppointments(ctx, page, size)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToListAppointmentsResponse(appointments, total), nil
} 