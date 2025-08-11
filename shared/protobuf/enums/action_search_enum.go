package shared_enum

type EActionSearch int32

const (
	ActionSearchCreate           EActionSearch = 10
	ActionSearchUpdate           EActionSearch = 20
	ActionSearchDelete           EActionSearch = 30
	ActionSearchDecentralization EActionSearch = 40
	ActionSearchRuleChange       EActionSearch = 50
)

var ActionSearchMap = map[EActionSearch]string{
	ActionSearchCreate:           "Tạo",
	ActionSearchUpdate:           "Cập nhật",
	ActionSearchDelete:           "Xóa",
	ActionSearchDecentralization: "Phân quyền",
	ActionSearchRuleChange:       "Thay đổi quy tắc",
}
