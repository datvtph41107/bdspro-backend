package postgre

import (
	"context"
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

type SeoIntentHypothesisPostgres struct {
	db *gorm.DB
}

var _ repo.SeoIntentHypothesisRepo = (*SeoIntentHypothesisPostgres)(nil)

func NewSeoIntentHypothesisPostgres(db *gorm.DB) repo.SeoIntentHypothesisRepo {
	return &SeoIntentHypothesisPostgres{db: db}
}

type seoIntentHypothesisRow struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	ObjectiveID uint64 `gorm:"column:objective_id"`
	UserNeedID  uint64 `gorm:"column:user_need_id"`

	Code string `gorm:"column:code"`
	Name string `gorm:"column:name"`

	GoalStatement   string `gorm:"column:goal_statement"`
	PrimaryQuestion string `gorm:"column:primary_question"`
	ExpectedResult  string `gorm:"column:expected_result"`

	SearchAction string `gorm:"column:search_action"`
	SubjectType  string `gorm:"column:subject_type"`
	ScopeType    string `gorm:"column:scope_type"`

	Status     string `gorm:"column:status"`
	Confidence string `gorm:"column:confidence"`

	OwnerCode   string     `gorm:"column:owner_code"`
	ReviewDueAt *time.Time `gorm:"column:review_due_at"`
	LockVersion uint32     `gorm:"column:lock_version"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (seoIntentHypothesisRow) TableName() string {
	return "seo_intent_hypothesis"
}

type seoIntentEvidenceLinkRow struct {
	ObjectiveID uint64 `gorm:"column:objective_id"`
	UserNeedID  uint64 `gorm:"column:user_need_id"`

	IntentHypothesisID uint64 `gorm:"column:intent_hypothesis_id"`
	QueryEvidenceID    uint64 `gorm:"column:query_evidence_id"`

	RelationType  string `gorm:"column:relation_type"`
	RelevanceNote string `gorm:"column:relevance_note"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (seoIntentEvidenceLinkRow) TableName() string {
	return "seo_intent_evidence_link"
}

type linkedIntentQueryEvidenceRow struct {
	LinkObjectiveID     uint64         `gorm:"column:link_objective_id"`
	LinkUserNeedID      uint64         `gorm:"column:link_user_need_id"`
	IntentHypothesisID  uint64         `gorm:"column:intent_hypothesis_id"`
	LinkQueryEvidenceID uint64         `gorm:"column:link_query_evidence_id"`
	RelationType        string         `gorm:"column:relation_type"`
	RelevanceNote       string         `gorm:"column:relevance_note"`
	LinkCreatedAt       time.Time      `gorm:"column:link_created_at"`
	LinkUpdatedAt       time.Time      `gorm:"column:link_updated_at"`
	EvidenceID          uint64         `gorm:"column:evidence_id"`
	EvidenceObjectiveID uint64         `gorm:"column:evidence_objective_id"`
	EvidenceUserNeedID  uint64         `gorm:"column:evidence_user_need_id"`
	RawQuery            string         `gorm:"column:raw_query"`
	NormalizedQuery     string         `gorm:"column:normalized_query"`
	QueryHash           string         `gorm:"column:query_hash"`
	SourceSystem        string         `gorm:"column:source_system"`
	SourceKind          string         `gorm:"column:source_kind"`
	ObservationType     string         `gorm:"column:observation_type"`
	SourceRecordKey     string         `gorm:"column:source_record_key"`
	SourceReference     string         `gorm:"column:source_reference"`
	Locale              string         `gorm:"column:locale"`
	CountryCode         string         `gorm:"column:country_code"`
	DeviceCategory      string         `gorm:"column:device_category"`
	ObservedFrom        *time.Time     `gorm:"column:observed_from"`
	ObservedTo          *time.Time     `gorm:"column:observed_to"`
	ValidUntil          *time.Time     `gorm:"column:valid_until"`
	Impressions         *int64         `gorm:"column:impressions"`
	Clicks              *int64         `gorm:"column:clicks"`
	CTR                 *float64       `gorm:"column:ctr"`
	AvgPosition         *float64       `gorm:"column:avg_position"`
	SearchVolume        *int64         `gorm:"column:search_volume"`
	EvidenceConfidence  string         `gorm:"column:evidence_confidence"`
	EvidenceStatus      string         `gorm:"column:evidence_status"`
	Metadata            datatypes.JSON `gorm:"column:metadata"`
	EvidenceCreatedAt   time.Time      `gorm:"column:evidence_created_at"`
	EvidenceUpdatedAt   time.Time      `gorm:"column:evidence_updated_at"`
}

func (r *SeoIntentHypothesisPostgres) Create(
	ctx context.Context,
	intent *seo_strategy.SeoIntentHypothesis,
) error {
	if intent == nil {
		return repo.ErrSeoIntentHypothesisInvalid
	}
	if err := seo_strategy.ValidateIntentHypothesisDefinition(intent); err != nil {
		return fmt.Errorf("%w: %v", repo.ErrSeoIntentHypothesisInvalid, err)
	}
	row := intentHypothesisToRow(intent)
	if err := impl.GetDB(ctx, r.db).WithContext(ctx).Create(&row).Error; err != nil {
		return mapIntentPostgresError(err)
	}
	*intent = *intentHypothesisRowToDomain(&row)
	return nil
}

func (r *SeoIntentHypothesisPostgres) GetByID(
	ctx context.Context,
	id uint64,
) (*seo_strategy.SeoIntentHypothesis, error) {
	var row seoIntentHypothesisRow
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoIntentHypothesisNotFound
	}
	if err != nil {
		return nil, err
	}
	return intentHypothesisRowToDomain(&row), nil
}

func (r *SeoIntentHypothesisPostgres) GetByCode(
	ctx context.Context,
	userNeedID uint64,
	code string,
) (*seo_strategy.SeoIntentHypothesis, error) {
	var row seoIntentHypothesisRow
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Where(
			"user_need_id = ? AND code = ? AND deleted_at IS NULL",
			userNeedID,
			seo_strategy.NormalizeStrategyCode(code),
		).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoIntentHypothesisNotFound
	}
	if err != nil {
		return nil, err
	}
	return intentHypothesisRowToDomain(&row), nil
}

func (r *SeoIntentHypothesisPostgres) List(
	ctx context.Context,
	filter repo.SeoIntentHypothesisFilter,
) ([]seo_strategy.SeoIntentHypothesis, int64, error) {
	size := filter.Size
	if size == 0 || size > 100 {
		size = 20
	}
	query := impl.GetDB(ctx, r.db).WithContext(ctx).
		Model(&seoIntentHypothesisRow{}).
		Where("deleted_at IS NULL")

	if filter.ObjectiveID > 0 {
		query = query.Where("objective_id = ?", filter.ObjectiveID)
	}
	if filter.UserNeedID > 0 {
		query = query.Where("user_need_id = ?", filter.UserNeedID)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			`code ILIKE ? OR name ILIKE ? OR goal_statement ILIKE ? OR primary_question ILIKE ? OR expected_result ILIKE ?`,
			like, like, like, like, like,
		)
	}
	if filter.SearchAction != "" {
		query = query.Where("search_action = ?", string(filter.SearchAction))
	}
	if filter.SubjectType != "" {
		query = query.Where("subject_type = ?", string(filter.SubjectType))
	}
	if filter.ScopeType != "" {
		query = query.Where("scope_type = ?", string(filter.ScopeType))
	}
	if filter.Status != "" {
		query = query.Where("status = ?", string(filter.Status))
	}
	if filter.Confidence != "" {
		query = query.Where("confidence = ?", string(filter.Confidence))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []seoIntentHypothesisRow
	if err := query.
		Order("updated_at DESC, id DESC").
		Offset(int(filter.Page * size)).
		Limit(int(size)).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]seo_strategy.SeoIntentHypothesis, 0, len(rows))
	for i := range rows {
		items = append(items, *intentHypothesisRowToDomain(&rows[i]))
	}
	return items, total, nil
}

func (r *SeoIntentHypothesisPostgres) UpdateDefinition(
	ctx context.Context,
	intent *seo_strategy.SeoIntentHypothesis,
	expectedVersion uint32,
) error {
	if intent == nil || expectedVersion == 0 {
		return repo.ErrSeoIntentHypothesisInvalid
	}
	if err := seo_strategy.ValidateIntentHypothesisDefinition(intent); err != nil {
		return fmt.Errorf("%w: %v", repo.ErrSeoIntentHypothesisInvalid, err)
	}
	result := impl.GetDB(ctx, r.db).WithContext(ctx).
		Model(&seoIntentHypothesisRow{}).
		Where(
			"id = ? AND lock_version = ? AND deleted_at IS NULL",
			intent.ID,
			expectedVersion,
		).
		Updates(map[string]interface{}{
			"name":             intent.Name,
			"goal_statement":   intent.GoalStatement,
			"primary_question": intent.PrimaryQuestion,
			"expected_result":  intent.ExpectedResult,
			"search_action":    string(intent.SearchAction),
			"subject_type":     string(intent.SubjectType),
			"scope_type":       string(intent.ScopeType),
			"confidence":       string(intent.Confidence),
			"owner_code":       intent.OwnerCode,
			"review_due_at":    intent.ReviewDueAt,
			"lock_version":     gorm.Expr("lock_version + 1"),
			"updated_at":       gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return mapIntentPostgresError(result.Error)
	}
	if result.RowsAffected == 0 {
		return r.resolveIntentUpdateFailure(ctx, intent.ID)
	}
	return nil
}

func (r *SeoIntentHypothesisPostgres) UpdateStatus(
	ctx context.Context,
	id uint64,
	status seo_strategy.IntentHypothesisStatus,
	expectedVersion uint32,
) error {
	result := impl.GetDB(ctx, r.db).WithContext(ctx).
		Model(&seoIntentHypothesisRow{}).
		Where("id = ? AND lock_version = ? AND deleted_at IS NULL", id, expectedVersion).
		Updates(map[string]interface{}{
			"status":       string(status),
			"lock_version": gorm.Expr("lock_version + 1"),
			"updated_at":   gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return mapIntentPostgresError(result.Error)
	}
	if result.RowsAffected == 0 {
		return r.resolveIntentUpdateFailure(ctx, id)
	}
	return nil
}

func (r *SeoIntentHypothesisPostgres) LinkEvidence(
	ctx context.Context,
	link *seo_strategy.SeoIntentEvidenceLink,
) error {
	if err := seo_strategy.ValidateIntentEvidenceLink(link); err != nil {
		return fmt.Errorf("%w: %v", repo.ErrSeoIntentHypothesisInvalid, err)
	}
	row := seoIntentEvidenceLinkRow{
		ObjectiveID:        link.ObjectiveID,
		UserNeedID:         link.UserNeedID,
		IntentHypothesisID: link.IntentHypothesisID,
		QueryEvidenceID:    link.QueryEvidenceID,
		RelationType:       string(link.RelationType),
		RelevanceNote:      link.RelevanceNote,
		CreatedAt:          link.CreatedAt,
		UpdatedAt:          link.UpdatedAt,
	}
	if err := impl.GetDB(ctx, r.db).WithContext(ctx).Create(&row).Error; err != nil {
		return mapIntentEvidencePostgresError(err)
	}
	return nil
}

func (r *SeoIntentHypothesisPostgres) UnlinkEvidence(
	ctx context.Context,
	intentID uint64,
	queryEvidenceID uint64,
) error {
	result := impl.GetDB(ctx, r.db).WithContext(ctx).
		Where(
			"intent_hypothesis_id = ? AND query_evidence_id = ?",
			intentID,
			queryEvidenceID,
		).
		Delete(&seoIntentEvidenceLinkRow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repo.ErrSeoIntentEvidenceNotLinked
	}
	return nil
}

func (r *SeoIntentHypothesisPostgres) ListLinkedEvidence(
	ctx context.Context,
	intentID uint64,
) ([]seo_strategy.LinkedIntentQueryEvidence, error) {
	var rows []linkedIntentQueryEvidenceRow
	err := impl.GetDB(ctx, r.db).WithContext(ctx).Raw(`
		SELECT
			l.objective_id AS link_objective_id,
			l.user_need_id AS link_user_need_id,
			l.intent_hypothesis_id,
			l.query_evidence_id AS link_query_evidence_id,
			l.relation_type,
			l.relevance_note,
			l.created_at AS link_created_at,
			l.updated_at AS link_updated_at,
			e.id AS evidence_id,
			e.objective_id AS evidence_objective_id,
			e.user_need_id AS evidence_user_need_id,
			e.raw_query,
			e.normalized_query,
			e.query_hash,
			e.source_system,
			e.source_kind,
			e.observation_type,
			e.source_record_key,
			e.source_reference,
			e.locale,
			e.country_code,
			e.device_category,
			e.observed_from,
			e.observed_to,
			e.valid_until,
			e.impressions,
			e.clicks,
			e.ctr,
			e.avg_position,
			e.search_volume,
			e.confidence AS evidence_confidence,
			e.status AS evidence_status,
			e.metadata,
			e.created_at AS evidence_created_at,
			e.updated_at AS evidence_updated_at
		FROM seo_intent_evidence_link l
		JOIN seo_query_evidence e ON e.id = l.query_evidence_id
		WHERE l.intent_hypothesis_id = ?
		ORDER BY l.created_at ASC, e.id ASC
	`, intentID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]seo_strategy.LinkedIntentQueryEvidence, 0, len(rows))
	for i := range rows {
		row := &rows[i]
		evidence, err := queryEvidenceRowToDomain(&seoQueryEvidenceRow{
			ID:              row.EvidenceID,
			ObjectiveID:     row.EvidenceObjectiveID,
			UserNeedID:      row.EvidenceUserNeedID,
			RawQuery:        row.RawQuery,
			NormalizedQuery: row.NormalizedQuery,
			QueryHash:       row.QueryHash,
			SourceSystem:    row.SourceSystem,
			SourceKind:      row.SourceKind,
			ObservationType: row.ObservationType,
			SourceRecordKey: row.SourceRecordKey,
			SourceReference: row.SourceReference,
			Locale:          row.Locale,
			CountryCode:     row.CountryCode,
			DeviceCategory:  row.DeviceCategory,
			ObservedFrom:    row.ObservedFrom,
			ObservedTo:      row.ObservedTo,
			ValidUntil:      row.ValidUntil,
			Impressions:     row.Impressions,
			Clicks:          row.Clicks,
			CTR:             row.CTR,
			AvgPosition:     row.AvgPosition,
			SearchVolume:    row.SearchVolume,
			Confidence:      row.EvidenceConfidence,
			Status:          row.EvidenceStatus,
			Metadata:        row.Metadata,
			CreatedAt:       row.EvidenceCreatedAt,
			UpdatedAt:       row.EvidenceUpdatedAt,
		})
		if err != nil {
			return nil, err
		}
		items = append(items, seo_strategy.LinkedIntentQueryEvidence{
			Link: seo_strategy.SeoIntentEvidenceLink{
				ObjectiveID:        row.LinkObjectiveID,
				UserNeedID:         row.LinkUserNeedID,
				IntentHypothesisID: row.IntentHypothesisID,
				QueryEvidenceID:    row.LinkQueryEvidenceID,
				RelationType:       seo_strategy.IntentEvidenceRelation(row.RelationType),
				RelevanceNote:      row.RelevanceNote,
				CreatedAt:          row.LinkCreatedAt,
				UpdatedAt:          row.LinkUpdatedAt,
			},
			Evidence: *evidence,
		})
	}
	return items, nil
}

func (r *SeoIntentHypothesisPostgres) resolveIntentUpdateFailure(
	ctx context.Context,
	id uint64,
) error {
	var count int64
	if err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Model(&seoIntentHypothesisRow{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return repo.ErrSeoIntentHypothesisNotFound
	}
	return repo.ErrSeoIntentHypothesisConflict
}

func intentHypothesisToRow(intent *seo_strategy.SeoIntentHypothesis) seoIntentHypothesisRow {
	return seoIntentHypothesisRow{
		ID:              intent.ID,
		ObjectiveID:     intent.ObjectiveID,
		UserNeedID:      intent.UserNeedID,
		Code:            intent.Code,
		Name:            intent.Name,
		GoalStatement:   intent.GoalStatement,
		PrimaryQuestion: intent.PrimaryQuestion,
		ExpectedResult:  intent.ExpectedResult,
		SearchAction:    string(intent.SearchAction),
		SubjectType:     string(intent.SubjectType),
		ScopeType:       string(intent.ScopeType),
		Status:          string(intent.Status),
		Confidence:      string(intent.Confidence),
		OwnerCode:       intent.OwnerCode,
		ReviewDueAt:     intent.ReviewDueAt,
		LockVersion:     intent.LockVersion,
		CreatedAt:       intent.CreatedAt,
		UpdatedAt:       intent.UpdatedAt,
		DeletedAt:       intent.DeletedAt,
	}
}

func intentHypothesisRowToDomain(row *seoIntentHypothesisRow) *seo_strategy.SeoIntentHypothesis {
	if row == nil {
		return nil
	}
	return &seo_strategy.SeoIntentHypothesis{
		ID:              row.ID,
		ObjectiveID:     row.ObjectiveID,
		UserNeedID:      row.UserNeedID,
		Code:            row.Code,
		Name:            row.Name,
		GoalStatement:   row.GoalStatement,
		PrimaryQuestion: row.PrimaryQuestion,
		ExpectedResult:  row.ExpectedResult,
		SearchAction:    seo_strategy.IntentSearchAction(row.SearchAction),
		SubjectType:     seo_strategy.IntentSubjectType(row.SubjectType),
		ScopeType:       seo_strategy.IntentScopeType(row.ScopeType),
		Status:          seo_strategy.IntentHypothesisStatus(row.Status),
		Confidence:      seo_strategy.ObjectiveConfidence(row.Confidence),
		OwnerCode:       row.OwnerCode,
		ReviewDueAt:     row.ReviewDueAt,
		LockVersion:     row.LockVersion,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		DeletedAt:       row.DeletedAt,
	}
}

func mapIntentPostgresError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case "23505":
		return fmt.Errorf("%w: %s", repo.ErrSeoIntentHypothesisConflict, pgErr.ConstraintName)
	case "23503":
		return fmt.Errorf("%w: %s", repo.ErrSeoIntentEvidenceMismatch, pgErr.ConstraintName)
	case "23514", "22001", "22P02":
		return fmt.Errorf("%w: %s", repo.ErrSeoIntentHypothesisInvalid, pgErr.ConstraintName)
	default:
		return err
	}
}

func mapIntentEvidencePostgresError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case "23505":
		return fmt.Errorf("%w: %s", repo.ErrSeoIntentEvidenceAlreadyLinked, pgErr.ConstraintName)
	case "23503":
		return fmt.Errorf("%w: %s", repo.ErrSeoIntentEvidenceMismatch, pgErr.ConstraintName)
	case "23514", "22001", "22P02":
		return fmt.Errorf("%w: %s", repo.ErrSeoIntentHypothesisInvalid, pgErr.ConstraintName)
	default:
		return err
	}
}
