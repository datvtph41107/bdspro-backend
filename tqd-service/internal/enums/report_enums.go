package enums

type ReportType uint32

const (
	ReportTypeParcel  ReportType = 10
	ReportTypeRegion  ReportType = 20
	ReportTypeCompare ReportType = 30
)

type ReportProfile uint32

const (
	ReportProfileCompact  ReportProfile = 10
	ReportProfileDetailed ReportProfile = 20
	ReportProfileCompare  ReportProfile = 30
	ReportProfileAnalysis ReportProfile = 40
)

type ReportFormat uint32

const (
	ReportFormatPDF  ReportFormat = 10
	ReportFormatHTML ReportFormat = 20
	ReportFormatCSV  ReportFormat = 30
)

type ReportStatus uint32

const (
	ReportStatusPending    ReportStatus = 10
	ReportStatusProcessing ReportStatus = 20
	ReportStatusCompleted  ReportStatus = 30
	ReportStatusFailed     ReportStatus = 40
)
