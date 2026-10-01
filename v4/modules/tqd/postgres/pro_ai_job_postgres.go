package postgres

import (
	_db "common/db"
	"context"
	"strings"
	"time"

	qh "userv4/modules/tqd/entity/qh"
	qh_dto "userv4/modules/tqd/dto/qh"
	"userv4/modules/tqd/enums"
	repo "userv4/modules/tqd/repository"

	"gorm.io/gorm"
)

type ProAIJobRepo struct {
	db *_db.TransactionRepo
}

func NewProAIJobRepo(db *_db.TransactionRepo) repo.IProAIJobRepo {
	return &ProAIJobRepo{db: db}
}

func (r *ProAIJobRepo) Create(ctx context.Context, entity *qh.ProAIJob) error {
	return r.db.GetDB(ctx).Create(entity).Error
}

func (r *ProAIJobRepo) GetByID(ctx context.Context, id uint64) (*qh.ProAIJob, error) {
	var entity qh.ProAIJob
	err := r.db.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ProAIJobRepo) GetList(ctx context.Context, req *qh_dto.ListProAIJobsRequest) ([]qh.ProAIJob, int64, error) {
	var jobs []qh.ProAIJob
	var total int64

	db := r.db.GetDB(ctx).Model(&qh.ProAIJob{}).Where("deleted_at IS NULL")

	if req.Search != "" {
		searchTerm := "%" + strings.ToLower(req.Search) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(source_folder_name) LIKE ?", searchTerm, searchTerm)
	}
	if req.JobType != "" {
		db = db.Where("job_type = ?", req.JobType)
	}
	if req.ProcessStatus != "" {
		statuses := strings.Split(req.ProcessStatus, ",")
		trimmed := make([]string, 0, len(statuses))
		for _, s := range statuses {
			s = strings.TrimSpace(s)
			if s != "" {
				trimmed = append(trimmed, s)
			}
		}
		if len(trimmed) == 1 {
			db = db.Where("process_status = ?", trimmed[0])
		} else if len(trimmed) > 1 {
			db = db.Where("process_status IN ?", trimmed)
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("created_at DESC").
		Find(&jobs).Error
	return jobs, total, err
}

func (r *ProAIJobRepo) UpdateProcessStatusByProjectID(ctx context.Context, planningProjectID uint64, status uint32) error {
	now := time.Now()
	updates := map[string]interface{}{
		"process_status": status,
		"updated_at":     now,
	}
	switch status {
	case enums.PlanningProcessStatusClassified, enums.PlanningProcessStatusFailed:
		updates["completed_at"] = gorm.Expr("COALESCE(completed_at, ?)", now)
	case enums.PlanningProcessStatusApproved:
		updates["completed_at"] = gorm.Expr("COALESCE(completed_at, ?)", now)
		updates["approved_at"] = gorm.Expr("COALESCE(approved_at, ?)", now)
	case enums.PlanningProcessStatusPending:
		updates["completed_at"] = nil
		updates["approved_at"] = nil
	}
	return r.db.GetDB(ctx).Model(&qh.ProAIJob{}).
		Where("planning_project_id = ? AND deleted_at IS NULL", planningProjectID).
		Updates(updates).Error
}
