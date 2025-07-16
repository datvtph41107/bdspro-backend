package dto

type BdsproPost struct {
	Id              uint64
	ProductId       uint64
	ExpiredAt       string
	Status          uint32
	Visibility      uint32
	Hidden          bool
	Content         string
	TransactionType uint32
	Title           string
	NumDate         int32
	Type            string
	Price           float64
	PriceType       uint32
	PackageVisible  uint32
	Like            uint64
	Comment         uint64
	NumView         uint64
	OwnerId         uint64
	OwnerType       uint32
	MediaList       []*PostMedia
	MediaUrl        string
	MediaType       string
}

type PostMedia struct {
	Id          uint64
	Url         string
	Type        string
	OrderNumber int32
}
