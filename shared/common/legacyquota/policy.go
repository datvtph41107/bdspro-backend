package legacyquota

import (
	"context"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type PoolConfig struct {
	DefaultLimit int64    `yaml:"default_limit"`
	Services     []string `yaml:"services"`
}

type PolicyConfig struct {
	Pools   map[string]PoolConfig `yaml:"pools"`
	Actions map[string]int64      `yaml:"actions"`
}

type DefaultPolicyEngine struct {
	mu       sync.RWMutex
	pools    map[string]PoolConfig
	services map[string]string // service -> pool
	actions  map[string]int64
}

func NewDefaultPolicyEngine(configPath string) (*DefaultPolicyEngine, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy file %s: %v", configPath, err)
	}

	var cfg PolicyConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse policy file: %v", err)
	}

	engine := &DefaultPolicyEngine{
		pools:    cfg.Pools,
		services: make(map[string]string),
		actions:  cfg.Actions,
	}

	for poolName, poolCfg := range cfg.Pools {
		for _, svc := range poolCfg.Services {
			engine.services[svc] = poolName
		}
	}

	return engine, nil
}

func (e *DefaultPolicyEngine) Resolve(service, action string) (pool string, cost int64, err error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	pool, ok := e.services[service]
	if !ok {
		return "", 0, fmt.Errorf("unknown service: %s", service)
	}

	key := service + "." + action
	cost, ok = e.actions[key]
	if !ok {
		cost = -1
	}
	return pool, cost, nil
}

func (e *DefaultPolicyEngine) GetUserLimit(ctx context.Context, userID, pool string) (int64, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if cfg, ok := e.pools[pool]; ok {
		return cfg.DefaultLimit, nil
	}
	return 0, fmt.Errorf("unknown pool: %s", pool)
}
