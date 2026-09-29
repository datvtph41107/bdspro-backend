package service

import (
	_errors "common/errors"

	"google.golang.org/grpc/codes"
)

var (
	Unauthenticated             = _errors.MustSpec(300001, "SOCIAL_UNAUTHENTICATED", "Unauthorized", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	GroupIDRequired             = _errors.MustSpec(300002, "SOCIAL_GROUP_ID_REQUIRED", "groupId is required", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	CommentUnavailable          = _errors.MustSpec(300003, "SOCIAL_COMMENT_UNAVAILABLE", "bài viết không tồn tại hoặc bị giới hạn bình luận", codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	CommentCreateFieldsRequired = _errors.MustSpec(300004, "SOCIAL_COMMENT_CREATE_FIELDS_REQUIRED", "content and news feed id are required", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	CommentUpdateFieldsRequired = _errors.MustSpec(300005, "SOCIAL_COMMENT_UPDATE_FIELDS_REQUIRED", "content, news feed id and id are required", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	IDRequired                  = _errors.MustSpec(300006, "SOCIAL_ID_REQUIRED", "id is required", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OwnerTypeInvalid            = _errors.MustSpec(300007, "SOCIAL_OWNER_TYPE_INVALID", "ownerOf is user/group/organization", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OrganizationIDRequired      = _errors.MustSpec(300008, "SOCIAL_ORGANIZATION_ID_REQUIRED", "organizationId is required", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	UserIDRequired              = _errors.MustSpec(300009, "SOCIAL_USER_ID_REQUIRED", "userId is required", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PostForwardDenied           = _errors.MustSpec(
		300010, "SOCIAL_POST_FORWARD_DENIED", "bạn không có quyền tạo/forward tin đăng này",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	FriendTagDenied = _errors.MustSpec(
		300011, "SOCIAL_FRIEND_TAG_DENIED", "not friends",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	NewsFeedUpdateDenied = _errors.MustSpec(
		300012, "SOCIAL_NEWS_FEED_UPDATE_DENIED", "bài viết không tồn tại hoặc bạn không có quyền cập nhật",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	NewsFeedDeleteDenied = _errors.MustSpec(
		300013, "SOCIAL_NEWS_FEED_DELETE_DENIED", "bạn không có quyền xóa bài viết này",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	AdminAccessRequired = _errors.MustSpec(
		300014, "SOCIAL_ADMIN_ACCESS_REQUIRED", "Bạn không có quyền admin",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	ShareFieldsRequired = _errors.MustSpec(
		300015, "SOCIAL_SHARE_FIELDS_REQUIRED", "newsFeedId or shareType is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
)
