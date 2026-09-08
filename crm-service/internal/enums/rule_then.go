package enums

type ERuleTrigger int

const (
	RuleThenReminder    ERuleTrigger = 10
	RuleThenLabel       ERuleTrigger = 20
	RuleThenAssignLead  ERuleTrigger = 30
	RuleThenSwitchStage ERuleTrigger = 40
)

var RuleThenMap = map[ERuleTrigger]string{
	RuleThenReminder:    "Nhắc nhở chu kỳ",
	RuleThenLabel:       "Gắn nhãn nguy cơ",
	RuleThenAssignLead:  "Gán nhân viên",
	RuleThenSwitchStage: "Chuyển trạng thái",
}