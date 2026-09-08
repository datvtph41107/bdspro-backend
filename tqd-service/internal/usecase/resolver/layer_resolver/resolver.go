package layer_resolver

import (
	"context"
	"sync"

	"tqd/internal/interface/repo"
	"tqd/internal/usecase/resolver/layer_resolver/db"
	"tqd/internal/usecase/resolver/layer_resolver/engine"
	"tqd/internal/usecase/resolver/layer_resolver/types"
)

type ILayerResolver interface {
	Resolve(ctx context.Context, candidates []*types.LayerCandidate, mode string) *types.ResolveResult
	ReloadConfig(ctx context.Context) error
}

type Resolver struct {
	engine   *engine.ResolverEngine
	dbLoader *db.ConfigLoader
	mu       sync.RWMutex
}

func NewResolverWithDB(configRepo repo.IQHLayerResolverConfigRepo) *Resolver {
	return &Resolver{
		dbLoader: db.NewConfigLoader(configRepo),
	}
}

func (r *Resolver) getEngine(ctx context.Context) *engine.ResolverEngine {
	r.mu.RLock()
	if r.engine != nil {
		defer r.mu.RUnlock()
		return r.engine
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	config, err := r.dbLoader.LoadConfig(ctx)
	if err != nil {
		config = types.DefaultConfig()
	}
	r.engine = engine.NewResolverEngine(config)

	return r.engine
}

func (r *Resolver) Resolve(ctx context.Context, candidates []*types.LayerCandidate, mode string) *types.ResolveResult {
	eng := r.getEngine(ctx)
	resolveMode := types.ResolveMode(mode)
	if !resolveMode.IsValid() {
		resolveMode = types.ModeDetail
	}
	return eng.Resolve(candidates, resolveMode)
}

func (r *Resolver) ReloadConfig(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.dbLoader.Reload(ctx); err != nil {
		return err
	}

	config, err := r.dbLoader.LoadConfig(ctx)
	if err != nil {
		return err
	}
	r.engine = engine.NewResolverEngine(config)

	return nil
}
