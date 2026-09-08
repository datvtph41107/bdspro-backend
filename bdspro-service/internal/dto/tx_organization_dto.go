package dto

type TxAddProductToDealRequest struct {
	DealID    uint64 `json:"deal_id"`
	ProductID uint64 `json:"product_id"`
}
