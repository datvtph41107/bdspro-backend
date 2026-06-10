package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	_redis "common/redis"
)

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

func (s *GrpcProfileService) genTileSession(ctx context.Context) (sessionK uint64, sessionEncryptKey string, expiresInSec int64, err error) {
	if s.Redis == nil {
		return 0, "", 0, fmt.Errorf("redis not configured")
	}

	sessionEncryptKey, err = s.Redis.GetTileSessionEncryptKey(s.sessionKey)
	if err == nil && sessionEncryptKey != "" && s.sessionKey != 0 {
		return s.sessionKey, sessionEncryptKey, int64(_redis.TileSessionTTL.Seconds()), nil
	}

	sessionK, err = randomSessionK()
	if err != nil {
		return 0, "", 0, err
	}
	s.sessionKey = sessionK

	aesKey := make([]byte, 32)
	if _, err = rand.Read(aesKey); err != nil {
		return 0, "", 0, fmt.Errorf("generate aes key: %w", err)
	}

	sessionEncryptKey = base64.StdEncoding.EncodeToString(aesKey)
	expiresAt := time.Now().Add(_redis.TileSessionTTL)

	if err = s.Redis.SaveTileSession(sessionK, sessionEncryptKey, expiresAt); err != nil {
		return 0, "", 0, err
	}

	return sessionK, sessionEncryptKey, int64(time.Until(expiresAt).Seconds()), nil
}
