package entity

type DealOfBranch struct {
	ID       uint64 `gorm:"primaryKey;column:id" json:"id"`
	DealID   uint64 `gorm:"column:deal_id;not null" json:"dealId"`
	BranchID uint64 `gorm:"column:branch_id;not null" json:"branchId"`
}

func (DealOfBranch) TableName() string {
	return "deal_of_branch"
}
