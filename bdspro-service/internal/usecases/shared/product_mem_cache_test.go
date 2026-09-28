package shared_usecase

import (
	"testing"

	"bdspro/internal/dto"
)

func TestUserProductsMemCacheKeyIncludesPageAndSize(t *testing.T) {
	req0 := dto.ProductSearchRequest{}
	req0.Page = 0
	req0.Size = 5
	req0.Text = "a"

	req1 := req0
	req1.Page = 1

	key0 := userProductsMemCacheKeyFromDTO(42, req0)
	key1 := userProductsMemCacheKeyFromDTO(42, req1)
	if key0 == key1 {
		t.Fatalf("expected distinct cache keys for page 0 vs 1, got %#v", key0)
	}

	req1SameSize := req1
	req1SameSize.Size = 10
	keySize := userProductsMemCacheKeyFromDTO(42, req1SameSize)
	if key1 == keySize {
		t.Fatalf("expected distinct cache keys for different size, got %#v", key1)
	}
}
