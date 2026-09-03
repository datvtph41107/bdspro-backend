package redisquota

import (
	commonmetering "common/metering"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"tqd/internal/access"
)

const (
	reservationPrefix = "quota:reservation:"
	expiredSetKey     = "quota:reservation:expires"
)

func usageKey(
	subject access.Subject,
	meterCode commonmetering.Code,
	periodStart time.Time,
	periodEnd time.Time,
) string {
	return fmt.Sprintf(
		"quota:usage:%s:%s:%s:%d:%d",
		subject.Type,
		subject.ID,
		meterCode,
		periodStart.Unix(),
		periodEnd.Unix(),
	)
}

func commandKey(usageKey, command string) string {
	sum := sha256.Sum256([]byte(usageKey + "\n" + command))
	return "quota:command:" + hex.EncodeToString(sum[:])
}

func reservationKey(id string) string {
	return reservationPrefix + id
}

func reservationID(key string) string {
	if len(key) <= len(reservationPrefix) || key[:len(reservationPrefix)] != reservationPrefix {
		return ""
	}
	return key[len(reservationPrefix):]
}
