package organization

import "time"

type Organization struct {
	ID                 uint64
	Code               string
	Name               string
	Type               string
	TaxCode            string
	Phone              string
	Email              string
	Address            string
	Website            string
	Description        string
	Status             string
	VerificationStatus string
	WarningLevel       string
	WarningCount       uint32
	OwnerProfileID     uint64
	CreatedBy          uint64
	UpdatedBy          uint64
	CreatedAt          time.Time
	UpdatedAt          time.Time
	ArchivedAt         *time.Time
	ActiveMemberCount  uint64
}

type ListQuery struct {
	Page               uint32
	Size               uint32
	SearchText         string
	Type               string
	Status             string
	VerificationStatus string
	WarningLevel       string
}

type Update struct {
	ID                 uint64
	ActorID            uint64
	Code               *string
	Name               *string
	Type               *string
	TaxCode            *string
	Phone              *string
	Email              *string
	Address            *string
	Website            *string
	Description        *string
	Status             *string
	VerificationStatus *string
	WarningLevel       *string
}

type MemberCheck struct {
	ProfileID uint64
	IsMember  bool
	JoinedAt  time.Time
	MemberID  uint64
	RoleKey   string
	Status    string
}
