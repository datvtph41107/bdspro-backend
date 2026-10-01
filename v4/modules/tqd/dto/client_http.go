package dto

// Client HTTP request shapes for /v2/tqd/client routes (form/path tags for httpapi).

// --- Administrative units ---

type ClientListAdministrativeUnitsRequest struct {
	Page              uint32 `form:"page"`
	Size              uint32 `form:"size"`
	Sort              string `form:"sort"`
	Search            string `form:"search"`
	PlanningProjectID uint64 `form:"planningProjectId"`
	CompositeIDs      string `form:"compositeIds"`
}

type ClientCompositeIDRequest struct {
	CompositeID uint64 `path:"compositeId"`
}

type ClientFollowedAdministrativeUnitsRequest struct {
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

// --- Layers ---

type ClientListLayersRequest struct {
	Page              uint32  `form:"page"`
	Size              uint32  `form:"size"`
	Search            *string `form:"search"`
	Type              *uint32 `form:"type"`
	DocID             uint64  `form:"docId"`
	PlanningProjectID uint64  `form:"planningProjectId"`
	Timestamp         int64   `form:"timestamp"`
}

type ClientGetLayerRequest struct {
	ID        uint64 `path:"id"`
	Timestamp int64  `form:"timestamp"`
}

type ClientBuildLayerTileRequest struct {
	LayerID   uint64  `path:"layerId" json:"-"`
	MinZoom   uint32  `json:"minZoom"`
	MaxZoom   uint32  `json:"maxZoom"`
	OutputDir *string `json:"outputDir"`
}

// --- Regions ---

type ClientListRegionsRequest struct {
	LayerID   *uint64  `form:"layerId"`
	MinLng    *float64 `form:"minLng"`
	MinLat    *float64 `form:"minLat"`
	MaxLng    *float64 `form:"maxLng"`
	MaxLat    *float64 `form:"maxLat"`
	Page      uint32   `form:"page"`
	Size      uint32   `form:"size"`
	Timestamp int64    `form:"timestamp"`
}

type ClientFindRegionByPointRequest struct {
	Lat       float64 `form:"lat"`
	Lng       float64 `form:"lng"`
	LayerID   *uint64 `form:"layerId"`
	Timestamp int64   `form:"timestamp"`
}

type ClientFindRegionsByBBoxRequest struct {
	MinLng    float64 `form:"minLng"`
	MinLat    float64 `form:"minLat"`
	MaxLng    float64 `form:"maxLng"`
	MaxLat    float64 `form:"maxLat"`
	LayerID   *uint64 `form:"layerId"`
	LabelID   *uint64 `form:"labelId"`
	Page      int     `form:"page"`
	Size      int     `form:"size"`
	Timestamp int64   `form:"timestamp"`
}

// --- Parcels ---

type ClientSearchParcelsRequest struct {
	Text string `form:"text"`
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

// --- Families ---

type ClientListFamiliesRequest struct {
	Page      uint32  `form:"page"`
	Size      uint32  `form:"size"`
	Search    *string `form:"search"`
	Timestamp int64   `form:"timestamp"`
}

type ClientBuildFamilyTileRequest struct {
	FamilyID  uint64  `path:"familyId" json:"-"`
	MinZoom   uint32  `json:"minZoom"`
	MaxZoom   uint32  `json:"maxZoom"`
	OutputDir *string `json:"outputDir"`
}

// --- Labels ---

type ClientListLabelsByLayerRequest struct {
	LayerID   uint64 `path:"layerId"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Timestamp int64  `form:"timestamp"`
}

type ClientGetLabelRequest struct {
	ID        uint64 `path:"id"`
	Timestamp int64  `form:"timestamp"`
}

// --- Land use ---

type ClientListLandUseRequest struct {
	LayerID uint64 `path:"layerId"`
	Page    uint32 `form:"page"`
	Size    uint32 `form:"size"`
}

type ClientCreateLandUseRequest struct {
	LayerID   uint64 `path:"layerId" json:"-"`
	LandUseID uint64 `json:"landUseId"`
}

// --- Legends ---

type ClientListLegendsRequest struct {
	LayerID uint64 `path:"layerId"`
	Page    uint32 `form:"page"`
	Size    uint32 `form:"size"`
}

// --- Planning client ---

type ClientListPlanningProjectsRequest struct {
	Page           uint32 `form:"page"`
	Size           int32  `form:"size"`
	Search         string `form:"search"`
	Keyword        string `form:"keyword"`
	PlanningType   string `form:"planningType"`
	PlanningLevel  string `form:"planningLevel"`
	ValidityStatus string `form:"validityStatus"`
	ProcessStatus  string `form:"processStatus"`
	JurisdictionID uint64 `form:"jurisdictionId"`
	Area           string `form:"area"`
	LegalStatus    string `form:"legalStatus"`
	HasLayer       *int32 `form:"hasLayer"`
	SortBy         string `form:"sortBy"`
	SortOrder      string `form:"sortOrder"`
	Decision       string `form:"decision"`
	UpdatedFrom    string `form:"updatedFrom"`
	UpdatedTo      string `form:"updatedTo"`
	TagIDs         string `form:"tagIds"`
	ProvinceIDs    string `form:"provinceIds"`
	AdmID          uint64 `form:"admId"`
	ParcelID       uint64 `form:"parcelId"`
	LayerIDs       string `form:"layerIds"`
}

type ClientPlanningProjectIDRequest struct {
	ProjectID uint64 `path:"project_id"`
	Page      uint32 `form:"page"`
	Size      int32  `form:"size"`
	Keyword   string `form:"keyword"`
	DocType   string `form:"documentType"`
	Status    string `form:"status"`
}

type ClientFollowPlanningProjectRequest struct {
	ProjectID uint64 `path:"project_id" json:"-"`
	Note      string `json:"note"`
}

type ClientListFollowedPlanningProjectsRequest struct {
	Page   uint32 `form:"page"`
	Size   int32  `form:"size"`
	Search string `form:"search"`
}

type ClientListParcelPlanningRelationRequest struct {
	ParcelID uint64 `form:"parcelId"`
	Page     uint32 `form:"page"`
	Size     int32  `form:"size"`
}

type ClientPlanningDocumentIDRequest struct {
	DocumentID uint64 `path:"document_id"`
}

type ClientPlanningDocumentHistoryRequest struct {
	PlanningDocumentID uint64 `path:"planningDocumentId"`
	Page               uint32 `form:"page"`
	Size               uint32 `form:"size"`
}

type ClientLayerPlanDocRelationRequest struct {
	ParcelID uint64 `path:"parcel_id"`
	Page     uint32 `form:"page"`
	Size     int32  `form:"size"`
	DocType  string `form:"documentType"`
}

// ClientPlanningDocumentsTwoSegRequest disambiguates ServeMux-conflicting
// planning-documents/{relation-layer|view|{id}}/{parcel_id|document_id|histories} paths.
type ClientPlanningDocumentsTwoSegRequest struct {
	Seg1    string `path:"seg1"`
	Seg2    string `path:"seg2"`
	Page    uint32 `form:"page"`
	Size    int32  `form:"size"`
	DocType string `form:"documentType"`
}

type ClientPlanningProjectProjectionRequest struct {
	Value string `path:"value"`
}

type ClientNearbyAreasRequest struct {
	Lat    float64 `form:"lat"`
	Lng    float64 `form:"lng"`
	Radius float64 `form:"radius"`
}

type ClientPlanningDocumentsByAdmRequest struct {
	AdmID uint64 `form:"admId"`
	Page  uint32 `form:"page"`
	Size  int32  `form:"size"`
}
