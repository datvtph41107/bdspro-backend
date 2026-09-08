package postgresjob

import (
	"common/operation"
	"context"
	"errors"
	"time"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/usecase/generatedreport/processing"

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

func (s *Store) SaveJob(ctx context.Context, job processing.Job) (bool, error) {
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

// FindJobByReportID returns the durable processing job owned by one report.
//
// report_id is unique in report_jobs, so this lookup is used by Report
// acceptance to prove that an already accepted report has its execution job.
func (s *Store) FindJobByReportID(
	ctx context.Context,
	reportID uint64,
) (processing.Job, bool, error) {
	if s == nil || s.db == nil {
		return processing.Job{}, false, errors.New(
			"report job database is not configured",
		)
	}
	if reportID == 0 {
		return processing.Job{}, false, errors.New(
			"report ID is required",
		)
	}

	var row model
	result := s.db.WithContext(ctx).
		Where("report_id = ?", reportID).
		Limit(1).
		Find(&row)

	if result.Error != nil {
		return processing.Job{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return processing.Job{}, false, nil
	}

	return fromModel(row), true, nil
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
) (processing.Job, bool, error) {
	if s == nil || s.db == nil || workerID == "" {
		return processing.Job{}, false, errors.New("report job store is not configured")
	}

	var claimed model
	found := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model

		// An empty queue is normal worker state, not a database error.
		// Find + RowsAffected keeps idle polling silent while preserving
		// real PostgreSQL failures and the single-query SKIP LOCKED claim.
		result := tx.
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ?", processing.StatusPending).
			Where("available_at <= ?", now).
			Order("available_at ASC, created_at ASC").
			Limit(1).
			Find(&row)

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		row.Status = string(processing.StatusRunning)
		row.Attempts++
		row.ClaimVersion++
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
		return processing.Job{}, false, err
	}
	if !found {
		return processing.Job{}, false, nil
	}
	return fromModel(claimed), true, nil
}

/**
 * CompleteJob đánh dấu job hoàn tất và đưa report sang ready trong cùng transaction.
 */
func (s *Store) CompleteJob(
	ctx context.Context,
	job processing.Job,
	output processing.Output,
	now time.Time,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model{}).
			Where(
				"id = ? AND status = ? AND locked_by = ? AND claim_version = ?",
				job.ID,
				processing.StatusRunning,
				job.LockedBy,
				job.ClaimVersion,
			).
			Updates(map[string]any{
				"status":     processing.StatusCompleted,
				"locked_at":  nil,
				"locked_by":  "",
				"last_error": "",
				"updated_at": now,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return processing.ErrClaimLost
		}

		reportUpdates := map[string]any{
			"status":        enums.GeneratedReportStatusReady,
			"thumbnail_url": output.ThumbnailURL,
			"image_url":     output.ImageURL,
			"pdf_url":       output.PDFURL,
			"share_url":     output.ShareURL,
			"file_size":     output.FileSize,
			"error_message": "",
			"updated_at":    now,
		}
		// Compatibility generators may not report a format yet. In that case
		// preserve the accepted value instead of replacing it with empty text.
		if output.Format != "" {
			reportUpdates["format"] = output.Format
		}

		result = tx.Model(&qh_domain.QHUserReported{}).
			Where("id = ? AND user_id = ? AND deleted_at IS NULL", job.ReportID, job.UserID).
			Updates(reportUpdates)
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
	job processing.Job,
	lastError string,
	availableAt time.Time,
) error {
	result := s.db.WithContext(ctx).
		Model(&model{}).
		Where(
			"id = ? AND status = ? AND locked_by = ? AND claim_version = ?",
			job.ID,
			processing.StatusRunning,
			job.LockedBy,
			job.ClaimVersion,
		).
		Updates(map[string]any{
			"status":       processing.StatusPending,
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
		return processing.ErrClaimLost
	}
	return nil
}

/**
 * FailJob đánh dấu cả job và report là failed trong cùng transaction.
 */
func (s *Store) FailJob(
	ctx context.Context,
	job processing.Job,
	lastError string,
	now time.Time,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model{}).
			Where(
				"id = ? AND status = ? AND locked_by = ? AND claim_version = ?",
				job.ID,
				processing.StatusRunning,
				job.LockedBy,
				job.ClaimVersion,
			).
			Updates(map[string]any{
				"status":     processing.StatusFailed,
				"locked_at":  nil,
				"locked_by":  "",
				"last_error": lastError,
				"updated_at": now,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return processing.ErrClaimLost
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
			Where("status = ?", processing.StatusRunning).
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
				"status":       processing.StatusPending,
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

func toModel(job processing.Job) model {
	return model{
		ID:           job.ID,
		ReportID:     job.ReportID,
		UserID:       job.UserID,
		Operation:    string(job.Operation),
		OperationID:  job.OperationID,
		CommandKey:   job.CommandKey,
		Status:       string(job.Status),
		Attempts:     job.Attempts,
		AvailableAt:  job.AvailableAt,
		LockedAt:     job.LockedAt,
		LockedBy:     job.LockedBy,
		ClaimVersion: job.ClaimVersion,
		LastError:    job.LastError,
		CreatedAt:    job.CreatedAt,
		UpdatedAt:    job.UpdatedAt,
	}
}

func fromModel(row model) processing.Job {
	return processing.Job{
		ID:           row.ID,
		ReportID:     row.ReportID,
		UserID:       row.UserID,
		Operation:    operation.Code(row.Operation),
		OperationID:  row.OperationID,
		CommandKey:   row.CommandKey,
		Status:       processing.Status(row.Status),
		Attempts:     row.Attempts,
		AvailableAt:  row.AvailableAt,
		LockedAt:     row.LockedAt,
		LockedBy:     row.LockedBy,
		ClaimVersion: row.ClaimVersion,
		LastError:    row.LastError,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
