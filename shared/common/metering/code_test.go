package metering

import "testing"

func TestParseCodeUsesStableIdentifierGrammar(t *testing.T) {
	t.Parallel()

	code, err := ParseCode(
		"workspace.report_generation.accepted",
	)
	if err != nil {
		t.Fatalf("ParseCode() error = %v", err)
	}
	if code != Code("workspace.report_generation.accepted") {
		t.Fatalf("code = %q", code)
	}
}

func TestParseCodeKeepsExistingWhitespaceNormalization(t *testing.T) {
	t.Parallel()

	code, err := ParseCode(
		" workspace.report_generation.accepted ",
	)
	if err != nil {
		t.Fatalf("ParseCode() error = %v", err)
	}
	if code != Code("workspace.report_generation.accepted") {
		t.Fatalf("code = %q", code)
	}
}

func TestParseCodeRejectsInvalidMeter(t *testing.T) {
	t.Parallel()

	for _, value := range []string{
		"",
		"meter",
		"Workspace.Report",
		"workspace..accepted",
	} {
		value := value
		t.Run(value, func(t *testing.T) {
			t.Parallel()

			if _, err := ParseCode(value); err == nil {
				t.Fatalf("ParseCode(%q) error = nil", value)
			}
		})
	}
}
