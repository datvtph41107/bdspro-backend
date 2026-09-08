package enums

// EIntentCategory - Danh mục intent
type EIntentCategory string

const (
	IntentCategoryProduct     EIntentCategory = "product"
	IntentCategoryPost        EIntentCategory = "post"
	IntentCategoryAsset       EIntentCategory = "asset"
	IntentCategoryProject     EIntentCategory = "project"
	IntentCategoryTransaction EIntentCategory = "transaction"
	IntentCategoryContact     EIntentCategory = "contact"
	IntentCategoryGeneral     EIntentCategory = "general"
)

// EIntentAction - Hành động của intent
type EIntentAction string

const (
	// Product actions
	IntentActionProductDetail  EIntentAction = "product.detail"
	IntentActionProductSearch  EIntentAction = "product.search"
	IntentActionProductCreate  EIntentAction = "product.create"
	IntentActionProductUpdate  EIntentAction = "product.update"
	IntentActionProductDelete  EIntentAction = "product.delete"
	IntentActionProductList    EIntentAction = "product.list"
	IntentActionProductArchive EIntentAction = "product.archive"
	IntentActionProductShare   EIntentAction = "product.share"
	IntentActionProductExport  EIntentAction = "product.export"

	// Post actions
	IntentActionPostDetail  EIntentAction = "post.detail"
	IntentActionPostSearch  EIntentAction = "post.search"
	IntentActionPostCreate  EIntentAction = "post.create"
	IntentActionPostUpdate  EIntentAction = "post.update"
	IntentActionPostDelete  EIntentAction = "post.delete"
	IntentActionPostList    EIntentAction = "post.list"
	IntentActionPostPublish EIntentAction = "post.publish"

	// Asset actions
	IntentActionAssetDetail  EIntentAction = "asset.detail"
	IntentActionAssetSearch  EIntentAction = "asset.search"
	IntentActionAssetCreate  EIntentAction = "asset.create"
	IntentActionAssetUpdate  EIntentAction = "asset.update"
	IntentActionAssetDelete  EIntentAction = "asset.delete"
	IntentActionAssetList    EIntentAction = "asset.list"
	IntentActionAssetArchive EIntentAction = "asset.archive"

	// Project actions
	IntentActionProjectDetail EIntentAction = "project.detail"
	IntentActionProjectSearch EIntentAction = "project.search"
	IntentActionProjectCreate EIntentAction = "project.create"
	IntentActionProjectUpdate EIntentAction = "project.update"
	IntentActionProjectDelete EIntentAction = "project.delete"
	IntentActionProjectList   EIntentAction = "project.list"

	// Transaction actions
	IntentActionTransactionDetail EIntentAction = "transaction.detail"
	IntentActionTransactionCreate EIntentAction = "transaction.create"
	IntentActionTransactionUpdate EIntentAction = "transaction.update"
	IntentActionTransactionList   EIntentAction = "transaction.list"

	// Contact actions
	IntentActionContactDetail EIntentAction = "contact.detail"
	IntentActionContactSearch EIntentAction = "contact.search"
	IntentActionContactCreate EIntentAction = "contact.create"
	IntentActionContactUpdate EIntentAction = "contact.update"
	IntentActionContactDelete EIntentAction = "contact.delete"

	// General actions
	IntentActionHelp      EIntentAction = "general.help"
	IntentActionGreeting  EIntentAction = "general.greeting"
	IntentActionUnknown   EIntentAction = "general.unknown"
	IntentActionFallback  EIntentAction = "general.fallback"
	IntentActionAnalyze   EIntentAction = "general.analyze"
	IntentActionSummarize EIntentAction = "general.summarize"
	IntentActionTranslate EIntentAction = "general.translate"
	IntentActionCalculate EIntentAction = "general.calculate"
	IntentActionRecommend EIntentAction = "general.recommend"
)

// IntentMetadata - Metadata cho intent
type IntentMetadata struct {
	Category    EIntentCategory
	Action      EIntentAction
	Description string
	Examples    []string
	Keywords    []string
}

// IntentRegistry - Registry chứa tất cả intents
var IntentRegistry = map[EIntentAction]IntentMetadata{
	// Product intents
	IntentActionProductDetail: {
		Category:    IntentCategoryProduct,
		Action:      IntentActionProductDetail,
		Description: "Xem chi tiết sản phẩm",
		Examples:    []string{"cho tôi xem sản phẩm SP001", "chi tiết sản phẩm này", "thông tin sản phẩm"},
		Keywords:    []string{"xem", "chi tiết", "thông tin", "sản phẩm"},
	},
	IntentActionProductSearch: {
		Category:    IntentCategoryProduct,
		Action:      IntentActionProductSearch,
		Description: "Tìm kiếm sản phẩm",
		Examples:    []string{"tìm nhà ở Hà Nội", "sản phẩm giá 5 tỷ", "nhà 3 tầng"},
		Keywords:    []string{"tìm", "tìm kiếm", "search", "nhà", "sản phẩm"},
	},
	IntentActionProductCreate: {
		Category:    IntentCategoryProduct,
		Action:      IntentActionProductCreate,
		Description: "Tạo sản phẩm mới",
		Examples:    []string{"tạo sản phẩm", "thêm sản phẩm mới", "đăng bán nhà"},
		Keywords:    []string{"tạo", "thêm", "đăng", "mới"},
	},
	IntentActionProductUpdate: {
		Category:    IntentCategoryProduct,
		Action:      IntentActionProductUpdate,
		Description: "Cập nhật sản phẩm",
		Examples:    []string{"cập nhật giá sản phẩm", "sửa thông tin", "chỉnh sửa"},
		Keywords:    []string{"cập nhật", "sửa", "chỉnh sửa", "thay đổi"},
	},
	IntentActionProductDelete: {
		Category:    IntentCategoryProduct,
		Action:      IntentActionProductDelete,
		Description: "Xóa sản phẩm",
		Examples:    []string{"xóa sản phẩm", "gỡ sản phẩm"},
		Keywords:    []string{"xóa", "gỡ", "delete"},
	},

	// Post intents
	IntentActionPostCreate: {
		Category:    IntentCategoryPost,
		Action:      IntentActionPostCreate,
		Description: "Tạo bài đăng",
		Examples:    []string{"đăng tin", "tạo bài viết", "post bài"},
		Keywords:    []string{"đăng", "post", "tạo bài"},
	},
	IntentActionPostSearch: {
		Category:    IntentCategoryPost,
		Action:      IntentActionPostSearch,
		Description: "Tìm kiếm bài đăng",
		Examples:    []string{"tìm bài đăng", "bài viết về nhà"},
		Keywords:    []string{"tìm", "bài đăng", "bài viết"},
	},

	// Asset intents
	IntentActionAssetCreate: {
		Category:    IntentCategoryAsset,
		Action:      IntentActionAssetCreate,
		Description: "Tạo tài sản",
		Examples:    []string{"tạo tài sản", "thêm tài sản mới"},
		Keywords:    []string{"tài sản", "asset", "tạo"},
	},

	// General intents
	IntentActionHelp: {
		Category:    IntentCategoryGeneral,
		Action:      IntentActionHelp,
		Description: "Yêu cầu trợ giúp",
		Examples:    []string{"help", "giúp tôi", "hướng dẫn"},
		Keywords:    []string{"help", "giúp", "hướng dẫn", "trợ giúp"},
	},
	IntentActionGreeting: {
		Category:    IntentCategoryGeneral,
		Action:      IntentActionGreeting,
		Description: "Chào hỏi",
		Examples:    []string{"xin chào", "hello", "chào bạn"},
		Keywords:    []string{"xin chào", "hello", "hi", "chào"},
	},
	IntentActionAnalyze: {
		Category:    IntentCategoryGeneral,
		Action:      IntentActionAnalyze,
		Description: "Phân tích dữ liệu",
		Examples:    []string{"phân tích văn bản này", "analyze"},
		Keywords:    []string{"phân tích", "analyze"},
	},
}

// GetIntentByAction - Lấy intent metadata theo action
func GetIntentByAction(action EIntentAction) (IntentMetadata, bool) {
	metadata, exists := IntentRegistry[action]
	return metadata, exists
}

// GetIntentsByCategory - Lấy danh sách intents theo category
func GetIntentsByCategory(category EIntentCategory) []IntentMetadata {
	var intents []IntentMetadata
	for _, metadata := range IntentRegistry {
		if metadata.Category == category {
			intents = append(intents, metadata)
		}
	}
	return intents
}

// IsValidIntent - Kiểm tra intent có hợp lệ không
func IsValidIntent(action EIntentAction) bool {
	_, exists := IntentRegistry[action]
	return exists
}
