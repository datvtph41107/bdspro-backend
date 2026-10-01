package entity

import "time"

// SearchLatencyLog stores one flushed latency batch for a text-search API path.
type SearchLatencyLog struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	Method        string    `gorm:"size:16;not null;index;column:method"`
	Path          string    `gorm:"size:512;not null;index;column:path"`
	Service       string    `gorm:"size:64;not null;index;column:service"`
	Count         int64     `gorm:"not null;column:count"`
	SumDurationMs int64     `gorm:"not null;column:sum_duration_ms"`
	MinDurationMs int64     `gorm:"not null;column:min_duration_ms"`
	MaxDurationMs int64     `gorm:"not null;column:max_duration_ms"`
	P50DurationMs int64     `gorm:"not null;default:0;column:p50_duration_ms"`
	P95DurationMs int64     `gorm:"not null;default:0;column:p95_duration_ms"`
	P99DurationMs int64     `gorm:"not null;default:0;column:p99_duration_ms"`
	Timestamp     time.Time `gorm:"not null;index;column:timestamp"`
}

func (SearchLatencyLog) TableName() string { return "search_latency_log" }
