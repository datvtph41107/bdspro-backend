package redisquota

import (
	commonaccess "common/access"
	"common/businessoperation"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	reservationPrefix = "quota:reservation:"
	expiredSetKey     = "quota:reservation:expires"
)

func usageKey(
	subject commonaccess.Subject,
	operation businessoperation.Code,
	periodStart time.Time,
	periodEnd time.Time,
) string {
	return fmt.Sprintf(
		"quota:usage:%s:%s:%s:%d:%d",
		subject.Type,
		subject.ID,
		operation,
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
