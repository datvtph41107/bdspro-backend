package requeststats

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Store persists buffered request stats/audit rows into the local database.
type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) FlushCounts(ctx context.Context, logs []*RequestStatsLog) error {
	if s == nil || s.db == nil || len(logs) == 0 {
		return nil
	}
	db := s.db.WithContext(ctx)
	const batchSize = 100
	for i := 0; i < len(logs); i += batchSize {
		end := i + batchSize
		if end > len(logs) {
			end = len(logs)
		}
		batch := logs[i:end]
		placeholders := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*7)
		for _, item := range batch {
			placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?)")
			args = append(args,
				item.Method,
				item.Path,
				item.Service,
				item.StatusCode,
				item.Count,
				item.Timestamp,
				item.BucketHour,
			)
		}
		sql := fmt.Sprintf(`
			INSERT INTO request_stats_log (method, path, service, status_code, count, timestamp, bucket_hour)
			VALUES %s
			ON CONFLICT (method, path, service, status_code, bucket_hour)
			DO UPDATE SET
				count = request_stats_log.count + EXCLUDED.count,
				"timestamp" = EXCLUDED."timestamp"
		`, strings.Join(placeholders, ", "))
		if err := db.Exec(sql, args...).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) FlushUserInteractions(ctx context.Context, logs []*UserInteractionLog) error {
	if s == nil || s.db == nil || len(logs) == 0 {
		return nil
	}
	db := s.db.WithContext(ctx)
	const batchSize = 100
	for i := 0; i < len(logs); i += batchSize {
		end := i + batchSize
		if end > len(logs) {
			end = len(logs)
		}
		batch := logs[i:end]
		placeholders := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*5)
		for _, item := range batch {
			placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
			args = append(args,
				item.ProfileID,
				item.OrganizationID,
				item.Count,
				item.Timestamp,
				item.BucketHour,
			)
		}
		sql := fmt.Sprintf(`
			INSERT INTO user_interaction_log (profile_id, organization_id, count, timestamp, bucket_hour)
			VALUES %s
			ON CONFLICT (profile_id, organization_id, bucket_hour)
			DO UPDATE SET
				count = user_interaction_log.count + EXCLUDED.count,
				"timestamp" = EXCLUDED."timestamp"
		`, strings.Join(placeholders, ", "))
		if err := db.Exec(sql, args...).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) FlushAuditLogs(ctx context.Context, logs []*RequestAuditLog) error {
	if s == nil || s.db == nil || len(logs) == 0 {
		return nil
	}
	return s.db.WithContext(ctx).CreateInBatches(logs, 100).Error
}
