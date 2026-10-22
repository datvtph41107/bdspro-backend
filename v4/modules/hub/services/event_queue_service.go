package services

import (
	"context"
	"strconv"
	"time"

	_errors "common/errors"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/entity"
	"mdn4/modules/hub/repository"
)

// EventQueueService owns event queue CRUD.
type EventQueueService interface {
	List(ctx context.Context, text string, page, size int) ([]*dto.EventQueueItem, int64, error)
	GetByID(ctx context.Context, id uint64) (*dto.EventQueueItem, error)
	Create(ctx context.Context, req dto.EventQueueSaveRequest) (*dto.EventQueueItem, error)
	Update(ctx context.Context, id uint64, req dto.EventQueueSaveRequest) (*dto.EventQueueItem, error)
	Delete(ctx context.Context, id uint64) error
	UpdateStatus(ctx context.Context, req dto.EventQueueStatusRequest) (*dto.EventQueueItem, error)
}

// @bind: mdn4/modules/hub/services.EventQueueService
type LocalEventQueueService struct {
	repo repository.EventQueueRepository
}

func NewLocalEventQueueService(repo repository.EventQueueRepository) *LocalEventQueueService {
	return &LocalEventQueueService{repo: repo}
}

func (s *LocalEventQueueService) List(ctx context.Context, text string, page, size int) ([]*dto.EventQueueItem, int64, error) {
	rows, total, err := s.repo.Search(ctx, text, page, size)
	if err != nil {
		return nil, 0, err
	}
	return mapEventQueueItems(rows), total, nil
}

func (s *LocalEventQueueService) GetByID(ctx context.Context, id uint64) (*dto.EventQueueItem, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, _errors.NotFoundException("event queue not found")
	}
	return mapEventQueueItem(row), nil
}

func (s *LocalEventQueueService) Create(ctx context.Context, req dto.EventQueueSaveRequest) (*dto.EventQueueItem, error) {
	if req.EventType == "" || req.EventName == "" {
		return nil, _errors.BadRequestException("eventType and eventName are required")
	}
	eq := &entity.EventQueue{
		EventType:  req.EventType,
		EventName:  req.EventName,
		EventData:  req.EventData,
		Metadata:   req.Metadata,
		Status:     entity.EEventQueueStatusPending,
		MaxRetries: 3,
	}
	if req.Status != nil {
		eq.Status = *req.Status
	}
	if req.ProfileID != nil {
		eq.ProfileID = req.ProfileID
	}
	if req.OrganizationID != nil {
		eq.OrganizationID = req.OrganizationID
	}
	if req.MaxRetries != nil {
		eq.MaxRetries = *req.MaxRetries
	}
	if req.ScheduledAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err == nil {
			eq.ScheduledAt = &t
		}
	}
	if err := s.repo.Create(ctx, eq); err != nil {
		return nil, err
	}
	return mapEventQueueItem(eq), nil
}

func (s *LocalEventQueueService) Update(ctx context.Context, id uint64, req dto.EventQueueSaveRequest) (*dto.EventQueueItem, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, _errors.NotFoundException("event queue not found")
	}
	existing.EventType = req.EventType
	existing.EventName = req.EventName
	existing.EventData = req.EventData
	existing.Metadata = req.Metadata
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.ProfileID != nil {
		existing.ProfileID = req.ProfileID
	}
	if req.OrganizationID != nil {
		existing.OrganizationID = req.OrganizationID
	}
	if req.MaxRetries != nil {
		existing.MaxRetries = *req.MaxRetries
	}
	if req.RetryCount != nil {
		existing.RetryCount = *req.RetryCount
	}
	if req.ErrorMessage != nil {
		existing.ErrorMessage = req.ErrorMessage
	}
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return mapEventQueueItem(existing), nil
}

func (s *LocalEventQueueService) Delete(ctx context.Context, id uint64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return _errors.NotFoundException("event queue not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *LocalEventQueueService) UpdateStatus(ctx context.Context, req dto.EventQueueStatusRequest) (*dto.EventQueueItem, error) {
	if !req.Status.IsValid() {
		return nil, _errors.BadRequestException("invalid status")
	}
	row, err := s.repo.UpdateStatus(ctx, req.ID, req.Status, req.ErrorMessage)
	if err != nil {
		return nil, err
	}
	return mapEventQueueItem(row), nil
}

func mapEventQueueItem(eq *entity.EventQueue) *dto.EventQueueItem {
	if eq == nil {
		return nil
	}
	item := &dto.EventQueueItem{
		ID: eq.ID, EventType: eq.EventType, EventName: eq.EventName,
		EventData: eq.EventData, Metadata: eq.Metadata, Status: eq.Status,
		ProfileID: eq.ProfileID, OrganizationID: eq.OrganizationID,
		RetryCount: eq.RetryCount, MaxRetries: eq.MaxRetries, ErrorMessage: eq.ErrorMessage,
	}
	if eq.ScheduledAt != nil {
		s := eq.ScheduledAt.Format(time.RFC3339)
		item.ScheduledAt = &s
	}
	if eq.ProcessedAt != nil {
		s := eq.ProcessedAt.Format(time.RFC3339)
		item.ProcessedAt = &s
	}
	if eq.CreatedAt != nil {
		s := strconv.FormatInt(eq.CreatedAt.UnixMilli(), 10)
		item.CreatedAt = &s
	}
	if eq.UpdatedAt != nil {
		s := strconv.FormatInt(eq.UpdatedAt.UnixMilli(), 10)
		item.UpdatedAt = &s
	}
	return item
}

func mapEventQueueItems(rows []*entity.EventQueue) []*dto.EventQueueItem {
	out := make([]*dto.EventQueueItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapEventQueueItem(r))
	}
	return out
}
