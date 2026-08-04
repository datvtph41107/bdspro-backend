package operationregistry

import "testing"

func TestDefinitionValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		definition Definition
		wantError  bool
	}{
		{
			name: "active metered definition",
			definition: Definition{
				Code:        "parcel.lookup",
				ModuleOwner: "parcel",
				FeatureCode: "parcel.detail",
				Metering: Metering{
					MeterCode:        "parcel.lookup_count",
					ChargeMode:       ChargeModeReserveThenCommit,
					SuccessCondition: SuccessBusinessOutcome,
					FailurePolicy:    FailureCancelReservation,
				},
				Status:          StatusActive,
				SemanticVersion: "1.0.0",
				Description:     "Resolves one parcel lookup business outcome.",
			},
		},
		{
			name: "unmetered definition",
			definition: Definition{
				Code:        "system.health_read",
				ModuleOwner: "system",
				FeatureCode: "system.health",
				Metering: Metering{
					ChargeMode:       ChargeModeUnmetered,
					SuccessCondition: SuccessHandlerSucceeded,
					FailurePolicy:    FailureNoCharge,
				},
				Status:          StatusActive,
				SemanticVersion: "1.0.0",
				Description:     "Reads service health without commercial metering.",
			},
		},
		{
			name: "technical style is rejected",
			definition: Definition{
				Code:        "FindParcelByLocation",
				ModuleOwner: "parcel",
				FeatureCode: "parcel.detail",
				Metering: Metering{
					ChargeMode:       ChargeModeUnmetered,
					SuccessCondition: SuccessHandlerSucceeded,
					FailurePolicy:    FailureNoCharge,
				},
				Status:          StatusActive,
				SemanticVersion: "1.0.0",
				Description:     "Invalid technical name.",
			},
			wantError: true,
		},
		{
			name: "metered operation requires meter",
			definition: Definition{
				Code:        "parcel.lookup",
				ModuleOwner: "parcel",
				FeatureCode: "parcel.detail",
				Metering: Metering{
					ChargeMode:       ChargeModeCommitOnSuccess,
					SuccessCondition: SuccessBusinessOutcome,
					FailurePolicy:    FailureNoCharge,
				},
				Status:          StatusActive,
				SemanticVersion: "1.0.0",
				Description:     "Missing meter.",
			},
			wantError: true,
		},
		{
			name: "active replacement is rejected",
			definition: Definition{
				Code:        "parcel.lookup",
				ModuleOwner: "parcel",
				FeatureCode: "parcel.detail",
				Metering: Metering{
					ChargeMode:       ChargeModeUnmetered,
					SuccessCondition: SuccessHandlerSucceeded,
					FailurePolicy:    FailureNoCharge,
				},
				Status:          StatusActive,
				SemanticVersion: "1.0.0",
				Replacement:     "parcel.resolve",
				Description:     "Invalid active replacement.",
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.definition.Validate()
			if (err != nil) != test.wantError {
				t.Fatalf("Validate() error = %v, wantError = %v", err, test.wantError)
			}
		})
	}
}
