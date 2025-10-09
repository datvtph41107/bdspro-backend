package dto

type AdminDealSearchRequest struct {
	Page           int
	Size           int
	Keyword        *string
	Status         *uint32
	DealType       *uint32
	OrganizationID *uint64
	GroupID        *uint64
	OwnerID        *uint64
	OwnerType      *uint32
	StartDate      *string
	EndDate        *string
	MinAmount      *float64
	MaxAmount      *float64
}
