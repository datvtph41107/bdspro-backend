package dto

type IdentifierDTO struct {
	ID             string  `json:"id"`
	LandParcelCode string  `json:"land_parcel_code"`
	ProjectCode    string  `json:"project_code"`
	ProvinceID     string  `json:"province_id"`
	DistrictID     string  `json:"district_id"`
	WardID         string  `json:"ward_id"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Type           string  `json:"type"`
	LegalStatus    string  `json:"legal_status"`
	CurrentOwnerID string  `json:"current_owner_id"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	Version        int     `json:"version"`
}

type SearchCriteriaDTO struct {
	ProvinceID  string  `json:"province_id"`
	DistrictID  string  `json:"district_id"`
	WardID      string  `json:"ward_id"`
	Type        string  `json:"type"`
	LegalStatus string  `json:"legal_status"`
	OwnerID     string  `json:"owner_id"`
	Keyword     string  `json:"keyword"`
	CreatedFrom string  `json:"created_from"`
	CreatedTo   string  `json:"created_to"`
	RadiusKm    float64 `json:"radius_km"`
	CenterLat   float64 `json:"center_lat"`
	CenterLng   float64 `json:"center_lng"`
}

type PaginationDTO struct {
	Page  int32 `json:"page"`
	Limit int32 `json:"limit"`
}

type SearchResultDTO struct {
	Identifiers []*IdentifierDTO `json:"identifiers"`
	Total       int64            `json:"total"`
	Page        int32            `json:"page"`
	Limit       int32            `json:"limit"`
}
