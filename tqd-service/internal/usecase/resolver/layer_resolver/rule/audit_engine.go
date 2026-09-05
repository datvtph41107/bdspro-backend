package rule

import (
	"context"
	"fmt"
	"time"

	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// AuditRepository — interface cho Audit Engine
type AuditRepository interface {
	Insert(ctx context.Context, entry *types.AuditEntry) error
	Query(ctx context.Context, q *types.AuditQuery) ([]*types.AuditEntry, int64, error)
	QueryByParcel(ctx context.Context, parcelID uint64) ([]*types.AuditEntry, error)
}

// AuditEngine — Engine ghi nhận và truy vấn audit log
type AuditEngine struct {
	config    *types.AuditConfig
	auditRepo AuditRepository
}

// NewAuditEngine tạo AuditEngine
func NewAuditEngine(config *types.AuditConfig, auditRepo AuditRepository) *AuditEngine {
	if config == nil {
		config = types.DefaultAuditConfig()
	}
	return &AuditEngine{config: config, auditRepo: auditRepo}
}

// Name trả về tên engine
func (e *AuditEngine) Name() string { return "AuditEngine" }

// RecordChange ghi nhận 1 thay đổi
func (e *AuditEngine) RecordChange(ctx context.Context, entry *types.AuditEntry) error {
	entry.Timestamp = time.Now()

	if actor, ok := ctx.Value("actor").(string); ok {
		entry.Actor = actor
	}
	if requestID, ok := ctx.Value("request_id").(string); ok {
		entry.RequestID = requestID
	}

	return e.auditRepo.Insert(ctx, entry)
}

// RecordWithDiff ghi nhận thay đổi kèm diff tự động
func (e *AuditEngine) RecordWithDiff(
	ctx context.Context,
	entityType string, entityID uint64,
	action types.AuditAction,
	oldObj, newObj any,
	reason string,
) error {
	diffs := e.computeDiff(oldObj, newObj)
	return e.RecordChange(ctx, &types.AuditEntry{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Changes:    diffs,
		Reason:     reason,
	})
}

// QueryHistory truy vấn lịch sử thay đổi
func (e *AuditEngine) QueryHistory(ctx context.Context, q *types.AuditQuery) ([]*types.AuditEntry, int64, error) {
	return e.auditRepo.Query(ctx, q)
}

// GetChangeTimeline lấy timeline thay đổi của 1 entity
func (e *AuditEngine) GetChangeTimeline(ctx context.Context, entityType string, entityID uint64) ([]*types.AuditEntry, error) {
	entries, _, err := e.auditRepo.Query(ctx, &types.AuditQuery{
		EntityType: entityType,
		EntityID:   entityID,
		PageSize:   100,
	})
	return entries, err
}

// computeDiff so sánh 2 object, trả về danh sách field thay đổi
func (e *AuditEngine) computeDiff(oldObj, newObj any) []types.FieldDiff {
	// Simplified: trả về empty, implementation thực tế dùng reflection hoặc library
	_ = oldObj
	_ = newObj
	return nil
}

// IsActionEnabled kiểm tra action có được phép audit không
func (e *AuditEngine) IsActionEnabled(action types.AuditAction) bool {
	for _, a := range e.config.EnabledActions {
		if a == action {
			return true
		}
	}
	return false
}

// Validate kiểm tra config
func (e *AuditEngine) Validate() error {
	if e.auditRepo == nil {
		return fmt.Errorf("AuditEngine: auditRepo is nil")
	}
	return nil
}
