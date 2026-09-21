package service

import (
	_errors "common/errors"

	"google.golang.org/grpc/codes"
)

var (
	ReportAlreadySubmitted               = _errors.MustSpec(410001, "CRM_REPORT_ALREADY_SUBMITTED", "bạn đã gửi báo cáo cho đối tượng này", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ReportUpdateDenied                   = _errors.MustSpec(410002, "CRM_REPORT_UPDATE_DENIED", "Bạn không có quyền cập nhật báo cáo này", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	ReportDeleteDenied                   = _errors.MustSpec(410003, "CRM_REPORT_DELETE_DENIED", "Bạn không có quyền xóa báo cáo này", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	SelfContactNotAllowed                = _errors.MustSpec(410004, "CRM_SELF_CONTACT_NOT_ALLOWED", "Bạn không thể tạo liên hệ với chính mình", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ReportAlreadyRejected                = _errors.MustSpec(410005, "CRM_REPORT_ALREADY_REJECTED", "Báo cáo đã bị từ chối trước đó", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ReportAlreadyApproved                = _errors.MustSpec(410006, "CRM_REPORT_ALREADY_APPROVED", "Báo cáo đã được duyệt trước đó", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ReportNotFound                       = _errors.MustSpec(410007, "CRM_REPORT_NOT_FOUND", "Báo cáo không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SupportDepartmentInvalid             = _errors.MustSpec(410008, "CRM_SUPPORT_DEPARTMENT_INVALID", "bộ phận xử lý không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SupportTicketUpdateContentRequired   = _errors.MustSpec(410009, "CRM_SUPPORT_TICKET_UPDATE_CONTENT_REQUIRED", "Cần ít nhất tiêu đề hoặc mô tả để cập nhật", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SupportTicketClosedAssignmentDenied  = _errors.MustSpec(410010, "CRM_SUPPORT_TICKET_CLOSED_ASSIGNMENT_DENIED", "cannot assign closed ticket", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOCanonicalURLInvalid               = _errors.MustSpec(410011, "CRM_SEO_CANONICAL_URL_INVALID", "canonical URL không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOCanonicalURLAlreadyExists         = _errors.MustSpec(410012, "CRM_SEO_CANONICAL_URL_ALREADY_EXISTS", "canonicalUrl đã tồn tại", codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409))
	SEOCanonicalURLRequiredForGenerate   = _errors.MustSpec(410013, "CRM_SEO_CANONICAL_URL_REQUIRED_FOR_GENERATE", "canonicalUrl là bắt buộc để generate", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(422))
	SEOPublishedPageRequiredForGenerate  = _errors.MustSpec(410014, "CRM_SEO_PUBLISHED_PAGE_REQUIRED_FOR_GENERATE", "chỉ generate SEO page đã published", codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409))
	AdminOpportunityScopeDenied          = _errors.MustSpec(410015, "CRM_ADMIN_OPPORTUNITY_SCOPE_DENIED", "Cơ hội không thuộc phạm vi Admin", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	SEOGeneratedAtInvalid                = _errors.MustSpec(410016, "CRM_SEO_GENERATED_AT_INVALID", "generatedAt phải theo RFC3339 hoặc YYYY-MM-DD", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	IDInvalid                            = _errors.MustSpec(410017, "CRM_ID_INVALID", "id không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AdminIDInvalid                       = _errors.MustSpec(410018, "CRM_ADMIN_ID_INVALID", "ID không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOInternalLinkAlreadyExists         = _errors.MustSpec(410019, "CRM_SEO_INTERNAL_LINK_ALREADY_EXISTS", "internal link đã tồn tại", codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409))
	AdminOpportunityCloseResultInvalid   = _errors.MustSpec(410020, "CRM_ADMIN_OPPORTUNITY_CLOSE_RESULT_INVALID", "Kết quả đóng không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	CustomerNotFound                     = _errors.MustSpec(410021, "CRM_CUSTOMER_NOT_FOUND", "Khách hàng không tồn tại", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SupportTicketClosedRoutingDenied     = _errors.MustSpec(410022, "CRM_SUPPORT_TICKET_CLOSED_ROUTING_DENIED", "không thể chuyển tuyến ticket đã đóng", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OpportunityNotFound                  = _errors.MustSpec(410023, "CRM_OPPORTUNITY_NOT_FOUND", "Không tìm thấy cơ hội", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	ContactNotFound                      = _errors.MustSpec(410024, "CRM_CONTACT_NOT_FOUND", "Không tìm thấy liên hệ", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	LeadNotFound                         = _errors.MustSpec(410025, "CRM_LEAD_NOT_FOUND", "Lead không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	AdminContactScopeDenied              = _errors.MustSpec(410026, "CRM_ADMIN_CONTACT_SCOPE_DENIED", "Liên hệ không thuộc phạm vi Admin", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	ReportReasonNotFound                 = _errors.MustSpec(410027, "CRM_REPORT_REASON_NOT_FOUND", "lý do báo cáo không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	OpportunityCloseReasonRequired       = _errors.MustSpec(410028, "CRM_OPPORTUNITY_CLOSE_REASON_REQUIRED", "Lý do đóng là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOMetadataJSONInvalid               = _errors.MustSpec(410029, "CRM_SEO_METADATA_JSON_INVALID", "metadata phải là JSON hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOMetadataObjectRequired            = _errors.MustSpec(410030, "CRM_SEO_METADATA_OBJECT_REQUIRED", "metadata phải là JSON object", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OpportunityAssigneeInvalid           = _errors.MustSpec(410031, "CRM_OPPORTUNITY_ASSIGNEE_INVALID", "Người phụ trách không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOOriginURLRequired                 = _errors.MustSpec(410032, "CRM_SEO_ORIGIN_URL_REQUIRED", "originUrl không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OwnerIDRequired                      = _errors.MustSpec(410033, "CRM_OWNER_ID_REQUIRED", "ownerId là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OwnerTypeInvalid                     = _errors.MustSpec(410034, "CRM_OWNER_TYPE_INVALID", "OwnerOf is invalid", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOSitemapPageStateInvalid           = _errors.MustSpec(410035, "CRM_SEO_SITEMAP_PAGE_STATE_INVALID", "page trong sitemap phải published=true và isIndex=true", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOArchivedPageFlagsInvalid          = _errors.MustSpec(410036, "CRM_SEO_ARCHIVED_PAGE_FLAGS_INVALID", "page_status=archived không được published/index/sitemap", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOPublishedStatusFlagRequired       = _errors.MustSpec(410037, "CRM_SEO_PUBLISHED_STATUS_FLAG_REQUIRED", "page_status=published yêu cầu published=true", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOPageStatusInvalid                 = _errors.MustSpec(410038, "CRM_SEO_PAGE_STATUS_INVALID", "pageStatus không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOParcelNotFound                    = _errors.MustSpec(410039, "CRM_SEO_PARCEL_NOT_FOUND", "parcel không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEOParcelIDRequired                  = _errors.MustSpec(410040, "CRM_SEO_PARCEL_ID_REQUIRED", "parcelId không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOParentChildSame                   = _errors.MustSpec(410041, "CRM_SEO_PARENT_CHILD_SAME", "parentSeoId không được trùng childSeoId", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOParentChildRequired               = _errors.MustSpec(410042, "CRM_SEO_PARENT_CHILD_REQUIRED", "parentSeoId và childSeoId là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PhoneRequired                        = _errors.MustSpec(410043, "CRM_PHONE_REQUIRED", "phone is required", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOPublicDocumentRequired            = _errors.MustSpec(410044, "CRM_SEO_PUBLIC_DOCUMENT_REQUIRED", "public SEO document là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(422))
	SEOPublishedPageFieldsRequired       = _errors.MustSpec(410045, "CRM_SEO_PUBLISHED_PAGE_FIELDS_REQUIRED", "published page cần slug, canonicalUrl và title", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOPublishedAtInvalid                = _errors.MustSpec(410046, "CRM_SEO_PUBLISHED_AT_INVALID", "publishedAt phải theo RFC3339 hoặc YYYY-MM-DD", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEORefIDInvalid                      = _errors.MustSpec(410047, "CRM_SEO_REF_ID_INVALID", "refId không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEORefIDRequired                     = _errors.MustSpec(410048, "CRM_SEO_REF_ID_REQUIRED", "refId là bắt buộc khi refType khác 0", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOSourceRefIDRequired               = _errors.MustSpec(410049, "CRM_SEO_SOURCE_REF_ID_REQUIRED", "refId là bắt buộc với source entity", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEORefLastSyncedAtInvalid            = _errors.MustSpec(410050, "CRM_SEO_REF_LAST_SYNCED_AT_INVALID", "refLastSyncedAt phải theo RFC3339 hoặc YYYY-MM-DD", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEORefResolverUnsupported            = _errors.MustSpec(410051, "CRM_SEO_REF_RESOLVER_UNSUPPORTED", "refType hoặc resolver không được hỗ trợ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEORefTypeInvalid                    = _errors.MustSpec(410052, "CRM_SEO_REF_TYPE_INVALID", "refType không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RegistryInvalid                      = _errors.MustSpec(410053, "CRM_REGISTRY_INVALID", "registry không đúng", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RegistryNotFound                     = _errors.MustSpec(410054, "CRM_REGISTRY_NOT_FOUND", "registry không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	RegistryNameNotFound                 = _errors.MustSpec(410055, "CRM_REGISTRY_NAME_NOT_FOUND", "registry name không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEORenderStatusInvalid               = _errors.MustSpec(410056, "CRM_SEO_RENDER_STATUS_INVALID", "renderStatus không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RequestInvalid                       = _errors.MustSpec(410057, "CRM_REQUEST_INVALID", "request không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOSecretKeyInvalid                  = _errors.MustSpec(410058, "CRM_SEO_SECRET_KEY_INVALID", "secretKey không hợp lệ", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	SEOCanonicalContractMismatch         = _errors.MustSpec(410059, "CRM_SEO_CANONICAL_CONTRACT_MISMATCH", "SEO canonical contract mismatch", codes.Unavailable, _errors.LegacyHTTP200(), _errors.LegacyCode(503))
	SEODomainUnpublished                 = _errors.MustSpec(410060, "CRM_SEO_DOMAIN_UNPUBLISHED", "seo domain chưa publish", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEODomainIndexDenied                 = _errors.MustSpec(410061, "CRM_SEO_DOMAIN_INDEX_DENIED", "seo domain không cho index", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEODomainNotFound                    = _errors.MustSpec(410062, "CRM_SEO_DOMAIN_NOT_FOUND", "seo domain không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEOInternalLinkNotFound              = _errors.MustSpec(410063, "CRM_SEO_INTERNAL_LINK_NOT_FOUND", "seo internal link không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEOPublicRouteContractMissing        = _errors.MustSpec(410064, "CRM_SEO_PUBLIC_ROUTE_CONTRACT_MISSING", "SEO module chưa có public route contract", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOModuleUnsupported                 = _errors.MustSpec(410065, "CRM_SEO_MODULE_UNSUPPORTED", "SEO module không được hỗ trợ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOSourcePageAlreadyExists           = _errors.MustSpec(410066, "CRM_SEO_SOURCE_PAGE_ALREADY_EXISTS", "SEO page cho source này đã tồn tại", codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409))
	SEOPageNotPublic                     = _errors.MustSpec(410067, "CRM_SEO_PAGE_NOT_PUBLIC", "seo page chưa được xuất bản công khai", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEOPageRenderIncomplete              = _errors.MustSpec(410068, "CRM_SEO_PAGE_RENDER_INCOMPLETE", "seo page chưa render thành công", codes.Unavailable, _errors.LegacyHTTP200(), _errors.LegacyCode(503))
	SEOPageArchived                      = _errors.MustSpec(410069, "CRM_SEO_PAGE_ARCHIVED", "seo page đã archived", codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(410))
	SEOPageNotFound                      = _errors.MustSpec(410070, "CRM_SEO_PAGE_NOT_FOUND", "seo page không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEOPublicPageNotFound                = _errors.MustSpec(410071, "CRM_SEO_PUBLIC_PAGE_NOT_FOUND", "SEO page không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEORelativeAlreadyExists             = _errors.MustSpec(410072, "CRM_SEO_RELATIVE_ALREADY_EXISTS", "seo relative đã tồn tại", codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409))
	SEORelativeNotFound                  = _errors.MustSpec(410073, "CRM_SEO_RELATIVE_NOT_FOUND", "seo relative không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEORenderUsecaseNotConfigured        = _errors.MustSpec(410074, "CRM_SEO_RENDER_USECASE_NOT_CONFIGURED", "SeoRenderUsecase chưa được cấu hình", codes.Unavailable, _errors.LegacyHTTP200(), _errors.LegacyCode(503))
	SEOSitemapLastedAtInvalid            = _errors.MustSpec(410075, "CRM_SEO_SITEMAP_LASTED_AT_INVALID", "siteMapLastedAt phải theo RFC3339 hoặc YYYY-MM-DD", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOSitemapPriorityInvalid            = _errors.MustSpec(410076, "CRM_SEO_SITEMAP_PRIORITY_INVALID", "sitemapPriority phải nằm trong khoảng 0..1", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOSlugRequired                      = _errors.MustSpec(410077, "CRM_SEO_SLUG_REQUIRED", "slug không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOSlugInvalid                       = _errors.MustSpec(410078, "CRM_SEO_SLUG_INVALID", "slug không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOPlanningSourceNotFound            = _errors.MustSpec(410079, "CRM_SEO_PLANNING_SOURCE_NOT_FOUND", "source đồ án quy hoạch không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEOSourceNotFound                    = _errors.MustSpec(410080, "CRM_SEO_SOURCE_NOT_FOUND", "source không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEOParcelSourceNotFound              = _errors.MustSpec(410081, "CRM_SEO_PARCEL_SOURCE_NOT_FOUND", "source parcel không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEOProjectSourceNotFound             = _errors.MustSpec(410082, "CRM_SEO_PROJECT_SOURCE_NOT_FOUND", "source project không tồn tại hoặc AdminProjectService chưa có API detail theo id", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SEOSourceKeyRequired                 = _errors.MustSpec(410083, "CRM_SEO_SOURCE_KEY_REQUIRED", "sourceKey là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOSourceUpdatedAtInvalid            = _errors.MustSpec(410084, "CRM_SEO_SOURCE_UPDATED_AT_INVALID", "sourceUpdatedAt phải theo RFC3339 hoặc YYYY-MM-DD", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OpportunityStageInvalid              = _errors.MustSpec(410085, "CRM_OPPORTUNITY_STAGE_INVALID", "Stage không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	StageIDRequired                      = _errors.MustSpec(410086, "CRM_STAGE_ID_REQUIRED", "stageId is required", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SupportTicketStatusTransitionInvalid = _errors.MustSpec(410087, "CRM_SUPPORT_TICKET_STATUS_TRANSITION_INVALID", "STATUS_TRANSITION_INVALID", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	TagNotFound                          = _errors.MustSpec(410088, "CRM_TAG_NOT_FOUND", "Tag không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	ContactNameRequired                  = _errors.MustSpec(410089, "CRM_CONTACT_NAME_REQUIRED", "Tên liên hệ là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	TagNameRequired                      = _errors.MustSpec(410090, "CRM_TAG_NAME_REQUIRED", "Tên tag không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	TagPayloadRequired                   = _errors.MustSpec(410091, "CRM_TAG_PAYLOAD_REQUIRED", "Thiếu thông tin tag", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	FollowUpTimeInvalid                  = _errors.MustSpec(410092, "CRM_FOLLOW_UP_TIME_INVALID", "Thời điểm follow-up không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SupportTicketNotFound                = _errors.MustSpec(410093, "CRM_SUPPORT_TICKET_NOT_FOUND", "ticket not found", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SupportTicketTitleRequired           = _errors.MustSpec(410094, "CRM_SUPPORT_TICKET_TITLE_REQUIRED", "Tiêu đề không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOTitleRequired                     = _errors.MustSpec(410095, "CRM_SEO_TITLE_REQUIRED", "title không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOTitleRequiredForGenerate          = _errors.MustSpec(410096, "CRM_SEO_TITLE_REQUIRED_FOR_GENERATE", "title là bắt buộc để generate", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(422))
	SEOTitleRequiredForRender            = _errors.MustSpec(410097, "CRM_SEO_TITLE_REQUIRED_FOR_RENDER", "title là bắt buộc để render", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(422))
	SEOInternalLinkTitleRequired         = _errors.MustSpec(410098, "CRM_SEO_INTERNAL_LINK_TITLE_REQUIRED", "title là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	TQDProviderNotConfigured             = _errors.MustSpec(410099, "CRM_TQD_PROVIDER_NOT_CONFIGURED", "TQD provider chưa được cấu hình", codes.Unavailable, _errors.LegacyHTTP200(), _errors.LegacyCode(503))
	OpportunityStateActionDenied         = _errors.MustSpec(410100, "CRM_OPPORTUNITY_STATE_ACTION_DENIED", "Trạng thái cơ hội không cho phép thao tác", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OpportunityStateInvalid              = _errors.MustSpec(410101, "CRM_OPPORTUNITY_STATE_INVALID", "Trạng thái cơ hội không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OpportunityContactOrCustomerRequired = _errors.MustSpec(410102, "CRM_OPPORTUNITY_CONTACT_OR_CUSTOMER_REQUIRED", "Vui lòng chọn liên hệ hoặc nhập tên khách hàng", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SupportTicketStatusNoteRequired      = _errors.MustSpec(410103, "CRM_SUPPORT_TICKET_STATUS_NOTE_REQUIRED", "Vui lòng nhập ghi chú cập nhật trạng thái", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SupportTicketRoutingReasonRequired   = _errors.MustSpec(410104, "CRM_SUPPORT_TICKET_ROUTING_REASON_REQUIRED", "Vui lòng nhập lý do chuyển tuyến", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ExportSensitiveDataReasonRequired    = _errors.MustSpec(410105, "CRM_EXPORT_SENSITIVE_DATA_REASON_REQUIRED", "Vui lòng nhập lý do xuất dữ liệu vì có thể chứa thông tin nhạy cảm", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RequestCannotBePerformed             = _errors.MustSpec(410106, "CRM_REQUEST_CANNOT_BE_PERFORMED", "Yêu cầu không thể thực hiện", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOSourceUnsupported                 = _errors.MustSpec(410107, "CRM_SEO_SOURCE_UNSUPPORTED", "seo source chưa hỗ trợ resolve", codes.Unimplemented, _errors.LegacyHTTP200(), _errors.LegacyCode(501))
	SEOMetadataFieldJSONInvalid          = _errors.MustSpec(410108, "CRM_SEO_METADATA_FIELD_JSON_INVALID", "field phải là JSON hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	SEOModuleContractMismatch            = _errors.MustSpec(410109, "CRM_SEO_MODULE_CONTRACT_MISMATCH", "SEO module contract mismatch", codes.Unavailable, _errors.LegacyHTTP200(), _errors.LegacyCode(503))
	SEOResolverUnavailable               = _errors.MustSpec(410110, "CRM_SEO_RESOLVER_UNAVAILABLE", "resolver SEO chưa cấu hình client nguồn", codes.Unavailable, _errors.LegacyHTTP200(), _errors.LegacyCode(503))
	SEOResolverNotImplemented            = _errors.MustSpec(410111, "CRM_SEO_RESOLVER_NOT_IMPLEMENTED", "resolver SEO chưa có adapter production", codes.Unimplemented, _errors.LegacyHTTP200(), _errors.LegacyCode(501))
	UserNotFound                         = _errors.MustSpec(
		410112, "CRM_USER_NOT_FOUND", "Người dùng không tồn tại",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	SelfBlockNotAllowed = _errors.MustSpec(
		410113, "CRM_SELF_BLOCK_NOT_ALLOWED", "Không thể chặn chính mình",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	SelfUnblockNotAllowed = _errors.MustSpec(
		410114, "CRM_SELF_UNBLOCK_NOT_ALLOWED", "Không thể bỏ chặn chính mình",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	SelfRelationshipNotAllowed = _errors.MustSpec(
		410115, "CRM_SELF_RELATIONSHIP_NOT_ALLOWED", "Bạn không thể theo dõi/kết bạn với chính mình",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	UnblockRequired = _errors.MustSpec(
		410116, "CRM_UNBLOCK_REQUIRED", "Vui lòng bỏ chặn",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(505),
	)
	DefaultStageMutationDenied = _errors.MustSpec(
		410117, "CRM_DEFAULT_STAGE_MUTATION_DENIED", "Giai đoạn của quy trình mặc định không thể thao tác",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	StageHasCustomers = _errors.MustSpec(
		410118, "CRM_STAGE_HAS_CUSTOMERS", "Giai đoạn này đang có khách hàng",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	FriendRequestInputInvalid = _errors.MustSpec(
		410119, "CRM_FRIEND_REQUEST_INPUT_INVALID", "Thông tin không đúng vui lòng kiểm tra lại",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PipelineAccessDenied = _errors.MustSpec(
		410120, "CRM_PIPELINE_ACCESS_DENIED", "Bạn không có quyền truy cập quy trình",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	DefaultPipelineUpdateDenied = _errors.MustSpec(
		410121, "CRM_DEFAULT_PIPELINE_UPDATE_DENIED", "Bạn không có quyền cập nhật quy trình mặc định",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PipelineHasCustomers = _errors.MustSpec(
		410122, "CRM_PIPELINE_HAS_CUSTOMERS", "Quy trình này đang có khách hàng",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	DefaultPipelineAccessDenied = _errors.MustSpec(
		410123, "CRM_DEFAULT_PIPELINE_ACCESS_DENIED", "Bạn không có quyền truy cập quy trình mặc định",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ContactAccessDenied = _errors.MustSpec(
		410124, "CRM_CONTACT_ACCESS_DENIED", "Bạn không có quyền truy cập",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	CustomerAlreadyExists = _errors.MustSpec(
		410125, "CRM_CUSTOMER_ALREADY_EXISTS", "Khách hàng đã tồn tại",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409),
	)
	SelfFollowNotAllowed = _errors.MustSpec(
		410126, "CRM_SELF_FOLLOW_NOT_ALLOWED", "Không thể tự follow chính mình",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	AlreadyFollowing = _errors.MustSpec(
		410127, "CRM_ALREADY_FOLLOWING", "Đã follow người này rồi",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	NotFollowing = _errors.MustSpec(
		410128, "CRM_NOT_FOLLOWING", "Bạn chưa follow người này",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	StageLimitExceeded = _errors.MustSpec(
		410129, "CRM_STAGE_LIMIT_EXCEEDED", "Số lượng stage vượt quá 20",
		codes.ResourceExhausted, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	FriendGroupDeleteDenied = _errors.MustSpec(
		410130, "CRM_FRIEND_GROUP_DELETE_DENIED", "bạn không có quyền xóa nhóm này",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(401),
	)
	RuleFieldsRequired = _errors.MustSpec(
		410131, "CRM_RULE_FIELDS_REQUIRED", "ownerId, ownerType, ruleName, condition, trigger, conditionValue and triggerValue are required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	IDRequired = _errors.MustSpec(
		410132, "CRM_ID_REQUIRED", "id is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	StageCreateFieldsRequired = _errors.MustSpec(
		410133, "CRM_STAGE_CREATE_FIELDS_REQUIRED", "pipelineId and stageName are required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	PipelineOwnerFieldsRequired = _errors.MustSpec(
		410134, "CRM_PIPELINE_OWNER_FIELDS_REQUIRED", "ownerId and ownerType are required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	LeadStageFieldsRequired = _errors.MustSpec(
		410135, "CRM_LEAD_STAGE_FIELDS_REQUIRED", "leadId and stageID is required",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	ContactAlreadyExists = _errors.MustSpec(
		410136, "CRM_CONTACT_ALREADY_EXISTS", "Contact already exists for this profile",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409),
	)
	FriendshipAlreadyExists = _errors.MustSpec(
		410137, "CRM_FRIENDSHIP_ALREADY_EXISTS", "already friends",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409),
	)
	FriendRequestAlreadySent = _errors.MustSpec(
		410138, "CRM_FRIEND_REQUEST_ALREADY_SENT", "friend request already sent",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(409),
	)
	FriendRequestPendingIncoming = _errors.MustSpec(
		410139, "CRM_FRIEND_REQUEST_PENDING_INCOMING", "you have a pending request from this user, please accept or reject it",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(409),
	)
	FriendRequestAlreadyProcessed = _errors.MustSpec(
		410140, "CRM_FRIEND_REQUEST_ALREADY_PROCESSED", "request already processed",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(409),
	)
	FriendRequestNotPending = _errors.MustSpec(
		410141, "CRM_FRIEND_REQUEST_NOT_PENDING", "cannot cancel non-pending request",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(409),
	)
	FriendRequestReceiverDenied = _errors.MustSpec(
		410142, "CRM_FRIEND_REQUEST_RECEIVER_DENIED", "you are not the receiver of this request",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	FriendRequestNotFound = _errors.MustSpec(
		410143, "CRM_FRIEND_REQUEST_NOT_FOUND", "friend request not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	FriendshipNotFound = _errors.MustSpec(
		410144, "CRM_FRIENDSHIP_NOT_FOUND", "friendship not found",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
)
