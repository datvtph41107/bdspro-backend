package transformer

import (
	appointmentpb "pb/types/appointment"
	pb_bdspro "pb/types/bdspro"
	organizationpb "pb/types/organization"
	"time"

	"github.com/tuandq2112/appointment-service/internal/domain/entity"
	"github.com/tuandq2112/appointment-service/internal/dto"
)

type AppointmentTransformer interface {
	CreateAppointmentRequestToEntity(req *appointmentpb.CreateAppointmentRequest) *entity.Appointment
	UpdateAppointmentRequestToEntity(req *appointmentpb.UpdateAppointmentRequest) *entity.Appointment
	EntityToAppointmentResponse(appointment *entity.Appointment, deal *organizationpb.GroupDeal, products []*pb_bdspro.ProductAttachment, dealTransactions []*organizationpb.DealTransactionResponse, participants []*organizationpb.DealInvitationMember) *appointmentpb.Appointment
	EntityToListAppointmentsResponse(appointments []*dto.AppointmentDTO, total int32) *appointmentpb.ListAppointmentsResponse
}

type appointmentTransformer struct{}

func NewAppointmentTransformer() AppointmentTransformer {
	return &appointmentTransformer{}
}

func (t *appointmentTransformer) CreateAppointmentRequestToEntity(req *appointmentpb.CreateAppointmentRequest) *entity.Appointment {
	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)

	return &entity.Appointment{
		Title:              req.Title,
		Description:        req.Description,
		StartTime:          startTime,
		EndTime:            endTime,
		Status:             req.Status,
		Place:              req.Place,
		ParticipantIds:     req.ParticipantIds,
		DealId:             req.DealId,
		ProductIds:         req.ProductIds,
		ReminderSchedule:   req.ReminderSchedule,
		TransactionIds:     req.TransactionIds,
	}
}

func (t *appointmentTransformer) UpdateAppointmentRequestToEntity(req *appointmentpb.UpdateAppointmentRequest) *entity.Appointment {
	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)
	return &entity.Appointment{
		Id:                 req.Id,
		Title:              req.Title,
		Description:        req.Description,
		StartTime:          startTime,
		EndTime:            endTime,
		Status:             req.Status,
		Place:              req.Place,
		ParticipantIds:     req.ParticipantIds,
		DealId:             req.DealId,
		ProductIds:         req.ProductIds,
		ReminderSchedule:   req.ReminderSchedule,
		TransactionIds:     req.TransactionIds,
	}
}

func (t *appointmentTransformer) EntityToAppointmentResponse(appointment *entity.Appointment, deal *organizationpb.GroupDeal, products []*pb_bdspro.ProductAttachment, dealTransactions []*organizationpb.DealTransactionResponse, participants []*organizationpb.DealInvitationMember) *appointmentpb.Appointment {
	return &appointmentpb.Appointment{
		Id:                 appointment.Id,
		Title:              appointment.Title,
		Description:        appointment.Description,
		StartTime:          appointment.StartTime.Format(time.RFC3339),
		EndTime:            appointment.EndTime.Format(time.RFC3339),
		Status:             appointment.Status,
		Place:              appointment.Place,
		ParticipantIds:     appointment.ParticipantIds,
		DealId:             appointment.DealId,
		ProductIds:         appointment.ProductIds,
		ReminderSchedule:   appointment.ReminderSchedule,
		TransactionIds:     appointment.TransactionIds,
		Deal:               deal,
		Products:           products,
		DealTransactions:   dealTransactions,
		Participants:       participants,
	}
}

func (t *appointmentTransformer) EntityToListAppointmentsResponse(appointments []*dto.AppointmentDTO, total int32) *appointmentpb.ListAppointmentsResponse {
	appointmentsResponse := make([]*appointmentpb.Appointment, len(appointments))
	for i, appointment := range appointments {
		appointmentsResponse[i] = t.EntityToAppointmentResponse(appointment.Appointment, appointment.Deal, appointment.Products, appointment.DealTransactions, appointment.Participants)
	}
	return &appointmentpb.ListAppointmentsResponse{
		Data:  appointmentsResponse,
		Total: total,
	}
}
	