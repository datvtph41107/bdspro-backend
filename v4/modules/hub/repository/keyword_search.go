package repository

import (
	"context"
	"time"

	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"

	"gorm.io/gorm"
)

// KeywordSearchRepository manages keyword_search_log persistence.
type KeywordSearchRepository interface {
	FlushCounts(ctx context.Context, logs []*entity.KeywordSearchLog) error
	CountFromTime(ctx context.Context, fromTime time.Time, source string) (uint64, error)
	GetTopKeywords(ctx context.Context, fromTime time.Time, source string, limit int) ([]dto.KeywordSearchItem, error)
}

type keywordSearchPostgres struct{ db *gorm.DB }

func NewKeywordSearchPostgres(db *gorm.DB) KeywordSearchRepository {
	return &keywordSearchPostgres{db: db}
}

func (r *keywordSearchPostgres) FlushCounts(ctx context.Context, logs []*entity.KeywordSearchLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(logs, 100).Error
}

func (r *keywordSearchPostgres) CountFromTime(ctx context.Context, fromTime time.Time, source string) (uint64, error) {
	var total int64
	query := `SELECT COALESCE(SUM(count), 0) FROM keyword_search_log WHERE "timestamp" >= ?`
	args := []interface{}{fromTime}
	if source != "" {
		query += ` AND source = ?`
		args = append(args, source)
	}
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&total).Error
	return uint64(total), err
}

func (r *keywordSearchPostgres) GetTopKeywords(ctx context.Context, fromTime time.Time, source string, limit int) ([]dto.KeywordSearchItem, error) {
	type row struct {
		Keyword string `gorm:"column:keyword"`
		Source  string `gorm:"column:source"`
		Count   uint64 `gorm:"column:count"`
	}
	query := `SELECT keyword, source, SUM(count)::bigint AS count FROM keyword_search_log WHERE "timestamp" >= ?`
	args := []interface{}{fromTime}
	if source != "" {
		query += ` AND source = ?`
		args = append(args, source)
	}
	query += ` GROUP BY keyword, source ORDER BY count DESC, keyword ASC LIMIT ?`
	args = append(args, limit)
	var rows []row
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]dto.KeywordSearchItem, 0, len(rows))
	for _, r := range rows {
		result = append(result, dto.KeywordSearchItem{Keyword: r.Keyword, Source: r.Source, Count: r.Count})
	}
	return result, nil
}
