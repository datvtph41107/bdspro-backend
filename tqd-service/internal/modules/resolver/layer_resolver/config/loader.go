package config

import (
	_ "embed"
	"sync"
	"tqd/internal/modules/resolver/layer_resolver/types"

	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var configYAML []byte

var (
	instance *Config
	once     sync.Once
)

// Config - Cấu trúc config từ YAML
type Config struct {
	Weights        *WeightsConfig             `yaml:"weights"`
	LegalScores    map[uint32]float64         `yaml:"legal_scores"`
	ConflictMatrix map[string]map[string]bool `yaml:"conflict_matrix"`
	GroupPriority  map[string]int             `yaml:"group_priority"`
	Thresholds     *ThresholdsConfig          `yaml:"thresholds"`
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

// GetConfig - Singleton, load config từ YAML
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := yaml.Unmarshal(configYAML, instance); err != nil {
			panic("failed to load layer resolver config: " + err.Error())
		}
	})
	return instance
}

// ToTypesConfig - Chuyển đổi sang types.Config (dùng cho resolver engine)
func (c *Config) ToTypesConfig() *types.Config {
	return &types.Config{
		LegalWeight:          c.Weights.LegalWeight,
		PriorityWeight:       c.Weights.PriorityWeight,
		AreaWeight:           c.Weights.AreaWeight,
		LegalScores:          c.LegalScores,
		ConflictMatrix:       c.ConflictMatrix,
		DefaultGroupPriority: c.GroupPriority,
	}
}
