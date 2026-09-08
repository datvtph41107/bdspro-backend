package request

import (
	"strings"
	"testing"
)

func TestIsValidRequestID(
	t *testing.T,
) {
	t.Parallel()

	testCases := []struct {
		name  string
		value string
		valid bool
	}{
		{
			name:  "server generated",
			value: "req_0123456789abcdef",
			valid: true,
		},
		{
			name:  "uuid",
			value: "550e8400-e29b-41d4-a716-446655440000",
			valid: true,
		},
		{
			name:  "colon and dot",
			value: "client.web:request.123",
			valid: true,
		},
		{
			name:  "empty",
			value: "",
			valid: false,
		},
		{
			name:  "space",
			value: "request id",
			valid: false,
		},
		{
			name:  "newline",
			value: "request\nid",
			valid: false,
		},
		{
			name:  "slash",
			value: "request/id",
			valid: false,
		},
		{
			name:  "too long",
			value: strings.Repeat("a", MaxRequestIDLength+1),
			valid: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual := IsValidRequestID(
				testCase.value,
			)

			if actual != testCase.valid {
				t.Fatalf(
					"IsValidRequestID(%q) = %v, want %v",
					testCase.value,
					actual,
					testCase.valid,
				)
			}
		})
	}
}

func TestResolveRequestID(
	t *testing.T,
) {
	t.Parallel()

	generate := func() string {
		return "req_generated_123"
	}

	testCases := []struct {
		name       string
		values     []string
		wantID     string
		wantSource RequestIDSource
		wantReason RequestIDReason
	}{
		{
			name:       "missing",
			values:     nil,
			wantID:     "req_generated_123",
			wantSource: RequestIDSourceGenerated,
			wantReason: RequestIDReasonMissing,
		},
		{
			name:       "valid",
			values:     []string{" client-request-123 "},
			wantID:     "client-request-123",
			wantSource: RequestIDSourceClient,
			wantReason: RequestIDReasonValid,
		},
		{
			name:       "invalid",
			values:     []string{"unavaiable exist"},
			wantID:     "req_generated_123",
			wantSource: RequestIDSourceGenerated,
			wantReason: RequestIDReasonInvalid,
		},
		{
			name:       "multiple",
			values:     []string{"req-a", "req-b"},
			wantID:     "req_generated_123",
			wantSource: RequestIDSourceGenerated,
			wantReason: RequestIDReasonMultiple,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual := resolveRequestID(
				testCase.values,
				generate,
			)

			if actual.ID != testCase.wantID {
				t.Fatalf(
					"ID = %q, want %q",
					actual.ID,
					testCase.wantID,
				)
			}

			if actual.Source != testCase.wantSource {
				t.Fatalf(
					"Source = %q, want %q",
					actual.Source,
					testCase.wantSource,
				)
			}

			if actual.Reason != testCase.wantReason {
				t.Fatalf(
					"Reason = %q, want %q",
					actual.Reason,
					testCase.wantReason,
				)
			}
		})
	}
}
