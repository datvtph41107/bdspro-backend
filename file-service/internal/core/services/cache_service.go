package core_services

import (
	"context"
	"encoding/json"
	"file/internal/core/domain"
	"file/internal/ports"
	"fmt"
	"time"
)

type CacheService struct {
	cacheRepo ports.CacheRepository
}

func NewCacheService(cacheRepo ports.CacheRepository) *CacheService {
	return &CacheService{cacheRepo: cacheRepo}
}

// File-related cache operations
func (s *CacheService) GetFile(ctx context.Context, fileID string) (*domain.File, error) {
	key := fmt.Sprintf("file:%s", fileID)
	var file domain.File

	data, err := s.cacheRepo.Get(ctx, key)
	if err != nil || data == nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

func (s *CacheService) SetFile(ctx context.Context, file *domain.File, ttl time.Duration) error {
	key := fmt.Sprintf("file:%d", file.ID)
	data, err := json.Marshal(file)
	if err != nil {
		return err
	}

	return s.cacheRepo.Set(ctx, key, data, ttl)
}

func (s *CacheService) DeleteFile(ctx context.Context, fileID string) error {
	key := fmt.Sprintf("file:%s", fileID)
	return s.cacheRepo.Delete(ctx, key)
}

// Access control cache
func (s *CacheService) CheckAccess(ctx context.Context, fileID string, userID uint64, permission string) (bool, error) {
	key := fmt.Sprintf("access:%s:%d:%s", fileID, userID, permission)

	exists, err := s.cacheRepo.Exists(ctx, key)
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	// Cache miss - check from repository
	// This would be filled by the main service
	return false, nil
}

func (s *CacheService) SetAccess(ctx context.Context, fileID string, userID uint64, permission string, ttl time.Duration) error {
	key := fmt.Sprintf("access:%s:%d:%s", fileID, userID, permission)
	return s.cacheRepo.Set(ctx, key, []byte("1"), ttl)
}

func (s *CacheService) RevokeAccess(ctx context.Context, fileID string, userID uint64) error {
	patterns := []string{
		fmt.Sprintf("access:%s:%d:read", fileID, userID),
		fmt.Sprintf("access:%s:%d:write", fileID, userID),
		fmt.Sprintf("access:%s:%d:delete", fileID, userID),
	}

	for _, pattern := range patterns {
		s.cacheRepo.Delete(ctx, pattern)
	}

	return nil
}

// Rate limiting
func (s *CacheService) RateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, int, error) {
	current, err := s.cacheRepo.Incr(ctx, key)
	if err != nil {
		return false, 0, err
	}

	if current == 1 {
		// First request, set expiration
		s.cacheRepo.Set(ctx, key, []byte("1"), window)
	}

	return current > int64(limit), int(current), nil
}

// Session management
func (s *CacheService) GetSession(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	data, err := s.cacheRepo.Get(ctx, key)
	if err != nil || data == nil {
		return nil, err
	}
	var session map[string]interface{}
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *CacheService) SetSession(
	ctx context.Context,
	sessionID string,
	data map[string]interface{},
	ttl time.Duration,
) error {
	key := fmt.Sprintf("session:%s", sessionID)
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.cacheRepo.Set(ctx, key, raw, ttl)
}

func (s *CacheService) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.cacheRepo.Delete(ctx, key)
}

// Distributed lock
func (s *CacheService) AcquireLock(ctx context.Context, lockKey string, ttl time.Duration) (bool, error) {
	key := fmt.Sprintf("lock:%s", lockKey)
	return s.cacheRepo.SetNX(ctx, key, []byte("1"), ttl)
}

func (s *CacheService) ReleaseLock(ctx context.Context, lockKey string) error {
	key := fmt.Sprintf("lock:%s", lockKey)
	return s.cacheRepo.Delete(ctx, key)
}

// Cache statistics
func (s *CacheService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	// Get cache hit/miss ratios
	// This would require tracking hits/misses
	return map[string]interface{}{
		"total_keys":  0,
		"memory_used": 0,
	}, nil
}

// Cache warming for frequently accessed files
func (s *CacheService) WarmCache(ctx context.Context, fileIDs []string) error {
	// Pre-load files into cache
	// This would be called during startup or periodically
	// for _, fileID := range fileIDs {
	//  Load file and cache it
	// 	Implementation depends on your repository
	// }

	return nil
}

// Cache cleanup
func (s *CacheService) CleanExpired(ctx context.Context) error {
	// Redis handles TTL automatically
	// For other caches, implement cleanup logic
	return nil
}
