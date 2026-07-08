package dto

type SeoWorkerGenerateRequest struct {
	Limit       int
	TriggerType string
	Reason      string
}

type SeoWorkerGenerateResult struct {
	Scanned int
	Success int
	Failed  int
	Skipped int
	Errors  []string
}
