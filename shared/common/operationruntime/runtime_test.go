package operationruntime

import (
	_operationbinding "common/operationbinding"
	_operationcontext "common/operationcontext"
	_operationregistry "common/operationregistry"

	"context"
	"errors"
	"testing"
)

const (
	testOperation _operationregistry.Code = "workspace.report.generate"
	testBinding   _operationregistry.Code = "workspace.report.generate_transport"
	testResolver  _operationregistry.Code = "workspace.report.resolve"
)

func testRuntime(t *testing.T, disposition _operationbinding.Disposition, resolvers map[_operationregistry.Code]Resolver) *Runtime {
	t.Helper()

	operations, err := _operationregistry.New(_operationregistry.Catalog{
		Version: "1.0.0",
		Definitions: []_operationregistry.Definition{
			{
				Code:        testOperation,
				ModuleOwner: "workspace",
				FeatureCode: "workspace.generated_reports",
				Metering: _operationregistry.Metering{
					MeterCode:        "workspace.report_generation.accepted",
					ChargeMode:       _operationregistry.ChargeModeReserveThenCommit,
					SuccessCondition: _operationregistry.SuccessJobAccepted,
					FailurePolicy:    _operationregistry.FailureCancelReservation,
				},
				Status:          _operationregistry.StatusActive,
				SemanticVersion: "1.0.0",
				Description:     "Test operation.",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	binding := _operationbinding.Binding{
		Code:        testBinding,
		ModuleOwner: "workspace",
		Transports: []_operationbinding.Transport{
			_operationbinding.GRPCUnary("/types.tqd.MapWorkspaceService/CreateGeneratedReport"),
		},
		Disposition: disposition,
		Description: "Test binding.",
	}
	if disposition == _operationbinding.DispositionDirect {
		binding.OperationCode = testOperation
	}
	if disposition == _operationbinding.DispositionResolver {
		binding.ResolverCode = testResolver
	}

	bindings, err := _operationbinding.New(_operationbinding.Catalog{
		Version:  "1.0.0",
		Bindings: []_operationbinding.Binding{binding},
	}, operations)
	if err != nil {
		t.Fatal(err)
	}

	runtime, err := New(operations, bindings, resolvers)
	if err != nil {
		t.Fatal(err)
	}
	return runtime
}

func TestResolveDirect(t *testing.T) {
	runtime := testRuntime(t, _operationbinding.DispositionDirect, nil)
	selection, ok, err := runtime.Resolve(
		context.Background(),
		struct{}{},
		_operationbinding.GRPCUnary("/types.tqd.MapWorkspaceService/CreateGeneratedReport"),
	)
	if err != nil || !ok {
		t.Fatalf("Resolve() ok=%v err=%v", ok, err)
	}
	if selection.Code != testOperation || selection.BindingCode != testBinding {
		t.Fatalf("selection = %+v", selection)
	}
	if selection.Source != _operationcontext.SourceDirectBinding {
		t.Fatalf("source = %q", selection.Source)
	}
}

func TestResolveUnknownTransport(t *testing.T) {
	runtime := testRuntime(t, _operationbinding.DispositionDirect, nil)
	_, ok, err := runtime.Resolve(
		context.Background(), nil,
		_operationbinding.GRPCUnary("/types.tqd.MapWorkspaceService/ListGeneratedReports"),
	)
	if err != nil || ok {
		t.Fatalf("Resolve() ok=%v err=%v", ok, err)
	}
}

func TestResolvePreservesExistingOperation(t *testing.T) {
	runtime := testRuntime(t, _operationbinding.DispositionDirect, nil)
	selection, err := _operationcontext.New(testOperation, testBinding, _operationcontext.SourceDirectBinding)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := _operationcontext.Bind(context.Background(), selection)
	if err != nil {
		t.Fatal(err)
	}

	resolved, ok, err := runtime.Resolve(
		ctx, nil,
		_operationbinding.GRPCUnary("/types.tqd.MapWorkspaceService/ListGeneratedReports"),
	)
	if err != nil || !ok || resolved != selection {
		t.Fatalf("Resolve() selection=%+v ok=%v err=%v", resolved, ok, err)
	}
}

func TestRequireReadyRejectsMissingResolver(t *testing.T) {
	runtime := testRuntime(t, _operationbinding.DispositionResolver, nil)
	err := runtime.RequireReady([]_operationbinding.Transport{
		_operationbinding.GRPCUnary("/types.tqd.MapWorkspaceService/CreateGeneratedReport"),
	})
	if !errors.Is(err, ErrResolverNotConfigured) {
		t.Fatalf("RequireReady() error = %v", err)
	}
}

func TestResolveDomainResolver(t *testing.T) {
	runtime := testRuntime(t, _operationbinding.DispositionResolver, map[_operationregistry.Code]Resolver{
		testResolver: func(context.Context, interface{}) (_operationregistry.Code, error) {
			return testOperation, nil
		},
	})
	selection, ok, err := runtime.Resolve(
		context.Background(), nil,
		_operationbinding.GRPCUnary("/types.tqd.MapWorkspaceService/CreateGeneratedReport"),
	)
	if err != nil || !ok {
		t.Fatalf("Resolve() ok=%v err=%v", ok, err)
	}
	if selection.Source != _operationcontext.SourceDomainResolver {
		t.Fatalf("source = %q", selection.Source)
	}
}
