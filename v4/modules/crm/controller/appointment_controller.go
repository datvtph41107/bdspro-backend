package controller

import (
	"context"
	"net/http"
	"time"

	dashboard "base/dashboard"
	_dto "common/domain/dto"
	_utils "common/utils"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"userv4/internal/httpapi"
	domain "userv4/modules/crm/entity"
	"userv4/modules/crm/mapper"
	"userv4/modules/crm/services"
	"userv4/modules/crm/validator"
)

// AppointmentController owns /v2.1/appointment/*, ported from crm-service AppointmentHandler.
// GET /v2.1/appointment/list/by-product/{id} stays in modules/bdspro.
type AppointmentController struct {
	usecase     services.AppointmentUsecase
	transformer mapper.AppointmentTransformer
	validator   validator.AppointmentValidator
	stats       dashboard.DashboardStatsJob
	sync        *_utils.SyncUtil
}

func NewAppointmentController(
	usecase services.AppointmentUsecase,
	transformer mapper.AppointmentTransformer,
	validator validator.AppointmentValidator,
	stats dashboard.DashboardStatsJob,
	sync *_utils.SyncUtil,
) *AppointmentController {
	return &AppointmentController{
		usecase:     usecase,
		transformer: transformer,
		validator:   validator,
		stats:       stats,
		sync:        sync,
	}
}

func (c *AppointmentController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.HandleProto(mux, http.MethodPost, "/v2.1/appointment", c.CreateAppointment)
	httpapi.HandleProto(mux, http.MethodGet, "/v2.1/appointment", c.ListAppointments)
	httpapi.HandleProto(mux, http.MethodGet, "/v2.1/appointment/admin/list", c.AdminListAppointments)
	httpapi.HandleProto(mux, http.MethodGet, "/v2.1/appointment/stats", c.GetAppointmentStats)
	httpapi.HandleProto(mux, http.MethodGet, "/v2.1/appointment/list/by-deal/{id}", c.GetAppointmentsByDeal)
	httpapi.HandleProto(mux, http.MethodGet, "/v2.1/appointment/{id}", c.GetAppointment)
	httpapi.HandleProto(mux, http.MethodPut, "/v2.1/appointment/{id}", c.UpdateAppointment)
	httpapi.HandleProto(mux, http.MethodDelete, "/v2.1/appointment/{id}", c.DeleteAppointment)
	httpapi.HandleProto(mux, http.MethodPut, "/v2.1/appointment/{id}/status/{action}", c.UpdateAppointmentStatus)
}

func invalidArgument(err error) error {
	return status.Error(codes.InvalidArgument, err.Error())
}

func (c *AppointmentController) CreateAppointment(ctx context.Context, req *crmpb.CreateAppointmentRequest) (*crmpb.Appointment, error) {
	if err := c.validator.ValidateCreateAppointmentRequest(req); err != nil {
		return nil, invalidArgument(err)
	}
	created, err := c.usecase.CreateAppointment(ctx, c.transformer.CreateAppointmentRequestToEntity(req))
	if err != nil {
		return nil, err
	}
	productIDs := req.ProductIds
	if created != nil && len(created.ProductIds) > 0 {
		productIDs = created.ProductIds
	}
	c.bumpListSync(ctx, productIDs, participantIDs(created)...)
	return c.transformer.EntityToAppointmentResponse(created, nil, nil, nil, nil), nil
}

func (c *AppointmentController) GetAppointment(ctx context.Context, req *crmpb.GetAppointmentRequest) (*crmpb.Appointment, error) {
	if err := c.validator.ValidateGetAppointmentRequest(req); err != nil {
		return nil, invalidArgument(err)
	}
	a, err := c.usecase.GetAppointment(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return c.transformer.EntityToAppointmentResponse(a.Appointment, a.Deal, a.Products, a.Participants, a.DealContract), nil
}

func (c *AppointmentController) UpdateAppointment(ctx context.Context, req *crmpb.UpdateAppointmentRequest) (*crmpb.Appointment, error) {
	if err := c.validator.ValidateUpdateAppointmentRequest(req); err != nil {
		return nil, invalidArgument(err)
	}
	updated, err := c.usecase.UpdateAppointment(ctx, c.transformer.UpdateAppointmentRequestToEntity(req))
	if err != nil {
		return nil, err
	}
	c.bumpListSync(ctx, nil, participantIDs(updated)...)
	return c.transformer.EntityToAppointmentResponse(updated, nil, nil, nil, nil), nil
}

func (c *AppointmentController) DeleteAppointment(ctx context.Context, req *crmpb.DeleteAppointmentRequest) (*crmpb.DeleteAppointmentResponse, error) {
	if err := c.validator.ValidateDeleteAppointmentRequest(req); err != nil {
		return nil, invalidArgument(err)
	}
	existing, err := c.usecase.GetAppointment(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if err := c.usecase.DeleteAppointment(ctx, req.Id); err != nil {
		return nil, err
	}
	if existing != nil {
		c.bumpListSync(ctx, nil, participantIDs(existing.Appointment)...)
	} else {
		c.bumpListSync(ctx, nil)
	}
	return &crmpb.DeleteAppointmentResponse{Id: req.Id}, nil
}

func (c *AppointmentController) ListAppointments(ctx context.Context, req *crmpb.ListAppointmentsRequest) (*crmpb.ListAppointmentsResponse, error) {
	if err := c.validator.ValidateListAppointmentsRequest(req); err != nil {
		return nil, invalidArgument(err)
	}

	page := int32(0)
	if req.Page != nil {
		page = *req.Page
	}
	syncKey := ""
	if c.sync != nil {
		syncKey = c.sync.GetKey(ctx, _utils.SyncKeyAppointmentMe, 0)
		if !c.sync.HasUpdated(ctx, syncKey, req.Timestamp) && page == 0 {
			return &crmpb.ListAppointmentsResponse{Timestamp: c.sync.GetTimestamp(ctx, syncKey)}, nil
		}
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

	filter := make(map[string]any)
	title := req.Title
	if title == "" && req.Text != "" {
		title = req.Text
	}
	if title != "" {
		filter["title"] = title
	}
	if req.Status != 0 {
		filter["status"] = req.Status
	}
	if len(req.DealIds) > 0 {
		dealIDs := make([]uint32, len(req.DealIds))
		copy(dealIDs, req.DealIds)
		filter["deal_ids"] = dealIDs
	}
	if req.TransactionId != nil && *req.TransactionId > 0 {
		filter["transaction_id"] = *req.TransactionId
	}
	if len(req.ParticipantIds) > 0 {
		filter["participant_ids"] = _utils.ParseInt64Array(req.ParticipantIds)
	}
	if req.FromDate != nil && *req.FromDate != "" {
		filter["from_date"] = *req.FromDate
	}
	if req.ToDate != nil && *req.ToDate != "" {
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
		if statuses := _utils.ParseInt64Array(*req.StatusFilters); len(statuses) > 0 {
			filter["status_filters"] = statuses
		}
	}

	appointments, total, err := c.usecase.ListAppointments(ctx, filter, pagable)
	if err != nil {
		return nil, err
	}

	// A full page-0 response refreshes the per-user sync timestamp returned to the client.
	ts := int64(0)
	if c.sync != nil && syncKey != "" && page == 0 {
		ts = time.Now().UnixMilli()
		_ = c.sync.PutTimestamp(ctx, syncKey, ts)
	}
	resp := c.transformer.EntityToListAppointmentsResponse(appointments, int32(total))
	resp.Timestamp = ts
	return resp, nil
}

func (c *AppointmentController) AdminListAppointments(ctx context.Context, req *crmpb.AdminListAppointmentsRequest) (*crmpb.AdminListAppointmentsResponse, error) {
	if err := c.validator.ValidateAdminListAppointmentsRequest(req); err != nil {
		return nil, invalidArgument(err)
	}
	var pagable _dto.Pagable
	if req.Page != nil && req.Size != nil {
		pagable = _dto.Pagable{Page: uint32(*req.Page), Size: uint32(*req.Size)}
	}
	if req.Sort != nil && *req.Sort != "" {
		pagable.Sort = *req.Sort
	}
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
	appointments, total, err := c.usecase.ListAppointments(ctx, filter, pagable)
	if err != nil {
		return nil, err
	}
	return c.transformer.EntityToAdminListAppointmentsResponse(appointments, int32(total)), nil
}

func (c *AppointmentController) GetAppointmentStats(ctx context.Context, req *sharepb.GetStatsRequest) (*sharepb.GetStatsResponse, error) {
	firstDate, err := time.Parse(time.DateOnly, req.FirstDate)
	if err != nil {
		return nil, invalidArgument(err)
	}
	secondDate, err := time.Parse(time.DateOnly, req.SecondDate)
	if err != nil {
		return nil, invalidArgument(err)
	}
	first, second, err := c.stats.GetStats(ctx, firstDate, secondDate)
	if err != nil {
		return nil, err
	}
	return &sharepb.GetStatsResponse{
		FirstDate:       first.CalculateTime.Format(time.DateOnly),
		SecondDate:      second.CalculateTime.Format(time.DateOnly),
		FirstDateCount:  uint32(first.Count),
		SecondDateCount: uint32(second.Count),
	}, nil
}

func (c *AppointmentController) GetAppointmentsByDeal(ctx context.Context, req *sharepb.IdRequest) (*crmpb.ListAppointmentsResponse, error) {
	appointments, total, err := c.usecase.GetAppointmentsByDeal(ctx, uint32(req.Id), 0, 100)
	if err != nil {
		return nil, err
	}
	return c.transformer.EntityToListAppointmentsResponse(appointments, int32(total)), nil
}

func (c *AppointmentController) UpdateAppointmentStatus(ctx context.Context, req *crmpb.UpdateAppointmentStatusRequest) (*crmpb.Appointment, error) {
	if err := c.validator.ValidateUpdateAppointmentStatusRequest(req); err != nil {
		return nil, invalidArgument(err)
	}
	updated, err := c.usecase.UpdateAppointmentStatus(ctx, req.Id, req.Action)
	if err != nil {
		return nil, err
	}
	c.bumpListSync(ctx, nil, participantIDs(updated)...)
	a, err := c.usecase.GetAppointment(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return c.transformer.EntityToAppointmentResponse(updated, a.Deal, a.Products, a.Participants, a.DealContract), nil
}

func participantIDs(appointment *domain.Appointment) []uint64 {
	if appointment == nil {
		return nil
	}
	ids := make([]uint64, 0, len(appointment.ParticipantIds))
	for _, id := range appointment.ParticipantIds {
		if id > 0 {
			ids = append(ids, uint64(id))
		}
	}
	return ids
}

func (c *AppointmentController) bumpListSync(ctx context.Context, productIDs []int64, participantUserIDs ...uint64) {
	if c.sync == nil {
		return
	}
	now := time.Now().UnixMilli()
	_ = c.sync.PutTimestamp(ctx, c.sync.GetKey(ctx, _utils.SyncKeyAppointmentMe, 0), now)
	profileID := _utils.GetProfileIdWithContext(ctx)
	for _, userID := range participantUserIDs {
		if userID == 0 || userID == profileID {
			continue
		}
		_ = c.sync.PutTimestamp(ctx, c.sync.GetKeyForProfile(userID, 0, _utils.SyncKeyAppointmentMe), now)
	}
	seen := make(map[uint64]struct{}, len(productIDs))
	for _, productID := range productIDs {
		if productID <= 0 {
			continue
		}
		id := uint64(productID)
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		_ = c.sync.PutTimestamp(ctx, c.sync.GetKey(ctx, _utils.SyncKeyProductAppointments, id), now)
	}
}
