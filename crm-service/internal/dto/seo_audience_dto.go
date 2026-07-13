package dto

import "time"

type CreateSeoAudienceRequest struct {
	Code string
	Name string
}

type UpdateSeoAudienceRequest struct {
	ID              uint64
	ExpectedVersion uint32
	Name            string
	Description     string
	DecisionRole    string
	KnowledgeLevel  string
	OwnerCode       string
	ReviewDueAt     *time.Time
}

type SeoAudienceCommandRequest struct {
	ID              uint64
	ExpectedVersion uint32
	Reason          string
}

type SeoAudienceListRequest struct {
	Keyword      string
	Status       string
	DecisionRole string
	Page         uint32
	Size         uint32
}
