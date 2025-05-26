package enums

type ProductVisibility uint

const (
	ProductVisiblePrivate  ProductVisibility = 10
	ProductVisibleInternal ProductVisibility = 20
	ProductVisiblePublic   ProductVisibility = 30
)
