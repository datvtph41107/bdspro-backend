package cmd

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"path"
	"strconv"
	"strings"

	_crypto "common/pkg/crypto"
	_redis "common/redis"

	"github.com/redis/go-redis/v9"
)

const (
	tileEncryptedHeader    = "X-Content-Encrypted"
	tileEncryptedHeaderVal = "aes-256-ctr"
	tileSessionIDQuery     = "sId"
)

func tilesetNameFromPath(urlPath string) string {
	urlPath = strings.TrimPrefix(urlPath, "/")
	if urlPath == "" {
		return ""
	}
	seg := urlPath
	if i := strings.IndexByte(urlPath, '/'); i >= 0 {
		seg = urlPath[:i]
	}
	seg = strings.TrimSuffix(seg, ".pmtiles")
	if ext := path.Ext(seg); ext != "" {
		seg = strings.TrimSuffix(seg, ext)
	}
	return seg
}

func isCIEncryptedTileset(tileset string) bool {
	return strings.HasPrefix(tileset, "ci_")
}

func sIdFromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}
	// if v := strings.TrimSpace(r.Header.Get(tileSessionIDHeader)); v != "" {
	// 	return v
	// }
	return strings.TrimSpace(r.URL.Query().Get(tileSessionIDQuery))
}

// tileEncryptorFromSessionID lấy sessionEncryptKey từ Redis ss:k:{sessionK} theo sId (sessionK).
func tileEncryptorFromSessionID(ctx context.Context, redisSvc *_redis.RedisService, sId string) (*_crypto.AESEncryptor, error) {
	if redisSvc == nil {
		return nil, errors.New("redis not configured")
	}
	if sId == "" {
		return nil, redis.Nil
	}
	sessionK, err := strconv.ParseUint(sId, 10, 64)
	if err != nil || sessionK == 0 {
		return nil, redis.Nil
	}

	sessionEncryptKey, err := redisSvc.GetTileSessionEncryptKey(sessionK)
	if errors.Is(err, redis.Nil) {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}

	aesKey, err := base64.StdEncoding.DecodeString(sessionEncryptKey)
	if err != nil {
		return nil, err
	}
	return _crypto.NewAESEncryptorFromKey(aesKey)
}

func encryptTileBody(ctx context.Context, enc *_crypto.AESEncryptor, body []byte) ([]byte, error) {
	if len(body) == 0 {
		return body, nil
	}
	return enc.EncryptBytesCTR(ctx, body)
}

func stripHeadersAfterEncrypt(headers map[string]string) {
	delete(headers, "Content-Length")
	delete(headers, "Content-Encoding")
	delete(headers, "ETag")
}

func writeTileUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"sId invalid or session expired"}`))
}
