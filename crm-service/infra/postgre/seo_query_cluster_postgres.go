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

type SeoQueryClusterPostgres struct {
	db *gorm.DB
}

var _ repo.SeoQueryClusterRepo = (*SeoQueryClusterPostgres)(nil)

func NewSeoQueryClusterPostgres(db *gorm.DB) repo.SeoQueryClusterRepo {
	return &SeoQueryClusterPostgres{db: db}
}

type seoQueryClusterRow struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	ObjectiveID        uint64 `gorm:"column:objective_id"`
	UserNeedID         uint64 `gorm:"column:user_need_id"`
	IntentHypothesisID uint64 `gorm:"column:intent_hypothesis_id"`

	Code string `gorm:"column:code"`
	Name string `gorm:"column:name"`

	SearchPromise string `gorm:"column:search_promise"`
	InclusionRule string `gorm:"column:inclusion_rule"`
	ExclusionRule string `gorm:"column:exclusion_rule"`

	Status     string `gorm:"column:status"`
	Confidence string `gorm:"column:confidence"`

	OwnerCode   string     `gorm:"column:owner_code"`
	ReviewDueAt *time.Time `gorm:"column:review_due_at"`
	LockVersion uint32     `gorm:"column:lock_version"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (seoQueryClusterRow) TableName() string {
	return "seo_query_cluster"
}

type seoQueryClusterMemberRow struct {
	ObjectiveID uint64 `gorm:"column:objective_id"`
	UserNeedID  uint64 `gorm:"column:user_need_id"`

	ClusterID       uint64 `gorm:"column:cluster_id"`
	QueryEvidenceID uint64 `gorm:"column:query_evidence_id"`

	MembershipType string `gorm:"column:membership_type"`
	Rationale      string `gorm:"column:rationale"`
	Status         string `gorm:"column:status"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (seoQueryClusterMemberRow) TableName() string {
	return "seo_query_cluster_member"
}

type linkedQueryClusterMemberRow struct {
	MemberObjectiveID     uint64         `gorm:"column:member_objective_id"`
	MemberUserNeedID      uint64         `gorm:"column:member_user_need_id"`
	ClusterID             uint64         `gorm:"column:cluster_id"`
	MemberQueryEvidenceID uint64         `gorm:"column:member_query_evidence_id"`
	MembershipType        string         `gorm:"column:membership_type"`
	Rationale             string         `gorm:"column:rationale"`
	MemberStatus          string         `gorm:"column:member_status"`
	MemberCreatedAt       time.Time      `gorm:"column:member_created_at"`
	MemberUpdatedAt       time.Time      `gorm:"column:member_updated_at"`
	EvidenceID            uint64         `gorm:"column:evidence_id"`
	EvidenceObjectiveID   uint64         `gorm:"column:evidence_objective_id"`
	EvidenceUserNeedID    uint64         `gorm:"column:evidence_user_need_id"`
	RawQuery              string         `gorm:"column:raw_query"`
	NormalizedQuery       string         `gorm:"column:normalized_query"`
	QueryHash             string         `gorm:"column:query_hash"`
	SourceSystem          string         `gorm:"column:source_system"`
	SourceKind            string         `gorm:"column:source_kind"`
	ObservationType       string         `gorm:"column:observation_type"`
	SourceRecordKey       string         `gorm:"column:source_record_key"`
	SourceReference       string         `gorm:"column:source_reference"`
	Locale                string         `gorm:"column:locale"`
	CountryCode           string         `gorm:"column:country_code"`
	DeviceCategory        string         `gorm:"column:device_category"`
	ObservedFrom          *time.Time     `gorm:"column:observed_from"`
	ObservedTo            *time.Time     `gorm:"column:observed_to"`
	ValidUntil            *time.Time     `gorm:"column:valid_until"`
	Impressions           *int64         `gorm:"column:impressions"`
	Clicks                *int64         `gorm:"column:clicks"`
	CTR                   *float64       `gorm:"column:ctr"`
	AvgPosition           *float64       `gorm:"column:avg_position"`
	SearchVolume          *int64         `gorm:"column:search_volume"`
	EvidenceConfidence    string         `gorm:"column:evidence_confidence"`
	EvidenceStatus        string         `gorm:"column:evidence_status"`
	Metadata              datatypes.JSON `gorm:"column:metadata"`
	EvidenceCreatedAt     time.Time      `gorm:"column:evidence_created_at"`
	EvidenceUpdatedAt     time.Time      `gorm:"column:evidence_updated_at"`
}

func (r *SeoQueryClusterPostgres) Create(
	ctx context.Context,
	cluster *seo_strategy.SeoQueryCluster,
) error {
	if cluster == nil {
		return repo.ErrSeoQueryClusterInvalid
	}
	if err := seo_strategy.ValidateQueryClusterDefinition(cluster); err != nil {
		return fmt.Errorf("%w: %v", repo.ErrSeoQueryClusterInvalid, err)
	}
	row := queryClusterToRow(cluster)
	if err := impl.GetDB(ctx, r.db).WithContext(ctx).Create(&row).Error; err != nil {
		return mapQueryClusterPostgresError(err)
	}
	*cluster = *queryClusterRowToDomain(&row)
	return nil
}

func (r *SeoQueryClusterPostgres) GetByID(
	ctx context.Context,
	id uint64,
) (*seo_strategy.SeoQueryCluster, error) {
	var row seoQueryClusterRow
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoQueryClusterNotFound
	}
	if err != nil {
		return nil, err
	}
	return queryClusterRowToDomain(&row), nil
}

func (r *SeoQueryClusterPostgres) GetByCode(
	ctx context.Context,
	userNeedID uint64,
	code string,
) (*seo_strategy.SeoQueryCluster, error) {
	var row seoQueryClusterRow
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Where(
			"user_need_id = ? AND code = ? AND deleted_at IS NULL",
			userNeedID,
			seo_strategy.NormalizeStrategyCode(code),
		).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrSeoQueryClusterNotFound
	}
	if err != nil {
		return nil, err
	}
	return queryClusterRowToDomain(&row), nil
}

func (r *SeoQueryClusterPostgres) List(
	ctx context.Context,
	filter repo.SeoQueryClusterFilter,
) ([]seo_strategy.SeoQueryCluster, int64, error) {
	size := filter.Size
	if size == 0 || size > 100 {
		size = 20
	}
	query := impl.GetDB(ctx, r.db).WithContext(ctx).
		Model(&seoQueryClusterRow{}).
		Where("deleted_at IS NULL")

	if filter.ObjectiveID > 0 {
		query = query.Where("objective_id = ?", filter.ObjectiveID)
	}
	if filter.UserNeedID > 0 {
		query = query.Where("user_need_id = ?", filter.UserNeedID)
	}
	if filter.IntentHypothesisID > 0 {
		query = query.Where("intent_hypothesis_id = ?", filter.IntentHypothesisID)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			`code ILIKE ? OR name ILIKE ? OR search_promise ILIKE ? OR inclusion_rule ILIKE ? OR exclusion_rule ILIKE ?`,
			like, like, like, like, like,
		)
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
	var rows []seoQueryClusterRow
	if err := query.
		Order("updated_at DESC, id DESC").
		Offset(int(filter.Page * size)).
		Limit(int(size)).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]seo_strategy.SeoQueryCluster, 0, len(rows))
	for i := range rows {
		items = append(items, *queryClusterRowToDomain(&rows[i]))
	}
	return items, total, nil
}

func (r *SeoQueryClusterPostgres) UpdateDefinition(
	ctx context.Context,
	cluster *seo_strategy.SeoQueryCluster,
	expectedVersion uint32,
) error {
	if cluster == nil || expectedVersion == 0 {
		return repo.ErrSeoQueryClusterInvalid
	}
	if err := seo_strategy.ValidateQueryClusterDefinition(cluster); err != nil {
		return fmt.Errorf("%w: %v", repo.ErrSeoQueryClusterInvalid, err)
	}
	result := impl.GetDB(ctx, r.db).WithContext(ctx).
		Model(&seoQueryClusterRow{}).
		Where("id = ? AND lock_version = ? AND deleted_at IS NULL", cluster.ID, expectedVersion).
		Updates(map[string]interface{}{
			"name":           cluster.Name,
			"search_promise": cluster.SearchPromise,
			"inclusion_rule": cluster.InclusionRule,
			"exclusion_rule": cluster.ExclusionRule,
			"confidence":     string(cluster.Confidence),
			"owner_code":     cluster.OwnerCode,
			"review_due_at":  cluster.ReviewDueAt,
			"lock_version":   gorm.Expr("lock_version + 1"),
			"updated_at":     gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return mapQueryClusterPostgresError(result.Error)
	}
	if result.RowsAffected == 0 {
		return r.resolveQueryClusterUpdateFailure(ctx, impl.GetDB(ctx, r.db), cluster.ID)
	}
	return nil
}

func (r *SeoQueryClusterPostgres) UpdateStatus(
	ctx context.Context,
	id uint64,
	status seo_strategy.QueryClusterStatus,
	expectedVersion uint32,
) error {
	result := impl.GetDB(ctx, r.db).WithContext(ctx).
		Model(&seoQueryClusterRow{}).
		Where("id = ? AND lock_version = ? AND deleted_at IS NULL", id, expectedVersion).
		Updates(map[string]interface{}{
			"status":       string(status),
			"lock_version": gorm.Expr("lock_version + 1"),
			"updated_at":   gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return mapQueryClusterPostgresError(result.Error)
	}
	if result.RowsAffected == 0 {
		return r.resolveQueryClusterUpdateFailure(ctx, impl.GetDB(ctx, r.db), id)
	}
	return nil
}

func (r *SeoQueryClusterPostgres) UpsertMember(
	ctx context.Context,
	member *seo_strategy.SeoQueryClusterMember,
	expectedVersion uint32,
) error {
	if expectedVersion == 0 {
		return repo.ErrSeoQueryClusterInvalid
	}
	if err := seo_strategy.ValidateQueryClusterMember(member); err != nil {
		return fmt.Errorf("%w: %v", repo.ErrSeoQueryClusterInvalid, err)
	}

	return impl.GetDB(ctx, r.db).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&seoQueryClusterRow{}).
			Where("id = ? AND lock_version = ? AND deleted_at IS NULL", member.ClusterID, expectedVersion).
			Updates(map[string]interface{}{
				"lock_version": gorm.Expr("lock_version + 1"),
				"updated_at":   gorm.Expr("NOW()"),
			})
		if result.Error != nil {
			return mapQueryClusterPostgresError(result.Error)
		}
		if result.RowsAffected == 0 {
			return r.resolveQueryClusterUpdateFailure(ctx, tx, member.ClusterID)
		}

		err := tx.Exec(`
			INSERT INTO seo_query_cluster_member (
				objective_id,
				user_need_id,
				cluster_id,
				query_evidence_id,
				membership_type,
				rationale,
				status,
				created_at,
				updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (cluster_id, query_evidence_id)
			DO UPDATE SET
				objective_id = EXCLUDED.objective_id,
				user_need_id = EXCLUDED.user_need_id,
				membership_type = EXCLUDED.membership_type,
				rationale = EXCLUDED.rationale,
				status = 'active',
				updated_at = EXCLUDED.updated_at
		`,
			member.ObjectiveID,
			member.UserNeedID,
			member.ClusterID,
			member.QueryEvidenceID,
			string(member.MembershipType),
			member.Rationale,
			string(member.Status),
			member.CreatedAt,
			member.UpdatedAt,
		).Error
		if err != nil {
			return mapQueryClusterMembershipPostgresError(err)
		}
		return nil
	})
}

func (r *SeoQueryClusterPostgres) RemoveMember(
	ctx context.Context,
	clusterID uint64,
	queryEvidenceID uint64,
	expectedVersion uint32,
) error {
	if clusterID == 0 || queryEvidenceID == 0 || expectedVersion == 0 {
		return repo.ErrSeoQueryClusterInvalid
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&seoQueryClusterRow{}).
			Where("id = ? AND lock_version = ? AND deleted_at IS NULL", clusterID, expectedVersion).
			Updates(map[string]interface{}{
				"lock_version": gorm.Expr("lock_version + 1"),
				"updated_at":   gorm.Expr("NOW()"),
			})
		if result.Error != nil {
			return mapQueryClusterPostgresError(result.Error)
		}
		if result.RowsAffected == 0 {
			return r.resolveQueryClusterUpdateFailure(ctx, tx, clusterID)
		}

		memberResult := tx.Model(&seoQueryClusterMemberRow{}).
			Where("cluster_id = ? AND query_evidence_id = ? AND status = ?", clusterID, queryEvidenceID, string(seo_strategy.QueryClusterMemberStatusActive)).
			Updates(map[string]interface{}{
				"status":     string(seo_strategy.QueryClusterMemberStatusRemoved),
				"updated_at": gorm.Expr("NOW()"),
			})
		if memberResult.Error != nil {
			return mapQueryClusterMembershipPostgresError(memberResult.Error)
		}
		if memberResult.RowsAffected == 0 {
			return repo.ErrSeoQueryClusterMemberNotFound
		}
		return nil
	})
}

func (r *SeoQueryClusterPostgres) ListMembers(
	ctx context.Context,
	clusterID uint64,
) ([]seo_strategy.LinkedQueryClusterMember, error) {
	var rows []linkedQueryClusterMemberRow
	err := impl.GetDB(ctx, r.db).WithContext(ctx).Raw(`
		SELECT
			m.objective_id AS member_objective_id,
			m.user_need_id AS member_user_need_id,
			m.cluster_id,
			m.query_evidence_id AS member_query_evidence_id,
			m.membership_type,
			m.rationale,
			m.status AS member_status,
			m.created_at AS member_created_at,
			m.updated_at AS member_updated_at,
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
		FROM seo_query_cluster_member m
		JOIN seo_query_evidence e ON e.id = m.query_evidence_id
		WHERE m.cluster_id = ?
		ORDER BY
			CASE m.status WHEN 'active' THEN 0 ELSE 1 END,
			CASE m.membership_type WHEN 'primary' THEN 0 WHEN 'included' THEN 1 ELSE 2 END,
			m.created_at ASC,
			e.id ASC
	`, clusterID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]seo_strategy.LinkedQueryClusterMember, 0, len(rows))
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
		items = append(items, seo_strategy.LinkedQueryClusterMember{
			Member: seo_strategy.SeoQueryClusterMember{
				ObjectiveID:     row.MemberObjectiveID,
				UserNeedID:      row.MemberUserNeedID,
				ClusterID:       row.ClusterID,
				QueryEvidenceID: row.MemberQueryEvidenceID,
				MembershipType:  seo_strategy.QueryClusterMembershipType(row.MembershipType),
				Rationale:       row.Rationale,
				Status:          seo_strategy.QueryClusterMemberStatus(row.MemberStatus),
				CreatedAt:       row.MemberCreatedAt,
				UpdatedAt:       row.MemberUpdatedAt,
			},
			Evidence: *evidence,
		})
	}
	return items, nil
}

func (r *SeoQueryClusterPostgres) resolveQueryClusterUpdateFailure(
	ctx context.Context,
	db *gorm.DB,
	id uint64,
) error {
	var count int64
	if err := db.WithContext(ctx).
		Model(&seoQueryClusterRow{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return repo.ErrSeoQueryClusterNotFound
	}
	return repo.ErrSeoQueryClusterConflict
}

func queryClusterToRow(cluster *seo_strategy.SeoQueryCluster) seoQueryClusterRow {
	return seoQueryClusterRow{
		ID:                 cluster.ID,
		ObjectiveID:        cluster.ObjectiveID,
		UserNeedID:         cluster.UserNeedID,
		IntentHypothesisID: cluster.IntentHypothesisID,
		Code:               cluster.Code,
		Name:               cluster.Name,
		SearchPromise:      cluster.SearchPromise,
		InclusionRule:      cluster.InclusionRule,
		ExclusionRule:      cluster.ExclusionRule,
		Status:             string(cluster.Status),
		Confidence:         string(cluster.Confidence),
		OwnerCode:          cluster.OwnerCode,
		ReviewDueAt:        cluster.ReviewDueAt,
		LockVersion:        cluster.LockVersion,
		CreatedAt:          cluster.CreatedAt,
		UpdatedAt:          cluster.UpdatedAt,
		DeletedAt:          cluster.DeletedAt,
	}
}

func queryClusterRowToDomain(row *seoQueryClusterRow) *seo_strategy.SeoQueryCluster {
	if row == nil {
		return nil
	}
	return &seo_strategy.SeoQueryCluster{
		ID:                 row.ID,
		ObjectiveID:        row.ObjectiveID,
		UserNeedID:         row.UserNeedID,
		IntentHypothesisID: row.IntentHypothesisID,
		Code:               row.Code,
		Name:               row.Name,
		SearchPromise:      row.SearchPromise,
		InclusionRule:      row.InclusionRule,
		ExclusionRule:      row.ExclusionRule,
		Status:             seo_strategy.QueryClusterStatus(row.Status),
		Confidence:         seo_strategy.ObjectiveConfidence(row.Confidence),
		OwnerCode:          row.OwnerCode,
		ReviewDueAt:        row.ReviewDueAt,
		LockVersion:        row.LockVersion,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
		DeletedAt:          row.DeletedAt,
	}
}

func mapQueryClusterPostgresError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case "23505":
		return fmt.Errorf("%w: %s", repo.ErrSeoQueryClusterConflict, pgErr.ConstraintName)
	case "23503":
		return fmt.Errorf("%w: %s", repo.ErrSeoQueryClusterMemberMismatch, pgErr.ConstraintName)
	case "23514", "22001", "22P02":
		return fmt.Errorf("%w: %s", repo.ErrSeoQueryClusterInvalid, pgErr.ConstraintName)
	default:
		return err
	}
}

func mapQueryClusterMembershipPostgresError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case "23505":
		return fmt.Errorf("%w: %s", repo.ErrSeoQueryClusterMembershipConflict, pgErr.ConstraintName)
	case "23503":
		return fmt.Errorf("%w: %s", repo.ErrSeoQueryClusterMemberMismatch, pgErr.ConstraintName)
	case "23514", "22001", "22P02":
		return fmt.Errorf("%w: %s", repo.ErrSeoQueryClusterInvalid, pgErr.ConstraintName)
	default:
		return err
	}
}
