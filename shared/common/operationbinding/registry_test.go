package operationbinding

import (
	_operationregistry "common/operationregistry"

	"errors"
	"testing"
)

func TestRegistryResolvesDirectBinding(t *testing.T) {
	t.Parallel()

	operations := testOperationRegistry(t)
	registry, err := New(Catalog{
		Version: "1.0.0",
		Bindings: []Binding{
			{
				Code:        "parcel.lookup_transport",
				ModuleOwner: "parcel",
				Transports: []Transport{
					GRPCUnary("/tqdpb.ParcelService/ResolveMapTargetByLocation"),
					HTTPRoute("POST", "/v2/tqd/map-target/by-location"),
				},
				Disposition:   DispositionDirect,
				OperationCode: "parcel.lookup",
				Description:   "Maps parcel lookup entrypoints to one stable operation.",
			},
		},
	}, operations)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	code, err := registry.ResolveDirect(GRPCUnary("/tqdpb.ParcelService/ResolveMapTargetByLocation"))
	if err != nil {
		t.Fatalf("ResolveDirect() error = %v", err)
	}
	if code != "parcel.lookup" {
		t.Fatalf("ResolveDirect() = %q, want parcel.lookup", code)
	}

	httpBinding, ok := registry.LookupTransport(HTTPRoute("post", "/v2/tqd/map-target/by-location"))
	if !ok || httpBinding.Code != "parcel.lookup_transport" {
		t.Fatalf("LookupTransport() = %#v, %v", httpBinding, ok)
	}
}

func TestRegistryRequiresResolver(t *testing.T) {
	t.Parallel()

	registry, err := New(Catalog{
		Version: "1.0.0",
		Bindings: []Binding{
			{
				Code:         "parcel.planning_transport",
				ModuleOwner:  "parcel",
				Transports:   []Transport{GRPCUnary("/tqdpb.ParcelService/GetParcelPlanning")},
				Disposition:  DispositionResolver,
				ResolverCode: "parcel.resolve_planning_operation",
				Description:  "Uses request fields to select the planning operation.",
			},
		},
	}, nil)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = registry.ResolveDirect(GRPCUnary("/tqdpb.ParcelService/GetParcelPlanning"))
	if !errors.Is(err, ErrResolverRequired) {
		t.Fatalf("ResolveDirect() error = %v, want ErrResolverRequired", err)
	}
}

func TestRegistryRejectsDuplicateTransport(t *testing.T) {
	t.Parallel()

	transport := GRPCUnary("/tqdpb.ParcelService/GetParcelPlanning")
	_, err := New(Catalog{
		Version: "1.0.0",
		Bindings: []Binding{
			{
				Code:        "parcel.first_review",
				ModuleOwner: "parcel",
				Transports:  []Transport{transport},
				Disposition: DispositionReviewRequired,
				Description: "First review record.",
			},
			{
				Code:        "parcel.second_review",
				ModuleOwner: "parcel",
				Transports:  []Transport{transport},
				Disposition: DispositionReviewRequired,
				Description: "Second review record.",
			},
		},
	}, nil)
	if err == nil {
		t.Fatal("New() error = nil, want duplicate transport error")
	}
}

func TestRegistryRejectsOwnerMismatch(t *testing.T) {
	t.Parallel()

	_, err := New(Catalog{
		Version: "1.0.0",
		Bindings: []Binding{
			{
				Code:          "planning.lookup_transport",
				ModuleOwner:   "planning",
				Transports:    []Transport{GRPCUnary("/tqdpb.ParcelService/ResolveMapTargetByLocation")},
				Disposition:   DispositionDirect,
				OperationCode: "parcel.lookup",
				Description:   "Invalid owner mapping.",
			},
		},
	}, testOperationRegistry(t))
	if err == nil {
		t.Fatal("New() error = nil, want owner mismatch")
	}
}

func TestRegistryCoverageDistinguishesReviewedFromResolvable(t *testing.T) {
	t.Parallel()

	transport := GRPCUnary("/tqdpb.ParcelService/FindParcelsByPolygon")
	registry, err := New(Catalog{
		Version: "1.0.0",
		Bindings: []Binding{
			{
				Code:        "parcel.polygon_review",
				ModuleOwner: "parcel",
				Transports:  []Transport{transport},
				Disposition: DispositionReviewRequired,
				Description: "Business operation boundary is not approved.",
			},
		},
	}, nil)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := registry.RequireReviewed([]Transport{transport}); err != nil {
		t.Fatalf("RequireReviewed() error = %v", err)
	}
	if err := registry.RequireResolvable([]Transport{transport}); !errors.Is(err, ErrNotResolvable) {
		t.Fatalf("RequireResolvable() error = %v, want ErrNotResolvable", err)
	}
}

func TestRegistryReturnsImmutableSnapshots(t *testing.T) {
	t.Parallel()

	registry, err := New(Catalog{
		Version: "1.0.0",
		Bindings: []Binding{
			{
				Code:        "parcel.review",
				ModuleOwner: "parcel",
				Transports:  []Transport{GRPCUnary("/tqdpb.ParcelService/FindParcelByLocation")},
				Disposition: DispositionReviewRequired,
				Description: "Review entry.",
			},
		},
	}, nil)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	snapshot := registry.Bindings()
	snapshot[0].Transports[0].Method = "/changed.Service/Method"

	binding, ok := registry.LookupCode("parcel.review")
	if !ok {
		t.Fatal("LookupCode() missing binding")
	}
	if binding.Transports[0].Method != "/tqdpb.ParcelService/FindParcelByLocation" {
		t.Fatalf("registry mutated: %q", binding.Transports[0].Method)
	}
}

func testOperationRegistry(t *testing.T) *_operationregistry.Registry {
	t.Helper()
	registry, err := _operationregistry.New(_operationregistry.Catalog{
		Version: "1.0.0",
		Definitions: []_operationregistry.Definition{
			{
				Code:        "parcel.lookup",
				ModuleOwner: "parcel",
				FeatureCode: "parcel.detail",
				Metering: _operationregistry.Metering{
					ChargeMode:       _operationregistry.ChargeModeUnmetered,
					SuccessCondition: _operationregistry.SuccessBusinessOutcome,
					FailurePolicy:    _operationregistry.FailureNoCharge,
				},
				Status:          _operationregistry.StatusActive,
				SemanticVersion: "1.0.0",
				Description:     "Resolves one parcel lookup business outcome.",
			},
		},
	})
	if err != nil {
		t.Fatalf("operationregistry.New() error = %v", err)
	}
	return registry
}
