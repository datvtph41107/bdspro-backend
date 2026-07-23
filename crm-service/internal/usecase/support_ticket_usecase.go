package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
	"encoding/json"
	"fmt"
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
		return nil, _errors.ReturnError(404, "ticket not found")
	}
	notes, _ := u.repo.ListNotes(ctx, id)
	events, _ := u.repo.ListEvents(ctx, id, 50)
	return toDetail(ticket, notes, events), nil
}

func (u *SupportTicketUsecase) Assign(ctx context.Context, id uint64, req *dto.SupportTicketAssignRequest) error {
	actorID := _utils.GetProfileIdWithContext(ctx)
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return _errors.ReturnError(404, "ticket not found")
	}
	if ticket.Status == enums.SupportTicketStatusClosed {
		return _errors.ReturnError(400, "cannot assign closed ticket")
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
		return _errors.ReturnError(404, "ticket not found")
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
		return _errors.ReturnError(404, "ticket not found")
	}
	to := enums.SupportTicketStatus(req.Status)
	if !enums.CanTransitionSupportTicket(ticket.Status, to) {
		return _errors.ReturnError(400, "STATUS_TRANSITION_INVALID")
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

func (u *SupportTicketUsecase) AddNote(ctx context.Context, id uint64, req *dto.SupportTicketNoteRequest) error {
	actorID := _utils.GetProfileIdWithContext(ctx)
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return _errors.ReturnError(404, "ticket not found")
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
	return u.UpdateStatus(ctx, id, &dto.SupportTicketStatusRequest{
		Status: uint32(enums.SupportTicketStatusClosed),
		Note:   req.Note,
	})
}

func (u *SupportTicketUsecase) UpdateImages(ctx context.Context, id uint64, images []string) error {
	actorID := _utils.GetProfileIdWithContext(ctx)
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return _errors.ReturnError(404, "ticket not found")
	}
	existing := unmarshalImages(ticket.Images)
	before, _ := json.Marshal(existing)
	merged := mergeUniqueStrings(existing, images)
	b, err := json.Marshal(merged)
	if err != nil {
		return err
	}
	ticket.Images = datatypes.JSON(b)
	if err := u.repo.Update(ctx, ticket); err != nil {
		return err
	}
	after, _ := json.Marshal(merged)
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

func mergeUniqueStrings(existing, incoming []string) []string {
	seen := make(map[string]struct{}, len(existing)+len(incoming))
	out := make([]string, 0, len(existing)+len(incoming))
	for _, s := range existing {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	for _, s := range incoming {
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

func toItem(t domain.SupportTicket) dto.SupportTicketItemResponse {
	return dto.SupportTicketItemResponse{
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
		Images:          unmarshalImages(t.Images),
		UpdatedAt:       t.UpdatedAt,
		CreatedAt:       t.CreatedAt,
	}
}

func toDetail(t *domain.SupportTicket, notes []domain.SupportTicketNote, events []domain.SupportTicketEvent) *dto.SupportTicketDetailResponse {
	noteRes := make([]dto.SupportTicketNoteResponse, len(notes))
	for i, n := range notes {
		noteRes[i] = dto.SupportTicketNoteResponse{ID: n.ID, AuthorID: n.AuthorID, Body: n.Body, CreatedAt: n.CreatedAt}
	}
	eventRes := make([]dto.SupportTicketEventResponse, len(events))
	for i, e := range events {
		eventRes[i] = dto.SupportTicketEventResponse{ID: e.ID, ActorID: e.ActorID, Action: e.Action, Note: e.Note, CreatedAt: e.CreatedAt}
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
		Images:          unmarshalImages(t.Images),
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
		ClosedAt:        t.ClosedAt,
		Notes:           noteRes,
		Events:          eventRes,
	}
}
