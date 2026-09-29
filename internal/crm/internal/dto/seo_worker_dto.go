package dto

// SeoWorkerGenerateRequest cấu hình một lượt tạo HTML theo lô.
type SeoWorkerGenerateRequest struct {
	Limit       int
	TriggerType string
	Reason      string
}

// SeoWorkerGenerateResult tổng hợp kết quả của một lượt tạo HTML.
type SeoWorkerGenerateResult struct {
	Scanned int
	Success int
	Failed  int
	Skipped int
	Errors  []string
}
