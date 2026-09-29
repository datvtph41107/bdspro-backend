package cache

import (
	"fmt"
	"time"
)

// ==================== KEY PATTERNS ====================

func UserSyncTimeKey(profileID uint64) string {
	return fmt.Sprintf("sync:user:%d:time", profileID)
}

func UserSyncProductsKey(profileID uint64) string {
	return fmt.Sprintf("sync:user:%d:products", profileID)
}

func UserProductsChangeKey(profileID uint64) string {
	return fmt.Sprintf("sync:user:%d:changes", profileID)
}

// 🔥 tombstone (nếu tách riêng)
func UserRemovedKey(profileID uint64) string {
	return fmt.Sprintf("sync:user:%d:removed", profileID)
}

func ProductKey(productID uint64) string {
	return fmt.Sprintf("product:%d", productID)
}

func ProductUpdatedKey(productID uint64) string {
	return fmt.Sprintf("product:%d:updated", productID)
}

// ==================== TTL CONSTANTS ====================

const (
	// ⚠️ cân nhắc KHÔNG expire change log
	TTLUserSync = 24 * time.Hour

	// product cache
	TTLProduct = 5 * time.Minute

	// distributed lock
	TTLLock = 10 * time.Second
)
