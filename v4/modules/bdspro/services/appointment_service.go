package services

import (
	"context"
	"fmt"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	"mdn4/internal/interface/providers"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/repository"
	bdspropb "pb/types/bdspro"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"github.com/lib/pq"
)

type AppointmentService interface {
	GetByID(ctx context.Context, id uint64) (*crmpb.Appointment, error)
	ListMe(ctx context.Context, req dto.AppointmentListRequest) (*crmpb.ListAppointmentsResponse, error)
	ListByProduct(ctx context.Context, productID uint64, req dto.SyncByIDRequest) (*crmpb.ListAppointmentsResponse, error)
	CreateAppointment(ctx context.Context, title, description, place, startTime, endTime string, participantIDs []int64, dealID uint32, productIDs []int64, status uint32) (uint64, error)
	DeleteAppointment(ctx context.Context, id uint64) error
}

// @bind: mdn4/modules/bdspro/services.AppointmentService
type LocalAppointmentService struct {
	sync     *_utils.SyncUtil
	repo     repository.AppointmentRepository
	products repository.ProductRepository
	deals    repository.DealRepository
	profile  providers.ProfileProvider
}

func NewLocalAppointmentService(
	sync *_utils.SyncUtil,
	repo repository.AppointmentRepository,
	products repository.ProductRepository,
	deals repository.DealRepository,
	profile providers.ProfileProvider,
) *LocalAppointmentService {
	return &LocalAppointmentService{
		sync:     sync,
		repo:     repo,
		products: products,
		deals:    deals,
		profile:  profile,
	}
}

func (s *LocalAppointmentService) GetByID(ctx context.Context, id uint64) (*crmpb.Appointment, error) {
	if id == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	appt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if appt == nil {
		return nil, _errors.ReturnError(404, "appointment not found")
	}
	data, err := s.mapAppointmentsEnriched(ctx, []*entity.Appointment{appt})
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, _errors.ReturnError(404, "appointment not found")
	}
	return data[0], nil
}

func (s *LocalAppointmentService) ListMe(ctx context.Context, req dto.AppointmentListRequest) (*crmpb.ListAppointmentsResponse, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ReturnError(401, "profileId not found in context")
	}

	var syncKey string
	if s.sync != nil {
		syncKey = s.sync.GetKey(ctx, _utils.SyncKeyAppointmentMe, 0)
		if !s.sync.HasUpdated(ctx, syncKey, req.Timestamp) && req.Page == 0 {
			return &crmpb.ListAppointmentsResponse{
				Timestamp: s.sync.GetTimestamp(ctx, syncKey),
			}, nil
		}
	}

	list, total, err := s.repo.ListForMe(ctx, profileID, req)
	if err != nil {
		return nil, err
	}

	data, err := s.mapAppointmentsEnriched(ctx, list)
	if err != nil {
		return nil, err
	}

	resp := &crmpb.ListAppointmentsResponse{
		Data:  data,
		Total: int32(total),
	}
	if s.sync != nil && syncKey != "" && req.Page == 0 {
		now := time.Now().UnixMilli()
		_ = s.sync.PutTimestamp(ctx, syncKey, now)
		resp.Timestamp = now
	}
	return resp, nil
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

// CreateAppointment creates an appointment and inserts participants.
// startTime and endTime are RFC3339 strings. Host (createdBy) is injected from JWT.
func (s *LocalAppointmentService) CreateAppointment(
	ctx context.Context,
	title, description, place, startTime, endTime string,
	participantIDs []int64,
	dealID uint32,
	productIDs []int64,
	status uint32,
) (uint64, error) {
	hostID := _utils.GetProfileIdWithContext(ctx)

	appt := &entity.Appointment{
		Title:       title,
		Description: description,
		Place:       place,
		DealID:      uint64(dealID),
		Status:      int32(status),
		ProductIDs:  pq.Int64Array(append([]int64(nil), productIDs...)),
	}
	if hostID > 0 {
		appt.CreatedBy = &hostID
	}
	if startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			appt.StartTime = &t
		}
	}
	if endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			appt.EndTime = &t
		}
	}

	// Ensure host is among participants.
	hostPresent := false
	for _, pid := range participantIDs {
		if uint64(pid) == hostID {
			hostPresent = true
			break
		}
	}
	merged := append([]int64(nil), participantIDs...)
	if !hostPresent && hostID > 0 {
		merged = append(merged, int64(hostID))
	}
	appt.ParticipantIDs = pq.Int64Array(merged)

	if err := s.repo.Create(ctx, appt); err != nil {
		return 0, err
	}

	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyAppointmentMe, 0)
		_ = s.sync.PutTimestamp(ctx, key, 0)
		for _, pid := range productIDs {
			pKey := s.sync.GetKey(ctx, _utils.SyncKeyProductAppointments, uint64(pid))
			_ = s.sync.PutTimestamp(ctx, pKey, 0)
		}
	}

	return appt.ID, nil
}

// DeleteAppointment soft-deletes an appointment by id.
func (s *LocalAppointmentService) DeleteAppointment(ctx context.Context, id uint64) error {
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return err
	}
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyAppointmentMe, 0)
		_ = s.sync.PutTimestamp(ctx, key, 0)
	}
	return nil
}

func (s *LocalAppointmentService) mapAppointmentsEnriched(ctx context.Context, list []*entity.Appointment) ([]*crmpb.Appointment, error) {
	if len(list) == 0 {
		return []*crmpb.Appointment{}, nil
	}

	productIDsSet := map[uint64]struct{}{}
	dealIDsSet := map[uint64]struct{}{}
	profileIDsSet := map[uint64]struct{}{}
	for _, a := range list {
		for _, id := range a.ProductIDs {
			if id > 0 {
				productIDsSet[uint64(id)] = struct{}{}
			}
		}
		if a.DealID > 0 {
			dealIDsSet[a.DealID] = struct{}{}
		}
		for _, p := range a.Participants {
			if p.UserID > 0 {
				profileIDsSet[p.UserID] = struct{}{}
			}
		}
	}

	productsMap := map[uint64]*bdspropb.ProductAttachment{}
	if s.products != nil && len(productIDsSet) > 0 {
		ids := make([]uint64, 0, len(productIDsSet))
		for id := range productIDsSet {
			ids = append(ids, id)
		}
		rows, err := s.products.ListByIDs(ctx, ids)
		if err != nil {
			return nil, fmt.Errorf("list products: %w", err)
		}
		for _, p := range rows {
			if p == nil {
				continue
			}
			att := &bdspropb.ProductAttachment{
				Id:         p.ID,
				Name:       p.Name,
				Code:       p.Code,
				Area:       p.Area,
				RentStatus: p.RentStatus,
				SaleStatus: p.SaleStatus,
				SourceType: int32(p.SourceType),
			}
			if p.Price != nil && p.Price.SalePrice != nil {
				att.Price = *p.Price.SalePrice
			}
			if p.AssetID != nil {
				att.AssetId = p.AssetID
			}
			productsMap[p.ID] = att
		}
	}

	dealsMap := map[uint64]*bdspropb.GroupDeal{}
	if s.deals != nil && len(dealIDsSet) > 0 {
		ids := make([]uint64, 0, len(dealIDsSet))
		for id := range dealIDsSet {
			ids = append(ids, id)
		}
		rows, err := s.deals.ListByIDs(ctx, ids)
		if err != nil {
			return nil, fmt.Errorf("list deals: %w", err)
		}
		for _, d := range rows {
			if d == nil {
				continue
			}
			dealsMap[d.ID] = mapDealToPB(d)
		}
	}

	profilesMap := map[uint64]*sharepb.ProfileItem{}
	if s.profile != nil && len(profileIDsSet) > 0 {
		ids := make([]uint64, 0, len(profileIDsSet))
		for id := range profileIDsSet {
			ids = append(ids, id)
		}
		if m, err := s.profile.GetMapProfileByIds(ctx, ids); err == nil {
			for id, p := range m {
				if p == nil {
					continue
				}
				profilesMap[id] = &sharepb.ProfileItem{
					Id:       id,
					FullName: p.FullName,
					Avatar:   p.Avatar,
					Phone:    p.Phone,
				}
			}
		}
	}

	out := make([]*crmpb.Appointment, 0, len(list))
	for _, a := range list {
		pb := mapAppointmentToPB(a)
		if a.DealID > 0 {
			pb.Deal = dealsMap[a.DealID]
		}
		for _, pid := range a.ProductIDs {
			if att := productsMap[uint64(pid)]; att != nil {
				pb.Products = append(pb.Products, att)
			}
		}
		participants := make([]*sharepb.ProfileItem, 0, len(a.Participants))
		for _, p := range a.Participants {
			if item := profilesMap[p.UserID]; item != nil {
				participants = append(participants, item)
			} else if p.UserID > 0 {
				participants = append(participants, &sharepb.ProfileItem{Id: p.UserID})
			}
		}
		pb.Participants = participants
		out = append(out, pb)
	}
	return out, nil
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
