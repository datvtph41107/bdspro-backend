package fileauthorization

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeAccessLookup struct {
	exists   bool
	err      error
	calls    int
	fileID   int64
	accessID int64
}

func (f *fakeAccessLookup) ExistsAccess(
	fileID int64,
	accessID int64,
) (bool, error) {
	f.calls++
	f.fileID = fileID
	f.accessID = accessID
	return f.exists, f.err
}

func TestSignedAccessAllowsVerifiedTokenWithDurableEvidence(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	expiresAt := now.Add(time.Minute).Unix()

	lookup := &fakeAccessLookup{exists: true}
	access := NewSignedAccess(
		lookup,
		testSignedReadXORKey,
		testSignedReadSignatureKey,
	)

	allowed, err := access.Allows(
		"files/secure/document/2026-08-25/uuid/42.pdf",
		validSignedReadToken(7, expiresAt),
		now,
	)

	require.NoError(t, err)
	assert.True(t, allowed)
	assert.Equal(t, 1, lookup.calls)
	assert.Equal(t, int64(42), lookup.fileID)
	assert.Equal(t, int64(7), lookup.accessID)
}

func TestSignedAccessRejectsInvalidTokenBeforeAccessLookup(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)

	lookup := &fakeAccessLookup{exists: true}
	access := NewSignedAccess(
		lookup,
		testSignedReadXORKey,
		testSignedReadSignatureKey,
	)

	allowed, err := access.Allows(
		"files/secure/document/2026-08-25/uuid/42.pdf",
		"invalid-token",
		now,
	)

	require.NoError(t, err)
	assert.False(t, allowed)
	assert.Zero(t, lookup.calls)
}

func TestSignedAccessRejectsPathWithoutFileIDBeforeLookup(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	expiresAt := now.Add(time.Minute).Unix()

	lookup := &fakeAccessLookup{exists: true}
	access := NewSignedAccess(
		lookup,
		testSignedReadXORKey,
		testSignedReadSignatureKey,
	)

	allowed, err := access.Allows(
		"files/secure/document/no-file-id.pdf",
		validSignedReadToken(7, expiresAt),
		now,
	)

	require.NoError(t, err)
	assert.False(t, allowed)
	assert.Zero(t, lookup.calls)
}

func TestSignedAccessPreservesLookupFailure(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	expiresAt := now.Add(time.Minute).Unix()
	lookupErr := errors.New("postgres unavailable")

	lookup := &fakeAccessLookup{err: lookupErr}
	access := NewSignedAccess(
		lookup,
		testSignedReadXORKey,
		testSignedReadSignatureKey,
	)

	allowed, err := access.Allows(
		"files/secure/document/2026-08-25/uuid/42.pdf",
		validSignedReadToken(7, expiresAt),
		now,
	)

	assert.False(t, allowed)
	require.ErrorIs(t, err, lookupErr)
	assert.Equal(t, 1, lookup.calls)
}
