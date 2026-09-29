package service

import (
	_errors "common/errors"

	"google.golang.org/grpc/codes"
)

var (
	UserGuideNotFound = _errors.MustSpec(
		420001, "HUB_USER_GUIDE_NOT_FOUND", "Không tìm thấy user guide",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	IDRequired = _errors.MustSpec(
		420002, "HUB_ID_REQUIRED", "Thiếu ID",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	OwnerIDRequired = _errors.MustSpec(
		420003, "HUB_OWNER_ID_REQUIRED", "owner_id is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	SystemConfigGroupInvalid = _errors.MustSpec(
		420004, "HUB_SYSTEM_CONFIG_GROUP_INVALID", "Group key không hợp lệ",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	SystemConfigItemGroupInvalid = _errors.MustSpec(
		420005, "HUB_SYSTEM_CONFIG_ITEM_GROUP_INVALID", "Group key không hợp lệ cho config",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	SystemConfigDefaultNotFound = _errors.MustSpec(
		420006, "HUB_SYSTEM_CONFIG_DEFAULT_NOT_FOUND", "Không tìm thấy config mặc định cho group",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	SystemConfigKeyRequired = _errors.MustSpec(
		420007, "HUB_SYSTEM_CONFIG_KEY_REQUIRED", "Key không được để trống",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	SystemConfigNotFound = _errors.MustSpec(
		420008, "HUB_SYSTEM_CONFIG_NOT_FOUND", "Không tìm thấy config",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	ErrorLogMessageRequired = _errors.MustSpec(
		420009, "HUB_ERROR_LOG_MESSAGE_REQUIRED", "message là bắt buộc",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
)
