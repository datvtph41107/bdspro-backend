package request

import (
	"errors"
	"strings"
	"testing"
)

func TestNormalizeOperationID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "unchanged",
			in:   "op-business-123",
			want: "op-business-123",
		},
		{
			name: "trim surrounding whitespace",
			in:   "  op-business-123  ",
			want: "op-business-123",
		},
		{
			name: "empty",
			in:   "",
			want: "",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := NormalizeOperationID(
				testCase.in,
			)

			if got != testCase.want {
				t.Fatalf(
					"NormalizeOperationID(%q) = %q, want %q",
					testCase.in,
					got,
					testCase.want,
				)
			}
		})
	}
}

func TestIsValidOperationID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "generated style",
			value: "op_abcdef123456",
			want:  true,
		},
		{
			name:  "caller owned safe identifier",
			value: "checkout:order-123",
			want:  true,
		},
		{
			name:  "all supported separators",
			value: "operation.test_value:123-abc",
			want:  true,
		},
		{
			name:  "maximum length",
			value: strings.Repeat("a", MaxOperationIDLength),
			want:  true,
		},
		{
			name:  "empty",
			value: "",
			want:  false,
		},
		{
			name:  "whitespace only",
			value: "   ",
			want:  false,
		},
		{
			name:  "surrounding whitespace is not validation",
			value: " op-business-123 ",
			want:  false,
		},
		{
			name:  "embedded whitespace",
			value: "op business",
			want:  false,
		},
		{
			name:  "slash",
			value: "op/business",
			want:  false,
		},
		{
			name:  "unicode",
			value: "op_đồ-án",
			want:  false,
		},
		{
			name: "too long",
			value: strings.Repeat(
				"a",
				MaxOperationIDLength+1,
			),
			want: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := IsValidOperationID(
				testCase.value,
			)

			if got != testCase.want {
				t.Fatalf(
					"IsValidOperationID(%q) = %v, want %v",
					testCase.value,
					got,
					testCase.want,
				)
			}
		})
	}
}

func TestAcceptOperationID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		values []string
		wantID string
		wantOK bool
	}{
		{
			name: "one valid",
			values: []string{
				"op-business-123",
			},
			wantID: "op-business-123",
			wantOK: true,
		},
		{
			name: "one valid normalized",
			values: []string{
				"  op-business-123  ",
			},
			wantID: "op-business-123",
			wantOK: true,
		},
		{
			name:   "missing",
			values: nil,
			wantID: "",
			wantOK: false,
		},
		{
			name: "invalid",
			values: []string{
				"unsafe operation id",
			},
			wantID: "",
			wantOK: false,
		},
		{
			name: "multiple distinct",
			values: []string{
				"op-a",
				"op-b",
			},
			wantID: "",
			wantOK: false,
		},
		{
			name: "multiple identical",
			values: []string{
				"op-a",
				"op-a",
			},
			wantID: "",
			wantOK: false,
		},
		{
			name: "one normalized empty value",
			values: []string{
				"   ",
			},
			wantID: "",
			wantOK: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			gotID, gotOK := AcceptOperationID(
				testCase.values,
			)

			if gotOK != testCase.wantOK {
				t.Fatalf(
					"AcceptOperationID(%v) ok = %v, want %v",
					testCase.values,
					gotOK,
					testCase.wantOK,
				)
			}

			if gotID != testCase.wantID {
				t.Fatalf(
					"AcceptOperationID(%v) ID = %q, want %q",
					testCase.values,
					gotID,
					testCase.wantID,
				)
			}
		})
	}
}

func TestResolveOperationIDMissingGeneratesFallback(t *testing.T) {
	t.Parallel()

	generateCalls := 0

	generate := func() (string, error) {
		generateCalls++

		return "op-generated-123", nil
	}

	decision, err := resolveOperationID(
		nil,
		generate,
	)

	if err != nil {
		t.Fatalf(
			"resolveOperationID() error = %v",
			err,
		)
	}

	if generateCalls != 1 {
		t.Fatalf(
			"generator calls = %d, want 1",
			generateCalls,
		)
	}

	if decision.ID != "op-generated-123" {
		t.Fatalf(
			"decision ID = %q, want %q",
			decision.ID,
			"op-generated-123",
		)
	}

	if decision.Source !=
		OperationIDSourceGeneratedFallback {

		t.Fatalf(
			"decision source = %q, want %q",
			decision.Source,
			OperationIDSourceGeneratedFallback,
		)
	}

	if decision.Reason !=
		OperationIDReasonMissing {

		t.Fatalf(
			"decision reason = %q, want %q",
			decision.Reason,
			OperationIDReasonMissing,
		)
	}
}

func TestResolveOperationIDValidCallerValueDoesNotGenerate(t *testing.T) {
	t.Parallel()

	generateCalls := 0

	generate := func() (string, error) {
		generateCalls++

		return "op-should-not-exist", nil
	}

	decision, err := resolveOperationID(
		[]string{
			"  partner.operation-123  ",
		},
		generate,
	)

	if err != nil {
		t.Fatalf(
			"resolveOperationID() error = %v",
			err,
		)
	}

	if generateCalls != 0 {
		t.Fatalf(
			"generator calls = %d, want 0",
			generateCalls,
		)
	}

	if decision.ID !=
		"partner.operation-123" {

		t.Fatalf(
			"decision ID = %q, want %q",
			decision.ID,
			"partner.operation-123",
		)
	}

	if decision.Source !=
		OperationIDSourceCaller {

		t.Fatalf(
			"decision source = %q, want %q",
			decision.Source,
			OperationIDSourceCaller,
		)
	}

	if decision.Reason !=
		OperationIDReasonValid {

		t.Fatalf(
			"decision reason = %q, want %q",
			decision.Reason,
			OperationIDReasonValid,
		)
	}
}

func TestResolveOperationIDInvalidValueDoesNotGenerate(t *testing.T) {
	t.Parallel()

	generateCalls := 0

	generate := func() (string, error) {
		generateCalls++

		return "op-generated-incorrectly", nil
	}

	decision, err := resolveOperationID(
		[]string{
			"unsafe operation id",
		},
		generate,
	)

	if !errors.Is(
		err,
		ErrInvalidOperationID,
	) {
		t.Fatalf(
			"resolveOperationID() error = %v, want %v",
			err,
			ErrInvalidOperationID,
		)
	}

	if generateCalls != 0 {
		t.Fatalf(
			"generator calls = %d, want 0",
			generateCalls,
		)
	}

	if decision !=
		(OperationIDDecision{}) {

		t.Fatalf(
			"decision = %+v, want zero decision",
			decision,
		)
	}
}

func TestResolveOperationIDMultipleValuesDoNotGenerate(t *testing.T) {
	t.Parallel()

	generateCalls := 0

	generate := func() (string, error) {
		generateCalls++

		return "op-generated-incorrectly", nil
	}

	decision, err := resolveOperationID(
		[]string{
			"op-a",
			"op-b",
		},
		generate,
	)

	if !errors.Is(
		err,
		ErrMultipleOperationIDs,
	) {
		t.Fatalf(
			"resolveOperationID() error = %v, want %v",
			err,
			ErrMultipleOperationIDs,
		)
	}

	if generateCalls != 0 {
		t.Fatalf(
			"generator calls = %d, want 0",
			generateCalls,
		)
	}

	if decision !=
		(OperationIDDecision{}) {

		t.Fatalf(
			"decision = %+v, want zero decision",
			decision,
		)
	}
}

func TestResolveOperationIDRepeatedValueIsStillMultiple(t *testing.T) {
	t.Parallel()

	generateCalls := 0

	generate := func() (string, error) {
		generateCalls++

		return "op-generated-incorrectly", nil
	}

	_, err := resolveOperationID(
		[]string{
			"op-same",
			"op-same",
		},
		generate,
	)

	if !errors.Is(
		err,
		ErrMultipleOperationIDs,
	) {
		t.Fatalf(
			"resolveOperationID() error = %v, want %v",
			err,
			ErrMultipleOperationIDs,
		)
	}

	if generateCalls != 0 {
		t.Fatalf(
			"generator calls = %d, want 0",
			generateCalls,
		)
	}
}

func TestResolveOperationIDWhitespaceOnlyValueIsInvalid(
	t *testing.T,
) {
	t.Parallel()

	generateCalls := 0

	_, err :=
		resolveOperationID(
			[]string{
				"   ",
			},
			func() (string, error) {
				generateCalls++

				return "op-generated", nil
			},
		)

	if !errors.Is(
		err,
		ErrInvalidOperationID,
	) {
		t.Fatalf(
			"error = %v, want %v",
			err,
			ErrInvalidOperationID,
		)
	}

	if generateCalls != 0 {
		t.Fatalf(
			"generator calls = %d, want 0",
			generateCalls,
		)
	}
}

func TestResolveOperationIDMissingPropagatesGenerationFailure(
	t *testing.T,
) {
	t.Parallel()

	generateCalls := 0

	wantErr :=
		errors.New(
			"entropy unavailable",
		)

	generate :=
		func() (
			string,
			error,
		) {
			generateCalls++

			return "",
				wantErr
		}

	decision, err :=
		resolveOperationID(
			nil,
			generate,
		)

	if !errors.Is(
		err,
		wantErr,
	) {
		t.Fatalf(
			"error = %v, want %v",
			err,
			wantErr,
		)
	}

	if generateCalls != 1 {
		t.Fatalf(
			"generator calls = %d, want 1",
			generateCalls,
		)
	}

	if decision !=
		(OperationIDDecision{}) {

		t.Fatalf(
			"decision = %+v, want zero decision",
			decision,
		)
	}
}

func TestResolveOperationIDRejectsInvalidGeneratedValue(
	t *testing.T,
) {
	t.Parallel()

	generate :=
		func() (
			string,
			error,
		) {
			return "unsafe operation id",
				nil
		}

	decision, err :=
		resolveOperationID(
			nil,
			generate,
		)

	if !errors.Is(
		err,
		ErrOperationIDGeneration,
	) {
		t.Fatalf(
			"error = %v, want %v",
			err,
			ErrOperationIDGeneration,
		)
	}

	if decision !=
		(OperationIDDecision{}) {

		t.Fatalf(
			"decision = %+v, want zero decision",
			decision,
		)
	}
}
