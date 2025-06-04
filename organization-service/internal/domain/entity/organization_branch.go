package entity

import "time"

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
}
