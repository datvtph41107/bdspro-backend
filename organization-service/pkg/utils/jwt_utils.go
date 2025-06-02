package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func GetUserFromContext(ctx context.Context, key string) any {
	if ctx == nil {
		return nil
	}
	claims, ok := ctx.Value(key).(jwt.MapClaims)
	if !ok {
		return nil
	}

	return claims
}

func GetUserID(ctx context.Context, key string) uint32 {
	claims := GetUserFromContext(ctx, key)
	if claims == nil {
		return 0
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return 0
	}

	id, ok := mapClaims["id"].(float64)
	if !ok {
		return 0
	}

	return uint32(id)
}

func DecodeJWTPayload(tokenStr string) (map[string]interface{}, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid token")
	}

	payloadSegment := parts[1]
	payloadSegment += strings.Repeat("=", (4-len(payloadSegment)%4)%4)

	decoded, err := base64.URLEncoding.DecodeString(payloadSegment)
	if err != nil {
		return nil, err
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, err
	}

	return claims, nil
}
