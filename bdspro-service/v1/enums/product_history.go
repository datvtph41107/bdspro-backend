package enums

type ProductHistory uint

const (
	ProductHistory_NewChild ProductHistory = 1
	ProductHistory_Merge    ProductHistory = 2
	ProductHistory_MergeAll ProductHistory = 3
	ProductHistory_Devide   ProductHistory = 4
)
