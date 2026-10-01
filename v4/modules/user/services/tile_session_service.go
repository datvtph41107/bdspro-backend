package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	_redis "common/redis"
)

// TileSessionResult is the payload for POST /v2/user/profile/session.
type TileSessionResult struct {
	SessionK          string
	SessionEncryptKey string
	ExpiresIn         int64
}

// TileSessionService owns map-tile sessionK + AES key generation (Redis TTL).
type TileSessionService interface {
	GenTileSession(ctx context.Context) (*TileSessionResult, error)
}

// LocalTileSessionService stores short-lived tile decrypt keys in Redis.
type LocalTileSessionService struct {
	redis *_redis.RedisService

	mu         sync.Mutex
	sessionKey uint64
}

func NewLocalTileSessionService(redis *_redis.RedisService) *LocalTileSessionService {
	return &LocalTileSessionService{redis: redis}
}

func (s *LocalTileSessionService) GenTileSession(ctx context.Context) (*TileSessionResult, error) {
	if s == nil || s.redis == nil {
		return nil, fmt.Errorf("redis not configured")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if key, err := s.redis.GetTileSessionEncryptKey(s.sessionKey); err == nil && key != "" && s.sessionKey != 0 {
		return &TileSessionResult{
			SessionK:          fmt.Sprintf("%d", s.sessionKey),
			SessionEncryptKey: key,
			ExpiresIn:         int64(_redis.TileSessionTTL.Seconds()),
		}, nil
	}

	sessionK, err := randomSessionK()
	if err != nil {
		return nil, err
	}
	s.sessionKey = sessionK

	aesKey := make([]byte, 32)
	if _, err = rand.Read(aesKey); err != nil {
		return nil, fmt.Errorf("generate aes key: %w", err)
	}
	sessionEncryptKey := base64.StdEncoding.EncodeToString(aesKey)
	expiresAt := time.Now().Add(_redis.TileSessionTTL)

	if err = s.redis.SaveTileSession(sessionK, sessionEncryptKey, expiresAt); err != nil {
		return nil, err
	}

	return &TileSessionResult{
		SessionK:          fmt.Sprintf("%d", sessionK),
		SessionEncryptKey: sessionEncryptKey,
		ExpiresIn:         int64(time.Until(expiresAt).Seconds()),
	}, nil
}

func randomSessionK() (uint64, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, fmt.Errorf("generate sessionK: %w", err)
	}
	n := binary.BigEndian.Uint64(b[:])
	if n == 0 {
		n = 1
	}
	return n, nil
}
