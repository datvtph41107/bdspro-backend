package dto

import (
	"encoding/json"
	"time"
)

type CreateSeoObjectiveRequest struct {
	Code string
	Name string
}

type UpdateSeoObjectiveDefinitionRequest struct {
	ID              uint64
	ExpectedVersion uint32

	Name string

	AudienceSummary  string
	Situation        string
	ProblemStatement string

	UserOutcome     string
	ProductOutcome  string
	BusinessOutcome string

	NonGoals    []string
	Constraints []string

	Confidence string
	OwnerCode  string

	ReviewDueAt *time.Time
}

type AddSeoObjectiveEvidenceRequest struct {
	ObjectiveID uint64

	EvidenceType string
	Stance       string

	SourceSystem    string
	SourceReference string

	Finding    string
	Confidence string

	ObservedAt *time.Time
	ValidUntil *time.Time

	Metadata json.RawMessage
}

type UpsertSeoObjectiveMetricRequest struct {
	ObjectiveID uint64

	MetricCode string
	MetricName string

	MetricLayer string
	Description string

	SourceSystem      string
	CalculationMethod string

	Unit            string
	TargetDirection string

	BaselineValue *float64
	TargetValue   *float64

	IsPrimary bool

	InstrumentationStatus string
	Status                string
}

type SeoObjectiveCommandRequest struct {
	ID              uint64
	ExpectedVersion uint32
	Reason          string
}

type SeoObjectiveListRequest struct {
	Keyword    string
	Status     string
	Confidence string

	Page uint32
	Size uint32
}
