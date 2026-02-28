package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// XOR encryption for file paths
func XOREncrypt(text, key string) string {
	var result []byte
	keyBytes := []byte(key)

	for i := 0; i < len(text); i++ {
		result = append(result, text[i]^keyBytes[i%len(keyBytes)])
	}

	return base64.URLEncoding.EncodeToString(result)
}

func XORDecrypt(encrypted, key string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	var result []byte
	keyBytes := []byte(key)

	for i := 0; i < len(data); i++ {
		result = append(result, data[i]^keyBytes[i%len(keyBytes)])
	}

	return string(result), nil
}

// Generate HMAC signature
func GenerateHMAC(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// Verify HMAC signature
func VerifyHMAC(data, signature, key string) bool {
	expected := GenerateHMAC(data, key)
	return hmac.Equal([]byte(signature), []byte(expected))
}

// Generate random string
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}

// Generate UUID v4
func GenerateUUID() string {
	u := make([]byte, 16)
	rand.Read(u)

	// Set version (4) and variant (RFC4122)
	u[6] = (u[6] & 0x0f) | 0x40
	u[8] = (u[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

// Generate signed URL
func GenerateSignedURL(path string, userID uint64, key string, expiry int64) string {
	data := fmt.Sprintf("%s:%d:%d", path, userID, expiry)
	signature := GenerateHMAC(data, key)
	return fmt.Sprintf("%s:%s", data, signature)
}

// Parse signed URL
func ParseSignedURL(signedURL, key string) (path string, userID uint64, expiry int64, valid bool) {
	parts := strings.Split(signedURL, ":")
	if len(parts) != 4 {
		return "", 0, 0, false
	}

	path = parts[0]
	fmt.Sscanf(parts[1], "%d", &userID)
	fmt.Sscanf(parts[2], "%d", &expiry)
	signature := parts[3]

	// Verify signature
	expectedData := fmt.Sprintf("%s:%d:%d", path, userID, expiry)
	expectedSignature := GenerateHMAC(expectedData, key)

	return path, userID, expiry, hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// Generate secure token for file access
func GenerateAccessToken(fileID string, userID uint64, permission string, expiry time.Duration) string {
	expiryTime := time.Now().Add(expiry).Unix()
	data := fmt.Sprintf("%s:%d:%s:%d", fileID, userID, permission, expiryTime)

	// Encrypt with XOR then base64
	encrypted := XOREncrypt(data, "access-token-key")

	return base64.URLEncoding.EncodeToString([]byte(encrypted))
}

// Validate access token
func ValidateAccessToken(token, fileID string, userID uint64, permission string) bool {
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	data, err := XORDecrypt(string(decoded), "access-token-key")
	if err != nil {
		return false
	}

	parts := strings.Split(data, ":")
	if len(parts) != 4 {
		return false
	}

	tokenFileID := parts[0]
	tokenUserID := parts[1]
	tokenPermission := parts[2]
	expiryStr := parts[3]

	var expiry int64
	fmt.Sscanf(expiryStr, "%d", &expiry)

	// Check expiry
	if time.Now().Unix() > expiry {
		return false
	}

	// Check data matches
	if tokenFileID != fileID || tokenPermission != permission {
		return false
	}

	var tokenUID uint64
	fmt.Sscanf(tokenUserID, "%d", &tokenUID)

	return tokenUID == userID
}
