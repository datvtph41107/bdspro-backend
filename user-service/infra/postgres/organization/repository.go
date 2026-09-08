package organization

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	domain "user/internal/domain/organization"
	usecase "user/internal/usecase/organization"
)

// @bind: user/internal/usecase/organization.Repository
type Repository struct{ database *gorm.DB }

func NewRepository(database *gorm.DB) *Repository { return &Repository{database: database} }

// projectionRow is infrastructure-owned. Explicit column mapping keeps SQL
// naming concerns out of the domain and avoids silent zero values while GORM
// scans computed/aliased projection fields.
type projectionRow struct {
	ID                 uint64     `gorm:"column:id"`
	Code               string     `gorm:"column:code"`
	Name               string     `gorm:"column:name"`
	Type               string     `gorm:"column:type"`
	TaxCode            string     `gorm:"column:tax_code"`
	Phone              string     `gorm:"column:phone"`
	Email              string     `gorm:"column:email"`
	Address            string     `gorm:"column:address"`
	Website            string     `gorm:"column:website"`
	Description        string     `gorm:"column:description"`
	Status             string     `gorm:"column:status"`
	VerificationStatus string     `gorm:"column:verification_status"`
	WarningLevel       string     `gorm:"column:warning_level"`
	WarningCount       uint32     `gorm:"column:warning_count"`
	OwnerProfileID     uint64     `gorm:"column:owner_profile_id"`
	CreatedBy          uint64     `gorm:"column:created_by"`
	UpdatedBy          uint64     `gorm:"column:updated_by"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"`
	ArchivedAt         *time.Time `gorm:"column:archived_at"`
	ActiveMemberCount  uint64     `gorm:"column:active_member_count"`
}

func (row projectionRow) domain() domain.Organization {
	return domain.Organization{
		ID: row.ID, Code: row.Code, Name: row.Name, Type: row.Type,
		TaxCode: row.TaxCode, Phone: row.Phone, Email: row.Email,
		Address: row.Address, Website: row.Website, Description: row.Description,
		Status: row.Status, VerificationStatus: row.VerificationStatus,
		WarningLevel: row.WarningLevel, WarningCount: row.WarningCount,
		OwnerProfileID: row.OwnerProfileID, CreatedBy: row.CreatedBy, UpdatedBy: row.UpdatedBy,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt, ArchivedAt: row.ArchivedAt,
		ActiveMemberCount: row.ActiveMemberCount,
	}
}

func projectionSelect() string {
	return `o.id, o.code, o.name, o.type, COALESCE(o.tax_code, '') AS tax_code,
		COALESCE(o.phone, '') AS phone, COALESCE(o.email, '') AS email,
		COALESCE(o.address, '') AS address, COALESCE(o.website, '') AS website,
		COALESCE(o.description, '') AS description, o.status, o.verification_status,
		o.warning_level, o.warning_count, o.owner_profile_id, o.created_by, o.updated_by,
		o.created_at, o.updated_at, o.archived_at,
		(SELECT COUNT(*) FROM organization_members m
		 WHERE m.organization_id = o.id AND m.status = 'active') AS active_member_count`
}

func mapWriteError(operation string, err error) error {
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return fmt.Errorf("%w: code or tax code already exists", usecase.ErrConflict)
	}
	return fmt.Errorf("%s organization: %w", operation, err)
}

func (r *Repository) List(ctx context.Context, query domain.ListQuery) ([]domain.Organization, uint64, error) {
	base := r.database.WithContext(ctx).Table("organizations AS o").Where("o.archived_at IS NULL")
	if query.SearchText != "" {
		pattern := "%" + strings.ToLower(query.SearchText) + "%"
		base = base.Where("lower(o.name) LIKE ? OR lower(o.code) LIKE ? OR lower(COALESCE(o.tax_code, '')) LIKE ? OR lower(COALESCE(o.email, '')) LIKE ? OR COALESCE(o.phone, '') LIKE ?", pattern, pattern, pattern, pattern, "%"+query.SearchText+"%")
	}
	if query.Type != "" {
		base = base.Where("o.type = ?", query.Type)
	}
	if query.Status != "" {
		base = base.Where("o.status = ?", query.Status)
	}
	if query.VerificationStatus != "" {
		base = base.Where("o.verification_status = ?", query.VerificationStatus)
	}
	if query.WarningLevel != "" {
		base = base.Where("o.warning_level = ?", query.WarningLevel)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count organizations: %w", err)
	}
	var rows []projectionRow
	if err := base.Select(projectionSelect()).Order("o.created_at DESC, o.id DESC").
		Offset(int(query.Page * query.Size)).Limit(int(query.Size)).Scan(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list organizations: %w", err)
	}
	organizations := make([]domain.Organization, 0, len(rows))
	for _, row := range rows {
		organizations = append(organizations, row.domain())
	}
	return organizations, uint64(total), nil
}

func (r *Repository) Get(ctx context.Context, id uint64) (domain.Organization, error) {
	var row projectionRow
	result := r.database.WithContext(ctx).Table("organizations AS o").Select(projectionSelect()).
		Where("o.id = ? AND o.archived_at IS NULL", id).Scan(&row)
	if result.Error != nil {
		return domain.Organization{}, fmt.Errorf("get organization: %w", result.Error)
	}
	if result.RowsAffected == 0 || row.ID == 0 {
		return domain.Organization{}, usecase.ErrNotFound
	}
	return row.domain(), nil
}

func (r *Repository) ListForProfile(ctx context.Context, profileID uint64) ([]domain.Organization, error) {
	var rows []projectionRow
	if err := r.database.WithContext(ctx).Table("organizations AS o").
		Select(projectionSelect()).
		Joins("JOIN organization_members member ON member.organization_id = o.id").
		Where("member.profile_id = ? AND member.status = 'active' AND member.removed_at IS NULL", profileID).
		Where("o.status = 'active' AND o.archived_at IS NULL").
		Order("o.name ASC, o.id ASC").Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list organizations for profile: %w", err)
	}
	organizations := make([]domain.Organization, 0, len(rows))
	for _, row := range rows {
		organizations = append(organizations, row.domain())
	}
	return organizations, nil
}

func (r *Repository) Create(ctx context.Context, organization domain.Organization) (domain.Organization, error) {
	err := r.database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inserted struct {
			ID uint64
		}
		row := tx.Raw(`
			INSERT INTO organizations
			(code, name, type, tax_code, phone, email, address, website, description,
			 status, verification_status, warning_level, warning_count,
			 owner_profile_id, created_by, updated_by)
			VALUES (?, ?, ?, NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''),
			        NULLIF(?, ''), NULLIF(?, ''), ?, ?, ?, 0, ?, ?, ?)
			RETURNING id`, organization.Code, organization.Name, organization.Type,
			organization.TaxCode, organization.Phone, organization.Email,
			organization.Address, organization.Website, organization.Description,
			organization.Status, organization.VerificationStatus, organization.WarningLevel,
			organization.OwnerProfileID, organization.CreatedBy, organization.CreatedBy,
		).Scan(&inserted)
		if row.Error != nil {
			return row.Error
		}
		if inserted.ID == 0 {
			return errors.New("create organization returned no identity")
		}
		organization.ID = inserted.ID
		return tx.Exec(`
			INSERT INTO organization_members
			(organization_id, profile_id, role_key, status, joined_at)
			VALUES (?, ?, 'owner', 'active', NOW())`, organization.ID, organization.OwnerProfileID).Error
	})
	if err != nil {
		return domain.Organization{}, mapWriteError("create", err)
	}
	return r.Get(ctx, organization.ID)
}

func (r *Repository) Update(ctx context.Context, update domain.Update) (domain.Organization, error) {
	values := map[string]any{"updated_by": update.ActorID, "updated_at": time.Now()}
	set := func(key string, value *string) {
		if value != nil {
			values[key] = strings.TrimSpace(*value)
		}
	}
	set("code", update.Code)
	set("name", update.Name)
	set("type", update.Type)
	set("tax_code", update.TaxCode)
	set("phone", update.Phone)
	set("email", update.Email)
	set("address", update.Address)
	set("website", update.Website)
	set("description", update.Description)
	set("status", update.Status)
	set("verification_status", update.VerificationStatus)
	set("warning_level", update.WarningLevel)

	result := r.database.WithContext(ctx).Table("organizations").
		Where("id = ? AND archived_at IS NULL", update.ID).Updates(values)
	if result.Error != nil {
		return domain.Organization{}, mapWriteError("update", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.Organization{}, usecase.ErrNotFound
	}
	return r.Get(ctx, update.ID)
}

func (r *Repository) Archive(ctx context.Context, id, actorID uint64) error {
	return r.database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var activeSubscriptions int64
		if err := tx.Table("catalog_subscriptions").
			Where("subject_kind = 'organization' AND subject_id = ? AND status IN ('active', 'past_due')", fmt.Sprint(id)).
			Count(&activeSubscriptions).Error; err != nil {
			return fmt.Errorf("check organization subscriptions: %w", err)
		}
		if activeSubscriptions > 0 {
			return usecase.ErrActiveSubscription
		}
		result := tx.Exec(`
			UPDATE organizations
			SET status = 'archived', archived_at = NOW(), updated_at = NOW(), updated_by = ?
			WHERE id = ? AND archived_at IS NULL`, actorID, id)
		if result.Error != nil {
			return fmt.Errorf("archive organization: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return usecase.ErrNotFound
		}
		return tx.Exec(`
			UPDATE organization_members SET status = 'removed', removed_at = NOW(), updated_at = NOW()
			WHERE organization_id = ? AND status <> 'removed'`, id).Error
	})
}

func (r *Repository) CheckMembers(ctx context.Context, organizationID uint64, profileIDs []uint64) ([]domain.MemberCheck, error) {
	if len(profileIDs) == 0 {
		return []domain.MemberCheck{}, nil
	}
	type memberRow struct {
		MemberID  uint64     `gorm:"column:member_id"`
		ProfileID uint64     `gorm:"column:profile_id"`
		JoinedAt  *time.Time `gorm:"column:joined_at"`
		RoleKey   string     `gorm:"column:role_key"`
		Status    string     `gorm:"column:status"`
	}
	var rows []memberRow
	if err := r.database.WithContext(ctx).Table("organization_members AS member").
		Select("member.id AS member_id, member.profile_id, member.joined_at, member.role_key, member.status").
		Joins("JOIN organizations organization ON organization.id = member.organization_id").
		Where("member.organization_id = ? AND member.profile_id IN ?", organizationID, profileIDs).
		Where("member.status = 'active' AND member.removed_at IS NULL").
		Where("organization.status = 'active' AND organization.archived_at IS NULL").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("check organization members: %w", err)
	}
	found := make(map[uint64]domain.MemberCheck, len(rows))
	for _, row := range rows {
		joinedAt := time.Time{}
		if row.JoinedAt != nil {
			joinedAt = *row.JoinedAt
		}
		found[row.ProfileID] = domain.MemberCheck{
			ProfileID: row.ProfileID, IsMember: true, JoinedAt: joinedAt,
			MemberID: row.MemberID, RoleKey: row.RoleKey, Status: row.Status,
		}
	}
	checks := make([]domain.MemberCheck, 0, len(profileIDs))
	for _, profileID := range profileIDs {
		if check, exists := found[profileID]; exists {
			checks = append(checks, check)
			continue
		}
		checks = append(checks, domain.MemberCheck{ProfileID: profileID})
	}
	return checks, nil
}

var _ usecase.Repository = (*Repository)(nil)
