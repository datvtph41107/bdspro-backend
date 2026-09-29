package config

import (
	_ "embed"
	"fmt"
	"sync"
	"tqd/internal/usecase/resolver/layer_resolver/types"

	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var configYAML []byte

var (
	instance    *Config
	instanceErr error
	once        sync.Once
)

// Config - Cấu trúc config từ YAML
type Config struct {
	Weights         *WeightsConfig               `yaml:"weights"`
	LegalScores     map[uint32]float64           `yaml:"legal_scores"`
	ConflictMatrix  map[string]map[string]bool   `yaml:"conflict_matrix"`
	GroupPriority   map[string]int               `yaml:"group_priority"`
	Thresholds      *ThresholdsConfig            `yaml:"thresholds"`
	RiskWeights     map[string]*RiskWeightConfig `yaml:"risk_weights"`
	RiskThresholds  []*RiskThresholdConfig       `yaml:"risk_thresholds"`
	ConflictLabels  *ConflictLabelsConfig        `yaml:"conflict_labels"`
	WarningMessages *WarningMessagesConfig       `yaml:"warning_messages"`
	AssessMessages  *AssessMessagesConfig        `yaml:"assess_messages"`
	StatusLabels    *StatusLabelsConfig          `yaml:"status_labels"`
	LandTypeNames   map[uint64]string            `yaml:"land_type_names"`
	SummaryLabels   *SummaryLabelsConfig         `yaml:"summary_labels"`
}

type WeightsConfig struct {
	LegalWeight    float64 `yaml:"legal_weight"`
	PriorityWeight float64 `yaml:"priority_weight"`
	AreaWeight     float64 `yaml:"area_weight"`
}

type ThresholdsConfig struct {
	MinOverlapPercent        float64 `yaml:"min_overlap_percent"`
	ConflictWarningThreshold float64 `yaml:"conflict_warning_threshold"`
	HighRiskThreshold        float64 `yaml:"high_risk_threshold"`
	CriticalRiskThreshold    float64 `yaml:"critical_risk_threshold"`
}

type RiskWeightConfig struct {
	Multiplier float64 `yaml:"multiplier"`
	Label      string  `yaml:"label"`
}

type RiskThresholdConfig struct {
	MinScore       float64 `yaml:"min_score"`
	Level          string  `yaml:"level"`
	CanBuild       bool    `yaml:"can_build"`
	Recommendation string  `yaml:"recommendation"`
}

type ConflictLabelsConfig struct {
	Type           string `yaml:"type"`
	Severity       string `yaml:"severity"`
	Recommendation string `yaml:"recommendation"`
	ActionRequired string `yaml:"action_required"`
}

type WarningMessagesConfig struct {
	OverlapExceeded string `yaml:"overlap_exceeded"`
	MultiLayer      string `yaml:"multi_layer"`
	Conflict        string `yaml:"conflict"`
}

type AssessMessagesConfig struct {
	NoPlanningImpact    string `yaml:"no_planning_impact"`
	NoSignificantImpact string `yaml:"no_significant_impact"`
}

type StatusLabelsConfig struct {
	Prohibited     string `yaml:"prohibited"`
	Conditional    string `yaml:"conditional"`
	Allowed        string `yaml:"allowed"`
	ConflictSuffix string `yaml:"conflict_suffix"`
}

type SummaryLabelsConfig struct {
	ConflictSuffix string `yaml:"conflict_suffix"`
}

func parseConfig(data []byte) (*Config, error) {
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("decode embedded layer resolver config: %w", err)
	}
	return cfg, nil
}

// GetConfig loads the embedded resolver defaults once and preserves any decode
// error so callers can fail or degrade explicitly instead of panicking in a
// request path.
func GetConfig() (*Config, error) {
	once.Do(func() {
		instance, instanceErr = parseConfig(configYAML)
	})
	return instance, instanceErr
}

// ToTypesConfig - Chuyển đổi sang types.Config (dùng cho resolver engine)
func (c *Config) ToTypesConfig() *types.Config {
	legalWeight := types.DefaultLegalWeight
	priorityWeight := types.DefaultPriorityWeight
	areaWeight := types.DefaultAreaWeight
	if c.Weights != nil {
		if c.Weights.LegalWeight > 0 {
			legalWeight = c.Weights.LegalWeight
		}
		if c.Weights.PriorityWeight > 0 {
			priorityWeight = c.Weights.PriorityWeight
		}
		if c.Weights.AreaWeight > 0 {
			areaWeight = c.Weights.AreaWeight
		}
	}

	legalScores := c.LegalScores
	if legalScores == nil {
		legalScores = types.DefaultLegalScores
	}

	conflictMatrix := c.ConflictMatrix
	if conflictMatrix == nil {
		conflictMatrix = make(map[string]map[string]bool)
	}

	groupPriority := c.GroupPriority
	if groupPriority == nil {
		groupPriority = make(map[string]int)
	}

	cfg := &types.Config{
		LegalWeight:          legalWeight,
		PriorityWeight:       priorityWeight,
		AreaWeight:           areaWeight,
		LegalScores:          legalScores,
		ConflictMatrix:       conflictMatrix,
		DefaultGroupPriority: groupPriority,
		DefaultPriority:      types.DefaultPriorityValue,
		MultiLayerPenalty:    types.DefaultMultiLayerPenaltyValue,
	}

	if c.RiskWeights != nil {
		cfg.LandUseRiskWeights = make(map[string]*types.LandUseRiskWeight, len(c.RiskWeights))
		for k, v := range c.RiskWeights {
			cfg.LandUseRiskWeights[k] = &types.LandUseRiskWeight{Multiplier: v.Multiplier, Label: v.Label}
		}
	}
	if cfg.LandUseRiskWeights == nil {
		cfg.LandUseRiskWeights = types.DefaultLandUseRiskWeights
	}
	cfg.DefaultRiskWeight = types.DefaultRiskWeightValue

	if c.RiskThresholds != nil {
		for _, t := range c.RiskThresholds {
			cfg.RiskThresholds = append(cfg.RiskThresholds, &types.RiskThreshold{
				MinScore: t.MinScore, Level: t.Level, CanBuild: t.CanBuild, Recommendation: t.Recommendation,
			})
		}
	}
	if len(cfg.RiskThresholds) == 0 {
		cfg.RiskThresholds = types.DefaultRiskThresholds
	}

	if c.ConflictLabels != nil {
		cfg.ConflictType = c.ConflictLabels.Type
		cfg.ConflictSeverity = c.ConflictLabels.Severity
		cfg.ConflictRecommendation = c.ConflictLabels.Recommendation
		cfg.ConflictActionRequired = c.ConflictLabels.ActionRequired
	}
	if cfg.ConflictType == "" {
		cfg.ConflictType = types.DefaultConflictType
	}
	if cfg.ConflictSeverity == "" {
		cfg.ConflictSeverity = types.DefaultConflictSeverity
	}
	if cfg.ConflictRecommendation == "" {
		cfg.ConflictRecommendation = types.DefaultConflictRecommendation
	}
	if cfg.ConflictActionRequired == "" {
		cfg.ConflictActionRequired = types.DefaultConflictActionRequired
	}

	if c.WarningMessages != nil {
		cfg.WarningOverlapExceeded = c.WarningMessages.OverlapExceeded
		cfg.WarningMultiLayer = c.WarningMessages.MultiLayer
		cfg.WarningConflict = c.WarningMessages.Conflict
	}
	if cfg.WarningOverlapExceeded == "" {
		cfg.WarningOverlapExceeded = types.DefaultWarningOverlapExceeded
	}
	if cfg.WarningMultiLayer == "" {
		cfg.WarningMultiLayer = types.DefaultWarningMultiLayer
	}
	if cfg.WarningConflict == "" {
		cfg.WarningConflict = types.DefaultWarningConflict
	}

	if c.AssessMessages != nil {
		cfg.NoPlanningImpactMsg = c.AssessMessages.NoPlanningImpact
		cfg.NoSignificantImpactMsg = c.AssessMessages.NoSignificantImpact
	}
	if cfg.NoPlanningImpactMsg == "" {
		cfg.NoPlanningImpactMsg = types.DefaultNoPlanningImpactMsg
	}
	if cfg.NoSignificantImpactMsg == "" {
		cfg.NoSignificantImpactMsg = types.DefaultNoSignificantImpactMsg
	}

	if c.StatusLabels != nil {
		cfg.StatusProhibited = c.StatusLabels.Prohibited
		cfg.StatusConditional = c.StatusLabels.Conditional
		cfg.StatusAllowed = c.StatusLabels.Allowed
		cfg.ConflictStatusSuffix = c.StatusLabels.ConflictSuffix
	}
	if cfg.StatusProhibited == "" {
		cfg.StatusProhibited = types.DefaultStatusProhibited
	}
	if cfg.StatusConditional == "" {
		cfg.StatusConditional = types.DefaultStatusConditional
	}
	if cfg.StatusAllowed == "" {
		cfg.StatusAllowed = types.DefaultStatusAllowed
	}
	if cfg.ConflictStatusSuffix == "" {
		cfg.ConflictStatusSuffix = types.DefaultConflictStatusSuffix
	}

	if c.LandTypeNames != nil {
		cfg.LandTypeNames = c.LandTypeNames
	}
	if cfg.LandTypeNames == nil {
		cfg.LandTypeNames = types.DefaultLandTypeNames
	}

	if c.SummaryLabels != nil {
		cfg.ConflictSummarySuffix = c.SummaryLabels.ConflictSuffix
	}
	if cfg.ConflictSummarySuffix == "" {
		cfg.ConflictSummarySuffix = types.DefaultConflictSummarySuffix
	}

	return cfg
}
