package modules

import (
	"bdspro/internal/enums"
	"fmt"
	"math"
)

type FieldChange struct {
	FieldPath  string      `json:"fieldPath"`
	OldValue   interface{} `json:"oldValue"`
	NewValue   interface{} `json:"newValue"`
	ChangeType string      `json:"changeType"`
}

type ImpactRule struct {
	FieldPath        string            `json:"fieldPath"`
	FieldCategory    FieldCategory     `json:"fieldCategory"`
	CriticalLevel    enums.ImpactLevel `json:"criticalLevel"`
	Threshold        *float64          `json:"threshold,omitempty"`
	AffectedEntities []string          `json:"affectedEntities"`
	Description      string            `json:"description"`
}

type FieldCategory int

const (
	FieldCategoryCore      FieldCategory = 10
	FieldCategoryLegal     FieldCategory = 20
	FieldCategoryGeometry  FieldCategory = 30
	FieldCategoryPhysical  FieldCategory = 40
	FieldCategoryMetadata  FieldCategory = 50
	FieldCategoryLifecycle FieldCategory = 60
)

var ImpactRules = []ImpactRule{
	{
		FieldPath:        "property_identify_id",
		FieldCategory:    FieldCategoryCore,
		CriticalLevel:    enums.ImpactLevelCritical,
		AffectedEntities: []string{"product", "listing", "asset"},
		Description:      "Thay đổi định danh BĐS",
	},
	{
		FieldPath:        "info.property_type_id",
		FieldCategory:    FieldCategoryCore,
		CriticalLevel:    enums.ImpactLevelMajor,
		AffectedEntities: []string{"product", "listing"},
		Description:      "Thay đổi loại BĐS",
	},
	{
		FieldPath:        "location.latitude",
		FieldCategory:    FieldCategoryGeometry,
		CriticalLevel:    enums.ImpactLevelCritical,
		AffectedEntities: []string{"product", "listing", "asset"},
		Threshold:        float64Ptr(0.0001),
		Description:      "Thay đổi tọa độ",
	},
	{
		FieldPath:        "location.longitude",
		FieldCategory:    FieldCategoryGeometry,
		CriticalLevel:    enums.ImpactLevelCritical,
		AffectedEntities: []string{"product", "listing", "asset"},
		Threshold:        float64Ptr(0.0001),
		Description:      "Thay đổi tọa độ",
	},
	{
		FieldPath:        "land_info.area_total",
		FieldCategory:    FieldCategoryPhysical,
		CriticalLevel:    enums.ImpactLevelModerate,
		AffectedEntities: []string{"product", "listing"},
		Threshold:        float64Ptr(0.1),
		Description:      "Thay đổi diện tích",
	},
	{
		FieldPath:        "lineage_status",
		FieldCategory:    FieldCategoryLifecycle,
		CriticalLevel:    enums.ImpactLevelCritical,
		AffectedEntities: []string{"product", "listing", "asset"},
		Description:      "Thay đổi trạng thái lineage",
	},
}

type RuleEngine interface {
	EvaluateImpact(changes []FieldChange, actorRole enums.ActorRole) (*ImpactEvaluationResult, error)
	GetRulesForField(fieldPath string) []ImpactRule
}

type ImpactEvaluationResult struct {
	ImpactLevel      enums.ImpactLevel `json:"impactLevel"`
	AffectedEntities map[string]bool   `json:"affectedEntities"`
	CriticalChanges  []FieldChange     `json:"criticalChanges"`
	Warnings         []string          `json:"warnings"`
}

type ruleEngineImpl struct {
	rules []ImpactRule
}

func NewRuleEngine() RuleEngine {
	return &ruleEngineImpl{rules: ImpactRules}
}

func (e *ruleEngineImpl) EvaluateImpact(changes []FieldChange, actorRole enums.ActorRole) (*ImpactEvaluationResult, error) {
	result := &ImpactEvaluationResult{
		ImpactLevel:      enums.ImpactLevelNone,
		AffectedEntities: make(map[string]bool),
		CriticalChanges:  []FieldChange{},
		Warnings:         []string{},
	}

	for _, change := range changes {
		applicableRules := e.getRulesForField(change.FieldPath)

		for _, rule := range applicableRules {
			level := rule.CriticalLevel

			if rule.Threshold != nil && change.OldValue != nil && change.NewValue != nil {
				percentChange := calculatePercentChange(change.OldValue, change.NewValue)
				if percentChange > *rule.Threshold {
					level = upgradeImpactLevel(level)
					result.Warnings = append(result.Warnings,
						formatWarning(rule, percentChange))
				}
			}

			if level > result.ImpactLevel {
				result.ImpactLevel = level
			}

			for _, entity := range rule.AffectedEntities {
				result.AffectedEntities[entity] = true
			}

			if level >= enums.ImpactLevelMajor {
				result.CriticalChanges = append(result.CriticalChanges, change)
			}
		}
	}

	if actorRole.CanBypassRestrictions() && result.ImpactLevel > enums.ImpactLevelModerate {
		result.ImpactLevel = enums.ImpactLevelModerate
		result.Warnings = append(result.Warnings, "Admin override: impact level reduced")
	}

	return result, nil
}

func (e *ruleEngineImpl) GetRulesForField(fieldPath string) []ImpactRule {
	var rules []ImpactRule
	for _, rule := range e.rules {
		if matchesFieldPath(rule.FieldPath, fieldPath) {
			rules = append(rules, rule)
		}
	}
	return rules
}

func (e *ruleEngineImpl) getRulesForField(fieldPath string) []ImpactRule {
	return e.GetRulesForField(fieldPath)
}

func calculatePercentChange(oldVal, newVal interface{}) float64 {
	oldFloat, oldOk := toFloat64(oldVal)
	newFloat, newOk := toFloat64(newVal)
	if !oldOk || !newOk || oldFloat == 0 {
		return 0
	}
	return math.Abs((newFloat - oldFloat) / oldFloat)
}

func upgradeImpactLevel(level enums.ImpactLevel) enums.ImpactLevel {
	switch level {
	case enums.ImpactLevelMinor:
		return enums.ImpactLevelModerate
	case enums.ImpactLevelModerate:
		return enums.ImpactLevelMajor
	case enums.ImpactLevelMajor:
		return enums.ImpactLevelCritical
	default:
		return level
	}
}

func formatWarning(rule ImpactRule, percentChange float64) string {
	return fmt.Sprintf("%s: thay đổi %.1f%%", rule.Description, percentChange*100)
}

func matchesFieldPath(rulePath, actualPath string) bool {
	return rulePath == actualPath
}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case *float64:
		if val != nil {
			return *val, true
		}
	case *float32:
		if val != nil {
			return float64(*val), true
		}
	}
	return 0, false
}

func float64Ptr(v float64) *float64 {
	return &v
}
