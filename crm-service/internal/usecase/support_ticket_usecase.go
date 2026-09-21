package usecase

import (
	"context"
	"crm/internal"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_errors "common/errors"
	_utils "common/utils"

	"gorm.io/datatypes"
)

type SupportTicketUsecase struct {
	repo repo.SupportTicketRepo
}

func NewSupportTicketUsecase(repo repo.SupportTicketRepo) *SupportTicketUsecase {
	return &SupportTicketUsecase{repo: repo}
}

func (u *SupportTicketUsecase) Create(ctx context.Context, req *dto.SupportTicketCreateRequest) (*dto.SupportTicketDetailResponse, error) {
	actorID := _utils.GetProfileIdWithContext(ctx)
	code, err := u.nextCode(ctx)
	if err != nil {
		return nil, err
	}
	imagesJSON := datatypes.JSON([]byte("[]"))
	if len(req.Images) > 0 {
		b, mErr := json.Marshal(req.Images)
		if mErr != nil {
			return nil, mErr
		}
		imagesJSON = datatypes.JSON(b)
	}
	ticket := &domain.SupportTicket{
		Code:            code,
		Title:           req.Title,
		Description:     req.Description,
		IssueType:       dto.DefaultIssueType(req.IssueType),
		Status:          enums.SupportTicketStatusNew,
		Priority:        dto.DefaultPriority(req.Priority),
		Product:         dto.DefaultProduct(req.Product),
		Source:          dto.DefaultSource(req.Source),
		RelatedUserID:   req.RelatedUserID,
		RelatedReportID: req.RelatedReportID,
		Images:          imagesJSON,
	}
	if actorID > 0 {
		ticket.CreatedBy = &actorID
	}
	ticket, err = u.repo.Create(ctx, ticket)
	if err != nil {
		return nil, err
	}
	_ = u.repo.AddEvent(ctx, &domain.SupportTicketEvent{
		TicketID: ticket.ID,
		ActorID:  actorID,
		Action:   "create",
		Note:     "Ticket created",
	})
	return u.GetByID(ctx, ticket.ID)
}

func (u *SupportTicketUsecase) nextCode(ctx context.Context) (string, error) {
	prefix := fmt.Sprintf("TK-%s-", time.Now().Format("20060102"))
	n, err := u.repo.CountByCodePrefix(ctx, prefix)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%04d", prefix, n+1), nil
}

func (u *SupportTicketUsecase) List(ctx context.Context, req *dto.SupportTicketListRequest) (*dto.SupportTicketListResponse, error) {
	tickets, total, err := u.repo.List(ctx, req)
	if err != nil {
		return nil, err
	}
	items := make([]dto.SupportTicketItemResponse, len(tickets))
	for i, t := range tickets {
		items[i] = toItem(t)
	}
	return &dto.SupportTicketListResponse{Data: items, Total: total}, nil
}

func (u *SupportTicketUsecase) Summary(ctx context.Context) (*dto.SupportTicketSummaryResponse, error) {
	return u.repo.Summary(ctx)
}

func (u *SupportTicketUsecase) GetByID(ctx context.Context, id uint64) (*dto.SupportTicketDetailResponse, error) {
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, _errors.ReturnError(service.SupportTicketNotFound)
	}
	notes, _ := u.repo.ListNotes(ctx, id)
	events, _ := u.repo.ListEvents(ctx, id, 50)
	return toDetail(ticket, notes, events), nil
}

func (u *SupportTicketUsecase) Update(ctx context.Context, id uint64, req *dto.SupportTicketUpdateRequest) error {
	actorID := _utils.GetProfileIdWithContext(ctx)
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return _errors.ReturnError(service.SupportTicketNotFound)
	}
	if req == nil || (req.Title == nil && req.Description == nil) {
		return _errors.ReturnError(service.SupportTicketUpdateContentRequired)
	}
	before, _ := json.Marshal(map[string]string{
		"title":       ticket.Title,
		"description": ticket.Description,
	})
	if req.Title != nil {
		title := *req.Title
		if title == "" {
			return _errors.ReturnError(service.SupportTicketTitleRequired)
		}
		ticket.Title = title
	}
	if req.Description != nil {
		ticket.Description = *req.Description
	}
	if err := u.repo.Update(ctx, ticket); err != nil {
		return err
	}
	after, _ := json.Marshal(map[string]string{
		"title":       ticket.Title,
		"description": ticket.Description,
	})
	return u.repo.AddEvent(ctx, &domain.SupportTicketEvent{
		TicketID:   id,
		ActorID:    actorID,
		Action:     "update",
		BeforeJSON: string(before),
		AfterJSON:  string(after),
		Note:       "Ticket content updated",
	})
}

func (u *SupportTicketUsecase) Assign(ctx context.Context, id uint64, req *dto.SupportTicketAssignRequest) error {
	actorID := _utils.GetProfileIdWithContext(ctx)
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return _errors.ReturnError(service.SupportTicketNotFound)
	}
	if ticket.Status == enums.SupportTicketStatusClosed {
		return _errors.ReturnError(service.SupportTicketClosedAssignmentDenied)
	}
	before, _ := json.Marshal(map[string]interface{}{"assigneeId": ticket.AssigneeID, "status": ticket.Status})
	ticket.AssigneeID = &req.AssigneeID
	if ticket.Status == enums.SupportTicketStatusNew {
		ticket.Status = enums.SupportTicketStatusInProgress
	}
	if err := u.repo.Update(ctx, ticket); err != nil {
		return err
	}
	after, _ := json.Marshal(map[string]interface{}{"assigneeId": ticket.AssigneeID, "status": ticket.Status})
	return u.repo.AddEvent(ctx, &domain.SupportTicketEvent{
		TicketID:   id,
		ActorID:    actorID,
		Action:     "assign",
		BeforeJSON: string(before),
		AfterJSON:  string(after),
		Note:       req.Note,
	})
}

func (u *SupportTicketUsecase) UpdatePriority(ctx context.Context, id uint64, req *dto.SupportTicketPriorityRequest) error {
	actorID := _utils.GetProfileIdWithContext(ctx)
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return _errors.ReturnError(service.SupportTicketNotFound)
	}
	before := fmt.Sprintf("%d", ticket.Priority)
	ticket.Priority = enums.SupportTicketPriority(req.Priority)
	if err := u.repo.Update(ctx, ticket); err != nil {
		return err
	}
	return u.repo.AddEvent(ctx, &domain.SupportTicketEvent{
		TicketID:   id,
		ActorID:    actorID,
		Action:     "priority",
		BeforeJSON: before,
		AfterJSON:  fmt.Sprintf("%d", ticket.Priority),
		Note:       req.Reason,
	})
}

func (u *SupportTicketUsecase) UpdateStatus(ctx context.Context, id uint64, req *dto.SupportTicketStatusRequest) error {
	actorID := _utils.GetProfileIdWithContext(ctx)
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return _errors.ReturnError(service.SupportTicketNotFound)
	}
	to := enums.SupportTicketStatus(req.Status)
	if !enums.CanTransitionSupportTicket(ticket.Status, to) {
		return _errors.ReturnError(service.SupportTicketStatusTransitionInvalid)
	}
	if strings.TrimSpace(req.Note) == "" {
		return _errors.ReturnError(service.SupportTicketStatusNoteRequired)
	}
	before := fmt.Sprintf("%d", ticket.Status)
	ticket.Status = to
	if to == enums.SupportTicketStatusClosed {
		now := time.Now()
		ticket.ClosedAt = &now
	}
	if err := u.repo.Update(ctx, ticket); err != nil {
		return err
	}
	return u.repo.AddEvent(ctx, &domain.SupportTicketEvent{
		TicketID:   id,
		ActorID:    actorID,
		Action:     "status",
		BeforeJSON: before,
		AfterJSON:  fmt.Sprintf("%d", to),
		Note:       req.Note,
	})
}

func (u *SupportTicketUsecase) Transfer(ctx context.Context, id uint64, req *dto.SupportTicketTransferRequest) error {
	actorID := _utils.GetProfileIdWithContext(ctx)
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return _errors.ReturnError(service.SupportTicketNotFound)
	}
	if ticket.Status == enums.SupportTicketStatusClosed {
		return _errors.ReturnError(service.SupportTicketClosedRoutingDenied)
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		return _errors.ReturnError(service.SupportTicketRoutingReasonRequired)
	}
	if !enums.IsValidSupportTicketHandlingTeam(req.HandlingTeam) {
		return _errors.ReturnError(service.SupportDepartmentInvalid)
	}
	toStatus := enums.SupportTicketStatusTransferred
	if !enums.CanTransitionSupportTicket(ticket.Status, toStatus) && ticket.Status != toStatus {
		return _errors.ReturnError(service.SupportTicketStatusTransitionInvalid)
	}
	before, _ := json.Marshal(map[string]interface{}{
		"status":       ticket.Status,
		"handlingTeam": ticket.HandlingTeam,
	})
	team := enums.SupportTicketHandlingTeam(req.HandlingTeam)
	ticket.HandlingTeam = &team
	ticket.Status = toStatus
	if err := u.repo.Update(ctx, ticket); err != nil {
		return err
	}
	after, _ := json.Marshal(map[string]interface{}{
		"status":       ticket.Status,
		"handlingTeam": ticket.HandlingTeam,
	})
	return u.repo.AddEvent(ctx, &domain.SupportTicketEvent{
		TicketID:   id,
		ActorID:    actorID,
		Action:     "transfer",
		BeforeJSON: string(before),
		AfterJSON:  string(after),
		Note:       reason,
	})
}

func (u *SupportTicketUsecase) AddNote(ctx context.Context, id uint64, req *dto.SupportTicketNoteRequest) error {
	actorID := _utils.GetProfileIdWithContext(ctx)
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return _errors.ReturnError(service.SupportTicketNotFound)
	}
	_, err = u.repo.AddNote(ctx, &domain.SupportTicketNote{
		TicketID: id,
		AuthorID: actorID,
		Body:     req.Body,
	})
	if err != nil {
		return err
	}
	return u.repo.AddEvent(ctx, &domain.SupportTicketEvent{
		TicketID: id,
		ActorID:  actorID,
		Action:   "note",
		Note:     "Internal note added",
	})
}

func (u *SupportTicketUsecase) Close(ctx context.Context, id uint64, req *dto.SupportTicketCloseRequest) error {
	note := strings.TrimSpace(req.Note)
	if note == "" {
		note = "Đóng ticket"
	}
	return u.UpdateStatus(ctx, id, &dto.SupportTicketStatusRequest{
		Status: uint32(enums.SupportTicketStatusClosed),
		Note:   note,
	})
}

func (u *SupportTicketUsecase) UpdateImages(ctx context.Context, id uint64, images []string) error {
	actorID := _utils.GetProfileIdWithContext(ctx)
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return _errors.ReturnError(service.SupportTicketNotFound)
	}
	existing := unmarshalImages(ticket.Images)
	before, _ := json.Marshal(existing)
	// Client gửi full array (thêm/xoá local rồi lưu) — thay thế, không merge
	next := sanitizeImageList(images)
	b, err := json.Marshal(next)
	if err != nil {
		return err
	}
	ticket.Images = datatypes.JSON(b)
	if err := u.repo.Update(ctx, ticket); err != nil {
		return err
	}
	after, _ := json.Marshal(next)
	return u.repo.AddEvent(ctx, &domain.SupportTicketEvent{
		TicketID:   id,
		ActorID:    actorID,
		Action:     "images",
		BeforeJSON: string(before),
		AfterJSON:  string(after),
		Note:       "Images updated",
	})
}

func unmarshalImages(raw datatypes.JSON) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return []string{}
	}
	return out
}

func sanitizeImageList(images []string) []string {
	seen := make(map[string]struct{}, len(images))
	out := make([]string, 0, len(images))
	for _, s := range images {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func mergeUniqueStrings(existing, incoming []string) []string {
	return sanitizeImageList(append(append([]string{}, existing...), incoming...))
}

func toItem(t domain.SupportTicket) dto.SupportTicketItemResponse {
	item := dto.SupportTicketItemResponse{
		ID:              t.ID,
		Code:            t.Code,
		Title:           t.Title,
		IssueType:       uint32(t.IssueType),
		Status:          uint32(t.Status),
		Priority:        uint32(t.Priority),
		AssigneeID:      t.AssigneeID,
		RelatedUserID:   t.RelatedUserID,
		RelatedReportID: t.RelatedReportID,
		Product:         uint32(t.Product),
		Source:          uint32(t.Source),
		Images:          unmarshalImages(t.Images),
		CreatedBy:       derefUint64(t.CreatedBy),
		UpdatedAt:       t.UpdatedAt,
		CreatedAt:       t.CreatedAt,
	}
	if t.HandlingTeam != nil {
		v := uint32(*t.HandlingTeam)
		item.HandlingTeam = &v
	}
	return item
}

func toDetail(t *domain.SupportTicket, notes []domain.SupportTicketNote, events []domain.SupportTicketEvent) *dto.SupportTicketDetailResponse {
	noteRes := make([]dto.SupportTicketNoteResponse, len(notes))
	for i, n := range notes {
		noteRes[i] = dto.SupportTicketNoteResponse{ID: n.ID, AuthorID: n.AuthorID, Body: n.Body, CreatedAt: n.CreatedAt}
	}
	eventRes := make([]dto.SupportTicketEventResponse, len(events))
	for i, e := range events {
		eventRes[i] = dto.SupportTicketEventResponse{
			ID: e.ID, ActorID: e.ActorID, Action: e.Action, Note: e.Note, CreatedAt: e.CreatedAt,
			BeforeJSON: e.BeforeJSON, AfterJSON: e.AfterJSON,
		}
	}
	return &dto.SupportTicketDetailResponse{
		ID:              t.ID,
		Code:            t.Code,
		Title:           t.Title,
		Description:     t.Description,
		IssueType:       uint32(t.IssueType),
		Status:          uint32(t.Status),
		Priority:        uint32(t.Priority),
		AssigneeID:      t.AssigneeID,
		RelatedUserID:   t.RelatedUserID,
		RelatedReportID: t.RelatedReportID,
		Product:         uint32(t.Product),
		Source:          uint32(t.Source),
		HandlingTeam:    handlingTeamPtr(t.HandlingTeam),
		Images:          unmarshalImages(t.Images),
		CreatedBy:       derefUint64(t.CreatedBy),
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
		ClosedAt:        t.ClosedAt,
		Notes:           noteRes,
		Events:          eventRes,
	}
}

func handlingTeamPtr(t *enums.SupportTicketHandlingTeam) *uint32 {
	if t == nil {
		return nil
	}
	v := uint32(*t)
	return &v
}

func derefUint64(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}
