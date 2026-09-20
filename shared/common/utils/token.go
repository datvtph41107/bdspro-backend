package _utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetJwtFromRequest(r *http.Request, headerKey, prefixToken string) string {
	bearerToken := r.Header.Get(headerKey)
	if bearerToken != "" && strings.HasPrefix(bearerToken, prefixToken) {
		return bearerToken[len(prefixToken):]
	}
	return ""
}

func ValidTokenByIssueAt(issuedAt time.Time, timestampStr string) bool {
	if timestampStr == "" {
		return true
	}

	// Parse Unix timestamp từ string
	validFromUnix, err := strconv.ParseInt(timestampStr, 10, 64)
	slog.Info(strings.TrimSuffix(fmt.Sprintln("validFromUnix", validFromUnix, err), "\n"))
	if err != nil {
		// Fallback: thử parse RFC3339 format cũ để tương thích ngược
		validDateFrom, err := time.Parse(time.RFC3339, timestampStr)
		slog.Info(strings.TrimSuffix(fmt.Sprintln("validDateFrom", validDateFrom, issuedAt), "\n"))
		if err != nil {
			return false
		}
		return !issuedAt.After(validDateFrom)
	}

	// Convert cả hai về Unix timestamp (UTC) để so sánh
	validFromTime := time.Unix(validFromUnix, 0).UTC()
	issuedAtUTC := issuedAt.UTC()
	slog.Info(strings.TrimSuffix(fmt.Sprintln("validFromTime", issuedAtUTC, validFromTime), "\n"))
	slog.Info(strings.TrimSuffix(fmt.Sprintln("issuedAtUTC.Before(validFromTime)", !issuedAtUTC.Before(validFromTime)), "\n"))

	// Token hợp lệ nếu IssuedAt >= ValidFrom
	// Nghĩa là token được issue sau thời điểm validFrom
	return !issuedAtUTC.Before(validFromTime)
}

func GenerateSignedUrl(info string, profileId int64, secretKey string, expiration int64) string {
	data := info + ":" + strconv.FormatInt(profileId, 10) + ":" + strconv.FormatInt(expiration, 10)
	return GenerateHmacSHA256(data, secretKey)
}

func GenerateHmacSHA256(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func XorEncode(text, key string) string {
	encryptedBytes := xorWithKey([]byte(text), []byte(key))
	return base64.RawURLEncoding.EncodeToString(encryptedBytes)
}

func XorDecode(encodedText, key string) (string, error) {
	decodedBytes, err := base64.RawURLEncoding.DecodeString(encodedText)
	if err != nil {
		return "", err
	}
	decryptedBytes := xorWithKey(decodedBytes, []byte(key))
	return string(decryptedBytes), nil
}

func xorWithKey(data, key []byte) []byte {
	result := make([]byte, len(data))
	for i := range data {
		result[i] = data[i] ^ key[i%len(key)]
	}
	return result
}

func GetFileIdFromPath(path string) (int64, error) {
	ext := GetFileExtension(path)
	pattern := regexp.MustCompile(`(\\d+)\\.` + ext + `$`)
	matches := pattern.FindStringSubmatch(path)
	if len(matches) > 1 {
		return strconv.ParseInt(matches[1], 10, 64)
	}
	return 0, errors.New("file ID not found in path")
}

func GetFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}
