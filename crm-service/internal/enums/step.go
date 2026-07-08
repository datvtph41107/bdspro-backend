package enums

type EStep int

const (
	StepNew       EStep = 10
	StepContacted EStep = 20
	StepNegotiate EStep = 30
	StepDeposite  EStep = 40
	StepClose     EStep = 50
)

var StepMap = map[EStep]string{
	StepNew:       "Mới",
	StepContacted: "Đã liên hệ",
	StepNegotiate: "Đàm phán",
	StepDeposite:  "Cọc",
	StepClose:     "Chốt",
}