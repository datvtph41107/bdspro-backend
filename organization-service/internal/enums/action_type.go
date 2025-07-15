package enums

type ActionType uint

const (
	// Deal related actions
	ActionTypeDealCreated        ActionType = 10
	ActionTypeDealUpdated        ActionType = 11
	ActionTypeDealStatusChanged  ActionType = 12
	ActionTypeDealCancelled      ActionType = 13
	
	// Commission related actions
	ActionTypeCommissionUpdated  ActionType = 20
	ActionTypeCommissionAdded    ActionType = 21
	
	// Member related actions
	ActionTypeMemberAdded        ActionType = 30
	ActionTypeMemberRemoved      ActionType = 31
	ActionTypeMemberRoleChanged  ActionType = 32
	
	// Internal note actions
	ActionTypeInternalNoteAdded  ActionType = 40
	ActionTypeInternalNoteEdited ActionType = 41
	
	// Investment related actions
	ActionTypeInvestmentAdded    ActionType = 50
	ActionTypeInvestmentApproved ActionType = 51
	ActionTypeInvestmentRejected ActionType = 52
)

var ActionTypeMap = map[ActionType]string{
	ActionTypeDealCreated:        "Tạo thương vụ",
	ActionTypeDealUpdated:        "Cập nhật thương vụ",
	ActionTypeDealStatusChanged:  "Thay đổi trạng thái thương vụ",
	ActionTypeDealCancelled:      "Hủy thương vụ",
	ActionTypeCommissionUpdated:  "Cập nhật hoa hồng",
	ActionTypeCommissionAdded:    "Thêm hoa hồng",
	ActionTypeMemberAdded:        "Thêm thành viên",
	ActionTypeMemberRemoved:      "Xóa thành viên",
	ActionTypeMemberRoleChanged:  "Thay đổi vai trò thành viên",
	ActionTypeInternalNoteAdded:  "Thêm ghi chú nội bộ",
	ActionTypeInternalNoteEdited: "Sửa ghi chú nội bộ",
	ActionTypeInvestmentAdded:    "Thêm đầu tư",
	ActionTypeInvestmentApproved: "Duyệt đầu tư",
	ActionTypeInvestmentRejected: "Từ chối đầu tư",
} 