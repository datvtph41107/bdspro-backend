package seo_strategy

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type AudienceDecisionRole string

const (
	AudienceDecisionRoleUnknown           AudienceDecisionRole = "unknown"
	AudienceDecisionRoleSelfDecisionMaker AudienceDecisionRole = "self_decision_maker"
	AudienceDecisionRoleAdvisor           AudienceDecisionRole = "advisor"
	AudienceDecisionRoleResearcher        AudienceDecisionRole = "researcher"
	AudienceDecisionRoleOperator          AudienceDecisionRole = "operator"
)

type AudienceKnowledgeLevel string

const (
	AudienceKnowledgeUnknown      AudienceKnowledgeLevel = "unknown"
	AudienceKnowledgeBeginner     AudienceKnowledgeLevel = "beginner"
	AudienceKnowledgeIntermediate AudienceKnowledgeLevel = "intermediate"
	AudienceKnowledgeAdvanced     AudienceKnowledgeLevel = "advanced"
	AudienceKnowledgeExpert       AudienceKnowledgeLevel = "expert"
)

type AudienceStatus string

const (
	AudienceStatusDraft    AudienceStatus = "draft"
	AudienceStatusActive   AudienceStatus = "active"
	AudienceStatusPaused   AudienceStatus = "paused"
	AudienceStatusArchived AudienceStatus = "archived"
)

type SeoAudienceSegment struct {
	ID uint64

	Code        string
	Name        string
	Description string

	DecisionRole   AudienceDecisionRole
	KnowledgeLevel AudienceKnowledgeLevel
	Status         AudienceStatus

	OwnerCode   string
	ReviewDueAt *time.Time
	LockVersion uint32

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type SeoAudienceAggregate struct {
	Audience  *SeoAudienceSegment
	Readiness AudienceReadiness
}

// CreateSeoAudienceSegmentDraft tạo domain object từ dữ liệu runtime.
// Không dùng New* vì New* được dành cho Wire providers trong repository này.
func CreateSeoAudienceSegmentDraft(code, name string, now time.Time) (*SeoAudienceSegment, error) {
	audience := &SeoAudienceSegment{
		Code:           NormalizeStrategyCode(code),
		Name:           strings.TrimSpace(name),
		DecisionRole:   AudienceDecisionRoleUnknown,
		KnowledgeLevel: AudienceKnowledgeUnknown,
		Status:         AudienceStatusDraft,
		LockVersion:    1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := ValidateAudienceIdentity(audience); err != nil {
		return nil, err
	}
	return audience, nil
}

func ValidateAudienceIdentity(audience *SeoAudienceSegment) error {
	if audience == nil {
		return errors.New("audience không tồn tại")
	}
	if err := ValidateStrategyCode(audience.Code); err != nil {
		return err
	}
	if strings.TrimSpace(audience.Name) == "" {
		return errors.New("audience name không được để trống")
	}
	return nil
}

func ValidateAudienceDefinition(audience *SeoAudienceSegment) error {
	if err := ValidateAudienceIdentity(audience); err != nil {
		return err
	}
	if !IsValidAudienceDecisionRole(audience.DecisionRole) {
		return fmt.Errorf("decision role không hợp lệ: %s", audience.DecisionRole)
	}
	if !IsValidAudienceKnowledgeLevel(audience.KnowledgeLevel) {
		return fmt.Errorf("knowledge level không hợp lệ: %s", audience.KnowledgeLevel)
	}
	if !IsValidAudienceStatus(audience.Status) {
		return fmt.Errorf("audience status không hợp lệ: %s", audience.Status)
	}
	return nil
}

func IsValidAudienceDecisionRole(value AudienceDecisionRole) bool {
	switch value {
	case AudienceDecisionRoleUnknown,
		AudienceDecisionRoleSelfDecisionMaker,
		AudienceDecisionRoleAdvisor,
		AudienceDecisionRoleResearcher,
		AudienceDecisionRoleOperator:
		return true
	default:
		return false
	}
}

func IsValidAudienceKnowledgeLevel(value AudienceKnowledgeLevel) bool {
	switch value {
	case AudienceKnowledgeUnknown,
		AudienceKnowledgeBeginner,
		AudienceKnowledgeIntermediate,
		AudienceKnowledgeAdvanced,
		AudienceKnowledgeExpert:
		return true
	default:
		return false
	}
}

func IsValidAudienceStatus(value AudienceStatus) bool {
	switch value {
	case AudienceStatusDraft, AudienceStatusActive, AudienceStatusPaused, AudienceStatusArchived:
		return true
	default:
		return false
	}
}
