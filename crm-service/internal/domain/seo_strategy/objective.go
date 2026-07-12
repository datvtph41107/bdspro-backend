package seo_strategy

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type ObjectiveStatus string

const (
	ObjectiveStatusDraft    ObjectiveStatus = "draft"
	ObjectiveStatusActive   ObjectiveStatus = "active"
	ObjectiveStatusPaused   ObjectiveStatus = "paused"
	ObjectiveStatusArchived ObjectiveStatus = "archived"
)

type ObjectiveConfidence string

const (
	ObjectiveConfidenceExploratory ObjectiveConfidence = "exploratory"
	ObjectiveConfidenceSupported   ObjectiveConfidence = "supported"
	ObjectiveConfidenceValidated   ObjectiveConfidence = "validated"
)

type EvidenceType string

const (
	EvidenceTypeAssumption     EvidenceType = "assumption"
	EvidenceTypeUserInterview  EvidenceType = "user_interview"
	EvidenceTypeSearchConsole  EvidenceType = "search_console"
	EvidenceTypeAnalytics      EvidenceType = "analytics"
	EvidenceTypeInternalSearch EvidenceType = "internal_search"
	EvidenceTypeSupportRequest EvidenceType = "support_request"
	EvidenceTypeCRM            EvidenceType = "crm"
	EvidenceTypeMarketResearch EvidenceType = "market_research"
	EvidenceTypeExperiment     EvidenceType = "experiment"
	EvidenceTypeExpertReview   EvidenceType = "expert_review"
)

type EvidenceStance string

const (
	EvidenceStanceSupports    EvidenceStance = "supports"
	EvidenceStanceContradicts EvidenceStance = "contradicts"
	EvidenceStanceContext     EvidenceStance = "context"
)

type EvidenceStatus string

const (
	EvidenceStatusActive     EvidenceStatus = "active"
	EvidenceStatusSuperseded EvidenceStatus = "superseded"
	EvidenceStatusRejected   EvidenceStatus = "rejected"
)

type MetricLayer string

const (
	MetricLayerSearch    MetricLayer = "search"
	MetricLayerTechnical MetricLayer = "technical"
	MetricLayerUser      MetricLayer = "user"
	MetricLayerProduct   MetricLayer = "product"
	MetricLayerBusiness  MetricLayer = "business"
	MetricLayerTrust     MetricLayer = "trust"
)

type TargetDirection string

const (
	TargetDirectionIncrease TargetDirection = "increase"
	TargetDirectionDecrease TargetDirection = "decrease"
	TargetDirectionMaintain TargetDirection = "maintain"
)

type InstrumentationStatus string

const (
	InstrumentationStatusUnknown      InstrumentationStatus = "unknown"
	InstrumentationStatusPlanned      InstrumentationStatus = "planned"
	InstrumentationStatusInstrumented InstrumentationStatus = "instrumented"
	InstrumentationStatusVerified     InstrumentationStatus = "verified"
)

type MetricStatus string

const (
	MetricStatusActive     MetricStatus = "active"
	MetricStatusDeprecated MetricStatus = "deprecated"
)

type SeoObjective struct {
	ID   uint64
	Code string
	Name string

	AudienceSummary  string
	Situation        string
	ProblemStatement string

	UserOutcome     string
	ProductOutcome  string
	BusinessOutcome string

	NonGoals    []string
	Constraints []string

	Status     ObjectiveStatus
	Confidence ObjectiveConfidence

	OwnerCode   string
	ReviewDueAt *time.Time

	LockVersion uint32

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type SeoObjectiveEvidence struct {
	ID          uint64
	ObjectiveID uint64

	Type   EvidenceType
	Stance EvidenceStance

	SourceSystem    string
	SourceReference string

	Finding    string
	Confidence ObjectiveConfidence

	ObservedAt *time.Time
	ValidUntil *time.Time

	Status   EvidenceStatus
	Metadata json.RawMessage

	CreatedAt time.Time
	UpdatedAt time.Time
}

type SeoObjectiveMetric struct {
	ID          uint64
	ObjectiveID uint64

	Code string
	Name string

	Layer       MetricLayer
	Description string

	SourceSystem      string
	CalculationMethod string

	Unit            string
	TargetDirection TargetDirection

	BaselineValue *float64
	TargetValue   *float64

	IsPrimary bool

	InstrumentationStatus InstrumentationStatus
	Status                MetricStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

type SeoObjectiveAggregate struct {
	Objective *SeoObjective
	Evidence  []SeoObjectiveEvidence
	Metrics   []SeoObjectiveMetric
	Readiness ObjectiveReadiness
}

var objectiveCodePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{2,119}$`)

// CreateSeoObjectiveDraft tạo domain draft từ dữ liệu request.
//
// Không dùng tiền tố New* vì Wire generator của dự án tự động
// xem mọi hàm New* là dependency-injection provider.
func CreateSeoObjectiveDraft(
	code string,
	name string,
	now time.Time,
) (*SeoObjective, error) {
	code = NormalizeObjectiveCode(code)
	name = strings.TrimSpace(name)

	objective := &SeoObjective{
		Code:        code,
		Name:        name,
		Status:      ObjectiveStatusDraft,
		Confidence:  ObjectiveConfidenceExploratory,
		LockVersion: 1,
		NonGoals:    []string{},
		Constraints: []string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := ValidateObjectiveIdentity(objective); err != nil {
		return nil, err
	}

	return objective, nil
}

func NormalizeObjectiveCode(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "-", "_")
	value = strings.Join(strings.Fields(value), "_")
	return value
}

func ValidateObjectiveIdentity(objective *SeoObjective) error {
	if objective == nil {
		return errors.New("objective không tồn tại")
	}

	if !objectiveCodePattern.MatchString(objective.Code) {
		return errors.New(
			"objective code phải dùng lower_snake_case và dài từ 3 đến 120 ký tự",
		)
	}

	if strings.TrimSpace(objective.Name) == "" {
		return errors.New("objective name không được để trống")
	}

	return nil
}

func ValidateObjectiveDefinition(objective *SeoObjective) error {
	if err := ValidateObjectiveIdentity(objective); err != nil {
		return err
	}

	if !IsValidObjectiveConfidence(objective.Confidence) {
		return fmt.Errorf(
			"objective confidence không hợp lệ: %s",
			objective.Confidence,
		)
	}

	return nil
}

func ValidateEvidence(evidence *SeoObjectiveEvidence) error {
	if evidence == nil {
		return errors.New("evidence không tồn tại")
	}

	if evidence.ObjectiveID == 0 {
		return errors.New("objectiveId không hợp lệ")
	}

	if !IsValidEvidenceType(evidence.Type) {
		return fmt.Errorf("evidence type không hợp lệ: %s", evidence.Type)
	}

	if !IsValidEvidenceStance(evidence.Stance) {
		return fmt.Errorf("evidence stance không hợp lệ: %s", evidence.Stance)
	}

	if !IsValidObjectiveConfidence(evidence.Confidence) {
		return fmt.Errorf(
			"evidence confidence không hợp lệ: %s",
			evidence.Confidence,
		)
	}

	if strings.TrimSpace(evidence.Finding) == "" {
		return errors.New("evidence finding không được để trống")
	}

	if len(evidence.Metadata) == 0 {
		evidence.Metadata = json.RawMessage(`{}`)
	}

	if !json.Valid(evidence.Metadata) {
		return errors.New("evidence metadata phải là JSON hợp lệ")
	}

	return nil
}

func ValidateMetric(metric *SeoObjectiveMetric) error {
	if metric == nil {
		return errors.New("metric không tồn tại")
	}

	if metric.ObjectiveID == 0 {
		return errors.New("objectiveId không hợp lệ")
	}

	metric.Code = NormalizeObjectiveCode(metric.Code)
	metric.Name = strings.TrimSpace(metric.Name)

	if !objectiveCodePattern.MatchString(metric.Code) {
		return errors.New("metric code phải dùng lower_snake_case")
	}

	if metric.Name == "" {
		return errors.New("metric name không được để trống")
	}

	if !IsValidMetricLayer(metric.Layer) {
		return fmt.Errorf("metric layer không hợp lệ: %s", metric.Layer)
	}

	if !IsValidTargetDirection(metric.TargetDirection) {
		return fmt.Errorf(
			"target direction không hợp lệ: %s",
			metric.TargetDirection,
		)
	}

	if !IsValidInstrumentationStatus(metric.InstrumentationStatus) {
		return fmt.Errorf(
			"instrumentation status không hợp lệ: %s",
			metric.InstrumentationStatus,
		)
	}

	if metric.Status != MetricStatusActive &&
		metric.Status != MetricStatusDeprecated {
		return fmt.Errorf("metric status không hợp lệ: %s", metric.Status)
	}

	return nil
}

func IsValidObjectiveConfidence(value ObjectiveConfidence) bool {
	switch value {
	case ObjectiveConfidenceExploratory,
		ObjectiveConfidenceSupported,
		ObjectiveConfidenceValidated:
		return true
	default:
		return false
	}
}

func IsValidEvidenceType(value EvidenceType) bool {
	switch value {
	case EvidenceTypeAssumption,
		EvidenceTypeUserInterview,
		EvidenceTypeSearchConsole,
		EvidenceTypeAnalytics,
		EvidenceTypeInternalSearch,
		EvidenceTypeSupportRequest,
		EvidenceTypeCRM,
		EvidenceTypeMarketResearch,
		EvidenceTypeExperiment,
		EvidenceTypeExpertReview:
		return true
	default:
		return false
	}
}

func IsValidEvidenceStance(value EvidenceStance) bool {
	switch value {
	case EvidenceStanceSupports,
		EvidenceStanceContradicts,
		EvidenceStanceContext:
		return true
	default:
		return false
	}
}

func IsValidMetricLayer(value MetricLayer) bool {
	switch value {
	case MetricLayerSearch,
		MetricLayerTechnical,
		MetricLayerUser,
		MetricLayerProduct,
		MetricLayerBusiness,
		MetricLayerTrust:
		return true
	default:
		return false
	}
}

func IsValidTargetDirection(value TargetDirection) bool {
	switch value {
	case TargetDirectionIncrease,
		TargetDirectionDecrease,
		TargetDirectionMaintain:
		return true
	default:
		return false
	}
}

func IsValidInstrumentationStatus(value InstrumentationStatus) bool {
	switch value {
	case InstrumentationStatusUnknown,
		InstrumentationStatusPlanned,
		InstrumentationStatusInstrumented,
		InstrumentationStatusVerified:
		return true
	default:
		return false
	}
}
