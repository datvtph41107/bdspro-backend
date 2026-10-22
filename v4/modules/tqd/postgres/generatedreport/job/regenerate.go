package postgresjob

import (
	"context"
	"fmt"
	qh "mdn4/modules/tqd/entity/qh"
	"mdn4/modules/tqd/enums"
	"mdn4/modules/tqd/services/generatedreport/processing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RegenerateReport requeues the original accepted job, preserving its command,
// usage and File effect identity. It never creates another report or usage event.
// Lock order matches worker completion: job first, then its report.
func (s *Store) RegenerateReport(ctx context.Context, userID, reportID uint64, now time.Time) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("report job database is not configured")
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var job model
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("report_id = ? AND user_id = ?", reportID, userID).First(&job).Error; err != nil {
			return err
		}
		var report qh.QHUserReported
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", reportID, userID).First(&report).Error; err != nil {
			return err
		}
		// A response-loss retry observes the existing scheduled/running recovery.
		if job.Status == string(processing.StatusPending) || job.Status == string(processing.StatusRunning) {
			return nil
		}
		if !report.Status.CanRegenerate() ||
			(job.Status != string(processing.StatusCompleted) && job.Status != string(processing.StatusFailed)) {
			return fmt.Errorf("report cannot regenerate in current status")
		}
		if err := tx.Model(&job).Updates(map[string]any{
			"status": processing.StatusPending, "attempts": 0, "available_at": now,
			"locked_at": nil, "locked_by": "", "last_error": "", "updated_at": now,
		}).Error; err != nil {
			return err
		}
		return tx.Model(&report).Updates(map[string]any{
			"status": enums.GeneratedReportStatusProcessing, "error_message": "", "updated_at": now, "expires_at": nil,
		}).Error
	})
}
