package fileauthorization

import (
	"fmt"
	"testing"
	"time"

	_utils "common/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSignedReadXORKey       = "xor-test-key"
	testSignedReadSignatureKey = "signature-test-key"
)

func encodeSignedReadToken(
	profileID int64,
	expiresAt int64,
	signature string,
) string {
	envelope := fmt.Sprintf(
		"%s:%d:%d",
		signature,
		profileID,
		expiresAt,
	)
	return _utils.XorEncode(envelope, testSignedReadXORKey)
}

func validSignedReadToken(
	profileID int64,
	expiresAt int64,
) string {
	signature := _utils.GenerateSignedUrl(
		"info",
		profileID,
		testSignedReadSignatureKey,
		expiresAt,
	)
	return encodeSignedReadToken(
		profileID,
		expiresAt,
		signature,
	)
}

func TestVerifySignedReadTokenAcceptsCurrentWireContract(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	expiresAt := now.Add(time.Minute).Unix()

	encoded := validSignedReadToken(42, expiresAt)

	token, ok := VerifySignedReadToken(
		encoded,
		testSignedReadXORKey,
		testSignedReadSignatureKey,
		now,
	)

	require.True(t, ok)
	assert.Equal(t, int64(42), token.ProfileID)
	assert.Equal(t, expiresAt, token.ExpiresAt)
}

func TestVerifySignedReadTokenRejectsTamperedSignature(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	expiresAt := now.Add(time.Minute).Unix()

	encoded := encodeSignedReadToken(
		42,
		expiresAt,
		"forged-signature",
	)

	_, ok := VerifySignedReadToken(
		encoded,
		testSignedReadXORKey,
		testSignedReadSignatureKey,
		now,
	)

	assert.False(t, ok)
}

func TestVerifySignedReadTokenRejectsTamperedProfile(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	expiresAt := now.Add(time.Minute).Unix()

	signatureForProfile42 := _utils.GenerateSignedUrl(
		"info",
		42,
		testSignedReadSignatureKey,
		expiresAt,
	)

	// Keep the valid signature for profile 42 but alter the envelope identity.
	encoded := encodeSignedReadToken(
		43,
		expiresAt,
		signatureForProfile42,
	)

	_, ok := VerifySignedReadToken(
		encoded,
		testSignedReadXORKey,
		testSignedReadSignatureKey,
		now,
	)

	assert.False(t, ok)
}

func TestVerifySignedReadTokenRejectsTamperedExpiry(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	originalExpiry := now.Add(time.Minute).Unix()

	signatureForOriginalExpiry := _utils.GenerateSignedUrl(
		"info",
		42,
		testSignedReadSignatureKey,
		originalExpiry,
	)

	encoded := encodeSignedReadToken(
		42,
		originalExpiry+3600,
		signatureForOriginalExpiry,
	)

	_, ok := VerifySignedReadToken(
		encoded,
		testSignedReadXORKey,
		testSignedReadSignatureKey,
		now,
	)

	assert.False(t, ok)
}

func TestVerifySignedReadTokenRejectsExpiredToken(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)

	encoded := validSignedReadToken(
		42,
		now.Add(-time.Second).Unix(),
	)

	_, ok := VerifySignedReadToken(
		encoded,
		testSignedReadXORKey,
		testSignedReadSignatureKey,
		now,
	)

	assert.False(t, ok)
}

func TestVerifySignedReadTokenRejectsMalformedEnvelope(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)

	encoded := _utils.XorEncode(
		"not:a:valid:envelope",
		testSignedReadXORKey,
	)

	_, ok := VerifySignedReadToken(
		encoded,
		testSignedReadXORKey,
		testSignedReadSignatureKey,
		now,
	)

	assert.False(t, ok)
}

func TestVerifySignedReadTokenRejectsInvalidEncoding(t *testing.T) {
	_, ok := VerifySignedReadToken(
		"%%%not-base64%%%",
		testSignedReadXORKey,
		testSignedReadSignatureKey,
		time.Unix(1_700_000_000, 0),
	)

	assert.False(t, ok)
}

func TestVerifySignedReadTokenRejectsMissingKeys(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	expiresAt := now.Add(time.Minute).Unix()
	encoded := validSignedReadToken(42, expiresAt)

	_, ok := VerifySignedReadToken(
		encoded,
		"",
		testSignedReadSignatureKey,
		now,
	)
	assert.False(t, ok)

	_, ok = VerifySignedReadToken(
		encoded,
		testSignedReadXORKey,
		"",
		now,
	)
	assert.False(t, ok)
}

func TestSignedReadIssuerPreservesCurrentWireContract(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	ttl := 5 * time.Minute

	issuer := NewSignedReadIssuer(
		testSignedReadXORKey,
		testSignedReadSignatureKey,
		ttl,
	)

	issued := issuer.Issue(42, now)
	expiresAt := now.Add(ttl).Unix()

	expectedSignature := _utils.GenerateSignedUrl(
		"info",
		42,
		testSignedReadSignatureKey,
		expiresAt,
	)

	expectedEncoded := encodeSignedReadToken(
		42,
		expiresAt,
		expectedSignature,
	)

	assert.Equal(t, expectedEncoded, issued.Encoded)
	assert.Equal(t, int64(42), issued.ProfileID)
	assert.Equal(t, expiresAt, issued.ExpiresAt)
}

func TestSignedReadIssuerRoundTripsThroughVerifier(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)

	issuer := NewSignedReadIssuer(
		testSignedReadXORKey,
		testSignedReadSignatureKey,
		5*time.Minute,
	)

	issued := issuer.Issue(42, now)

	verified, ok := VerifySignedReadToken(
		issued.Encoded,
		testSignedReadXORKey,
		testSignedReadSignatureKey,
		now,
	)

	require.True(t, ok)
	assert.Equal(t, issued.ProfileID, verified.ProfileID)
	assert.Equal(t, issued.ExpiresAt, verified.ExpiresAt)
}
