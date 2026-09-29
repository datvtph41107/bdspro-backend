// internal/modules/resolver/layer_resolver/db/loader.go
package db

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"tqd/internal/interface/repo"
	"tqd/internal/usecase/resolver/layer_resolver/config"
	"tqd/internal/usecase/resolver/layer_resolver/types"
)

type ConfigLoader struct {
	configRepo repo.IQHLayerResolverConfigRepo
	cache      map[string]json.RawMessage
	mu         sync.RWMutex
	lastLoad   time.Time
	ttl        time.Duration
}

func NewConfigLoader(configRepo repo.IQHLayerResolverConfigRepo) *ConfigLoader {
	return &ConfigLoader{
		configRepo: configRepo,
		cache:      make(map[string]json.RawMessage),
		ttl:        5 * time.Minute,
	}
}

func (l *ConfigLoader) LoadConfig(ctx context.Context) (*types.Config, error) {
	l.mu.RLock()
	if time.Since(l.lastLoad) < l.ttl && len(l.cache) > 0 {
		defer l.mu.RUnlock()
		return l.buildConfig()
	}
	l.mu.RUnlock()

	l.mu.Lock()
	defer l.mu.Unlock()

	dbConfigs, err := l.configRepo.GetAll(ctx)
	if err != nil {
		// DB override failure falls back to the embedded defaults, but invalid
		// embedded configuration is surfaced rather than panicking.
		yamlConfig, configErr := config.GetConfig()
		if configErr != nil {
			return nil, configErr
		}
		return yamlConfig.ToTypesConfig(), nil
	}

	for k, v := range dbConfigs {
		l.cache[k] = v
	}
	l.lastLoad = time.Now()

	return l.buildConfig()
}

func (l *ConfigLoader) buildConfig() (*types.Config, error) {
	yamlConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	cfg := yamlConfig.ToTypesConfig()

	// Override with DB configs
	if val, ok := l.cache["weights"]; ok {
		var weights struct {
			LegalWeight    float64 `json:"legal_weight"`
			PriorityWeight float64 `json:"priority_weight"`
			AreaWeight     float64 `json:"area_weight"`
		}
		if err := json.Unmarshal(val, &weights); err == nil {
			if weights.LegalWeight > 0 {
				cfg.LegalWeight = weights.LegalWeight
			}
			if weights.PriorityWeight > 0 {
				cfg.PriorityWeight = weights.PriorityWeight
			}
			if weights.AreaWeight > 0 {
				cfg.AreaWeight = weights.AreaWeight
			}
		}
	}

	if val, ok := l.cache["legal_scores"]; ok {
		var scores map[uint32]float64
		if err := json.Unmarshal(val, &scores); err == nil {
			for k, v := range scores {
				cfg.LegalScores[k] = v
			}
		}
	}

	if val, ok := l.cache["conflict_matrix"]; ok {
		var matrix map[string]map[string]bool
		if err := json.Unmarshal(val, &matrix); err == nil {
			cfg.ConflictMatrix = matrix
		}
	}

	if val, ok := l.cache["group_priority"]; ok {
		var priorities map[string]int
		if err := json.Unmarshal(val, &priorities); err == nil {
			cfg.DefaultGroupPriority = priorities
		}
	}

	if val, ok := l.cache["risk_weights"]; ok {
		var rw map[string]*types.LandUseRiskWeight
		if err := json.Unmarshal(val, &rw); err == nil {
			cfg.LandUseRiskWeights = rw
		}
	}

	if val, ok := l.cache["risk_thresholds"]; ok {
		var rt []*types.RiskThreshold
		if err := json.Unmarshal(val, &rt); err == nil {
			cfg.RiskThresholds = rt
		}
	}

	if val, ok := l.cache["conflict_labels"]; ok {
		var cl struct {
			Type           string `json:"type"`
			Severity       string `json:"severity"`
			Recommendation string `json:"recommendation"`
			ActionRequired string `json:"action_required"`
		}
		if err := json.Unmarshal(val, &cl); err == nil {
			if cl.Type != "" {
				cfg.ConflictType = cl.Type
			}
			if cl.Severity != "" {
				cfg.ConflictSeverity = cl.Severity
			}
			if cl.Recommendation != "" {
				cfg.ConflictRecommendation = cl.Recommendation
			}
			if cl.ActionRequired != "" {
				cfg.ConflictActionRequired = cl.ActionRequired
			}
		}
	}

	if val, ok := l.cache["warning_messages"]; ok {
		var wm struct {
			OverlapExceeded string `json:"overlap_exceeded"`
			MultiLayer      string `json:"multi_layer"`
			Conflict        string `json:"conflict"`
		}
		if err := json.Unmarshal(val, &wm); err == nil {
			if wm.OverlapExceeded != "" {
				cfg.WarningOverlapExceeded = wm.OverlapExceeded
			}
			if wm.MultiLayer != "" {
				cfg.WarningMultiLayer = wm.MultiLayer
			}
			if wm.Conflict != "" {
				cfg.WarningConflict = wm.Conflict
			}
		}
	}

	if val, ok := l.cache["assess_messages"]; ok {
		var am struct {
			NoPlanningImpact    string `json:"no_planning_impact"`
			NoSignificantImpact string `json:"no_significant_impact"`
		}
		if err := json.Unmarshal(val, &am); err == nil {
			if am.NoPlanningImpact != "" {
				cfg.NoPlanningImpactMsg = am.NoPlanningImpact
			}
			if am.NoSignificantImpact != "" {
				cfg.NoSignificantImpactMsg = am.NoSignificantImpact
			}
		}
	}

	if val, ok := l.cache["status_labels"]; ok {
		var sl struct {
			Prohibited     string `json:"prohibited"`
			Conditional    string `json:"conditional"`
			Allowed        string `json:"allowed"`
			ConflictSuffix string `json:"conflict_suffix"`
		}
		if err := json.Unmarshal(val, &sl); err == nil {
			if sl.Prohibited != "" {
				cfg.StatusProhibited = sl.Prohibited
			}
			if sl.Conditional != "" {
				cfg.StatusConditional = sl.Conditional
			}
			if sl.Allowed != "" {
				cfg.StatusAllowed = sl.Allowed
			}
			if sl.ConflictSuffix != "" {
				cfg.ConflictStatusSuffix = sl.ConflictSuffix
			}
		}
	}

	if val, ok := l.cache["land_type_names"]; ok {
		var ltn map[uint64]string
		if err := json.Unmarshal(val, &ltn); err == nil {
			cfg.LandTypeNames = ltn
		}
	}

	if val, ok := l.cache["summary_labels"]; ok {
		var sl struct {
			ConflictSuffix string `json:"conflict_suffix"`
		}
		if err := json.Unmarshal(val, &sl); err == nil {
			if sl.ConflictSuffix != "" {
				cfg.ConflictSummarySuffix = sl.ConflictSuffix
			}
		}
	}

	return cfg, nil
}

func (l *ConfigLoader) Reload(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	dbConfigs, err := l.configRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	l.cache = make(map[string]json.RawMessage)
	for k, v := range dbConfigs {
		l.cache[k] = v
	}
	l.lastLoad = time.Now()

	return nil
}
