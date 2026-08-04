package operationcatalog

import _operationregistry "common/operationregistry"

const (
	WorkspaceReportGenerateOperation _operationregistry.Code = "workspace.report.generate"
	WorkspaceReportAcceptedMeter     _operationregistry.Code = "workspace.report_generation.accepted"
)

/** TQDWorkspaceV1 returns approved workspace operations. */
func TQDWorkspaceV1() _operationregistry.Catalog {
	return _operationregistry.Catalog{
		Version: "1.0.0",
		Definitions: []_operationregistry.Definition{
			{
				Code:        WorkspaceReportGenerateOperation,
				ModuleOwner: "workspace",
				FeatureCode: "workspace.generated_reports",
				Metering: _operationregistry.Metering{
					MeterCode:        WorkspaceReportAcceptedMeter,
					ChargeMode:       _operationregistry.ChargeModeReserveThenCommit,
					SuccessCondition: _operationregistry.SuccessJobAccepted,
					FailurePolicy:    _operationregistry.FailureCancelReservation,
				},
				Status:          _operationregistry.StatusActive,
				SemanticVersion: "1.0.0",
				Description:     "Accept one authenticated request to start generating a workspace report.",
			},
		},
	}
}
