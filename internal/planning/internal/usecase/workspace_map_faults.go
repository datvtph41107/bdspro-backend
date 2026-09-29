package usecase

import (
	_errors "common/errors"
	"tqd/internal"
)

func workspaceUserIDRequired() error {
	return _errors.ReturnError(service.WorkspaceUserIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "user_id", Description: "is required"}))
}
func workspaceParcelIDRequired() error {
	return _errors.ReturnError(service.WorkspaceParcelIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "parcel_id", Description: "is required"}))
}
func workspaceFollowOrParcelIDRequired() error {
	return _errors.ReturnError(service.WorkspaceFollowOrParcelIDRequired)
}
func workspaceEntityIDRequired() error {
	return _errors.ReturnError(service.WorkspaceEntityIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "entity_id", Description: "is required"}))
}
func workspaceHistoryIDRequired() error {
	return _errors.ReturnError(service.WorkspaceHistoryIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "history_id", Description: "is required"}))
}
func workspaceReportIDRequired() error {
	return _errors.ReturnError(service.WorkspaceReportIDRequired, _errors.WithViolations(_errors.FieldViolation{Field: "report_id", Description: "is required"}))
}
func workspaceParcelNotFound() error { return _errors.ReturnError(service.WorkspaceParcelNotFound) }
func workspaceRegionNotFound() error { return _errors.ReturnError(service.WorkspaceRegionNotFound) }
func workspaceReportNotFound() error { return _errors.ReturnError(service.WorkspaceReportNotFound) }
func workspaceReportRegenerationNotAllowed() error {
	return _errors.ReturnError(service.WorkspaceReportRegenerationNotAllowed)
}
func workspaceReportShareNotReady() error {
	return _errors.ReturnError(service.WorkspaceReportShareNotReady)
}
