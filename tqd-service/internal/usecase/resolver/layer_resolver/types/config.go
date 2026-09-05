package types

// RiskThreshold - ngưỡng phân loại rủi ro (data-driven)
type RiskThreshold struct {
	MinScore       float64
	Level          string
	CanBuild       bool
	Recommendation string
}

// LandUseRiskWeight - trọng số rủi ro cho từng loại đất
type LandUseRiskWeight struct {
	Multiplier float64
	Label      string
}

// Config - Cấu hình cho resolver engine
type Config struct {
	LegalWeight          float64
	PriorityWeight       float64
	AreaWeight           float64
	LegalScores          map[uint32]float64
	ConflictMatrix       map[string]map[string]bool
	DefaultGroupPriority map[string]int

	LandUseRiskWeights map[string]*LandUseRiskWeight
	DefaultRiskWeight  *LandUseRiskWeight
	MultiLayerPenalty  float64
	RiskThresholds     []*RiskThreshold
	DefaultPriority    int

	ConflictType           string
	ConflictSeverity       string
	ConflictRecommendation string
	ConflictActionRequired string

	WarningOverlapExceeded string
	WarningMultiLayer      string
	WarningConflict        string

	NoPlanningImpactMsg    string
	NoSignificantImpactMsg string

	StatusProhibited     string
	StatusConditional    string
	StatusAllowed        string
	ConflictStatusSuffix string

	LandTypeNames         map[uint64]string
	ConflictSummarySuffix string
}

// DefaultConfig - Config mặc định (fallback), tham chiếu constants.go
func DefaultConfig() *Config {
	return &Config{
		LegalWeight:            DefaultLegalWeight,
		PriorityWeight:         DefaultPriorityWeight,
		AreaWeight:             DefaultAreaWeight,
		LegalScores:            DefaultLegalScores,
		ConflictMatrix:         make(map[string]map[string]bool),
		DefaultGroupPriority:   make(map[string]int),
		LandUseRiskWeights:     DefaultLandUseRiskWeights,
		DefaultRiskWeight:      DefaultRiskWeightValue,
		MultiLayerPenalty:      DefaultMultiLayerPenaltyValue,
		RiskThresholds:         DefaultRiskThresholds,
		DefaultPriority:        DefaultPriorityValue,
		ConflictType:           DefaultConflictType,
		ConflictSeverity:       DefaultConflictSeverity,
		ConflictRecommendation: DefaultConflictRecommendation,
		ConflictActionRequired: DefaultConflictActionRequired,
		WarningOverlapExceeded: DefaultWarningOverlapExceeded,
		WarningMultiLayer:      DefaultWarningMultiLayer,
		WarningConflict:        DefaultWarningConflict,
		NoPlanningImpactMsg:    DefaultNoPlanningImpactMsg,
		NoSignificantImpactMsg: DefaultNoSignificantImpactMsg,
		StatusProhibited:       DefaultStatusProhibited,
		StatusConditional:      DefaultStatusConditional,
		StatusAllowed:          DefaultStatusAllowed,
		ConflictStatusSuffix:   DefaultConflictStatusSuffix,
		LandTypeNames:          DefaultLandTypeNames,
		ConflictSummarySuffix:  DefaultConflictSummarySuffix,
	}
}
