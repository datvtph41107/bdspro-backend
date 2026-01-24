package domain

type ProductStats struct {
	ProductID uint64 `gorm:"primaryKey;column:product_id;comment:ID sản phẩm"`

	TotalViews uint64 `gorm:"column:total_views;not null;default:0;comment:Tổng số lượt xem"`

	DurationViews uint32 `gorm:"column:duration_views;not null;default:0"`

	LastViewEventTime int64 `gorm:"column:last_view_event_time;not null;default:0;index"`
	LastUpdatedAt     int64 `gorm:"column:last_updated_at;not null;index"`

	NumOfInterested *uint64 `gorm:"column:num_of_interested;type:bigint"`
	ViewerIndex     *uint64 `gorm:"column:viewer_index;type:bigint"`

	PostCount  *uint64 `gorm:"column:post_count;type:bigint"`
	DealCount  *uint64 `gorm:"column:deal_count;type:bigint"`
	ShareCount *uint64 `gorm:"column:share_count;type:bigint"`
	AssetCount *uint64 `gorm:"column:asset_count;type:bigint"`
	BuildCount *uint64 `gorm:"column:build_count;type:bigint"`
}

func (ProductStats) TableName() string {
	return "product_stats"
}
