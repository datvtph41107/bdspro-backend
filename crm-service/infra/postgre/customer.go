package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"strings"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.CustomerRepo
type PostgreCustomer struct {
	db *gorm.DB
}

func NewPostgreCustomer(db *gorm.DB) *PostgreCustomer {
	return &PostgreCustomer{db: db}
}

func (r *PostgreCustomer) Search(c context.Context, ownerId uint64, ownerType enums.EOwnerType, dto dto.CustomerSearchDTO) ([]domain.CustomerEntity, int64, error) {
	total := int64(0)
	customers := []domain.CustomerEntity{}
	query := r.db.WithContext(c).
		Model(&domain.CustomerEntity{}).
		Where("deleted_at is null and owner_id = ? and owner_type = ?", ownerId, ownerType)

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

	if len(dto.ChargeIds) > 0 {
		query.Where("charge_id in (?)", dto.ChargeIds)
	}

	if dto.FromDate != nil {
		query.Where("created_at >= ?", dto.FromDate)
	}

	if dto.ToDate != nil {
		query.Where("created_at <= ?", dto.ToDate)
	}

	if err := query.
		Order("created_at desc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&customers).
		Error; err != nil {
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *PostgreCustomer) Create(c context.Context, entity *domain.CustomerEntity) (*domain.CustomerEntity, error) {
	err := r.db.WithContext(c).Create(entity).Error
	return entity, err
}

func (r *PostgreCustomer) Update(ctx context.Context, id uint64, customer *domain.CustomerEntity) (*domain.CustomerEntity, error) {
	err := r.db.WithContext(ctx).
		Model(&domain.CustomerEntity{}).
		Where("id = ? and deleted_at is null", id).
		Updates(customer).
		Error
	return customer, err
}

func (r *PostgreCustomer) Delete(ctx context.Context, id uint64) error {
	return r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (r *PostgreCustomer) GetByID(ctx context.Context, id uint64) (*domain.CustomerEntity, error) {
	var customer domain.CustomerEntity
	err := r.db.WithContext(ctx).
		Where("id = ? and deleted_at is null", id).
		First(&customer).Error
	return &customer, err
}

func (r *PostgreCustomer) Existed(c context.Context, entity *domain.CustomerEntity) (bool, error) {
	return false, nil
}

func (r *PostgreCustomer) UpdateNote(c context.Context, id uint64, note string) (*domain.CustomerEntity, error) {
	customer := domain.CustomerEntity{}
	err := r.db.WithContext(c).
		Model(&domain.CustomerEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("note", note).
		Error
	return &customer, err
}

func (r *PostgreCustomer) Assign(c context.Context, id uint64, dto *dto.CustomerDTO) (*domain.CustomerEntity, error) {
	customer := domain.CustomerEntity{}
	err := r.db.WithContext(c).
		Model(&domain.CustomerEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("charge_id", dto.ChargeId).
		Update("assign_note", dto.AssignNote).
		Error
	return &customer, err
}

func (r *PostgreCustomer) SwitchStage(c context.Context, id uint64, stageID *uint64) (*domain.CustomerEntity, error) {
	customer := domain.CustomerEntity{}
	err := r.db.WithContext(c).
		Model(&domain.CustomerEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("stage_id", stageID).
		Error
	return &customer, err
}

func (r *PostgreCustomer) ExistedWithStageID(c context.Context, stageID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(c).
		Model(&domain.CustomerEntity{}).
		Where("stage_id = ? and deleted_at is null", stageID).
		Count(&count).Error
	return count > 0, err
}

func (r *PostgreCustomer) ExistedWithPipelineID(c context.Context, pipelineID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(c).
		Model(&domain.CustomerEntity{}).
		Where("pipeline_id = ? and deleted_at is null", pipelineID).
		Count(&count).Error
	return count > 0, err
}

func (r *PostgreCustomer) GetOwnerWithRule(c context.Context) ([]dto.CustomerWithRuleDTO, error) {
	var customers []dto.CustomerWithRuleDTO

	query := `
	SELECT
		c.id AS customer_id,
		c.full_name,
		c.phone,
		c.email,
		c.address,
		c.note,
		c.birthday,
		c.pipeline_id,
		c.stage_id,
		c.charge_id,
		c.updated_at,
		c.owner_id,
		c.owner_type,
		s.stage_name,
		r.id AS rule_id,
		r.rule_name,
		r.condition,
		r.condition_value,
		r.trigger,
		r.trigger_value
		FROM customers c
		JOIN stage s ON c.stage_id = s.id
		JOIN rules r ON s.rule_id = r.id
		WHERE c.stage_id IS NOT NULL
		AND s.rule_id IS NOT NULL
		AND r.active = true;
	`
	err := r.db.WithContext(c).
		Raw(query).
		Scan(&customers).Error
	return customers, err
}
