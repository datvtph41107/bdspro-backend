package layer_resolver

import (
	repo "userv4/modules/tqd/repository"
)

// NewResolver - Khởi tạo resolver với database config
func NewResolver(configRepo repo.IQHLayerResolverConfigRepo) ILayerResolver {
	return NewResolverWithDB(configRepo)
}
