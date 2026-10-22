package repository

import (
	"context"
	"errors"
	"time"

	"mdn4/modules/hub/entity"

	"gorm.io/gorm"
)

// EventQueueRepository manages tb_event_queue persistence.
type EventQueueRepository interface {
	Search(ctx context.Context, text string, page, size int) ([]*entity.EventQueue, int64, error)
	GetByID(ctx context.Context, id uint64) (*entity.EventQueue, error)
	Create(ctx context.Context, eq *entity.EventQueue) error
	Update(ctx context.Context, eq *entity.EventQueue) error
	Delete(ctx context.Context, id uint64) error
	UpdateStatus(ctx context.Context, id uint64, status entity.EEventQueueStatus, errMsg *string) (*entity.EventQueue, error)
	GetPendingEvents(ctx context.Context, limit int) ([]*entity.EventQueue, error)
}

type eventQueuePostgres struct{ db *gorm.DB }

func NewEventQueuePostgres(db *gorm.DB) EventQueueRepository {
	return &eventQueuePostgres{db: db}
}

func (r *eventQueuePostgres) Search(ctx context.Context, text string, page, size int) ([]*entity.EventQueue, int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.EventQueue{}).Where("deleted_at IS NULL")
	if text != "" {
		like := "%" + text + "%"
		q = q.Where("event_name ILIKE ? OR event_type ILIKE ?", like, like)
	}
	var total int64
	q.Count(&total)
	offset, size := pageOffset(page, size)
	var rows []*entity.EventQueue
	err := q.Order("created_at DESC").Offset(offset).Limit(size).Find(&rows).Error
	return rows, total, err
}

func (r *eventQueuePostgres) GetByID(ctx context.Context, id uint64) (*entity.EventQueue, error) {
	var eq entity.EventQueue
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&eq).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &eq, err
}

func (r *eventQueuePostgres) Create(ctx context.Context, eq *entity.EventQueue) error {
	return r.db.WithContext(ctx).Create(eq).Error
}

func (r *eventQueuePostgres) Update(ctx context.Context, eq *entity.EventQueue) error {
	return r.db.WithContext(ctx).Save(eq).Error
}

func (r *eventQueuePostgres) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.EventQueue{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *eventQueuePostgres) UpdateStatus(ctx context.Context, id uint64, status entity.EEventQueueStatus, errMsg *string) (*entity.EventQueue, error) {
	updates := map[string]interface{}{"status": status}
	if status == entity.EEventQueueStatusCompleted || status == entity.EEventQueueStatusFailed {
		now := time.Now()
		updates["processed_at"] = &now
	}
	if errMsg != nil {
		updates["error_message"] = errMsg
	}
	if err := r.db.WithContext(ctx).Model(&entity.EventQueue{}).
		Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *eventQueuePostgres) GetPendingEvents(ctx context.Context, limit int) ([]*entity.EventQueue, error) {
	var rows []*entity.EventQueue
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&entity.EventQueue{}).
		Where("deleted_at IS NULL").
		Where("status = ?", entity.EEventQueueStatusPending).
		Where("(scheduled_at IS NULL OR scheduled_at <= ?)", now).
		Where("retry_count < max_retries").
		Order("created_at ASC").Limit(limit).Find(&rows).Error
	return rows, err
}
