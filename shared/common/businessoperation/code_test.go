package businessoperation

import "testing"

func TestCodeIsValid(t *testing.T) {
	tests := []struct {
		name  string
		code  Code
		valid bool
	}{
		{name: "report", code: WorkspaceReportGenerate, valid: true},
		{name: "parcel view", code: Code("parcel.view"), valid: true},
		{name: "single part", code: "report", valid: false},
		{name: "uppercase", code: "Workspace.Report", valid: false},
		{name: "space", code: "workspace.report generate", valid: false},
		{name: "empty part", code: "workspace..generate", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}
