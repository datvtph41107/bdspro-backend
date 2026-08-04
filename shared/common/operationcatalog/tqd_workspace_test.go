package operationcatalog

import (
	_operationregistry "common/operationregistry"

	"testing"
)

func TestTQDWorkspaceV1(t *testing.T) {
	registry, err := _operationregistry.New(TQDWorkspaceV1())
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	definition, err := registry.RequireActive(WorkspaceReportGenerateOperation)
	if err != nil {
		t.Fatalf("RequireActive() error = %v", err)
	}
	if definition.ModuleOwner != "workspace" {
		t.Fatalf("ModuleOwner = %q", definition.ModuleOwner)
	}
	if definition.Metering.SuccessCondition != _operationregistry.SuccessJobAccepted {
		t.Fatalf("SuccessCondition = %q", definition.Metering.SuccessCondition)
	}
	if definition.Metering.FailurePolicy != _operationregistry.FailureCancelReservation {
		t.Fatalf("FailurePolicy = %q", definition.Metering.FailurePolicy)
	}
}
