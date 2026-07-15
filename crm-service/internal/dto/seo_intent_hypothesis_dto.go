package dto

import "time"

type CreateSeoIntentHypothesisRequest struct {
	UserNeedID uint64
	Code       string
	Name       string
}

type UpdateSeoIntentHypothesisRequest struct {
	ID              uint64
	ExpectedVersion uint32

	Name            string
	GoalStatement   string
	PrimaryQuestion string
	ExpectedResult  string

	SearchAction string
	SubjectType  string
	ScopeType    string
	Confidence   string

	OwnerCode   string
	ReviewDueAt *time.Time
}

type SeoIntentHypothesisListRequest struct {
	ObjectiveID uint64
	UserNeedID  uint64

	Keyword      string
	SearchAction string
	SubjectType  string
	ScopeType    string
	Status       string
	Confidence   string

	Page uint32
	Size uint32
}

type LinkSeoIntentEvidenceRequest struct {
	IntentHypothesisID uint64
	QueryEvidenceID    uint64
	RelationType       string
	RelevanceNote      string
}

type UnlinkSeoIntentEvidenceRequest struct {
	IntentHypothesisID uint64
	QueryEvidenceID    uint64
}

type SeoIntentHypothesisCommandRequest struct {
	ID              uint64
	ExpectedVersion uint32
	Reason          string
}
