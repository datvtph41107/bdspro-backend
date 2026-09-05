package request

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestNewOperationIDFromEntropy(
	t *testing.T,
) {
	t.Parallel()

	entropy :=
		bytes.Repeat(
			[]byte{0xab},
			operationIDRandomBytes,
		)

	got, err :=
		newOperationID(
			bytes.NewReader(
				entropy,
			),
		)

	if err != nil {
		t.Fatalf(
			"newOperationID() error = %v",
			err,
		)
	}

	const want = "op_abababababababababababababababab"

	if got != want {
		t.Fatalf(
			"newOperationID() = %q, want %q",
			got,
			want,
		)
	}

	if !IsValidOperationID(
		got,
	) {
		t.Fatalf(
			"generated Operation-ID is invalid: %q",
			got,
		)
	}
}

func TestNewOperationIDRejectsShortEntropy(
	t *testing.T,
) {
	t.Parallel()

	reader :=
		bytes.NewReader(
			bytes.Repeat(
				[]byte{0xab},
				operationIDRandomBytes-1,
			),
		)

	got, err :=
		newOperationID(
			reader,
		)

	if !errors.Is(
		err,
		io.ErrUnexpectedEOF,
	) {
		t.Fatalf(
			"error = %v, want %v",
			err,
			io.ErrUnexpectedEOF,
		)
	}

	if got != "" {
		t.Fatalf(
			"generated ID = %q, want empty",
			got,
		)
	}
}

type failingOperationIDReader struct {
	err error
}

func (
	reader failingOperationIDReader,
) Read(
	[]byte,
) (
	int,
	error,
) {
	return 0,
		reader.err
}

func TestNewOperationIDPropagatesEntropyFailure(
	t *testing.T,
) {
	t.Parallel()

	wantErr :=
		errors.New(
			"entropy unavailable",
		)

	got, err :=
		newOperationID(
			failingOperationIDReader{
				err: wantErr,
			},
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

	if got != "" {
		t.Fatalf(
			"generated ID = %q, want empty",
			got,
		)
	}
}

func TestResolveOperationIDValidCallerValueDoesNotDependOnGenerator(
	t *testing.T,
) {
	t.Parallel()

	generate :=
		func() (
			string,
			error,
		) {
			return "",
				errors.New(
					"entropy unavailable",
				)
		}

	decision, err :=
		resolveOperationID(
			[]string{
				"op-caller-123",
			},
			generate,
		)

	if err != nil {
		t.Fatalf(
			"error = %v, want nil",
			err,
		)
	}

	if decision.ID !=
		"op-caller-123" {

		t.Fatalf(
			"ID = %q, want %q",
			decision.ID,
			"op-caller-123",
		)
	}
}
