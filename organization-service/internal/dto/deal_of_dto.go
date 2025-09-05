package dto

type DealOfGroupDTO struct {
	ID      uint64 `json:"id"`
	DealID  uint64 `json:"dealId"`
	GroupID uint64 `json:"groupId"`
}

type DealOfOrganizationDTO struct {
	ID             uint64 `json:"id"`
	DealID         uint64 `json:"dealId"`
	OrganizationID uint64 `json:"organizationId"`
}

type DealOfBranchDTO struct {
	ID       uint64 `json:"id"`
	DealID   uint64 `json:"dealId"`
	BranchID uint64 `json:"branchId"`
}
