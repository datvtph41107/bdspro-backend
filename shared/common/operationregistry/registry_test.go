package operationregistry

import (
	"errors"
	"reflect"
	"testing"
)

func TestRegistryBuildsImmutableSnapshots(t *testing.T) {
	t.Parallel()

	definitions := []Definition{
		operation("planning.report_generate", "planning"),
		operation("parcel.lookup", "parcel"),
	}
	registry, err := New(Catalog{Version: "1.0.0", Definitions: definitions})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	definitions[0].Description = "mutated caller value"
	got, ok := registry.Lookup("planning.report_generate")
	if !ok || got.Description == "mutated caller value" {
		t.Fatalf("registry retained mutable caller input: %+v", got)
	}

	snapshot := registry.Definitions()
	if gotCodes := []Code{snapshot[0].Code, snapshot[1].Code}; !reflect.DeepEqual(gotCodes, []Code{"parcel.lookup", "planning.report_generate"}) {
		t.Fatalf("Definitions() codes = %v", gotCodes)
	}

	snapshot[0].Description = "mutated snapshot"
	got, _ = registry.Lookup("parcel.lookup")
	if got.Description == "mutated snapshot" {
		t.Fatal("registry snapshot mutated internal state")
	}
}

func TestRegistryValidatesLifecycle(t *testing.T) {
	t.Parallel()

	active := operation("parcel.resolve", "parcel")
	deprecated := operation("parcel.lookup", "parcel")
	deprecated.Status = StatusDeprecated
	deprecated.Replacement = active.Code

	registry, err := New(Catalog{
		Version:     "1.1.0",
		Definitions: []Definition{deprecated, active},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if _, err := registry.RequireActive(deprecated.Code); !errors.Is(err, ErrOperationNotActive) {
		t.Fatalf("RequireActive() error = %v", err)
	}
	if got, err := registry.RequireActive(active.Code); err != nil || got.Code != active.Code {
		t.Fatalf("RequireActive() = %+v, %v", got, err)
	}
}

func TestRegistryRejectsInvalidCatalogs(t *testing.T) {
	t.Parallel()

	t.Run("duplicate", func(t *testing.T) {
		definition := operation("parcel.lookup", "parcel")
		_, err := New(Catalog{Version: "1.0.0", Definitions: []Definition{definition, definition}})
		if err == nil {
			t.Fatal("New() accepted duplicate code")
		}
	})

	t.Run("missing replacement", func(t *testing.T) {
		definition := operation("parcel.lookup", "parcel")
		definition.Status = StatusDeprecated
		definition.Replacement = "parcel.resolve"
		_, err := New(Catalog{Version: "1.0.0", Definitions: []Definition{definition}})
		if err == nil {
			t.Fatal("New() accepted missing replacement")
		}
	})

	t.Run("replacement cycle", func(t *testing.T) {
		left := operation("parcel.lookup", "parcel")
		left.Status = StatusDeprecated
		left.Replacement = "parcel.resolve"
		right := operation("parcel.resolve", "parcel")
		right.Status = StatusDeprecated
		right.Replacement = "parcel.lookup"
		_, err := New(Catalog{Version: "1.0.0", Definitions: []Definition{left, right}})
		if err == nil {
			t.Fatal("New() accepted replacement cycle")
		}
	})
}

func TestRegistryGroupsByModule(t *testing.T) {
	t.Parallel()

	registry, err := New(Catalog{
		Version: "1.0.0",
		Definitions: []Definition{
			operation("planning.report_generate", "planning"),
			operation("parcel.lookup", "parcel"),
			operation("parcel.export", "parcel"),
		},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	parcel := registry.ByModule("parcel")
	if got := []Code{parcel[0].Code, parcel[1].Code}; !reflect.DeepEqual(got, []Code{"parcel.export", "parcel.lookup"}) {
		t.Fatalf("ByModule() codes = %v", got)
	}
}

func operation(code Code, module ModuleCode) Definition {
	return Definition{
		Code:        code,
		ModuleOwner: module,
		FeatureCode: Code(string(module) + ".core"),
		Metering: Metering{
			MeterCode:        Code(string(code) + "_count"),
			ChargeMode:       ChargeModeCommitOnSuccess,
			SuccessCondition: SuccessBusinessOutcome,
			FailurePolicy:    FailureNoCharge,
		},
		Status:          StatusActive,
		SemanticVersion: "1.0.0",
		Description:     "Test operation definition.",
	}
}
