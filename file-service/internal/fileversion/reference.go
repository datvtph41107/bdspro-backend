package fileversion

import (
	"crypto/hmac"
	"errors"
	"fmt"
	"strings"

	_utils "common/utils"
)

const (
	referenceVersion = "v1"
	referenceDomain  = "file-version:"
)

// ReferenceCodec owns the opaque public reference for version artifacts. New
// references authenticate the XOR envelope; legacy p<XOR> references remain
// readable at this single compatibility boundary.
type ReferenceCodec struct {
	xorKey       string
	signatureKey string
}

func NewReferenceCodec(xorKey, signatureKey string) (*ReferenceCodec, error) {
	if strings.TrimSpace(xorKey) == "" {
		return nil, errors.New("version reference XOR key is required")
	}
	if strings.TrimSpace(signatureKey) == "" {
		return nil, errors.New("version reference signature key is required")
	}
	return &ReferenceCodec{xorKey: xorKey, signatureKey: signatureKey}, nil
}

func (c *ReferenceCodec) Encode(storedPath string) (string, error) {
	if c == nil {
		return "", errors.New("version reference codec is unavailable")
	}
	storedPath = strings.TrimSpace(storedPath)
	if storedPath == "" {
		return "", errors.New("version stored path is empty")
	}

	payload := _utils.XorEncode(storedPath, c.xorKey)
	mac := c.signature(payload)
	return referenceVersion + "." + payload + "." + mac, nil
}

func (c *ReferenceCodec) Decode(reference string) (string, error) {
	if c == nil {
		return "", errors.New("version reference codec is unavailable")
	}
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return "", errors.New("version reference is empty")
	}

	if strings.HasPrefix(reference, referenceVersion+".") {
		parts := strings.Split(reference, ".")
		if len(parts) != 3 || parts[0] != referenceVersion || parts[1] == "" || parts[2] == "" {
			return "", errors.New("version reference is malformed")
		}
		expected := c.signature(parts[1])
		if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
			return "", errors.New("version reference authentication failed")
		}
		decoded, err := _utils.XorDecode(parts[1], c.xorKey)
		if err != nil {
			return "", fmt.Errorf("decode version reference: %w", err)
		}
		if strings.TrimSpace(decoded) == "" {
			return "", errors.New("version reference path is empty")
		}
		return decoded, nil
	}

	// Compatibility only: pre-closure version uploads emitted p<XOR(path)>.
	if strings.HasPrefix(reference, "p") && len(reference) > 1 {
		decoded, err := _utils.XorDecode(reference[1:], c.xorKey)
		if err != nil {
			return "", fmt.Errorf("decode legacy version reference: %w", err)
		}
		if strings.TrimSpace(decoded) == "" {
			return "", errors.New("legacy version reference path is empty")
		}
		return decoded, nil
	}

	return "", errors.New("unsupported version reference")
}

func (c *ReferenceCodec) signature(payload string) string {
	return _utils.GenerateHmacSHA256(
		referenceDomain+referenceVersion+":"+payload,
		c.signatureKey,
	)
}
