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
