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
	ClientViewOnly bool   `json:"clientViewOnly"`
	// AdmID is the administrative unit composite ID ({sourceId}{1=tỉnh|2=phường})
	AdmID uint64 `json:"admId,omitempty"`
	// LayerIDs lọc đồ án có tài liệu gắn layer qua layer_plan_doc_relation.
	LayerIDs []uint64 `json:"layerIds,omitempty"`
}

type ListPlanningProjectLayersRequest struct {
	_dto.Pagable
	PlanningID uint64 `json:"planningId"`
}

type ListProAIJobsRequest struct {
	_dto.Pagable
	Search        string `json:"search"`
	JobType       string `json:"jobType"`
	ProcessStatus string `json:"processStatus"`
}

type ListPlanningDocumentsRequest struct {
	_dto.Pagable
	PlanningProjectID uint64   `json:"planningProjectId"`
	DocumentType      string   `json:"documentType"`
	ValidityStatus    string   `json:"validityStatus"`
	ProcessStatus     string   `json:"processStatus"`
	Keyword           string   `json:"keyword"`
	LayerID           uint64   `json:"layerId"`
	LayerIDs          []uint64 `json:"layerIds"`
	ParcelID          uint64   `json:"parcelId"`
}

type AssignPlanningDocumentsToLayerRequest struct {
	LayerID  uint64   `json:"layerId"`
	LayerIDs []uint64 `json:"layerIds"`
	DocIDs   []uint64 `json:"docIds"`
}

type UnassignPlanningProjectLayerRequest struct {
	PlanningID uint64 `json:"planningId"`
	LayerID    uint64 `json:"layerId"`
}

type ListPlanningEventsRequest struct {
	_dto.Pagable
	PlanningProjectID uint64 `json:"planningProjectId"`
	Ver               bool   `json:"ver"`
}

type ListPlanningDocumentHistoriesRequest struct {
	_dto.Pagable
	PlanningDocumentID uint64 `json:"planningDocumentId"`
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
	ProvinceID        *uint64 `json:"provinceId" gorm:"column:province_id"`
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

type PlanningDocumentSeoSource struct {
	PlanningDocumentID uint64  `json:"planningDocumentId" gorm:"column:planning_document_id"`
	Title              string  `json:"title" gorm:"column:title"`
	Slug               string  `json:"slug" gorm:"column:slug"`
	SeoID              *uint64 `json:"seoId" gorm:"column:seo_id"`
	Filepath           string  `json:"filepath" gorm:"column:filepath"`
	Thumbnail          string  `json:"thumbnail" gorm:"column:thumbnail"`
	Ext                string  `json:"ext" gorm:"column:ext"`
	ProvinceID         *uint64 `json:"provinceId" gorm:"column:province_id"`
}

type AdminGenPlanningDocumentSeoItem struct {
	PlanningDocumentID uint64 `json:"planningDocumentId"`
	SeoID              uint64 `json:"seoId"`
	CanonicalURL       string `json:"canonicalUrl"`
	Slug               string `json:"slug"`
	Created            bool   `json:"created"`
	Message            string `json:"message"`
}

type AdminGenPlanningDocumentSeoResult struct {
	Items   []AdminGenPlanningDocumentSeoItem `json:"items"`
	Total   int64                             `json:"total"`
	Created int64                             `json:"created"`
	Skipped int64                             `json:"skipped"`
	Failed  int64                             `json:"failed"`
	Errors  []string                          `json:"errors"`
}

type PlanningProjectCount struct {
	Timeline int64 `json:"timeline"`
	Tin      int64 `json:"tin"`
	Relation int64 `json:"relation"`
	Doc      int64 `json:"doc"`
	Layer    int64 `json:"layer"`
}
