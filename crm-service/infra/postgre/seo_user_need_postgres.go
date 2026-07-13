package postgre

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"crm/infra/impl"
	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/repo"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SeoUserNeedPostgres struct {
	db *gorm.DB
}

var _ repo.SeoUserNeedRepo = (*SeoUserNeedPostgres)(nil)

func NewSeoUserNeedPostgres(
	db *gorm.DB,
) repo.SeoUserNeedRepo {
	return &SeoUserNeedPostgres{
		db: db,
	}
}

type seoUserNeedRow struct {
	ID                 uint64     `gorm:"column:id;primaryKey"`
	ObjectiveID        uint64     `gorm:"column:objective_id"`
	AudienceSegmentID  uint64     `gorm:"column:audience_segment_id"`
	Code               string     `gorm:"column:code"`
	Name               string     `gorm:"column:name"`
	DecisionStage      string     `gorm:"column:decision_stage"`
	Situation          string     `gorm:"column:situation"`
	TriggerStatement   string     `gorm:"column:trigger_statement"`
	JobToBeDone        string     `gorm:"column:job_to_be_done"`
	BarrierStatement   string     `gorm:"column:barrier_statement"`
	CurrentAlternative string     `gorm:"column:current_alternative"`
	DesiredOutcome     string     `gorm:"column:desired_outcome"`
	FailureImpact      string     `gorm:"column:failure_impact"`
	RiskLevel          string     `gorm:"column:risk_level"`
	Priority           string     `gorm:"column:priority"`
	Status             string     `gorm:"column:status"`
	Confidence         string     `gorm:"column:confidence"`
	OwnerCode          string     `gorm:"column:owner_code"`
	ReviewDueAt        *time.Time `gorm:"column:review_due_at"`
	LockVersion        uint32     `gorm:"column:lock_version"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"`
	DeletedAt          *time.Time `gorm:"column:deleted_at"`
}

func (seoUserNeedRow) TableName() string { return "seo_user_need" }

type seoUserNeedEvidenceRow struct {
	ObjectiveID   uint64    `gorm:"column:objective_id"`
	UserNeedID    uint64    `gorm:"column:user_need_id"`
	EvidenceID    uint64    `gorm:"column:evidence_id"`
	RelationType  string    `gorm:"column:relation_type"`
	RelevanceNote string    `gorm:"column:relevance_note"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (seoUserNeedEvidenceRow) TableName() string { return "seo_user_need_evidence" }

type linkedEvidenceQueryRow struct {
	LinkObjectiveID     uint64         `gorm:"column:link_objective_id"`
	LinkUserNeedID      uint64         `gorm:"column:link_user_need_id"`
	LinkEvidenceID      uint64         `gorm:"column:link_evidence_id"`
	RelationType        string         `gorm:"column:relation_type"`
	RelevanceNote       string         `gorm:"column:relevance_note"`
	LinkCreatedAt       time.Time      `gorm:"column:link_created_at"`
	EvidenceID          uint64         `gorm:"column:evidence_id"`
	EvidenceObjectiveID uint64         `gorm:"column:evidence_objective_id"`
	EvidenceType        string         `gorm:"column:evidence_type"`
	EvidenceStance      string         `gorm:"column:evidence_stance"`
	SourceSystem        string         `gorm:"column:source_system"`
	SourceReference     string         `gorm:"column:source_reference"`
	Finding             string         `gorm:"column:finding"`
	EvidenceConfidence  string         `gorm:"column:evidence_confidence"`
	ObservedAt          *time.Time     `gorm:"column:observed_at"`
	ValidUntil          *time.Time     `gorm:"column:valid_until"`
	EvidenceStatus      string         `gorm:"column:evidence_status"`
	Metadata            datatypes.JSON `gorm:"column:metadata"`
	EvidenceCreatedAt   time.Time      `gorm:"column:evidence_created_at"`
	EvidenceUpdatedAt   time.Time      `gorm:"column:evidence_updated_at"`
}

func (r *SeoUserNeedPostgres) Create(ctx context.Context, need *seo_strategy.SeoUserNeed) error {
	row := userNeedToRow(need)
	if err := impl.GetDB(ctx, r.db).WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	*need = *userNeedRowToDomain(&row)
	return nil
}

func (r *SeoUserNeedPostgres) GetByID(ctx context.Context, id uint64) (*seo_strategy.SeoUserNeed, error) {
	var row seoUserNeedRow
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoUserNeedNotFound
	}
	if err != nil {
		return nil, err
	}
	return userNeedRowToDomain(&row), nil
}

func (r *SeoUserNeedPostgres) GetByCode(ctx context.Context, objectiveID uint64, code string) (*seo_strategy.SeoUserNeed, error) {
	var row seoUserNeedRow
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Where("objective_id = ? AND code = ? AND deleted_at IS NULL", objectiveID, seo_strategy.NormalizeStrategyCode(code)).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoUserNeedNotFound
	}
	if err != nil {
		return nil, err
	}
	return userNeedRowToDomain(&row), nil
}

func (r *SeoUserNeedPostgres) List(ctx context.Context, filter repo.SeoUserNeedFilter) ([]seo_strategy.SeoUserNeed, int64, error) {
	size := filter.Size
	if size == 0 || size > 100 {
		size = 20
	}
	query := impl.GetDB(ctx, r.db).WithContext(ctx).Model(&seoUserNeedRow{}).Where("deleted_at IS NULL")
	if filter.ObjectiveID > 0 {
		query = query.Where("objective_id = ?", filter.ObjectiveID)
	}
	if filter.AudienceID > 0 {
		query = query.Where("audience_segment_id = ?", filter.AudienceID)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(`code ILIKE ? OR name ILIKE ? OR situation ILIKE ? OR job_to_be_done ILIKE ? OR desired_outcome ILIKE ?`, like, like, like, like, like)
	}
	if status := strings.TrimSpace(filter.Status); status != "" {
		query = query.Where("status = ?", status)
	}
	if confidence := strings.TrimSpace(filter.Confidence); confidence != "" {
		query = query.Where("confidence = ?", confidence)
	}
	if priority := strings.TrimSpace(filter.Priority); priority != "" {
		query = query.Where("priority = ?", priority)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []seoUserNeedRow
	if err := query.Order("updated_at DESC, id DESC").Offset(int(filter.Page * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	items := make([]seo_strategy.SeoUserNeed, 0, len(rows))
	for i := range rows {
		items = append(items, *userNeedRowToDomain(&rows[i]))
	}
	return items, total, nil
}

func (r *SeoUserNeedPostgres) UpdateDefinition(ctx context.Context, need *seo_strategy.SeoUserNeed, expectedVersion uint32) error {
	result := impl.GetDB(ctx, r.db).WithContext(ctx).Model(&seoUserNeedRow{}).
		Where("id = ? AND lock_version = ? AND deleted_at IS NULL", need.ID, expectedVersion).
		Updates(map[string]interface{}{
			"name": need.Name, "decision_stage": string(need.DecisionStage),
			"situation": need.Situation, "trigger_statement": need.TriggerStatement,
			"job_to_be_done": need.JobToBeDone, "barrier_statement": need.BarrierStatement,
			"current_alternative": need.CurrentAlternative, "desired_outcome": need.DesiredOutcome,
			"failure_impact": need.FailureImpact, "risk_level": string(need.RiskLevel),
			"priority": string(need.Priority), "confidence": string(need.Confidence),
			"owner_code": need.OwnerCode, "review_due_at": need.ReviewDueAt,
			"lock_version": gorm.Expr("lock_version + 1"), "updated_at": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return r.resolveUpdateFailure(ctx, need.ID)
	}
	return nil
}

func (r *SeoUserNeedPostgres) UpdateStatus(ctx context.Context, id uint64, status seo_strategy.UserNeedStatus, expectedVersion uint32) error {
	result := impl.GetDB(ctx, r.db).WithContext(ctx).Model(&seoUserNeedRow{}).
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

func (r *SeoUserNeedPostgres) LinkEvidence(ctx context.Context, link *seo_strategy.SeoUserNeedEvidenceLink) error {
	row := seoUserNeedEvidenceRow{
		ObjectiveID: link.ObjectiveID, UserNeedID: link.UserNeedID, EvidenceID: link.EvidenceID,
		RelationType: string(link.RelationType), RelevanceNote: link.RelevanceNote, CreatedAt: link.CreatedAt,
	}
	err := impl.GetDB(ctx, r.db).WithContext(ctx).Create(&row).Error
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return repo.ErrSeoUserNeedEvidenceAlreadyLinked
	}
	return err
}

func (r *SeoUserNeedPostgres) UnlinkEvidence(ctx context.Context, userNeedID, evidenceID uint64) error {
	return impl.GetDB(ctx, r.db).WithContext(ctx).
		Where("user_need_id = ? AND evidence_id = ?", userNeedID, evidenceID).
		Delete(&seoUserNeedEvidenceRow{}).Error
}

func (r *SeoUserNeedPostgres) ListLinkedEvidence(ctx context.Context, userNeedID uint64) ([]seo_strategy.LinkedUserNeedEvidence, error) {
	var rows []linkedEvidenceQueryRow
	err := impl.GetDB(ctx, r.db).WithContext(ctx).Raw(`
		SELECT
			l.objective_id AS link_objective_id,
			l.user_need_id AS link_user_need_id,
			l.evidence_id AS link_evidence_id,
			l.relation_type,
			l.relevance_note,
			l.created_at AS link_created_at,
			e.id AS evidence_id,
			e.objective_id AS evidence_objective_id,
			e.evidence_type,
			e.stance AS evidence_stance,
			e.source_system,
			e.source_reference,
			e.finding,
			e.confidence AS evidence_confidence,
			e.observed_at,
			e.valid_until,
			e.status AS evidence_status,
			e.metadata,
			e.created_at AS evidence_created_at,
			e.updated_at AS evidence_updated_at
		FROM seo_user_need_evidence l
		INNER JOIN seo_objective_evidence e
			ON e.id = l.evidence_id AND e.objective_id = l.objective_id
		WHERE l.user_need_id = ?
		ORDER BY e.created_at DESC, e.id DESC
	`, userNeedID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]seo_strategy.LinkedUserNeedEvidence, 0, len(rows))
	for _, row := range rows {
		metadata := json.RawMessage(row.Metadata)
		if len(metadata) == 0 {
			metadata = json.RawMessage(`{}`)
		}
		result = append(result, seo_strategy.LinkedUserNeedEvidence{
			Link: seo_strategy.SeoUserNeedEvidenceLink{
				ObjectiveID: row.LinkObjectiveID, UserNeedID: row.LinkUserNeedID, EvidenceID: row.LinkEvidenceID,
				RelationType: seo_strategy.UserNeedEvidenceRelation(row.RelationType), RelevanceNote: row.RelevanceNote, CreatedAt: row.LinkCreatedAt,
			},
			Evidence: seo_strategy.SeoObjectiveEvidence{
				ID: row.EvidenceID, ObjectiveID: row.EvidenceObjectiveID,
				Type: seo_strategy.EvidenceType(row.EvidenceType), Stance: seo_strategy.EvidenceStance(row.EvidenceStance),
				SourceSystem: row.SourceSystem, SourceReference: row.SourceReference, Finding: row.Finding,
				Confidence: seo_strategy.ObjectiveConfidence(row.EvidenceConfidence), ObservedAt: row.ObservedAt, ValidUntil: row.ValidUntil,
				Status: seo_strategy.EvidenceStatus(row.EvidenceStatus), Metadata: metadata,
				CreatedAt: row.EvidenceCreatedAt, UpdatedAt: row.EvidenceUpdatedAt,
			},
		})
	}
	return result, nil
}

func (r *SeoUserNeedPostgres) resolveUpdateFailure(ctx context.Context, id uint64) error {
	var count int64
	if err := impl.GetDB(ctx, r.db).WithContext(ctx).Model(&seoUserNeedRow{}).
		Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return repo.ErrSeoUserNeedNotFound
	}
	return repo.ErrSeoUserNeedConflict
}

func userNeedToRow(item *seo_strategy.SeoUserNeed) seoUserNeedRow {
	return seoUserNeedRow{
		ID: item.ID, ObjectiveID: item.ObjectiveID, AudienceSegmentID: item.AudienceSegmentID,
		Code: item.Code, Name: item.Name, DecisionStage: string(item.DecisionStage),
		Situation: item.Situation, TriggerStatement: item.TriggerStatement, JobToBeDone: item.JobToBeDone,
		BarrierStatement: item.BarrierStatement, CurrentAlternative: item.CurrentAlternative,
		DesiredOutcome: item.DesiredOutcome, FailureImpact: item.FailureImpact,
		RiskLevel: string(item.RiskLevel), Priority: string(item.Priority), Status: string(item.Status), Confidence: string(item.Confidence),
		OwnerCode: item.OwnerCode, ReviewDueAt: item.ReviewDueAt, LockVersion: item.LockVersion,
		CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt, DeletedAt: item.DeletedAt,
	}
}

func userNeedRowToDomain(row *seoUserNeedRow) *seo_strategy.SeoUserNeed {
	if row == nil {
		return nil
	}
	return &seo_strategy.SeoUserNeed{
		ID: row.ID, ObjectiveID: row.ObjectiveID, AudienceSegmentID: row.AudienceSegmentID,
		Code: row.Code, Name: row.Name, DecisionStage: seo_strategy.UserNeedDecisionStage(row.DecisionStage),
		Situation: row.Situation, TriggerStatement: row.TriggerStatement, JobToBeDone: row.JobToBeDone,
		BarrierStatement: row.BarrierStatement, CurrentAlternative: row.CurrentAlternative,
		DesiredOutcome: row.DesiredOutcome, FailureImpact: row.FailureImpact,
		RiskLevel: seo_strategy.UserNeedRiskLevel(row.RiskLevel), Priority: seo_strategy.UserNeedPriority(row.Priority),
		Status: seo_strategy.UserNeedStatus(row.Status), Confidence: seo_strategy.ObjectiveConfidence(row.Confidence),
		OwnerCode: row.OwnerCode, ReviewDueAt: row.ReviewDueAt, LockVersion: row.LockVersion,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt, DeletedAt: row.DeletedAt,
	}
}
