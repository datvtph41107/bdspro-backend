package qh_dto

// ============================================================
// PARCEL DETAIL DTO — Data cho màn chi tiết thửa đất
// ============================================================
// LOGICAL:
// - Surface: Cần cấu trúc dữ liệu cho frontend
// - Root: Tách biệt domain khỏi API contract, enum thay vì boolean
// - Mechanism: DTO mapping từ domain → API response

// ParcelDetailResponseDTO response cho GET /parcels/:id
type ParcelDetailResponseDTO struct {
	Parcel             ParcelIdentityDTO         `json:"parcel"`
	CurrentUse         CurrentUseDTO             `json:"currentUse"`
	PlanningConclusion *PlanningConclusionDTO    `json:"planningConclusion,omitempty"` // ← pointer
	PrimaryUse         *PrimaryUseDTO            `json:"primaryUse,omitempty"`         // ← pointer
	Guidance           *GuidanceDTO              `json:"guidance,omitempty"`           // ← pointer
	LegalDocuments     *LegalDocumentsSummaryDTO `json:"legalDocuments,omitempty"`     // ← pointer
	Meta               *ResponseMetaDTO          `json:"_meta,omitempty"`
	Links              *LinksDTO                 `json:"_links,omitempty"`
}

// ParcelIdentityDTO — Định danh thửa đất
// LOGICAL: Identity domain — xác định đang xem thửa nào
type ParcelIdentityDTO struct {
	ID               string        `json:"id"`
	MapSheetNumber   string        `json:"mapSheetNumber"`
	LandParcelNumber string        `json:"landParcelNumber"`
	Address          string        `json:"address"`
	TotalAreaSqm     float64       `json:"totalAreaSqm"`
	Centroid         *CentroidDTO  `json:"centroid"`
	Ownership        *OwnershipDTO `json:"ownership"`
	Location         *LocationDTO  `json:"location"`
}

type CentroidDTO struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type OwnershipDTO struct {
	CertificateNumber string `json:"certificateNumber"`
	IssuedAt          string `json:"issuedAt"`
	IssuedBy          string `json:"issuedBy"`
	HolderName        string `json:"holderName"`
	LegalDocumentId   string `json:"legalDocumentId"`
}

type LocationDTO struct {
	Province     string `json:"province"`
	ProvinceCode string `json:"provinceCode"`
	WardCode     string `json:"wardCode"`
}

// CurrentUseDTO — Thực trạng sử dụng đất
// LOGICAL: Physical domain — đất hiện tại là gì?
type CurrentUseDTO struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	Color           string `json:"color"`
	BuildStatus     string `json:"buildStatus"`
	LegalDocumentId string `json:"legalDocumentId"`
}

// PlanningConclusionDTO — Kết luận quy hoạch
// LOGICAL: Planning domain — kết luận cuối cùng
type PlanningConclusionDTO struct {
	BuildStatus string                 `json:"buildStatus"`
	HasConflict bool                   `json:"hasConflict"`
	Alert       *AlertDTO              `json:"alert"`
	LegalBasis  *LegalBasisDTO         `json:"legalBasis"`
	Breakdown   []PlanningBreakdownDTO `json:"breakdown"`
	Statistics  *PlanningStatisticsDTO `json:"statistics"`
}

type AlertDTO struct {
	Level   string `json:"level"`
	Score   uint32 `json:"score"`
	Message string `json:"message"`
}

type LegalBasisDTO struct {
	Summary          string   `json:"summary"`
	LegalDocumentIds []string `json:"legalDocumentIds"`
}

// PlanningBreakdownDTO — Chi tiết từng loại đất
// LOGICAL: Breakdown domain — phân bổ diện tích theo loại đất
type PlanningBreakdownDTO struct {
	Rank             int      `json:"rank"`
	LandUseCode      string   `json:"landUseCode"`
	LandUseName      string   `json:"landUseName"`
	LandUseColor     string   `json:"landUseColor"`
	GroupCode        string   `json:"groupCode"`
	GroupName        string   `json:"groupName"`
	AreaSqm          float64  `json:"areaSqm"`
	Percent          float64  `json:"percent"`
	BuildStatus      string   `json:"buildStatus"`
	AlertLevel       string   `json:"alertLevel"`
	LegalDocumentIds []string `json:"legalDocumentIds"`
}

type PlanningStatisticsDTO struct {
	TotalLayersAffected int     `json:"totalLayersAffected"`
	TotalZonesAffected  int     `json:"totalZonesAffected"`
	TotalOverlapAreaSqm float64 `json:"totalOverlapAreaSqm"`
	TotalOverlapPercent float64 `json:"totalOverlapPercent"`
}

// PrimaryUseDTO — Quy hoạch chính
// LOGICAL: Primary domain — loại đất chiếm ưu thế nhất
type PrimaryUseDTO struct {
	Rule        string  `json:"rule"`
	LandUseCode string  `json:"landUseCode"`
	LandUseName string  `json:"landUseName"`
	Percent     float64 `json:"percent"`
	BuildStatus string  `json:"buildStatus"`
}

// GuidanceDTO — Hướng dẫn hành động
// LOGICAL: Guidance domain — user cần làm gì tiếp theo?
type GuidanceDTO struct {
	Warnings        []WarningDTO        `json:"warnings"`
	Recommendations []RecommendationDTO `json:"recommendations"`
}

type WarningDTO struct {
	AlertLevel      string         `json:"alertLevel"`
	LandUseCode     string         `json:"landUseCode"`
	AffectedAreaSqm float64        `json:"affectedAreaSqm"`
	AffectedPercent float64        `json:"affectedPercent"`
	Title           string         `json:"title"`
	Action          string         `json:"action"`
	LegalBasis      *LegalBasisDTO `json:"legalBasis"`
}

type RecommendationDTO struct {
	Type             string   `json:"type"`
	Message          string   `json:"message"`
	LegalDocumentIds []string `json:"legalDocumentIds,omitempty"`
}
