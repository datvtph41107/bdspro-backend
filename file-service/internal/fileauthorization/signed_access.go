package fileauthorization

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// AccessLookup is the narrow durable evidence required by ordinary
// private-file access. It deliberately exposes read capability only.
type AccessLookup interface {
	ExistsAccess(fileID int64, accessID int64) (bool, error)
}

// SignedAccess owns the ordinary signed private-file access mechanism:
//
//	valid signed token
//	    +
//	durable file_access evidence
//
// Global IAM fallback is intentionally not part of this type; that policy
// remains owned by AuthorizePrivateRead.
type SignedAccess struct {
	access       AccessLookup
	xorKey       string
	signatureKey string
}

func NewSignedAccess(
	access AccessLookup,
	xorKey string,
	signatureKey string,
) *SignedAccess {
	return &SignedAccess{
		access:       access,
		xorKey:       xorKey,
		signatureKey: signatureKey,
	}
}

func (s *SignedAccess) Allows(
	decodedPath string,
	encodedToken string,
	now time.Time,
) (bool, error) {
	if s == nil || s.access == nil {
		return false, nil
	}

	token, ok := VerifySignedReadToken(
		encodedToken,
		s.xorKey,
		s.signatureKey,
		now,
	)
	if !ok {
		return false, nil
	}

	fileID, ok := signedAccessFileID(decodedPath)
	if !ok {
		return false, nil
	}

	return s.access.ExistsAccess(
		fileID,
		token.ProfileID,
	)
}

// signedAccessFileID extracts the durable FileEntity ID from the canonical
// stored-file basename (<id>.<ext>). Authorization owns this interpretation
// locally rather than depending on a generic shared path helper.
func signedAccessFileID(path string) (int64, bool) {
	base := filepath.Base(filepath.Clean(path))
	ext := filepath.Ext(base)
	if ext == "" {
		return 0, false
	}

	name := strings.TrimSuffix(base, ext)
	fileID, err := strconv.ParseInt(name, 10, 64)
	if err != nil || fileID <= 0 {
		return 0, false
	}

	return fileID, true
}
