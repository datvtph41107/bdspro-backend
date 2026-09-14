package usecase

import "common/fault"

func workspaceUserIDRequired() error {
	return fault.Validation(
		"tqd.workspace.user_id_required",
		"user_id is required",
		fault.FieldViolation{
			Field:       "user_id",
			Description: "is required",
		},
	)
}

func workspaceParcelIDRequired() error {
	return fault.Validation(
		"tqd.workspace.parcel_id_required",
		"parcel_id is required",
		fault.FieldViolation{
			Field:       "parcel_id",
			Description: "is required",
		},
	)
}

func workspaceFollowOrParcelIDRequired() error {
	return fault.Validation(
		"tqd.workspace.follow_or_parcel_id_required",
		"follow_id or parcel_id is required",
	)
}

func workspaceEntityIDRequired() error {
	return fault.Validation(
		"tqd.workspace.entity_id_required",
		"entity_id is required",
		fault.FieldViolation{
			Field:       "entity_id",
			Description: "is required",
		},
	)
}

func workspaceHistoryIDRequired() error {
	return fault.Validation(
		"tqd.workspace.history_id_required",
		"history_id is required",
		fault.FieldViolation{
			Field:       "history_id",
			Description: "is required",
		},
	)
}

func workspaceReportIDRequired() error {
	return fault.Validation(
		"tqd.workspace.report_id_required",
		"report_id is required",
		fault.FieldViolation{
			Field:       "report_id",
			Description: "is required",
		},
	)
}

func workspaceParcelNotFound() error {
	return fault.New(
		fault.KindNotFound,
		"tqd.workspace.parcel_not_found",
		"parcel not found",
	)
}

func workspaceRegionNotFound() error {
	return fault.New(
		fault.KindNotFound,
		"tqd.workspace.region_not_found",
		"region not found",
	)
}

func workspaceReportNotFound() error {
	return fault.New(
		fault.KindNotFound,
		"tqd.workspace.report_not_found",
		"report not found",
	)
}

func workspaceReportRegenerationNotAllowed() error {
	return fault.New(
		fault.KindPrecondition,
		"tqd.workspace.report_regeneration_not_allowed",
		"report cannot regenerate in current status",
	)
}

func workspaceReportShareNotReady() error {
	return fault.New(
		fault.KindPrecondition,
		"tqd.workspace.report_share_not_ready",
		"report is not ready to share",
	)
}
