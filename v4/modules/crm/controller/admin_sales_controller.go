package controller

import (
	"context"
	"net/http"
	"time"

	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/crm/dto"
	domain "mdn4/modules/crm/entity"
	"mdn4/modules/crm/enums"
	"mdn4/modules/crm/services"
	crmpb "pb/types/crm"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// AdminSalesController owns /v2/crm/admin/sales/* (16 routes).
type AdminSalesController struct {
	uc *services.AdminSalesUsecase
}

func NewAdminSalesController(uc *services.AdminSalesUsecase) *AdminSalesController {
	return &AdminSalesController{uc: uc}
}

func (c *AdminSalesController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/admin/sales/opportunities", c.listOpportunities)
	httpapi.GET(mux, "/v2/crm/admin/sales/summary", c.summary)
	httpapi.GET(mux, "/v2/crm/admin/sales/contacts", c.listContacts)
	httpapi.POST(mux, "/v2/crm/admin/sales/contacts", c.createContact)
	httpapi.GET(mux, "/v2/crm/admin/sales/opportunities/{id}/quick-view", c.quickView)
	httpapi.GET(mux, "/v2/crm/admin/sales/opportunities/{id}", c.getOpportunity)
	httpapi.POST(mux, "/v2/crm/admin/sales/opportunities", c.createOpportunity)
	httpapi.PATCH(mux, "/v2/crm/admin/sales/opportunities/{id}", c.updateOpportunity)
	httpapi.DELETE(mux, "/v2/crm/admin/sales/opportunities/{id}", c.deleteOpportunity)
	httpapi.POST(mux, "/v2/crm/admin/sales/opportunities/{id}/assign", c.assign)
	httpapi.POST(mux, "/v2/crm/admin/sales/opportunities/{id}/stage", c.switchStage)
	httpapi.POST(mux, "/v2/crm/admin/sales/opportunities/{id}/note", c.updateNote)
	httpapi.PATCH(mux, "/v2/crm/admin/sales/opportunities/{id}/status", c.updateStatus)
	httpapi.POST(mux, "/v2/crm/admin/sales/opportunities/{id}/follow-ups", c.followUp)
	httpapi.POST(mux, "/v2/crm/admin/sales/opportunities/{id}/close", c.close)
	httpapi.POST(mux, "/v2/crm/admin/sales/exports", c.export)
}

func (c *AdminSalesController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "admin sales service not configured")
	}
	return nil
}

func derefBool(v *bool) bool {
	return v != nil && *v
}

func parseOptionalTime(raw *string) *time.Time {
	if raw == nil || *raw == "" {
		return nil
	}
	layouts := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, *raw); err == nil {
			return &t
		}
	}
	return nil
}

func (c *AdminSalesController) listOpportunities(ctx context.Context, req dto.AdminSalesListOpportunitiesRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	search := services.BuildAdminLeadSearch(
		req.Page, req.Size, req.Q, req.StageIDs, req.ChargePersonIDs, req.PipelineID,
		req.Sources, derefBool(req.UnassignedOnly),
		req.OpportunityStatuses, req.CustomerTypes, req.AdminSources, req.Probabilities, req.ChurnRisks,
		req.Queue, parseOptionalTime(req.FollowUpFrom), parseOptionalTime(req.FollowUpTo),
		req.MinExpectedValue, req.MaxExpectedValue,
		derefBool(req.UpgradeOnly), derefBool(req.RenewalOnly), req.InterestedPlan,
	)
	items, total, err := c.uc.List(ctx, search)
	if err != nil {
		return nil, err
	}
	out := make([]*crmpb.AdminOpportunityItem, 0, len(items))
	for i := range items {
		out = append(out, mapAdminLeadToItem(&items[i]))
	}
	resp := &crmpb.ListAdminOpportunitiesResponse{Data: out, Total: total}
	if req.ViewMode == "funnel" {
		funnel, err := c.uc.Funnel(ctx, search)
		if err != nil {
			return nil, err
		}
		resp.Funnel = make([]*crmpb.AdminFunnelBucket, 0, len(funnel))
		for _, b := range funnel {
			resp.Funnel = append(resp.Funnel, &crmpb.AdminFunnelBucket{
				Status: int32(b.Status), Label: b.Label, Count: b.Count, ExpectedValueSum: b.ExpectedValueSum,
			})
		}
	}
	return protoDataResult(resp)
}

func (c *AdminSalesController) summary(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	s, err := c.uc.Summary(ctx)
	if err != nil {
		return nil, err
	}
	return protoDataResult(&crmpb.AdminOpportunitiesSummaryResponse{
		Total: s.Total, OverdueFollowUp: s.OverdueFollowUp, Unassigned: s.Unassigned,
		WaitingQuote: s.WaitingQuote, WaitingApproval: s.WaitingApproval, WaitingPayment: s.WaitingPayment,
		Renewal: s.Renewal, Upgrade: s.Upgrade, ChurnRisk: s.ChurnRisk,
		StatusNew: s.StatusNew, Consulting: s.Consulting, Won: s.Won, Lost: s.Lost,
		ExpectedRevenue: s.ExpectedRevenue, WaitingPaymentValue: s.WaitingPaymentValue, WonValue: s.WonValue,
	})
}

func (c *AdminSalesController) listContacts(ctx context.Context, req dto.AdminSalesListContactsRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	items, total, err := c.uc.ListContacts(ctx, req.Page, req.Size, req.Q)
	if err != nil {
		return nil, err
	}
	out := make([]*crmpb.AdminContactItem, 0, len(items))
	for i := range items {
		out = append(out, mapAdminContactToItem(&items[i]))
	}
	return protoDataResult(&crmpb.ListAdminContactsResponse{Data: out, Total: total})
}

func (c *AdminSalesController) createContact(ctx context.Context, req dto.AdminSalesCreateContactRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	contact, err := c.uc.CreateContact(ctx, &dto.AdminContactSaveDTO{
		FullName: req.FullName, Phone: req.Phone, Email: req.Email, Company: req.Company,
		Avatar: req.Avatar, Address: req.Address, Zalo: req.Zalo, Note: req.Note,
	})
	if err != nil {
		return nil, err
	}
	return protoDataResult(&crmpb.AdminContactResponse{Data: mapAdminContactToItem(contact)})
}

func (c *AdminSalesController) getOpportunity(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	lead, err := c.uc.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	events, _ := c.uc.ListEvents(ctx, req.ID)
	return protoDataResult(mapAdminLeadToResponse(lead, events))
}

func (c *AdminSalesController) quickView(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	lead, err := c.uc.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	events, _ := c.uc.ListEvents(ctx, req.ID)
	if len(events) > 5 {
		events = events[:5]
	}
	item := mapAdminLeadToItem(lead)
	resp := &crmpb.AdminOpportunityQuickViewResponse{
		Data: item, Alerts: services.BuildOpportunityAlerts(lead), RecentEvents: mapAdminEvents(events),
	}
	if lead.Contact != nil {
		resp.Address = lead.Contact.Address
		resp.Zalo = lead.Contact.Zalo
	}
	return protoDataResult(resp)
}

func (c *AdminSalesController) createOpportunity(ctx context.Context, req dto.AdminSalesCreateOpportunityRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	lead, err := c.uc.Create(ctx, mapAdminSaveFromCreate(req))
	if err != nil {
		return nil, err
	}
	events, _ := c.uc.ListEvents(ctx, lead.ID)
	return protoDataResult(mapAdminLeadToResponse(lead, events))
}

func (c *AdminSalesController) updateOpportunity(ctx context.Context, req dto.AdminSalesUpdateOpportunityRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "ID không hợp lệ")
	}
	lead, err := c.uc.Update(ctx, req.ID, mapAdminSaveFromUpdate(req))
	if err != nil {
		return nil, err
	}
	events, _ := c.uc.ListEvents(ctx, lead.ID)
	return protoDataResult(mapAdminLeadToResponse(lead, events))
}

func (c *AdminSalesController) deleteOpportunity(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "ok")
}

func (c *AdminSalesController) assign(ctx context.Context, req dto.AdminSalesAssignRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Assign(ctx, req.ID, &dto.AdminOpportunityAssignDTO{
		ChargePersonID: req.ChargePersonID, AssignNote: req.AssignNote,
	}); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "assigned")
}

func (c *AdminSalesController) switchStage(ctx context.Context, req dto.AdminSalesSwitchStageRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.SwitchStage(ctx, req.ID, req.StageID, req.StageNote); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "stage updated")
}

func (c *AdminSalesController) updateNote(ctx context.Context, req dto.AdminSalesUpdateNoteRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.UpdateNote(ctx, req.ID, req.Note); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "note updated")
}

func (c *AdminSalesController) updateStatus(ctx context.Context, req dto.AdminSalesUpdateStatusRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.UpdateStatus(ctx, req.ID, &dto.AdminOpportunityStatusDTO{
		OpportunityStatus: enums.EOpportunityStatus(req.OpportunityStatus), Note: req.Note,
	}); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "status updated")
}

func (c *AdminSalesController) followUp(ctx context.Context, req dto.AdminSalesFollowUpRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	t := parseOptionalTime(&req.NextFollowUpAt)
	if t == nil {
		return nil, _errors.ReturnError(400, "Thời điểm follow-up không hợp lệ")
	}
	if err := c.uc.CreateFollowUp(ctx, req.ID, &dto.AdminOpportunityFollowUpDTO{
		NextFollowUpAt: *t, Content: req.Content, NextAction: req.NextAction,
	}); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "follow-up created")
}

func (c *AdminSalesController) close(ctx context.Context, req dto.AdminSalesCloseRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Close(ctx, req.ID, &dto.AdminOpportunityCloseDTO{
		Result: enums.EOpportunityStatus(req.Result), CloseReason: req.CloseReason,
	}); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "closed")
}

func (c *AdminSalesController) export(ctx context.Context, req dto.AdminSalesExportRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Export(ctx, req.Reason); err != nil {
		return nil, err
	}
	return submitResult(0, "export queued")
}

func toPBTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func mapAdminContactToItem(c *domain.Contact) *crmpb.AdminContactItem {
	if c == nil {
		return nil
	}
	return &crmpb.AdminContactItem{
		Id: c.ID, FullName: c.FullName, Phone: c.Phone, Email: c.Email, Company: c.Company,
		Avatar: c.Avatar, Address: c.Address, Zalo: c.Zalo, Note: c.Note,
		CreatedAt: toPBTimestamp(c.CreatedAt), UpdatedAt: toPBTimestamp(c.UpdatedAt),
	}
}

func mapAdminEvents(events []dto.AdminOpportunityEventResponse) []*crmpb.AdminOpportunityEventItem {
	out := make([]*crmpb.AdminOpportunityEventItem, 0, len(events))
	for _, e := range events {
		ev := &crmpb.AdminOpportunityEventItem{
			Id: e.ID, ActorId: e.ActorID, Action: e.Action, Note: e.Note, CreatedAt: toPBTimestamp(e.CreatedAt),
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
	item := &crmpb.AdminOpportunityItem{
		Id: lead.ID, ContactId: lead.ContactID, Source: int32(lead.Source), Priority: int32(lead.Priority),
		Note: lead.Note, AssignNote: lead.AssignNote, CreatedAt: toPBTimestamp(lead.CreatedAt), UpdatedAt: toPBTimestamp(lead.UpdatedAt),
		Code: lead.Code, Title: lead.Title, CustomerType: int32(lead.CustomerType), NeedSummary: lead.NeedSummary,
		InterestedPlan: lead.InterestedPlan, OpportunityStatus: int32(lead.OpportunityStatus), AdminSource: int32(lead.AdminSource),
		Probability: int32(lead.Probability), ChurnRisk: int32(lead.ChurnRisk), UpgradeSignal: lead.UpgradeSignal,
		RenewalSignal: lead.RenewalSignal, OwnerTeam: lead.OwnerTeam, ProposalRef: lead.ProposalRef,
		PaymentRequestRef: lead.PaymentRequestRef, SubscriptionRef: lead.SubscriptionRef, TicketRef: lead.TicketRef,
		Segment: lead.Segment, Region: lead.Region, Tags: lead.Tags, CloseReason: lead.CloseReason,
		Alerts: services.BuildOpportunityAlerts(lead), NextFollowUpAt: toPBTimestamp(lead.NextFollowUpAt), ClosedAt: toPBTimestamp(lead.ClosedAt),
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
	if lead.Contact != nil {
		item.FullName = lead.Contact.FullName
		item.Phone = lead.Contact.Phone
		item.Email = lead.Contact.Email
		item.Company = lead.Contact.Company
		item.Avatar = lead.Contact.Avatar
	}
	return item
}

func mapAdminLeadToResponse(lead *domain.LeadEntity, events []dto.AdminOpportunityEventResponse) *crmpb.AdminOpportunityResponse {
	resp := &crmpb.AdminOpportunityResponse{Data: mapAdminLeadToItem(lead), Events: mapAdminEvents(events)}
	if lead != nil && lead.Contact != nil {
		resp.Address = lead.Contact.Address
		resp.Zalo = lead.Contact.Zalo
	}
	return resp
}

func mapAdminSaveFromCreate(req dto.AdminSalesCreateOpportunityRequest) *dto.AdminOpportunitySaveDTO {
	in := &dto.AdminOpportunitySaveDTO{
		FullName: req.FullName, Phone: req.Phone, Email: req.Email, Company: req.Company, Avatar: req.Avatar,
		Address: req.Address, Zalo: req.Zalo, Note: req.Note, Source: enums.ESourceLead(req.Source),
		StageID: req.StageID, ChargePersonID: req.ChargePersonID, Priority: enums.EPriority(req.Priority),
		AssignNote: req.AssignNote, Title: req.Title, CustomerType: enums.ECustomerType(req.CustomerType),
		NeedSummary: req.NeedSummary, InterestedPlan: req.InterestedPlan,
		OpportunityStatus: enums.EOpportunityStatus(req.OpportunityStatus), AdminSource: enums.EAdminSalesSource(req.AdminSource),
		ExpectedValue: req.ExpectedValue, Probability: enums.EWinProbability(req.Probability),
		ChurnRisk: enums.EChurnRisk(req.ChurnRisk), UpgradeSignal: req.UpgradeSignal, RenewalSignal: req.RenewalSignal,
		OwnerTeam: req.OwnerTeam, Segment: req.Segment, Region: req.Region, Tags: req.Tags, ContactID: req.ContactID,
	}
	in.ExpectedCloseDate = parseOptionalTime(req.ExpectedCloseDate)
	in.NextFollowUpAt = parseOptionalTime(req.NextFollowUpAt)
	return in
}

func mapAdminSaveFromUpdate(req dto.AdminSalesUpdateOpportunityRequest) *dto.AdminOpportunitySaveDTO {
	in := &dto.AdminOpportunitySaveDTO{
		FullName: req.FullName, Phone: req.Phone, Email: req.Email, Company: req.Company, Avatar: req.Avatar,
		Address: req.Address, Zalo: req.Zalo, Note: req.Note, Source: enums.ESourceLead(req.Source),
		StageID: req.StageID, ChargePersonID: req.ChargePersonID, Priority: enums.EPriority(req.Priority),
		AssignNote: req.AssignNote, Title: req.Title, CustomerType: enums.ECustomerType(req.CustomerType),
		NeedSummary: req.NeedSummary, InterestedPlan: req.InterestedPlan,
		OpportunityStatus: enums.EOpportunityStatus(req.OpportunityStatus), AdminSource: enums.EAdminSalesSource(req.AdminSource),
		ExpectedValue: req.ExpectedValue, Probability: enums.EWinProbability(req.Probability),
		ChurnRisk: enums.EChurnRisk(req.ChurnRisk), UpgradeSignal: req.UpgradeSignal, RenewalSignal: req.RenewalSignal,
		OwnerTeam: req.OwnerTeam, Segment: req.Segment, Region: req.Region, Tags: req.Tags,
		ProposalRef: req.ProposalRef, PaymentRequestRef: req.PaymentRequestRef, SubscriptionRef: req.SubscriptionRef, TicketRef: req.TicketRef,
	}
	in.ExpectedCloseDate = parseOptionalTime(req.ExpectedCloseDate)
	in.NextFollowUpAt = parseOptionalTime(req.NextFollowUpAt)
	return in
}
