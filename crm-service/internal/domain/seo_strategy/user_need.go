package seo_strategy

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type UserNeedDecisionStage string

const (
	UserNeedStageUnknown      UserNeedDecisionStage = "unknown"
	UserNeedStageExploration  UserNeedDecisionStage = "exploration"
	UserNeedStageComparison   UserNeedDecisionStage = "comparison"
	UserNeedStageVerification UserNeedDecisionStage = "verification"
	UserNeedStageTransaction  UserNeedDecisionStage = "transaction"
	UserNeedStageOwnership    UserNeedDecisionStage = "ownership"
	UserNeedStageMonitoring   UserNeedDecisionStage = "monitoring"
)

type UserNeedRiskLevel string

const (
	UserNeedRiskUnknown  UserNeedRiskLevel = "unknown"
	UserNeedRiskLow      UserNeedRiskLevel = "low"
	UserNeedRiskMedium   UserNeedRiskLevel = "medium"
	UserNeedRiskHigh     UserNeedRiskLevel = "high"
	UserNeedRiskCritical UserNeedRiskLevel = "critical"
)

type UserNeedPriority string

const (
	UserNeedPriorityUnassessed UserNeedPriority = "unassessed"
	UserNeedPriorityLow        UserNeedPriority = "low"
	UserNeedPriorityMedium     UserNeedPriority = "medium"
	UserNeedPriorityHigh       UserNeedPriority = "high"
	UserNeedPriorityCritical   UserNeedPriority = "critical"
)

type UserNeedStatus string

const (
	UserNeedStatusDraft    UserNeedStatus = "draft"
	UserNeedStatusActive   UserNeedStatus = "active"
	UserNeedStatusPaused   UserNeedStatus = "paused"
	UserNeedStatusArchived UserNeedStatus = "archived"
)

type UserNeedEvidenceRelation string

const (
	UserNeedEvidenceSupports    UserNeedEvidenceRelation = "supports"
	UserNeedEvidenceContradicts UserNeedEvidenceRelation = "contradicts"
	UserNeedEvidenceContext     UserNeedEvidenceRelation = "context"
)

type SeoUserNeed struct {
	ID uint64

	ObjectiveID       uint64
	AudienceSegmentID uint64
	Code              string
	Name              string
	DecisionStage     UserNeedDecisionStage

	Situation          string
	TriggerStatement   string
	JobToBeDone        string
	BarrierStatement   string
	CurrentAlternative string
	DesiredOutcome     string
	FailureImpact      string

	RiskLevel  UserNeedRiskLevel
	Priority   UserNeedPriority
	Status     UserNeedStatus
	Confidence ObjectiveConfidence

	OwnerCode   string
	ReviewDueAt *time.Time
	LockVersion uint32

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type SeoUserNeedEvidenceLink struct {
	ObjectiveID   uint64
	UserNeedID    uint64
	EvidenceID    uint64
	RelationType  UserNeedEvidenceRelation
	RelevanceNote string
	CreatedAt     time.Time
}

type LinkedUserNeedEvidence struct {
	Link     SeoUserNeedEvidenceLink
	Evidence SeoObjectiveEvidence
}

type SeoUserNeedAggregate struct {
	Need      *SeoUserNeed
	Objective *SeoObjective
	Audience  *SeoAudienceSegment
	Evidence  []LinkedUserNeedEvidence
	Readiness UserNeedReadiness
}

// CreateSeoUserNeedDraft là domain factory runtime, không phải Wire provider.
func CreateSeoUserNeedDraft(objectiveID, audienceSegmentID uint64, code, name string, now time.Time) (*SeoUserNeed, error) {
	need := &SeoUserNeed{
		ObjectiveID:       objectiveID,
		AudienceSegmentID: audienceSegmentID,
		Code:              NormalizeStrategyCode(code),
		Name:              strings.TrimSpace(name),
		DecisionStage:     UserNeedStageUnknown,
		RiskLevel:         UserNeedRiskUnknown,
		Priority:          UserNeedPriorityUnassessed,
		Status:            UserNeedStatusDraft,
		Confidence:        ObjectiveConfidenceExploratory,
		LockVersion:       1,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := ValidateUserNeedIdentity(need); err != nil {
		return nil, err
	}
	return need, nil
}

func ValidateUserNeedIdentity(need *SeoUserNeed) error {
	if need == nil {
		return errors.New("user need không tồn tại")
	}
	if need.ObjectiveID == 0 {
		return errors.New("objectiveId không hợp lệ")
	}
	if need.AudienceSegmentID == 0 {
		return errors.New("audienceSegmentId không hợp lệ")
	}
	if err := ValidateStrategyCode(need.Code); err != nil {
		return err
	}
	if strings.TrimSpace(need.Name) == "" {
		return errors.New("user need name không được để trống")
	}
	return nil
}

func ValidateUserNeedDefinition(need *SeoUserNeed) error {
	if err := ValidateUserNeedIdentity(need); err != nil {
		return err
	}
	if !IsValidUserNeedStage(need.DecisionStage) {
		return fmt.Errorf("decision stage không hợp lệ: %s", need.DecisionStage)
	}
	if !IsValidUserNeedRisk(need.RiskLevel) {
		return fmt.Errorf("risk level không hợp lệ: %s", need.RiskLevel)
	}
	if !IsValidUserNeedPriority(need.Priority) {
		return fmt.Errorf("priority không hợp lệ: %s", need.Priority)
	}
	if !IsValidUserNeedStatus(need.Status) {
		return fmt.Errorf("user need status không hợp lệ: %s", need.Status)
	}
	if !IsValidObjectiveConfidence(need.Confidence) {
		return fmt.Errorf("confidence không hợp lệ: %s", need.Confidence)
	}
	return nil
}

func IsValidUserNeedStage(v UserNeedDecisionStage) bool {
	switch v {
	case UserNeedStageUnknown, UserNeedStageExploration, UserNeedStageComparison,
		UserNeedStageVerification, UserNeedStageTransaction, UserNeedStageOwnership, UserNeedStageMonitoring:
		return true
	default:
		return false
	}
}

func IsValidUserNeedRisk(v UserNeedRiskLevel) bool {
	switch v {
	case UserNeedRiskUnknown, UserNeedRiskLow, UserNeedRiskMedium, UserNeedRiskHigh, UserNeedRiskCritical:
		return true
	default:
		return false
	}
}

func IsValidUserNeedPriority(v UserNeedPriority) bool {
	switch v {
	case UserNeedPriorityUnassessed, UserNeedPriorityLow, UserNeedPriorityMedium,
		UserNeedPriorityHigh, UserNeedPriorityCritical:
		return true
	default:
		return false
	}
}

func IsValidUserNeedStatus(v UserNeedStatus) bool {
	switch v {
	case UserNeedStatusDraft, UserNeedStatusActive, UserNeedStatusPaused, UserNeedStatusArchived:
		return true
	default:
		return false
	}
}

func IsValidUserNeedEvidenceRelation(
	value UserNeedEvidenceRelation,
) bool {
	switch value {
	case UserNeedEvidenceSupports,
		UserNeedEvidenceContradicts,
		UserNeedEvidenceContext:
		return true
	default:
		return false
	}
}
