package dto

import tqdpb "pb/types/tqd"

// --- map-workspace ---

type WorkspacePageHTTPRequest struct {
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Cursor    string `form:"cursor"`
	Timestamp string `form:"timestamp"`
}

type FollowParcelHTTPRequest struct {
	ParcelID uint64 `json:"parcelId"`
	Note     string `json:"note"`
}

type RemoveFollowedByIDHTTPRequest struct {
	FollowID uint64 `path:"followId"`
}

type RemoveFollowedByParcelHTTPRequest struct {
	ParcelID uint64 `path:"parcelId"`
}

type AddViewHistoryHTTPRequest struct {
	EntityType uint32  `json:"entityType"`
	EntityID   uint64  `json:"entityId"`
	ParcelID   uint64  `json:"parcelId"`
	RegionID   uint64  `json:"regionId"`
	Source     uint32  `json:"source"`
	Zoom       float64 `json:"zoom"`
}

type TrackViewHistoryHTTPRequest struct {
	ClientEventID  string           `json:"clientEventId"`
	SessionID      string           `json:"sessionId"`
	DeviceID       string           `json:"deviceId"`
	EntityType     uint32           `json:"entityType"`
	EntityID       uint64           `json:"entityId"`
	ParcelID       uint64           `json:"parcelId"`
	RegionID       uint64           `json:"regionId"`
	Source         uint32           `json:"source"`
	SourceRef      string           `json:"sourceRef"`
	RouteName      string           `json:"routeName"`
	Zoom           float64          `json:"zoom"`
	Center         SpatialPointDTO  `json:"center"`
	ViewportBounds SpatialBoundsDTO `json:"viewportBounds"`
	VisibleMs      uint64           `json:"visibleMs"`
	CountIntent    bool             `json:"countIntent"`
	MetadataJSON   string           `json:"metadataJson"`
}

type RemoveViewHistoryHTTPRequest struct {
	HistoryID uint64 `path:"historyId"`
}

type ClearViewHistoryHTTPRequest struct{}

type WorkspaceReportIDHTTPRequest struct {
	ReportID uint64 `path:"reportId"`
}

type WorkspaceReportActionHTTPRequest struct {
	Seg string `path:"seg"`
}

type CreateGeneratedReportHTTPRequest struct {
	ReportType  uint32                   `json:"reportType"`
	EntityType  uint32                   `json:"entityType"`
	EntityID    uint64                   `json:"entityId"`
	ParcelID    uint64                   `json:"parcelId"`
	RegionID    uint64                   `json:"regionId"`
	Title       string                   `json:"title"`
	Subtitle    string                   `json:"subtitle"`
	Location    *WorkspaceLocationDTO    `json:"location"`
	Spatial     *ReportSpatialPreviewDTO `json:"spatial"`
	Comparison  *ReportComparisonDTO     `json:"comparison"`
	Format      string                   `json:"format"`
	Idempotency string                   `json:"idempotencyKey"`
	Metadata    map[string]any           `json:"metadata"`
}

// --- discovery ---

type DiscoveryIdentifyHTTPRequest struct {
	Point   *DiscoveryPointHTTP   `json:"point"`
	Context *DiscoveryContextHTTP `json:"context"`
	Options *DiscoveryOptionsHTTP `json:"options"`
}

type DiscoveryPointHTTP struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DiscoveryContextHTTP struct {
	Zoom           float64           `json:"zoom"`
	MapMode        string            `json:"mapMode"`
	ActiveLayerIDs []string          `json:"activeLayerIds"`
	PreferredKinds []string          `json:"preferredKinds"`
	Viewport       *SpatialBoundsDTO `json:"viewport"`
}

type DiscoveryOptionsHTTP struct {
	ToleranceMeters        float64 `json:"toleranceMeters"`
	Limit                  int32   `json:"limit"`
	IncludeGeometryPreview bool    `json:"includeGeometryPreview"`
}

type DiscoverySearchHTTPRequest struct {
	Q      string  `form:"q"`
	Text   string  `form:"text"`
	Lat    float64 `form:"lat"`
	Lng    float64 `form:"lng"`
	Limit  int32   `form:"limit"`
	Cursor string  `form:"cursor"`
	Kinds  string  `form:"kinds"`
}

type DiscoveryEntityKeyHTTPRequest struct {
	Key string `path:"key"`
}

// --- parcels / map-target / layers ---

type ParcelBatchHTTPRequest struct {
	Data []map[string]any `json:"data"`
}

type FindByLocationHTTPRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

type ResolveMapTargetByLocationHTTPRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
}

type ResolveMapTargetByPoiHTTPRequest struct {
	PoiID uint64 `path:"poi_id"`
}

type FindByPolygonHTTPRequest struct {
	Coordinates     []PolygonPointHTTP `json:"coordinates"`
	Intersect       bool               `json:"intersect"`
	MaxAreaHa       float64            `json:"maxAreaHa"`
	Page            uint32             `json:"page"`
	Size            uint32             `json:"size"`
	Z               *uint32            `json:"z"`
	LayerIDs        []uint64           `json:"layerIds"`
	LabelIDs        []uint64           `json:"labelIds"`
	LandUseIDs      []uint64           `json:"landUseIds"`
	LandUseCodes    []string           `json:"landUseCodes"`
	Keyword         string             `json:"keyword"`
	CanBuild        *bool              `json:"canBuild"`
	IncludeAnalysis bool               `json:"includeAnalysis"`
	Timestamp       int64              `json:"timestamp"`
}

type PolygonPointHTTP struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type ParcelIDHTTPRequest struct {
	ParcelID  uint64  `path:"parcelId"`
	Timestamp int64   `form:"timestamp"`
	Page      uint32  `form:"page"`
	Size      uint32  `form:"size"`
	Limit     int32   `form:"limit"`
	DocType   *uint32 `form:"docType"`
	Status    *uint32 `form:"status"`
	Mode      string  `form:"mode"`
	View      string  `form:"view"`
}

type ParcelPreviewIDHTTPRequest struct {
	ID uint64 `path:"id|parcelId"`
}

type ZoneGeometryHTTPRequest struct {
	ParcelID uint64 `path:"parcelId"`
	LayerID  uint64 `path:"layerId"`
	ZoneID   uint64 `path:"zoneId"`
}

type LayerHistoryHTTPRequest struct {
	LayerID uint64 `path:"layerId"`
}

type LayerLegalDocsHTTPRequest struct {
	LayerID   uint64             `path:"layerId" json:"-"`
	LegalDocs []*LayerLegalInput `json:"legalDocs"`
}

type LayerLegalDocsListHTTPRequest struct {
	LayerID uint64 `path:"layerId"`
}

type LayerLegalDocIDHTTPRequest struct {
	ID uint64 `path:"id"`
}

// --- public / qh public ---

type PublicMapsHTTPRequest struct{}

type CreateDefaultSeoHTTPRequest struct {
	RefType uint32 `json:"refType"`
	RefID   uint64 `json:"refId"`
}

type PublicProjectionValueHTTPRequest struct {
	Value string `path:"value"`
}

type PublicParcelQuickViewHTTPRequest struct {
	ParcelID uint64 `path:"parcel_id"`
}

type PublicPlanningIdentityHTTPRequest struct {
	Identity string `path:"identity"`
}

type PublicRelatedEntityHTTPRequest struct {
	Key string `path:"key"`
}

type PublicPoisHTTPRequest struct {
	Page     uint32  `form:"page"`
	Size     uint32  `form:"size"`
	Category *uint64 `form:"categoryId"`
	Lat      float64 `form:"lat"`
	Lng      float64 `form:"lng"`
	Radius   float64 `form:"radius"`
	Search   string  `form:"search"`
	Mode     string  `form:"mode"`
}

type PublicParcelsSearchHTTPRequest struct {
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

type PublicClientSearchHTTPRequest struct {
	Text string `form:"text"`
	Q    string `form:"q"`
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

type QhPublicConfigHTTPRequest struct {
	Config string `path:"config"`
}

type QhPublicFontHTTPRequest struct {
	Fontstack string `path:"fontstack"`
	Range     string `path:"range"`
}

type QhPublicLegendsAllHTTPRequest struct{}

type QhPublicParcelQuickLayersHTTPRequest struct {
	ParcelID uint64 `path:"parcelId"`
}

// --- subscriptions / notifications ---

type SubscriptionListHTTPRequest struct {
	Page      uint32  `form:"page"`
	Size      uint32  `form:"size"`
	Status    *uint32 `form:"status"`
	Timestamp int64   `form:"timestamp"`
}

type SubscriptionParcelIDHTTPRequest struct {
	ParcelID uint64 `path:"parcelId"`
}

type CreateParcelSubscriptionHTTPRequest struct {
	ParcelID      uint64                          `path:"parcelId" json:"-"`
	Scope         *tqdpb.SubscriptionScopeMessage `json:"scope"`
	TriggerConfig string                          `json:"triggerConfig"`
}

type CreateRegionSubscriptionHTTPRequest struct {
	RegionID      uint64                          `json:"regionId"`
	Scope         *tqdpb.SubscriptionScopeMessage `json:"scope"`
	TriggerConfig string                          `json:"triggerConfig"`
}

type SubscriptionIDHTTPRequest struct {
	SubscriptionID uint64 `path:"subscriptionId"`
}

type NotificationListHTTPRequest struct {
	Page uint32  `form:"page"`
	Size uint32  `form:"size"`
	Type *uint32 `form:"type"`
	Read *bool   `form:"read"`
}

type NotificationIDHTTPRequest struct {
	NotificationID uint64 `path:"notificationId"`
}

type MarkAllNotificationsReadHTTPRequest struct{}

type CommercialUsageHTTPRequest struct{}

// --- features / reports / properties ---

type FeaturesNearbyHTTPRequest struct {
	Lat    float64 `form:"lat"`
	Lng    float64 `form:"lng"`
	Radius int     `form:"radius"`
	Limit  int     `form:"limit"`
	Offset int     `form:"offset"`
}

type FeaturesPolygonHTTPRequest struct {
	Polygon map[string]any `json:"polygon"`
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
}

type CreateReportHTTPRequest struct {
	ProblemReport uint32  `json:"problemReport"`
	ReportType    uint32  `json:"reportType"`
	Description   string  `json:"description"`
	TargetID      *uint64 `json:"targetId"`
}

type ReportIDHTTPRequest struct {
	ReportID uint64 `path:"reportId"`
}

type ListReportsHTTPRequest struct {
	Page       uint32  `form:"page"`
	Size       uint32  `form:"size"`
	ReportType *uint32 `form:"reportType"`
	Status     *uint32 `form:"status"`
}

type PropertyIDHTTPRequest struct {
	ID uint64 `path:"id"`
}
