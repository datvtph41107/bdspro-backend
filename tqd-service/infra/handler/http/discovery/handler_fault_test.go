package discovery

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	_errors "common/errors"
	"tqd/internal"
)

func TestDiscoveryIdentifyValidationMapsToBadRequest(t *testing.T) {
	statusCode, payload := identifyHTTPProblem(
		_errors.ReturnError(service.DiscoveryCoordinateInvalid),
	)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if payload["code"] != "DISCOVERY_IDENTIFY_FAILED" {
		t.Fatalf("legacy code = %v", payload["code"])
	}
	if payload["error_code"] != "tqd.discovery.coordinate_invalid" {
		t.Fatalf("error_code = %v", payload["error_code"])
	}
	if payload["message"] != "invalid latitude/longitude" {
		t.Fatalf("message = %v", payload["message"])
	}
}

func TestDiscoveryIdentifyUnknownFailureDoesNotLeak(t *testing.T) {
	dependencyErr := errors.New("postgres password=secret discovery provider failed")
	statusCode, payload := identifyHTTPProblem(dependencyErr)

	if statusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", statusCode, http.StatusInternalServerError)
	}
	if payload["code"] != "DISCOVERY_IDENTIFY_FAILED" {
		t.Fatalf("legacy code = %v", payload["code"])
	}
	if payload["error_code"] != "" {
		t.Fatalf("technical error received application identity: %v", payload["error_code"])
	}

	message, _ := payload["message"].(string)
	if message != "internal server error" {
		t.Fatalf("message = %q", message)
	}
	if strings.Contains(message, "secret") || strings.Contains(message, "postgres") {
		t.Fatalf("dependency detail leaked: %q", message)
	}
}
