package postgre

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm/infra/impl"
	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/repo"

	"gorm.io/gorm"
)

type SeoAudiencePostgres struct {
	db *gorm.DB
}

var _ repo.SeoAudienceRepo = (*SeoAudiencePostgres)(nil)

func NewSeoAudiencePostgres(
	db *gorm.DB,
) repo.SeoAudienceRepo {
	return &SeoAudiencePostgres{
		db: db,
	}
}

type seoAudienceRow struct {
	ID             uint64     `gorm:"column:id;primaryKey"`
	Code           string     `gorm:"column:code"`
	Name           string     `gorm:"column:name"`
	Description    string     `gorm:"column:description"`
	DecisionRole   string     `gorm:"column:decision_role"`
	KnowledgeLevel string     `gorm:"column:knowledge_level"`
	Status         string     `gorm:"column:status"`
	OwnerCode      string     `gorm:"column:owner_code"`
	ReviewDueAt    *time.Time `gorm:"column:review_due_at"`
	LockVersion    uint32     `gorm:"column:lock_version"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
}

func (seoAudienceRow) TableName() string { return "seo_audience_segment" }

func (r *SeoAudiencePostgres) Create(ctx context.Context, audience *seo_strategy.SeoAudienceSegment) error {
	row := audienceToRow(audience)
	if err := impl.GetDB(ctx, r.db).WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	*audience = *audienceRowToDomain(&row)
	return nil
}

func (r *SeoAudiencePostgres) GetByID(ctx context.Context, id uint64) (*seo_strategy.SeoAudienceSegment, error) {
	var row seoAudienceRow
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoAudienceNotFound
	}
	if err != nil {
		return nil, err
	}
	return audienceRowToDomain(&row), nil
}

func (r *SeoAudiencePostgres) GetByCode(ctx context.Context, code string) (*seo_strategy.SeoAudienceSegment, error) {
	var row seoAudienceRow
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Where("code = ? AND deleted_at IS NULL", seo_strategy.NormalizeStrategyCode(code)).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoAudienceNotFound
	}
	if err != nil {
		return nil, err
	}
	return audienceRowToDomain(&row), nil
}

func (r *SeoAudiencePostgres) List(ctx context.Context, filter repo.SeoAudienceFilter) ([]seo_strategy.SeoAudienceSegment, int64, error) {
	size := filter.Size
	if size == 0 || size > 100 {
		size = 20
	}
	query := impl.GetDB(ctx, r.db).WithContext(ctx).
		Model(&seoAudienceRow{}).Where("deleted_at IS NULL")
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("code ILIKE ? OR name ILIKE ? OR description ILIKE ?", like, like, like)
	}
	if status := strings.TrimSpace(filter.Status); status != "" {
		query = query.Where("status = ?", status)
	}
	if role := strings.TrimSpace(filter.DecisionRole); role != "" {
		query = query.Where("decision_role = ?", role)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []seoAudienceRow
	if err := query.Order("updated_at DESC, id DESC").
		Offset(int(filter.Page * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	items := make([]seo_strategy.SeoAudienceSegment, 0, len(rows))
	for i := range rows {
		items = append(items, *audienceRowToDomain(&rows[i]))
	}
	return items, total, nil
}

func (r *SeoAudiencePostgres) UpdateDefinition(ctx context.Context, audience *seo_strategy.SeoAudienceSegment, expectedVersion uint32) error {
	result := impl.GetDB(ctx, r.db).WithContext(ctx).Model(&seoAudienceRow{}).
		Where("id = ? AND lock_version = ? AND deleted_at IS NULL", audience.ID, expectedVersion).
		Updates(map[string]interface{}{
			"name": audience.Name, "description": audience.Description,
			"decision_role": string(audience.DecisionRole), "knowledge_level": string(audience.KnowledgeLevel),
			"owner_code": audience.OwnerCode, "review_due_at": audience.ReviewDueAt,
			"lock_version": gorm.Expr("lock_version + 1"), "updated_at": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return r.resolveUpdateFailure(ctx, audience.ID)
	}
	return nil
}

func (r *SeoAudiencePostgres) UpdateStatus(ctx context.Context, id uint64, status seo_strategy.AudienceStatus, expectedVersion uint32) error {
	result := impl.GetDB(ctx, r.db).WithContext(ctx).Model(&seoAudienceRow{}).
		Where("id = ? AND lock_version = ? AND deleted_at IS NULL", id, expectedVersion).
		Updates(map[string]interface{}{
			"status": string(status), "lock_version": gorm.Expr("lock_version + 1"), "updated_at": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return r.resolveUpdateFailure(ctx, id)
	}
	return nil
}

func (r *SeoAudiencePostgres) resolveUpdateFailure(ctx context.Context, id uint64) error {
	var count int64
	if err := impl.GetDB(ctx, r.db).WithContext(ctx).Model(&seoAudienceRow{}).
		Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return repo.ErrSeoAudienceNotFound
	}
	return repo.ErrSeoAudienceConflict
}

func audienceToRow(item *seo_strategy.SeoAudienceSegment) seoAudienceRow {
	return seoAudienceRow{
		ID: item.ID, Code: item.Code, Name: item.Name, Description: item.Description,
		DecisionRole: string(item.DecisionRole), KnowledgeLevel: string(item.KnowledgeLevel),
		Status: string(item.Status), OwnerCode: item.OwnerCode, ReviewDueAt: item.ReviewDueAt,
		LockVersion: item.LockVersion, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt, DeletedAt: item.DeletedAt,
	}
}

func audienceRowToDomain(row *seoAudienceRow) *seo_strategy.SeoAudienceSegment {
	if row == nil {
		return nil
	}
	return &seo_strategy.SeoAudienceSegment{
		ID: row.ID, Code: row.Code, Name: row.Name, Description: row.Description,
		DecisionRole:   seo_strategy.AudienceDecisionRole(row.DecisionRole),
		KnowledgeLevel: seo_strategy.AudienceKnowledgeLevel(row.KnowledgeLevel),
		Status:         seo_strategy.AudienceStatus(row.Status), OwnerCode: row.OwnerCode, ReviewDueAt: row.ReviewDueAt,
		LockVersion: row.LockVersion, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt, DeletedAt: row.DeletedAt,
	}
}
