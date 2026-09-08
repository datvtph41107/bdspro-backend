package handler

import (
	"context"
	"crm/infra/client"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"
	"fmt"
	"strings"
	"time"

	_enum "common/domain/enum"
	_errors "common/errors"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	permSalesView = []uint32{
		_enum.ADMIN_KD_XEM,
		_enum.ADMIN_KD_TAO,
		_enum.ADMIN_KD_SUA,
		_enum.ADMIN_KD_XOA,
	}
	permSalesCreate = []uint32{_enum.ADMIN_KD_TAO}
	permSalesUpdate = []uint32{_enum.ADMIN_KD_SUA}
	permSalesDelete = []uint32{_enum.ADMIN_KD_XOA}
)

type AdminSalesHandler struct {
	crmpb.UnimplementedAdminSalesServiceServer
	usecase    *usecase.AdminSalesUsecase
	authClient *client.AuthClient
}

// @bind: crm/infra/handler.AdminSalesHandler
func NewAdminSalesHandler(uc *usecase.AdminSalesUsecase, authClient *client.AuthClient) *AdminSalesHandler {
	return &AdminSalesHandler{usecase: uc, authClient: authClient}
}

func (h *AdminSalesHandler) ListOpportunities(ctx context.Context, req *crmpb.ListAdminOpportunitiesRequest) (*crmpb.ListAdminOpportunitiesResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesView); err != nil {
		return nil, err
	}
	search := mapListSearch(req)
	items, total, err := h.usecase.List(ctx, search)
	if err != nil {
		return nil, err
	}
	out := make([]*crmpb.AdminOpportunityItem, 0, len(items))
	for i := range items {
		out = append(out, mapAdminLeadToItem(&items[i]))
	}
	resp := &crmpb.ListAdminOpportunitiesResponse{Data: out, Total: total}
	if strings.EqualFold(req.GetViewMode(), "funnel") {
		funnel, err := h.usecase.Funnel(ctx, search)
		if err != nil {
			return nil, err
		}
		resp.Funnel = make([]*crmpb.AdminFunnelBucket, 0, len(funnel))
		for _, b := range funnel {
			resp.Funnel = append(resp.Funnel, &crmpb.AdminFunnelBucket{
				Status:           int32(b.Status),
				Label:            b.Label,
				Count:            b.Count,
				ExpectedValueSum: b.ExpectedValueSum,
			})
		}
	}
	return resp, nil
}

func (h *AdminSalesHandler) GetOpportunitiesSummary(ctx context.Context, _ *emptypb.Empty) (*crmpb.AdminOpportunitiesSummaryResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesView); err != nil {
		return nil, err
	}
	s, err := h.usecase.Summary(ctx)
	if err != nil {
		return nil, err
	}
	return &crmpb.AdminOpportunitiesSummaryResponse{
		Total:               s.Total,
		OverdueFollowUp:     s.OverdueFollowUp,
		Unassigned:          s.Unassigned,
		WaitingQuote:        s.WaitingQuote,
		WaitingApproval:     s.WaitingApproval,
		WaitingPayment:      s.WaitingPayment,
		Renewal:             s.Renewal,
		Upgrade:             s.Upgrade,
		ChurnRisk:           s.ChurnRisk,
		StatusNew:           s.StatusNew,
		Consulting:          s.Consulting,
		Won:                 s.Won,
		Lost:                s.Lost,
		ExpectedRevenue:     s.ExpectedRevenue,
		WaitingPaymentValue: s.WaitingPaymentValue,
		WonValue:            s.WonValue,
	}, nil
}

func (h *AdminSalesHandler) ListContacts(ctx context.Context, req *crmpb.ListAdminContactsRequest) (*crmpb.ListAdminContactsResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesView); err != nil {
		return nil, err
	}
	items, total, err := h.usecase.ListContacts(ctx, req.Page, req.Size, req.GetQ())
	if err != nil {
		return nil, err
	}
	out := make([]*crmpb.AdminContactItem, 0, len(items))
	for i := range items {
		out = append(out, mapAdminContactToItem(&items[i]))
	}
	return &crmpb.ListAdminContactsResponse{Data: out, Total: total}, nil
}

func (h *AdminSalesHandler) CreateContact(ctx context.Context, req *crmpb.CreateAdminContactRequest) (*crmpb.AdminContactResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesCreate); err != nil {
		return nil, err
	}
	contact, err := h.usecase.CreateContact(ctx, &dto.AdminContactSaveDTO{
		FullName: req.FullName,
		Phone:    req.Phone,
		Email:    req.Email,
		Company:  req.Company,
		Avatar:   req.Avatar,
		Address:  req.Address,
		Zalo:     req.Zalo,
		Note:     req.Note,
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.AdminContactResponse{Data: mapAdminContactToItem(contact)}, nil
}

func (h *AdminSalesHandler) GetOpportunity(ctx context.Context, req *sharepb.IdRequest) (*crmpb.AdminOpportunityResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesView); err != nil {
		return nil, err
	}
	lead, err := h.usecase.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	events, _ := h.usecase.ListEvents(ctx, req.Id)
	h.enrichSalesEventActors(ctx, events)
	return mapAdminLeadToResponse(lead, events), nil
}

func (h *AdminSalesHandler) GetOpportunityQuickView(ctx context.Context, req *sharepb.IdRequest) (*crmpb.AdminOpportunityQuickViewResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesView); err != nil {
		return nil, err
	}
	lead, err := h.usecase.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	events, _ := h.usecase.ListEvents(ctx, req.Id)
	h.enrichSalesEventActors(ctx, events)
	if len(events) > 5 {
		events = events[:5]
	}
	item := mapAdminLeadToItem(lead)
	resp := &crmpb.AdminOpportunityQuickViewResponse{
		Data:   item,
		Alerts: usecase.BuildOpportunityAlerts(lead),
	}
	if lead.Contact != nil {
		resp.Address = lead.Contact.Address
		resp.Zalo = lead.Contact.Zalo
	}
	resp.RecentEvents = mapEvents(events)
	return resp, nil
}

func (h *AdminSalesHandler) CreateOpportunity(ctx context.Context, req *crmpb.CreateAdminOpportunityRequest) (*crmpb.AdminOpportunityResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesCreate); err != nil {
		return nil, err
	}
	lead, err := h.usecase.Create(ctx, mapAdminCreateReq(req))
	if err != nil {
		return nil, err
	}
	events, _ := h.usecase.ListEvents(ctx, lead.ID)
	h.enrichSalesEventActors(ctx, events)
	return mapAdminLeadToResponse(lead, events), nil
}

func (h *AdminSalesHandler) UpdateOpportunity(ctx context.Context, req *crmpb.UpdateAdminOpportunityRequest) (*crmpb.AdminOpportunityResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesUpdate); err != nil {
		return nil, err
	}
	if req.Id == 0 {
		return nil, _errors.ReturnError(400, "ID không hợp lệ")
	}
	lead, err := h.usecase.Update(ctx, req.Id, mapAdminUpdateReq(req))
	if err != nil {
		return nil, err
	}
	events, _ := h.usecase.ListEvents(ctx, lead.ID)
	h.enrichSalesEventActors(ctx, events)
	return mapAdminLeadToResponse(lead, events), nil
}

func (h *AdminSalesHandler) DeleteOpportunity(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesDelete); err != nil {
		return nil, err
	}
	if err := h.usecase.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Message: "ok"}, nil
}

func (h *AdminSalesHandler) AssignOpportunity(ctx context.Context, req *crmpb.AssignAdminOpportunityRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesUpdate); err != nil {
		return nil, err
	}
	if err := h.usecase.Assign(ctx, req.Id, &dto.AdminOpportunityAssignDTO{
		ChargePersonID: req.ChargePersonId,
		AssignNote:     req.AssignNote,
	}); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "assigned"}, nil
}

func (h *AdminSalesHandler) SwitchStage(ctx context.Context, req *crmpb.SwitchAdminOpportunityStageRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesUpdate); err != nil {
		return nil, err
	}
	if err := h.usecase.SwitchStage(ctx, req.Id, req.StageId, req.StageNote); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "stage updated"}, nil
}

func (h *AdminSalesHandler) UpdateNote(ctx context.Context, req *crmpb.UpdateAdminOpportunityNoteRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesUpdate); err != nil {
		return nil, err
	}
	if err := h.usecase.UpdateNote(ctx, req.Id, req.Note); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "note updated"}, nil
}

func (h *AdminSalesHandler) UpdateStatus(ctx context.Context, req *crmpb.UpdateAdminOpportunityStatusRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesUpdate); err != nil {
		return nil, err
	}
	if err := h.usecase.UpdateStatus(ctx, req.Id, &dto.AdminOpportunityStatusDTO{
		OpportunityStatus: enums.EOpportunityStatus(req.OpportunityStatus),
		Note:              req.Note,
	}); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "status updated"}, nil
}

func (h *AdminSalesHandler) CreateFollowUp(ctx context.Context, req *crmpb.CreateAdminOpportunityFollowUpRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesUpdate); err != nil {
		return nil, err
	}
	t, err := parseTime(req.NextFollowUpAt)
	if err != nil || t == nil {
		return nil, _errors.ReturnError(400, "Thời điểm follow-up không hợp lệ")
	}
	if err := h.usecase.CreateFollowUp(ctx, req.Id, &dto.AdminOpportunityFollowUpDTO{
		NextFollowUpAt: *t,
		Content:        req.Content,
		NextAction:     req.NextAction,
	}); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "follow-up created"}, nil
}

func (h *AdminSalesHandler) CloseOpportunity(ctx context.Context, req *crmpb.CloseAdminOpportunityRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesUpdate); err != nil {
		return nil, err
	}
	if err := h.usecase.Close(ctx, req.Id, &dto.AdminOpportunityCloseDTO{
		Result:      enums.EOpportunityStatus(req.Result),
		CloseReason: req.CloseReason,
	}); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "closed"}, nil
}

func (h *AdminSalesHandler) ExportOpportunities(ctx context.Context, req *crmpb.ExportAdminOpportunitiesRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.RequiredPermissions(ctx, permSalesView); err != nil {
		return nil, err
	}
	if err := h.usecase.Export(ctx, req.Reason); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Message: "Đã ghi nhận yêu cầu xuất dữ liệu (file export sẽ bổ sung)."}, nil
}

func mapListSearch(req *crmpb.ListAdminOpportunitiesRequest) dto.LeadSearchDTO {
	var pipelineID *uint64
	if req.PipelineId != nil && req.GetPipelineId() > 0 {
		v := req.GetPipelineId()
		pipelineID = &v
	}
	unassigned := req.UnassignedOnly != nil && req.GetUnassignedOnly()
	var followFrom, followTo *time.Time
	if req.FollowUpFrom != nil && req.GetFollowUpFrom() != "" {
		if t, err := parseTime(req.GetFollowUpFrom()); err == nil {
			followFrom = t
		}
	}
	if req.FollowUpTo != nil && req.GetFollowUpTo() != "" {
		if t, err := parseTime(req.GetFollowUpTo()); err == nil {
			followTo = t
		}
	}
	var minEV, maxEV *float64
	if req.MinExpectedValue != nil {
		v := req.GetMinExpectedValue()
		minEV = &v
	}
	if req.MaxExpectedValue != nil {
		v := req.GetMaxExpectedValue()
		maxEV = &v
	}
	upgradeOnly := req.UpgradeOnly != nil && req.GetUpgradeOnly()
	renewalOnly := req.RenewalOnly != nil && req.GetRenewalOnly()
	return usecase.BuildAdminLeadSearch(
		req.Page, req.Size, req.GetQ(), req.StageIds, req.ChargePersonIds, pipelineID, req.Sources, unassigned,
		req.OpportunityStatuses, req.CustomerTypes, req.AdminSources, req.Probabilities, req.ChurnRisks,
		req.GetQueue(), followFrom, followTo, minEV, maxEV, upgradeOnly, renewalOnly, req.GetInterestedPlan(),
	)
}

func mapAdminCreateReq(req *crmpb.CreateAdminOpportunityRequest) *dto.AdminOpportunitySaveDTO {
	in := &dto.AdminOpportunitySaveDTO{
		FullName:          req.FullName,
		Phone:             req.Phone,
		Email:             req.Email,
		Company:           req.Company,
		Avatar:            req.Avatar,
		Address:           req.Address,
		Zalo:              req.Zalo,
		Note:              req.Note,
		Source:            enums.ESourceLead(req.Source),
		Priority:          enums.EPriority(req.Priority),
		AssignNote:        req.AssignNote,
		Title:             req.Title,
		CustomerType:      enums.ECustomerType(req.CustomerType),
		NeedSummary:       req.NeedSummary,
		InterestedPlan:    req.InterestedPlan,
		OpportunityStatus: enums.EOpportunityStatus(req.OpportunityStatus),
		AdminSource:       enums.EAdminSalesSource(req.AdminSource),
		Probability:       enums.EWinProbability(req.Probability),
		ChurnRisk:         enums.EChurnRisk(req.ChurnRisk),
		UpgradeSignal:     req.UpgradeSignal,
		RenewalSignal:     req.RenewalSignal,
		OwnerTeam:         req.OwnerTeam,
		Segment:           req.Segment,
		Region:            req.Region,
		Tags:              req.Tags,
	}
	if req.StageId != nil && req.GetStageId() > 0 {
		v := req.GetStageId()
		in.StageID = &v
	}
	if req.ChargePersonId != nil && req.GetChargePersonId() > 0 {
		v := req.GetChargePersonId()
		in.ChargePersonID = &v
	}
	if req.ExpectedValue != nil {
		v := req.GetExpectedValue()
		in.ExpectedValue = &v
	}
	if req.ExpectedCloseDate != nil && req.GetExpectedCloseDate() != "" {
		if t, err := parseDate(req.GetExpectedCloseDate()); err == nil {
			in.ExpectedCloseDate = t
		}
	}
	if req.NextFollowUpAt != nil && req.GetNextFollowUpAt() != "" {
		if t, err := parseTime(req.GetNextFollowUpAt()); err == nil {
			in.NextFollowUpAt = t
		}
	}
	if req.ContactId != nil && req.GetContactId() > 0 {
		v := req.GetContactId()
		in.ContactID = &v
	}
	return in
}

func mapAdminUpdateReq(req *crmpb.UpdateAdminOpportunityRequest) *dto.AdminOpportunitySaveDTO {
	in := &dto.AdminOpportunitySaveDTO{
		FullName:          req.FullName,
		Phone:             req.Phone,
		Email:             req.Email,
		Company:           req.Company,
		Avatar:            req.Avatar,
		Address:           req.Address,
		Zalo:              req.Zalo,
		Note:              req.Note,
		Source:            enums.ESourceLead(req.Source),
		Priority:          enums.EPriority(req.Priority),
		AssignNote:        req.AssignNote,
		Title:             req.Title,
		CustomerType:      enums.ECustomerType(req.CustomerType),
		NeedSummary:       req.NeedSummary,
		InterestedPlan:    req.InterestedPlan,
		OpportunityStatus: enums.EOpportunityStatus(req.OpportunityStatus),
		AdminSource:       enums.EAdminSalesSource(req.AdminSource),
		Probability:       enums.EWinProbability(req.Probability),
		ChurnRisk:         enums.EChurnRisk(req.ChurnRisk),
		UpgradeSignal:     req.UpgradeSignal,
		RenewalSignal:     req.RenewalSignal,
		OwnerTeam:         req.OwnerTeam,
		Segment:           req.Segment,
		Region:            req.Region,
		Tags:              req.Tags,
		ProposalRef:       req.ProposalRef,
		PaymentRequestRef: req.PaymentRequestRef,
		SubscriptionRef:   req.SubscriptionRef,
		TicketRef:         req.TicketRef,
	}
	if req.StageId != nil && req.GetStageId() > 0 {
		v := req.GetStageId()
		in.StageID = &v
	}
	if req.ChargePersonId != nil {
		v := req.GetChargePersonId()
		in.ChargePersonID = &v
	}
	if req.ExpectedValue != nil {
		v := req.GetExpectedValue()
		in.ExpectedValue = &v
	}
	if req.ExpectedCloseDate != nil && req.GetExpectedCloseDate() != "" {
		if t, err := parseDate(req.GetExpectedCloseDate()); err == nil {
			in.ExpectedCloseDate = t
		}
	}
	if req.NextFollowUpAt != nil && req.GetNextFollowUpAt() != "" {
		if t, err := parseTime(req.GetNextFollowUpAt()); err == nil {
			in.NextFollowUpAt = t
		}
	}
	return in
}

func mapAdminContactToItem(c *domain.ContactEntity) *crmpb.AdminContactItem {
	if c == nil {
		return nil
	}
	return &crmpb.AdminContactItem{
		Id:        c.ID,
		FullName:  c.FullName,
		Phone:     c.Phone,
		Email:     c.Email,
		Company:   c.Company,
		Avatar:    c.Avatar,
		Address:   c.Address,
		Zalo:      c.Zalo,
		Note:      c.Note,
		CreatedAt: toPBTimestamp(c.CreatedAt),
		UpdatedAt: toPBTimestamp(c.UpdatedAt),
	}
}

func mapAdminLeadToResponse(lead *domain.LeadEntity, events []dto.AdminOpportunityEventResponse) *crmpb.AdminOpportunityResponse {
	item := mapAdminLeadToItem(lead)
	resp := &crmpb.AdminOpportunityResponse{
		Data:       item,
		Events:     mapEvents(events),
		Commercial: buildCommercialProfile(lead),
	}
	if lead != nil && lead.Contact != nil {
		resp.Address = lead.Contact.Address
		resp.Zalo = lead.Contact.Zalo
	}
	return resp
}

func mapEvents(events []dto.AdminOpportunityEventResponse) []*crmpb.AdminOpportunityEventItem {
	out := make([]*crmpb.AdminOpportunityEventItem, 0, len(events))
	for _, e := range events {
		ev := &crmpb.AdminOpportunityEventItem{
			Id:        e.ID,
			ActorId:   e.ActorID,
			Action:    e.Action,
			Note:      e.Note,
			CreatedAt: toPBTimestamp(e.CreatedAt),
		}
		if e.ActorName != "" {
			name := e.ActorName
			ev.ActorName = &name
		}
		if e.BeforeJSON != "" {
			b := e.BeforeJSON
			ev.BeforeJson = &b
		}
		if e.AfterJSON != "" {
			a := e.AfterJSON
			ev.AfterJson = &a
		}
		out = append(out, ev)
	}
	return out
}

func mapAdminLeadToItem(lead *domain.LeadEntity) *crmpb.AdminOpportunityItem {
	if lead == nil {
		return nil
	}
	alerts := usecase.BuildOpportunityAlerts(lead)
	item := &crmpb.AdminOpportunityItem{
		Id:                lead.ID,
		ContactId:         lead.ContactID,
		Source:            int32(lead.Source),
		Priority:          int32(lead.Priority),
		Note:              lead.Note,
		AssignNote:        lead.AssignNote,
		CreatedAt:         toPBTimestamp(lead.CreatedAt),
		UpdatedAt:         toPBTimestamp(lead.UpdatedAt),
		Code:              lead.Code,
		Title:             lead.Title,
		CustomerType:      int32(lead.CustomerType),
		NeedSummary:       lead.NeedSummary,
		InterestedPlan:    lead.InterestedPlan,
		OpportunityStatus: int32(lead.OpportunityStatus),
		AdminSource:       int32(lead.AdminSource),
		Probability:       int32(lead.Probability),
		ChurnRisk:         int32(lead.ChurnRisk),
		UpgradeSignal:     lead.UpgradeSignal,
		RenewalSignal:     lead.RenewalSignal,
		OwnerTeam:         lead.OwnerTeam,
		ProposalRef:       lead.ProposalRef,
		PaymentRequestRef: lead.PaymentRequestRef,
		SubscriptionRef:   lead.SubscriptionRef,
		TicketRef:         lead.TicketRef,
		Segment:           lead.Segment,
		Region:            lead.Region,
		Tags:              lead.Tags,
		CloseReason:       lead.CloseReason,
		Alerts:            alerts,
		NextFollowUpAt:    toPBTimestamp(lead.NextFollowUpAt),
		ClosedAt:          toPBTimestamp(lead.ClosedAt),
	}
	if lead.ExpectedValue != nil {
		item.ExpectedValue = lead.ExpectedValue
	}
	if lead.ExpectedCloseDate != nil {
		s := lead.ExpectedCloseDate.Format("2006-01-02")
		item.ExpectedCloseDate = &s
	}
	if lead.NextFollowUpAt != nil && lead.NextFollowUpAt.Before(time.Now()) && !lead.OpportunityStatus.IsClosed() {
		item.FollowUpOverdue = true
	}
	if lead.RelatedUserID != nil {
		item.RelatedUserId = lead.RelatedUserID
	}
	if lead.RelatedBusinessID != nil {
		item.RelatedBusinessId = lead.RelatedBusinessID
	}
	if lead.StageID != nil {
		item.StageId = lead.StageID
	}
	if lead.PipelineID != nil {
		item.PipelineId = lead.PipelineID
	}
	if lead.ChargePersonID != nil {
		item.ChargePersonId = lead.ChargePersonID
	}
	if lead.Stage != nil {
		item.StageName = lead.Stage.StageName
		if lead.Stage.Pipeline != nil {
			item.PipelineName = lead.Stage.Pipeline.PipelineName
		}
	}
	if lead.Pipeline != nil && item.PipelineName == "" {
		item.PipelineName = lead.Pipeline.PipelineName
	}
	if lead.Contact != nil {
		item.FullName = lead.Contact.FullName
		item.Phone = lead.Contact.Phone
		item.Email = lead.Contact.Email
		item.Company = lead.Contact.Company
		item.Avatar = lead.Contact.Avatar
	}
	return item
}

func buildCommercialProfile(lead *domain.LeadEntity) *crmpb.AdminCommercialProfile {
	if lead == nil {
		return nil
	}
	usage := 55
	renewal := 50
	payment := 70
	support := 75
	churn := 20
	upgrade := 30
	activity := "active"
	plan := firstNonEmptyStr(lead.InterestedPlan, lead.SubscriptionRef, "—")
	paymentLabel := "Bình thường"
	if lead.OpportunityStatus == enums.OpportunityStatusWaitingPayment {
		payment = 35
		paymentLabel = "Chờ thanh toán"
	}
	if lead.ChurnRisk == enums.ChurnRiskHigh {
		churn = 80
		activity = "low"
	} else if lead.ChurnRisk == enums.ChurnRiskMedium {
		churn = 55
	}
	if lead.UpgradeSignal {
		upgrade = 75
	}
	if lead.RenewalSignal {
		renewal = 75
	}
	if lead.OpportunityStatus.IsClosed() {
		activity = "inactive"
	}
	overview := firstNonEmptyStr(lead.NeedSummary, lead.Title, "Hồ sơ thương mại 360°")
	p := &crmpb.AdminCommercialProfile{
		UsageHealth:        int32(usage),
		RenewalHealth:      int32(renewal),
		PaymentHealth:      int32(payment),
		SupportHealth:      int32(support),
		ChurnScore:         int32(churn),
		UpgradeScore:       int32(upgrade),
		ActivityStatus:     activity,
		CurrentPlanLabel:   plan,
		PaymentStatusLabel: paymentLabel,
		OverviewNote:       overview,
		OpenTicketCount:    0,
	}
	if lead.TicketRef != "" {
		p.OpenTicketCount = 1
	}
	p.UsageRows = []*crmpb.AdminContextRow{
		{Key: "plan", Label: "Gói quan tâm", Value: plan},
		{Key: "upgrade", Label: "Tín hiệu nâng cấp", Value: boolLabel(lead.UpgradeSignal)},
		{Key: "usage", Label: "Mức sử dụng", Value: "Chờ liên kết module Usage", Status: "pending"},
	}
	p.PlanRows = []*crmpb.AdminContextRow{
		{Key: "subscription", Label: "Subscription ref", Value: emptyDash(lead.SubscriptionRef)},
		{Key: "interested", Label: "Gói quan tâm", Value: emptyDash(lead.InterestedPlan)},
	}
	p.ProposalRows = []*crmpb.AdminContextRow{
		{Key: "proposal", Label: "Báo giá liên quan", Value: emptyDash(lead.ProposalRef), Status: "coming_soon", Href: ""},
	}
	p.PaymentRows = []*crmpb.AdminContextRow{
		{Key: "payment", Label: "Yêu cầu thanh toán", Value: emptyDash(lead.PaymentRequestRef)},
		{Key: "status", Label: "Trạng thái", Value: paymentLabel},
	}
	p.SupportRows = []*crmpb.AdminContextRow{
		{Key: "ticket", Label: "Ticket liên quan", Value: emptyDash(lead.TicketRef)},
		{Key: "churn", Label: "Rủi ro rời bỏ", Value: enums.ChurnRiskMap[lead.ChurnRisk]},
	}
	p.RenewalRows = []*crmpb.AdminContextRow{
		{Key: "renewal", Label: "Tín hiệu gia hạn", Value: boolLabel(lead.RenewalSignal)},
		{Key: "upgrade", Label: "Tín hiệu nâng cấp", Value: boolLabel(lead.UpgradeSignal)},
		{Key: "churn", Label: "Rủi ro rời bỏ", Value: enums.ChurnRiskMap[lead.ChurnRisk]},
	}
	return p
}

func emptyDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}

func boolLabel(v bool) string {
	if v {
		return "Có"
	}
	return "Không"
}

func firstNonEmptyStr(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func parseTime(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty")
	}
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("invalid time")
}

func parseDate(s string) (*time.Time, error) {
	return parseTime(s)
}

func toPBTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}

func (h *AdminSalesHandler) enrichSalesEventActors(ctx context.Context, events []dto.AdminOpportunityEventResponse) {
	if len(events) == 0 || h.authClient == nil || h.authClient.UserClient == nil {
		return
	}
	idSet := map[uint64]struct{}{}
	for _, e := range events {
		if e.ActorID > 0 {
			idSet[e.ActorID] = struct{}{}
		}
	}
	if len(idSet) == 0 {
		return
	}
	ids := make([]uint64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	resp, err := h.authClient.UserClient.GetAuthAdminByIds(ctx, ids)
	if err != nil || resp == nil {
		for i := range events {
			if events[i].ActorID > 0 && events[i].ActorName == "" {
				events[i].ActorName = fmt.Sprintf("User #%d", events[i].ActorID)
			}
		}
		return
	}
	names := map[uint64]string{}
	for _, a := range resp.GetData() {
		if a == nil {
			continue
		}
		label := strings.TrimSpace(a.GetFullName())
		if label == "" {
			label = strings.TrimSpace(a.GetUsername())
		}
		if label == "" {
			continue
		}
		if uid := a.GetUserId(); uid > 0 {
			names[uid] = label
		}
		if aid := a.GetId(); aid > 0 {
			names[aid] = label
		}
	}
	for i := range events {
		id := events[i].ActorID
		if id == 0 {
			continue
		}
		if n := names[id]; n != "" {
			events[i].ActorName = n
		} else {
			events[i].ActorName = fmt.Sprintf("User #%d", id)
		}
	}
}
