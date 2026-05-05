package enums

type ProblemReport uint32

const (
	ProblemReport_FailData   ProblemReport = 10
	ProblemReport_MissData   ProblemReport = 20
	ProblemReport_NeedVerify ProblemReport = 30
	ProblemReport_Other      ProblemReport = 40
)

var EProblemReportNames = map[ProblemReport]string{
	ProblemReport_FailData:   "Sai dữ liệu",
	ProblemReport_MissData:   "Thiếu dữ liệu",
	ProblemReport_NeedVerify: "Cần xác minh",
	ProblemReport_Other:      "Khác",
}
