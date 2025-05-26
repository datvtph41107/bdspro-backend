package u_product_child

type InfoAreaResponse struct {
	Total    float64 `gorm:"column:total_area" json:"totalArea"`
	Avaiable float64 `gorm:"column:avaiable_area" json:"avaiableArea"`
	Area     float64 `gorm:"column:split_area" json:"splitArea"`
}
