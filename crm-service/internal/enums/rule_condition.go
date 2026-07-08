package enums

type ECondition int

const (
	ConditionDateNotHistory ECondition = 10
	ConditionStage          ECondition = 20
	ConditionTag            ECondition = 30
	ConditionBirth          ECondition = 40
	ConditionEventDate      ECondition = 50
)

var ConditionMap = map[ECondition]string{
	ConditionDateNotHistory: "Số ngày không có nhật ký",
	ConditionStage:          "Trạng thái",
	ConditionTag:            "Tag",
	ConditionBirth:          "Ngày sinh",
	ConditionEventDate:      "Ngày sự kiện",
}