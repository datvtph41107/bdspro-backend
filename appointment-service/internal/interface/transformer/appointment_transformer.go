package transformer

import (
	appointmentpb "pb/types/appointment"
	bdspropb "pb/types/bdspro"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"
	"time"

	"appointment/internal/domain/entity"
	"appointment/internal/dto"
)

type AppointmentTransformer interface {
	CreateAppointmentRequestToEntity(req *appointmentpb.CreateAppointmentRequest) *entity.Appointment
	UpdateAppointmentRequestToEntity(req *appointmentpb.UpdateAppointmentRequest) *entity.Appointment
	EntityToAppointmentResponse(appointment *entity.Appointment, deal *organizationpb.GroupDeal, products []*bdspropb.ProductAttachment, participants []*sharepb.ProfileItem) *appointmentpb.Appointment
	EntityToListAppointmentsResponse(appointments []*dto.AppointmentDTO, total int32) *appointmentpb.ListAppointmentsResponse
	EntityToAdminListAppointmentsResponse(appointments []*dto.AppointmentDTO, total int32) *appointmentpb.AdminListAppointmentsResponse
}

type appointmentTransformer struct{}

func NewAppointmentTransformer() AppointmentTransformer {
	return &appointmentTransformer{}
}

func (t *appointmentTransformer) CreateAppointmentRequestToEntity(req *appointmentpb.CreateAppointmentRequest) *entity.Appointment {
	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)

	return &entity.Appointment{
		Title:            req.Title,
		Description:      req.Description,
		StartTime:        startTime,
		EndTime:          endTime,
		Mode:             req.Mode,
		Place:            req.Place,
		ParticipantIds:   req.ParticipantIds,
		DealId:           req.DealId,
		ProductIds:       req.ProductIds,
		ReminderSchedule: req.ReminderSchedule,
		TransactionIds:   req.TransactionIds,
		OptionExtend:     req.OptionExtend,
		Status:           req.Status,
	}
}

func (t *appointmentTransformer) UpdateAppointmentRequestToEntity(req *appointmentpb.UpdateAppointmentRequest) *entity.Appointment {
	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)
	return &entity.Appointment{
		Id:               req.Id,
		Title:            req.Title,
		Description:      req.Description,
		StartTime:        startTime,
		EndTime:          endTime,
		Mode:             req.Mode,
		Place:            req.Place,
		ParticipantIds:   req.ParticipantIds,
		DealId:           req.DealId,
		ProductIds:       req.ProductIds,
		ReminderSchedule: req.ReminderSchedule,
		OptionExtend:     req.OptionExtend,
		TransactionIds:   req.TransactionIds,
		Status:           req.Status,
	}
}

func (t *appointmentTransformer) EntityToAppointmentResponse(appointment *entity.Appointment, deal *organizationpb.GroupDeal, products []*bdspropb.ProductAttachment, participants []*sharepb.ProfileItem) *appointmentpb.Appointment {
	return &appointmentpb.Appointment{
		Id:               appointment.Id,
		Title:            appointment.Title,
		Description:      appointment.Description,
		StartTime:        appointment.StartTime.Format(time.RFC3339),
		EndTime:          appointment.EndTime.Format(time.RFC3339),
		Mode:             appointment.Mode,
		Place:            appointment.Place,
		ParticipantIds:   appointment.ParticipantIds,
		DealId:           appointment.DealId,
		ProductIds:       appointment.ProductIds,
		ReminderSchedule: appointment.ReminderSchedule,
		OptionExtend:     appointment.OptionExtend,
		Status:           appointment.Status,
		TransactionIds:   appointment.TransactionIds,
		Deal:             deal,
		Products:         products,
		// DealTransactions: dealTransactions,
		Participants: participants,
	}
}

func (t *appointmentTransformer) EntityToListAppointmentsResponse(appointments []*dto.AppointmentDTO, total int32) *appointmentpb.ListAppointmentsResponse {
	appointmentsResponse := make([]*appointmentpb.Appointment, len(appointments))
	for i, appointment := range appointments {
		appointmentsResponse[i] = t.EntityToAppointmentResponse(appointment.Appointment, appointment.Deal, appointment.Products, appointment.Participants)
	}
	return &appointmentpb.ListAppointmentsResponse{
		Data:  appointmentsResponse,
		Total: total,
	}
}

func (t *appointmentTransformer) EntityToAdminListAppointmentsResponse(appointments []*dto.AppointmentDTO, total int32) *appointmentpb.AdminListAppointmentsResponse {
	appointmentsResponse := make([]*appointmentpb.Appointment, len(appointments))
	for i, appointment := range appointments {
		appointmentsResponse[i] = t.EntityToAppointmentResponse(appointment.Appointment, appointment.Deal, appointment.Products, appointment.Participants)
	}
	return &appointmentpb.AdminListAppointmentsResponse{
		Data:  appointmentsResponse,
		Total: total,
	}
}
