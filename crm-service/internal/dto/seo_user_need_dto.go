package dto

import "time"

type CreateSeoUserNeedRequest struct {
	ObjectiveID       uint64
	AudienceSegmentID uint64
	Code              string
	Name              string
}

type UpdateSeoUserNeedRequest struct {
	ID                 uint64
	ExpectedVersion    uint32
	Name               string
	DecisionStage      string
	Situation          string
	TriggerStatement   string
	JobToBeDone        string
	BarrierStatement   string
	CurrentAlternative string
	DesiredOutcome     string
	FailureImpact      string
	RiskLevel          string
	Priority           string
	Confidence         string
	OwnerCode          string
	ReviewDueAt        *time.Time
}

type LinkSeoUserNeedEvidenceRequest struct {
	UserNeedID    uint64
	EvidenceID    uint64
	RelationType  string
	RelevanceNote string
}

type UnlinkSeoUserNeedEvidenceRequest struct {
	UserNeedID uint64
	EvidenceID uint64
}

type SeoUserNeedCommandRequest struct {
	ID              uint64
	ExpectedVersion uint32
	Reason          string
}

type SeoUserNeedListRequest struct {
	ObjectiveID uint64
	AudienceID  uint64
	Keyword     string
	Status      string
	Confidence  string
	Priority    string
	Page        uint32
	Size        uint32
}
