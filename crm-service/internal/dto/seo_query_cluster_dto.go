package dto

import "time"

type CreateSeoQueryClusterRequest struct {
	IntentHypothesisID uint64
	Code               string
	Name               string
}

type UpdateSeoQueryClusterRequest struct {
	ID              uint64
	ExpectedVersion uint32

	Name          string
	SearchPromise string
	InclusionRule string
	ExclusionRule string
	Confidence    string

	OwnerCode   string
	ReviewDueAt *time.Time
}

type SeoQueryClusterListRequest struct {
	ObjectiveID        uint64
	UserNeedID         uint64
	IntentHypothesisID uint64

	Keyword    string
	Status     string
	Confidence string

	Page uint32
	Size uint32
}

type UpsertSeoQueryClusterMemberRequest struct {
	ClusterID       uint64
	QueryEvidenceID uint64
	ExpectedVersion uint32
	MembershipType  string
	Rationale       string
}

type RemoveSeoQueryClusterMemberRequest struct {
	ClusterID       uint64
	QueryEvidenceID uint64
	ExpectedVersion uint32
}

type SeoQueryClusterCommandRequest struct {
	ID              uint64
	ExpectedVersion uint32
	Reason          string
}
