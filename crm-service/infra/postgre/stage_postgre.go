package postgre

import (
	"context"
	"crm/infra/impl"
	"crm/internal/domain"
	"crm/internal/dto"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.StageRepo
type PostgreStage struct {
	db *gorm.DB
}

func NewPostgreStage(db *gorm.DB) *PostgreStage {
	return &PostgreStage{db: db}
}

func (r *PostgreStage) Search(c context.Context, pipelineId uint64, dto dto.StageSearchDTO) ([]domain.StageEntity, int64, error) {
	var entities []domain.StageEntity
	query := r.db.WithContext(c).Model(&domain.StageEntity{})

	query = query.Where("pipeline_id = ?", pipelineId)

	if dto.StageName != "" {
		query = query.Where("stage_name LIKE ?", "%"+dto.StageName+"%")
	}

	if err := query.
		Order("order_number asc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&entities).
		Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *PostgreStage) Create(c context.Context, entity *domain.StageEntity) (*domain.StageEntity, error) {
	if err := r.db.WithContext(c).Create(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *PostgreStage) Update(c context.Context, entity *domain.StageEntity) (*domain.StageEntity, error) {
	if err := r.db.WithContext(c).Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *PostgreStage) Delete(c context.Context, id uint64) error {
	return r.db.WithContext(c).
		Model(&domain.StageEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (r *PostgreStage) GetByID(c context.Context, id uint64) (*domain.StageEntity, error) {
	var entity domain.StageEntity
	if err := r.db.WithContext(c).
		Model(&domain.StageEntity{}).
		Where("id = ?", id).
		First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PostgreStage) GetOneStageDefault(c context.Context, pipelineId uint64) (*domain.StageEntity, error) {
	var entity domain.StageEntity
	if err := r.db.WithContext(c).
		Model(&domain.StageEntity{}).
		Where("pipeline_id = ? and deleted_at is null", pipelineId).
		Order("order_number asc").
		First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PostgreStage) BulkUpdate(ctx context.Context, removeIds []uint64, stages []domain.StageEntity) error {
	// 	if len(stages) == 0 {
	// 		return nil
	// 	}

	// 	var values []string
	// 	var args []interface{}
	// 	var ids []string // dùng cho WHERE NOT IN (...)

	// 	for _, s := range stages {
	// 		// ép kiểu từng phần tử trực tiếp trong VALUES
	// 		values = append(values, "(?::bigint, ?::text, ?::int, ?::text, ?::int)")
	// 		args = append(args, s.ID, s.StageName, s.Color, s.ColorRGB, s.OrderNumber)
	// 		ids = append(ids, fmt.Sprintf("%d", s.ID))
	// 	}

	// 	query := fmt.Sprintf(`
	// 		UPDATE stage AS s SET
	// 			stage_name = u.stage_name,
	// 			color = u.color,
	// 			color_rgb = u.color_rgb,
	// 			order_number = u.order_number
	// 		FROM (VALUES %s) AS u(id, stage_name, color, color_rgb, order_number)
	// 		WHERE u.id = s.id
	// 	`, strings.Join(values, ","))

	// 	deleteQuery := fmt.Sprintf(`
	// 	UPDATE stage SET deleted_at = NOW()
	// 	WHERE owner_id = ? AND owner_type = ? AND id NOT IN (%s)
	// `, strings.Join(ids, ","))

	// 	err := impl.GetDB(ctx, r.db).Exec(query, args...).Error
	// 	if err != nil {
	// 		return err
	// 	}

	// 	err = impl.GetDB(ctx, r.db).Exec(deleteQuery, profileId, ownerType, strings.Join(ids, ",")).Error
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return nil
	tx := impl.GetDB(ctx, r.db)

	// Soft delete
	if len(removeIds) > 0 {
		if err := tx.Model(&domain.StageEntity{}).
			Where("id IN ?", removeIds).
			Update("deleted_at", time.Now()).Error; err != nil {
			return err
		}
	}

	// Save (insert or update)
	if len(stages) > 0 {
		if err := tx.Save(stages).Error; err != nil {
			return err
		}
	}

	return nil
}
