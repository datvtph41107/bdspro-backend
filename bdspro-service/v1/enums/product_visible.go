package enums

type ProductVisibility uint

const (
	ProductVisiblePrivate  ProductVisibility = 1
	ProductVisibleInternal ProductVisibility = 2
	ProductVisiblePublic   ProductVisibility = 3
)
