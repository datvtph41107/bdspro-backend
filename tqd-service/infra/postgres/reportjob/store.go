package reportjobpostgres

import (
	"common/businessoperation"
	"context"
	"errors"
	"time"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/reportjob"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/**
 * Store dùng PostgreSQL làm durable queue cho generated report.
 */
type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) SaveJob(ctx context.Context, job reportjob.Job) (bool, error) {
	if s == nil || s.db == nil {
		return false, errors.New("report job database is not configured")
	}
	if job.ID == "" || job.ReportID == 0 || job.UserID == 0 || !job.Operation.IsValid() || job.OperationID == "" || job.CommandKey == "" {
		return false, errors.New("report job is invalid")
	}

	row := toModel(job)
	result := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&row)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}

/**
 * ClaimNextJob lấy một pending job bằng row lock.
 *
 * SKIP LOCKED cho phép nhiều worker chạy song song mà không lấy trùng job.
 */
func (s *Store) ClaimNextJob(
	ctx context.Context,
	workerID string,
	now time.Time,
) (reportjob.Job, bool, error) {
	if s == nil || s.db == nil || workerID == "" {
		return reportjob.Job{}, false, errors.New("report job store is not configured")
	}

	var claimed model
	found := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ?", reportjob.StatusPending).
			Where("available_at <= ?", now).
			Order("available_at ASC, created_at ASC").
			First(&row).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return err
		}

		row.Status = string(reportjob.StatusRunning)
		row.Attempts++
		row.LockedAt = &now
		row.LockedBy = workerID
		row.UpdatedAt = now
		if err := tx.Save(&row).Error; err != nil {
			return err
		}

		claimed = row
		found = true
		return nil
	})
	if err != nil {
		return reportjob.Job{}, false, err
	}
	if !found {
		return reportjob.Job{}, false, nil
	}
	return fromModel(claimed), true, nil
}

/**
 * CompleteJob đánh dấu job hoàn tất và đưa report sang ready trong cùng transaction.
 */
func (s *Store) CompleteJob(
	ctx context.Context,
	job reportjob.Job,
	output reportjob.Output,
	now time.Time,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model{}).
			Where("id = ? AND status = ?", job.ID, reportjob.StatusRunning).
			Updates(map[string]any{
				"status":     reportjob.StatusCompleted,
				"locked_at":  nil,
				"locked_by":  "",
				"last_error": "",
				"updated_at": now,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errors.New("report job is no longer running")
		}

		result = tx.Model(&qh_domain.QHUserReported{}).
			Where("id = ? AND user_id = ? AND deleted_at IS NULL", job.ReportID, job.UserID).
			Updates(map[string]any{
				"status":        enums.GeneratedReportStatusReady,
				"thumbnail_url": output.ThumbnailURL,
				"image_url":     output.ImageURL,
				"pdf_url":       output.PDFURL,
				"share_url":     output.ShareURL,
				"file_size":     output.FileSize,
				"error_message": "",
				"updated_at":    now,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errors.New("generated report was not found")
		}
		return nil
	})
}

func (s *Store) RetryJob(
	ctx context.Context,
	job reportjob.Job,
	lastError string,
	availableAt time.Time,
) error {
	result := s.db.WithContext(ctx).
		Model(&model{}).
		Where("id = ? AND status = ?", job.ID, reportjob.StatusRunning).
		Updates(map[string]any{
			"status":       reportjob.StatusPending,
			"available_at": availableAt,
			"locked_at":    nil,
			"locked_by":    "",
			"last_error":   lastError,
			"updated_at":   availableAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("report job is no longer running")
	}
	return nil
}

/**
 * FailJob đánh dấu cả job và report là failed trong cùng transaction.
 */
func (s *Store) FailJob(
	ctx context.Context,
	job reportjob.Job,
	lastError string,
	now time.Time,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model{}).
			Where("id = ? AND status = ?", job.ID, reportjob.StatusRunning).
			Updates(map[string]any{
				"status":     reportjob.StatusFailed,
				"locked_at":  nil,
				"locked_by":  "",
				"last_error": lastError,
				"updated_at": now,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errors.New("report job is no longer running")
		}

		result = tx.Model(&qh_domain.QHUserReported{}).
			Where("id = ? AND user_id = ? AND deleted_at IS NULL", job.ReportID, job.UserID).
			Updates(map[string]any{
				"status":        enums.GeneratedReportStatusFailed,
				"error_message": lastError,
				"updated_at":    now,
			})
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}

/**
 * ReleaseStaleJobs đưa các running job bị worker bỏ dở về pending.
 */
func (s *Store) ReleaseStaleJobs(
	ctx context.Context,
	lockedBefore time.Time,
	limit int,
) (int, error) {
	if limit <= 0 {
		limit = 100
	}

	count := 0
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rows []model
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ?", reportjob.StatusRunning).
			Where("locked_at IS NOT NULL AND locked_at < ?", lockedBefore).
			Order("locked_at ASC").
			Limit(limit).
			Find(&rows).Error
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}

		ids := make([]string, 0, len(rows))
		for _, row := range rows {
			ids = append(ids, row.ID)
		}
		result := tx.Model(&model{}).
			Where("id IN ?", ids).
			Updates(map[string]any{
				"status":       reportjob.StatusPending,
				"available_at": lockedBefore,
				"locked_at":    nil,
				"locked_by":    "",
				"updated_at":   lockedBefore,
			})
		if result.Error != nil {
			return result.Error
		}
		count = int(result.RowsAffected)
		return nil
	})
	return count, err
}

func toModel(job reportjob.Job) model {
	return model{
		ID:          job.ID,
		ReportID:    job.ReportID,
		UserID:      job.UserID,
		Operation:   string(job.Operation),
		OperationID: job.OperationID,
		CommandKey:  job.CommandKey,
		Status:      string(job.Status),
		Attempts:    job.Attempts,
		AvailableAt: job.AvailableAt,
		LockedAt:    job.LockedAt,
		LockedBy:    job.LockedBy,
		LastError:   job.LastError,
		CreatedAt:   job.CreatedAt,
		UpdatedAt:   job.UpdatedAt,
	}
}

func fromModel(row model) reportjob.Job {
	return reportjob.Job{
		ID:          row.ID,
		ReportID:    row.ReportID,
		UserID:      row.UserID,
		Operation:   businessoperation.Code(row.Operation),
		OperationID: row.OperationID,
		CommandKey:  row.CommandKey,
		Status:      reportjob.Status(row.Status),
		Attempts:    row.Attempts,
		AvailableAt: row.AvailableAt,
		LockedAt:    row.LockedAt,
		LockedBy:    row.LockedBy,
		LastError:   row.LastError,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
