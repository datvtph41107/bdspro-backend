package service

import (
	_errors "common/errors"

	"google.golang.org/grpc/codes"
)

var (
	AdminAccessRequired = _errors.MustSpec(
		430001, "NOTIFICATION_ADMIN_ACCESS_REQUIRED", "Không có quyền truy cập",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(401),
	)
	SelfWarningNotAllowed = _errors.MustSpec(
		430002, "NOTIFICATION_SELF_WARNING_NOT_ALLOWED", "Không thể gửi cảnh báo cho chính tài khoản của mình",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	WarningContentTooShort = _errors.MustSpec(
		430003, "NOTIFICATION_WARNING_CONTENT_TOO_SHORT", "Nội dung cảnh báo phải có ít nhất 20 ký tự",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	DuplicateWarningNotAllowed = _errors.MustSpec(
		430004, "NOTIFICATION_DUPLICATE_WARNING_NOT_ALLOWED", "Không được gửi 2 cảnh báo trùng nội dung trong vòng 24 giờ",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	WarningNotFound = _errors.MustSpec(
		430005, "NOTIFICATION_WARNING_NOT_FOUND", "Không tìm thấy cảnh báo",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	WarningReadDenied = _errors.MustSpec(
		430006, "NOTIFICATION_WARNING_READ_DENIED", "Không có quyền đánh dấu đã đọc cảnh báo này",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	WarningAcknowledgeDenied = _errors.MustSpec(
		430007, "NOTIFICATION_WARNING_ACKNOWLEDGE_DENIED", "Không có quyền xác nhận cảnh báo này",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	WarningTemplateNotFound = _errors.MustSpec(
		430008, "NOTIFICATION_WARNING_TEMPLATE_NOT_FOUND", "Không tìm thấy mẫu cảnh báo",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	WarningLogSelectorRequired = _errors.MustSpec(
		430009, "NOTIFICATION_WARNING_LOG_SELECTOR_REQUIRED", "Cần cung cấp WarningID hoặc TargetID",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	UserIDInvalid = _errors.MustSpec(
		430010, "NOTIFICATION_USER_ID_INVALID", "userId không hợp lệ",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	NotificationOwnerIdentityRequired = _errors.MustSpec(
		430011, "NOTIFICATION_OWNER_IDENTITY_REQUIRED", "Không thể tạo thông báo",
		codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401),
	)
	NotificationMetadataRequired = _errors.MustSpec(
		430012, "NOTIFICATION_METADATA_REQUIRED", "Thông tin không chính xác",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PropertyHistoryNotFound = _errors.MustSpec(
		430013, "NOTIFICATION_PROPERTY_HISTORY_NOT_FOUND", "History not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
)
