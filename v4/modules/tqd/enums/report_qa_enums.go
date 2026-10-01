package enums

type ReportQaStatus uint32

const (
	ReportQaStatusSubmitted      ReportQaStatus = 10
	ReportQaStatusVerifying      ReportQaStatus = 20
	ReportQaStatusNeedMoreInfo   ReportQaStatus = 30
	ReportQaStatusVerified       ReportQaStatus = 40
	ReportQaStatusWaitingDataFix ReportQaStatus = 50
	ReportQaStatusResolved       ReportQaStatus = 60
	ReportQaStatusClosed         ReportQaStatus = 70
	ReportQaStatusRejected       ReportQaStatus = 80
)

type ReportSeverity uint32

const (
	ReportSeverityLow       ReportSeverity = 10
	ReportSeverityMedium    ReportSeverity = 20
	ReportSeverityHigh      ReportSeverity = 30
	ReportSeverityCritical  ReportSeverity = 40
	ReportSeverityEmergency ReportSeverity = 50
)

var ValidReportQaTransitions = map[ReportQaStatus][]ReportQaStatus{
	ReportQaStatusSubmitted:      {ReportQaStatusVerifying, ReportQaStatusRejected, ReportQaStatusClosed},
	ReportQaStatusVerifying:      {ReportQaStatusNeedMoreInfo, ReportQaStatusVerified, ReportQaStatusRejected, ReportQaStatusWaitingDataFix},
	ReportQaStatusNeedMoreInfo:   {ReportQaStatusVerifying, ReportQaStatusRejected, ReportQaStatusClosed},
	ReportQaStatusVerified:       {ReportQaStatusWaitingDataFix, ReportQaStatusResolved, ReportQaStatusClosed},
	ReportQaStatusWaitingDataFix: {ReportQaStatusResolved, ReportQaStatusVerifying, ReportQaStatusClosed},
	ReportQaStatusResolved:       {ReportQaStatusClosed},
	ReportQaStatusRejected:       {ReportQaStatusClosed},
	ReportQaStatusClosed:         {},
}

func CanTransitionReportQa(from, to ReportQaStatus) bool {
	for _, s := range ValidReportQaTransitions[from] {
		if s == to {
			return true
		}
	}
	return false
}
