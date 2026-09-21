package handler_http

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	_errors "common/errors"
	"tqd/internal"
)

func TestImportRegionHTTPMapsNotFound(t *testing.T) {
	err := _errors.ReturnError(
		service.ImportLayerNotFound,
		_errors.WithPublicMessage("layer 7 not found"),
	)
	statusCode, body := importHTTPProblem(err)

	if statusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", statusCode, http.StatusNotFound)
	}
	if body["error_code"] != "tqd.import.layer_not_found" {
		t.Fatalf("error_code = %v", body["error_code"])
	}
}

func TestImportRegionHTTPMapsInProgressToConflict(t *testing.T) {
	err := _errors.ReturnError(
		service.ImportInProgress,
		_errors.WithPublicMessage("layer 7 already has an import in progress"),
	)
	statusCode, body := importHTTPProblem(err)

	if statusCode != http.StatusConflict {
		t.Fatalf("status = %d, want %d", statusCode, http.StatusConflict)
	}
	if body["error_code"] != "tqd.import.in_progress" {
		t.Fatalf("error_code = %v", body["error_code"])
	}
}

func TestImportRegionHTTPDoesNotLeakDependencyFailure(t *testing.T) {
	statusCode, body := importHTTPProblem(
		errors.New("postgres password=secret connection failed"),
	)

	if statusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", statusCode, http.StatusInternalServerError)
	}
	if body["error_code"] != "" {
		t.Fatalf("technical error received application identity: %v", body["error_code"])
	}
	message, _ := body["error"].(string)
	if message != "internal server error" {
		t.Fatalf("message = %q, want safe public message", message)
	}
	if strings.Contains(message, "secret") || strings.Contains(message, "postgres") {
		t.Fatalf("internal dependency detail leaked: %q", message)
	}
}
