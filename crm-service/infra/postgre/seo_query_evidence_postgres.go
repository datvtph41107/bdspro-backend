package postgre

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"crm/infra/impl"
	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/repo"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SeoQueryEvidencePostgres struct {
	db *gorm.DB
}

var _ repo.SeoQueryEvidenceRepo = (*SeoQueryEvidencePostgres)(nil)

func NewSeoQueryEvidencePostgres(db *gorm.DB) repo.SeoQueryEvidenceRepo {
	return &SeoQueryEvidencePostgres{db: db}
}

type seoQueryEvidenceRow struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	ObjectiveID uint64 `gorm:"column:objective_id"`
	UserNeedID  uint64 `gorm:"column:user_need_id"`

	RawQuery        string `gorm:"column:raw_query"`
	NormalizedQuery string `gorm:"column:normalized_query"`
	QueryHash       string `gorm:"column:query_hash"`

	SourceSystem    string `gorm:"column:source_system"`
	SourceKind      string `gorm:"column:source_kind"`
	ObservationType string `gorm:"column:observation_type"`
	SourceRecordKey string `gorm:"column:source_record_key"`
	SourceReference string `gorm:"column:source_reference"`

	Locale         string `gorm:"column:locale"`
	CountryCode    string `gorm:"column:country_code"`
	DeviceCategory string `gorm:"column:device_category"`

	ObservedFrom *time.Time `gorm:"column:observed_from"`
	ObservedTo   *time.Time `gorm:"column:observed_to"`
	ValidUntil   *time.Time `gorm:"column:valid_until"`

	Impressions  *int64   `gorm:"column:impressions"`
	Clicks       *int64   `gorm:"column:clicks"`
	CTR          *float64 `gorm:"column:ctr"`
	AvgPosition  *float64 `gorm:"column:avg_position"`
	SearchVolume *int64   `gorm:"column:search_volume"`

	Confidence string `gorm:"column:confidence"`
	Status     string `gorm:"column:status"`

	Metadata datatypes.JSON `gorm:"column:metadata"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (seoQueryEvidenceRow) TableName() string {
	return "seo_query_evidence"
}

func (r *SeoQueryEvidencePostgres) Upsert(
	ctx context.Context,
	evidence *seo_strategy.SeoQueryEvidence,
) error {
	if evidence == nil {
		return repo.ErrSeoQueryEvidenceInvalid
	}
	if err := seo_strategy.ValidateQueryEvidence(evidence); err != nil {
		return fmt.Errorf("%w: %v", repo.ErrSeoQueryEvidenceInvalid, err)
	}

	row := queryEvidenceToRow(evidence)
	var persisted seoQueryEvidenceRow

	result := impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Raw(
			`
			INSERT INTO seo_query_evidence (
				objective_id,
				user_need_id,
				raw_query,
				normalized_query,
				query_hash,
				source_system,
				source_kind,
				observation_type,
				source_record_key,
				source_reference,
				locale,
				country_code,
				device_category,
				observed_from,
				observed_to,
				valid_until,
				impressions,
				clicks,
				ctr,
				avg_position,
				search_volume,
				confidence,
				status,
				metadata,
				created_at,
				updated_at
			)
			VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)
			ON CONFLICT (
				user_need_id,
				source_system,
				source_record_key
			)
			WHERE status <> 'rejected'
			DO UPDATE SET
				objective_id = EXCLUDED.objective_id,
				raw_query = EXCLUDED.raw_query,
				normalized_query = EXCLUDED.normalized_query,
				query_hash = EXCLUDED.query_hash,
				source_kind = EXCLUDED.source_kind,
				observation_type = EXCLUDED.observation_type,
				source_reference = EXCLUDED.source_reference,
				locale = EXCLUDED.locale,
				country_code = EXCLUDED.country_code,
				device_category = EXCLUDED.device_category,
				observed_from = EXCLUDED.observed_from,
				observed_to = EXCLUDED.observed_to,
				valid_until = EXCLUDED.valid_until,
				impressions = EXCLUDED.impressions,
				clicks = EXCLUDED.clicks,
				ctr = EXCLUDED.ctr,
				avg_position = EXCLUDED.avg_position,
				search_volume = EXCLUDED.search_volume,
				confidence = EXCLUDED.confidence,
				metadata = EXCLUDED.metadata,
				updated_at = NOW()
			RETURNING
				id,
				objective_id,
				user_need_id,
				raw_query,
				normalized_query,
				query_hash,
				source_system,
				source_kind,
				observation_type,
				source_record_key,
				source_reference,
				locale,
				country_code,
				device_category,
				observed_from,
				observed_to,
				valid_until,
				impressions,
				clicks,
				ctr,
				avg_position,
				search_volume,
				confidence,
				status,
				metadata,
				created_at,
				updated_at
			`,
			row.ObjectiveID,
			row.UserNeedID,
			row.RawQuery,
			row.NormalizedQuery,
			row.QueryHash,
			row.SourceSystem,
			row.SourceKind,
			row.ObservationType,
			row.SourceRecordKey,
			row.SourceReference,
			row.Locale,
			row.CountryCode,
			row.DeviceCategory,
			row.ObservedFrom,
			row.ObservedTo,
			row.ValidUntil,
			row.Impressions,
			row.Clicks,
			row.CTR,
			row.AvgPosition,
			row.SearchVolume,
			row.Confidence,
			row.Status,
			row.Metadata,
			row.CreatedAt,
			row.UpdatedAt,
		).
		Scan(&persisted)

	if result.Error != nil {
		return mapSeoQueryEvidenceDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return repo.ErrSeoQueryEvidenceConflict
	}

	domainItem, err := queryEvidenceRowToDomain(&persisted)
	if err != nil {
		return err
	}
	*evidence = *domainItem
	return nil
}

func (r *SeoQueryEvidencePostgres) GetByID(
	ctx context.Context,
	id uint64,
) (*seo_strategy.SeoQueryEvidence, error) {
	if id == 0 {
		return nil, repo.ErrSeoQueryEvidenceNotFound
	}

	var row seoQueryEvidenceRow
	err := impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Where("id = ?", id).
		First(&row).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoQueryEvidenceNotFound
	}
	if err != nil {
		return nil, err
	}
	return queryEvidenceRowToDomain(&row)
}

func (r *SeoQueryEvidencePostgres) GetBySourceRecordKey(
	ctx context.Context,
	userNeedID uint64,
	sourceSystem string,
	sourceRecordKey string,
) (*seo_strategy.SeoQueryEvidence, error) {
	sourceSystem = strings.TrimSpace(sourceSystem)
	sourceRecordKey = strings.TrimSpace(sourceRecordKey)
	if userNeedID == 0 || sourceSystem == "" || sourceRecordKey == "" {
		return nil, repo.ErrSeoQueryEvidenceNotFound
	}

	var row seoQueryEvidenceRow
	err := impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Where(
			`user_need_id = ?
			 AND source_system = ?
			 AND source_record_key = ?
			 AND status <> 'rejected'`,
			userNeedID,
			sourceSystem,
			sourceRecordKey,
		).
		First(&row).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoQueryEvidenceNotFound
	}
	if err != nil {
		return nil, err
	}
	return queryEvidenceRowToDomain(&row)
}

func (r *SeoQueryEvidencePostgres) List(
	ctx context.Context,
	filter repo.SeoQueryEvidenceFilter,
) ([]seo_strategy.SeoQueryEvidence, int64, error) {
	size := filter.Size
	if size == 0 || size > 100 {
		size = 20
	}

	query := impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Model(&seoQueryEvidenceRow{})

	if filter.ObjectiveID > 0 {
		query = query.Where("objective_id = ?", filter.ObjectiveID)
	}
	if filter.UserNeedID > 0 {
		query = query.Where("user_need_id = ?", filter.UserNeedID)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		normalizedKeyword := seo_strategy.NormalizeSearchQuery(keyword)
		rawLike := "%" + keyword + "%"
		normalizedLike := "%" + normalizedKeyword + "%"
		query = query.Where(
			`raw_query ILIKE ?
			 OR normalized_query ILIKE ?
			 OR source_reference ILIKE ?`,
			rawLike,
			normalizedLike,
			rawLike,
		)
	}
	if sourceSystem := strings.TrimSpace(filter.SourceSystem); sourceSystem != "" {
		query = query.Where("source_system = ?", sourceSystem)
	}
	if filter.SourceKind != "" {
		query = query.Where("source_kind = ?", string(filter.SourceKind))
	}
	if filter.ObservationType != "" {
		query = query.Where("observation_type = ?", string(filter.ObservationType))
	}
	if filter.DeviceCategory != "" {
		query = query.Where("device_category = ?", string(filter.DeviceCategory))
	}
	if filter.Confidence != "" {
		query = query.Where("confidence = ?", string(filter.Confidence))
	}
	if filter.Status != "" {
		query = query.Where("status = ?", string(filter.Status))
	}
	if filter.CurrentAt != nil {
		query = query.
			Where("status = ?", string(seo_strategy.QueryEvidenceStatusActive)).
			Where("valid_until IS NULL OR valid_until >= ?", *filter.CurrentAt)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []seoQueryEvidenceRow
	if err := query.
		Order("observed_to DESC NULLS LAST, updated_at DESC, id DESC").
		Offset(int(filter.Page * size)).
		Limit(int(size)).
		Find(&rows).
		Error; err != nil {
		return nil, 0, err
	}

	items := make([]seo_strategy.SeoQueryEvidence, 0, len(rows))
	for i := range rows {
		item, err := queryEvidenceRowToDomain(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}
	return items, total, nil
}

func (r *SeoQueryEvidencePostgres) UpdateStatus(
	ctx context.Context,
	id uint64,
	expectedStatus seo_strategy.QueryEvidenceStatus,
	targetStatus seo_strategy.QueryEvidenceStatus,
) error {
	if id == 0 {
		return repo.ErrSeoQueryEvidenceNotFound
	}
	if !seo_strategy.IsValidQueryEvidenceStatus(expectedStatus) {
		return fmt.Errorf("%w: expected status %s", repo.ErrSeoQueryEvidenceInvalid, expectedStatus)
	}
	if !seo_strategy.IsValidQueryEvidenceStatus(targetStatus) {
		return fmt.Errorf("%w: target status %s", repo.ErrSeoQueryEvidenceInvalid, targetStatus)
	}

	result := impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Model(&seoQueryEvidenceRow{}).
		Where("id = ? AND status = ?", id, string(expectedStatus)).
		Updates(map[string]interface{}{
			"status":     string(targetStatus),
			"updated_at": gorm.Expr("NOW()"),
		})

	if result.Error != nil {
		return mapSeoQueryEvidenceDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return r.resolveQueryEvidenceStatusUpdateFailure(ctx, id)
	}
	return nil
}

func (r *SeoQueryEvidencePostgres) resolveQueryEvidenceStatusUpdateFailure(
	ctx context.Context,
	id uint64,
) error {
	var row struct {
		Status string `gorm:"column:status"`
	}

	err := impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Model(&seoQueryEvidenceRow{}).
		Select("status").
		Where("id = ?", id).
		Take(&row).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repo.ErrSeoQueryEvidenceNotFound
	}
	if err != nil {
		return err
	}
	return repo.ErrSeoQueryEvidenceConflict
}

func mapSeoQueryEvidenceDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23503":
		return fmt.Errorf(
			"%w: %s",
			repo.ErrSeoQueryEvidenceUserNeedMismatch,
			pgErr.ConstraintName,
		)
	case "23514", "22001", "22P02":
		return fmt.Errorf(
			"%w: %s",
			repo.ErrSeoQueryEvidenceInvalid,
			pgErr.ConstraintName,
		)
	case "23505":
		return fmt.Errorf(
			"%w: %s",
			repo.ErrSeoQueryEvidenceConflict,
			pgErr.ConstraintName,
		)
	default:
		return err
	}
}

func queryEvidenceToRow(item *seo_strategy.SeoQueryEvidence) seoQueryEvidenceRow {
	metadata := item.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}

	return seoQueryEvidenceRow{
		ID:              item.ID,
		ObjectiveID:     item.ObjectiveID,
		UserNeedID:      item.UserNeedID,
		RawQuery:        item.RawQuery,
		NormalizedQuery: item.NormalizedQuery,
		QueryHash:       item.QueryHash,
		SourceSystem:    item.SourceSystem,
		SourceKind:      string(item.SourceKind),
		ObservationType: string(item.ObservationType),
		SourceRecordKey: item.SourceRecordKey,
		SourceReference: item.SourceReference,
		Locale:          item.Locale,
		CountryCode:     item.CountryCode,
		DeviceCategory:  string(item.DeviceCategory),
		ObservedFrom:    item.ObservedFrom,
		ObservedTo:      item.ObservedTo,
		ValidUntil:      item.ValidUntil,
		Impressions:     item.Impressions,
		Clicks:          item.Clicks,
		CTR:             item.CTR,
		AvgPosition:     item.AvgPosition,
		SearchVolume:    item.SearchVolume,
		Confidence:      string(item.Confidence),
		Status:          string(item.Status),
		Metadata:        datatypes.JSON(metadata),
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}
}

func queryEvidenceRowToDomain(row *seoQueryEvidenceRow) (*seo_strategy.SeoQueryEvidence, error) {
	if row == nil {
		return nil, nil
	}

	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}

	item := &seo_strategy.SeoQueryEvidence{
		ID:              row.ID,
		ObjectiveID:     row.ObjectiveID,
		UserNeedID:      row.UserNeedID,
		RawQuery:        row.RawQuery,
		NormalizedQuery: row.NormalizedQuery,
		QueryHash:       row.QueryHash,
		SourceSystem:    row.SourceSystem,
		SourceKind:      seo_strategy.QueryEvidenceSourceKind(row.SourceKind),
		ObservationType: seo_strategy.QueryObservationType(row.ObservationType),
		SourceRecordKey: row.SourceRecordKey,
		SourceReference: row.SourceReference,
		Locale:          row.Locale,
		CountryCode:     row.CountryCode,
		DeviceCategory:  seo_strategy.QueryDeviceCategory(row.DeviceCategory),
		ObservedFrom:    row.ObservedFrom,
		ObservedTo:      row.ObservedTo,
		ValidUntil:      row.ValidUntil,
		Impressions:     row.Impressions,
		Clicks:          row.Clicks,
		CTR:             row.CTR,
		AvgPosition:     row.AvgPosition,
		SearchVolume:    row.SearchVolume,
		Confidence:      seo_strategy.ObjectiveConfidence(row.Confidence),
		Status:          seo_strategy.QueryEvidenceStatus(row.Status),
		Metadata:        metadata,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}

	if err := seo_strategy.ValidateQueryEvidence(item); err != nil {
		return nil, fmt.Errorf(
			"seo_query_evidence id=%d không hợp lệ: %w",
			row.ID,
			err,
		)
	}
	return item, nil
}
