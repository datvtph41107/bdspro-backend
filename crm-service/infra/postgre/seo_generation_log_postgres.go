package postgre

import (
	"context"
	"crm/infra/impl"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/repo"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SeoGenerationLogPostgres struct {
	db *gorm.DB
}

func NewSeoGenerationLogPostgres(db *gorm.DB) repo.SeoGenerationLogRepo {
	return &SeoGenerationLogPostgres{db: db}
}

func (r *SeoGenerationLogPostgres) Create(ctx context.Context, log *seo_domain.SeoGenerationLog) (*seo_domain.SeoGenerationLog, error) {
	if err := impl.GetDB(ctx, r.db).WithContext(ctx).Create(log).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, log.ID)
}

func (r *SeoGenerationLogPostgres) GetByID(ctx context.Context, id uint64) (*seo_domain.SeoGenerationLog, error) {
	var log seo_domain.SeoGenerationLog
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Preload("SeoDomain").
		First(&log, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

func (r *SeoGenerationLogPostgres) GetList(ctx context.Context, req *dto.SeoGenerationLogListRequest) ([]seo_domain.SeoGenerationLog, int64, error) {
	var items []seo_domain.SeoGenerationLog
	var total int64

	if req == nil {
		req = &dto.SeoGenerationLogListRequest{}
	}
	req.Normalize()

	query := impl.GetDB(ctx, r.db).WithContext(ctx).Model(&seo_domain.SeoGenerationLog{})

	if req.SeoDomainID != nil {
		query = query.Where("seo_domain_id = ?", *req.SeoDomainID)
	}
	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		query = query.Where("status = ?", strings.TrimSpace(*req.Status))
	}
	if req.TriggerType != nil && strings.TrimSpace(*req.TriggerType) != "" {
		query = query.Where("trigger_type = ?", strings.TrimSpace(*req.TriggerType))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimitAdmin()).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *SeoGenerationLogPostgres) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = gorm.Expr("NOW()")

	return impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Model(&seo_domain.SeoGenerationLog{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *SeoGenerationLogPostgres) MarkRunning(ctx context.Context, id uint64, startedAt time.Time) error {
	return r.Update(ctx, id, map[string]interface{}{
		"status":     seo_domain.SeoGenerationStatusRunning,
		"started_at": startedAt,
	})
}

func (r *SeoGenerationLogPostgres) MarkSuccess(ctx context.Context, id uint64, staticPath string, staticHash string, finishedAt time.Time) error {
	return r.Update(ctx, id, map[string]interface{}{
		"status":           seo_domain.SeoGenerationStatusSuccess,
		"static_html_path": staticPath,
		"static_html_hash": staticHash,
		"finished_at":      finishedAt,
		"error_message":    "",
	})
}

func (r *SeoGenerationLogPostgres) MarkFailed(ctx context.Context, id uint64, errorMessage string, finishedAt time.Time) error {
	return r.Update(ctx, id, map[string]interface{}{
		"status":        seo_domain.SeoGenerationStatusFailed,
		"error_message": errorMessage,
		"finished_at":   finishedAt,
	})
}
