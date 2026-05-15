// internal/modules/resolver/layer_resolver/db/loader.go
package db

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"tqd/internal/interface/repo"
	"tqd/internal/modules/resolver/layer_resolver/config"
	"tqd/internal/modules/resolver/layer_resolver/types"
)

type ConfigLoader struct {
	configRepo repo.ILayerResolverConfigRepo
	cache      map[string]json.RawMessage
	mu         sync.RWMutex
	lastLoad   time.Time
	ttl        time.Duration
}

func NewConfigLoader(configRepo repo.ILayerResolverConfigRepo) *ConfigLoader {
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
		return l.buildConfig(), nil
	}
	l.mu.RUnlock()

	l.mu.Lock()
	defer l.mu.Unlock()

	dbConfigs, err := l.configRepo.GetAll(ctx)
	if err != nil {
		// Fallback to YAML config
		return config.GetConfig().ToTypesConfig(), nil
	}

	for k, v := range dbConfigs {
		l.cache[k] = v
	}
	l.lastLoad = time.Now()

	return l.buildConfig(), nil
}

func (l *ConfigLoader) buildConfig() *types.Config {
	yamlConfig := config.GetConfig()
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

	return cfg
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
