package codeformat

import (
	"strings"
	"testing"
)

func TestIsValidAcceptsStableDottedIdentifiers(t *testing.T) {
	t.Parallel()

	for _, value := range []string{
		"workspace.report.generate",
		"workspace.report_generation.accepted",
		"parcel.view",
		"payment.refund_request",
	} {
		value := value
		t.Run(value, func(t *testing.T) {
			t.Parallel()

			if !IsValid(value) {
				t.Fatalf("IsValid(%q) = false", value)
			}
		})
	}
}

func TestIsValidRejectsMalformedIdentifiers(t *testing.T) {
	t.Parallel()

	for _, value := range []string{
		"",
		"report",
		"Workspace.Report",
		"workspace.report generate",
		"workspace..generate",
		"workspace.1report",
		" workspace.report.generate",
		"workspace.report.generate ",
		"workspace." + strings.Repeat("a", MaxLength),
	} {
		value := value
		t.Run(value, func(t *testing.T) {
			t.Parallel()

			if IsValid(value) {
				t.Fatalf("IsValid(%q) = true", value)
			}
		})
	}
}
