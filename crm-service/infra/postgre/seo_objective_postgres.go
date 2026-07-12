package postgre

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/repo"

	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeoObjectivePostgres struct {
	db *gorm.DB
}

func NewSeoObjectivePostgres(
	db *gorm.DB,
) repo.SeoObjectiveRepo {
	return &SeoObjectivePostgres{db: db}
}

type seoObjectiveRow struct {
	ID   uint64 `gorm:"column:id;primaryKey"`
	Code string `gorm:"column:code"`
	Name string `gorm:"column:name"`

	AudienceSummary  string `gorm:"column:audience_summary"`
	Situation        string `gorm:"column:situation"`
	ProblemStatement string `gorm:"column:problem_statement"`

	UserOutcome     string `gorm:"column:user_outcome"`
	ProductOutcome  string `gorm:"column:product_outcome"`
	BusinessOutcome string `gorm:"column:business_outcome"`

	NonGoals    pq.StringArray `gorm:"column:non_goals;type:text[]"`
	Constraints pq.StringArray `gorm:"column:constraints;type:text[]"`

	Status     string `gorm:"column:status"`
	Confidence string `gorm:"column:confidence"`

	OwnerCode   string     `gorm:"column:owner_code"`
	ReviewDueAt *time.Time `gorm:"column:review_due_at"`

	LockVersion uint32 `gorm:"column:lock_version"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (seoObjectiveRow) TableName() string {
	return "seo_objective"
}

type seoObjectiveEvidenceRow struct {
	ID          uint64 `gorm:"column:id;primaryKey"`
	ObjectiveID uint64 `gorm:"column:objective_id"`

	EvidenceType string `gorm:"column:evidence_type"`
	Stance       string `gorm:"column:stance"`

	SourceSystem    string `gorm:"column:source_system"`
	SourceReference string `gorm:"column:source_reference"`

	Finding    string `gorm:"column:finding"`
	Confidence string `gorm:"column:confidence"`

	ObservedAt *time.Time `gorm:"column:observed_at"`
	ValidUntil *time.Time `gorm:"column:valid_until"`

	Status   string         `gorm:"column:status"`
	Metadata datatypes.JSON `gorm:"column:metadata"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (seoObjectiveEvidenceRow) TableName() string {
	return "seo_objective_evidence"
}

type seoObjectiveMetricRow struct {
	ID          uint64 `gorm:"column:id;primaryKey"`
	ObjectiveID uint64 `gorm:"column:objective_id"`

	MetricCode string `gorm:"column:metric_code"`
	MetricName string `gorm:"column:metric_name"`

	MetricLayer string `gorm:"column:metric_layer"`
	Description string `gorm:"column:description"`

	SourceSystem      string `gorm:"column:source_system"`
	CalculationMethod string `gorm:"column:calculation_method"`

	Unit            string `gorm:"column:unit"`
	TargetDirection string `gorm:"column:target_direction"`

	BaselineValue *float64 `gorm:"column:baseline_value"`
	TargetValue   *float64 `gorm:"column:target_value"`

	IsPrimary bool `gorm:"column:is_primary"`

	InstrumentationStatus string `gorm:"column:instrumentation_status"`
	Status                string `gorm:"column:status"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (seoObjectiveMetricRow) TableName() string {
	return "seo_objective_metric"
}

func (r *SeoObjectivePostgres) CreateObjective(
	ctx context.Context,
	objective *seo_strategy.SeoObjective,
) error {
	row := objectiveToRow(objective)

	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}

	*objective = *objectiveRowToDomain(&row)
	return nil
}

func (r *SeoObjectivePostgres) GetObjectiveByID(
	ctx context.Context,
	id uint64,
) (*seo_strategy.SeoObjective, error) {
	var row seoObjectiveRow

	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoObjectiveNotFound
	}

	if err != nil {
		return nil, err
	}

	return objectiveRowToDomain(&row), nil
}

func (r *SeoObjectivePostgres) GetObjectiveByCode(
	ctx context.Context,
	code string,
) (*seo_strategy.SeoObjective, error) {
	var row seoObjectiveRow

	err := r.db.WithContext(ctx).
		Where(
			"code = ? AND deleted_at IS NULL",
			seo_strategy.NormalizeObjectiveCode(code),
		).
		First(&row).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoObjectiveNotFound
	}

	if err != nil {
		return nil, err
	}

	return objectiveRowToDomain(&row), nil
}

func (r *SeoObjectivePostgres) ListObjectives(
	ctx context.Context,
	filter repo.SeoObjectiveFilter,
) ([]seo_strategy.SeoObjective, int64, error) {
	size := filter.Size
	if size == 0 || size > 100 {
		size = 20
	}

	query := r.db.WithContext(ctx).
		Model(&seoObjectiveRow{}).
		Where("deleted_at IS NULL")

	keyword := strings.TrimSpace(filter.Keyword)
	if keyword != "" {
		query = query.Where(
			`code ILIKE ? OR
			 name ILIKE ? OR
			 problem_statement ILIKE ?`,
			"%"+keyword+"%",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
	}

	if status := strings.TrimSpace(filter.Status); status != "" {
		query = query.Where("status = ?", status)
	}

	if confidence := strings.TrimSpace(filter.Confidence); confidence != "" {
		query = query.Where("confidence = ?", confidence)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []seoObjectiveRow
	err := query.
		Order("updated_at DESC, id DESC").
		Offset(int(filter.Page * size)).
		Limit(int(size)).
		Find(&rows).
		Error

	if err != nil {
		return nil, 0, err
	}

	items := make([]seo_strategy.SeoObjective, 0, len(rows))
	for i := range rows {
		items = append(items, *objectiveRowToDomain(&rows[i]))
	}

	return items, total, nil
}

func (r *SeoObjectivePostgres) UpdateObjectiveDefinition(
	ctx context.Context,
	objective *seo_strategy.SeoObjective,
	expectedVersion uint32,
) error {
	updates := map[string]interface{}{
		"name":              objective.Name,
		"audience_summary":  objective.AudienceSummary,
		"situation":         objective.Situation,
		"problem_statement": objective.ProblemStatement,
		"user_outcome":      objective.UserOutcome,
		"product_outcome":   objective.ProductOutcome,
		"business_outcome":  objective.BusinessOutcome,
		"non_goals":         pq.StringArray(objective.NonGoals),
		"constraints":       pq.StringArray(objective.Constraints),
		"confidence":        string(objective.Confidence),
		"owner_code":        objective.OwnerCode,
		"review_due_at":     objective.ReviewDueAt,
		"lock_version":      gorm.Expr("lock_version + 1"),
		"updated_at":        gorm.Expr("NOW()"),
	}

	result := r.db.WithContext(ctx).
		Model(&seoObjectiveRow{}).
		Where(
			"id = ? AND lock_version = ? AND deleted_at IS NULL",
			objective.ID,
			expectedVersion,
		).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return r.resolveUpdateFailure(ctx, objective.ID)
	}

	return nil
}

func (r *SeoObjectivePostgres) UpdateObjectiveStatus(
	ctx context.Context,
	id uint64,
	status seo_strategy.ObjectiveStatus,
	expectedVersion uint32,
) error {
	result := r.db.WithContext(ctx).
		Model(&seoObjectiveRow{}).
		Where(
			"id = ? AND lock_version = ? AND deleted_at IS NULL",
			id,
			expectedVersion,
		).
		Updates(map[string]interface{}{
			"status":       string(status),
			"lock_version": gorm.Expr("lock_version + 1"),
			"updated_at":   gorm.Expr("NOW()"),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return r.resolveUpdateFailure(ctx, id)
	}

	return nil
}

func (r *SeoObjectivePostgres) CreateEvidence(
	ctx context.Context,
	evidence *seo_strategy.SeoObjectiveEvidence,
) error {
	row := evidenceToRow(evidence)

	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}

	*evidence = *evidenceRowToDomain(&row)
	return nil
}

func (r *SeoObjectivePostgres) ListEvidence(
	ctx context.Context,
	objectiveID uint64,
) ([]seo_strategy.SeoObjectiveEvidence, error) {
	var rows []seoObjectiveEvidenceRow

	err := r.db.WithContext(ctx).
		Where("objective_id = ?", objectiveID).
		Order("created_at DESC, id DESC").
		Find(&rows).
		Error

	if err != nil {
		return nil, err
	}

	items := make([]seo_strategy.SeoObjectiveEvidence, 0, len(rows))
	for i := range rows {
		items = append(items, *evidenceRowToDomain(&rows[i]))
	}

	return items, nil
}

func (r *SeoObjectivePostgres) UpsertMetric(
	ctx context.Context,
	metric *seo_strategy.SeoObjectiveMetric,
) error {
	return r.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			// Một objective chỉ được có một primary metric.
			if metric.IsPrimary &&
				metric.Status == seo_strategy.MetricStatusActive {
				if err := tx.
					Model(&seoObjectiveMetricRow{}).
					Where(
						"objective_id = ? AND metric_code <> ?",
						metric.ObjectiveID,
						metric.Code,
					).
					Updates(map[string]interface{}{
						"is_primary": false,
						"updated_at": gorm.Expr("NOW()"),
					}).
					Error; err != nil {
					return err
				}
			}

			row := metricToRow(metric)

			err := tx.
				Clauses(clause.OnConflict{
					Columns: []clause.Column{
						{Name: "objective_id"},
						{Name: "metric_code"},
					},
					DoUpdates: clause.AssignmentColumns([]string{
						"metric_name",
						"metric_layer",
						"description",
						"source_system",
						"calculation_method",
						"unit",
						"target_direction",
						"baseline_value",
						"target_value",
						"is_primary",
						"instrumentation_status",
						"status",
						"updated_at",
					}),
				}).
				Create(&row).
				Error

			if err != nil {
				return err
			}

			var persisted seoObjectiveMetricRow
			if err := tx.
				Where(
					"objective_id = ? AND metric_code = ?",
					metric.ObjectiveID,
					metric.Code,
				).
				First(&persisted).
				Error; err != nil {
				return err
			}

			*metric = *metricRowToDomain(&persisted)
			return nil
		},
	)
}

func (r *SeoObjectivePostgres) ListMetrics(
	ctx context.Context,
	objectiveID uint64,
) ([]seo_strategy.SeoObjectiveMetric, error) {
	var rows []seoObjectiveMetricRow

	err := r.db.WithContext(ctx).
		Where("objective_id = ?", objectiveID).
		Order("is_primary DESC, created_at ASC").
		Find(&rows).
		Error

	if err != nil {
		return nil, err
	}

	items := make([]seo_strategy.SeoObjectiveMetric, 0, len(rows))
	for i := range rows {
		items = append(items, *metricRowToDomain(&rows[i]))
	}

	return items, nil
}

func (r *SeoObjectivePostgres) resolveUpdateFailure(
	ctx context.Context,
	id uint64,
) error {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&seoObjectiveRow{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).
		Error; err != nil {
		return err
	}

	if count == 0 {
		return repo.ErrSeoObjectiveNotFound
	}

	return repo.ErrSeoObjectiveConflict
}

func objectiveToRow(
	item *seo_strategy.SeoObjective,
) seoObjectiveRow {
	return seoObjectiveRow{
		ID:               item.ID,
		Code:             item.Code,
		Name:             item.Name,
		AudienceSummary:  item.AudienceSummary,
		Situation:        item.Situation,
		ProblemStatement: item.ProblemStatement,
		UserOutcome:      item.UserOutcome,
		ProductOutcome:   item.ProductOutcome,
		BusinessOutcome:  item.BusinessOutcome,
		NonGoals:         pq.StringArray(item.NonGoals),
		Constraints:      pq.StringArray(item.Constraints),
		Status:           string(item.Status),
		Confidence:       string(item.Confidence),
		OwnerCode:        item.OwnerCode,
		ReviewDueAt:      item.ReviewDueAt,
		LockVersion:      item.LockVersion,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
		DeletedAt:        item.DeletedAt,
	}
}

func objectiveRowToDomain(
	row *seoObjectiveRow,
) *seo_strategy.SeoObjective {
	if row == nil {
		return nil
	}

	return &seo_strategy.SeoObjective{
		ID:               row.ID,
		Code:             row.Code,
		Name:             row.Name,
		AudienceSummary:  row.AudienceSummary,
		Situation:        row.Situation,
		ProblemStatement: row.ProblemStatement,
		UserOutcome:      row.UserOutcome,
		ProductOutcome:   row.ProductOutcome,
		BusinessOutcome:  row.BusinessOutcome,
		NonGoals:         []string(row.NonGoals),
		Constraints:      []string(row.Constraints),
		Status:           seo_strategy.ObjectiveStatus(row.Status),
		Confidence:       seo_strategy.ObjectiveConfidence(row.Confidence),
		OwnerCode:        row.OwnerCode,
		ReviewDueAt:      row.ReviewDueAt,
		LockVersion:      row.LockVersion,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
		DeletedAt:        row.DeletedAt,
	}
}

func evidenceToRow(
	item *seo_strategy.SeoObjectiveEvidence,
) seoObjectiveEvidenceRow {
	metadata := item.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}

	return seoObjectiveEvidenceRow{
		ID:              item.ID,
		ObjectiveID:     item.ObjectiveID,
		EvidenceType:    string(item.Type),
		Stance:          string(item.Stance),
		SourceSystem:    item.SourceSystem,
		SourceReference: item.SourceReference,
		Finding:         item.Finding,
		Confidence:      string(item.Confidence),
		ObservedAt:      item.ObservedAt,
		ValidUntil:      item.ValidUntil,
		Status:          string(item.Status),
		Metadata:        datatypes.JSON(metadata),
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}
}

func evidenceRowToDomain(
	row *seoObjectiveEvidenceRow,
) *seo_strategy.SeoObjectiveEvidence {
	if row == nil {
		return nil
	}

	return &seo_strategy.SeoObjectiveEvidence{
		ID:              row.ID,
		ObjectiveID:     row.ObjectiveID,
		Type:            seo_strategy.EvidenceType(row.EvidenceType),
		Stance:          seo_strategy.EvidenceStance(row.Stance),
		SourceSystem:    row.SourceSystem,
		SourceReference: row.SourceReference,
		Finding:         row.Finding,
		Confidence:      seo_strategy.ObjectiveConfidence(row.Confidence),
		ObservedAt:      row.ObservedAt,
		ValidUntil:      row.ValidUntil,
		Status:          seo_strategy.EvidenceStatus(row.Status),
		Metadata:        json.RawMessage(row.Metadata),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

func metricToRow(
	item *seo_strategy.SeoObjectiveMetric,
) seoObjectiveMetricRow {
	return seoObjectiveMetricRow{
		ID:                    item.ID,
		ObjectiveID:           item.ObjectiveID,
		MetricCode:            item.Code,
		MetricName:            item.Name,
		MetricLayer:           string(item.Layer),
		Description:           item.Description,
		SourceSystem:          item.SourceSystem,
		CalculationMethod:     item.CalculationMethod,
		Unit:                  item.Unit,
		TargetDirection:       string(item.TargetDirection),
		BaselineValue:         item.BaselineValue,
		TargetValue:           item.TargetValue,
		IsPrimary:             item.IsPrimary,
		InstrumentationStatus: string(item.InstrumentationStatus),
		Status:                string(item.Status),
		CreatedAt:             item.CreatedAt,
		UpdatedAt:             item.UpdatedAt,
	}
}

func metricRowToDomain(
	row *seoObjectiveMetricRow,
) *seo_strategy.SeoObjectiveMetric {
	if row == nil {
		return nil
	}

	return &seo_strategy.SeoObjectiveMetric{
		ID:                    row.ID,
		ObjectiveID:           row.ObjectiveID,
		Code:                  row.MetricCode,
		Name:                  row.MetricName,
		Layer:                 seo_strategy.MetricLayer(row.MetricLayer),
		Description:           row.Description,
		SourceSystem:          row.SourceSystem,
		CalculationMethod:     row.CalculationMethod,
		Unit:                  row.Unit,
		TargetDirection:       seo_strategy.TargetDirection(row.TargetDirection),
		BaselineValue:         row.BaselineValue,
		TargetValue:           row.TargetValue,
		IsPrimary:             row.IsPrimary,
		InstrumentationStatus: seo_strategy.InstrumentationStatus(row.InstrumentationStatus),
		Status:                seo_strategy.MetricStatus(row.Status),
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}
}
