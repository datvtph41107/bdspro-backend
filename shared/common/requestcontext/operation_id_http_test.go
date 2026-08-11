package requestcontext

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPrepareHTTPOperationIDMissingCommitsGeneratedIdentity(
	t *testing.T,
) {
	t.Parallel()

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"https://example.test/v2/resource",
			nil,
		)

	request.Header.Set(
		"X-Unrelated",
		"keep-me",
	)

	response :=
		httptest.NewRecorder()

	generateCalls := 0

	decision, err :=
		prepareHTTPOperationID(
			response,
			request,
			func() (
				string,
				error,
			) {
				generateCalls++

				return "op-generated-123",
					nil
			},
		)

	if err != nil {
		t.Fatalf(
			"prepare error = %v",
			err,
		)
	}

	if generateCalls != 1 {
		t.Fatalf(
			"generator calls = %d, want 1",
			generateCalls,
		)
	}

	if decision.ID !=
		"op-generated-123" {

		t.Fatalf(
			"decision ID = %q",
			decision.ID,
		)
	}

	if decision.Source !=
		OperationIDSourceGeneratedFallback {

		t.Fatalf(
			"source = %q, want %q",
			decision.Source,
			OperationIDSourceGeneratedFallback,
		)
	}

	if got :=
		request.Header.Get(
			OperationIDHeader,
		); got !=
		decision.ID {

		t.Fatalf(
			"request header = %q, want %q",
			got,
			decision.ID,
		)
	}

	if got :=
		response.Header().
			Get(
				OperationIDHeader,
			); got !=
		decision.ID {

		t.Fatalf(
			"response header = %q, want %q",
			got,
			decision.ID,
		)
	}

	contextID, ok :=
		OperationIDFromContext(
			request.Context(),
		)

	if !ok ||
		contextID !=
			decision.ID {

		t.Fatalf(
			"context ID = %q, %v; want %q",
			contextID,
			ok,
			decision.ID,
		)
	}

	if got :=
		request.Header.Get(
			"X-Unrelated",
		); got !=
		"keep-me" {

		t.Fatalf(
			"unrelated header = %q, want keep-me",
			got,
		)
	}
}

func TestPrepareHTTPOperationIDPreservesValidCallerIdentity(
	t *testing.T,
) {
	t.Parallel()

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"https://example.test/v2/resource",
			nil,
		)

	request.Header.Set(
		OperationIDHeader,
		"  partner.operation-123  ",
	)

	response :=
		httptest.NewRecorder()

	generateCalls := 0

	decision, err :=
		prepareHTTPOperationID(
			response,
			request,
			func() (
				string,
				error,
			) {
				generateCalls++

				return "op-should-not-exist",
					nil
			},
		)

	if err != nil {
		t.Fatalf(
			"prepare error = %v",
			err,
		)
	}

	if generateCalls != 0 {
		t.Fatalf(
			"generator calls = %d, want 0",
			generateCalls,
		)
	}

	const want = "partner.operation-123"

	if decision.ID != want {
		t.Fatalf(
			"decision ID = %q, want %q",
			decision.ID,
			want,
		)
	}

	if got :=
		request.Header.Get(
			OperationIDHeader,
		); got != want {

		t.Fatalf(
			"request header = %q, want %q",
			got,
			want,
		)
	}

	if got :=
		response.Header().
			Get(
				OperationIDHeader,
			); got != want {

		t.Fatalf(
			"response header = %q, want %q",
			got,
			want,
		)
	}

	if got, ok :=
		OperationIDFromContext(
			request.Context(),
		); !ok ||
		got != want {

		t.Fatalf(
			"context ID = %q, %v; want %q",
			got,
			ok,
			want,
		)
	}
}

func TestPrepareHTTPOperationIDInvalidDoesNotMutateBoundary(
	t *testing.T,
) {
	t.Parallel()

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"https://example.test/v2/resource",
			nil,
		)

	const raw = "unsafe operation id"

	request.Header.Set(
		OperationIDHeader,
		raw,
	)

	response :=
		httptest.NewRecorder()

	generateCalls := 0

	decision, err :=
		prepareHTTPOperationID(
			response,
			request,
			func() (
				string,
				error,
			) {
				generateCalls++

				return "op-generated-incorrectly",
					nil
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

	if decision !=
		(OperationIDDecision{}) {

		t.Fatalf(
			"decision = %+v, want zero",
			decision,
		)
	}

	if generateCalls != 0 {
		t.Fatalf(
			"generator calls = %d, want 0",
			generateCalls,
		)
	}

	if got :=
		request.Header.Get(
			OperationIDHeader,
		); got != raw {

		t.Fatalf(
			"request header mutated to %q, want original %q",
			got,
			raw,
		)
	}

	if got :=
		response.Header().
			Get(
				OperationIDHeader,
			); got != "" {

		t.Fatalf(
			"response unexpectedly contains Operation-ID %q",
			got,
		)
	}

	if got, ok :=
		OperationIDFromContext(
			request.Context(),
		); ok {

		t.Fatalf(
			"context unexpectedly contains Operation-ID %q",
			got,
		)
	}
}

func TestPrepareHTTPOperationIDMultipleValuesDoNotMutateBoundary(
	t *testing.T,
) {
	t.Parallel()

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"https://example.test/v2/resource",
			nil,
		)

	request.Header.Add(
		OperationIDHeader,
		"op-a",
	)

	request.Header.Add(
		OperationIDHeader,
		"op-b",
	)

	response :=
		httptest.NewRecorder()

	_, err :=
		prepareHTTPOperationID(
			response,
			request,
			func() (
				string,
				error,
			) {
				t.Fatal(
					"generator must not run for multiple values",
				)

				return "",
					nil
			},
		)

	if !errors.Is(
		err,
		ErrMultipleOperationIDs,
	) {
		t.Fatalf(
			"error = %v, want %v",
			err,
			ErrMultipleOperationIDs,
		)
	}

	got :=
		request.Header.Values(
			OperationIDHeader,
		)

	if len(got) != 2 ||
		got[0] != "op-a" ||
		got[1] != "op-b" {

		t.Fatalf(
			"request headers mutated: %v",
			got,
		)
	}

	if response.Header().
		Get(
			OperationIDHeader,
		) != "" {

		t.Fatal(
			"response Operation-ID was committed on failure",
		)
	}
}

func TestPrepareHTTPOperationIDGenerationFailureDoesNotMutateBoundary(
	t *testing.T,
) {
	t.Parallel()

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"https://example.test/v2/resource",
			nil,
		)

	request.Header = nil

	response :=
		httptest.NewRecorder()

	wantErr :=
		ErrOperationIDGeneration

	decision, err :=
		prepareHTTPOperationID(
			response,
			request,
			func() (
				string,
				error,
			) {
				return "",
					wantErr
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

	if decision !=
		(OperationIDDecision{}) {

		t.Fatalf(
			"decision = %+v, want zero",
			decision,
		)
	}

	if request.Header != nil {
		t.Fatalf(
			"request Header was allocated on failure: %v",
			request.Header,
		)
	}

	if response.Header().
		Get(
			OperationIDHeader,
		) != "" {

		t.Fatal(
			"response Operation-ID committed on generation failure",
		)
	}

	if _, ok :=
		OperationIDFromContext(
			request.Context(),
		); ok {

		t.Fatal(
			"context Operation-ID committed on generation failure",
		)
	}
}

func TestPrepareHTTPOperationIDContextConflictDoesNotPartiallyCommit(
	t *testing.T,
) {
	t.Parallel()

	parent, err :=
		BindOperationID(
			context.Background(),
			"op-existing-123",
		)

	if err != nil {
		t.Fatalf(
			"BindOperationID() error = %v",
			err,
		)
	}

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"https://example.test/v2/resource",
			nil,
		).
			WithContext(
				parent,
			)

	const raw = "  op-new-456  "

	request.Header.Set(
		OperationIDHeader,
		raw,
	)

	response :=
		httptest.NewRecorder()

	decision, err :=
		prepareHTTPOperationID(
			response,
			request,
			func() (
				string,
				error,
			) {
				t.Fatal(
					"generator must not run",
				)

				return "",
					nil
			},
		)

	if !errors.Is(
		err,
		ErrOperationIDContextConflict,
	) {
		t.Fatalf(
			"error = %v, want %v",
			err,
			ErrOperationIDContextConflict,
		)
	}

	if decision !=
		(OperationIDDecision{}) {

		t.Fatalf(
			"decision = %+v, want zero",
			decision,
		)
	}

	// Even normalization must not have been partially committed.
	if got :=
		request.Header.Get(
			OperationIDHeader,
		); got != raw {

		t.Fatalf(
			"request header partially mutated to %q, want original %q",
			got,
			raw,
		)
	}

	if got :=
		response.Header().
			Get(
				OperationIDHeader,
			); got != "" {

		t.Fatalf(
			"response partially committed Operation-ID %q",
			got,
		)
	}

	got, ok :=
		OperationIDFromContext(
			request.Context(),
		)

	if !ok ||
		got !=
			"op-existing-123" {

		t.Fatalf(
			"context ID = %q, %v; want existing ID",
			got,
			ok,
		)
	}
}

func TestPrepareHTTPOperationIDAllowsSameExistingContextIdentity(
	t *testing.T,
) {
	t.Parallel()

	parent, err :=
		BindOperationID(
			context.Background(),
			"op-same-123",
		)

	if err != nil {
		t.Fatal(err)
	}

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"https://example.test/v2/resource",
			nil,
		).
			WithContext(
				parent,
			)

	request.Header.Set(
		OperationIDHeader,
		"op-same-123",
	)

	response :=
		httptest.NewRecorder()

	decision, err :=
		prepareHTTPOperationID(
			response,
			request,
			func() (
				string,
				error,
			) {
				t.Fatal(
					"generator must not run",
				)

				return "",
					nil
			},
		)

	if err != nil {
		t.Fatalf(
			"prepare error = %v",
			err,
		)
	}

	if decision.ID !=
		"op-same-123" {

		t.Fatalf(
			"decision ID = %q",
			decision.ID,
		)
	}
}

func TestPrepareHTTPOperationIDPanicsOnNilRequest(
	t *testing.T,
) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal(
				"expected panic",
			)
		}
	}()

	_, _ =
		PrepareHTTPOperationID(
			httptest.NewRecorder(),
			nil,
		)
}
