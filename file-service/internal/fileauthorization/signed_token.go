package fileauthorization

import (
	"crypto/hmac"
	"fmt"
	"strconv"
	"strings"
	"time"

	_utils "common/utils"
)

// SignedReadToken is the verified identity/expiry carried by the existing
// private-file signed-read envelope.
//
// The signature intentionally preserves the current wire contract:
// HMAC("info:<profileID>:<expiry>", SignatureKey), then XOR envelope encoding.
type SignedReadToken struct {
	ProfileID int64
	ExpiresAt int64
}

// VerifySignedReadToken authenticates the existing private-file signed token.
//
// Malformed, expired, or tampered tokens are normal authorization denials and
// therefore return ok=false rather than an operational error.
func VerifySignedReadToken(
	encodedToken string,
	xorKey string,
	signatureKey string,
	now time.Time,
) (token SignedReadToken, ok bool) {
	if encodedToken == "" || xorKey == "" || signatureKey == "" {
		return SignedReadToken{}, false
	}

	decoded, err := _utils.XorDecode(encodedToken, xorKey)
	if err != nil {
		return SignedReadToken{}, false
	}

	parts := strings.Split(decoded, ":")
	if len(parts) != 3 {
		return SignedReadToken{}, false
	}

	profileID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return SignedReadToken{}, false
	}

	expiresAt, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return SignedReadToken{}, false
	}

	// Preserve the existing boundary semantics: expiry is invalid only after
	// the encoded second has passed.
	if now.Unix() > expiresAt {
		return SignedReadToken{}, false
	}

	expected := _utils.GenerateSignedUrl(
		"info",
		profileID,
		signatureKey,
		expiresAt,
	)

	if !hmac.Equal([]byte(parts[0]), []byte(expected)) {
		return SignedReadToken{}, false
	}

	return SignedReadToken{
		ProfileID: profileID,
		ExpiresAt: expiresAt,
	}, true
}

// IssuedSignedReadToken is the authenticated signed-read envelope returned to
// the HTTP boundary.
//
// Encoded preserves the existing wire format consumed by File clients.
type IssuedSignedReadToken struct {
	Encoded   string
	ProfileID int64
	ExpiresAt int64
}

// SignedReadIssuer owns creation of the same signed-read token contract that
// VerifySignedReadToken authenticates.
//
// Token semantics intentionally remain:
//
//	HMAC("info:<profileID>:<expiry>", SignatureKey)
//	    ↓
//	"<hmac>:<profileID>:<expiry>"
//	    ↓
//	XOR(XorCryptKey)
type SignedReadIssuer struct {
	xorKey       string
	signatureKey string
	ttl          time.Duration
}

func NewSignedReadIssuer(
	xorKey string,
	signatureKey string,
	ttl time.Duration,
) *SignedReadIssuer {
	return &SignedReadIssuer{
		xorKey:       xorKey,
		signatureKey: signatureKey,
		ttl:          ttl,
	}
}

func (i *SignedReadIssuer) Issue(
	profileID int64,
	now time.Time,
) IssuedSignedReadToken {
	expiresAt := now.Add(i.ttl).Unix()

	signature := _utils.GenerateSignedUrl(
		"info",
		profileID,
		i.signatureKey,
		expiresAt,
	)

	envelope := fmt.Sprintf(
		"%s:%d:%d",
		signature,
		profileID,
		expiresAt,
	)

	return IssuedSignedReadToken{
		Encoded:   _utils.XorEncode(envelope, i.xorKey),
		ProfileID: profileID,
		ExpiresAt: expiresAt,
	}
}
