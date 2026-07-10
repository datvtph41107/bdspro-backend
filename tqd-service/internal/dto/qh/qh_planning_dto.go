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
}

type ListPlanningDocumentsRequest struct {
	_dto.Pagable
	PlanningProjectID uint64 `json:"planningProjectId"`
	DocumentType      string `json:"documentType"`
	ValidityStatus    string `json:"validityStatus"`
}

type ListPlanningEventsRequest struct {
	_dto.Pagable
	PlanningProjectID uint64 `json:"planningProjectId"`
}
