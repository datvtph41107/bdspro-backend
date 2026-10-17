package enums

// DealHistoryEventType định nghĩa các loại sự kiện trong lịch sử thương vụ
type DealHistoryEventType uint32

const (
	// Deal lifecycle events
	DealHistoryEventCreate       DealHistoryEventType = iota + 1 // 1
	DealHistoryEventUpdate                                       // 2
	DealHistoryEventDelete                                       // 3
	DealHistoryEventStatusChange                                 // 4
	DealHistoryEventCancel                                       // 5

	// Member management events
	DealHistoryEventAddMember // 6

	// Setting events
	DealHistoryEventUpdateSetting          // 7
	DealHistoryEventUpdateAllowManualInput // 8

	// Investment events
	DealHistoryEventAddInvestment     // 9
	DealHistoryEventUpdateInvestment  // 10
	DealHistoryEventRemoveInvestment  // 11
	DealHistoryEventApproveInvestment // 12
	DealHistoryEventRejectInvestment  // 13

	// Invitation events
	DealHistoryEventSendInvitation    // 14
	DealHistoryEventAcceptInvitation  // 15
	DealHistoryEventRejectInvitation  // 16
	DealHistoryEventWithdrawFromDeal  // 17
	DealHistoryEventRemoveFromDeal    // 18
	DealHistoryEventConfirmWithdrawal // 19

	// Deal cost events
	DealHistoryEventAddDealCost    // 20
	DealHistoryEventUpdateDealCost // 21
	DealHistoryEventRemoveDealCost // 22

	// Appointment events (linked deal)
	DealHistoryEventCreateAppointment       // 23
	DealHistoryEventUpdateAppointment       // 24
	DealHistoryEventDeleteAppointment       // 25
	DealHistoryEventUpdateAppointmentStatus // 26

	// Deal product link events
	DealHistoryEventAddDealProduct    // 27
	DealHistoryEventRemoveDealProduct // 28

	DealHistoryEventUpdateMemberRole // 29

	// Internal note events
	DealHistoryEventAddInternalNote    // 30
	DealHistoryEventEditInternalNote   // 31

	DealHistoryEventUpdateMembersCommission // 32
)

// DealHistoryEventNameMap mapping từ event type sang tên hiển thị
var DealHistoryEventNameMap = map[DealHistoryEventType]string{
	DealHistoryEventCreate:                 "Tạo thương vụ",
	DealHistoryEventUpdate:                 "Cập nhật thương vụ",
	DealHistoryEventDelete:                 "Xóa thương vụ",
	DealHistoryEventStatusChange:           "Cập nhật trạng thái thương vụ",
	DealHistoryEventCancel:                 "Hủy thương vụ",
	DealHistoryEventAddMember:              "Thêm thành viên vào thương vụ",
	DealHistoryEventUpdateSetting:          "Cập nhật cài đặt thương vụ",
	DealHistoryEventUpdateAllowManualInput: "Cập nhật cho phép nhập tay",
	DealHistoryEventAddInvestment:          "Thêm khoản góp",
	DealHistoryEventUpdateInvestment:       "Cập nhật khoản góp",
	DealHistoryEventRemoveInvestment:       "Xóa khoản góp",
	DealHistoryEventApproveInvestment:      "Duyệt khoản góp",
	DealHistoryEventRejectInvestment:       "Từ chối khoản góp",
	DealHistoryEventSendInvitation:         "Gửi lời mời tham gia thương vụ",
	DealHistoryEventAcceptInvitation:       "Chấp nhận lời mời tham gia thương vụ",
	DealHistoryEventRejectInvitation:       "Từ chối lời mời tham gia thương vụ",
	DealHistoryEventWithdrawFromDeal:       "Rút khỏi thương vụ",
	DealHistoryEventRemoveFromDeal:         "Gỡ khỏi thương vụ",
	DealHistoryEventConfirmWithdrawal:      "Xác nhận rút khỏi thương vụ",
	DealHistoryEventAddDealCost:            "Thêm chi phí thương vụ",
	DealHistoryEventUpdateDealCost:         "Cập nhật chi phí thương vụ",
	DealHistoryEventRemoveDealCost:         "Xóa chi phí thương vụ",
	DealHistoryEventCreateAppointment:       "Tạo lịch hẹn thương vụ",
	DealHistoryEventUpdateAppointment:       "Cập nhật lịch hẹn thương vụ",
	DealHistoryEventDeleteAppointment:       "Xóa lịch hẹn thương vụ",
	DealHistoryEventUpdateAppointmentStatus: "Cập nhật trạng thái lịch hẹn",
	DealHistoryEventAddDealProduct:          "Thêm sản phẩm vào thương vụ",
	DealHistoryEventRemoveDealProduct:       "Gỡ sản phẩm khỏi thương vụ",
	DealHistoryEventUpdateMemberRole:        "Cập nhật vai trò thành viên",
	DealHistoryEventAddInternalNote:         "Thêm ghi chú nội bộ",
	DealHistoryEventEditInternalNote:            "Sửa ghi chú nội bộ",
	DealHistoryEventUpdateMembersCommission:     "Cập nhật hoa hồng thành viên",
}

// GetEventName trả về tên hiển thị của event type
func (e DealHistoryEventType) GetEventName() string {
	if name, exists := DealHistoryEventNameMap[e]; exists {
		return name
	}
	return "UNKNOWN_EVENT"
}
