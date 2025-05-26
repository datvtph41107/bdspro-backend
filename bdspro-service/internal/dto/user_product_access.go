package dto

type ProductAccessRequest struct {
	// ProductID  *uint64 `json:"parentId" binding:"required"`
	TargetID   *uint64 `json:"targetId" binding:"required"`
	TargetType uint    `json:"targetType"`
	Commission float32 `json:"commission"`
	Fields     string  `json:"fields"`
	// ProfileIDs []uint64 `json:"profileIds" binding:"required,notemptyarray"`
}

type ProductAccessBulkRequest struct {
	Datas   []ProductAccessRequest `json:"datas"`
	Deletes []ProductAccessRequest `json:"deletes"`
}

type ProductAccessUpdate struct {
	TargetID   uint64  `json:"targetId" binding:"required"`
	TargetType uint64  `json:"targetType" binding:"required"`
	ProductID  uint64  `json:"productId" binding:"required"`
	Commission float64 `json:"commission" binding:"required"`
}
