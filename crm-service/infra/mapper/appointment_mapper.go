package mapper

import (
	bdspropb "pb/types/bdspro"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
	"time"

	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
)

type AppointmentTransformer interface {
	CreateAppointmentRequestToEntity(req *crmpb.CreateAppointmentRequest) *domain.Appointment
	UpdateAppointmentRequestToEntity(req *crmpb.UpdateAppointmentRequest) *domain.Appointment
	EntityToAppointmentResponse(appointment *domain.Appointment, deal *bdspropb.GroupDeal, products []*bdspropb.ProductAttachment, participants []*sharepb.ProfileItem, dealContract *bdspropb.DealContractItem) *crmpb.Appointment
	EntityToListAppointmentsResponse(appointments []*dto.AppointmentDTO, total int32) *crmpb.ListAppointmentsResponse
	EntityToAdminListAppointmentsResponse(appointments []*dto.AppointmentDTO, total int32) *crmpb.AdminListAppointmentsResponse
}

type appointmentTransformer struct{}

func NewAppointmentTransformer() AppointmentTransformer {
	return &appointmentTransformer{}
}

func (t *appointmentTransformer) CreateAppointmentRequestToEntity(req *crmpb.CreateAppointmentRequest) *domain.Appointment {
	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)

	// Set default status to "Chờ xác nhận" if not provided
	status := enums.EAppointmentStatus(req.Status)
	if status == 0 {
		status = enums.AppointmentStatusPending
	}

	// Convert transactionSteps from []uint32 to pq.Int32Array
	transactionSteps := make([]int32, len(req.TransactionSteps))
	for i, step := range req.TransactionSteps {
		transactionSteps[i] = int32(step)
	}

	// Tạo AppointmentParticipants từ participantIds
	participants := make([]domain.AppointmentParticipant, 0, len(req.ParticipantIds))
	for _, participantId := range req.ParticipantIds {
		participants = append(participants, domain.AppointmentParticipant{
			UserId: uint32(participantId),
			Role:   domain.AppointmentParticipantRoleGuest,
			Status: domain.AppointmentParticipantStatusPending,
		})
	}

	return &domain.Appointment{
		Title:                   req.Title,
		Description:             req.Description,
		StartTime:               startTime,
		EndTime:                 endTime,
		Mode:                    req.Mode,
		Place:                   req.Place,
		ParticipantIds:          req.ParticipantIds,
		AppointmentParticipants: participants,
		DealId:                  req.DealId,
		ProductIds:              req.ProductIds,
		ReminderSchedule:        req.ReminderSchedule,
		TransactionId:           req.TransactionId,
		TransactionSteps:        transactionSteps,
		OptionExtend:            req.OptionExtend,
		Status:                  status,
	}
}

func (t *appointmentTransformer) UpdateAppointmentRequestToEntity(req *crmpb.UpdateAppointmentRequest) *domain.Appointment {
	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)

	// Tạo AppointmentParticipants từ participantIds
	participants := make([]domain.AppointmentParticipant, 0, len(req.ParticipantIds))
	for _, participantId := range req.ParticipantIds {
		participants = append(participants, domain.AppointmentParticipant{
			UserId: uint32(participantId),
			Role:   domain.AppointmentParticipantRoleGuest,
			Status: domain.AppointmentParticipantStatusPending,
		})
	}

	appointment := &domain.Appointment{
		Id:                      req.Id,
		Title:                   req.Title,
		Description:             req.Description,
		StartTime:               startTime,
		EndTime:                 endTime,
		Mode:                    req.Mode,
		Place:                   req.Place,
		ParticipantIds:          req.ParticipantIds,
		AppointmentParticipants: participants,
		DealId:                  req.DealId,
		ProductIds:              req.ProductIds,
		ReminderSchedule:        req.ReminderSchedule,
		OptionExtend:            req.OptionExtend,
		TransactionId:           req.TransactionId,
	}

	// Convert transactionSteps from []uint32 to pq.Int32Array if provided
	if len(req.TransactionSteps) > 0 {
		transactionSteps := make([]int32, len(req.TransactionSteps))
		for i, step := range req.TransactionSteps {
			transactionSteps[i] = int32(step)
		}
		appointment.TransactionSteps = transactionSteps
	}

	// Only update status if provided (non-zero)
	if req.Status != 0 {
		appointment.Status = enums.EAppointmentStatus(req.Status)
	}

	return appointment
}

func (t *appointmentTransformer) EntityToAppointmentResponse(appointment *domain.Appointment, deal *bdspropb.GroupDeal, products []*bdspropb.ProductAttachment, participants []*sharepb.ProfileItem, dealContract *bdspropb.DealContractItem) *crmpb.Appointment {
	// Convert transactionSteps from pq.Int32Array to []uint32
	transactionSteps := make([]uint32, len(appointment.TransactionSteps))
	for i, step := range appointment.TransactionSteps {
		transactionSteps[i] = uint32(step)
	}

	return &crmpb.Appointment{
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
		Status:           uint32(appointment.Status),
		TransactionId:    appointment.TransactionId,
		TransactionSteps: transactionSteps,
		Deal:             deal,
		Products:         products,
		DealContract:     dealContract,
		// DealTransactions: dealTransactions,
		Participants: participants,
	}
}

func (t *appointmentTransformer) EntityToListAppointmentsResponse(appointments []*dto.AppointmentDTO, total int32) *crmpb.ListAppointmentsResponse {
	appointmentsResponse := make([]*crmpb.Appointment, len(appointments))
	for i, appointment := range appointments {
		appointmentsResponse[i] = t.EntityToAppointmentResponse(appointment.Appointment, appointment.Deal, appointment.Products, appointment.Participants, appointment.DealContract)
	}
	return &crmpb.ListAppointmentsResponse{
		Data:  appointmentsResponse,
		Total: total,
	}
}

func (t *appointmentTransformer) EntityToAdminListAppointmentsResponse(appointments []*dto.AppointmentDTO, total int32) *crmpb.AdminListAppointmentsResponse {
	appointmentsResponse := make([]*crmpb.Appointment, len(appointments))
	for i, appointment := range appointments {
		appointmentsResponse[i] = t.EntityToAppointmentResponse(appointment.Appointment, appointment.Deal, appointment.Products, appointment.Participants, appointment.DealContract)
	}
	return &crmpb.AdminListAppointmentsResponse{
		Data:  appointmentsResponse,
		Total: total,
	}
}
