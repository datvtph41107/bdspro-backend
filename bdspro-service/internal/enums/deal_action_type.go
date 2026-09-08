package enums

type DealActionType uint

const (
	// Deal related actions
	ActionTypeDealCreated       DealActionType = 10
	ActionTypeDealUpdated       DealActionType = 11
	ActionTypeDealStatusChanged DealActionType = 12
	ActionTypeDealCancelled     DealActionType = 13

	// Commission related actions
	ActionTypeCommissionUpdated DealActionType = 20
	ActionTypeCommissionAdded   DealActionType = 21

	// Member related actions
	ActionTypeMemberAdded       DealActionType = 30
	ActionTypeMemberRemoved     DealActionType = 31
	ActionTypeMemberRoleChanged DealActionType = 32

	// Internal note actions
	ActionTypeInternalNoteAdded  DealActionType = 40
	ActionTypeInternalNoteEdited DealActionType = 41

	// Investment related actions
	ActionTypeInvestmentAdded    DealActionType = 50
	ActionTypeInvestmentApproved DealActionType = 51
	ActionTypeInvestmentRejected DealActionType = 52
)

var DealActionTypeMap = map[DealActionType]string{
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
