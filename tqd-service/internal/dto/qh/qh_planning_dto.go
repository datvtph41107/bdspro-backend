package qh_dto

import (
	_dto "common/domain/dto"
)

type ListPlanningProjectsRequest struct {
	_dto.Pagable
	Search         string `json:"search"`
	PlanningType   string `json:"planningType"`
	PlanningLevel  string `json:"planningLevel"`
	ValidityStatus string `json:"validityStatus"`
	ProcessStatus  string `json:"processStatus"`
}

type ListProAIJobsRequest struct {
	_dto.Pagable
	Search        string `json:"search"`
	JobType       string `json:"jobType"`
	ProcessStatus string `json:"processStatus"`
}

type ListPlanningDocumentsRequest struct {
	_dto.Pagable
	PlanningProjectID uint64 `json:"planningProjectId"`
	DocumentType      string `json:"documentType"`
	ValidityStatus    string `json:"validityStatus"`
	ProcessStatus     string `json:"processStatus"`
}

type ListPlanningEventsRequest struct {
	_dto.Pagable
	PlanningProjectID uint64 `json:"planningProjectId"`
	Ver               bool   `json:"ver"`
}

// FolderFileInput — 1 file trong folder đồ án admin upload.
type FolderFileInput struct {
	RelativePath string
	Content      []byte
	ContentType  string
}

type PlanningProjectSeoSource struct {
	PlanningProjectID uint64  `json:"planningProjectId" gorm:"column:planning_project_id"`
	Name              string  `json:"name" gorm:"column:name"`
	Slug              string  `json:"slug" gorm:"column:slug"`
	SeoID             *uint64 `json:"seoId" gorm:"column:seo_id"`
}

type AdminGenPlanningProjectSeoItem struct {
	PlanningProjectID uint64 `json:"planningProjectId"`
	SeoID             uint64 `json:"seoId"`
	CanonicalURL      string `json:"canonicalUrl"`
	Slug              string `json:"slug"`
	Created           bool   `json:"created"`
	Message           string `json:"message"`
}

type AdminGenPlanningProjectSeoResult struct {
	Items   []AdminGenPlanningProjectSeoItem `json:"items"`
	Total   int64                            `json:"total"`
	Created int64                            `json:"created"`
	Skipped int64                            `json:"skipped"`
	Failed  int64                            `json:"failed"`
	Errors  []string                         `json:"errors"`
}
