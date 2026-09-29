package service

import (
	_errors "common/errors"

	"google.golang.org/grpc/codes"
)

var (
	PermissionSnapshotUnavailable = _errors.MustSpec(200001, "AUTH_PERMISSION_SNAPSHOT_UNAVAILABLE", "permission snapshot chưa sẵn sàng hoặc đã quá hạn", codes.Unavailable, _errors.LegacyHTTP200(), _errors.LegacyCode(503))
	PermissionIDNotFound          = _errors.MustSpec(200002, "AUTH_PERMISSION_ID_NOT_FOUND", "permissionId không tồn tại trong snapshot", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PermissionCodeInvalid         = _errors.MustSpec(200003, "AUTH_PERMISSION_CODE_INVALID", "permission code không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PermissionCodeDenied          = _errors.MustSpec(200004, "AUTH_PERMISSION_CODE_DENIED", "permission code không tồn tại trong snapshot", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	RoleAssignmentUnavailable     = _errors.MustSpec(200005, "AUTH_ROLE_ASSIGNMENT_UNAVAILABLE", "không thể xác minh role assignment hiện tại", codes.Unavailable, _errors.LegacyHTTP200(), _errors.LegacyCode(503))
	PermissionMatchModeInvalid    = _errors.MustSpec(200006, "AUTH_PERMISSION_MATCH_MODE_INVALID", "permission match mode không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	Unauthenticated               = _errors.MustSpec(200007, "AUTH_UNAUTHENTICATED", "Unauthorized", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	PermissionDenied              = _errors.MustSpec(200008, "AUTH_PERMISSION_DENIED", "Bạn không có quyền truy cập", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
)
