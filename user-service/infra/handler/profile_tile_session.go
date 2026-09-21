package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	_tilesession "common/tilesession"
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
	if s.TileSessions == nil {
		return 0, "", 0, fmt.Errorf("tile session store not configured")
	}

	sessionK, err = randomSessionK()
	if err != nil {
		return 0, "", 0, err
	}
	aesKey := make([]byte, 32)
	if _, err = rand.Read(aesKey); err != nil {
		return 0, "", 0, fmt.Errorf("generate aes key: %w", err)
	}

	sessionEncryptKey = base64.StdEncoding.EncodeToString(aesKey)
	expiresAt := time.Now().Add(_tilesession.TTL)

	if err = s.TileSessions.Save(ctx, sessionK, sessionEncryptKey, expiresAt); err != nil {
		return 0, "", 0, err
	}

	return sessionK, sessionEncryptKey, int64(time.Until(expiresAt).Seconds()), nil
}
