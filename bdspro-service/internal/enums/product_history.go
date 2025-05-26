package enums

type ProductHistory uint

const (
	ProductHistory_NewChild ProductHistory = 10
	ProductHistory_Merge    ProductHistory = 20
	ProductHistory_MergeAll ProductHistory = 30
	ProductHistory_Devide   ProductHistory = 40
)
