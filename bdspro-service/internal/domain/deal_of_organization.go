package domain

type DealOfOrganization struct {
	ID             uint64 `gorm:"primaryKey;column:id" json:"id"`
	DealID         uint64 `gorm:"column:deal_id;not null" json:"dealId"`
	OrganizationID uint64 `gorm:"column:organization_id;not null" json:"organizationId"`
	BranchID       uint64 `gorm:"column:branch_id;not null" json:"branchId"`
}

func (DealOfOrganization) TableName() string {
	return "deal_of_organization"
}
