package types

import "time"

// ResolutionStatus — trạng thái kết quả phân giải
type ResolutionStatus string

const (
	StatusResolved             ResolutionStatus = "resolved"
	StatusResolvedWithWarning  ResolutionStatus = "resolved_with_warning"
	StatusResolvedWithConflict ResolutionStatus = "resolved_with_conflict"
	StatusFallbackResolved     ResolutionStatus = "fallback_resolved"
	StatusNoCandidate          ResolutionStatus = "no_candidate"
)

// ResolutionTrace — trace của toàn bộ pipeline
type ResolutionTrace struct {
	Steps       []TraceStep
	StartedAt   time.Time
	CompletedAt time.Time
}

// TraceStep — 1 bước trong pipeline
type TraceStep struct {
	StepName    string
	DurationMs  int64
	InputCount  int
	OutputCount int
	Status      string
	Details     map[string]any
}

// PipelineConfig — cấu hình cho Resolution Pipeline
type PipelineConfig struct {
	EnableTrace    bool
	EnableTieBreak bool
	EnableFallback bool
	TieBreakField  string // "overlap_percent", "legal_confidence", "priority"
}

// DefaultPipelineConfig — cấu hình pipeline mặc định
func DefaultPipelineConfig() *PipelineConfig {
	return &PipelineConfig{
		EnableTrace:    true,
		EnableTieBreak: true,
		EnableFallback: true,
		TieBreakField:  "overlap_percent",
	}
}
