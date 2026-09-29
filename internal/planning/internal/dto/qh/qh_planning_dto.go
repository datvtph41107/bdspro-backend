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
