package dto

// Shared HTTP request bindings for admin / directory / internal SEO routes.

type PageQuery struct {
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

type IDPath struct {
	ID uint64 `path:"id"`
}

type LayerIDPath struct {
	LayerID uint64 `path:"layerId|layer_id"`
}

type BatchIDPath struct {
	BatchID string `path:"batchId"`
}

type ErrorIDPath struct {
	ErrorID uint64 `path:"errorId"`
}

type CompositeIDPath struct {
	CompositeID uint64 `path:"compositeId"`
}

type WardNeighborPath struct {
	WardID         uint64 `path:"ward_id"`
	NeighborWardID uint64 `path:"neighbor_ward_id"`
}

type WardIDPath struct {
	WardID uint64 `path:"ward_id"`
}

type ParcelIDPath struct {
	ParcelID uint64 `path:"parcelId"`
}

type PlanningIDPath struct {
	PlanningID uint64 `path:"planningId"`
}

type PlanningLayerPath struct {
	PlanningID uint64 `path:"planningId"`
	LayerID    uint64 `path:"layerId"`
}

type PlanningDocHistoryPath struct {
	PlanningDocumentID uint64 `path:"planningDocumentId"`
	ID                 uint64 `path:"id"`
}

type SubscriptionIDPath struct {
	SubscriptionID uint64 `path:"subscriptionId"`
}

type DayOfWeekPath struct {
	DayOfWeek int `path:"dayOfWeek"`
}

type CategoryIDPath struct {
	CategoryID uint64 `path:"categoryId"`
}

type CodePath struct {
	Code string `path:"code"`
}

type LevelPath struct {
	Level int32 `path:"level"`
}

type SeoIDBody struct {
	ID    uint64 `path:"id|layerId|parcelId|planningProjectId|planningDocumentId" json:"-"`
	SeoID uint64 `json:"seoId"`
}

type AdminGenSeoIDsRequest struct {
	IDs []uint64 `json:"ids"`
}

type AdminGenAdminUnitSeoRequest struct {
	CompositeIDs []uint64 `json:"compositeIds"`
}

type AdminSyncAreaRequest struct {
	CompositeIDs []uint64 `json:"compositeIds"`
}

type AdminRegionProcessingStatusRequest struct {
	RegionID uint64 `path:"regionId" json:"-"`
	Status   uint32 `json:"status"`
	Notes    string `json:"notes"`
}

type AdminSyncRegionRequest struct {
	LayerID uint64 `json:"layerId"`
	LabelID uint64 `json:"labelId"`
}

type AdminImportPreviewRequest struct {
	LayerID uint64 `json:"layerId"`
	GeoJSON string `json:"geoJson"`
}

type AdminListImportErrorsRequest struct {
	LayerID uint64 `path:"layerId" json:"-"`
	Page    uint32 `form:"page"`
	Size    uint32 `form:"size"`
}

type EmptyJSON struct{}

type UpdateRatingRequest struct {
	ID          uint64  `path:"id" json:"-"`
	Rating      float64 `json:"rating"`
	ReviewCount uint32  `json:"reviewCount"`
}

type AdminQaStatusRequest struct {
	ID       uint64 `path:"id" json:"-"`
	QaStatus uint32 `json:"qaStatus"`
	Note     string `json:"note"`
}

type AdminSeverityRequest struct {
	ID       uint64 `path:"id" json:"-"`
	Severity uint32 `json:"severity"`
	Reason   string `json:"reason"`
}

type AdminAssignReportRequest struct {
	ID         uint64 `path:"id" json:"-"`
	AssigneeID uint64 `json:"assigneeId"`
}

type AdminCloseReportRequest struct {
	ID     uint64 `path:"id" json:"-"`
	Reject bool   `json:"reject"`
	Note   string `json:"note"`
}

type AdminReportImagesRequest struct {
	ID     uint64   `path:"id" json:"-"`
	Images []string `json:"images"`
}

type AdminReportLinkageRequest struct {
	ID      uint64                 `path:"id" json:"-"`
	Linkage map[string]interface{} `json:"linkage"`
}

type ListParcelSeoSourcesRequest struct {
	Limit uint32 `form:"limit"`
}

type CommercialUsageQuery struct {
	Page      uint32 `form:"page"`
	PageSize  uint32 `form:"pageSize"`
	ProfileID uint64 `form:"profileId"`
	Operation string `form:"operation"`
	MeterCode string `form:"meterCode"`
}

type GeneratedReportsQuery struct {
	Page      uint32 `form:"page"`
	PageSize  uint32 `form:"pageSize"`
	ProfileID uint64 `form:"profileId"`
	Status    uint32 `form:"status"`
	JobStatus string `form:"jobStatus"`
}
