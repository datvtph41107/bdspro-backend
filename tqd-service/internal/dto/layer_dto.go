package dto

type LayerSeoSource struct {
	LayerID      uint64  `json:"layerId" gorm:"column:layer_id"`
	DisplayName  string  `json:"displayName" gorm:"column:display_name"`
	SeoID        *uint64 `json:"seoId" gorm:"column:seo_id"`
	ThumbnailURL string  `json:"thumbnailUrl" gorm:"column:thumbnail_url"`
	ImageURL     string  `json:"imageUrl" gorm:"column:image_url"`
	Avatar       string  `json:"avatar" gorm:"column:avatar"`
	ProvinceID   *uint64 `json:"provinceId" gorm:"column:province_id"`
}

type AdminGenLayerSeoItem struct {
	LayerID      uint64 `json:"layerId"`
	SeoID        uint64 `json:"seoId"`
	CanonicalURL string `json:"canonicalUrl"`
	Slug         string `json:"slug"`
	Created      bool   `json:"created"`
	Message      string `json:"message"`
}

type AdminGenLayerSeoResult struct {
	Items   []AdminGenLayerSeoItem `json:"items"`
	Total   int64                  `json:"total"`
	Created int64                  `json:"created"`
	Skipped int64                  `json:"skipped"`
	Failed  int64                  `json:"failed"`
	Errors  []string               `json:"errors"`
}
