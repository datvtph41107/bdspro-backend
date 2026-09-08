package enums

type EHistoryAction uint32

const (
	EActionProductCreate EHistoryAction = 100 //"product_create"
	EActionProductUpdate EHistoryAction = 110 //"product_update"
	EActionProductDelete EHistoryAction = 120 //"product_delete"
	EActionAssetCreate   EHistoryAction = 200 //"asset_create"
	EActionAssetUpdate   EHistoryAction = 210 //"asset_update"
	EActionAssetDelete   EHistoryAction = 220 //"asset_delete"
	EActionAssetMerge    EHistoryAction = 230 //"asset_merge"
	EActionAssetSplit    EHistoryAction = 240 //"asset_split"
)

var HistoryActionMap = map[EHistoryAction]string{
	EActionProductCreate: "Tạo sản phẩm",
	EActionProductUpdate: "Cập nhật sản phẩm",
	EActionProductDelete: "Xóa sản phẩm",
	EActionAssetCreate:   "Tạo tài sản",
	EActionAssetUpdate:   "Cập nhật tài sản",
	EActionAssetDelete:   "Xóa tài sản",
	EActionAssetMerge:    "Gộp tài sản",
	EActionAssetSplit:    "Tách tài sản",
}

var HistoryAsset = []EHistoryAction{
	EActionAssetCreate,
	EActionAssetUpdate,
	EActionAssetDelete,
	EActionAssetMerge,
	EActionAssetSplit,
}

var HistoryProduct = []EHistoryAction{
	EActionProductCreate,
	EActionProductUpdate,
	EActionProductDelete,
}
