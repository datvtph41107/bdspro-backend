package enums

type ERule int

const (
	RuleCycle        ERule = 10
	RuleWarnLabel    ERule = 20
	RuleSwitchStatus ERule = 30
)

var ReminderMap = map[ERule]string{
	RuleCycle:        "Nhắc nhở chu kỳ",
	RuleWarnLabel:    "Gắn nhãn nguy cơ",
	RuleSwitchStatus: "Chuyển trạng thái",
}
