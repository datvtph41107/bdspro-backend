package postgre

import (
	base_enum "base/enum"
	_db "common/db"
	_routes "common/routes"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"strings"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.LeadRepo
type LeadPostgre struct {
	*_db.BaseRepo
}

func NewPostgreCustomer(db *gorm.DB) *LeadPostgre {
	return &LeadPostgre{BaseRepo: _db.NewBaseRepo(db)}
}

func (r *LeadPostgre) Search(c context.Context, ownerId uint64, ownerType base_enum.EOwnerOf, dto dto.LeadSearchDTO) ([]domain.LeadEntity, int64, error) {
	total := int64(0)
	customers := []domain.LeadEntity{}
	query := r.GetDB(c).
		Preload("Contact").
		Preload("Stage").
		Preload("Stage.Pipeline").
		Model(&domain.LeadEntity{}).
		Joins("JOIN tb_contact o ON o.id = customers.contact_id").
		Where("o.deleted_at is null and o.owner_id = ? and o.owner_of = ?", ownerId, ownerType)

	if dto.FullName != "" {
		query.Where("lower(full_name) like ?", "%"+strings.ToLower(dto.FullName)+"%")
	}

	if len(dto.Steps) > 0 {
		query.Where("step_id in (?)", dto.Steps)
	}

	if len(dto.StageIDs) > 0 {
		query.Where("stage_id in (?)", dto.StageIDs)
	}

	if len(dto.Sources) > 0 {
		query.Where("source_id in (?)", dto.Sources)
	}

	if len(dto.ChargePersonIds) > 0 {
		query.Where("charge_person_id in (?)", dto.ChargePersonIds)
	}

	if dto.FromDate != nil {
		query.Where("created_at >= ?", dto.FromDate)
	}

	if dto.ToDate != nil {
		query.Where("created_at <= ?", dto.ToDate)
	}

	if dto.LastUpdatedAt != nil {
		query.Where("updated_at >= ?", dto.LastUpdatedAt)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if dto.PipelineID != nil {
		query.Where("pipeline_id = ?", dto.PipelineID)
	}

	if err := query.
		Order("priority desc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&customers).
		Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *LeadPostgre) Create(ctx context.Context, entity *domain.LeadEntity) (*domain.LeadEntity, error) {
	err := r.GetDB(ctx).Create(entity).Error
	return entity, err
}

func (r *LeadPostgre) Update(ctx context.Context, id uint64, entity *domain.LeadEntity) (*domain.LeadEntity, error) {
	err := r.GetDB(ctx).
		Model(&domain.LeadEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("source", entity.Source).
		Update("assign_note", entity.AssignNote).
		Update("stage_id", entity.StageID).
		Update("priority", entity.Priority).
		Update("charge_person_id", entity.ChargePersonID).
		Update("charge_person_type", entity.ChargePersonType).
		Update("pipeline_id", entity.PipelineID).
		Error
	return entity, err
}

func (r *LeadPostgre) Delete(ctx context.Context, id uint64) error {
	return r.GetDB(ctx).
		WithContext(ctx).
		Model(&domain.LeadEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (r *LeadPostgre) GetByID(ctx context.Context, id uint64) (*domain.LeadEntity, error) {
	var customer domain.LeadEntity
	err := r.GetDB(ctx).
		Preload("Stage").
		Preload("Stage.Pipeline").
		Preload("Documents").
		Preload("ProductCares").
		Preload("Contact").
		Where("id = ? and deleted_at is null", id).
		First(&customer).Error

	if customer.Stage != nil {
		customer.Pipeline = customer.Stage.Pipeline
	}

	return &customer, err
}

func (r *LeadPostgre) Existed(c context.Context, entity *domain.LeadEntity) (bool, error) {
	return false, nil
}

func (r *LeadPostgre) UpdateNote(c context.Context, id uint64, note string) (*domain.LeadEntity, error) {
	customer := domain.LeadEntity{}
	err := r.GetDB(c).
		Model(&domain.LeadEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("note", note).
		Error
	return &customer, err
}

func (r *LeadPostgre) Assign(c context.Context, id uint64, dto *dto.LeadDTO) (*domain.LeadEntity, error) {
	customer := domain.LeadEntity{}
	err := r.GetDB(c).
		Model(&domain.LeadEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("charge_person_id", dto.ChargePersonId).
		Update("assign_note", dto.AssignNote).
		Error
	return &customer, err
}

func (r *LeadPostgre) SwitchStage(c context.Context, id uint64, stageID *uint64, note string) (*domain.LeadEntity, error) {
	customer := domain.LeadEntity{}
	err := r.GetDB(c).
		Model(&domain.LeadEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("stage_id", stageID).
		Update("stage_note", note).
		Error
	return &customer, err
}

func (r *LeadPostgre) ExistedWithStageID(c context.Context, stageID uint64) (bool, error) {
	var count int64
	err := r.GetDB(c).
		Model(&domain.LeadEntity{}).
		Where("stage_id = ? and deleted_at is null", stageID).
		Count(&count).Error
	return count > 0, err
}

func (r *LeadPostgre) ExistedWithPipelineID(c context.Context, pipelineID uint64) (bool, error) {
	var count int64
	err := r.GetDB(c).
		Model(&domain.LeadEntity{}).
		Where("pipeline_id = ? and deleted_at is null", pipelineID).
		Count(&count).Error
	return count > 0, err
}

func (r *LeadPostgre) GetOwnerWithRule(c context.Context) ([]dto.CustomerWithRuleDTO, error) {
	var customers []dto.CustomerWithRuleDTO

	query := `
	SELECT
		l.id AS lead_id,
		l.contact_id,
		l.pipeline_id,
		l.stage_id,
		l.charge_person_id,
		l.owner_id,
		l.owner_type,
		l.charge_person_type,
		l.priority,
		l.assign_note,
		l.stage_note,
		l.updated_at,

		-- Contact fields
		c.full_name,
		c.phone,
		c.email,
		c.address,
		c.birthday,
		c.note,

		-- Stage fields
		s.stage_name,

		-- Rule fields
		r.id AS rule_id,
		r.rule_name,
		r.condition,
		r.condition_value,
		r.trigger,
		r.trigger_value

	FROM customers l
	JOIN tb_contact c ON l.contact_id = c.id
	JOIN stage s ON l.stage_id = s.id
	JOIN rules r ON s.rule_id = r.id

	WHERE l.stage_id IS NOT NULL
	AND s.rule_id IS NOT NULL
	AND r.active = true;
	`
	err := r.GetDB(c).
		Raw(query).
		Scan(&customers).Error
	return customers, err
}

func (r *LeadPostgre) GetByContactID(ctx context.Context, contactID uint64) (*domain.LeadEntity, error) {
	var customer domain.LeadEntity
	err := r.GetDB(ctx).
		Where("contact_id = ? and deleted_at is null", contactID).
		First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &_routes.Except{
				Code:    404,
				Message: "Liên hệ không tồn tại",
			}
		}
		return nil, err
	}
	return &customer, err
}

// func (r *LeadPostgre) DetailByID(ctx context.Context, contactID uint64) (*domain.LeadEntity, error) {
// 	var customer domain.LeadEntity
// 	err := r.GetDB(ctx).
// 		Preload("Stage").
// 		Where("contact_id = ? and deleted_at is null", contactID).
// 		First(&customer).Error
// 	return &customer, err
// }

func (r *LeadPostgre) CareExpiredTime(c context.Context, organizationId uint64) (int64, error) {
	var count int64

	from := time.Now().AddDate(0, 0, -7) // 7 ngày trước
	to := from.AddDate(0, 0, 3)          // từ ngày đó + 3 ngày = khoảng 3 ngày

	err := r.GetDB(c).
		Debug().
		Model(&domain.LeadEntity{}).
		Where(`
			owner_id = ? 
			AND owner_type = ? 
			AND deleted_at IS NULL 
			AND updated_at BETWEEN ? AND ?`,
			organizationId, base_enum.EOwnerOfOrgnization, from, to).
		Count(&count).Error

	return count, err
}
