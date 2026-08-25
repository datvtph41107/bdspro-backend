package postgres

import (
	"context"
	"time"

	_db "common/db"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/repo"
)

type KeywordSearchPostgres struct {
	db *_db.TransactionRepo
}

// @bind: internal/repo/IKeywordSearchRepo
func NewKeywordSearchRepo(db *_db.TransactionRepo) repo.IKeywordSearchRepo {
	return &KeywordSearchPostgres{db: db}
}

func (r *KeywordSearchPostgres) FlushCounts(ctx context.Context, logs []*domain.KeywordSearchLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.GetDB(ctx).WithContext(ctx).CreateInBatches(logs, 100).Error
}

func (r *KeywordSearchPostgres) CountFromTime(ctx context.Context, fromTime time.Time, source string) (uint64, error) {
	var total int64
	query := `
		SELECT COALESCE(SUM(count), 0)
		FROM keyword_search_log
		WHERE "timestamp" >= ?
	`
	args := []interface{}{fromTime}
	if source != "" {
		query += ` AND source = ?`
		args = append(args, source)
	}
	err := r.db.GetDB(ctx).WithContext(ctx).Raw(query, args...).Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return uint64(total), nil
}

func (r *KeywordSearchPostgres) GetTopKeywords(ctx context.Context, fromTime time.Time, source string, limit int) ([]dto.KeywordSearchItem, error) {
	type row struct {
		Keyword string `gorm:"column:keyword"`
		Source  string `gorm:"column:source"`
		Count   uint64 `gorm:"column:count"`
	}

	query := `
		SELECT
			keyword,
			source,
			SUM(count)::bigint AS count
		FROM keyword_search_log
		WHERE "timestamp" >= ?
	`
	args := []interface{}{fromTime}
	if source != "" {
		query += ` AND source = ?`
		args = append(args, source)
	}
	query += `
		GROUP BY keyword, source
		ORDER BY count DESC, keyword ASC
		LIMIT ?
	`
	args = append(args, limit)

	var rows []row
	err := r.db.GetDB(ctx).WithContext(ctx).Raw(query, args...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.KeywordSearchItem, 0, len(rows))
	for _, item := range rows {
		result = append(result, dto.KeywordSearchItem{
			Keyword: item.Keyword,
			Source:  item.Source,
			Count:   item.Count,
		})
	}
	return result, nil
}
