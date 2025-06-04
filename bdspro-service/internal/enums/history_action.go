package enums

type EHistoryAction string

const (
	EActionProductCreate EHistoryAction = "product_create"
	EActionProductUpdate EHistoryAction = "product_update"
	EActionProductDelete EHistoryAction = "product_delete"
	EActionAssetCreate   EHistoryAction = "asset_create"
	EActionAssetUpdate   EHistoryAction = "asset_update"
	EActionAssetDelete   EHistoryAction = "asset_delete"
	EActionAssetMerge    EHistoryAction = "asset_merge"
)

var HistoryActionMap = map[EHistoryAction]string{
	EActionProductCreate: "Tạo sản phẩm",
	EActionProductUpdate: "Cập nhật sản phẩm",
	EActionProductDelete: "Xóa sản phẩm",
	EActionAssetCreate:   "Tạo tài sản",
	EActionAssetUpdate:   "Cập nhật tài sản",
	EActionAssetDelete:   "Xóa tài sản",
	EActionAssetMerge:    "Gộp tài sản",
}
