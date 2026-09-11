package quota

import commonoperation "common/operation"

const (
	AreaAnalysisOperation     commonoperation.Code = "workspace.area.analyze"
	ParcelDetailViewOperation commonoperation.Code = "workspace.parcel.detail.view"
	GeneratedReportOperation  commonoperation.Code = "workspace.report.generate"
)

// ProfileOperations is the complete set exposed in the customer usage
// projection. Public map tiles and public parcel previews are deliberately not
// commercial operations: opening or refreshing the map must never spend quota.
var ProfileOperations = []commonoperation.Code{
	AreaAnalysisOperation,
	ParcelDetailViewOperation,
	GeneratedReportOperation,
}
