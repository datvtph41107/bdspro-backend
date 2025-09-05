package entity

import "time"

type OrganizationBranchType string

const (
	OrganizationBranchTypeBranch     OrganizationBranchType = "branch"     // Chi nhánh
	OrganizationBranchTypeDepartment OrganizationBranchType = "department" // Phòng ban
	OrganizationBranchTypeStore      OrganizationBranchType = "store"      // Cửa hàng
)

type OrganizationBranch struct {
	Id             uint32
	OrganizationId uint32
	Name           string
	Address        string
	Phone          string
	Email          string
	ManagerId      uint32
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      uint32
	IsActive       bool
	Description    string
	Type           OrganizationBranchType
	TotalMember    uint32
	TotalAsset     uint32
	TotalDeal      uint32
}
