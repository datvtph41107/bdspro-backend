package shared_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"time"

	"github.com/jinzhu/copier"
)

type userProductsMemCacheKey struct {
	userID   uint64
	text     string
	archived string
}

type userProductsMemCacheEntry struct {
	list      []dto.ProductListResponse
	total     int64
	expiresAt time.Time
}

type productDetailMemCacheKey struct {
	profileID uint64
	productID uint64
}

type productDetailMemCacheEntry struct {
	product   *domain.Product
	expiresAt time.Time
}

func archivedKeyForUserProducts(archived *bool) string {
	if archived == nil {
		return "nil"
	}
	if *archived {
		return "true"
	}
	return "false"
}

func userProductsMemCacheKeyFromDTO(userID uint64, req dto.ProductSearchRequest) userProductsMemCacheKey {
	return userProductsMemCacheKey{
		userID:   userID,
		text:     req.Text,
		archived: archivedKeyForUserProducts(req.Archived),
	}
}

func (s *ProductUsecase) productMemCacheTTL() time.Duration {
	if s.CacheTTLProductDetail > 0 {
		return s.CacheTTLProductDetail
	}
	return 30 * time.Minute
}

func (s *ProductUsecase) getCachedUserProducts(userID uint64, req dto.ProductSearchRequest) ([]dto.ProductListResponse, int64, bool) {
	key := userProductsMemCacheKeyFromDTO(userID, req)
	now := time.Now()

	s.productMemCacheMu.RLock()
	entry, ok := s.userProductsByKey[key]
	s.productMemCacheMu.RUnlock()
	if !ok || now.After(entry.expiresAt) {
		return nil, 0, false
	}

	out := make([]dto.ProductListResponse, len(entry.list))
	copy(out, entry.list)
	return out, entry.total, true
}

func (s *ProductUsecase) setCachedUserProducts(userID uint64, req dto.ProductSearchRequest, list []dto.ProductListResponse, total int64) {
	key := userProductsMemCacheKeyFromDTO(userID, req)
	stored := make([]dto.ProductListResponse, len(list))
	copy(stored, list)

	s.productMemCacheMu.Lock()
	s.userProductsByKey[key] = userProductsMemCacheEntry{
		list:      stored,
		total:     total,
		expiresAt: time.Now().Add(s.productMemCacheTTL()),
	}
	s.productMemCacheMu.Unlock()
}

func (s *ProductUsecase) getCachedProductDetail(profileID, productID uint64) (*domain.Product, bool) {
	key := productDetailMemCacheKey{profileID: profileID, productID: productID}
	now := time.Now()

	s.productMemCacheMu.RLock()
	entry, ok := s.productDetailByKey[key]
	s.productMemCacheMu.RUnlock()
	if !ok || now.After(entry.expiresAt) || entry.product == nil {
		return nil, false
	}

	var out domain.Product
	if err := copier.Copy(&out, entry.product); err != nil {
		return nil, false
	}
	return &out, true
}

func (s *ProductUsecase) setCachedProductDetail(profileID, productID uint64, entity *domain.Product) {
	if entity == nil {
		return
	}
	var stored domain.Product
	if err := copier.Copy(&stored, entity); err != nil {
		return
	}

	key := productDetailMemCacheKey{profileID: profileID, productID: productID}
	s.productMemCacheMu.Lock()
	s.productDetailByKey[key] = productDetailMemCacheEntry{
		product:   &stored,
		expiresAt: time.Now().Add(s.productMemCacheTTL()),
	}
	s.productMemCacheMu.Unlock()
}

func (s *ProductUsecase) invalidateProductMemCache(ownerUserID, productID uint64) {
	s.productMemCacheMu.Lock()
	defer s.productMemCacheMu.Unlock()

	for k := range s.productDetailByKey {
		if k.productID == productID {
			delete(s.productDetailByKey, k)
		}
	}
	for k := range s.userProductsByKey {
		if k.userID == ownerUserID {
			delete(s.userProductsByKey, k)
		}
	}
}
