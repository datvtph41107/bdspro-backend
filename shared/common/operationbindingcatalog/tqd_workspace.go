package operationbindingcatalog

import (
	_operationbinding "common/operationbinding"
	_operationcatalog "common/operationcatalog"
	_operationregistry "common/operationregistry"
)

const WorkspaceReportGenerateBinding _operationregistry.Code = "workspace.report.generate_transport"

func workspaceReportGenerateTransports() []_operationbinding.Transport {
	return []_operationbinding.Transport{
		_operationbinding.GRPCUnary("/types.tqd.MapWorkspaceService/CreateGeneratedReport"),
		_operationbinding.HTTPRoute("POST", "/v2/tqd/map-workspace/reports"),
	}
}

/** TQDWorkspaceV1 returns approved workspace transport bindings. */
func TQDWorkspaceV1() _operationbinding.Catalog {
	return _operationbinding.Catalog{
		Version: "1.0.0",
		Bindings: []_operationbinding.Binding{
			{
				Code:          WorkspaceReportGenerateBinding,
				ModuleOwner:   "workspace",
				Transports:    workspaceReportGenerateTransports(),
				Disposition:   _operationbinding.DispositionDirect,
				OperationCode: _operationcatalog.WorkspaceReportGenerateOperation,
				Description:   "CreateGeneratedReport always represents one accepted report-generation request.",
			},
		},
	}
}

/** TQDWorkspaceRequiredTransportsV1 returns startup coverage. */
func TQDWorkspaceRequiredTransportsV1() []_operationbinding.Transport {
	return workspaceReportGenerateTransports()
}
