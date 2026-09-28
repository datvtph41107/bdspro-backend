package services

import (
	"context"
	"time"

	_utils "common/utils"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/repository"
)

type AppointmentService interface {
	ListByProduct(ctx context.Context, productID uint64, req dto.SyncByIDRequest) (*crmpb.ListAppointmentsResponse, error)
}

type LocalAppointmentService struct {
	sync *_utils.SyncUtil
	repo repository.AppointmentRepository
}

func NewLocalAppointmentService(sync *_utils.SyncUtil, repo repository.AppointmentRepository) *LocalAppointmentService {
	return &LocalAppointmentService{sync: sync, repo: repo}
}

func (s *LocalAppointmentService) ListByProduct(ctx context.Context, productID uint64, req dto.SyncByIDRequest) (*crmpb.ListAppointmentsResponse, error) {
	if productID == 0 {
		return &crmpb.ListAppointmentsResponse{}, nil
	}
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyProductAppointments, productID)
		if !s.sync.HasUpdated(ctx, key, req.Timestamp) && req.Page == 0 && req.Timestamp > 0 {
			return &crmpb.ListAppointmentsResponse{}, nil
		}
	}
	list, total, err := s.repo.ListByProductID(ctx, productID, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	data := make([]*crmpb.Appointment, 0, len(list))
	for _, a := range list {
		data = append(data, mapAppointmentToPB(a))
	}
	return &crmpb.ListAppointmentsResponse{Data: data, Total: int32(total)}, nil
}

func mapAppointmentToPB(a *entity.Appointment) *crmpb.Appointment {
	if a == nil {
		return nil
	}
	out := &crmpb.Appointment{
		Id:               uint32(a.ID),
		Title:            a.Title,
		Place:            a.Place,
		DealId:           uint32(a.DealID),
		Description:      a.Description,
		ReminderSchedule: uint32(a.ReminderSchedule),
		Mode:             uint32(a.Mode),
		OptionExtend:     uint32(a.OptionExtend),
		Status:           uint32(a.Status),
		TransactionId:    a.TransactionID,
		ParticipantIds:   append([]int64(nil), a.ParticipantIDs...),
		ProductIds:       append([]int64(nil), a.ProductIDs...),
	}
	if a.StartTime != nil {
		out.StartTime = a.StartTime.Format(time.RFC3339)
	}
	if a.EndTime != nil {
		out.EndTime = a.EndTime.Format(time.RFC3339)
	}
	for _, step := range a.TransactionSteps {
		out.TransactionSteps = append(out.TransactionSteps, uint32(step))
	}
	for _, p := range a.Participants {
		if p.UserID > 0 {
			out.Participants = append(out.Participants, &sharepb.ProfileItem{Id: p.UserID})
		}
	}
	return out
}
