package layer_resolver

import (
	"tqd/internal/interface/repo"
)

// NewResolver - Khởi tạo resolver với database config
func NewResolver(configRepo repo.ILayerResolverConfigRepo) ILayerResolver {
	return NewResolverWithDB(configRepo)
}
