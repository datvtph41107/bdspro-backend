package operation

import (
	"errors"
	"strings"
	"testing"
)

func TestParseAcceptsStableDottedCode(t *testing.T) {
	t.Parallel()

	for _, raw := range []string{
		"workspace.report.generate",
		"parcel.view",
		"payment.refund_request",
	} {
		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(raw)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", raw, err)
			}
			if got != Code(raw) {
				t.Fatalf("Parse(%q) = %q", raw, got)
			}
		})
	}
}

func TestParseRejectsInvalidCode(t *testing.T) {
	t.Parallel()

	tooLong := "workspace." + strings.Repeat("a", MaxCodeLength)
	tests := []string{
		"",
		"report",
		"Workspace.Report",
		"workspace.report generate",
		"workspace..generate",
		"workspace.1report",
		" workspace.report.generate",
		tooLong,
	}

	for _, raw := range tests {
		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			if _, err := Parse(raw); !errors.Is(err, ErrInvalidCode) {
				t.Fatalf(
					"Parse(%q) error = %v, want ErrInvalidCode",
					raw,
					err,
				)
			}
		})
	}
}
