package handler

import (
	"context"
	"crm/infra/client"
	"crm/internal/dto"
	"crm/internal/usecase"
	"fmt"
	"strings"

	_dto "common/domain/dto"
	_utils "common/utils"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"google.golang.org/protobuf/types/known/emptypb"
)

// Quyền theo hành động. Key ADMIN_HT_SUA/ADMIN_HT_XOA là key cũ trước khi tách
// Xử lý / Đóng, vẫn nhận cho tới khi mọi môi trường chạy script rename permission.
var (
	permTicketView = []string{
		"ADMIN_HT_XEM",
		"ADMIN_HT_XU_LY", "ADMIN_HT_SUA",
		"ADMIN_HT_DONG", "ADMIN_HT_XOA",
	}
	permTicketCreate  = []string{"ADMIN_HT_TAO"}
	permTicketProcess = []string{"ADMIN_HT_XU_LY", "ADMIN_HT_SUA"}
	permTicketClose   = []string{"ADMIN_HT_DONG", "ADMIN_HT_XOA"}
)

type SupportTicketHandler struct {
	crmpb.UnimplementedSupportTicketServiceServer
	usecase    *usecase.SupportTicketUsecase
	authClient *client.AuthClient
}

// @bind: crm/infra/handler.SupportTicketHandler
func NewSupportTicketHandler(usecase *usecase.SupportTicketUsecase, authClient *client.AuthClient) *SupportTicketHandler {
	return &SupportTicketHandler{usecase: usecase, authClient: authClient}
}

func (h *SupportTicketHandler) ListSupportTickets(ctx context.Context, req *crmpb.ListSupportTicketsRequest) (*crmpb.ListSupportTicketsResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketView); err != nil {
		return nil, err
	}
	listReq := &dto.SupportTicketListRequest{
		Pagable: _dto.Pagable{Page: req.Page, Size: req.Size},
		Q:       req.GetQ(),
	}
	if req.Status != nil {
		v := req.GetStatus()
		listReq.Status = &v
	}
	if req.Priority != nil {
		v := req.GetPriority()
		listReq.Priority = &v
	}
	if req.IssueType != nil {
		v := req.GetIssueType()
		listReq.IssueType = &v
	}
	if req.AssigneeId != nil {
		v := req.GetAssigneeId()
		listReq.AssigneeID = &v
	}
	if req.Product != nil {
		v := req.GetProduct()
		listReq.Product = &v
	}
	if req.Source != nil {
		v := req.GetSource()
		listReq.Source = &v
	}
	if req.HandlingTeam != nil {
		v := req.GetHandlingTeam()
		listReq.HandlingTeam = &v
	}
	if req.UnassignedOnly != nil && req.GetUnassignedOnly() {
		listReq.UnassignedOnly = true
	}
	res, err := h.usecase.List(ctx, listReq)
	if err != nil {
		return nil, err
	}
	data := make([]*crmpb.SupportTicketItem, len(res.Data))
	for i, item := range res.Data {
		data[i] = toProtoItem(item)
	}
	return &crmpb.ListSupportTicketsResponse{Data: data, Total: res.Total}, nil
}

func (h *SupportTicketHandler) GetSupportTicketSummary(ctx context.Context, _ *emptypb.Empty) (*crmpb.SupportTicketSummaryResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketView); err != nil {
		return nil, err
	}
	res, err := h.usecase.Summary(ctx)
	if err != nil {
		return nil, err
	}
	return &crmpb.SupportTicketSummaryResponse{
		Total:             res.Total,
		StatusNew:         res.StatusNew,
		StatusInProgress:  res.StatusInProgress,
		StatusWaitingUser: res.StatusWaitingUser,
		StatusResolved:    res.StatusResolved,
		StatusClosed:      res.StatusClosed,
		StatusTransferred: res.StatusTransferred,
		Unassigned:        res.Unassigned,
		PriorityCritical:  res.PriorityCritical,
	}, nil
}

func (h *SupportTicketHandler) GetSupportTicket(ctx context.Context, req *sharepb.IdRequest) (*crmpb.GetSupportTicketResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketView); err != nil {
		return nil, err
	}
	res, err := h.usecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	h.enrichTicketEventActors(ctx, res)
	return &crmpb.GetSupportTicketResponse{Ticket: toProtoDetail(res)}, nil
}

func (h *SupportTicketHandler) CreateSupportTicket(ctx context.Context, req *crmpb.CreateSupportTicketRequest) (*crmpb.GetSupportTicketResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketCreate); err != nil {
		return nil, err
	}
	createReq := &dto.SupportTicketCreateRequest{
		Title:       req.Title,
		Description: req.Description,
		IssueType:   req.IssueType,
		Priority:    req.Priority,
		Product:     req.Product,
		Source:      req.Source,
		Images:      req.Images,
	}
	if req.RelatedUserId != nil {
		v := req.GetRelatedUserId()
		createReq.RelatedUserID = &v
	}
	if req.RelatedReportId != nil {
		v := req.GetRelatedReportId()
		createReq.RelatedReportID = &v
	}
	res, err := h.usecase.Create(ctx, createReq)
	if err != nil {
		return nil, err
	}
	return &crmpb.GetSupportTicketResponse{Ticket: toProtoDetail(res)}, nil
}

func (h *SupportTicketHandler) UpdateSupportTicket(ctx context.Context, req *crmpb.UpdateSupportTicketRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketProcess); err != nil {
		return nil, err
	}
	updateReq := &dto.SupportTicketUpdateRequest{}
	if req.Title != nil {
		v := req.GetTitle()
		updateReq.Title = &v
	}
	if req.Description != nil {
		v := req.GetDescription()
		updateReq.Description = &v
	}
	if err := h.usecase.Update(ctx, req.Id, updateReq); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "updated"}, nil
}

func (h *SupportTicketHandler) AssignSupportTicket(ctx context.Context, req *crmpb.AssignSupportTicketRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketProcess); err != nil {
		return nil, err
	}
	err := h.usecase.Assign(ctx, req.Id, &dto.SupportTicketAssignRequest{
		AssigneeID: req.AssigneeId,
		Note:       req.GetNote(),
	})
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "assigned"}, nil
}

func (h *SupportTicketHandler) UpdateSupportTicketPriority(ctx context.Context, req *crmpb.UpdateSupportTicketPriorityRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketProcess); err != nil {
		return nil, err
	}
	err := h.usecase.UpdatePriority(ctx, req.Id, &dto.SupportTicketPriorityRequest{
		Priority: req.Priority,
		Reason:   req.Reason,
	})
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "priority updated"}, nil
}

func (h *SupportTicketHandler) UpdateSupportTicketStatus(ctx context.Context, req *crmpb.UpdateSupportTicketStatusRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketProcess); err != nil {
		return nil, err
	}
	err := h.usecase.UpdateStatus(ctx, req.Id, &dto.SupportTicketStatusRequest{
		Status: req.Status,
		Note:   req.Note,
	})
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "status updated"}, nil
}

func (h *SupportTicketHandler) AddSupportTicketNote(ctx context.Context, req *crmpb.AddSupportTicketNoteRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketProcess); err != nil {
		return nil, err
	}
	err := h.usecase.AddNote(ctx, req.Id, &dto.SupportTicketNoteRequest{Body: req.Body})
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "note added"}, nil
}

func (h *SupportTicketHandler) CloseSupportTicket(ctx context.Context, req *crmpb.CloseSupportTicketRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketClose); err != nil {
		return nil, err
	}
	err := h.usecase.Close(ctx, req.Id, &dto.SupportTicketCloseRequest{Note: req.Note})
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "closed"}, nil
}

func (h *SupportTicketHandler) UpdateSupportTicketImages(ctx context.Context, req *crmpb.UpdateSupportTicketImagesRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketProcess); err != nil {
		return nil, err
	}
	err := h.usecase.UpdateImages(ctx, req.Id, req.Images)
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "images updated"}, nil
}

func (h *SupportTicketHandler) TransferSupportTicket(ctx context.Context, req *crmpb.TransferSupportTicketRequest) (*sharepb.SubmitResponse, error) {
	if err := h.authClient.HasPermissions(ctx, permTicketProcess); err != nil {
		return nil, err
	}
	err := h.usecase.Transfer(ctx, req.Id, &dto.SupportTicketTransferRequest{
		HandlingTeam: req.HandlingTeam,
		Reason:       req.Reason,
	})
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "transferred"}, nil
}

func toProtoItem(item dto.SupportTicketItemResponse) *crmpb.SupportTicketItem {
	out := &crmpb.SupportTicketItem{
		Id:        item.ID,
		Code:      item.Code,
		Title:     item.Title,
		IssueType: item.IssueType,
		Status:    item.Status,
		Priority:  item.Priority,
		Product:   item.Product,
		Source:    item.Source,
		Images:    item.Images,
		CreatedBy: item.CreatedBy,
		UpdatedAt: _utils.FormatTimeToString(item.UpdatedAt),
		CreatedAt: _utils.FormatTimeToString(item.CreatedAt),
	}
	if item.AssigneeID != nil {
		out.AssigneeId = item.AssigneeID
	}
	if item.RelatedUserID != nil {
		out.RelatedUserId = item.RelatedUserID
	}
	if item.RelatedReportID != nil {
		out.RelatedReportId = item.RelatedReportID
	}
	if item.HandlingTeam != nil {
		out.HandlingTeam = item.HandlingTeam
	}
	return out
}

func toProtoDetail(d *dto.SupportTicketDetailResponse) *crmpb.SupportTicketDetail {
	if d == nil {
		return nil
	}
	notes := make([]*crmpb.SupportTicketNoteItem, len(d.Notes))
	for i, n := range d.Notes {
		notes[i] = &crmpb.SupportTicketNoteItem{
			Id: n.ID, AuthorId: n.AuthorID, Body: n.Body, CreatedAt: _utils.FormatTimeToString(n.CreatedAt),
		}
	}
	events := make([]*crmpb.SupportTicketEventItem, len(d.Events))
	for i, e := range d.Events {
		ev := &crmpb.SupportTicketEventItem{
			Id: e.ID, ActorId: e.ActorID, Action: e.Action, Note: e.Note, CreatedAt: _utils.FormatTimeToString(e.CreatedAt),
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
		events[i] = ev
	}
	out := &crmpb.SupportTicketDetail{
		Id:          d.ID,
		Code:        d.Code,
		Title:       d.Title,
		Description: d.Description,
		IssueType:   d.IssueType,
		Status:      d.Status,
		Priority:    d.Priority,
		Product:     d.Product,
		Source:      d.Source,
		Images:      d.Images,
		CreatedBy:   d.CreatedBy,
		CreatedAt:   _utils.FormatTimeToString(d.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(d.UpdatedAt),
		Notes:       notes,
		Events:      events,
	}
	if d.AssigneeID != nil {
		out.AssigneeId = d.AssigneeID
	}
	if d.RelatedUserID != nil {
		out.RelatedUserId = d.RelatedUserID
	}
	if d.RelatedReportID != nil {
		out.RelatedReportId = d.RelatedReportID
	}
	if d.HandlingTeam != nil {
		out.HandlingTeam = d.HandlingTeam
	}
	if d.ClosedAt != nil {
		s := _utils.FormatTimeToString(d.ClosedAt)
		out.ClosedAt = &s
	}
	return out
}

func (h *SupportTicketHandler) enrichTicketEventActors(ctx context.Context, d *dto.SupportTicketDetailResponse) {
	if d == nil || len(d.Events) == 0 || h.authClient == nil || h.authClient.UserClient == nil {
		return
	}
	idSet := map[uint64]struct{}{}
	for _, e := range d.Events {
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
		for i := range d.Events {
			if d.Events[i].ActorID > 0 && d.Events[i].ActorName == "" {
				d.Events[i].ActorName = fmt.Sprintf("User #%d", d.Events[i].ActorID)
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
	for i := range d.Events {
		id := d.Events[i].ActorID
		if id == 0 {
			continue
		}
		if n := names[id]; n != "" {
			d.Events[i].ActorName = n
		} else {
			d.Events[i].ActorName = fmt.Sprintf("User #%d", id)
		}
	}
}
