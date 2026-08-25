package domain

import "time"

// KeywordSearchLog stores one flushed keyword-search count batch.
type KeywordSearchLog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Keyword   string    `gorm:"size:255;not null;index;column:keyword"`
	Source    string    `gorm:"size:64;not null;index;column:source"`
	Count     int64     `gorm:"not null;column:count"`
	Timestamp time.Time `gorm:"not null;index;column:timestamp"`
}

func (KeywordSearchLog) TableName() string {
	return "keyword_search_log"
}
