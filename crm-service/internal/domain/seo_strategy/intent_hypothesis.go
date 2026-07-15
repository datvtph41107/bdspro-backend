package seo_strategy

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type IntentSearchAction string

const (
	IntentSearchActionUnknown  IntentSearchAction = "unknown"
	IntentSearchActionLearn    IntentSearchAction = "learn"
	IntentSearchActionFind     IntentSearchAction = "find"
	IntentSearchActionCompare  IntentSearchAction = "compare"
	IntentSearchActionVerify   IntentSearchAction = "verify"
	IntentSearchActionTransact IntentSearchAction = "transact"
	IntentSearchActionMonitor  IntentSearchAction = "monitor"
	IntentSearchActionNavigate IntentSearchAction = "navigate"
)

type IntentSubjectType string

const (
	IntentSubjectUnknown            IntentSubjectType = "unknown"
	IntentSubjectParcel             IntentSubjectType = "parcel"
	IntentSubjectAdministrativeUnit IntentSubjectType = "administrative_unit"
	IntentSubjectPlanningRegion     IntentSubjectType = "planning_region"
	IntentSubjectPlanningDocument   IntentSubjectType = "planning_document"
	IntentSubjectMapLayer           IntentSubjectType = "map_layer"
	IntentSubjectProject            IntentSubjectType = "project"
	IntentSubjectReport             IntentSubjectType = "report"
	IntentSubjectMixed              IntentSubjectType = "mixed"
)

type IntentScopeType string

const (
	IntentScopeUnknown        IntentScopeType = "unknown"
	IntentScopeSpecificEntity IntentScopeType = "specific_entity"
	IntentScopeArea           IntentScopeType = "area"
	IntentScopeTopic          IntentScopeType = "topic"
	IntentScopeTask           IntentScopeType = "task"
	IntentScopeMixed          IntentScopeType = "mixed"
)

type IntentHypothesisStatus string

const (
	IntentHypothesisStatusDraft    IntentHypothesisStatus = "draft"
	IntentHypothesisStatusActive   IntentHypothesisStatus = "active"
	IntentHypothesisStatusPaused   IntentHypothesisStatus = "paused"
	IntentHypothesisStatusArchived IntentHypothesisStatus = "archived"
)

type IntentEvidenceRelation string

const (
	IntentEvidenceSupports    IntentEvidenceRelation = "supports"
	IntentEvidenceContradicts IntentEvidenceRelation = "contradicts"
	IntentEvidenceContext     IntentEvidenceRelation = "context"
)

type SeoIntentHypothesis struct {
	ID uint64

	ObjectiveID uint64
	UserNeedID  uint64

	Code string
	Name string

	GoalStatement   string
	PrimaryQuestion string
	ExpectedResult  string

	SearchAction IntentSearchAction
	SubjectType  IntentSubjectType
	ScopeType    IntentScopeType

	Status     IntentHypothesisStatus
	Confidence ObjectiveConfidence

	OwnerCode   string
	ReviewDueAt *time.Time
	LockVersion uint32

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type SeoIntentEvidenceLink struct {
	ObjectiveID uint64
	UserNeedID  uint64

	IntentHypothesisID uint64
	QueryEvidenceID    uint64

	RelationType  IntentEvidenceRelation
	RelevanceNote string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type LinkedIntentQueryEvidence struct {
	Link     SeoIntentEvidenceLink
	Evidence SeoQueryEvidence
}

type SeoIntentHypothesisAggregate struct {
	Intent    *SeoIntentHypothesis
	UserNeed  *SeoUserNeed
	Evidence  []LinkedIntentQueryEvidence
	Readiness IntentHypothesisReadiness
}

// CreateSeoIntentHypothesisDraft là domain factory runtime, không phải Wire provider.
func CreateSeoIntentHypothesisDraft(
	objectiveID uint64,
	userNeedID uint64,
	code string,
	name string,
	now time.Time,
) (*SeoIntentHypothesis, error) {
	intent := &SeoIntentHypothesis{
		ObjectiveID:  objectiveID,
		UserNeedID:   userNeedID,
		Code:         NormalizeStrategyCode(code),
		Name:         strings.TrimSpace(name),
		SearchAction: IntentSearchActionUnknown,
		SubjectType:  IntentSubjectUnknown,
		ScopeType:    IntentScopeUnknown,
		Status:       IntentHypothesisStatusDraft,
		Confidence:   ObjectiveConfidenceExploratory,
		LockVersion:  1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := ValidateIntentHypothesisIdentity(intent); err != nil {
		return nil, err
	}
	return intent, nil
}

func ValidateIntentHypothesisIdentity(intent *SeoIntentHypothesis) error {
	if intent == nil {
		return errors.New("intent hypothesis không tồn tại")
	}
	if intent.ObjectiveID == 0 {
		return errors.New("objectiveId không hợp lệ")
	}
	if intent.UserNeedID == 0 {
		return errors.New("userNeedId không hợp lệ")
	}
	intent.Code = NormalizeStrategyCode(intent.Code)
	if err := ValidateStrategyCode(intent.Code); err != nil {
		return err
	}
	intent.Name = strings.TrimSpace(intent.Name)
	if intent.Name == "" {
		return errors.New("intent hypothesis name không được để trống")
	}
	if len([]rune(intent.Name)) > 500 {
		return errors.New("intent hypothesis name không được dài quá 500 ký tự")
	}
	return nil
}

func ValidateIntentHypothesisDefinition(intent *SeoIntentHypothesis) error {
	if err := ValidateIntentHypothesisIdentity(intent); err != nil {
		return err
	}

	intent.GoalStatement = strings.TrimSpace(intent.GoalStatement)
	intent.PrimaryQuestion = strings.TrimSpace(intent.PrimaryQuestion)
	intent.ExpectedResult = strings.TrimSpace(intent.ExpectedResult)
	intent.OwnerCode = strings.TrimSpace(intent.OwnerCode)

	if len([]rune(intent.GoalStatement)) > 4000 {
		return errors.New("goal statement không được dài quá 4000 ký tự")
	}
	if len([]rune(intent.PrimaryQuestion)) > 2000 {
		return errors.New("primary question không được dài quá 2000 ký tự")
	}
	if len([]rune(intent.ExpectedResult)) > 4000 {
		return errors.New("expected result không được dài quá 4000 ký tự")
	}
	if len([]rune(intent.OwnerCode)) > 120 {
		return errors.New("owner code không được dài quá 120 ký tự")
	}
	if !IsValidIntentSearchAction(intent.SearchAction) {
		return fmt.Errorf("search action không hợp lệ: %s", intent.SearchAction)
	}
	if !IsValidIntentSubjectType(intent.SubjectType) {
		return fmt.Errorf("subject type không hợp lệ: %s", intent.SubjectType)
	}
	if !IsValidIntentScopeType(intent.ScopeType) {
		return fmt.Errorf("scope type không hợp lệ: %s", intent.ScopeType)
	}
	if !IsValidIntentHypothesisStatus(intent.Status) {
		return fmt.Errorf("intent hypothesis status không hợp lệ: %s", intent.Status)
	}
	if !IsValidObjectiveConfidence(intent.Confidence) {
		return fmt.Errorf("intent hypothesis confidence không hợp lệ: %s", intent.Confidence)
	}
	if intent.LockVersion == 0 {
		return errors.New("lock version phải lớn hơn 0")
	}
	return nil
}

func ValidateIntentEvidenceLink(link *SeoIntentEvidenceLink) error {
	if link == nil {
		return errors.New("intent evidence link không tồn tại")
	}
	if link.ObjectiveID == 0 || link.UserNeedID == 0 ||
		link.IntentHypothesisID == 0 || link.QueryEvidenceID == 0 {
		return errors.New("intent evidence link identity không hợp lệ")
	}
	if !IsValidIntentEvidenceRelation(link.RelationType) {
		return fmt.Errorf("intent evidence relation không hợp lệ: %s", link.RelationType)
	}
	link.RelevanceNote = strings.TrimSpace(link.RelevanceNote)
	if len([]rune(link.RelevanceNote)) > 4000 {
		return errors.New("relevance note không được dài quá 4000 ký tự")
	}
	return nil
}

func IsValidIntentSearchAction(value IntentSearchAction) bool {
	switch value {
	case IntentSearchActionUnknown, IntentSearchActionLearn, IntentSearchActionFind,
		IntentSearchActionCompare, IntentSearchActionVerify, IntentSearchActionTransact,
		IntentSearchActionMonitor, IntentSearchActionNavigate:
		return true
	default:
		return false
	}
}

func IsValidIntentSubjectType(value IntentSubjectType) bool {
	switch value {
	case IntentSubjectUnknown, IntentSubjectParcel, IntentSubjectAdministrativeUnit,
		IntentSubjectPlanningRegion, IntentSubjectPlanningDocument, IntentSubjectMapLayer,
		IntentSubjectProject, IntentSubjectReport, IntentSubjectMixed:
		return true
	default:
		return false
	}
}

func IsValidIntentScopeType(value IntentScopeType) bool {
	switch value {
	case IntentScopeUnknown, IntentScopeSpecificEntity, IntentScopeArea,
		IntentScopeTopic, IntentScopeTask, IntentScopeMixed:
		return true
	default:
		return false
	}
}

func IsValidIntentHypothesisStatus(value IntentHypothesisStatus) bool {
	switch value {
	case IntentHypothesisStatusDraft, IntentHypothesisStatusActive,
		IntentHypothesisStatusPaused, IntentHypothesisStatusArchived:
		return true
	default:
		return false
	}
}

func IsValidIntentEvidenceRelation(value IntentEvidenceRelation) bool {
	switch value {
	case IntentEvidenceSupports, IntentEvidenceContradicts, IntentEvidenceContext:
		return true
	default:
		return false
	}
}

func ParseIntentSearchAction(value string) (IntentSearchAction, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return IntentSearchActionUnknown, nil
	}
	result := IntentSearchAction(value)
	if !IsValidIntentSearchAction(result) {
		return "", fmt.Errorf("search action không hợp lệ: %s", value)
	}
	return result, nil
}

func ParseIntentSubjectType(value string) (IntentSubjectType, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return IntentSubjectUnknown, nil
	}
	result := IntentSubjectType(value)
	if !IsValidIntentSubjectType(result) {
		return "", fmt.Errorf("subject type không hợp lệ: %s", value)
	}
	return result, nil
}

func ParseIntentScopeType(value string) (IntentScopeType, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return IntentScopeUnknown, nil
	}
	result := IntentScopeType(value)
	if !IsValidIntentScopeType(result) {
		return "", fmt.Errorf("scope type không hợp lệ: %s", value)
	}
	return result, nil
}

func ParseIntentHypothesisStatus(value string) (IntentHypothesisStatus, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return IntentHypothesisStatusDraft, nil
	}
	result := IntentHypothesisStatus(value)
	if !IsValidIntentHypothesisStatus(result) {
		return "", fmt.Errorf("intent hypothesis status không hợp lệ: %s", value)
	}
	return result, nil
}

func ParseIntentEvidenceRelation(value string) (IntentEvidenceRelation, error) {
	result := IntentEvidenceRelation(strings.TrimSpace(value))
	if !IsValidIntentEvidenceRelation(result) {
		return "", fmt.Errorf("intent evidence relation không hợp lệ: %s", value)
	}
	return result, nil
}
