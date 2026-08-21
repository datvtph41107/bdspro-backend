package request

import (
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewRequestIDHasCanonicalShape(
	t *testing.T,
) {
	t.Parallel()

	requestID := NewRequestID()

	if !strings.HasPrefix(
		requestID,
		"req_",
	) {
		t.Fatalf(
			"request ID = %q, missing req_ prefix",
			requestID,
		)
	}

	if !IsValidRequestID(requestID) {
		t.Fatalf(
			"generated request ID is invalid: %q",
			requestID,
		)
	}
}

func TestNewRequestIDUsesFallbackWhenReadFails(
	t *testing.T,
) {
	t.Parallel()

	fixedTime := time.Unix(
		1_700_000_000,
		123,
	)

	var sequence atomic.Uint64

	requestID := newRequestID(
		func([]byte) (int, error) {
			return 0, errors.New("entropy unavailable")
		},
		func() time.Time {
			return fixedTime
		},
		&sequence,
	)

	if !IsValidRequestID(requestID) {
		t.Fatalf(
			"fallback request ID is invalid: %q",
			requestID,
		)
	}

	if !strings.HasPrefix(
		requestID,
		"req_",
	) {
		t.Fatalf(
			"fallback request ID = %q",
			requestID,
		)
	}
}

func TestNewRequestIDIsUniqueAcrossSample(
	t *testing.T,
) {
	t.Parallel()

	const sampleSize = 2_000

	seen := make(
		map[string]struct{},
		sampleSize,
	)

	for index := 0; index < sampleSize; index++ {
		requestID := NewRequestID()

		if _, exists := seen[requestID]; exists {
			t.Fatalf(
				"duplicate request ID: %q",
				requestID,
			)
		}

		seen[requestID] = struct{}{}
	}
}
