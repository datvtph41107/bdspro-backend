package postgres

import (
	_utils "common/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

type OrganizationMemberModel struct {
	ID             uint32    `gorm:"primaryKey;autoIncrement"`
	OrganizationID uint32    `gorm:"not null;index:idx_org_members_organization_id,priority:1;index:idx_org_members_user_org,priority:2;index:idx_org_members_status_org,priority:2"`
	UserID         uint64    `gorm:"not null;index:idx_org_members_user_org,priority:1"`
	Role           uint32    `gorm:"not null;index:idx_org_members_role"`
	RoleId         uint64    `gorm:"index:idx_org_members_role_id"`
	RoleKey        uint32    `gorm:"size:50;index:idx_org_members_role_key"`
	Status         uint32    `gorm:"default:1;index:idx_org_members_status_org,priority:1"`
	JoinedAt       time.Time `gorm:"autoCreateTime"`
	RemovedAt      *time.Time
	DeletedAt      *time.Time `gorm:"index"`
	IsDeleted      bool       `gorm:"default:false"`

	Organization *entity.Organization
}

type OrganizationMemberRoleModel struct {
	ID          uint32   `gorm:"column:id"`
	UserId      uint64   `gorm:"column:user_id"`
	Role        string   `gorm:"column:role"`
	Permissions []string `gorm:"column:permissions"`
}

func (OrganizationMemberModel) TableName() string {
	return "organization_members"
}

func OrganizationMemberModelToEntity(model *OrganizationMemberModel) *entity.OrganizationMember {
	roleID := model.RoleId
	if roleID == 0 {
		// `role` là cột production legacy. Đồng bộ đọc giúp dữ liệu cũ tiếp tục
		// chạy trong khi mọi write mới luôn ghi cả role và role_id.
		roleID = uint64(model.Role)
	}
	return &entity.OrganizationMember{
		ID:             model.ID,
		OrganizationID: model.OrganizationID,
		UserID:         uint64(model.UserID),
		RoleId:         roleID,
		RoleKey:        model.RoleKey,
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
		RoleKey:        entity.RoleKey,
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
		UserID         uint64     `gorm:"column:user_id"`
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
		// Role: &entity.OrganizationRole{
		// 	Id:   result.RoleId,
		// 	Name: result.RoleName,
		// 	Key:  result.RoleKey,
		// },
	}

	// if result.Permissions != "" {
	// 	permissions := strings.Split(result.Permissions, ",")
	// 	member.Role.Permissions = make([]*entity.OrganizationPermission, len(permissions))
	// 	for i, permission := range permissions {
	// 		member.Role.Permissions[i] = &entity.OrganizationPermission{
	// 			Key: permission,
	// 		}
	// 	}
	// }

	return member, nil
}

func (r *OrganizationMemberPostgresRepository) FindByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (*entity.OrganizationMember, error) {
	return r.findByUserIDAndOrganizationID(ctx, uint64(userId), organizationId)
}

// FindByUserIDAndOrganizationID giữ nguyên profile ID 64-bit cho các RPC
// internal mới; method legacy uint32 phía trên vẫn còn để không phá contract.
func (r *OrganizationMemberPostgresRepository) FindByUserIDAndOrganizationID(ctx context.Context, userID uint64, organizationID uint32) (*entity.OrganizationMember, error) {
	return r.findByUserIDAndOrganizationID(ctx, userID, organizationID)
}

func (r *OrganizationMemberPostgresRepository) findByUserIDAndOrganizationID(ctx context.Context, userID uint64, organizationID uint32) (*entity.OrganizationMember, error) {
	var organizationMemberModel OrganizationMemberModel
	err := r.db.WithContext(ctx).Where("user_id = ? AND organization_id = ? AND is_deleted = ?", userID, organizationID, false).First(&organizationMemberModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return OrganizationMemberModelToEntity(&organizationMemberModel), nil
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
	// Chỉ cập nhật status và roleKey
	updateData := map[string]interface{}{
		"status":   uint32(organizationMember.Status),
		"role_key": organizationMember.RoleKey,
	}

	err := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).
		Where("id = ? and deleted_at is null", organizationMember.ID).
		Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	// Lấy lại entity đã cập nhật
	updatedModel := &OrganizationMemberModel{}
	err = r.db.WithContext(ctx).Where("id = ?", organizationMember.ID).First(updatedModel).Error
	if err != nil {
		return nil, err
	}

	return OrganizationMemberModelToEntity(updatedModel), nil
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
		UserID         uint64          `gorm:"column:user_id"`
		OrganizationID uint32          `gorm:"column:organization_id"`
		RoleId         uint32          `gorm:"column:role"`
		Status         uint32          `gorm:"column:status"`
		JoinedAt       time.Time       `gorm:"column:joined_at"`
		RemovedAt      *time.Time      `gorm:"column:removed_at"`
		RoleKey        uint32          `gorm:"column:role_key"`
		RoleName       string          `gorm:"column:role_name"`
		Permissions    string          `gorm:"column:permissions"`
		DealMember     json.RawMessage `gorm:"column:deal_member"`
	}
	profileId := _utils.GetProfileIdWithContext(ctx)

	dealSelect := ""
	if dealId != nil {
		dealSelect = ", json_build_object('id', di.id,'role_id', di.role_id,'role_key', di.role_key,'status', di.status) as deal_member "
	}
	query := r.db.WithContext(ctx).
		Debug().
		Table("organization_members om").
		Select(fmt.Sprintf(`om.id, 
		om.user_id, 
		om.organization_id, 
		om.role, 
		om.role_key, 
		om.status, 
		om.joined_at, 
		om.removed_at, 
		or2.name as role_name, 
		string_agg(op.key, ',') as permissions
		%s`, dealSelect)).
		Joins("left join organization_roles or2 on om.role = or2.id").
		Joins("left join role_permissions rp on or2.id = rp.organization_role_model_id").
		Joins("left join organization_permissions op on rp.organization_permission_model_id = op.id")

	groupBy := "om.id, om.user_id, om.organization_id, om.role, om.status, om.joined_at, om.removed_at, or2.key, or2.name"
	if dealId != nil {
		query = query.Joins("left join deal_members di on om.user_id = di.member_id and di.deal_id = ?", dealId)
		groupBy += ", di.id, di.role_id, di.role_key, di.status"
	}

	query = query.
		Where("om.organization_id = ? AND om.deleted_at is null and om.user_id <> ?", organizationId, profileId).
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
			RoleKey:        result.RoleKey,
			Status:         entity.OrganizationMemberStatus(result.Status),
			JoinedAt:       result.JoinedAt,
			RemovedAt:      result.RemovedAt,
			RoleName:       result.RoleName,
			// Role: &entity.OrganizationRole{
			// 	Id:   result.RoleId,
			// 	Name: result.RoleName,
			// 	Key:  result.RoleKey,
			// },
		}

		// if result.Permissions != "" {
		// 	permissions := strings.Split(result.Permissions, ",")
		// 	member.Role.Permissions = make([]*entity.OrganizationPermission, len(permissions))
		// 	for j, permission := range permissions {
		// 		member.Role.Permissions[j] = &entity.OrganizationPermission{
		// 			Key: permission,
		// 		}
		// 	}
		// }

		if result.DealMember != nil {
			dealMember := &entity.DealMember{}
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

// CountMembersByOrganizationIds - Lấy count members của nhiều organizations trong 1 query
func (r *OrganizationMemberPostgresRepository) CountMembersByOrganizationIds(ctx context.Context, organizationIds []uint32) (map[uint32]uint32, error) {
	if len(organizationIds) == 0 {
		return make(map[uint32]uint32), nil
	}

	type Result struct {
		OrganizationID uint32 `gorm:"column:organization_id"`
		Count          uint32 `gorm:"column:count"`
	}

	var results []Result
	err := r.db.WithContext(ctx).
		Model(&OrganizationMemberModel{}).
		Select("organization_id, COUNT(*) as count").
		Where("organization_id IN ? AND is_deleted = ?", organizationIds, false).
		Group("organization_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Tạo map kết quả
	countMap := make(map[uint32]uint32)
	for _, result := range results {
		countMap[result.OrganizationID] = result.Count
	}

	// Đảm bảo tất cả organizationIds đều có trong map (với count = 0 nếu không có members)
	for _, orgId := range organizationIds {
		if _, exists := countMap[orgId]; !exists {
			countMap[orgId] = 0
		}
	}

	return countMap, nil
}
