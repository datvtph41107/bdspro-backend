package service

import (
	_errors "common/errors"

	"google.golang.org/grpc/codes"
)

var (
	Unauthenticated = _errors.MustSpec(
		800001, "TQD_UNAUTHENTICATED", "unauthorized",
		codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401),
	)

	DirectoryCategoryNameRequired = _errors.MustSpec(
		800002, "TQD_DIRECTORY_CATEGORY_NAME_REQUIRED", "name is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	DirectoryCategoryNameTooLong = _errors.MustSpec(
		800003, "TQD_DIRECTORY_CATEGORY_NAME_TOO_LONG", "Name must not exceed 255 characters",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	DirectoryCategoryCodeRequired = _errors.MustSpec(
		800004, "TQD_DIRECTORY_CATEGORY_CODE_REQUIRED", "code is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	DirectoryCategoryCodeTooLong = _errors.MustSpec(
		800005, "TQD_DIRECTORY_CATEGORY_CODE_TOO_LONG", "Code must not exceed 50 characters",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	DirectoryCategoryCodeFormatInvalid = _errors.MustSpec(
		800006, "TQD_DIRECTORY_CATEGORY_CODE_FORMAT_INVALID", "Code must contain only uppercase letters, numbers, and underscores",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	DirectoryCategoryDescriptionTooLong = _errors.MustSpec(
		800007, "TQD_DIRECTORY_CATEGORY_DESCRIPTION_TOO_LONG", "Description must not exceed 500 characters",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)

	PlanningNameRequired = _errors.MustSpec(
		800008, "TQD_PLANNING_NAME_REQUIRED", "name is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PlanningCodeRequired = _errors.MustSpec(
		800009, "TQD_PLANNING_CODE_REQUIRED", "code is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PlanningTypeRequired = _errors.MustSpec(
		800010, "TQD_PLANNING_TYPE_REQUIRED", "planningType is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PlanningLevelRequired = _errors.MustSpec(
		800011, "TQD_PLANNING_LEVEL_REQUIRED", "planningLevel is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PlanningValidityStatusRequired = _errors.MustSpec(
		800012, "TQD_PLANNING_VALIDITY_STATUS_REQUIRED", "validityStatus is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PlanningProjectIDRequired = _errors.MustSpec(
		800013, "TQD_PLANNING_PROJECT_ID_REQUIRED", "planningProjectId is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PlanningTitleRequired = _errors.MustSpec(
		800014, "TQD_PLANNING_TITLE_REQUIRED", "title is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PlanningDocumentTypeRequired = _errors.MustSpec(
		800015, "TQD_PLANNING_DOCUMENT_TYPE_REQUIRED", "documentType is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)

	ParcelIDRequired = _errors.MustSpec(
		800016, "TQD_PARCEL_ID_REQUIRED", "parcel_id is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	RegionIDRequired = _errors.MustSpec(
		800017, "TQD_REGION_ID_REQUIRED", "region_id is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	SubscriptionIDRequired = _errors.MustSpec(
		800018, "TQD_SUBSCRIPTION_ID_REQUIRED", "subscription_id is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	NotificationIDRequired = _errors.MustSpec(
		800019, "TQD_NOTIFICATION_ID_REQUIRED", "notification_id is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ParcelNotFound = _errors.MustSpec(
		800020, "TQD_PARCEL_NOT_FOUND", "parcel not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	RegionNotFound = _errors.MustSpec(
		800021, "TQD_REGION_NOT_FOUND", "region not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	SubscriptionAlreadyExists = _errors.MustSpec(
		800022, "TQD_SUBSCRIPTION_ALREADY_EXISTS", "subscription already exists",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	SubscriptionScopeInvalid = _errors.MustSpec(
		800023, "TQD_SUBSCRIPTION_SCOPE_INVALID", "invalid subscription scope",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	SubscriptionNotFound = _errors.MustSpec(
		800024, "TQD_SUBSCRIPTION_NOT_FOUND", "subscription not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	SubscriptionPermissionDenied = _errors.MustSpec(
		800025, "TQD_SUBSCRIPTION_PERMISSION_DENIED", "permission denied",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)

	ReportIDRequired = _errors.MustSpec(
		800026, "TQD_REPORT_ID_REQUIRED", "report_id is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ReportNotFound = _errors.MustSpec(
		800027, "TQD_REPORT_NOT_FOUND", "report not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	ReportOperationInvalid = _errors.MustSpec(
		800028, "TQD_REPORT_OPERATION_INVALID", "report operation is invalid",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ReportPermissionDenied = _errors.MustSpec(
		800029, "TQD_REPORT_PERMISSION_DENIED", "permission denied",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	ReportTargetDataInvalid = _errors.MustSpec(
		800030, "TQD_REPORT_TARGET_DATA_INVALID", "invalid targetData json",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)

	OneHouseIDRequired = _errors.MustSpec(
		800031, "TQD_ONE_HOUSE_ID_REQUIRED", "id is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	OneHouseNotFound = _errors.MustSpec(
		800032, "TQD_ONE_HOUSE_NOT_FOUND", "one_house not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	AdminReportIDInvalid = _errors.MustSpec(
		800033, "TQD_ADMIN_REPORT_ID_INVALID", "invalid id",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PermissionAuthorityUnavailable = _errors.MustSpec(
		800034, "TQD_PERMISSION_AUTHORITY_UNAVAILABLE", "permission authority unavailable",
		codes.Unavailable, _errors.LegacyHTTP200(), _errors.LegacyCode(503),
	)
	ReportNotReady = _errors.MustSpec(
		800035, "TQD_REPORT_NOT_READY", "report not ready",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ReportAssigneeIDRequired = _errors.MustSpec(
		800036, "TQD_REPORT_ASSIGNEE_ID_REQUIRED", "assigneeId is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ReportStatusTransitionInvalid = _errors.MustSpec(
		800037, "TQD_REPORT_STATUS_TRANSITION_INVALID", "STATUS_TRANSITION_INVALID",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	OneHouseIdentifierRequired = _errors.MustSpec(
		800038, "TQD_ONE_HOUSE_IDENTIFIER_REQUIRED", "id or uuid is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	ImportLayerIDRequired                 = _errors.MustSpec(800039, "TQD_IMPORT_LAYER_ID_REQUIRED", "layerId is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.import.layer_id_required"))
	ImportErrorIDRequired                 = _errors.MustSpec(800040, "TQD_IMPORT_ERROR_ID_REQUIRED", "errorId is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.import.error_id_required"))
	ImportLayerNotFound                   = _errors.MustSpec(800041, "TQD_IMPORT_LAYER_NOT_FOUND", "layer not found", codes.NotFound, _errors.LegacyProblemCode("tqd.import.layer_not_found"))
	ImportInProgress                      = _errors.MustSpec(800042, "TQD_IMPORT_IN_PROGRESS", "import already in progress", codes.FailedPrecondition, _errors.LegacyProblemCode("tqd.import.in_progress"))
	ImportErrorNotFound                   = _errors.MustSpec(800043, "TQD_IMPORT_ERROR_NOT_FOUND", "import error not found", codes.NotFound, _errors.LegacyProblemCode("tqd.import.error_not_found"))
	ImportRetryInProgress                 = _errors.MustSpec(800044, "TQD_IMPORT_RETRY_IN_PROGRESS", "import retry already in progress", codes.FailedPrecondition, _errors.LegacyProblemCode("tqd.import.retry_in_progress"))
	ImportRetryNotAllowed                 = _errors.MustSpec(800045, "TQD_IMPORT_RETRY_NOT_ALLOWED", "import retry is not allowed", codes.FailedPrecondition, _errors.LegacyProblemCode("tqd.import.retry_not_allowed"))
	ImportRetryLockConflict               = _errors.MustSpec(800046, "TQD_IMPORT_RETRY_LOCK_CONFLICT", "import retry lock conflict", codes.Aborted, _errors.LegacyProblemCode("tqd.import.retry_lock_conflict"))
	RegionLayerIDRequired                 = _errors.MustSpec(800047, "TQD_REGION_LAYER_ID_REQUIRED", "layerId is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.region.layer_id_required"))
	RegionLabelIDRequired                 = _errors.MustSpec(800048, "TQD_REGION_LABEL_ID_REQUIRED", "labelId is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.region.label_id_required"))
	RegionRecordIDRequired                = _errors.MustSpec(800049, "TQD_REGION_RECORD_ID_REQUIRED", "id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.region.id_required"))
	RegionNameRequired                    = _errors.MustSpec(800050, "TQD_REGION_NAME_REQUIRED", "name is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.region.name_required"))
	RegionUpdateRequired                  = _errors.MustSpec(800051, "TQD_REGION_UPDATE_REQUIRED", "at least one field to update is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.region.update_required"))
	RegionNameEmpty                       = _errors.MustSpec(800052, "TQD_REGION_NAME_EMPTY", "name cannot be empty", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.region.name_empty"))
	RegionStatusInvalid                   = _errors.MustSpec(800053, "TQD_REGION_STATUS_INVALID", "invalid region status", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.region.status_invalid"))
	RegionGeometryInvalid                 = _errors.MustSpec(800054, "TQD_REGION_GEOMETRY_INVALID", "invalid geometry", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.region.geometry_invalid"))
	RegionLayerNotFound                   = _errors.MustSpec(800055, "TQD_REGION_LAYER_NOT_FOUND", "layer not found", codes.NotFound, _errors.LegacyProblemCode("tqd.region.layer_not_found"))
	RegionRecordNotFound                  = _errors.MustSpec(800056, "TQD_REGION_RECORD_NOT_FOUND", "region not found", codes.NotFound, _errors.LegacyProblemCode("tqd.region.not_found"))
	RegionLabelNotFound                   = _errors.MustSpec(800057, "TQD_REGION_LABEL_NOT_FOUND", "label not found", codes.NotFound, _errors.LegacyProblemCode("tqd.region.label_not_found"))
	RegionSyncReferenceNotFound           = _errors.MustSpec(800058, "TQD_REGION_SYNC_REFERENCE_NOT_FOUND", "region sync reference not found", codes.NotFound, _errors.LegacyProblemCode("tqd.region.sync_reference_not_found"))
	LegendPayloadRequired                 = _errors.MustSpec(800059, "TQD_LEGEND_PAYLOAD_REQUIRED", "payload is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.legend.payload_required"))
	LegendIDRequired                      = _errors.MustSpec(800060, "TQD_LEGEND_ID_REQUIRED", "id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.legend.id_required"))
	LegendLayerIDRequired                 = _errors.MustSpec(800061, "TQD_LEGEND_LAYER_ID_REQUIRED", "layerId is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.legend.layer_id_required"))
	LegendLabelIDRequired                 = _errors.MustSpec(800062, "TQD_LEGEND_LABEL_ID_REQUIRED", "labelId is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.legend.label_id_required"))
	LegendLayerNotFound                   = _errors.MustSpec(800063, "TQD_LEGEND_LAYER_NOT_FOUND", "layer not found", codes.NotFound, _errors.LegacyProblemCode("tqd.legend.layer_not_found"))
	LegendLabelNotFound                   = _errors.MustSpec(800064, "TQD_LEGEND_LABEL_NOT_FOUND", "label not found", codes.NotFound, _errors.LegacyProblemCode("tqd.legend.label_not_found"))
	LegendRecordNotFound                  = _errors.MustSpec(800065, "TQD_LEGEND_RECORD_NOT_FOUND", "legend record not found", codes.NotFound, _errors.LegacyProblemCode("tqd.legend.record_not_found"))
	LegendDuplicate                       = _errors.MustSpec(800066, "TQD_LEGEND_DUPLICATE", "legend already exists", codes.AlreadyExists, _errors.LegacyProblemCode("tqd.legend.duplicate"))
	LegendTypeInvalid                     = _errors.MustSpec(800067, "TQD_LEGEND_TYPE_INVALID", "invalid legend type", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.legend.legend_type_invalid"))
	LegendGeometryTypeInvalid             = _errors.MustSpec(800068, "TQD_LEGEND_GEOMETRY_TYPE_INVALID", "invalid geometry type", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.legend.geometry_type_invalid"))
	LandUsePayloadRequired                = _errors.MustSpec(800069, "TQD_LAND_USE_PAYLOAD_REQUIRED", "payload is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.land_use.payload_required"))
	LandUseIDRequired                     = _errors.MustSpec(800070, "TQD_LAND_USE_ID_REQUIRED", "id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.land_use.id_required"))
	LandUseLayerIDRequired                = _errors.MustSpec(800071, "TQD_LAND_USE_LAYER_ID_REQUIRED", "layerId is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.land_use.layer_id_required"))
	LandUseValueIDRequired                = _errors.MustSpec(800072, "TQD_LAND_USE_VALUE_ID_REQUIRED", "landUseId is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.land_use.land_use_id_required"))
	LandUseGroupIDRequired                = _errors.MustSpec(800073, "TQD_LAND_USE_GROUP_ID_REQUIRED", "landUseGroupId is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.land_use.group_id_required"))
	LandUseNameRequired                   = _errors.MustSpec(800074, "TQD_LAND_USE_NAME_REQUIRED", "name is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.land_use.name_required"))
	LandUseLayerNotFound                  = _errors.MustSpec(800075, "TQD_LAND_USE_LAYER_NOT_FOUND", "layer not found", codes.NotFound, _errors.LegacyProblemCode("tqd.land_use.layer_not_found"))
	LandUseNotFound                       = _errors.MustSpec(800076, "TQD_LAND_USE_NOT_FOUND", "land use not found", codes.NotFound, _errors.LegacyProblemCode("tqd.land_use.not_found"))
	LandUseRecordNotFound                 = _errors.MustSpec(800077, "TQD_LAND_USE_RECORD_NOT_FOUND", "land use record not found", codes.NotFound, _errors.LegacyProblemCode("tqd.land_use.record_not_found"))
	LandUseDuplicate                      = _errors.MustSpec(800078, "TQD_LAND_USE_DUPLICATE", "land use already exists in layer", codes.AlreadyExists, _errors.LegacyProblemCode("tqd.land_use.duplicate"))
	LabelSourceNotFound                   = _errors.MustSpec(800079, "TQD_QH_LABEL_SOURCE_NOT_FOUND", "source label not found", codes.NotFound, _errors.LegacyProblemCode("tqd.qh_label.source_not_found"))
	LabelSourceLayerMismatch              = _errors.MustSpec(800080, "TQD_QH_LABEL_SOURCE_LAYER_MISMATCH", "source label belongs to another layer", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_label.source_layer_mismatch"))
	LabelNameConflict                     = _errors.MustSpec(800081, "TQD_QH_LABEL_NAME_CONFLICT", "label name already exists in layer", codes.AlreadyExists, _errors.LegacyProblemCode("tqd.qh_label.name_conflict"))
	LayerFamilyPayloadRequired            = _errors.MustSpec(800082, "TQD_QH_LAYER_FAMILY_PAYLOAD_REQUIRED", "payload is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_layer_family.payload_required"))
	LayerFamilyNameRequired               = _errors.MustSpec(800083, "TQD_QH_LAYER_FAMILY_NAME_REQUIRED", "name is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_layer_family.name_required"))
	LayerFamilyIDRequired                 = _errors.MustSpec(800084, "TQD_QH_LAYER_FAMILY_ID_REQUIRED", "id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_layer_family.id_required"))
	LayerFamilyNotFound                   = _errors.MustSpec(800085, "TQD_QH_LAYER_FAMILY_NOT_FOUND", "layer family not found", codes.NotFound, _errors.LegacyProblemCode("tqd.qh_layer_family.not_found"))
	AuthorityPayloadRequired              = _errors.MustSpec(800086, "TQD_QH_AUTHORITY_ISSURING_PAYLOAD_REQUIRED", "payload is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_authority_issuring.payload_required"))
	AuthorityNameRequired                 = _errors.MustSpec(800087, "TQD_QH_AUTHORITY_ISSURING_NAME_REQUIRED", "name is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_authority_issuring.name_required"))
	AuthorityCodeRequired                 = _errors.MustSpec(800088, "TQD_QH_AUTHORITY_ISSURING_CODE_REQUIRED", "code is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_authority_issuring.code_required"))
	AuthorityIDRequired                   = _errors.MustSpec(800089, "TQD_QH_AUTHORITY_ISSURING_ID_REQUIRED", "id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_authority_issuring.id_required"))
	AuthorityNotFound                     = _errors.MustSpec(800090, "TQD_QH_AUTHORITY_ISSURING_NOT_FOUND", "authority issuring not found", codes.NotFound, _errors.LegacyProblemCode("tqd.qh_authority_issuring.not_found"))
	AuthorityCodeConflict                 = _errors.MustSpec(800091, "TQD_QH_AUTHORITY_ISSURING_CODE_CONFLICT", "authority issuring code already exists", codes.AlreadyExists, _errors.LegacyProblemCode("tqd.qh_authority_issuring.code_conflict"))
	WorkspaceUserIDRequired               = _errors.MustSpec(800092, "TQD_WORKSPACE_USER_ID_REQUIRED", "user_id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.workspace.user_id_required"))
	WorkspaceParcelIDRequired             = _errors.MustSpec(800093, "TQD_WORKSPACE_PARCEL_ID_REQUIRED", "parcel_id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.workspace.parcel_id_required"))
	WorkspaceFollowOrParcelIDRequired     = _errors.MustSpec(800094, "TQD_WORKSPACE_FOLLOW_OR_PARCEL_ID_REQUIRED", "follow_id or parcel_id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.workspace.follow_or_parcel_id_required"))
	WorkspaceEntityIDRequired             = _errors.MustSpec(800095, "TQD_WORKSPACE_ENTITY_ID_REQUIRED", "entity_id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.workspace.entity_id_required"))
	WorkspaceHistoryIDRequired            = _errors.MustSpec(800096, "TQD_WORKSPACE_HISTORY_ID_REQUIRED", "history_id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.workspace.history_id_required"))
	WorkspaceReportIDRequired             = _errors.MustSpec(800097, "TQD_WORKSPACE_REPORT_ID_REQUIRED", "report_id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.workspace.report_id_required"))
	WorkspaceParcelNotFound               = _errors.MustSpec(800098, "TQD_WORKSPACE_PARCEL_NOT_FOUND", "parcel not found", codes.NotFound, _errors.LegacyProblemCode("tqd.workspace.parcel_not_found"))
	WorkspaceRegionNotFound               = _errors.MustSpec(800099, "TQD_WORKSPACE_REGION_NOT_FOUND", "region not found", codes.NotFound, _errors.LegacyProblemCode("tqd.workspace.region_not_found"))
	WorkspaceReportNotFound               = _errors.MustSpec(800100, "TQD_WORKSPACE_REPORT_NOT_FOUND", "report not found", codes.NotFound, _errors.LegacyProblemCode("tqd.workspace.report_not_found"))
	WorkspaceReportRegenerationNotAllowed = _errors.MustSpec(800101, "TQD_WORKSPACE_REPORT_REGENERATION_NOT_ALLOWED", "report cannot regenerate in current status", codes.FailedPrecondition, _errors.LegacyProblemCode("tqd.workspace.report_regeneration_not_allowed"))
	WorkspaceReportShareNotReady          = _errors.MustSpec(800102, "TQD_WORKSPACE_REPORT_SHARE_NOT_READY", "report is not ready to share", codes.FailedPrecondition, _errors.LegacyProblemCode("tqd.workspace.report_share_not_ready"))
	DiscoveryCoordinateInvalid            = _errors.MustSpec(800103, "TQD_DISCOVERY_COORDINATE_INVALID", "invalid latitude/longitude", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.discovery.coordinate_invalid"))

	AuthorityRequestRequired        = _errors.MustSpec(800104, "TQD_QH_AUTHORITY_ISSURING_REQUEST_REQUIRED", "request is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_authority_issuring.request_required"))
	AuthorityUpdateFieldsRequired   = _errors.MustSpec(800105, "TQD_QH_AUTHORITY_ISSURING_UPDATE_FIELDS_REQUIRED", "at least one field to update is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_authority_issuring.update_fields_required"))
	LayerFamilyRequestRequired      = _errors.MustSpec(800106, "TQD_QH_LAYER_FAMILY_REQUEST_REQUIRED", "request is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_layer_family.request_required"))
	LayerFamilyUpdateFieldsRequired = _errors.MustSpec(800107, "TQD_QH_LAYER_FAMILY_UPDATE_FIELDS_REQUIRED", "at least one field to update is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_layer_family.update_fields_required"))
	LayerFamilyFamilyIDRequired     = _errors.MustSpec(800108, "TQD_QH_LAYER_FAMILY_FAMILY_ID_REQUIRED", "family_id is required", codes.InvalidArgument, _errors.LegacyProblemCode("tqd.qh_layer_family.family_id_required"))
	LayerFamilyBuildInProgress      = _errors.MustSpec(800109, "TQD_QH_LAYER_FAMILY_BUILD_IN_PROGRESS", "family tile build is already in progress", codes.ResourceExhausted, _errors.LegacyProblemCode("tqd.qh_layer_family.build_in_progress"))
	WorkspaceRecordNotFound         = _errors.MustSpec(800110, "TQD_WORKSPACE_RECORD_NOT_FOUND", "record not found", codes.NotFound, _errors.LegacyProblemCode("tqd.workspace.record_not_found"))
)
