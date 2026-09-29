package service

import (
	_errors "common/errors"

	"google.golang.org/grpc/codes"
)

var (
	OwnerIDRequired             = _errors.MustSpec(400001, "BDSPRO_OWNER_ID_REQUIRED", "owner_id is required", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PostQuotaExceeded           = _errors.MustSpec(400002, "BDSPRO_POST_QUOTA_EXCEEDED", "Số lượng tạo tin đã đạt giới hạn. Vui lòng nâng cấp gói", codes.ResourceExhausted, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PostStillActive             = _errors.MustSpec(400003, "BDSPRO_POST_STILL_ACTIVE", "Tin đăng vẫn còn hạn", codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	DealNameRequired            = _errors.MustSpec(400004, "BDSPRO_DEAL_NAME_REQUIRED", "name is required", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	InvalidPagable              = _errors.MustSpec(400005, "BDSPRO_PAGABLE_INVALID", "invalid pagable", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PropertyNameRequired        = _errors.MustSpec(400006, "BDSPRO_PROPERTY_NAME_REQUIRED", "Tên tính chất không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PropertyTypeNameRequired    = _errors.MustSpec(400007, "BDSPRO_PROPERTY_TYPE_NAME_REQUIRED", "Tên loại không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ProjectNameRequired         = _errors.MustSpec(400008, "BDSPRO_PROJECT_NAME_REQUIRED", "Tên dự án không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	DeveloperIDRequired         = _errors.MustSpec(400009, "BDSPRO_DEVELOPER_ID_REQUIRED", "Developer ID không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AssetNameRequired           = _errors.MustSpec(400010, "BDSPRO_ASSET_NAME_REQUIRED", "Tên tài sản không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	DealInvitationAlreadyExists = _errors.MustSpec(400011, "BDSPRO_DEAL_INVITATION_ALREADY_EXISTS", "invitation already exists for this user", codes.AlreadyExists, _errors.LegacyProblemCode("bdspro.deal_invitation.already_exists"))
	ProductChildSplitOnly       = _errors.MustSpec(
		400012, "BDSPRO_PRODUCT_CHILD_SPLIT_ONLY", "Chia sản phẩm con chỉ áp dụng đối với sản phẩm con",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	AssetChildMergeOnly = _errors.MustSpec(
		400013, "BDSPRO_ASSET_CHILD_MERGE_ONLY", "Chỉ gộp tài sản con",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	AreaInvalid = _errors.MustSpec(
		400014, "BDSPRO_AREA_INVALID", "Diện tích không hợp lệ",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ProductDepositTransitionDenied = _errors.MustSpec(
		400015, "BDSPRO_PRODUCT_DEPOSIT_TRANSITION_DENIED", "Không thể cập nhật trạng thái sản phẩm thành cọc",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	AssetSplitBlockedByProductLinks = _errors.MustSpec(
		400016, "BDSPRO_ASSET_SPLIT_BLOCKED_BY_PRODUCT_LINKS", "Không thể tách do tài sản đang liên kết với nhiều sản phẩm",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ProductChildSplitIneligible = _errors.MustSpec(
		400017, "BDSPRO_PRODUCT_CHILD_SPLIT_INELIGIBLE", "Sản phẩm con không đạt điều kiện để chia",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PostAlreadyPublished = _errors.MustSpec(
		400018, "BDSPRO_POST_ALREADY_PUBLISHED", "Tin đăng đã được đăng",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ProductStatusUnchanged = _errors.MustSpec(
		400019, "BDSPRO_PRODUCT_STATUS_UNCHANGED", "Trạng thái trùng với trạng thái cũ",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ProductsAlreadyLinkedToAsset = _errors.MustSpec(
		400020, "BDSPRO_PRODUCTS_ALREADY_LINKED_TO_ASSET", "Tất cả các sản phẩm đã được liên kết với tài sản này",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ProductChildMergeIneligible = _errors.MustSpec(
		400021, "BDSPRO_PRODUCT_CHILD_MERGE_INELIGIBLE", "Tồn tại sản phẩm con không đạt điều kiện để gộp",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	AccessDenied = _errors.MustSpec(
		400022, "BDSPRO_ACCESS_DENIED", "Bạn không có quyền truy cập",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(401),
	)
	ProductDeleteDenied = _errors.MustSpec(
		400023, "BDSPRO_PRODUCT_DELETE_DENIED", "Không thể xóa sản phẩm này",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(401),
	)
	ParentProductNotFound = _errors.MustSpec(
		400024, "BDSPRO_PARENT_PRODUCT_NOT_FOUND", "Không tìm thấy sản phẩm cha",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	DealNotFound = _errors.MustSpec(
		400025, "BDSPRO_DEAL_NOT_FOUND", "Không tìm thấy thương vụ",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	AssetNotFound = _errors.MustSpec(
		400026, "BDSPRO_ASSET_NOT_FOUND", "Tài sản không tồn tại",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	ProductNotFound = _errors.MustSpec(
		400027, "BDSPRO_PRODUCT_NOT_FOUND", "Sản phẩm không tồn tại",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	ProductOwnershipDenied = _errors.MustSpec(
		400028, "BDSPRO_PRODUCT_OWNERSHIP_DENIED", "Sản phẩm không thuộc quyền sở hữu của bạn",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	AssetOwnershipDenied = _errors.MustSpec(
		400029, "BDSPRO_ASSET_OWNERSHIP_DENIED", "Bạn không có quyền truy cập",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(401),
	)
	AssetOwnershipOrTradeDenied = _errors.MustSpec(
		400030, "BDSPRO_ASSET_OWNERSHIP_OR_TRADE_DENIED", "Bạn không sở hữu tài sản hoặc tài sản đang trong giao dịch chưa hoàn tất",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(401),
	)
	DistributionRevoked = _errors.MustSpec(
		400031, "BDSPRO_DISTRIBUTION_REVOKED", "distribution has been revoked",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	DealNotModifiable = _errors.MustSpec(
		400032, "BDSPRO_DEAL_NOT_MODIFIABLE", "deal is not modifiable",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(409),
	)
	DealProductAlreadyExists = _errors.MustSpec(
		400033, "BDSPRO_DEAL_PRODUCT_ALREADY_EXISTS", "product already exists in deal",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409),
	)
	PropertyReportAlreadyPending = _errors.MustSpec(
		400034, "BDSPRO_PROPERTY_REPORT_ALREADY_PENDING", "you already have a pending report for this property",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409),
	)
	PropertyAlreadyAdded = _errors.MustSpec(
		400035, "BDSPRO_PROPERTY_ALREADY_ADDED", "Đã thêm BĐS này",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409),
	)
	ProductPriceNotChannel = _errors.MustSpec(
		400036, "BDSPRO_PRODUCT_PRICE_NOT_CHANNEL", "Price is not channel price",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	ProductPriceNotSet = _errors.MustSpec(
		400037, "BDSPRO_PRODUCT_PRICE_NOT_SET", "Product price has not been set",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	ProductAccessDenied = _errors.MustSpec(
		400038, "BDSPRO_PRODUCT_ACCESS_DENIED", "User has no access to product",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	ProductSourceUpdateDenied = _errors.MustSpec(
		400039, "BDSPRO_PRODUCT_SOURCE_UPDATE_DENIED", "You do not have permission to update this product source",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	PropertyUpdateDenied = _errors.MustSpec(
		400040, "BDSPRO_PROPERTY_UPDATE_DENIED", "You don't have permission to update this property",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	ProductNoteAccessDenied = _errors.MustSpec(
		400041, "BDSPRO_PRODUCT_NOTE_ACCESS_DENIED", "permission denied",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	DealMembershipDenied = _errors.MustSpec(
		400042, "BDSPRO_DEAL_MEMBERSHIP_DENIED", "user is not member of deal",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	DistributionPriceNotFound = _errors.MustSpec(
		400043, "BDSPRO_DISTRIBUTION_PRICE_NOT_FOUND", "Distribute price not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	DistributionNotFound = _errors.MustSpec(
		400044, "BDSPRO_DISTRIBUTION_NOT_FOUND", "Distribution not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	PropertyNotFound = _errors.MustSpec(
		400045, "BDSPRO_PROPERTY_NOT_FOUND", "Property not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	ProductPriceNotFound = _errors.MustSpec(
		400046, "BDSPRO_PRODUCT_PRICE_NOT_FOUND", "Price not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	ProductUserRelationNotFound = _errors.MustSpec(
		400047, "BDSPRO_PRODUCT_USER_RELATION_NOT_FOUND", "Product user relation not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	PropertyLineageNotFound = _errors.MustSpec(
		400048, "BDSPRO_PROPERTY_LINEAGE_NOT_FOUND", "lineage not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	ProductNoteNotFound = _errors.MustSpec(
		400049, "BDSPRO_PRODUCT_NOTE_NOT_FOUND", "note not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
)
