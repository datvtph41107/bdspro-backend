package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

type OrganizationMemberModel struct {
	ID             uint32    `gorm:"primaryKey;autoIncrement"`
	OrganizationID uint32    `gorm:"not null;index:idx_org_members_organization_id,priority:1;index:idx_org_members_user_org,priority:2;index:idx_org_members_status_org,priority:2"`
	UserID         uint32    `gorm:"not null;index:idx_org_members_user_org,priority:1"`
	Role           uint32    `gorm:"not null;index:idx_org_members_role"`
	RoleId         uint64    `gorm:"index:idx_org_members_role_id"`
	Status         uint32    `gorm:"default:1;index:idx_org_members_status_org,priority:1"`
	JoinedAt       time.Time `gorm:"autoCreateTime"`
	RemovedAt      *time.Time
	DeletedAt      *time.Time `gorm:"index"`
	IsDeleted      bool       `gorm:"default:false"`

	Organization *OrganizationModel
}

type OrganizationMemberRoleModel struct {
	ID          uint32   `gorm:"column:id"`
	UserId      uint32   `gorm:"column:user_id"`
	Role        string   `gorm:"column:role"`
	Permissions []string `gorm:"column:permissions"`
}

func (OrganizationMemberModel) TableName() string {
	return "organization_members"
}

func OrganizationMemberModelToEntity(model *OrganizationMemberModel) *entity.OrganizationMember {
	return &entity.OrganizationMember{
		ID:             model.ID,
		OrganizationID: model.OrganizationID,
		UserID:         model.UserID,
		RoleId:         uint64(model.RoleId),
		Status:         entity.OrganizationMemberStatus(model.Status),
		JoinedAt:       model.JoinedAt,
		RemovedAt:      model.RemovedAt,
	}
}

func OrganizationMemberEntityToModel(entity *entity.OrganizationMember) *OrganizationMemberModel {
	return &OrganizationMemberModel{
		ID:             entity.ID,
		OrganizationID: entity.OrganizationID,
		UserID:         entity.UserID,
		Role:           uint32(entity.RoleId),
		RoleId:         entity.RoleId,
		Status:         uint32(entity.Status),
		JoinedAt:       entity.JoinedAt,
		RemovedAt:      entity.RemovedAt,
		IsDeleted:      false,
	}
}

func OrganizationMemberModelsToEntities(models []*OrganizationMemberModel) []*entity.OrganizationMember {
	entities := make([]*entity.OrganizationMember, len(models))
	for i, model := range models {
		entities[i] = OrganizationMemberModelToEntity(model)
	}
	return entities
}

// @bind: organization/internal/domain/repository.OrganizationMemberRepository
type OrganizationMemberPostgresRepository struct {
	db *gorm.DB
}

func NewOrganizationMemberPostgresRepository(db *gorm.DB) *OrganizationMemberPostgresRepository {
	return &OrganizationMemberPostgresRepository{db: db}
}

func (r *OrganizationMemberPostgresRepository) FindByUserIdAndOrganizationIdWithRole(ctx context.Context, userId, organizationId uint32) (*entity.OrganizationMember, error) {
	var result struct {
		ID             uint32     `gorm:"column:id"`
		UserID         uint32     `gorm:"column:user_id"`
		OrganizationID uint32     `gorm:"column:organization_id"`
		RoleId         uint32     `gorm:"column:role"`
		Status         uint32     `gorm:"column:status"`
		JoinedAt       time.Time  `gorm:"column:joined_at"`
		RemovedAt      *time.Time `gorm:"column:removed_at"`
		RoleKey        string     `gorm:"column:role_key"`
		RoleName       string     `gorm:"column:role_name"`
		Permissions    string     `gorm:"column:permissions"`
	}

	err := r.db.WithContext(ctx).
		Table("organization_members om").
		Select("om.id, om.user_id, om.organization_id, om.role, om.status, om.joined_at, om.removed_at, or2.key as role_key, or2.name as role_name, string_agg(op.key, ',') as permissions").
		Joins("left join organization_roles or2 on om.role = or2.id").
		Joins("left join role_permissions rp on or2.id = rp.organization_role_model_id").
		Joins("left join organization_permissions op on rp.organization_permission_model_id = op.id").
		Where("om.user_id = ? AND om.organization_id = ? AND om.is_deleted = ?", userId, organizationId, false).
		Group("om.id, om.user_id, om.organization_id, om.role, om.status, om.joined_at, om.removed_at, or2.key, or2.name").
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	member := &entity.OrganizationMember{
		ID:             result.ID,
		OrganizationID: result.OrganizationID,
		UserID:         result.UserID,
		RoleId:         uint64(result.RoleId),
		Status:         entity.OrganizationMemberStatus(result.Status),
		JoinedAt:       result.JoinedAt,
		RemovedAt:      result.RemovedAt,
		Role: &entity.OrganizationRole{
			Id:   result.RoleId,
			Name: result.RoleName,
			Key:  result.RoleKey,
		},
	}

	if result.Permissions != "" {
		permissions := strings.Split(result.Permissions, ",")
		member.Role.Permissions = make([]*entity.OrganizationPermission, len(permissions))
		for i, permission := range permissions {
			member.Role.Permissions[i] = &entity.OrganizationPermission{
				Key: permission,
			}
		}
	}

	return member, nil
}

func (r *OrganizationMemberPostgresRepository) FindByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (*entity.OrganizationMember, error) {
	var organizationMemberModel *OrganizationMemberModel
	err := r.db.WithContext(ctx).Where("user_id = ? AND organization_id = ? AND is_deleted = ?", userId, organizationId, false).First(&organizationMemberModel).Error
	if err != nil {
		return nil, err
	}
	return OrganizationMemberModelToEntity(organizationMemberModel), nil
}

func (r *OrganizationMemberPostgresRepository) CreateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	model := OrganizationMemberEntityToModel(organizationMember)
	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationMemberModelToEntity(model), nil
}

func (r *OrganizationMemberPostgresRepository) CreateOrganizationMemberBatch(ctx context.Context, organizationMembers []*entity.OrganizationMember) ([]*entity.OrganizationMember, error) {
	models := make([]*OrganizationMemberModel, len(organizationMembers))
	for i, member := range organizationMembers {
		models[i] = OrganizationMemberEntityToModel(member)
	}

	err := r.db.WithContext(ctx).Create(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.OrganizationMember, len(models))
	for i, model := range models {
		entities[i] = OrganizationMemberModelToEntity(model)
	}

	return entities, nil
}

func (r *OrganizationMemberPostgresRepository) FindByUserIdsAndOrganizationId(ctx context.Context, userIds []uint64, organizationId uint32) ([]*entity.OrganizationMember, error) {
	var organizationMembers []*OrganizationMemberModel
	err := r.db.WithContext(ctx).Where("user_id IN ? AND organization_id = ? AND is_deleted = ?", userIds, organizationId, false).Find(&organizationMembers).Error
	if err != nil {
		return nil, err
	}
	return OrganizationMemberModelsToEntities(organizationMembers), nil
}

func (r *OrganizationMemberPostgresRepository) FindById(ctx context.Context, id uint32) (*entity.OrganizationMember, error) {
	var organizationMemberModel *OrganizationMemberModel
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&organizationMemberModel).Error
	if err != nil {
		return nil, err
	}
	return OrganizationMemberModelToEntity(organizationMemberModel), nil
}

func (r *OrganizationMemberPostgresRepository) UpdateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	model := OrganizationMemberEntityToModel(organizationMember)
	err := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).Where("id = ?", organizationMember.ID).Updates(model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationMemberModelToEntity(model), nil
}

func (r *OrganizationMemberPostgresRepository) DeleteOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	err := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).Where("id = ?", organizationMember.ID).Delete(&OrganizationMemberModel{}).Error
	if err != nil {
		return nil, err
	}
	return organizationMember, nil
}

func (r *OrganizationMemberPostgresRepository) FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationMember, error) {
	var organizationMembers []*OrganizationMemberModel
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationId).Find(&organizationMembers).Error
	if err != nil {
		return nil, err
	}
	return OrganizationMemberModelsToEntities(organizationMembers), nil
}

func (r *OrganizationMemberPostgresRepository) CountMemberByOrganizationId(ctx context.Context, organizationId uint32) (uint32, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).Where("organization_id = ?", organizationId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}

func (r *OrganizationMemberPostgresRepository) FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int, roleId *uint32, dealId *uint64) ([]*entity.OrganizationMember, uint32, error) {
	var results []struct {
		ID             uint32          `gorm:"column:id"`
		UserID         uint32          `gorm:"column:user_id"`
		OrganizationID uint32          `gorm:"column:organization_id"`
		RoleId         uint32          `gorm:"column:role"`
		Status         uint32          `gorm:"column:status"`
		JoinedAt       time.Time       `gorm:"column:joined_at"`
		RemovedAt      *time.Time      `gorm:"column:removed_at"`
		RoleKey        string          `gorm:"column:role_key"`
		RoleName       string          `gorm:"column:role_name"`
		Permissions    string          `gorm:"column:permissions"`
		DealMember     json.RawMessage `gorm:"column:deal_member"`
	}

	dealSelect := ""
	if dealId != nil {
		dealSelect = ", row_to_json(di) as deal_member "
	}
	query := r.db.WithContext(ctx).
		Debug().
		Table("organization_members om").
		Select(fmt.Sprintf(`om.id, 
		om.user_id, 
		om.organization_id, 
		om.role, 
		om.status, 
		om.joined_at, 
		om.removed_at, 
		or2.key as role_key, 
		or2.name as role_name, 
		string_agg(op.key, ',') as permissions
		%s`, dealSelect)).
		Joins("left join organization_roles or2 on om.role = or2.id").
		Joins("left join role_permissions rp on or2.id = rp.organization_role_model_id").
		Joins("left join organization_permissions op on rp.organization_permission_model_id = op.id")

	groupBy := "om.id, om.user_id, om.organization_id, om.role, om.status, om.joined_at, om.removed_at, or2.key, or2.name"
	if dealId != nil {
		query = query.Joins("left join deal_invitations di on om.user_id = di.invitee_id and di.deal_id = ?", dealId)
		groupBy += ", di.id"
	}

	query = query.
		Where("om.organization_id = ? AND om.deleted_at is null", organizationId).
		Group(groupBy)

	// Add role filter if roleId is provided
	if roleId != nil {
		query = query.Where("om.role = ?", *roleId)
	}

	if err := query.
		Order("om.id desc").
		Offset(page * size).
		Limit(size).
		Find(&results).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]*entity.OrganizationMember, len(results))
	for i, result := range results {
		member := &entity.OrganizationMember{
			ID:             result.ID,
			OrganizationID: result.OrganizationID,
			UserID:         result.UserID,
			RoleId:         uint64(result.RoleId),
			Status:         entity.OrganizationMemberStatus(result.Status),
			JoinedAt:       result.JoinedAt,
			RemovedAt:      result.RemovedAt,
			RoleName:       result.RoleName,
			Role: &entity.OrganizationRole{
				Id:   result.RoleId,
				Name: result.RoleName,
				Key:  result.RoleKey,
			},
		}

		if result.Permissions != "" {
			permissions := strings.Split(result.Permissions, ",")
			member.Role.Permissions = make([]*entity.OrganizationPermission, len(permissions))
			for j, permission := range permissions {
				member.Role.Permissions[j] = &entity.OrganizationPermission{
					Key: permission,
				}
			}
		}

		if result.DealMember != nil {
			dealMember := &entity.DealInvitation{}
			err := json.Unmarshal(result.DealMember, dealMember)
			if err != nil {
				return nil, 0, err
			}
			member.DealMember = dealMember
		}

		entities[i] = member
	}

	var count int64
	countQuery := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).Where("organization_id = ? AND is_deleted = ?", organizationId, false)
	if roleId != nil {
		countQuery = countQuery.Where("role = ?", *roleId)
	}
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(count), nil
}

func (r *OrganizationMemberPostgresRepository) AddParticipant(ctx context.Context, participant *entity.OrganizationMember) error {
	model := OrganizationMemberEntityToModel(participant)
	model.IsDeleted = false
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *OrganizationMemberPostgresRepository) HasExactConversationWithUserIDs(ctx context.Context, userIDs []uint64) (uint64, error) {
	var conversationID uint64
	err := r.db.WithContext(ctx).
		Model(&OrganizationMemberModel{}).
		Select("conversation_id").
		Group("conversation_id").
		Having("COUNT(*) = ? AND COUNT(CASE WHEN user_id IN ? AND is_deleted = ? THEN 1 END) = ?", len(userIDs), userIDs, false, len(userIDs)).
		Limit(1).
		Scan(&conversationID).Error
	if err != nil {
		return 0, err
	}
	return conversationID, nil
}

// Add soft delete function
func (r *OrganizationMemberPostgresRepository) DeleteMember(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
			"removed_at": time.Now(),
		}).Error
}
