package operationbindingcatalog

import (
	_operationbinding "common/operationbinding"
	_operationcatalog "common/operationcatalog"
	_operationregistry "common/operationregistry"

	"testing"
)

func TestTQDWorkspaceV1(t *testing.T) {
	operations, err := _operationregistry.New(_operationcatalog.TQDWorkspaceV1())
	if err != nil {
		t.Fatalf("operation registry error = %v", err)
	}
	bindings, err := _operationbinding.New(TQDWorkspaceV1(), operations)
	if err != nil {
		t.Fatalf("binding registry error = %v", err)
	}
	if err := bindings.RequireResolvable(TQDWorkspaceRequiredTransportsV1()); err != nil {
		t.Fatalf("RequireResolvable() error = %v", err)
	}

	operation, err := bindings.ResolveDirect(
		_operationbinding.GRPCUnary("/types.tqd.MapWorkspaceService/CreateGeneratedReport"),
	)
	if err != nil {
		t.Fatalf("ResolveDirect() error = %v", err)
	}
	if operation != _operationcatalog.WorkspaceReportGenerateOperation {
		t.Fatalf("operation = %q", operation)
	}
}

func TestTQDWorkspaceV1ReturnsIndependentTransportSlices(t *testing.T) {
	first := TQDWorkspaceV1()
	second := TQDWorkspaceV1()
	first.Bindings[0].Transports[0].Method = "/types.tqd.OtherService/Changed"
	if second.Bindings[0].Transports[0].Method != "/types.tqd.MapWorkspaceService/CreateGeneratedReport" {
		t.Fatal("catalog transport slice is shared")
	}
}
