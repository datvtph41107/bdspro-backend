package postgres

import (
	_db "common/db"
	_provider "common/provider"
	_utils "common/utils"
	"context"
	"errors"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/enums"
	"time"

	"gorm.io/gorm"
)

// @bind: hub/internal/repo.IEventQueueRepo
type EventQueueRepo struct {
	_provider.CrudRepo[domain.EventQueueEntity]
}

func NewEventQueueRepo(db *_db.TransactionRepo) *EventQueueRepo {
	repo := &EventQueueRepo{}
	repo.Init(repo, db)
	return repo
}

// GetByID normalizes only record-not-found to absence so the application layer
// owns the business meaning; unrelated database failures remain infrastructure errors.
func (r *EventQueueRepo) GetByID(c context.Context, id uint64) (*domain.EventQueueEntity, error) {
	var entity domain.EventQueueEntity
	err := r.GetDB(c).
		Where("id = ? and deleted_at is null", id).
		First(&entity).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *EventQueueRepo) Search(c context.Context, searchDTO dto.EventQueueSearchDTO) ([]domain.EventQueueEntity, int64, error) {
	var events []domain.EventQueueEntity

	query := r.GetDB(c).Model(&domain.EventQueueEntity{}).
		Where("deleted_at IS NULL")

	// Filter by text (search in event_name, event_type)
	if searchDTO.Text != nil && *searchDTO.Text != "" {
		query = query.Where("(event_name ILIKE ? OR event_type ILIKE ?)", "%"+*searchDTO.Text+"%", "%"+*searchDTO.Text+"%")
	}

	// Filter by profileId
	if searchDTO.ProfileID != nil {
		query = query.Where("profile_id = ?", *searchDTO.ProfileID)
	}

	// Filter by organizationId
	if searchDTO.OrganizationID != nil {
		query = query.Where("organization_id = ?", *searchDTO.OrganizationID)
	}

	// Filter by status
	if len(searchDTO.Status) > 0 {
		query = query.Where("status IN (?)", searchDTO.Status)
	}

	// Filter by event type
	if len(searchDTO.EventType) > 0 {
		query = query.Where("event_type IN (?)", searchDTO.EventType)
	}

	// Filter by date range
	if searchDTO.FromDate != nil && *searchDTO.FromDate != "" {
		fromDate := _utils.ParseStringToTime(*searchDTO.FromDate)
		if fromDate != nil {
			query = query.Where("created_at >= ?", fromDate)
		}
	}

	if searchDTO.ToDate != nil && *searchDTO.ToDate != "" {
		toDate := _utils.ParseStringToTime(*searchDTO.ToDate)
		if toDate != nil {
			query = query.Where("created_at <= ?", toDate)
		}
	}

	// Count total
	var total int64
	query.Count(&total)

	// Apply pagination
	result := query.
		Order("created_at DESC").
		Offset(searchDTO.GetOffset()).
		Limit(searchDTO.GetLimit()).
		Find(&events)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return events, total, nil
}

func (r *EventQueueRepo) UpdateStatus(c context.Context, id uint64, status uint32, errorMessage *string) (*domain.EventQueueEntity, error) {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == uint32(enums.EEventQueueStatusCompleted) || status == uint32(enums.EEventQueueStatusFailed) {
		now := time.Now()
		updates["processed_at"] = &now
	}

	if errorMessage != nil {
		updates["error_message"] = errorMessage
	}

	result := r.GetDB(c).
		Model(&domain.EventQueueEntity{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	// Get updated entity
	var event domain.EventQueueEntity
	if err := r.GetDB(c).Where("id = ?", id).First(&event).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventQueueRepo) GetPendingEvents(c context.Context, limit int) ([]domain.EventQueueEntity, error) {
	var events []domain.EventQueueEntity

	now := time.Now()

	result := r.GetDB(c).
		Model(&domain.EventQueueEntity{}).
		Where("deleted_at IS NULL").
		Where("status = ?", enums.EEventQueueStatusPending).
		Where("(scheduled_at IS NULL OR scheduled_at <= ?)", now).
		Where("retry_count < max_retries").
		Order("created_at ASC").
		Limit(limit).
		Find(&events)

	if result.Error != nil {
		return nil, result.Error
	}

	return events, nil
}

func (r *EventQueueRepo) BeforeSave(c context.Context, id *uint64, entity *domain.EventQueueEntity) error {
	// Validation logic if needed
	return nil
}

func (r *EventQueueRepo) AfterSave(c context.Context, id *uint64, entity *domain.EventQueueEntity) error {
	// Post-save logic if needed (cache, audit, etc.)
	return nil
}
