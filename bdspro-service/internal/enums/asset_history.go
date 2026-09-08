package enums

type AssetHistory uint

const (
	AssetHistory_Merge    AssetHistory = 10
	AssetHistory_Slit     AssetHistory = 20
	AssetHistory_New      AssetHistory = 30
	AssetHistory_Update   AssetHistory = 40
	AssetHistory_Archived AssetHistory = 50
	AssetHistory_Deleted  AssetHistory = 60
	AssetHistory_UpPost   AssetHistory = 70
	// AssetHistory_UpPost   AssetHistory = 7
	// AssetHistory_MergeAll AssetHistory = 3
)

var AssetHistoryNames = map[AssetHistory]string{
	AssetHistory_Merge:    "Gộp",
	AssetHistory_Slit:     "Chia",
	AssetHistory_New:      "Tạo",
	AssetHistory_Update:   "Cập nhật",
	AssetHistory_Archived: "Lưu trữ",
	AssetHistory_Deleted:  "Xóa",
	AssetHistory_UpPost:   "Đăng tin",
}
