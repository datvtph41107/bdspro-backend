package transformer

import (
	appointmentpb "pb/types/appointment"

	"github.com/tuandq2112/appointment/internal/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AppointmentTransformer interface {
	CreateAppointmentRequestToEntity(req *appointmentpb.CreateAppointmentRequest) *entity.Appointment
	UpdateAppointmentRequestToEntity(req *appointmentpb.UpdateAppointmentRequest) *entity.Appointment
	EntityToAppointmentResponse(appointment *entity.Appointment) *appointmentpb.Appointment
	EntityToListAppointmentsResponse(appointments []*entity.Appointment, total int32) *appointmentpb.ListAppointmentsResponse
}

type appointmentTransformer struct{}

func NewAppointmentTransformer() AppointmentTransformer {
	return &appointmentTransformer{}
}

func (t *appointmentTransformer) CreateAppointmentRequestToEntity(req *appointmentpb.CreateAppointmentRequest) *entity.Appointment {
	return &entity.Appointment{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime.AsTime(),
		EndTime:     req.EndTime.AsTime(),
		Status:      req.Status,
	}
}

func (t *appointmentTransformer) UpdateAppointmentRequestToEntity(req *appointmentpb.UpdateAppointmentRequest) *entity.Appointment {
	return &entity.Appointment{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime.AsTime(),
		EndTime:     req.EndTime.AsTime(),
		Status:      req.Status,
	}
}

func (t *appointmentTransformer) EntityToAppointmentResponse(appointment *entity.Appointment) *appointmentpb.Appointment {
	return &appointmentpb.Appointment{
		Id:          appointment.ID,
		UserId:      appointment.UserID,
		Title:       appointment.Title,
		Description: appointment.Description,
		StartTime:   timestamppb.New(appointment.StartTime),
		EndTime:     timestamppb.New(appointment.EndTime),
		Status:      appointment.Status,
		CreatedBy:   appointment.CreatedBy,
		UpdatedBy:   appointment.UpdatedBy,
		CreatedAt:   timestamppb.New(appointment.CreatedAt),
		UpdatedAt:   timestamppb.New(appointment.UpdatedAt),
	}
}

func (t *appointmentTransformer) EntityToListAppointmentsResponse(appointments []*entity.Appointment, total int32) *appointmentpb.ListAppointmentsResponse {
	appointmentsResponse := make([]*appointmentpb.Appointment, len(appointments))
	for i, appointment := range appointments {
		appointmentsResponse[i] = t.EntityToAppointmentResponse(appointment)
	}
	return &appointmentpb.ListAppointmentsResponse{
		Appointments: appointmentsResponse,
		Total:        total,
	}
} 