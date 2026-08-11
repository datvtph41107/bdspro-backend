package accesspostgres

import "common/businessoperation"

/**
 * CatalogCodes nối business operation với các code đang lưu trong product catalog.
 *
 * Mapping này thuộc access, không thuộc G3 business operation.
 */
type CatalogCodes struct {
	FeatureCode string
	UsageCode   string
}

func catalogCodesForOperation(operation businessoperation.Code) (CatalogCodes, bool) {
	switch operation {
	case businessoperation.WorkspaceReportGenerate:
		return CatalogCodes{
			FeatureCode: "workspace.report.generate",
			UsageCode:   "workspace.report_generation.accepted",
		}, true
	default:
		return CatalogCodes{}, false
	}
}
