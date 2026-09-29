package utils

import (
	"context"

	"relay/constants"
)

func GetCurrentUserID(ctx context.Context) uint64 {
	if ctx == nil {
		return 0
	}
	claims := ctx.Value(constants.CONTEXT_USER_KEY)
	ok := claims != nil
	if !ok {
		return 0
	}

	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		return 0
	}
	if idFloat, ok := claimsMap["profile"].(float64); ok {

		id := uint64(idFloat)
		return id
	}

	return 0
}

func GetCurrentUser(ctx context.Context) any {
	claims, ok := ctx.Value(constants.CONTEXT_USER_KEY).(any)
	if !ok {
		return nil
	}

	return claims
}
