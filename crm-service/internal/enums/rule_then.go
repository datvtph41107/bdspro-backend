package enums

type ERule int

const (
	RuleThenReminder    ERule = 10
	RuleThenLabel       ERule = 20
	RuleThenAssignLead  ERule = 30
	RuleThenSwitchStage ERule = 40
)

var RuleThenMap = map[ERule]string{
	RuleThenReminder:    "Nhắc nhở chu kỳ",
	RuleThenLabel:       "Gắn nhãn nguy cơ",
	RuleThenAssignLead:  "Gán nhân viên",
	RuleThenSwitchStage: "Chuyển trạng thái",
}
