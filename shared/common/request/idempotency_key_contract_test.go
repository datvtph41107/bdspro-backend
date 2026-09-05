package request

import (
	"errors"
	"strings"
	"testing"
)

func TestIsValidIdempotencyKey(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "opaque safe value", value: "idem-AbC_123.v1:retry", want: true},
		{name: "maximum length", value: strings.Repeat("a", MaxIdempotencyKeyLength), want: true},
		{name: "empty", value: "", want: false},
		{name: "too long", value: strings.Repeat("a", MaxIdempotencyKeyLength+1), want: false},
		{name: "space", value: "unsafe key", want: false},
		{name: "slash", value: "unsafe/key", want: false},
		{name: "newline", value: "unsafe\nkey", want: false},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			if got := IsValidIdempotencyKey(testCase.value); got != testCase.want {
				t.Fatalf("IsValidIdempotencyKey(%q) = %v, want %v", testCase.value, got, testCase.want)
			}
		})
	}
}

func TestParseIdempotencyKey(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		values  []string
		wantKey string
		wantOK  bool
		wantErr error
	}{
		{name: "missing is optional", values: nil},
		{name: "one valid value", values: []string{"  idem-command-123  "}, wantKey: "idem-command-123", wantOK: true},
		{name: "one invalid value", values: []string{"unsafe command key"}, wantErr: ErrInvalidIdempotencyKey},
		{name: "two different values", values: []string{"idem-first", "idem-second"}, wantErr: ErrMultipleIdempotencyKeys},
		{name: "two repeated values", values: []string{"idem-repeat", "idem-repeat"}, wantErr: ErrMultipleIdempotencyKeys},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			key, ok, err := ParseIdempotencyKey(testCase.values)
			if !errors.Is(err, testCase.wantErr) {
				t.Fatalf("ParseIdempotencyKey() error = %v, want %v", err, testCase.wantErr)
			}
			if key != testCase.wantKey || ok != testCase.wantOK {
				t.Fatalf("ParseIdempotencyKey() = %q, %v; want %q, %v", key, ok, testCase.wantKey, testCase.wantOK)
			}
		})
	}
}
